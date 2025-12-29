package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/insuranceownrisks"
	"multi-onboarding/utils"
)

type repository struct{}

// NewInsuranceOwnRisksRepository creates a new insurance own risks repository
func NewInsuranceOwnRisksRepository() insuranceownrisks.Repository {
	return &repository{}
}

// getDB gets the database connection
func (r *repository) getDB(c echo.Context) *sql.DB {
	return utils.GetDB()
}

// GetProductsByInsuranceCode gets all active products for an insurance with electric vehicle flag
func (r *repository) GetProductsByInsuranceCode(c echo.Context, insuranceCode string) ([]insuranceownrisks.ProductInfo, error) {
	db := r.getDB(c)
	query := `SELECT code, is_electric_vehicle FROM vehicle_service_development.products WHERE insurance_code = ? AND is_active = 1`
	rows, err := db.Query(query, insuranceCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []insuranceownrisks.ProductInfo
	for rows.Next() {
		var code string
		var isElectricVehicle sql.NullInt64
		if err := rows.Scan(&code, &isElectricVehicle); err != nil {
			continue
		}
		electricValue := 0
		if isElectricVehicle.Valid && isElectricVehicle.Int64 == 1 {
			electricValue = 1
		}
		products = append(products, insuranceownrisks.ProductInfo{
			Code:              code,
			IsElectricVehicle: electricValue,
		})
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

// GetAllOwnRisksDraft gets all own risk drafts
func (r *repository) GetAllOwnRisksDraft(c echo.Context, insuranceCode string) ([]utils.OwnRisksDraft, error) {
	draftData, err := r.LoadDraftData(c)
	if err != nil {
		return nil, err
	}

	if insuranceCode != "" {
		var filtered []utils.OwnRisksDraft
		for _, item := range draftData.OwnRisks {
			if item.InsuranceCode == insuranceCode {
				filtered = append(filtered, item)
			}
		}
		return filtered, nil
	}

	return draftData.OwnRisks, nil
}

// UpdateOwnRiskDraftInMemory updates a draft own risk in memory
func (r *repository) UpdateOwnRiskDraftInMemory(c echo.Context, draftData *utils.DraftData, req insuranceownrisks.UpdateOwnRiskDraftRequest) error {
	found := false
	for i, risk := range draftData.OwnRisks {
		if risk.ID == req.ID {
			draftData.OwnRisks[i].InsuranceCode = req.InsuranceCode
			draftData.OwnRisks[i].ProductType = req.ProductType
			draftData.OwnRisks[i].Code = req.Code
			draftData.OwnRisks[i].Title = req.Title
			if req.Value != nil {
				draftData.OwnRisks[i].Value = strconv.FormatFloat(*req.Value, 'f', -1, 64)
			}
			draftData.OwnRisks[i].ValueType = req.ValueType
			if req.Description != nil {
				draftData.OwnRisks[i].Description = *req.Description
			}
			if req.IsMandatory != nil {
				if *req.IsMandatory {
					draftData.OwnRisks[i].IsMandatory = 1
				} else {
					draftData.OwnRisks[i].IsMandatory = 0
				}
			}
			if req.IsActive != nil {
				if *req.IsActive {
					draftData.OwnRisks[i].IsActive = 1
				} else {
					draftData.OwnRisks[i].IsActive = 0
				}
			}
			if req.IsElectricVehicle != nil {
				draftData.OwnRisks[i].IsElectricVehicle = *req.IsElectricVehicle
			}
			draftData.OwnRisks[i].Timestamp = utils.GetTimestamp()
			found = true
			break
		}
	}

	if !found {
		return &notFoundError{Message: "Own risk draft not found"}
	}

	return nil
}

// DeleteOwnRiskDraftFromMemory deletes a draft own risk from memory
func (r *repository) DeleteOwnRiskDraftFromMemory(c echo.Context, draftData *utils.DraftData, id int) error {
	found := false
	for i, risk := range draftData.OwnRisks {
		if risk.ID == id {
			draftData.OwnRisks = append(draftData.OwnRisks[:i], draftData.OwnRisks[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return &notFoundError{Message: "Own risk draft not found"}
	}

	return nil
}

// ClearOwnRisksDraftInMemory clears all draft own risks from memory
func (r *repository) ClearOwnRisksDraftInMemory(c echo.Context, draftData *utils.DraftData) {
	draftData.OwnRisks = []utils.OwnRisksDraft{}
}

// GetAllOwnRisks gets all confirmed own risks
func (r *repository) GetAllOwnRisks(c echo.Context, insuranceCode string) ([]map[string]interface{}, error) {
	db := r.getDB(c)
	var query string
	var args []interface{}

	if insuranceCode != "" {
		query = `SELECT id, insurance_code, product_type, code, title, value, value_type, description, is_mandatory, is_active, created_at 
			FROM vehicle_service_development.insurance_own_risks 
			WHERE insurance_code = ?`
		args = []interface{}{insuranceCode}
	} else {
		query = `SELECT id, insurance_code, product_type, code, title, value, value_type, description, is_mandatory, is_active, created_at 
			FROM vehicle_service_development.insurance_own_risks`
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ownRisks []map[string]interface{}
	for rows.Next() {
		var id int
		var insuranceCode, productType, code, title, value string
		var valueType, description sql.NullString
		var isMandatory, isActive sql.NullInt64
		var createdAt sql.NullString

		if err := rows.Scan(&id, &insuranceCode, &productType, &code, &title, &value, &valueType, &description, &isMandatory, &isActive, &createdAt); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		ownRisk := map[string]interface{}{
			"id":             id,
			"insurance_code": insuranceCode,
			"product_type":   productType,
			"code":           code,
			"title":          title,
			"value":          value,
		}

		if valueType.Valid {
			ownRisk["value_type"] = valueType.String
		} else {
			ownRisk["value_type"] = ""
		}

		if description.Valid {
			ownRisk["description"] = description.String
		} else {
			ownRisk["description"] = ""
		}

		if isMandatory.Valid {
			ownRisk["is_mandatory"] = isMandatory.Int64
		} else {
			ownRisk["is_mandatory"] = 0
		}

		if isActive.Valid {
			ownRisk["is_active"] = isActive.Int64
		} else {
			ownRisk["is_active"] = 1
		}

		if createdAt.Valid {
			ownRisk["created_at"] = createdAt.String
		}

		ownRisks = append(ownRisks, ownRisk)
	}

	return ownRisks, rows.Err()
}

// GetOwnRiskBeforeUpdate gets own risk data before update
func (r *repository) GetOwnRiskBeforeUpdate(c echo.Context, id int) (*insuranceownrisks.OwnRiskBeforeUpdate, error) {
	db := r.getDB(c)
	query := `SELECT insurance_code, product_type, code, title, value, value_type, description, is_mandatory, is_active 
		FROM vehicle_service_development.insurance_own_risks WHERE id = ?`

	var before insuranceownrisks.OwnRiskBeforeUpdate
	err := db.QueryRow(query, id).Scan(
		&before.InsuranceCode, &before.ProductType, &before.Code, &before.Title, &before.Value, &before.ValueType, &before.Description,
		&before.IsMandatory, &before.IsActive,
	)
	if err != nil {
		return nil, err
	}

	return &before, nil
}

// UpdateOwnRisk updates a confirmed own risk
func (r *repository) UpdateOwnRisk(c echo.Context, id int, updateData map[string]interface{}) (int64, error) {
	db := r.getDB(c)
	if db == nil {
		return 0, fmt.Errorf("database connection is nil")
	}

	// Helper function to safely get string from map
	getString := func(key string, defaultValue string) string {
		if val, ok := updateData[key]; ok && val != nil {
			if str, ok := val.(string); ok {
				return str
			}
			return fmt.Sprintf("%v", val)
		}
		return defaultValue
	}

	// Helper function to safely get int from map
	getInt := func(key string, defaultValue int) int {
		if val, ok := updateData[key]; ok && val != nil {
			if i, ok := val.(int); ok {
				return i
			}
			if i, ok := val.(int64); ok {
				return int(i)
			}
			if f, ok := val.(float64); ok {
				return int(f)
			}
			if str, ok := val.(string); ok {
				if parsed, err := strconv.Atoi(str); err == nil {
					return parsed
				}
			}
		}
		return defaultValue
	}

	// Get current data first to preserve unchanged fields
	beforeData, err := r.GetOwnRiskBeforeUpdate(c, id)
	if err != nil {
		return 0, fmt.Errorf("failed to get current data: %v", err)
	}

	// Use provided values or fallback to current values
	productType := getString("product_type", "")
	if productType == "" && beforeData.ProductType.Valid {
		productType = beforeData.ProductType.String
	}

	code := getString("code", "")
	if code == "" && beforeData.Code.Valid {
		code = beforeData.Code.String
	}

	title := getString("title", "")
	if title == "" && beforeData.Title.Valid {
		title = beforeData.Title.String
	}

	value := getString("value", "")
	if value == "" && beforeData.Value.Valid {
		value = fmt.Sprintf("%.2f", beforeData.Value.Float64)
	}

	valueType := getString("value_type", "")
	if valueType == "" && beforeData.ValueType.Valid {
		valueType = beforeData.ValueType.String
	}

	description := getString("description", "")
	if description == "" && beforeData.Description.Valid {
		description = beforeData.Description.String
	}

	isMandatory := getInt("is_mandatory", 0)
	if isMandatory == 0 && beforeData.IsMandatory.Valid {
		if beforeData.IsMandatory.Bool {
			isMandatory = 1
		}
	}

	isActive := getInt("is_active", 1)
	if isActive == 0 && beforeData.IsActive.Valid {
		if !beforeData.IsActive.Bool {
			isActive = 0
		} else {
			isActive = 1
		}
	}

	query := `UPDATE vehicle_service_development.insurance_own_risks 
		SET product_type = ?, code = ?, title = ?, value = ?, value_type = ?, description = ?, is_mandatory = ?, is_active = ?
		WHERE id = ?`

	result, err := db.Exec(query, productType, code, title, value, valueType, description, isMandatory, isActive, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// InsertOwnRisk inserts an own risk into database
func (r *repository) InsertOwnRisk(c echo.Context, risk utils.OwnRisksDraft) error {
	db := r.getDB(c)
	query := `INSERT INTO vehicle_service_development.insurance_own_risks 
		(insurance_code, product_type, code, title, value, value_type, description, is_mandatory, is_active, ordering, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`

	// Convert value string to float64
	valueFloat := 0.0
	if risk.Value != "" {
		if val, err := strconv.ParseFloat(risk.Value, 64); err == nil {
			valueFloat = val
		}
	}

	// Set default values
	isMandatory := risk.IsMandatory
	if isMandatory == 0 {
		isMandatory = 1
	}

	isActive := risk.IsActive
	if isActive == 0 {
		isActive = 1
	}

	_, err := db.Exec(query,
		risk.InsuranceCode,
		risk.ProductType,
		risk.Code,
		risk.Title,
		valueFloat,
		risk.ValueType,
		risk.Description,
		isMandatory,
		isActive,
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


