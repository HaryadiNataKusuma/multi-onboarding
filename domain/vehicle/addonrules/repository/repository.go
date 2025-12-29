package repository

import (
	"database/sql"
	"log"
	"strings"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/addonrules"
	"multi-onboarding/utils"
)

type repository struct{}

// NewAddonRulesRepository creates a new addon rules repository
func NewAddonRulesRepository() addonrules.Repository {
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

// GetActiveAddons gets all active addons
func (r *repository) GetActiveAddons(c echo.Context) ([]string, error) {
	db := r.getDB(c)
	query := `SELECT code FROM vehicle_service_development.addons WHERE is_active = 1`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	addonMap := make(map[string]bool)
	var addons []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			continue
		}
		if !addonMap[code] {
			addons = append(addons, code)
			addonMap[code] = true
		}
	}

	return addons, rows.Err()
}

// GetLatestAddons gets addons created in the last 24 hours
func (r *repository) GetLatestAddons(c echo.Context) ([]string, error) {
	db := r.getDB(c)
	query := `SELECT code FROM vehicle_service_development.addons WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 DAY)`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addons []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			continue
		}
		addons = append(addons, code)
	}

	return addons, rows.Err()
}

// LoadDraftData loads draft data from JSON file
func (r *repository) LoadDraftData(c echo.Context) (*utils.DraftData, error) {
	return utils.LoadDraftData(c)
}

// SaveDraftData saves draft data to JSON file
func (r *repository) SaveDraftData(c echo.Context, data *utils.DraftData) error {
	return utils.SaveDraftData(c, data)
}

// GetAllAddonRulesDraft gets all addon rule drafts
func (r *repository) GetAllAddonRulesDraft(c echo.Context, insuranceCode string) ([]utils.AddonRuleDraft, error) {
	draftData, err := r.LoadDraftData(c)
	if err != nil {
		return nil, err
	}

	if insuranceCode != "" {
		var filtered []utils.AddonRuleDraft
		for _, rule := range draftData.AddonRules {
			if strings.EqualFold(rule.InsuranceCode, insuranceCode) {
				filtered = append(filtered, rule)
			}
		}
		return filtered, nil
	}

	return draftData.AddonRules, nil
}

// UpdateAddonRuleDraftInMemory updates a draft addon rule in memory
func (r *repository) UpdateAddonRuleDraftInMemory(c echo.Context, draftData *utils.DraftData, req addonrules.UpdateAddonRuleDraftRequest) error {
	found := false
	for i, rule := range draftData.AddonRules {
		if rule.ID == req.ID {
			draftData.AddonRules[i].ProductCode = req.ProductCode
			draftData.AddonRules[i].InsuranceCode = req.InsuranceCode
			draftData.AddonRules[i].Rules = req.Rules
			draftData.AddonRules[i].Timestamp = utils.GetTimestamp()
			found = true
			break
		}
	}

	if !found {
		return &notFoundError{Message: "Addon rule draft not found"}
	}

	return nil
}

// DeleteAddonRuleDraftFromMemory deletes a draft addon rule from memory
func (r *repository) DeleteAddonRuleDraftFromMemory(c echo.Context, draftData *utils.DraftData, id int) error {
	found := false
	for i, rule := range draftData.AddonRules {
		if rule.ID == id {
			draftData.AddonRules = append(draftData.AddonRules[:i], draftData.AddonRules[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return &notFoundError{Message: "Addon rule draft not found"}
	}

	return nil
}

// ClearAddonRulesDraftInMemory clears all draft addon rules from memory
func (r *repository) ClearAddonRulesDraftInMemory(c echo.Context, draftData *utils.DraftData) {
	draftData.AddonRules = []utils.AddonRuleDraft{}
}

// GetAllAddonRules gets all confirmed addon rules
func (r *repository) GetAllAddonRules(c echo.Context, insuranceCode string) ([]map[string]interface{}, error) {
	db := r.getDB(c)
	var query string
	var rows *sql.Rows
	var err error

	if insuranceCode != "" {
		query = `SELECT id, addon_code, product_code, insurance_code, rules, created_by, created_at 
			FROM vehicle_service_development.addon_rules WHERE insurance_code = ?`
		rows, err = db.Query(query, insuranceCode)
	} else {
		query = `SELECT id, addon_code, product_code, insurance_code, rules, created_by, created_at 
			FROM vehicle_service_development.addon_rules`
		rows, err = db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rulesList []map[string]interface{}
	for rows.Next() {
		var id, createdBy int
		var addonCode, productCode, insuranceCode, rulesJSON, createdAt string

		err := rows.Scan(&id, &addonCode, &productCode, &insuranceCode, &rulesJSON, &createdBy, &createdAt)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		rule := map[string]interface{}{
			"id":             id,
			"addon_code":     addonCode,
			"product_code":   productCode,
			"insurance_code": insuranceCode,
			"rules":          rulesJSON,
			"created_by":     createdBy,
			"created_at":     createdAt,
		}
		rulesList = append(rulesList, rule)
	}

	return rulesList, rows.Err()
}

// GetAddonRuleBeforeUpdate gets addon rule data before update
func (r *repository) GetAddonRuleBeforeUpdate(c echo.Context, id int) (*addonrules.AddonRuleBeforeUpdate, error) {
	db := r.getDB(c)
	query := `SELECT addon_code, product_code, insurance_code, rules, created_by FROM vehicle_service_development.addon_rules WHERE id = ?`

	var before addonrules.AddonRuleBeforeUpdate
	err := db.QueryRow(query, id).Scan(&before.AddonCode, &before.ProductCode, &before.InsuranceCode, &before.Rules, &before.CreatedBy)
	if err != nil {
		return nil, err
	}

	return &before, nil
}

// UpdateAddonRule updates a confirmed addon rule
func (r *repository) UpdateAddonRule(c echo.Context, id int, updateData map[string]interface{}) (int64, error) {
	db := r.getDB(c)
	var setParts []string
	var values []interface{}

	if productCode, ok := updateData["product_code"].(string); ok {
		setParts = append(setParts, "product_code = ?")
		values = append(values, productCode)
	}
	if insuranceCode, ok := updateData["insurance_code"].(string); ok {
		setParts = append(setParts, "insurance_code = ?")
		values = append(values, insuranceCode)
	}
	if rules, ok := updateData["rules"].(string); ok {
		setParts = append(setParts, "rules = ?")
		values = append(values, rules)
	}

	if len(setParts) == 0 {
		return 0, &validationError{Message: "No fields to update"}
	}

	values = append(values, id)
	query := "UPDATE vehicle_service_development.addon_rules SET " + joinStrings(setParts, ", ") + " WHERE id = ?"

	result, err := db.Exec(query, values...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// InsertAddonRule inserts an addon rule into database
func (r *repository) InsertAddonRule(c echo.Context, addonCode, productCode, insuranceCode, rules string, createdBy int) error {
	db := r.getDB(c)
	query := `INSERT INTO vehicle_service_development.addon_rules (
		addon_code, product_code, insurance_code, rules, created_by, created_at
	) VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	_, err := db.Exec(query, addonCode, productCode, insuranceCode, rules, createdBy)
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


