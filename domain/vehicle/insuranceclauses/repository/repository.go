package repository

import (
	"database/sql"
	"log"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/insuranceclauses"
	"multi-onboarding/utils"
)

type repository struct{}

// NewInsuranceClausesRepository creates a new insurance clauses repository
func NewInsuranceClausesRepository() insuranceclauses.Repository {
	return &repository{}
}

// getDB gets the database connection
func (r *repository) getDB(c echo.Context) *sql.DB {
	return utils.GetDB()
}

// GetProductsByInsuranceCode gets all active products for an insurance
func (r *repository) GetProductsByInsuranceCode(c echo.Context, insuranceCode string) ([]string, error) {
	db := r.getDB(c)
	query := `SELECT code FROM vehicle_service_development.products WHERE insurance_code = ? AND is_active = 1`
	rows, err := db.Query(query, insuranceCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			continue
		}
		products = append(products, code)
	}

	return products, rows.Err()
}

// LoadDraftData loads draft data from JSON file
func (r *repository) LoadDraftData(c echo.Context) (*utils.DraftData, error) {
	return utils.LoadDraftData(c)
}

// SaveDraftData saves draft data to JSON file
func (r *repository) SaveDraftData(c echo.Context, data *utils.DraftData) error {
	return utils.SaveDraftData(c, data)
}

// GetAllClausesDraft gets all clause drafts
func (r *repository) GetAllClausesDraft(c echo.Context, insuranceCode string) ([]utils.ClausesDraft, error) {
	draftData, err := r.LoadDraftData(c)
	if err != nil {
		return nil, err
	}

	if insuranceCode != "" {
		var filtered []utils.ClausesDraft
		for _, item := range draftData.Clauses {
			if item.InsuranceCode == insuranceCode {
				filtered = append(filtered, item)
			}
		}
		return filtered, nil
	}

	return draftData.Clauses, nil
}

// UpdateClauseDraftInMemory updates a draft clause in memory
func (r *repository) UpdateClauseDraftInMemory(c echo.Context, draftData *utils.DraftData, req insuranceclauses.UpdateClauseDraftRequest) error {
	found := false
	for i, clause := range draftData.Clauses {
		if clause.ID == req.ID {
			draftData.Clauses[i].InsuranceCode = req.InsuranceCode
			draftData.Clauses[i].ProductType = req.ProductType
			draftData.Clauses[i].Code = req.Code
			draftData.Clauses[i].Title = req.Title
			draftData.Clauses[i].Content = req.Content
			if req.Description != nil {
				draftData.Clauses[i].Description = *req.Description
			}
			if req.IsMandatory != nil {
				if *req.IsMandatory {
					draftData.Clauses[i].IsMandatory = 1
				} else {
					draftData.Clauses[i].IsMandatory = 0
				}
			}
			if req.IsActive != nil {
				if *req.IsActive {
					draftData.Clauses[i].IsActive = 1
				} else {
					draftData.Clauses[i].IsActive = 0
				}
			}
			draftData.Clauses[i].Timestamp = utils.GetTimestamp()
			found = true
			break
		}
	}

	if !found {
		return &notFoundError{Message: "Clauses draft not found"}
	}

	return nil
}

// DeleteClauseDraftFromMemory deletes a draft clause from memory
func (r *repository) DeleteClauseDraftFromMemory(c echo.Context, draftData *utils.DraftData, id int) error {
	found := false
	var newClauses []utils.ClausesDraft
	for _, r := range draftData.Clauses {
		if r.ID != id {
			newClauses = append(newClauses, r)
		} else {
			found = true
		}
	}

	if !found {
		return &notFoundError{Message: "Clauses draft not found"}
	}

	draftData.Clauses = newClauses
	return nil
}

// ClearClausesDraftInMemory clears all draft clauses from memory
func (r *repository) ClearClausesDraftInMemory(c echo.Context, draftData *utils.DraftData) {
	draftData.Clauses = []utils.ClausesDraft{}
}

// GetAllClauses gets all confirmed clauses
func (r *repository) GetAllClauses(c echo.Context, insuranceCode string) ([]map[string]interface{}, error) {
	db := r.getDB(c)
	var query string
	var args []interface{}

	if insuranceCode != "" {
		query = `SELECT id, insurance_code, product_type, code, title, content, description, is_mandatory, is_active, created_at 
			FROM vehicle_service_development.insurance_clauses 
			WHERE insurance_code = ?`
		args = []interface{}{insuranceCode}
	} else {
		query = `SELECT id, insurance_code, product_type, code, title, content, description, is_mandatory, is_active, created_at 
			FROM vehicle_service_development.insurance_clauses`
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clauses []map[string]interface{}
	for rows.Next() {
		var id int
		var insuranceCode, productType, code, title, content string
		var description sql.NullString
		var isMandatory, isActive sql.NullInt64
		var createdAt sql.NullString

		if err := rows.Scan(&id, &insuranceCode, &productType, &code, &title, &content, &description, &isMandatory, &isActive, &createdAt); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		clause := map[string]interface{}{
			"id":             id,
			"insurance_code": insuranceCode,
			"product_type":   productType,
			"code":           code,
			"title":          title,
			"content":        content,
		}

		if description.Valid {
			clause["description"] = description.String
		} else {
			clause["description"] = ""
		}

		if isMandatory.Valid {
			clause["is_mandatory"] = isMandatory.Int64
		} else {
			clause["is_mandatory"] = 0
		}

		if isActive.Valid {
			clause["is_active"] = isActive.Int64
		} else {
			clause["is_active"] = 1
		}

		if createdAt.Valid {
			clause["created_at"] = createdAt.String
		}

		clauses = append(clauses, clause)
	}

	return clauses, rows.Err()
}

// GetClauseBeforeUpdate gets clause data before update
func (r *repository) GetClauseBeforeUpdate(c echo.Context, id int) (*insuranceclauses.ClauseBeforeUpdate, error) {
	db := r.getDB(c)
	query := `SELECT insurance_code, product_type, code, title, content, description, is_mandatory, is_active 
		FROM vehicle_service_development.insurance_clauses WHERE id = ?`

	var before insuranceclauses.ClauseBeforeUpdate
	err := db.QueryRow(query, id).Scan(
		&before.InsuranceCode, &before.ProductType, &before.Code, &before.Title, &before.Content, &before.Description,
		&before.IsMandatory, &before.IsActive,
	)
	if err != nil {
		return nil, err
	}

	return &before, nil
}

// UpdateClause updates a confirmed clause
func (r *repository) UpdateClause(c echo.Context, id int, updateData map[string]interface{}) (int64, error) {
	db := r.getDB(c)
	query := `UPDATE vehicle_service_development.insurance_clauses 
		SET product_type = ?, code = ?, title = ?, content = ?, description = ?, is_mandatory = ?, is_active = ?
		WHERE id = ?`

	productType := updateData["product_type"].(string)
	code := updateData["code"].(string)
	title := updateData["title"].(string)
	content := updateData["content"].(string)
	description := updateData["description"].(string)
	isMandatory := updateData["is_mandatory"].(int)
	isActive := updateData["is_active"].(int)

	result, err := db.Exec(query, productType, code, title, content, description, isMandatory, isActive, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// InsertClause inserts a clause into database
func (r *repository) InsertClause(c echo.Context, clause utils.ClausesDraft) error {
	db := r.getDB(c)
	query := `INSERT INTO vehicle_service_development.insurance_clauses 
		(insurance_code, product_type, code, title, content, description, is_mandatory, is_active, ordering, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`

	_, err := db.Exec(query,
		clause.InsuranceCode,
		clause.ProductType,
		clause.Code,
		clause.Title,
		clause.Content,
		clause.Description,
		clause.IsMandatory,
		clause.IsActive,
		0,
	)

	return err
}

// InsertHistory inserts a history record
func (r *repository) InsertHistory(c echo.Context, userName, section, action string, recordID int, recordType, dataBefore, dataAfter string) error {
	db := r.getDB(c)
	query := `INSERT INTO vehicle_service_development.histories (user_name, section, action, record_id, record_type, data_before, data_after, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`
	_, err := db.Exec(query, userName, section, action, recordID, recordType, dataBefore, dataAfter)
	if err != nil {
		log.Printf("Error inserting history: %v", err)
		return err
	}
	return nil
}

// Error types
type notFoundError struct {
	Message string
}

func (e *notFoundError) Error() string {
	return e.Message
}


