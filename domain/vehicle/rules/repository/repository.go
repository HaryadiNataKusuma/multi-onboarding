package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/rules"
	"multi-onboarding/utils"
)

type repository struct{}

// NewRulesRepository creates a new rules repository
func NewRulesRepository() rules.Repository {
	return &repository{}
}

// getDB gets the database connection
func (r *repository) getDB(c echo.Context) *sql.DB {
	db := utils.GetDB()
	if db == nil {
		log.Printf("ERROR: Database connection is nil")
	}
	return db
}

// GetProductsByInsuranceCode gets all active products for an insurance
func (r *repository) GetProductsByInsuranceCode(c echo.Context, insuranceCode string) ([]rules.ProductInfo, error) {
	db := r.getDB(c)
	log.Printf("GetProductsByInsuranceCode called with insurance_code: %s", insuranceCode)

	query := `SELECT code, insurance_code FROM vehicle_service_development.products WHERE insurance_code = ? AND is_active = 1`
	log.Printf("Executing query: %s with param: %s", query, insuranceCode)

	rows, err := db.Query(query, insuranceCode)
	if err != nil {
		log.Printf("ERROR: Failed to query products: %v", err)
		return nil, fmt.Errorf("failed to query products: %v", err)
	}
	defer rows.Close()

	var products []rules.ProductInfo
	for rows.Next() {
		var code, insCode string
		if err := rows.Scan(&code, &insCode); err != nil {
			log.Printf("WARNING: Failed to scan product row: %v", err)
			continue
		}
		products = append(products, rules.ProductInfo{
			Code:          code,
			InsuranceCode: insCode,
		})
		log.Printf("Found product: code=%s, insurance_code=%s", code, insCode)
	}

	if err := rows.Err(); err != nil {
		log.Printf("ERROR: Row iteration error: %v", err)
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	log.Printf("Total products found: %d", len(products))
	return products, nil
}

// GetAdminFeeByProductCode gets admin fee for a product
func (r *repository) GetAdminFeeByProductCode(c echo.Context, productCode string) (float64, error) {
	db := r.getDB(c)
	query := `SELECT admin_fee FROM quotation_service_local.plan_commissions 
		WHERE plan_code = ? AND is_active = 1 
		ORDER BY version DESC LIMIT 1`
	var adminFee sql.NullFloat64
	err := db.QueryRow(query, productCode).Scan(&adminFee)
	if err != nil {
		return 0, err
	}
	if adminFee.Valid {
		return adminFee.Float64, nil
	}
	return 0, nil
}

// LoadDraftData loads draft data from JSON file
func (r *repository) LoadDraftData(c echo.Context) (*utils.DraftData, error) {
	return utils.LoadDraftData(c)
}

// SaveDraftData saves draft data to JSON file
func (r *repository) SaveDraftData(c echo.Context, data *utils.DraftData) error {
	return utils.SaveDraftData(c, data)
}

// GetAllRulesDraft gets all rule drafts
func (r *repository) GetAllRulesDraft(c echo.Context, insuranceCode string) ([]utils.RuleDraft, error) {
	draftData, err := r.LoadDraftData(c)
	if err != nil {
		return nil, err
	}

	if insuranceCode != "" {
		var filtered []utils.RuleDraft
		for _, rule := range draftData.Rules {
			if rule.InsuranceCode == insuranceCode {
				filtered = append(filtered, rule)
			}
		}
		return filtered, nil
	}

	return draftData.Rules, nil
}

// UpdateRuleDraftInMemory updates a draft rule in memory
func (r *repository) UpdateRuleDraftInMemory(c echo.Context, draftData *utils.DraftData, req rules.UpdateRuleDraftRequest) error {
	found := false
	for i, rule := range draftData.Rules {
		if rule.ID == req.ID {
			rule.ID = req.ID
			rule.ProductCode = req.ProductCode
			rule.InsuranceCode = req.InsuranceCode
			rule.VehicleType = req.VehicleType
			rule.VehicleCategory = req.VehicleCategory
			rule.StartVehicleValue = req.StartVehicleValue
			rule.EndVehicleValue = req.EndVehicleValue
			rule.LimitVehicleAge = req.LimitVehicleAge
			rule.AdminFee = req.AdminFee
			rule.HardcopyAdminFee = req.HardcopyAdminFee
			rule.RegionID = req.RegionID
			rule.BasePremiumValue = req.BasePremiumValue
			rule.LoadingFeePremiumValue = req.LoadingFeePremiumValue
			rule.CommercialUsageValue = req.CommercialUsageValue
			rule.StartLoadingAge = req.StartLoadingAge
			rule.AdditionalPremium = req.AdditionalPremium
			rule.TypeAdditionalPremium = req.TypeAdditionalPremium
			rule.Rules = req.Rules
			rule.BasePremiumType = req.BasePremiumType
			rule.BaseLoadingPremiumValue = req.BaseLoadingPremiumValue
			rule.Timestamp = utils.GetTimestamp()
			draftData.Rules[i] = rule
			found = true
			break
		}
	}

	if !found {
		return &notFoundError{Message: "Rule draft not found"}
	}

	return nil
}

// DeleteRuleDraftFromMemory deletes a draft rule from memory
func (r *repository) DeleteRuleDraftFromMemory(c echo.Context, draftData *utils.DraftData, id int) error {
	found := false
	var updatedRules []utils.RuleDraft
	for _, rule := range draftData.Rules {
		if rule.ID == id {
			found = true
			continue
		}
		updatedRules = append(updatedRules, rule)
	}

	if !found {
		return &notFoundError{Message: "Rule draft not found"}
	}

	draftData.Rules = updatedRules
	return nil
}

// ClearRulesDraftInMemory clears all draft rules from memory
func (r *repository) ClearRulesDraftInMemory(c echo.Context, draftData *utils.DraftData) {
	draftData.Rules = []utils.RuleDraft{}
}

// GetAllRules gets all confirmed rules
func (r *repository) GetAllRules(c echo.Context, insuranceCode string) ([]map[string]interface{}, error) {
	db := r.getDB(c)
	var query string
	var rows *sql.Rows
	var err error

	if insuranceCode != "" {
		query = `SELECT id, product_code, insurance_code, vehicle_type, vehicle_category,
			start_vehicle_value, end_vehicle_value, limit_vehicle_age,
			admin_fee, hardcopy_admin_fee, region_id,
			base_premium_value, loading_fee_premium_value, commercial_usage_value,
			start_loading_age, additional_premium, type_additional_premium,
			rules, created_by, created_at, base_premium_type, base_loading_premium_value
			FROM vehicle_service_development.product_rules WHERE insurance_code = ?`
		rows, err = db.Query(query, insuranceCode)
	} else {
		query = `SELECT id, product_code, insurance_code, vehicle_type, vehicle_category,
			start_vehicle_value, end_vehicle_value, limit_vehicle_age,
			admin_fee, hardcopy_admin_fee, region_id,
			base_premium_value, loading_fee_premium_value, commercial_usage_value,
			start_loading_age, additional_premium, type_additional_premium,
			rules, created_by, created_at, base_premium_type, base_loading_premium_value
			FROM vehicle_service_development.product_rules`
		rows, err = db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rulesList []map[string]interface{}
	for rows.Next() {
		var id, limitVehicleAge, regionID, startLoadingAge, createdBy int
		var productCode, insuranceCode, vehicleType, vehicleCategory, typeAdditionalPremium, rulesJSON, basePremiumType string
		var startVehicleValue, endVehicleValue, adminFee, hardcopyAdminFee, basePremiumValue,
			loadingFeePremiumValue, commercialUsageValue, additionalPremium, baseLoadingPremiumValue *float64
		var createdAt string

		err := rows.Scan(
			&id, &productCode, &insuranceCode, &vehicleType, &vehicleCategory,
			&startVehicleValue, &endVehicleValue, &limitVehicleAge,
			&adminFee, &hardcopyAdminFee, &regionID,
			&basePremiumValue, &loadingFeePremiumValue, &commercialUsageValue,
			&startLoadingAge, &additionalPremium, &typeAdditionalPremium,
			&rulesJSON, &createdBy, &createdAt, &basePremiumType, &baseLoadingPremiumValue,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		rule := map[string]interface{}{
			"id":                         id,
			"product_code":               productCode,
			"insurance_code":             insuranceCode,
			"vehicle_type":               vehicleType,
			"vehicle_category":           vehicleCategory,
			"start_vehicle_value":        startVehicleValue,
			"end_vehicle_value":          endVehicleValue,
			"limit_vehicle_age":          limitVehicleAge,
			"admin_fee":                  adminFee,
			"hardcopy_admin_fee":         hardcopyAdminFee,
			"region_id":                  regionID,
			"base_premium_value":         basePremiumValue,
			"loading_fee_premium_value":  loadingFeePremiumValue,
			"commercial_usage_value":     commercialUsageValue,
			"start_loading_age":          startLoadingAge,
			"additional_premium":         additionalPremium,
			"type_additional_premium":    typeAdditionalPremium,
			"rules":                      rulesJSON,
			"created_by":                 createdBy,
			"created_at":                 createdAt,
			"base_premium_type":          basePremiumType,
			"base_loading_premium_value": baseLoadingPremiumValue,
		}
		rulesList = append(rulesList, rule)
	}

	return rulesList, rows.Err()
}

// GetRuleBeforeUpdate gets rule data before update
func (r *repository) GetRuleBeforeUpdate(c echo.Context, id int) (*rules.RuleBeforeUpdate, error) {
	db := r.getDB(c)
	query := `SELECT product_code, insurance_code, vehicle_type, vehicle_category, 
		start_vehicle_value, end_vehicle_value, limit_vehicle_age, admin_fee, hardcopy_admin_fee, 
		region_id, base_premium_value, loading_fee_premium_value, commercial_usage_value, 
		start_loading_age, additional_premium, type_additional_premium, rules, base_premium_type, 
		base_loading_premium_value 
		FROM vehicle_service_development.product_rules WHERE id = ?`

	var before rules.RuleBeforeUpdate
	err := db.QueryRow(query, id).Scan(
		&before.ProductCode, &before.InsuranceCode, &before.VehicleType, &before.VehicleCategory,
		&before.StartVehicleValue, &before.EndVehicleValue, &before.LimitVehicleAge,
		&before.AdminFee, &before.HardcopyAdminFee, &before.RegionID,
		&before.BasePremiumValue, &before.LoadingFeePremiumValue, &before.CommercialUsageValue,
		&before.StartLoadingAge, &before.AdditionalPremium, &before.TypeAdditionalPremium,
		&before.Rules, &before.BasePremiumType, &before.BaseLoadingPremiumValue,
	)
	if err != nil {
		return nil, err
	}

	return &before, nil
}

// UpdateRule updates a confirmed rule
func (r *repository) UpdateRule(c echo.Context, id int, updateData map[string]interface{}) (int64, error) {
	db := r.getDB(c)
	var setParts []string
	var values []interface{}

	fields := []string{
		"start_vehicle_value", "end_vehicle_value", "limit_vehicle_age",
		"admin_fee", "hardcopy_admin_fee", "region_id",
		"base_premium_value", "loading_fee_premium_value", "commercial_usage_value",
		"start_loading_age", "additional_premium", "type_additional_premium",
		"rules", "base_premium_type", "base_loading_premium_value",
	}

	for _, field := range fields {
		if val, ok := updateData[field]; ok {
			setParts = append(setParts, field+" = ?")
			values = append(values, val)
		}
	}

	if len(setParts) == 0 {
		return 0, &validationError{Message: "No fields to update"}
	}

	values = append(values, id)
	query := "UPDATE vehicle_service_development.product_rules SET " + joinStrings(setParts, ", ") + " WHERE id = ?"

	result, err := db.Exec(query, values...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// DeleteRule deletes a confirmed rule
func (r *repository) DeleteRule(c echo.Context, id int) (int64, error) {
	db := r.getDB(c)
	query := "DELETE FROM vehicle_service_development.product_rules WHERE id = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// InsertRule inserts a rule into database
func (r *repository) InsertRule(c echo.Context, rule utils.RuleDraft) error {
	db := r.getDB(c)
	query := `INSERT INTO vehicle_service_development.product_rules (
		product_code, insurance_code, vehicle_type, vehicle_category,
		start_vehicle_value, end_vehicle_value, limit_vehicle_age,
		admin_fee, hardcopy_admin_fee, region_id,
		base_premium_value, loading_fee_premium_value, commercial_usage_value,
		start_loading_age, additional_premium, type_additional_premium,
		rules, created_by, created_at, base_premium_type, base_loading_premium_value
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, ?, ?)`

	basePremiumType := 0
	if rule.BasePremiumType != "" {
		if val, err := strconv.Atoi(rule.BasePremiumType); err == nil {
			basePremiumType = val
		}
	}

	_, err := db.Exec(query,
		rule.ProductCode,
		rule.InsuranceCode,
		rule.VehicleType,
		rule.VehicleCategory,
		rule.StartVehicleValue,
		rule.EndVehicleValue,
		rule.LimitVehicleAge,
		rule.AdminFee,
		rule.HardcopyAdminFee,
		rule.RegionID,
		rule.BasePremiumValue,
		rule.LoadingFeePremiumValue,
		rule.CommercialUsageValue,
		rule.StartLoadingAge,
		rule.AdditionalPremium,
		rule.TypeAdditionalPremium,
		rule.Rules,
		rule.CreatedBy,
		basePremiumType,
		rule.BaseLoadingPremiumValue,
	)

	return err
}

// BulkUpdateRules performs bulk update on confirmed rules
func (r *repository) BulkUpdateRules(c echo.Context, filter map[string]interface{}, updateData map[string]interface{}) (int64, error) {
	db := r.getDB(c)
	var whereParts []string
	var whereValues []interface{}

	if productCode, ok := filter["product_code"].(string); ok && productCode != "" {
		whereParts = append(whereParts, "product_code = ?")
		whereValues = append(whereValues, productCode)
	}
	if regionID, ok := filter["region_id"].(*int); ok && regionID != nil {
		whereParts = append(whereParts, "region_id = ?")
		whereValues = append(whereValues, *regionID)
	}
	if vehicleCategory, ok := filter["vehicle_category"].(string); ok && vehicleCategory != "" {
		whereParts = append(whereParts, "vehicle_category = ?")
		whereValues = append(whereValues, vehicleCategory)
	}

	if len(whereParts) == 0 {
		return 0, &validationError{Message: "At least one filter condition is required"}
	}

	var setParts []string
	var setValues []interface{}

	if rules, ok := updateData["rules"].(*string); ok && rules != nil {
		setParts = append(setParts, "rules = ?")
		setValues = append(setValues, *rules)
	}
	if limitVehicleAge, ok := updateData["limit_vehicle_age"].(*int); ok && limitVehicleAge != nil {
		setParts = append(setParts, "limit_vehicle_age = ?")
		setValues = append(setValues, *limitVehicleAge)
	}
	if productCode, ok := updateData["product_code"].(string); ok && productCode != "" {
		setParts = append(setParts, "product_code = ?")
		setValues = append(setValues, productCode)
	}

	if len(setParts) == 0 {
		return 0, &validationError{Message: "At least one update field is required"}
	}

	query := "UPDATE vehicle_service_development.product_rules SET " + joinStrings(setParts, ", ") + " WHERE " + joinStrings(whereParts, " AND ")
	allValues := append(setValues, whereValues...)

	result, err := db.Exec(query, allValues...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
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

// Helper functions
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// Error types
type notFoundError struct {
	Message string
}

func (e *notFoundError) Error() string {
	return e.Message
}

type validationError struct {
	Message string
}

func (e *validationError) Error() string {
	return e.Message
}

