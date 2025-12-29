package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"multi-onboarding/domain/Travel/template"
	"multi-onboarding/utils"
)

type repository struct {
	mysqlSess *sql.DB
}

// NewTemplateRepository creates a new template repository
func NewTemplateRepository() template.Repository {
	return &repository{
		mysqlSess: utils.GetDB(),
	}
}

// Draft operations

func (r *repository) AddDraft(draft *template.TemplateDraft) error {
	var insuranceJSON []byte
	if draft.InsuranceCode != "" {
		insuranceArray := []string{draft.InsuranceCode}
		var err error
		insuranceJSON, err = json.Marshal(insuranceArray)
		if err != nil {
			return fmt.Errorf("failed to marshal insurance: %w", err)
		}
	} else if draft.Insurance != "" {
		insuranceArray := []string{draft.Insurance}
		var err error
		insuranceJSON, err = json.Marshal(insuranceArray)
		if err != nil {
			return fmt.Errorf("failed to marshal insurance: %w", err)
		}
	} else {
		insuranceJSON = []byte("[]")
	}

	query := `
		INSERT INTO template_drafts (id, locale, value, product_code, insurance, created_by, created_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
			value = VALUES(value),
			product_code = VALUES(product_code),
			insurance = VALUES(insurance),
			created_by = VALUES(created_by)
	`

	_, err := r.mysqlSess.Exec(query,
		draft.ID,
		draft.Locale,
		draft.Value,
		draft.ProductCode,
		insuranceJSON,
		draft.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to add draft: %w", err)
	}

	draft.CreatedAt = time.Now()
	return nil
}

func (r *repository) GetAllDrafts(insuranceCode string) ([]template.TemplateDraft, error) {
	var query string
	var args []interface{}

	if insuranceCode != "" {
		query = `
			SELECT id, locale, value, product_code, insurance, created_by, created_at
			FROM travel_service_development.template_drafts
			WHERE JSON_CONTAINS(insurance, ?)
			ORDER BY created_at DESC
		`
		insuranceJSON, _ := json.Marshal([]string{insuranceCode})
		args = []interface{}{string(insuranceJSON)}
	} else {
		query = `
			SELECT id, locale, value, product_code, insurance, created_by, created_at
			FROM travel_service_development.template_drafts
			ORDER BY created_at DESC
		`
		args = []interface{}{}
	}

	rows, err := r.mysqlSess.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query drafts: %w", err)
	}
	defer rows.Close()

	drafts := []template.TemplateDraft{}
	for rows.Next() {
		var draft template.TemplateDraft
		var insuranceJSON sql.NullString
		err := rows.Scan(
			&draft.ID,
			&draft.Locale,
			&draft.Value,
			&draft.ProductCode,
			&insuranceJSON,
			&draft.CreatedBy,
			&draft.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan draft: %w", err)
		}

		if insuranceJSON.Valid && insuranceJSON.String != "" {
			var insuranceArray []string
			if err := json.Unmarshal([]byte(insuranceJSON.String), &insuranceArray); err == nil {
				if len(insuranceArray) > 0 {
					draft.Insurance = insuranceArray[0]
					draft.InsuranceCode = insuranceArray[0]
				}
			}
		}

		drafts = append(drafts, draft)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return drafts, nil
}

func (r *repository) GetDraftByID(id string) (*template.TemplateDraft, error) {
	query := `
		SELECT id, locale, value, product_code, insurance, created_by, created_at
		FROM template_drafts
		WHERE id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	draft := &template.TemplateDraft{}
	var insuranceJSON sql.NullString
	err := r.mysqlSess.QueryRow(query, id).Scan(
		&draft.ID,
		&draft.Locale,
		&draft.Value,
		&draft.ProductCode,
		&insuranceJSON,
		&draft.CreatedBy,
		&draft.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("draft not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get draft: %w", err)
	}

	if insuranceJSON.Valid && insuranceJSON.String != "" {
		var insuranceArray []string
		if err := json.Unmarshal([]byte(insuranceJSON.String), &insuranceArray); err == nil {
			if len(insuranceArray) > 0 {
				draft.Insurance = insuranceArray[0]
				draft.InsuranceCode = insuranceArray[0]
			}
		}
	}

	return draft, nil
}

func (r *repository) UpdateDraft(id string, draft *template.TemplateDraft) error {
	existingDraft, err := r.GetDraftByID(id)
	if err != nil {
		return fmt.Errorf("failed to get existing draft: %w", err)
	}

	var insuranceJSON []byte
	insuranceCode := ""
	if draft.InsuranceCode != "" {
		insuranceCode = draft.InsuranceCode
	} else if draft.Insurance != "" {
		insuranceCode = draft.Insurance
	} else if existingDraft.InsuranceCode != "" {
		insuranceCode = existingDraft.InsuranceCode
	} else if existingDraft.Insurance != "" {
		insuranceCode = existingDraft.Insurance
	}

	if insuranceCode != "" {
		insuranceArray := []string{insuranceCode}
		var marshalErr error
		insuranceJSON, marshalErr = json.Marshal(insuranceArray)
		if marshalErr != nil {
			return fmt.Errorf("failed to marshal insurance: %w", marshalErr)
		}
	} else {
		insuranceJSON = []byte("[]")
	}

	productCode := draft.ProductCode
	if productCode == "" && existingDraft.ProductCode != "" {
		productCode = existingDraft.ProductCode
	}

	query := `
		UPDATE template_drafts
		SET value = ?, product_code = ?, insurance = ?, created_by = ?
		WHERE id = ?
	`

	result, err := r.mysqlSess.Exec(query,
		draft.Value,
		productCode,
		insuranceJSON,
		draft.CreatedBy,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update draft: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("draft not found")
	}

	draft.ID = id
	draft.InsuranceCode = insuranceCode
	draft.ProductCode = productCode
	return nil
}

func (r *repository) DeleteDraft(id string) error {
	query := `DELETE FROM template_drafts WHERE id = ?`

	result, err := r.mysqlSess.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete draft: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("draft not found")
	}

	return nil
}

func (r *repository) ClearDrafts() error {
	query := `DELETE FROM template_drafts`

	_, err := r.mysqlSess.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to clear drafts: %w", err)
	}

	return nil
}

func (r *repository) ConfirmDrafts(insuranceCodeFilter string) error {
	drafts, err := r.GetAllDrafts(insuranceCodeFilter)
	if err != nil {
		return fmt.Errorf("failed to get drafts: %w", err)
	}

	for _, draft := range drafts {
		t := &template.Template{
			Locale:    draft.Locale,
			ID:        draft.ID,
			Value:     draft.Value,
			CreatedBy: draft.CreatedBy,
		}

		if err := r.Create(t); err != nil {
			return fmt.Errorf("failed to create template from draft: %w", err)
		}
	}

	// Clear drafts after confirmation
	if err := r.ClearDrafts(); err != nil {
		return fmt.Errorf("failed to clear drafts: %w", err)
	}

	return nil
}

// Database operations

// TemplateWithProduct represents template with product information
type TemplateWithProduct struct {
	template.Template
	InsuranceCode string `json:"insurance_code"`
	ProductName   string `json:"product_name"`
}

func (r *repository) GetAll(insuranceCode string) ([]template.Template, error) {
	var query string
	var args []interface{}

	// JOIN with products table to get insurance_code and product name
	// Template ID matches either summary or insurance_detail in products table
	if insuranceCode != "" {
		query = `
			SELECT DISTINCT 
				t.locale, t.id, t.value, t.created_by, t.created_at, t.updated_by, t.updated_at,
				p.insurance_code, p.name as product_name
			FROM travel_service_development.templates t
			INNER JOIN travel_service_development.products p ON (t.id = p.summary OR t.id = p.insurance_detail)
			WHERE p.insurance_code = ?
			ORDER BY t.created_at DESC
		`
		args = []interface{}{insuranceCode}
	} else {
		query = `
			SELECT DISTINCT 
				t.locale, t.id, t.value, t.created_by, t.created_at, t.updated_by, t.updated_at,
				p.insurance_code, p.name as product_name
			FROM travel_service_development.templates t
			INNER JOIN travel_service_development.products p ON (t.id = p.summary OR t.id = p.insurance_detail)
			ORDER BY t.created_at DESC
		`
		args = []interface{}{}
	}

	rows, err := r.mysqlSess.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query templates: %w", err)
	}
	defer rows.Close()

	templates := []template.Template{}
	for rows.Next() {
		var t template.Template
		var updatedBy sql.NullInt64
		var updatedAt sql.NullTime
		var insuranceCode sql.NullString
		var productName sql.NullString

		err := rows.Scan(
			&t.Locale,
			&t.ID,
			&t.Value,
			&t.CreatedBy,
			&t.CreatedAt,
			&updatedBy,
			&updatedAt,
			&insuranceCode,
			&productName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}

		if updatedBy.Valid {
			t.UpdatedBy = &updatedBy.Int64
		}
		if updatedAt.Valid {
			t.UpdatedAt = &updatedAt.Time
		}

		templates = append(templates, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return templates, nil
}

func (r *repository) GetByID(locale, id string) (*template.Template, error) {
	query := `
		SELECT locale, id, value, created_by, created_at, updated_by, updated_at
		FROM templates
		WHERE locale = ? AND id = ?
	`

	var t template.Template
	var updatedBy sql.NullInt64
	var updatedAt sql.NullTime

	err := r.mysqlSess.QueryRow(query, locale, id).Scan(
		&t.Locale,
		&t.ID,
		&t.Value,
		&t.CreatedBy,
		&t.CreatedAt,
		&updatedBy,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("template not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	if updatedBy.Valid {
		t.UpdatedBy = &updatedBy.Int64
	}
	if updatedAt.Valid {
		t.UpdatedAt = &updatedAt.Time
	}

	return &t, nil
}

func (r *repository) Create(t *template.Template) error {
	query := `
		INSERT INTO templates (locale, id, value, created_by, created_at)
		VALUES (?, ?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
			value = VALUES(value),
			created_by = VALUES(created_by)
	`

	_, err := r.mysqlSess.Exec(query,
		t.Locale,
		t.ID,
		t.Value,
		t.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	t.CreatedAt = time.Now()
	return nil
}

func (r *repository) Update(locale, id string, t *template.Template) error {
	// Get data before update
	beforeData, err := r.GetTemplateBeforeUpdate(locale, id)
	if err != nil {
		beforeData = &template.TemplateBeforeUpdate{
			Locale: t.Locale,
			ID:     t.ID,
			Value:  t.Value,
		}
	}

	// Build changed fields
	changedBefore := make(map[string]interface{})
	changedAfter := make(map[string]interface{})
	updateFields := []string{}
	updateValues := []interface{}{}

	if beforeData.Value != t.Value {
		changedBefore["value"] = beforeData.Value
		changedAfter["value"] = t.Value
		updateFields = append(updateFields, "value = ?")
		updateValues = append(updateValues, t.Value)
	}

	if len(updateFields) == 0 {
		return nil // No changes
	}

	var updatedBy interface{}
	if t.UpdatedBy != nil {
		updatedBy = *t.UpdatedBy
		updateFields = append(updateFields, "updated_by = ?")
		updateValues = append(updateValues, updatedBy)
	}

	updateFields = append(updateFields, "updated_at = NOW()")
	updateValues = append(updateValues, locale, id)

	query := fmt.Sprintf(`
		UPDATE templates
		SET %s
		WHERE locale = ? AND id = ?
	`, strings.Join(updateFields, ", "))

	result, err := r.mysqlSess.Exec(query, updateValues...)
	if err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("template not found")
	}

	now := time.Now()
	t.UpdatedAt = &now

	// Save history if there are changes
	if len(changedBefore) > 0 {
		beforeJSON, _ := json.Marshal(changedBefore)
		afterJSON, _ := json.Marshal(changedAfter)
		_ = r.InsertHistory("system", "templates", "UPDATE", id, "template", string(beforeJSON), string(afterJSON))
	}

	return nil
}

func (r *repository) Delete(locale, id string) error {
	query := `DELETE FROM templates WHERE locale = ? AND id = ?`

	result, err := r.mysqlSess.Exec(query, locale, id)
	if err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("template not found")
	}

	return nil
}

func (r *repository) GetTemplateBeforeUpdate(locale, id string) (*template.TemplateBeforeUpdate, error) {
	query := `SELECT locale, id, value FROM templates WHERE locale = ? AND id = ?`

	var before template.TemplateBeforeUpdate
	err := r.mysqlSess.QueryRow(query, locale, id).Scan(
		&before.Locale,
		&before.ID,
		&before.Value,
	)
	if err != nil {
		return nil, err
	}

	return &before, nil
}

func (r *repository) InsertHistory(userName, section, action string, recordID string, recordType, dataBefore, dataAfter string) error {
	historyQuery := `INSERT INTO histories (user_name, section, action, record_id, record_type, data_before, data_after, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	_, err := r.mysqlSess.Exec(historyQuery, userName, section, action, recordID, recordType, dataBefore, dataAfter)
	if err != nil {
		return fmt.Errorf("failed to insert history: %w", err)
	}

	return nil
}
