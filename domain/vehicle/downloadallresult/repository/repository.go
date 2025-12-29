package repository

import (
	"database/sql"
	"log"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/downloadallresult"
	"multi-onboarding/utils"
)

type repository struct{}

// NewDownloadAllResultRepository creates a new download all result repository
func NewDownloadAllResultRepository() downloadallresult.Repository {
	return &repository{}
}

// getDB gets the database connection
func (r *repository) getDB(c echo.Context) *sql.DB {
	return utils.GetDB()
}

// GetInsurances gets insurances data
func (r *repository) GetInsurances(c echo.Context, insuranceCode string) ([]downloadallresult.InsuranceRow, error) {
	db := r.getDB(c)
	query := `SELECT id, code, name, status, logo, created_at, updated_at 
		FROM vehicle_service_development.insurances 
		WHERE code = ? AND status = 1`

	rows, err := db.Query(query, insuranceCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []downloadallresult.InsuranceRow
	for rows.Next() {
		var row downloadallresult.InsuranceRow
		var createdAt, updatedAt sql.NullString

		if err := rows.Scan(&row.ID, &row.Code, &row.Name, &row.Status, &row.Logo, &createdAt, &updatedAt); err != nil {
			continue
		}

		if createdAt.Valid {
			row.CreatedAt = createdAt.String
		}
		if updatedAt.Valid {
			row.UpdatedAt = updatedAt.String
		}

		results = append(results, row)
	}

	return results, rows.Err()
}

// GetProducts gets products data
func (r *repository) GetProducts(c echo.Context, insuranceCode string) ([]downloadallresult.ProductRow, error) {
	db := r.getDB(c)
	query := `SELECT id, code, insurance_code, insurance_product_code, vehicle_type, 
		is_insurance_api, name, logo, is_active, is_insurer_driven, 
		is_personal_usage, is_electric_vehicle, summary, config, 
		insurance_detail, protection_detail, how_to_claim, created_by, base_product,
		created_at, updated_at 
		FROM vehicle_service_development.products 
		WHERE insurance_code = ? AND is_active = 1
		ORDER BY id`

	rows, err := db.Query(query, insuranceCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []downloadallresult.ProductRow
	for rows.Next() {
		var row downloadallresult.ProductRow
		var logo, summary, config, insuranceDetail, protectionDetail, howToClaim sql.NullString
		var createdAt, updatedAt sql.NullString

		if err := rows.Scan(&row.ID, &row.Code, &row.InsuranceCode, &row.InsuranceProductCode, &row.VehicleType,
			&row.IsInsuranceAPI, &row.Name, &logo, &row.IsActive, &row.IsInsurerDriven,
			&row.IsPersonalUsage, &row.IsElectricVehicle, &summary, &config,
			&insuranceDetail, &protectionDetail, &howToClaim, &row.CreatedBy, &row.BaseProduct,
			&createdAt, &updatedAt); err != nil {
			log.Printf("Error scanning product row: %v", err)
			continue
		}

		if logo.Valid {
			row.Logo = logo.String
		}
		if summary.Valid {
			row.Summary = summary.String
		}
		if config.Valid {
			row.Config = config.String
		}
		if insuranceDetail.Valid {
			row.InsuranceDetail = insuranceDetail.String
		}
		if protectionDetail.Valid {
			row.ProtectionDetail = protectionDetail.String
		}
		if howToClaim.Valid {
			row.HowToClaim = howToClaim.String
		}
		if createdAt.Valid {
			row.CreatedAt = createdAt.String
		}
		if updatedAt.Valid {
			row.UpdatedAt = updatedAt.String
		}

		results = append(results, row)
	}

	return results, rows.Err()
}

// GetTemplates gets templates data
func (r *repository) GetTemplates(c echo.Context, insuranceCode string) ([]downloadallresult.TemplateRow, error) {
	db := r.getDB(c)
	query := `SELECT id, locale, template_id, value, insurance_code, created_by, created_at 
		FROM vehicle_service_development.templates 
		WHERE insurance_code = ?`

	rows, err := db.Query(query, insuranceCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []downloadallresult.TemplateRow
	for rows.Next() {
		var row downloadallresult.TemplateRow
		var createdAt sql.NullString

		if err := rows.Scan(&row.ID, &row.Locale, &row.TemplateID, &row.Value, &row.InsuranceCode, &row.CreatedBy, &createdAt); err != nil {
			continue
		}

		if createdAt.Valid {
			row.CreatedAt = createdAt.String
		}

		results = append(results, row)
	}

	return results, rows.Err()
}

// GetCommissions gets commissions data
func (r *repository) GetCommissions(c echo.Context, insuranceCode string) ([]downloadallresult.CommissionRow, error) {
	db := r.getDB(c)
	query := `SELECT id, product_code, insurance_code, agent_level, corporate_id, 
		basic_commission, bonus_point, note, start_date, end_date, created_at, updated_at 
		FROM vehicle_service_development.commissions 
		WHERE insurance_code = ?`

	rows, err := db.Query(query, insuranceCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []downloadallresult.CommissionRow
	for rows.Next() {
		var row downloadallresult.CommissionRow
		var startDate, endDate, createdAt, updatedAt sql.NullString

		if err := rows.Scan(&row.ID, &row.ProductCode, &row.InsuranceCode, &row.AgentLevel, &row.CorporateID,
			&row.BasicCommission, &row.BonusPoint, &row.Note, &startDate, &endDate, &createdAt, &updatedAt); err != nil {
			continue
		}

		if startDate.Valid {
			row.StartDate = startDate.String
		}
		if endDate.Valid {
			row.EndDate = endDate.String
		}
		if createdAt.Valid {
			row.CreatedAt = createdAt.String
		}
		if updatedAt.Valid {
			row.UpdatedAt = updatedAt.String
		}

		results = append(results, row)
	}

	return results, rows.Err()
}

// GetPlanCommissions gets plan commissions data
func (r *repository) GetPlanCommissions(c echo.Context, insuranceCode string) ([]downloadallresult.PlanCommissionRow, error) {
	db := r.getDB(c)
	query := `SELECT id, plan_code, product_id, insurer_id, commission_percentage, commission_vat_type,
		af_percentage, af_vat_type, admin_fee, hardcopy_fee, is_active, version, created_at 
		FROM vehicle_service_development.plan_commissions 
		WHERE is_active = 1 AND plan_code LIKE ?`

	rows, err := db.Query(query, "%"+insuranceCode+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []downloadallresult.PlanCommissionRow
	for rows.Next() {
		var row downloadallresult.PlanCommissionRow
		var createdAt sql.NullString

		if err := rows.Scan(&row.ID, &row.PlanCode, &row.ProductID, &row.InsurerID, &row.CommissionPercentage, &row.CommissionVatType,
			&row.AFPercentage, &row.AFVatType, &row.AdminFee, &row.HardcopyFee, &row.IsActive, &row.Version, &createdAt); err != nil {
			continue
		}

		if createdAt.Valid {
			row.CreatedAt = createdAt.String
		}

		results = append(results, row)
	}

	return results, rows.Err()
}

// GetDefaultConfigProducts gets default config products data
func (r *repository) GetDefaultConfigProducts(c echo.Context, insuranceCode string) ([]downloadallresult.DefaultConfigProductRow, error) {
	db := r.getDB(c)
	query := `SELECT id, level, sequence, kind, category, insurer, product, selected, commission, 
		point, upline_bonus, bonus, max_discount, ` + "`order`" + `, created_at, updated_at 
		FROM vehicle_service_development.default_config_products 
		WHERE insurer = ?`

	rows, err := db.Query(query, insuranceCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []downloadallresult.DefaultConfigProductRow
	for rows.Next() {
		var row downloadallresult.DefaultConfigProductRow
		var createdAt, updatedAt sql.NullString

		if err := rows.Scan(&row.ID, &row.Level, &row.Sequence, &row.Kind, &row.Category, &row.Insurer, &row.Product, &row.Selected,
			&row.Commission, &row.Point, &row.UplineBonus, &row.Bonus, &row.MaxDiscount, &row.Order, &createdAt, &updatedAt); err != nil {
			log.Printf("Error scanning default config products row: %v", err)
			continue
		}

		if createdAt.Valid {
			row.CreatedAt = createdAt.String
		}
		if updatedAt.Valid {
			row.UpdatedAt = updatedAt.String
		}

		results = append(results, row)
	}

	return results, rows.Err()
}

// GetProductRules gets product rules data
func (r *repository) GetProductRules(c echo.Context, insuranceCode string) ([]downloadallresult.ProductRuleRow, error) {
	db := r.getDB(c)
	query := `SELECT product_code, insurance_code, vehicle_type, vehicle_category,
		start_vehicle_value, end_vehicle_value, limit_vehicle_age, admin_fee, hardcopy_admin_fee,
		region_id, base_premium_value, loading_fee_premium_value, commercial_usage_value,
		start_loading_age, additional_premium, type_additional_premium, rules, created_by,
		created_at, base_premium_type, base_loading_premium_value
		FROM vehicle_service_development.product_rules 
		WHERE insurance_code = ?`

	rows, err := db.Query(query, insuranceCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []downloadallresult.ProductRuleRow
	for rows.Next() {
		var row downloadallresult.ProductRuleRow
		var rules, createdAt sql.NullString

		if err := rows.Scan(&row.ProductCode, &row.InsuranceCode, &row.VehicleType, &row.VehicleCategory,
			&row.StartVehicleValue, &row.EndVehicleValue, &row.LimitVehicleAge, &row.AdminFee, &row.HardcopyAdminFee,
			&row.RegionID, &row.BasePremiumValue, &row.LoadingFeePremiumValue, &row.CommercialUsageValue,
			&row.StartLoadingAge, &row.AdditionalPremium, &row.TypeAdditionalPremium, &rules, &row.CreatedBy,
			&createdAt, &row.BasePremiumType, &row.BaseLoadingPremiumValue); err != nil {
			continue
		}

		if rules.Valid {
			row.Rules = rules.String
		}
		if createdAt.Valid {
			row.CreatedAt = createdAt.String
		}

		results = append(results, row)
	}

	return results, rows.Err()
}

// GetAddons gets addons data
func (r *repository) GetAddons(c echo.Context) ([]downloadallresult.AddonRow, error) {
	db := r.getDB(c)
	query := `SELECT a1.id, a1.code, a1.name, a1.is_active, a1.created_at
		FROM vehicle_service_development.addons a1
		INNER JOIN (
			SELECT code, MAX(created_at) as max_created_at
			FROM vehicle_service_development.addons
			WHERE is_active = 1
			GROUP BY code
		) a2 ON a1.code = a2.code AND a1.created_at = a2.max_created_at
		WHERE a1.is_active = 1
		ORDER BY a1.code`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []downloadallresult.AddonRow
	for rows.Next() {
		var row downloadallresult.AddonRow
		var createdAt sql.NullString

		if err := rows.Scan(&row.ID, &row.Code, &row.Name, &row.IsActive, &createdAt); err != nil {
			continue
		}

		if createdAt.Valid {
			row.CreatedAt = createdAt.String
		}

		results = append(results, row)
	}

	return results, rows.Err()
}

// GetAddonRules gets addon rules data
func (r *repository) GetAddonRules(c echo.Context, insuranceCode string) ([]downloadallresult.AddonRuleRow, error) {
	db := r.getDB(c)
	query := `SELECT addon_code, product_code, insurance_code, rules, created_by 
		FROM vehicle_service_development.addon_rules 
		WHERE insurance_code = ?`

	rows, err := db.Query(query, insuranceCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []downloadallresult.AddonRuleRow
	for rows.Next() {
		var row downloadallresult.AddonRuleRow
		var rules sql.NullString

		if err := rows.Scan(&row.AddonCode, &row.ProductCode, &row.InsuranceCode, &rules, &row.CreatedBy); err != nil {
			continue
		}

		if rules.Valid {
			row.Rules = rules.String
		}

		results = append(results, row)
	}

	return results, rows.Err()
}

// GetInsuranceOwnRisks gets insurance own risks data
func (r *repository) GetInsuranceOwnRisks(c echo.Context, insuranceCode string) ([]downloadallresult.InsuranceOwnRiskRow, error) {
	db := r.getDB(c)
	query := `SELECT insurance_code, product_type, code, title, value, value_type, 
		description, is_mandatory, is_active, is_electric_vehicle, created_by
		FROM vehicle_service_development.insurance_own_risks 
		WHERE insurance_code = ? AND is_active = 1`

	rows, err := db.Query(query, insuranceCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []downloadallresult.InsuranceOwnRiskRow
	for rows.Next() {
		var row downloadallresult.InsuranceOwnRiskRow
		var description sql.NullString
		var isMandatory, isActive bool
		var isElectricVehicle sql.NullInt64

		if err := rows.Scan(&row.InsuranceCode, &row.ProductType, &row.Code, &row.Title, &row.Value, &row.ValueType,
			&description, &isMandatory, &isActive, &isElectricVehicle, &row.CreatedBy); err != nil {
			continue
		}

		if description.Valid {
			row.Description = description.String
		}
		if isElectricVehicle.Valid {
			row.IsElectricVehicle = int(isElectricVehicle.Int64)
		}

		// Convert bool to int
		if isMandatory {
			row.IsMandatory = 1
		} else {
			row.IsMandatory = 0
		}
		if isActive {
			row.IsActive = 1
		} else {
			row.IsActive = 0
		}

		results = append(results, row)
	}

	return results, rows.Err()
}

// GetInsuranceClauses gets insurance clauses data
func (r *repository) GetInsuranceClauses(c echo.Context, insuranceCode string) ([]downloadallresult.InsuranceClauseRow, error) {
	db := r.getDB(c)
	query := `SELECT insurance_code, product_type, code, title, content, 
		description, is_mandatory, is_active, created_by
		FROM vehicle_service_development.insurance_clauses 
		WHERE insurance_code = ? AND is_active = 1`

	rows, err := db.Query(query, insuranceCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []downloadallresult.InsuranceClauseRow
	for rows.Next() {
		var row downloadallresult.InsuranceClauseRow
		var description sql.NullString
		var isMandatory, isActive bool

		if err := rows.Scan(&row.InsuranceCode, &row.ProductType, &row.Code, &row.Title, &row.Content,
			&description, &isMandatory, &isActive, &row.CreatedBy); err != nil {
			continue
		}

		if description.Valid {
			row.Description = description.String
		}

		// Convert bool to int
		if isMandatory {
			row.IsMandatory = 1
		} else {
			row.IsMandatory = 0
		}
		if isActive {
			row.IsActive = 1
		} else {
			row.IsActive = 0
		}

		results = append(results, row)
	}

	return results, rows.Err()
}


