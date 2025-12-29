package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"multi-onboarding/domain/Property/insurance"
	"multi-onboarding/utils"
)

type repository struct {
	mysqlSess *sql.DB
	draftData []insurance.InsuranceDraft
	draftMu   sync.RWMutex
	nextID    int
}

// NewInsuranceRepository creates a new insurance repository
func NewInsuranceRepository() insurance.Repository {
	return &repository{
		mysqlSess: utils.GetDB(),
		draftData: make([]insurance.InsuranceDraft, 0),
		nextID:    1,
	}
}

// Draft operations (in-memory)

func (r *repository) AddDraft(draft *insurance.InsuranceDraft) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	draft.ID = r.nextID
	r.nextID++
	draft.Timestamp = time.Now()
	r.draftData = append(r.draftData, *draft)
	return nil
}

func (r *repository) GetAllDrafts() ([]insurance.InsuranceDraft, error) {
	r.draftMu.RLock()
	defer r.draftMu.RUnlock()

	result := make([]insurance.InsuranceDraft, len(r.draftData))
	copy(result, r.draftData)
	return result, nil
}

func (r *repository) GetDraftByID(id int) (*insurance.InsuranceDraft, error) {
	r.draftMu.RLock()
	defer r.draftMu.RUnlock()

	for i := range r.draftData {
		if r.draftData[i].ID == id {
			result := r.draftData[i]
			return &result, nil
		}
	}

	return nil, fmt.Errorf("draft not found")
}

func (r *repository) UpdateDraft(id int, draft *insurance.InsuranceDraft) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	for i := range r.draftData {
		if r.draftData[i].ID == id {
			draft.ID = id
			draft.Timestamp = time.Now()
			r.draftData[i] = *draft
			return nil
		}
	}

	return fmt.Errorf("draft not found")
}

func (r *repository) DeleteDraft(id int) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	for i := range r.draftData {
		if r.draftData[i].ID == id {
			r.draftData = append(r.draftData[:i], r.draftData[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("draft not found")
}

func (r *repository) ClearDrafts() error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	r.draftData = make([]insurance.InsuranceDraft, 0)
	r.nextID = 1
	return nil
}

func (r *repository) ConfirmDrafts() error {
	// Get all drafts
	drafts, err := r.GetAllDrafts()
	if err != nil {
		return fmt.Errorf("failed to get drafts: %w", err)
	}

	// Insert each draft into property_service_development.insurances table
	for _, draft := range drafts {
		ins := &insurance.Insurance{
			Code:   draft.Code,
			Name:   draft.Name,
			Status: draft.Status,
			Logo:   draft.Logo,
		}

		if err := r.Create(ins); err != nil {
			return fmt.Errorf("failed to confirm draft %d: %w", draft.ID, err)
		}
	}

	// Clear drafts after confirmation
	if err := r.ClearDrafts(); err != nil {
		return fmt.Errorf("failed to clear drafts: %w", err)
	}

	return nil
}

// Database operations - using property_service_development.insurances

func (r *repository) Create(ins *insurance.Insurance) error {
	// Insert into property_service_development.insurances table
	// Use is_active column (table uses is_active, not status)
	// Include created_by with default value 1
	// Include country with default value 'ID'
	query := `
		INSERT INTO property_service_development.insurances (code, name, is_active, logo, country, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, 'ID', 1, NOW(), NOW())
	`

	result, err := r.mysqlSess.Exec(query, ins.Code, ins.Name, ins.Status, ins.Logo)
	if err != nil {
		return fmt.Errorf("failed to create insurance: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	ins.ID = int(id)
	ins.CreatedAt = time.Now()
	ins.UpdatedAt = time.Now()
	return nil
}

func (r *repository) GetByID(id int) (*insurance.Insurance, error) {
	// Query from property_service_development.insurances
	// Use is_active column (table uses is_active, not status)
	query := `
		SELECT id, code, name, is_active, logo, created_at, updated_at
		FROM property_service_development.insurances
		WHERE id = ?
	`

	ins := &insurance.Insurance{}
	var statusValue int
	err := r.mysqlSess.QueryRow(query, id).Scan(
		&ins.ID,
		&ins.Code,
		&ins.Name,
		&statusValue,
		&ins.Logo,
		&ins.CreatedAt,
		&ins.UpdatedAt,
	)
	
	ins.Status = statusValue

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("insurance not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get insurance: %w", err)
	}

	return ins, nil
}

func (r *repository) GetAll() ([]insurance.Insurance, error) {
	// Query from property_service_development.insurances
	// Use is_active column (table uses is_active, not status)
	query := `
		SELECT id, code, name, is_active, logo, created_at, updated_at
		FROM property_service_development.insurances
		ORDER BY created_at DESC
	`

	rows, err := r.mysqlSess.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query insurances: %w", err)
	}
	defer rows.Close()

	insurances := []insurance.Insurance{}
	count := 0
	for rows.Next() {
		var ins insurance.Insurance
		var statusValue int
		err := rows.Scan(
			&ins.ID,
			&ins.Code,
			&ins.Name,
			&statusValue,
			&ins.Logo,
			&ins.CreatedAt,
			&ins.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan insurance: %w", err)
		}
		ins.Status = statusValue
		insurances = append(insurances, ins)
		count++
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Always return an array, even if empty
	if insurances == nil {
		insurances = []insurance.Insurance{}
	}

	// Log for debugging
	if count == 0 {
		fmt.Printf("⚠️  GetAll Property insurances: No data found in database\n")
	} else {
		fmt.Printf("✅ GetAll Property insurances: Found %d insurances\n", count)
	}

	return insurances, nil
}

func (r *repository) Update(id int, ins *insurance.Insurance) error {
	// Get data before update
	beforeData, err := r.GetInsuranceBeforeUpdate(id)
	if err != nil {
		beforeData = &insurance.InsuranceBeforeUpdate{
			Code:   ins.Code,
			Name:   ins.Name,
			Status: ins.Status,
			Logo:   ins.Logo,
		}
	}

	// Build changed fields
	changedBefore := make(map[string]interface{})
	changedAfter := make(map[string]interface{})
	updateFields := []string{}
	updateValues := []interface{}{}

	if beforeData.Code != ins.Code {
		changedBefore["code"] = beforeData.Code
		changedAfter["code"] = ins.Code
		updateFields = append(updateFields, "code = ?")
		updateValues = append(updateValues, ins.Code)
	}
	if beforeData.Name != ins.Name {
		changedBefore["name"] = beforeData.Name
		changedAfter["name"] = ins.Name
		updateFields = append(updateFields, "name = ?")
		updateValues = append(updateValues, ins.Name)
	}
	if beforeData.Status != ins.Status {
		changedBefore["status"] = beforeData.Status
		changedAfter["status"] = ins.Status
		// Use is_active column (table uses is_active, not status)
		updateFields = append(updateFields, "is_active = ?")
		updateValues = append(updateValues, ins.Status)
	}
	if beforeData.Logo != ins.Logo {
		changedBefore["logo"] = beforeData.Logo
		changedAfter["logo"] = ins.Logo
		updateFields = append(updateFields, "logo = ?")
		updateValues = append(updateValues, ins.Logo)
	}

	if len(updateFields) == 0 {
		return nil // No changes
	}

	updateFields = append(updateFields, "updated_at = NOW()")
	updateValues = append(updateValues, id)

	// Build update query (using is_active column) - update property_service_development.insurances
	queryStr := strings.Join(updateFields, ", ")
	query := fmt.Sprintf(`
		UPDATE property_service_development.insurances
		SET %s
		WHERE id = ?
	`, queryStr)

	result, err := r.mysqlSess.Exec(query, updateValues...)
	if err != nil {
		return fmt.Errorf("failed to update insurance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("insurance not found")
	}

	ins.ID = id

	// Save history if there are changes
	if len(changedBefore) > 0 {
		beforeJSON, _ := json.Marshal(changedBefore)
		afterJSON, _ := json.Marshal(changedAfter)
		_ = r.InsertHistory("system", "property_insurances", "UPDATE", id, "insurance", string(beforeJSON), string(afterJSON))
	}

	return nil
}

func (r *repository) Delete(id int) error {
	query := `DELETE FROM property_service_development.insurances WHERE id = ?`

	result, err := r.mysqlSess.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete insurance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("insurance not found")
	}

	return nil
}

func (r *repository) GetInsuranceBeforeUpdate(id int) (*insurance.InsuranceBeforeUpdate, error) {
	// Use is_active column (table uses is_active, not status)
	query := `SELECT code, name, is_active, logo FROM property_service_development.insurances WHERE id = ?`

	var before insurance.InsuranceBeforeUpdate
	var statusValue int
	err := r.mysqlSess.QueryRow(query, id).Scan(
		&before.Code,
		&before.Name,
		&statusValue,
		&before.Logo,
	)
	
	before.Status = statusValue
	if err != nil {
		return nil, err
	}

	return &before, nil
}

func (r *repository) InsertHistory(userName, section, action string, recordID int, recordType, dataBefore, dataAfter string) error {
	// Insert history to property_service_development.histories (or shared histories table)
	// Assuming histories table exists in property_service_development or use travel_service_development.histories
	historyQuery := `INSERT INTO property_service_development.histories (user_name, section, action, record_id, record_type, data_before, data_after, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	_, err := r.mysqlSess.Exec(historyQuery, userName, section, action, recordID, recordType, dataBefore, dataAfter)
	if err != nil {
		// If property_service_development.histories doesn't exist, try travel_service_development.histories
		historyQuery = `INSERT INTO travel_service_development.histories (user_name, section, action, record_id, record_type, data_before, data_after, created_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`
		_, err = r.mysqlSess.Exec(historyQuery, userName, section, action, recordID, recordType, dataBefore, dataAfter)
		if err != nil {
			return fmt.Errorf("failed to insert history: %w", err)
		}
	}

	return nil
}

