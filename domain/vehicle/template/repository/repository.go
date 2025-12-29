package repository

import (
	"database/sql"
	"log"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/template"
	"multi-onboarding/utils"
)

type repository struct{}

// NewTemplateRepository creates a new template repository
func NewTemplateRepository() template.Repository {
	return &repository{}
}

// getDB gets the database connection
func (r *repository) getDB(c echo.Context) *sql.DB {
	return utils.GetDB()
}

// LoadDraftData loads draft data from JSON file
func (r *repository) LoadDraftData(c echo.Context) (*utils.DraftData, error) {
	return utils.LoadDraftData(c)
}

// SaveDraftData saves draft data to JSON file
func (r *repository) SaveDraftData(c echo.Context, data *utils.DraftData) error {
	return utils.SaveDraftData(c, data)
}

// GetAllTemplates retrieves all templates from database with JOIN to products
func (r *repository) GetAllTemplates(c echo.Context, insuranceCode string) ([]map[string]interface{}, error) {
	db := r.getDB(c)
	var query string
	var rows *sql.Rows
	var err error

	if insuranceCode != "" {
		query = `SELECT DISTINCT 
			t.id, t.locale, t.value, t.created_by, t.created_at, t.updated_at,
			p.insurance_code, p.code as product_code
		FROM vehicle_service_development.templates t
		INNER JOIN vehicle_service_development.products p ON (t.id = p.summary OR t.id = p.insurance_detail)
		WHERE p.insurance_code = ?
		ORDER BY t.id`
		rows, err = db.Query(query, insuranceCode)
	} else {
		query = `SELECT DISTINCT 
			t.id, t.locale, t.value, t.created_by, t.created_at, t.updated_at,
			p.insurance_code, p.code as product_code
		FROM vehicle_service_development.templates t
		INNER JOIN vehicle_service_development.products p ON (t.id = p.summary OR t.id = p.insurance_detail)
		ORDER BY t.id`
		rows, err = db.Query(query)
	}

	if err != nil {
		log.Printf("Error querying templates: %v", err)
		return nil, err
	}
	defer rows.Close()

	var templates []map[string]interface{}
	for rows.Next() {
		var createdBy int
		var id, locale, value, createdAt string
		var updatedAt sql.NullString
		var insuranceCode, productCode sql.NullString

		err := rows.Scan(&id, &locale, &value, &createdBy, &createdAt, &updatedAt, &insuranceCode, &productCode)
		if err != nil {
			log.Printf("Error scanning template row: %v", err)
			continue
		}

		log.Printf("Template row scanned - id: %s, locale: %s, value length: %d", id, locale, len(value))

		if id == "" {
			log.Printf("WARNING: Template id is empty, skipping this row")
			continue
		}

		tmpl := map[string]interface{}{
			"id":         id,
			"locale":     locale,
			"value":      value,
			"created_by": createdBy,
			"created_at": createdAt,
		}

		log.Printf("Template mapped - id: %s, will be returned in response", id)

		if updatedAt.Valid {
			tmpl["updated_at"] = updatedAt.String
		}

		if insuranceCode.Valid {
			tmpl["insurance_code"] = insuranceCode.String
			tmpl["insurance"] = insuranceCode.String
		}
		if productCode.Valid {
			tmpl["product_code"] = productCode.String
		}

		templates = append(templates, tmpl)
	}

	return templates, rows.Err()
}

// GetTemplateBeforeUpdate retrieves template data before update for history
func (r *repository) GetTemplateBeforeUpdate(c echo.Context, locale, id string) (*template.TemplateBeforeUpdate, error) {
	db := r.getDB(c)
	beforeQuery := `SELECT locale, id, value, created_by FROM vehicle_service_development.templates WHERE locale = ? AND id = ?`

	var beforeLocale, beforeTemplateID, beforeValue sql.NullString
	var beforeCreatedBy sql.NullInt64

	err := db.QueryRow(beforeQuery, locale, id).Scan(&beforeLocale, &beforeTemplateID, &beforeValue, &beforeCreatedBy)
	if err != nil {
		return nil, err
	}

	return &template.TemplateBeforeUpdate{
		Locale:     beforeLocale,
		TemplateID: beforeTemplateID,
		Value:      beforeValue,
		CreatedBy:  beforeCreatedBy,
	}, nil
}

// UpdateTemplate updates a template in the database
func (r *repository) UpdateTemplate(c echo.Context, oldLocale, oldID string, newLocale, newTemplateID, value string, createdBy int) (int64, error) {
	db := r.getDB(c)
	query := `UPDATE vehicle_service_development.templates SET locale = ?, id = ?, value = ?, created_by = ?, updated_at = CURRENT_TIMESTAMP WHERE locale = ? AND id = ?`

	result, err := db.Exec(query, newLocale, newTemplateID, value, createdBy, oldLocale, oldID)
	if err != nil {
		log.Printf("Error updating template %s: %v", oldID, err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return 0, err
	}

	return rowsAffected, nil
}

// InsertHistory inserts a history record
func (r *repository) InsertHistory(c echo.Context, userName, section, action string, recordID string, recordType, dataBefore, dataAfter string) error {
	db := r.getDB(c)
	historyQuery := `INSERT INTO vehicle_service_development.histories (user_name, section, action, record_id, record_type, data_before, data_after, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	_, err := db.Exec(historyQuery, userName, section, action, recordID, recordType, dataBefore, dataAfter)
	if err != nil {
		log.Printf("Error inserting history: %v", err)
		return err
	}

	return nil
}

// InsertTemplate inserts a template into database
func (r *repository) InsertTemplate(c echo.Context, locale, templateID, value string, createdBy int) error {
	db := r.getDB(c)
	query := `INSERT INTO vehicle_service_development.templates (locale, id, value, created_by, created_at) 
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)`

	_, err := db.Exec(query, locale, templateID, value, createdBy)
	if err != nil {
		log.Printf("Error inserting template %s: %v", templateID, err)
		return err
	}

	return nil
}


