package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"multi-onboarding/model/Travel"
)

// InsuranceRepository defines the interface for insurance operations
type InsuranceRepository interface {
	// Draft operations (in-memory)
	AddDraft(insurance *model.InsuranceDraft) error
	GetAllDrafts() ([]model.InsuranceDraft, error)
	GetDraftByID(id int) (*model.InsuranceDraft, error)
	UpdateDraft(id int, insurance *model.InsuranceDraft) error
	DeleteDraft(id int) error
	ClearDrafts() error
	ConfirmDrafts() error

	// Database operations
	Create(insurance *model.Insurance) error
	GetByID(id int) (*model.Insurance, error)
	GetAll() ([]model.Insurance, error)
	Update(id int, insurance *model.Insurance) error
	Delete(id int) error
}

// MySQLInsuranceRepository implements InsuranceRepository with MySQL database storage
type MySQLInsuranceRepository struct {
	db         *sql.DB
	historyRepo HistoryRepository
	draftData []model.InsuranceDraft
	draftMu   sync.RWMutex
	nextID    int
}

// NewMySQLInsuranceRepository creates a new MySQL insurance repository
func NewMySQLInsuranceRepository(db *sql.DB) *MySQLInsuranceRepository {
	return &MySQLInsuranceRepository{
		db:        db,
		draftData: make([]model.InsuranceDraft, 0),
		nextID:    1,
	}
}

// SetHistoryRepository sets the history repository for saving update histories
func (r *MySQLInsuranceRepository) SetHistoryRepository(historyRepo HistoryRepository) {
	r.historyRepo = historyRepo
}

// Draft operations (in-memory)

func (r *MySQLInsuranceRepository) AddDraft(insurance *model.InsuranceDraft) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	insurance.ID = r.nextID
	r.nextID++
	insurance.Timestamp = time.Now()
	r.draftData = append(r.draftData, *insurance)
	return nil
}

func (r *MySQLInsuranceRepository) GetAllDrafts() ([]model.InsuranceDraft, error) {
	r.draftMu.RLock()
	defer r.draftMu.RUnlock()

	result := make([]model.InsuranceDraft, len(r.draftData))
	copy(result, r.draftData)
	return result, nil
}

func (r *MySQLInsuranceRepository) GetDraftByID(id int) (*model.InsuranceDraft, error) {
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

func (r *MySQLInsuranceRepository) UpdateDraft(id int, insurance *model.InsuranceDraft) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	for i := range r.draftData {
		if r.draftData[i].ID == id {
			insurance.ID = id
			insurance.Timestamp = time.Now()
			r.draftData[i] = *insurance
			return nil
		}
	}

	return fmt.Errorf("draft not found")
}

func (r *MySQLInsuranceRepository) DeleteDraft(id int) error {
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

func (r *MySQLInsuranceRepository) ClearDrafts() error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	r.draftData = make([]model.InsuranceDraft, 0)
	r.nextID = 1
	return nil
}

func (r *MySQLInsuranceRepository) ConfirmDrafts() error {
	// Get all drafts
	drafts, err := r.GetAllDrafts()
	if err != nil {
		return fmt.Errorf("failed to get drafts: %w", err)
	}

	// Insert each draft into insurances table
	for _, draft := range drafts {
		insurance := &model.Insurance{
			Code:   draft.Code,
			Name:   draft.Name,
			Status: draft.Status,
			Logo:   draft.Logo,
		}

		if err := r.Create(insurance); err != nil {
			return fmt.Errorf("failed to confirm draft %d: %w", draft.ID, err)
		}
	}

	// Clear drafts after confirmation
	if err := r.ClearDrafts(); err != nil {
		return fmt.Errorf("failed to clear drafts: %w", err)
	}

	return nil
}

// Database operations

func (r *MySQLInsuranceRepository) Create(insurance *model.Insurance) error {
	query := `
		INSERT INTO insurances (code, name, is_active, logo, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, 1, NOW(), NOW())
	`

	result, err := r.db.Exec(query, insurance.Code, insurance.Name, insurance.Status, insurance.Logo)
	if err != nil {
		return fmt.Errorf("failed to create insurance: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	insurance.ID = int(id)
	return nil
}

func (r *MySQLInsuranceRepository) GetByID(id int) (*model.Insurance, error) {
	query := `
		SELECT id, code, name, is_active, logo, created_at, updated_at
		FROM insurances
		WHERE id = ?
	`

	insurance := &model.Insurance{}
	var isActive int
	err := r.db.QueryRow(query, id).Scan(
		&insurance.ID,
		&insurance.Code,
		&insurance.Name,
		&isActive,
		&insurance.Logo,
		&insurance.CreatedAt,
		&insurance.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("insurance not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get insurance: %w", err)
	}

	insurance.Status = isActive
	return insurance, nil
}

func (r *MySQLInsuranceRepository) GetAll() ([]model.Insurance, error) {
	query := `
		SELECT id, code, name, is_active, logo, created_at, updated_at
		FROM insurances
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query insurances: %w", err)
	}
	defer rows.Close()

	insurances := []model.Insurance{}
	for rows.Next() {
		var insurance model.Insurance
		var isActive int
		err := rows.Scan(
			&insurance.ID,
			&insurance.Code,
			&insurance.Name,
			&isActive,
			&insurance.Logo,
			&insurance.CreatedAt,
			&insurance.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan insurance: %w", err)
		}
		insurance.Status = isActive
		insurances = append(insurances, insurance)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return insurances, nil
}

func (r *MySQLInsuranceRepository) Update(id int, insurance *model.Insurance) error {
	// Get data before update for comparison and history
	beforeInsurance, err := r.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get insurance before update: %w", err)
	}

	// Build changed fields map (only fields that are different)
	changedBefore := make(map[string]interface{})
	changedAfter := make(map[string]interface{})
	updateFields := []string{}
	updateValues := []interface{}{}

	// Compare each field and build update query dynamically
	if beforeInsurance.Code != insurance.Code {
		changedBefore["code"] = beforeInsurance.Code
		changedAfter["code"] = insurance.Code
		updateFields = append(updateFields, "code = ?")
		updateValues = append(updateValues, insurance.Code)
	}
	if beforeInsurance.Name != insurance.Name {
		changedBefore["name"] = beforeInsurance.Name
		changedAfter["name"] = insurance.Name
		updateFields = append(updateFields, "name = ?")
		updateValues = append(updateValues, insurance.Name)
	}
	if beforeInsurance.Status != insurance.Status {
		changedBefore["is_active"] = beforeInsurance.Status
		changedAfter["is_active"] = insurance.Status
		updateFields = append(updateFields, "is_active = ?")
		updateValues = append(updateValues, insurance.Status)
	}
	if beforeInsurance.Logo != insurance.Logo {
		changedBefore["logo"] = beforeInsurance.Logo
		changedAfter["logo"] = insurance.Logo
		updateFields = append(updateFields, "logo = ?")
		updateValues = append(updateValues, insurance.Logo)
	}

	// If no fields changed, return early (no update needed)
	if len(updateFields) == 0 {
		return nil
	}

	// Always update updated_at
	updateFields = append(updateFields, "updated_at = NOW()")
	updateValues = append(updateValues, id)

	// Build dynamic UPDATE query
	query := fmt.Sprintf(`
		UPDATE travel_service_development.insurances
		SET %s
		WHERE id = ?
	`, strings.Join(updateFields, ", "))

	result, err := r.db.Exec(query, updateValues...)
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

	insurance.ID = id

	// Save history only if there are changes and historyRepo is available
	if r.historyRepo != nil && len(changedBefore) > 0 {
		beforeJSON, err := json.Marshal(changedBefore)
		if err != nil {
			fmt.Printf("Warning: failed to marshal changed before data: %v\n", err)
		} else {
			afterJSON, err := json.Marshal(changedAfter)
			if err != nil {
				fmt.Printf("Warning: failed to marshal changed after data: %v\n", err)
			} else {
				history := &model.History{
					UserName:   "System",
					Section:    "insurances",
					Action:     "UPDATE",
					RecordID:   id,
					RecordType: "insurance",
					DataBefore: string(beforeJSON),
					DataAfter:  string(afterJSON),
				}
				if err := r.historyRepo.Create(history); err != nil {
					// Log error but don't fail the update
					fmt.Printf("Warning: failed to save history: %v\n", err)
				} else {
					fmt.Printf("DEBUG: History saved successfully for insurance ID %d with %d changed field(s)\n", id, len(changedBefore))
				}
			}
		}
	}

	return nil
}

func (r *MySQLInsuranceRepository) Delete(id int) error {
	query := `DELETE FROM insurances WHERE id = ?`

	result, err := r.db.Exec(query, id)
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

