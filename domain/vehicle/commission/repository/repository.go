package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/commission"
	"multi-onboarding/utils"
)

type repository struct{}

// NewCommissionRepository creates a new commission repository
func NewCommissionRepository() commission.Repository {
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

// GetProductsByInsuranceCode gets all active products for an insurance
func (r *repository) GetProductsByInsuranceCode(c echo.Context, insuranceCode string) ([]string, error) {
	db := r.getDB(c)
	query := `SELECT code FROM vehicle_service_development.products WHERE insurance_code = ? AND is_active = 1`
	rows, err := db.Query(query, insuranceCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productCodes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			continue
		}
		productCodes = append(productCodes, code)
	}

	return productCodes, rows.Err()
}

// ClearDraftTables clears all commission draft tables
func (r *repository) ClearDraftTables(c echo.Context) error {
	db := r.getDB(c)
	clearQueries := []string{
		"DELETE FROM commission_service_local.commissions_draft",
		"DELETE FROM quotation_service_local.plan_commissions_draft",
		"DELETE FROM agent_service_local.default_config_products_draft",
	}

	for _, clearQuery := range clearQueries {
		_, err := db.Exec(clearQuery)
		if err != nil {
			return err
		}
	}

	return nil
}

// ClearDraftTablesForInsurance clears draft tables for a specific insurance code
func (r *repository) ClearDraftTablesForInsurance(c echo.Context, insuranceCode string) error {
	db := r.getDB(c)
	log.Printf("Clearing draft tables for insurance_code: %s", insuranceCode)

	// Clear commissions_draft for this insurance code
	query1 := `DELETE FROM commission_service_local.commissions_draft WHERE insurance_code = ?`
	_, err := db.Exec(query1, insuranceCode)
	if err != nil {
		log.Printf("ERROR: Failed to clear commissions_draft: %v", err)
		return fmt.Errorf("failed to clear commissions_draft: %v", err)
	}

	// Clear plan_commissions_draft for products with this insurance code
	query2 := `DELETE pcd FROM quotation_service_local.plan_commissions_draft pcd
		INNER JOIN vehicle_service_development.products p ON pcd.plan_code = p.code
		WHERE p.insurance_code = ?`
	_, err = db.Exec(query2, insuranceCode)
	if err != nil {
		log.Printf("ERROR: Failed to clear plan_commissions_draft: %v", err)
		return fmt.Errorf("failed to clear plan_commissions_draft: %v", err)
	}

	// Clear default_config_products_draft for this insurance code
	query3 := `DELETE FROM agent_service_local.default_config_products_draft WHERE insurer = ?`
	_, err = db.Exec(query3, insuranceCode)
	if err != nil {
		log.Printf("ERROR: Failed to clear default_config_products_draft: %v", err)
		return fmt.Errorf("failed to clear default_config_products_draft: %v", err)
	}

	log.Printf("SUCCESS: Cleared all draft tables for insurance_code: %s", insuranceCode)
	return nil
}

// InsertCommissionDraft inserts a commission draft
func (r *repository) InsertCommissionDraft(c echo.Context, productCode, insuranceCode, agentLevel string, corporateID int, basicCommission float64, note string) error {
	db := r.getDB(c)
	query := `INSERT INTO commission_service_local.commissions_draft 
		(product_code, insurance_code, agent_level, corporate_id, basic_commission, note, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`
	_, err := db.Exec(query, productCode, insuranceCode, agentLevel, corporateID, basicCommission, note)
	return err
}

// InsertPlanCommissionDraft inserts a plan commission draft
func (r *repository) InsertPlanCommissionDraft(c echo.Context, productCode, commissionVatType string, basicCommission, afPercentage float64, afVatType string, adminFee float64) error {
	db := r.getDB(c)
	query := `INSERT INTO quotation_service_local.plan_commissions_draft 
		(plan_code, commission_percentage, commission_vat_type, af_percentage, af_vat_type, admin_fee, is_active, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, 1, CURRENT_TIMESTAMP)`
	_, err := db.Exec(query, productCode, basicCommission, commissionVatType, afPercentage, afVatType, adminFee)
	return err
}

// InsertConfigProductDraft inserts a config product draft
func (r *repository) InsertConfigProductDraft(c echo.Context, level, insurer, product string, basicCommission float64) error {
	db := r.getDB(c)
	query := `INSERT INTO agent_service_local.default_config_products_draft 
		(level, insurer, product, commission, max_discount, created_at) 
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`
	_, err := db.Exec(query, level, insurer, product, basicCommission, basicCommission)
	return err
}

// GetCommissionDraft gets commission drafts
func (r *repository) GetCommissionDraft(c echo.Context, insuranceCode string) ([]map[string]interface{}, error) {
	db := r.getDB(c)
	var query string
	var rows *sql.Rows
	var err error

	if insuranceCode == "" {
		query = `SELECT id, product_code, insurance_code, agent_level, corporate_id, basic_commission, 
			bonus_point, note, created_at
			FROM commission_service_local.commissions_draft
			ORDER BY created_at DESC`
		rows, err = db.Query(query)
	} else {
		query = `SELECT id, product_code, insurance_code, agent_level, corporate_id, basic_commission, 
			bonus_point, note, created_at
			FROM commission_service_local.commissions_draft
			WHERE insurance_code = ?
			ORDER BY created_at DESC`
		rows, err = db.Query(query, insuranceCode)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commissions []map[string]interface{}
	for rows.Next() {
		var id, corporateID int
		var productCode, insuranceCode, agentLevel string
		var basicCommission float64
		var bonusPoint sql.NullFloat64
		var note sql.NullString
		var createdAt interface{}

		err := rows.Scan(&id, &productCode, &insuranceCode, &agentLevel, &corporateID,
			&basicCommission, &bonusPoint, &note, &createdAt)
		if err != nil {
			continue
		}

		bonusPointValue := 0.0
		if bonusPoint.Valid {
			bonusPointValue = bonusPoint.Float64
		}
		noteValue := ""
		if note.Valid {
			noteValue = note.String
		}

		commissions = append(commissions, map[string]interface{}{
			"id":               id,
			"product_code":     productCode,
			"insurance_code":   insuranceCode,
			"agent_level":      agentLevel,
			"corporate_id":     corporateID,
			"basic_commission": basicCommission,
			"bonus_point":      bonusPointValue,
			"note":             noteValue,
			"created_at":       createdAt,
		})
	}

	return commissions, rows.Err()
}

// GetPlanCommissionDraft gets plan commission drafts
func (r *repository) GetPlanCommissionDraft(c echo.Context, insuranceCode string) ([]map[string]interface{}, error) {
	db := r.getDB(c)
	var query string
	var rows *sql.Rows
	var err error

	if insuranceCode == "" {
		query = `SELECT id, plan_code, commission_percentage, commission_vat_type,
			af_percentage, af_vat_type, admin_fee, is_active, created_at
			FROM quotation_service_local.plan_commissions_draft
			ORDER BY created_at DESC`
		rows, err = db.Query(query)
	} else {
		query = `SELECT id, plan_code, commission_percentage, commission_vat_type,
			af_percentage, af_vat_type, admin_fee, is_active, created_at
			FROM quotation_service_local.plan_commissions_draft
			WHERE plan_code LIKE ?
			ORDER BY created_at DESC`
		rows, err = db.Query(query, "%"+insuranceCode+"%")
	}
	if err != nil {
		log.Printf("ERROR: Failed to query plan commission draft: %v", err)
		return nil, err
	}
	defer rows.Close()

	var planCommissions []map[string]interface{}
	for rows.Next() {
		var id int
		var planCode, commissionVatType, afVatType string
		var commissionPercentage, afPercentage, adminFee float64
		var isActive bool
		var createdAt interface{}

		err := rows.Scan(&id, &planCode, &commissionPercentage, &commissionVatType,
			&afPercentage, &afVatType, &adminFee, &isActive, &createdAt)
		if err != nil {
			log.Printf("ERROR: Failed to scan plan commission draft row: %v", err)
			continue
		}

		planCommissions = append(planCommissions, map[string]interface{}{
			"id":                    id,
			"plan_code":             planCode,
			"commission_percentage": commissionPercentage,
			"commission_vat_type":   commissionVatType,
			"af_percentage":         afPercentage,
			"af_vat_type":           afVatType,
			"admin_fee":             adminFee,
			"is_active":             isActive,
			"created_at":            createdAt,
		})
	}

	return planCommissions, rows.Err()
}

// GetConfigProductDraft gets config product drafts
func (r *repository) GetConfigProductDraft(c echo.Context, insuranceCode string) ([]map[string]interface{}, error) {
	db := r.getDB(c)
	var query string
	var rows *sql.Rows
	var err error

	if insuranceCode == "" {
		query = `SELECT id, level, insurer, product, commission, max_discount, created_at
			FROM agent_service_local.default_config_products_draft
			ORDER BY created_at DESC`
		rows, err = db.Query(query)
	} else {
		query = `SELECT id, level, insurer, product, commission, max_discount, created_at
			FROM agent_service_local.default_config_products_draft
			WHERE insurer = ?
			ORDER BY created_at DESC`
		rows, err = db.Query(query, insuranceCode)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configProducts []map[string]interface{}
	for rows.Next() {
		var id int
		var level, insurer, product string
		var commission, maxDiscount float64
		var createdAt interface{}

		err := rows.Scan(&id, &level, &insurer, &product, &commission, &maxDiscount, &createdAt)
		if err != nil {
			continue
		}

		configProducts = append(configProducts, map[string]interface{}{
			"id":           id,
			"level":        level,
			"insurer":      insurer,
			"product":      product,
			"commission":   commission,
			"max_discount": maxDiscount,
			"created_at":   createdAt,
		})
	}

	return configProducts, rows.Err()
}

// MoveCommissionsDraftToConfirmed moves commissions from draft to confirmed
func (r *repository) MoveCommissionsDraftToConfirmed(c echo.Context) (int64, error) {
	db := r.getDB(c)
	query := `INSERT INTO commission_service_local.commissions 
		(product_code, insurance_code, agent_level, corporate_id, basic_commission, bonus_point, note, created_at) 
		SELECT product_code, insurance_code, agent_level, corporate_id, basic_commission, bonus_point, note, CURRENT_TIMESTAMP 
		FROM commission_service_local.commissions_draft`
	result, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// MovePlanCommissionsDraftToConfirmed moves plan commissions from draft to confirmed
func (r *repository) MovePlanCommissionsDraftToConfirmed(c echo.Context) (int64, error) {
	db := r.getDB(c)
	query := `
		INSERT INTO quotation_service_local.plan_commissions 
		(plan_code, product_id, insurer_id, commission_percentage, commission_vat_type, 
		 af_percentage, af_vat_type, admin_fee, hardcopy_fee, is_active, version, created_at) 
		SELECT 
			pcd.plan_code,
			2 as product_id,
			qi.id as insurer_id,
			pcd.commission_percentage,
			pcd.commission_vat_type,
			pcd.af_percentage,
			pcd.af_vat_type,
			pcd.admin_fee,
			0 as hardcopy_fee,
			pcd.is_active,
			1 as version,
			CURRENT_TIMESTAMP
		FROM quotation_service_local.plan_commissions_draft pcd
		INNER JOIN vehicle_service_development.products p ON pcd.plan_code = p.code
		INNER JOIN vehicle_service_development.insurances qi ON p.insurance_code = qi.code
		WHERE p.is_active = 1`
	result, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// MoveConfigProductsDraftToConfirmed moves config products from draft to confirmed
func (r *repository) MoveConfigProductsDraftToConfirmed(c echo.Context) (int64, error) {
	db := r.getDB(c)
	query := `INSERT INTO agent_service_local.default_config_products 
		(level, insurer, product, commission, max_discount, kind, category, point, upline_bonus, bonus, selected, created_at) 
		SELECT level, insurer, product, commission, max_discount, 'PL' as kind, 'MV_CAR' as category, 0 as point, 0 as upline_bonus, 0 as bonus, 1 as selected, CURRENT_TIMESTAMP 
		FROM agent_service_local.default_config_products_draft`
	result, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// UpdateProductRulesHardcopyAdminFee updates hardcopy_admin_fee in product_rules
func (r *repository) UpdateProductRulesHardcopyAdminFee(c echo.Context) error {
	db := r.getDB(c)
	query := `
		UPDATE vehicle_service_development.product_rules pr
		INNER JOIN quotation_service_local.plan_commissions pc ON pr.product_code = pc.plan_code
		SET pr.hardcopy_admin_fee = pc.admin_fee, pr.updated_at = CURRENT_TIMESTAMP
		WHERE pc.is_active = 1
	`
	_, err := db.Exec(query)
	return err
}

// GetConfirmedCommissions gets confirmed commissions
func (r *repository) GetConfirmedCommissions(c echo.Context, insuranceCode, productCode string) ([]map[string]interface{}, error) {
	db := r.getDB(c)
	query := `
		SELECT id, product_code, insurance_code, agent_level, corporate_id, basic_commission, 
		       bonus_point, note, start_date, end_date, created_at, updated_at
		FROM commission_service_local.commissions
	`
	params := []interface{}{}
	if insuranceCode != "" {
		query += " WHERE insurance_code = ?"
		params = append(params, insuranceCode)
	}
	if productCode != "" {
		if insuranceCode != "" {
			query += " AND product_code = ?"
		} else {
			query += " WHERE product_code = ?"
		}
		params = append(params, productCode)
	}
	query += " ORDER BY created_at DESC"

	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commissions []map[string]interface{}
	for rows.Next() {
		var id, corporateID int
		var productCode, insuranceCode, agentLevel, note string
		var basicCommission, bonusPoint float64
		var startDate, endDate, createdAt, updatedAt sql.NullString

		err := rows.Scan(&id, &productCode, &insuranceCode, &agentLevel, &corporateID,
			&basicCommission, &bonusPoint, &note, &startDate, &endDate, &createdAt, &updatedAt)
		if err != nil {
			continue
		}

		comm := map[string]interface{}{
			"id":               id,
			"product_code":     productCode,
			"insurance_code":   insuranceCode,
			"agent_level":      agentLevel,
			"corporate_id":     corporateID,
			"basic_commission": basicCommission,
			"bonus_point":      bonusPoint,
			"note":             note,
		}

		if startDate.Valid {
			comm["start_date"] = startDate.String
		}
		if endDate.Valid {
			comm["end_date"] = endDate.String
		}
		if createdAt.Valid {
			comm["created_at"] = createdAt.String
		}
		if updatedAt.Valid {
			comm["updated_at"] = updatedAt.String
		}

		commissions = append(commissions, comm)
	}

	return commissions, rows.Err()
}

// GetConfirmedPlanCommissions gets confirmed plan commissions
func (r *repository) GetConfirmedPlanCommissions(c echo.Context, insuranceCode, productCode string) ([]map[string]interface{}, error) {
	db := r.getDB(c)
	query := `
		SELECT id, plan_code, product_id, insurer_id, commission_percentage, commission_vat_type,
		       af_percentage, af_vat_type, admin_fee, hardcopy_fee, is_active, version, created_at
		FROM quotation_service_local.plan_commissions
		WHERE is_active = 1
	`
	params := []interface{}{}
	if insuranceCode != "" || productCode != "" {
		query += " AND ("
		if insuranceCode != "" {
			query += "plan_code LIKE ?"
			params = append(params, "%"+insuranceCode+"%")
		}
		if productCode != "" {
			if insuranceCode != "" {
				query += " OR plan_code = ?"
			} else {
				query += "plan_code = ?"
			}
			params = append(params, productCode)
		}
		query += ")"
	}
	query += " ORDER BY created_at DESC"

	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var planCommissions []map[string]interface{}
	for rows.Next() {
		var id, productID, insurerID, version int
		var planCode, commissionVatType, afVatType string
		var commissionPercentage, afPercentage, adminFee, hardcopyFee float64
		var isActive bool
		var createdAt sql.NullString

		err := rows.Scan(&id, &planCode, &productID, &insurerID, &commissionPercentage, &commissionVatType,
			&afPercentage, &afVatType, &adminFee, &hardcopyFee, &isActive, &version, &createdAt)
		if err != nil {
			continue
		}

		planCommissions = append(planCommissions, map[string]interface{}{
			"id":                    id,
			"plan_code":             planCode,
			"product_id":            productID,
			"insurer_id":            insurerID,
			"commission_percentage": commissionPercentage,
			"commission_vat_type":   commissionVatType,
			"af_percentage":         afPercentage,
			"af_vat_type":           afVatType,
			"admin_fee":             adminFee,
			"hardcopy_fee":          hardcopyFee,
			"is_active":             isActive,
			"version":               version,
			"created_at":            createdAt.String,
		})
	}

	return planCommissions, rows.Err()
}

// GetConfirmedConfigProducts gets confirmed config products
func (r *repository) GetConfirmedConfigProducts(c echo.Context, insuranceCode, productCode string) ([]map[string]interface{}, error) {
	db := r.getDB(c)
	query := `
		SELECT id, level, sequence, kind, category, insurer, product, selected, commission, 
		       point, upline_bonus, bonus, max_discount, ` + "`order`" + `, created_at, updated_at
		FROM agent_service_local.default_config_products
	`
	params := []interface{}{}
	if insuranceCode != "" {
		query += " WHERE insurer = ?"
		params = append(params, insuranceCode)
	}
	if productCode != "" {
		if insuranceCode != "" {
			query += " AND product = ?"
		} else {
			query += " WHERE product = ?"
		}
		params = append(params, productCode)
	}
	query += " ORDER BY created_at DESC"

	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configProducts []map[string]interface{}
	for rows.Next() {
		var id, sequence, order, selected int
		var level, kind, category, insurer, product string
		var commission, point, uplineBonus, bonus, maxDiscount float64
		var createdAt, updatedAt sql.NullString

		err := rows.Scan(&id, &level, &sequence, &kind, &category, &insurer, &product, &selected,
			&commission, &point, &uplineBonus, &bonus, &maxDiscount, &order, &createdAt, &updatedAt)
		if err != nil {
			continue
		}

		configProducts = append(configProducts, map[string]interface{}{
			"id":           id,
			"level":        level,
			"sequence":     sequence,
			"kind":         kind,
			"category":     category,
			"insurer":      insurer,
			"product":      product,
			"selected":     selected,
			"commission":   commission,
			"point":        point,
			"upline_bonus": uplineBonus,
			"bonus":        bonus,
			"max_discount": maxDiscount,
			"order":        order,
			"created_at":   createdAt.String,
			"updated_at":   updatedAt.String,
		})
	}

	return configProducts, rows.Err()
}

// GetPlanCommissionByIDs gets plan commissions by IDs
func (r *repository) GetPlanCommissionByIDs(c echo.Context, ids []int) ([]commission.PlanCommissionRecord, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	db := r.getDB(c)
	placeholders := make([]string, len(ids))
	queryArgs := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		queryArgs[i] = id
	}

	query := `
		SELECT id, plan_code, product_id, insurer_id, commission_percentage, commission_vat_type,
		       af_percentage, af_vat_type, admin_fee, hardcopy_fee, is_active, version
		FROM quotation_service_local.plan_commissions
		WHERE id IN (` + strings.Join(placeholders, ",") + `)`

	rows, err := db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []commission.PlanCommissionRecord
	for rows.Next() {
		var rec commission.PlanCommissionRecord
		err := rows.Scan(&rec.ID, &rec.PlanCode, &rec.ProductID, &rec.InsurerID,
			&rec.CommissionPercentage, &rec.CommissionVatType,
			&rec.AfPercentage, &rec.AfVatType, &rec.AdminFee, &rec.HardcopyFee,
			&rec.IsActive, &rec.Version)
		if err != nil {
			continue
		}
		records = append(records, rec)
	}

	return records, rows.Err()
}

// GetInsurerIDByPlanCode gets insurer ID by plan code
func (r *repository) GetInsurerIDByPlanCode(c echo.Context, planCode string) (int, error) {
	db := r.getDB(c)
	query := `
		SELECT qi.id
		FROM vehicle_service_development.insurances qi
		INNER JOIN vehicle_service_development.products p ON qi.code = p.insurance_code
		WHERE p.code = ? AND p.is_active = 1
		LIMIT 1`
	var insurerID int
	err := db.QueryRow(query, planCode).Scan(&insurerID)
	return insurerID, err
}

// GetInsuranceCodeByPlanCode gets insurance code by plan code
func (r *repository) GetInsuranceCodeByPlanCode(c echo.Context, planCode string) (string, error) {
	db := r.getDB(c)
	query := `
		SELECT p.insurance_code
		FROM vehicle_service_development.products p
		WHERE p.code = ? AND p.is_active = 1
		LIMIT 1`
	var insuranceCode string
	err := db.QueryRow(query, planCode).Scan(&insuranceCode)
	return insuranceCode, err
}

// InsertPlanCommission inserts a new plan commission
func (r *repository) InsertPlanCommission(c echo.Context, planCode string, insurerID, productID int, commissionPercentage float64, commissionVatType string, afPercentage float64, afVatType string, adminFee, hardcopyFee float64, version int) error {
	db := r.getDB(c)
	query := `
		INSERT INTO quotation_service_local.plan_commissions 
		(plan_code, product_id, insurer_id, commission_percentage, commission_vat_type,
		 af_percentage, af_vat_type, admin_fee, hardcopy_fee, is_active, version, created_at)
		VALUES (?, 2, ?, ?, ?, ?, ?, ?, ?, 1, ?, CURRENT_TIMESTAMP)`
	_, err := db.Exec(query, planCode, insurerID, commissionPercentage, commissionVatType,
		afPercentage, afVatType, adminFee, hardcopyFee, version)
	return err
}

// DeactivatePlanCommissions deactivates plan commissions by plan code
func (r *repository) DeactivatePlanCommissions(c echo.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	db := r.getDB(c)
	planCodes, err := r.getPlanCodesByIDs(c, ids)
	if err != nil {
		return err
	}

	for _, planCode := range planCodes {
		query := `UPDATE quotation_service_local.plan_commissions SET is_active = 0 WHERE plan_code = ? AND is_active = 1`
		_, err := db.Exec(query, planCode)
		if err != nil {
			return err
		}
	}

	return nil
}

// getPlanCodesByIDs helper to get plan codes by IDs
func (r *repository) getPlanCodesByIDs(c echo.Context, ids []int) ([]string, error) {
	db := r.getDB(c)
	placeholders := make([]string, len(ids))
	queryArgs := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		queryArgs[i] = id
	}

	query := `SELECT DISTINCT plan_code FROM quotation_service_local.plan_commissions WHERE id IN (` + strings.Join(placeholders, ",") + `)`
	rows, err := db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var planCodes []string
	for rows.Next() {
		var planCode string
		if err := rows.Scan(&planCode); err != nil {
			continue
		}
		planCodes = append(planCodes, planCode)
	}

	return planCodes, rows.Err()
}

// UpdateProductRulesHardcopyAdminFeeByPlanCode updates hardcopy_admin_fee by plan code
func (r *repository) UpdateProductRulesHardcopyAdminFeeByPlanCode(c echo.Context, planCode string, adminFee float64) error {
	db := r.getDB(c)
	query := `
		UPDATE vehicle_service_development.product_rules 
		SET hardcopy_admin_fee = ?, updated_at = CURRENT_TIMESTAMP
		WHERE product_code = ?`
	_, err := db.Exec(query, adminFee, planCode)
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


