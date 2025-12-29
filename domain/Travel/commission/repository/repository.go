package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"multi-onboarding/domain/Travel/commission"
	"multi-onboarding/utils"
)

type repository struct {
	mysqlSess *sql.DB
	draftData []commission.CommissionDraft
	draftMu   sync.RWMutex
	nextID    int
}

// NewCommissionRepository creates a new commission repository
func NewCommissionRepository() commission.Repository {
	return &repository{
		mysqlSess: utils.GetDB(),
		draftData: make([]commission.CommissionDraft, 0),
		nextID:    1,
	}
}

// Draft operations (in-memory)

func (r *repository) AddDraft(draft *commission.CommissionDraft) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	draft.ID = r.nextID
	r.nextID++
	draft.CreatedAt = time.Now()
	r.draftData = append(r.draftData, *draft)
	return nil
}

func (r *repository) GetAllDrafts() ([]commission.CommissionDraft, error) {
	r.draftMu.RLock()
	defer r.draftMu.RUnlock()

	result := make([]commission.CommissionDraft, len(r.draftData))
	copy(result, r.draftData)
	return result, nil
}

func (r *repository) GetDraftByID(id int) (*commission.CommissionDraft, error) {
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

func (r *repository) UpdateDraft(id int, draft *commission.CommissionDraft) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	for i := range r.draftData {
		if r.draftData[i].ID == id {
			draft.ID = id
			draft.CreatedAt = time.Now()
			r.draftData[i] = *draft
			return nil
		}
	}

	return fmt.Errorf("draft not found")
}

func (r *repository) DeleteDraft(id int) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	for i, draft := range r.draftData {
		if draft.ID == id {
			r.draftData = append(r.draftData[:i], r.draftData[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("draft not found")
}

func (r *repository) ClearDrafts() error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	r.draftData = make([]commission.CommissionDraft, 0)
	r.nextID = 1
	return nil
}

func (r *repository) ConfirmDrafts(createdBy int64) error {
	r.draftMu.Lock()
	drafts := make([]commission.CommissionDraft, len(r.draftData))
	copy(drafts, r.draftData)
	r.draftMu.Unlock()

	for _, draft := range drafts {
		c := &commission.Commission{
			ProductCode:          draft.ProductCode,
			CommissionPercentage: draft.CommissionPercentage,
			CommissionVATType:    draft.CommissionVATType,
			AFPercentage:         draft.AFPercentage,
			AFVATType:            draft.AFVATType,
			AdminFee:             draft.AdminFee,
			CreatedBy:            createdBy,
		}

		if err := r.Create(c); err != nil {
			return fmt.Errorf("failed to create commission from draft: %w", err)
		}
	}

	// Clear drafts after confirmation
	if err := r.ClearDrafts(); err != nil {
		return fmt.Errorf("failed to clear drafts: %w", err)
	}

	return nil
}

// Database operations

func (r *repository) GetAll() ([]commission.Commission, error) {
	query := `
		SELECT id, product_code, commission_percentage, commission_vat_type, 
		       af_percentage, af_vat_type, admin_fee, created_by, created_at, updated_by, updated_at
		FROM commissions
		ORDER BY created_at DESC
	`

	rows, err := r.mysqlSess.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commissions []commission.Commission
	for rows.Next() {
		var c commission.Commission
		var updatedBy sql.NullInt64
		var updatedAt sql.NullTime

		err := rows.Scan(
			&c.ID,
			&c.ProductCode,
			&c.CommissionPercentage,
			&c.CommissionVATType,
			&c.AFPercentage,
			&c.AFVATType,
			&c.AdminFee,
			&c.CreatedBy,
			&c.CreatedAt,
			&updatedBy,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if updatedBy.Valid {
			c.UpdatedBy = &updatedBy.Int64
		}
		if updatedAt.Valid {
			c.UpdatedAt = &updatedAt.Time
		}

		commissions = append(commissions, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return commissions, nil
}

func (r *repository) GetByID(id int64) (*commission.Commission, error) {
	query := `
		SELECT id, product_code, commission_percentage, commission_vat_type, 
		       af_percentage, af_vat_type, admin_fee, created_by, created_at, updated_by, updated_at
		FROM commissions
		WHERE id = ?
	`

	var c commission.Commission
	var updatedBy sql.NullInt64
	var updatedAt sql.NullTime

	err := r.mysqlSess.QueryRow(query, id).Scan(
		&c.ID,
		&c.ProductCode,
		&c.CommissionPercentage,
		&c.CommissionVATType,
		&c.AFPercentage,
		&c.AFVATType,
		&c.AdminFee,
		&c.CreatedBy,
		&c.CreatedAt,
		&updatedBy,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("commission not found")
	}
	if err != nil {
		return nil, err
	}

	if updatedBy.Valid {
		c.UpdatedBy = &updatedBy.Int64
	}
	if updatedAt.Valid {
		c.UpdatedAt = &updatedAt.Time
	}

	return &c, nil
}

func (r *repository) GetByProductCode(productCode string) (*commission.Commission, error) {
	query := `
		SELECT id, product_code, commission_percentage, commission_vat_type, 
		       af_percentage, af_vat_type, admin_fee, created_by, created_at, updated_by, updated_at
		FROM commissions
		WHERE product_code = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	var c commission.Commission
	var updatedBy sql.NullInt64
	var updatedAt sql.NullTime

	err := r.mysqlSess.QueryRow(query, productCode).Scan(
		&c.ID,
		&c.ProductCode,
		&c.CommissionPercentage,
		&c.CommissionVATType,
		&c.AFPercentage,
		&c.AFVATType,
		&c.AdminFee,
		&c.CreatedBy,
		&c.CreatedAt,
		&updatedBy,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("commission not found")
	}
	if err != nil {
		return nil, err
	}

	if updatedBy.Valid {
		c.UpdatedBy = &updatedBy.Int64
	}
	if updatedAt.Valid {
		c.UpdatedAt = &updatedAt.Time
	}

	return &c, nil
}

func (r *repository) Create(c *commission.Commission) error {
	query := `
		INSERT INTO commissions 
		(product_code, commission_percentage, commission_vat_type, af_percentage, af_vat_type, admin_fee, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.mysqlSess.Exec(query,
		c.ProductCode,
		c.CommissionPercentage,
		c.CommissionVATType,
		c.AFPercentage,
		c.AFVATType,
		c.AdminFee,
		c.CreatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create commission: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	c.ID = id
	c.CreatedAt = time.Now()
	now := time.Now()
	c.UpdatedAt = &now

	return nil
}

func (r *repository) Update(id int64, c *commission.Commission) error {
	// Get data before update
	beforeData, err := r.GetCommissionBeforeUpdate(id)
	if err != nil {
		beforeData = &commission.CommissionBeforeUpdate{
			ProductCode:          c.ProductCode,
			CommissionPercentage: c.CommissionPercentage,
			CommissionVATType:    c.CommissionVATType,
			AFPercentage:         c.AFPercentage,
			AFVATType:            c.AFVATType,
			AdminFee:             c.AdminFee,
		}
	}

	// Build changed fields
	changedBefore := make(map[string]interface{})
	changedAfter := make(map[string]interface{})
	updateFields := []string{}
	updateValues := []interface{}{}

	if beforeData.ProductCode != c.ProductCode {
		changedBefore["product_code"] = beforeData.ProductCode
		changedAfter["product_code"] = c.ProductCode
		updateFields = append(updateFields, "product_code = ?")
		updateValues = append(updateValues, c.ProductCode)
	}
	if beforeData.CommissionPercentage != c.CommissionPercentage {
		changedBefore["commission_percentage"] = beforeData.CommissionPercentage
		changedAfter["commission_percentage"] = c.CommissionPercentage
		updateFields = append(updateFields, "commission_percentage = ?")
		updateValues = append(updateValues, c.CommissionPercentage)
	}
	if beforeData.CommissionVATType != c.CommissionVATType {
		changedBefore["commission_vat_type"] = beforeData.CommissionVATType
		changedAfter["commission_vat_type"] = c.CommissionVATType
		updateFields = append(updateFields, "commission_vat_type = ?")
		updateValues = append(updateValues, c.CommissionVATType)
	}
	if beforeData.AFPercentage != c.AFPercentage {
		changedBefore["af_percentage"] = beforeData.AFPercentage
		changedAfter["af_percentage"] = c.AFPercentage
		updateFields = append(updateFields, "af_percentage = ?")
		updateValues = append(updateValues, c.AFPercentage)
	}
	if beforeData.AFVATType != c.AFVATType {
		changedBefore["af_vat_type"] = beforeData.AFVATType
		changedAfter["af_vat_type"] = c.AFVATType
		updateFields = append(updateFields, "af_vat_type = ?")
		updateValues = append(updateValues, c.AFVATType)
	}
	if beforeData.AdminFee != c.AdminFee {
		changedBefore["admin_fee"] = beforeData.AdminFee
		changedAfter["admin_fee"] = c.AdminFee
		updateFields = append(updateFields, "admin_fee = ?")
		updateValues = append(updateValues, c.AdminFee)
	}

	if len(updateFields) == 0 {
		return nil // No changes
	}

	var updatedBy interface{}
	if c.UpdatedBy != nil {
		updatedBy = *c.UpdatedBy
		updateFields = append(updateFields, "updated_by = ?")
		updateValues = append(updateValues, updatedBy)
	}

	updateFields = append(updateFields, "updated_at = NOW()")
	updateValues = append(updateValues, id)

	query := fmt.Sprintf(`
		UPDATE commissions
		SET %s
		WHERE id = ?
	`, strings.Join(updateFields, ", "))

	result, err := r.mysqlSess.Exec(query, updateValues...)
	if err != nil {
		return fmt.Errorf("failed to update commission: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("commission not found")
	}

	c.ID = id
	now := time.Now()
	c.UpdatedAt = &now

	// Save history if there are changes
	if len(changedBefore) > 0 {
		beforeJSON, _ := json.Marshal(changedBefore)
		afterJSON, _ := json.Marshal(changedAfter)
		_ = r.InsertHistory("system", "commissions", "UPDATE", id, "commission", string(beforeJSON), string(afterJSON))
	}

	return nil
}

func (r *repository) Delete(id int64) error {
	query := `DELETE FROM commissions WHERE id = ?`

	result, err := r.mysqlSess.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("commission not found")
	}

	return nil
}

func (r *repository) GetCommissionBeforeUpdate(id int64) (*commission.CommissionBeforeUpdate, error) {
	query := `
		SELECT product_code, commission_percentage, commission_vat_type, 
		       af_percentage, af_vat_type, admin_fee
		FROM commissions
		WHERE id = ?
	`

	var before commission.CommissionBeforeUpdate
	err := r.mysqlSess.QueryRow(query, id).Scan(
		&before.ProductCode,
		&before.CommissionPercentage,
		&before.CommissionVATType,
		&before.AFPercentage,
		&before.AFVATType,
		&before.AdminFee,
	)
	if err != nil {
		return nil, err
	}

	return &before, nil
}

func (r *repository) InsertHistory(userName, section, action string, recordID int64, recordType, dataBefore, dataAfter string) error {
	historyQuery := `INSERT INTO histories (user_name, section, action, record_id, record_type, data_before, data_after, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	_, err := r.mysqlSess.Exec(historyQuery, userName, section, action, recordID, recordType, dataBefore, dataAfter)
	if err != nil {
		return fmt.Errorf("failed to insert history: %w", err)
	}

	return nil
}

