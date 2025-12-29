package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	model "multi-onboarding/model/Travel"
)

// TemplateRepository defines the interface for template operations
type TemplateRepository interface {
	// Draft operations
	AddDraft(template *model.TemplateDraft) error
	GetAllDrafts(insuranceCode string) ([]model.TemplateDraft, error)
	GetDraftByID(id, locale string) (*model.TemplateDraft, error)
	UpdateDraft(id, locale string, template *model.TemplateDraft) error
	DeleteDraft(id, locale string) error
	ClearDrafts() error
	ConfirmDrafts() error

	// Database operations
	Create(template *model.Template) error
	GetByID(id, locale string) (*model.Template, error)
	GetAll(insuranceCode string) ([]model.Template, error)
	Update(id, locale string, template *model.Template) error
	Delete(id, locale string) error
}

// MySQLTemplateRepository implements TemplateRepository with MySQL database storage
type MySQLTemplateRepository struct {
	db          *sql.DB
	historyRepo HistoryRepository
}

// NewMySQLTemplateRepository creates a new MySQL template repository
func NewMySQLTemplateRepository(db *sql.DB) *MySQLTemplateRepository {
	return &MySQLTemplateRepository{
		db: db,
	}
}

// SetHistoryRepository sets the history repository for saving update histories
func (r *MySQLTemplateRepository) SetHistoryRepository(historyRepo HistoryRepository) {
	r.historyRepo = historyRepo
}

// Draft operations

func (r *MySQLTemplateRepository) AddDraft(template *model.TemplateDraft) error {
	// Convert insurance_code to JSON array if provided
	var insuranceJSON []byte
	if template.InsuranceCode != "" {
		insuranceArray := []string{template.InsuranceCode}
		var err error
		insuranceJSON, err = json.Marshal(insuranceArray)
		if err != nil {
			return fmt.Errorf("failed to marshal insurance: %w", err)
		}
	} else if template.Insurance != "" {
		// If insurance is already a string, convert to array
		insuranceArray := []string{template.Insurance}
		var err error
		insuranceJSON, err = json.Marshal(insuranceArray)
		if err != nil {
			return fmt.Errorf("failed to marshal insurance: %w", err)
		}
	} else {
		insuranceJSON = []byte("[]")
	}

	// Explicitly use travel_service_development.template_drafts
	query := `
		INSERT INTO travel_service_development.template_drafts (id, locale, value, product_code, insurance, created_by, created_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
			value = VALUES(value),
			product_code = VALUES(product_code),
			insurance = VALUES(insurance),
			created_by = VALUES(created_by)
	`

	_, err := r.db.Exec(query,
		template.ID,
		template.Locale,
		template.Value,
		template.ProductCode,
		insuranceJSON,
		template.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to add draft: %w", err)
	}

	template.CreatedAt = time.Now()
	return nil
}

func (r *MySQLTemplateRepository) GetAllDrafts(insuranceCode string) ([]model.TemplateDraft, error) {
	// Query from travel_service_development.template_drafts
	var query string
	var args []interface{}

	if insuranceCode != "" {
		// Explicitly use travel_service_development.template_drafts
		query = `
			SELECT id, locale, value, product_code, insurance, created_by, created_at
			FROM travel_service_development.template_drafts
			WHERE JSON_CONTAINS(insurance, ?)
			ORDER BY created_at DESC
		`
		insuranceJSON, _ := json.Marshal([]string{insuranceCode})
		args = []interface{}{string(insuranceJSON)}
	} else {
		// Explicitly use travel_service_development.template_drafts
		query = `
			SELECT id, locale, value, product_code, insurance, created_by, created_at
			FROM travel_service_development.template_drafts
			ORDER BY created_at DESC
		`
		args = []interface{}{}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query drafts: %w", err)
	}
	defer rows.Close()

	drafts := []model.TemplateDraft{}
	for rows.Next() {
		var draft model.TemplateDraft
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

		// Parse insurance JSON
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

func (r *MySQLTemplateRepository) GetDraftByID(id, locale string) (*model.TemplateDraft, error) {
	// Explicitly use travel_service_development.template_drafts
	query := `
		SELECT id, locale, value, product_code, insurance, created_by, created_at
		FROM travel_service_development.template_drafts
		WHERE id = ? AND locale = ?
	`

	draft := &model.TemplateDraft{}
	var insuranceJSON sql.NullString
	err := r.db.QueryRow(query, id, locale).Scan(
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

	// Parse insurance JSON
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

func (r *MySQLTemplateRepository) UpdateDraft(id, locale string, template *model.TemplateDraft) error {
	// First, get existing draft to preserve insurance_code if not provided
	existingDraft, err := r.GetDraftByID(id, locale)
	if err != nil {
		return fmt.Errorf("failed to get existing draft: %w", err)
	}

	// Convert insurance_code to JSON array
	// Priority: template.InsuranceCode > template.Insurance > existingDraft.InsuranceCode
	var insuranceJSON []byte
	insuranceCode := ""
	if template.InsuranceCode != "" {
		insuranceCode = template.InsuranceCode
	} else if template.Insurance != "" {
		insuranceCode = template.Insurance
	} else if existingDraft.InsuranceCode != "" {
		// Preserve existing insurance_code if not provided
		insuranceCode = existingDraft.InsuranceCode
	} else if existingDraft.Insurance != "" {
		// Preserve existing insurance if not provided
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
		// If no insurance code at all, use empty array
		insuranceJSON = []byte("[]")
	}

	// Preserve product_code if not provided
	productCode := template.ProductCode
	if productCode == "" && existingDraft.ProductCode != "" {
		productCode = existingDraft.ProductCode
	}

	// Explicitly use travel_service_development.template_drafts
	query := `
		UPDATE travel_service_development.template_drafts
		SET value = ?, product_code = ?, insurance = ?, created_by = ?
		WHERE id = ? AND locale = ?
	`

	result, err := r.db.Exec(query,
		template.Value,
		productCode,
		insuranceJSON,
		template.CreatedBy,
		id,
		locale,
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

	template.ID = id
	template.Locale = locale
	template.InsuranceCode = insuranceCode
	template.ProductCode = productCode
	return nil
}

func (r *MySQLTemplateRepository) DeleteDraft(id, locale string) error {
	// Explicitly use travel_service_development.template_drafts
	query := `DELETE FROM travel_service_development.template_drafts WHERE id = ? AND locale = ?`

	result, err := r.db.Exec(query, id, locale)
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

func (r *MySQLTemplateRepository) ClearDrafts() error {
	// Explicitly use travel_service_development.template_drafts
	query := `DELETE FROM travel_service_development.template_drafts`

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to clear drafts: %w", err)
	}

	return nil
}

func (r *MySQLTemplateRepository) ConfirmDrafts() error {
	// Get all drafts
	drafts, err := r.GetAllDrafts("")
	if err != nil {
		return fmt.Errorf("failed to get drafts: %w", err)
	}

	// Insert each draft into templates table
	for _, draft := range drafts {
		template := &model.Template{
			Locale:    draft.Locale,
			ID:        draft.ID,
			Value:     draft.Value,
			CreatedBy: draft.CreatedBy,
		}

		if err := r.Create(template); err != nil {
			return fmt.Errorf("failed to confirm draft %s/%s: %w", draft.ID, draft.Locale, err)
		}
	}

	// Clear drafts after confirmation
	if err := r.ClearDrafts(); err != nil {
		return fmt.Errorf("failed to clear drafts: %w", err)
	}

	return nil
}

// Database operations

func (r *MySQLTemplateRepository) Create(template *model.Template) error {
	// Explicitly use travel_service_development.templates
	query := `
		INSERT INTO travel_service_development.templates (locale, id, value, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE
			value = VALUES(value),
			updated_by = VALUES(updated_by),
			updated_at = NOW()
	`

	_, err := r.db.Exec(query,
		template.Locale,
		template.ID,
		template.Value,
		template.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	template.CreatedAt = time.Now()
	return nil
}

func (r *MySQLTemplateRepository) GetByID(id, locale string) (*model.Template, error) {
	// Explicitly use travel_service_development.templates
	query := `
		SELECT locale, id, value, created_by, created_at, updated_by, updated_at
		FROM travel_service_development.templates
		WHERE locale = ? AND id = ?
	`

	template := &model.Template{}
	var updatedBy sql.NullInt64
	var updatedAt sql.NullTime
	err := r.db.QueryRow(query, locale, id).Scan(
		&template.Locale,
		&template.ID,
		&template.Value,
		&template.CreatedBy,
		&template.CreatedAt,
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
		template.UpdatedBy = &updatedBy.Int64
	}
	if updatedAt.Valid {
		template.UpdatedAt = &updatedAt.Time
	}

	return template, nil
}

func (r *MySQLTemplateRepository) GetAll(insuranceCode string) ([]model.Template, error) {
	// Query from travel_service_development.templates with JOIN to products
	// to get insurance_code and product_name
	var query string
	var args []interface{}

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

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query templates: %w", err)
	}
	defer rows.Close()

	templates := []model.Template{}
	for rows.Next() {
		var template model.Template
		var updatedBy sql.NullInt64
		var updatedAt sql.NullTime
		var insuranceCode sql.NullString
		var productName sql.NullString

		err := rows.Scan(
			&template.Locale,
			&template.ID,
			&template.Value,
			&template.CreatedBy,
			&template.CreatedAt,
			&updatedBy,
			&updatedAt,
			&insuranceCode,
			&productName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}

		if updatedBy.Valid {
			template.UpdatedBy = &updatedBy.Int64
		}
		if updatedAt.Valid {
			template.UpdatedAt = &updatedAt.Time
		}

		templates = append(templates, template)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return templates, nil
}

func (r *MySQLTemplateRepository) Update(id, locale string, template *model.Template) error {
	// Get data before update for comparison and history
	beforeTemplate, err := r.GetByID(id, locale)
	if err != nil {
		return fmt.Errorf("failed to get template before update: %w", err)
	}

	// Build changed fields map (only fields that are different)
	changedBefore := make(map[string]interface{})
	changedAfter := make(map[string]interface{})
	updateFields := []string{}
	updateValues := []interface{}{}

	// Compare each field and build update query dynamically
	if beforeTemplate.Value != template.Value {
		changedBefore["value"] = beforeTemplate.Value
		changedAfter["value"] = template.Value
		updateFields = append(updateFields, "value = ?")
		updateValues = append(updateValues, template.Value)
	}

	// If no fields changed, return early (no update needed)
	if len(updateFields) == 0 {
		return nil
	}

	// Always update updated_by and updated_at
	var updatedBy interface{}
	if template.UpdatedBy != nil {
		updatedBy = *template.UpdatedBy
	} else {
		updatedBy = nil
	}
	updateFields = append(updateFields, "updated_by = ?", "updated_at = NOW()")
	updateValues = append(updateValues, updatedBy)
	updateValues = append(updateValues, locale, id)

	// Build dynamic UPDATE query
	query := fmt.Sprintf(`
		UPDATE travel_service_development.templates
		SET %s
		WHERE locale = ? AND id = ?
	`, strings.Join(updateFields, ", "))

	result, err := r.db.Exec(query, updateValues...)
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

	template.ID = id
	template.Locale = locale
	now := time.Now()
	template.UpdatedAt = &now

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
				// Use composite key for record_id (id + locale)
				// Since record_id is int, we'll use 0 and store identifier in data
				history := &model.History{
					UserName:   "System",
					Section:    "templates",
					Action:     "UPDATE",
					RecordID:   0, // Templates use composite key (id + locale)
					RecordType: "template",
					DataBefore: string(beforeJSON),
					DataAfter:  string(afterJSON),
				}

				if err := r.historyRepo.Create(history); err != nil {
					// Log error but don't fail the update
					fmt.Printf("Warning: failed to save history: %v\n", err)
				} else {
					fmt.Printf("DEBUG: History saved successfully for template %s/%s with %d changed field(s)\n", id, locale, len(changedBefore))
				}
			}
		}
	}

	return nil
}

func (r *MySQLTemplateRepository) Delete(id, locale string) error {
	// Explicitly use travel_service_development.templates
	query := `DELETE FROM travel_service_development.templates WHERE locale = ? AND id = ?`

	result, err := r.db.Exec(query, locale, id)
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
