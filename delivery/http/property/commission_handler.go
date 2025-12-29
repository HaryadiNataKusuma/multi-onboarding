package http

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

// PropertyCommissionHandler handles HTTP requests for property commissions
type PropertyCommissionHandler struct {
	db *sql.DB
}

// NewPropertyCommissionHandler creates a new property commission handler
func NewPropertyCommissionHandler(db *sql.DB) *PropertyCommissionHandler {
	return &PropertyCommissionHandler{
		db: db,
	}
}

// GenerateCommissionRequest represents the request for generating commissions
type GenerateCommissionRequest struct {
	InsuranceCode        string   `json:"insurance_code"`
	ProductCodes         []string `json:"product_codes"`
	CommissionPercentage float64  `json:"commission_percentage"`
	CommissionVATType    string   `json:"commission_vat_type"`
	AFPercentage         float64  `json:"af_percentage"`
	AFVATType            string   `json:"af_vat_type"`
}

// GenerateCommissions handles POST /api/property/commissions/generate
// This endpoint generates commission data to 3 tables:
// 1. commission_service_development.commissions
// 2. agent_service_development.default_config_products
// 3. quotation_service_development.plan_commissions
func (h *PropertyCommissionHandler) GenerateCommissions(c echo.Context) error {
	var request GenerateCommissionRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if request.InsuranceCode == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "insurance code is required",
		})
	}
	if len(request.ProductCodes) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "at least one product code is required",
		})
	}
	if request.CommissionVATType != "EXCLUSIVE" && request.CommissionVATType != "INCLUSIVE" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "commission VAT type must be EXCLUSIVE or INCLUSIVE",
		})
	}
	if request.AFVATType != "EXCLUSIVE" && request.AFVATType != "INCLUSIVE" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "AF VAT type must be EXCLUSIVE or INCLUSIVE",
		})
	}

	// Validate that all product codes exist in property_service_development.products table
	query := `SELECT code FROM property_service_development.products WHERE code IN (`
	args := make([]interface{}, len(request.ProductCodes))
	for i, code := range request.ProductCodes {
		if i > 0 {
			query += ","
		}
		query += "?"
		args[i] = code
	}
	query += ")"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to validate product codes: " + err.Error(),
		})
	}
	defer rows.Close()

	validProductCodes := make(map[string]bool)
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			continue
		}
		validProductCodes[code] = true
	}

	// Check each product code
	var invalidCodes []string
	for _, productCode := range request.ProductCodes {
		if !validProductCodes[productCode] {
			invalidCodes = append(invalidCodes, productCode)
		}
	}

	if len(invalidCodes) > 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("invalid product codes (not found in property_service_development.products): %v", invalidCodes),
		})
	}

	// Get hardcopy_admin_fee for each product
	productHardcopyFees := make(map[string]int)
	for _, productCode := range request.ProductCodes {
		var hardcopyFee sql.NullInt64
		query := `SELECT hardcopy_admin_fee FROM property_service_development.products WHERE code = ?`
		err := h.db.QueryRow(query, productCode).Scan(&hardcopyFee)
		if err == nil && hardcopyFee.Valid {
			productHardcopyFees[productCode] = int(hardcopyFee.Int64)
		} else {
			// Fallback to admin_fee if hardcopy_admin_fee is null
			var adminFee sql.NullInt64
			fallbackQuery := `SELECT admin_fee FROM property_service_development.products WHERE code = ?`
			if err := h.db.QueryRow(fallbackQuery, productCode).Scan(&adminFee); err == nil && adminFee.Valid {
				productHardcopyFees[productCode] = int(adminFee.Int64)
			} else {
				productHardcopyFees[productCode] = 0
			}
		}
	}

	// Agent levels
	agentLevels := []string{
		"GREEN", "SILVER", "GOLD", "DIAMOND", "PLATINUM", "CORPORATE", "CORPORATE2",
		"CORPORATE3", "CORPORATE4", "TIED", "SKB", "SHOPDRIVE", "MVPARTNERSHIP1", "DIRECTPROPERTY",
	}

	// Start transaction
	tx, err := h.db.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to start transaction: " + err.Error(),
		})
	}
	defer tx.Rollback()

	// For each product code
	for _, productCode := range request.ProductCodes {
		hardcopyFee := productHardcopyFees[productCode]

		// For each agent level
		for _, agentLevel := range agentLevels {
			// 1. Insert into commission_service_development.commissions
			// Set corporate_id: 1 for CORPORATE, CORPORATE2, CORPORATE3, CORPORATE4, otherwise 0
			corporateLevels := []string{"CORPORATE", "CORPORATE2", "CORPORATE3", "CORPORATE4"}
			corporateId := 0
			for _, cl := range corporateLevels {
				if agentLevel == cl {
					corporateId = 1
					break
				}
			}

			query1 := `
				INSERT INTO commission_service_development.commissions 
				(product_code, insurance_code, agent_level, corporate_id, basic_commission, bonus_point, note, start_date, end_date, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, 0, ?, NOW(), NOW(), NOW(), NOW())
				ON DUPLICATE KEY UPDATE 
					basic_commission = VALUES(basic_commission),
					corporate_id = VALUES(corporate_id),
					updated_at = NOW()
			`
			note := "Commission for " + productCode + " - " + agentLevel
			_, err = tx.Exec(query1, productCode, request.InsuranceCode, agentLevel, corporateId, request.CommissionPercentage, note)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to insert commission: " + err.Error(),
				})
			}

			// 2. Insert into agent_service_development.default_config_products
			query2 := `
				INSERT INTO agent_service_development.default_config_products 
				(level, sequence, product, commission, max_discount, insurer, kind, category, selected, renewal_count, point, upline_bonus, bonus, ` + "`order`" + `, created_at, updated_at)
				VALUES (?, 9, ?, ?, ?, ?, 'PL', 'PR', 1, 0, 0, 0, 0, 1000000, NOW(), NOW())
				ON DUPLICATE KEY UPDATE 
					commission = VALUES(commission),
					max_discount = VALUES(max_discount),
					kind = VALUES(kind),
					category = VALUES(category),
					updated_at = NOW()
			`
			_, err = tx.Exec(query2, agentLevel, productCode, request.CommissionPercentage, request.CommissionPercentage, request.InsuranceCode)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to insert default_config: " + err.Error(),
				})
			}
		}

		// 3. Insert into quotation_service_development.plan_commissions (one row per product, not per agent level)
		query3 := `
			INSERT INTO quotation_service_development.plan_commissions 
			(plan_code, product_id, commission_percentage, commission_vat_type, af_percentage, af_vat_type, hardcopy_fee, is_active, version, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, 1, 1, NOW())
			ON DUPLICATE KEY UPDATE 
				commission_percentage = VALUES(commission_percentage),
				commission_vat_type = VALUES(commission_vat_type),
				af_percentage = VALUES(af_percentage),
				af_vat_type = VALUES(af_vat_type),
				hardcopy_fee = VALUES(hardcopy_fee)
		`
		_, err = tx.Exec(query3, productCode, 5, request.CommissionPercentage, request.CommissionVATType, request.AFPercentage, request.AFVATType, hardcopyFee)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to insert plan_commission: " + err.Error(),
			})
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to commit transaction: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Commissions generated successfully",
	})
}

// CheckDuplicateProductCodes handles POST /api/property/commissions/check-duplicates
func (h *PropertyCommissionHandler) CheckDuplicateProductCodes(c echo.Context) error {
	var request struct {
		ProductCodes []string `json:"product_codes"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if len(request.ProductCodes) == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"has_duplicates": false,
			"duplicates":     []string{},
		})
	}

	// Check for duplicates in quotation_service_development.plan_commissions
	// Only check for product codes that exist in property_service_development.products
	query := `
		SELECT DISTINCT pc.plan_code 
		FROM quotation_service_development.plan_commissions pc
		INNER JOIN property_service_development.products p ON pc.plan_code = p.code
		WHERE pc.plan_code IN (`

	args := make([]interface{}, len(request.ProductCodes))
	for i, code := range request.ProductCodes {
		if i > 0 {
			query += ","
		}
		query += "?"
		args[i] = code
	}
	query += ") AND pc.is_active = 1"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to check duplicates: " + err.Error(),
		})
	}
	defer rows.Close()

	var duplicates []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			continue
		}
		duplicates = append(duplicates, code)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"has_duplicates": len(duplicates) > 0,
		"duplicates":     duplicates,
	})
}

// GetCommissionsFromTable handles GET /api/property/commissions/table/commissions
func (h *PropertyCommissionHandler) GetCommissionsFromTable(c echo.Context) error {
	insuranceCode := c.QueryParam("insurance_code")

	query := `
		SELECT c.id, c.product_code, c.insurance_code, c.agent_level, c.corporate_id, 
		       c.basic_commission, c.bonus_point, c.note, c.start_date, c.end_date, 
		       c.created_at, c.updated_at
		FROM commission_service_development.commissions c
		INNER JOIN property_service_development.products p ON c.product_code = p.code
		WHERE 1=1`

	var args []interface{}
	if insuranceCode != "" {
		query += " AND c.insurance_code = ?"
		args = append(args, insuranceCode)
	}
	query += " ORDER BY c.id DESC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch commissions: " + err.Error(),
		})
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id, corporateId int
		var productCode, insuranceCode, agentLevel, note string
		var basicCommission, bonusPoint float64
		var startDate, endDate, createdAt, updatedAt string

		err := rows.Scan(&id, &productCode, &insuranceCode, &agentLevel, &corporateId,
			&basicCommission, &bonusPoint, &note, &startDate, &endDate, &createdAt, &updatedAt)
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"id":               id,
			"product_code":     productCode,
			"insurance_code":   insuranceCode,
			"agent_level":      agentLevel,
			"corporate_id":     corporateId,
			"basic_commission": basicCommission,
			"bonus_point":      bonusPoint,
			"note":             note,
			"start_date":       startDate,
			"end_date":         endDate,
			"created_at":       createdAt,
			"updated_at":       updatedAt,
		})
	}

	return c.JSON(http.StatusOK, results)
}

// GetPlanCommissionsFromTable handles GET /api/property/commissions/table/plan_commissions
func (h *PropertyCommissionHandler) GetPlanCommissionsFromTable(c echo.Context) error {
	insuranceCode := c.QueryParam("insurance_code")

	query := `
		SELECT pc.id, pc.plan_code, pc.product_id, pc.insurer_id, 
		       pc.commission_percentage, pc.commission_vat_type, 
		       pc.af_percentage, pc.af_vat_type, 
		       pc.hardcopy_fee, pc.is_active, pc.version, pc.created_at
		FROM quotation_service_development.plan_commissions pc
		INNER JOIN property_service_development.products p ON pc.plan_code = p.code
		WHERE 1=1`

	var args []interface{}
	if insuranceCode != "" {
		// Extract insurance code from product_code (format: PR-{INSURANCE_CODE}-... or PR-HOME-{INSURANCE_CODE}-...)
		query += " AND (p.code LIKE ? OR p.code LIKE ?)"
		args = append(args, "PR-"+insuranceCode+"-%", "PR-HOME-"+insuranceCode+"-%")
	}
	query += " ORDER BY pc.id DESC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch plan commissions: " + err.Error(),
		})
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id, productId int
		var planCode, commissionVATType, afVATType string
		var insurerId sql.NullInt64
		var commissionPercentage, afPercentage, hardcopyFee float64
		var isActive int
		var version int
		var createdAt string

		err := rows.Scan(&id, &planCode, &productId, &insurerId,
			&commissionPercentage, &commissionVATType,
			&afPercentage, &afVATType,
			&hardcopyFee, &isActive, &version, &createdAt)
		if err != nil {
			continue
		}

		result := map[string]interface{}{
			"id":                    id,
			"plan_code":             planCode,
			"product_code":          planCode, // For compatibility
			"product_id":            productId,
			"commission_percentage": commissionPercentage,
			"commission_vat_type":   commissionVATType,
			"af_percentage":         afPercentage,
			"af_vat_type":           afVATType,
			"hardcopy_fee":          hardcopyFee,
			"is_active":             isActive,
			"version":               version,
			"created_at":            createdAt,
		}

		if insurerId.Valid {
			result["insurer_id"] = insurerId.Int64
		} else {
			result["insurer_id"] = nil
		}

		results = append(results, result)
	}

	return c.JSON(http.StatusOK, results)
}

// GetDefaultConfigFromTable handles GET /api/property/commissions/table/default_config_products
func (h *PropertyCommissionHandler) GetDefaultConfigFromTable(c echo.Context) error {
	insuranceCode := c.QueryParam("insurance_code")

	// Get all property product codes first
	productCodesQuery := `SELECT code FROM property_service_development.products`
	productRows, err := h.db.Query(productCodesQuery)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch product codes: " + err.Error(),
		})
	}
	defer productRows.Close()

	propertyProductCodes := make(map[string]bool)
	for productRows.Next() {
		var code string
		if err := productRows.Scan(&code); err != nil {
			continue
		}
		propertyProductCodes[code] = true
	}

	if len(propertyProductCodes) == 0 {
		return c.JSON(http.StatusOK, []interface{}{})
	}

	// Build query with IN clause for product codes
	query := `
		SELECT dcp.id, dcp.level, dcp.sequence, dcp.kind, dcp.category, 
		       dcp.insurer, dcp.product, dcp.payment_frequency, dcp.policy_year, 
		       dcp.selected, dcp.renewal_count, dcp.commission, dcp.point, 
		       dcp.upline_bonus, dcp.bonus, dcp.max_discount, dcp.` + "`order`" + `, dcp.created_at
		FROM agent_service_development.default_config_products dcp
		WHERE dcp.product IN (`

	args := make([]interface{}, 0)
	productCodesList := make([]string, 0)
	for code := range propertyProductCodes {
		productCodesList = append(productCodesList, code)
	}

	for i, code := range productCodesList {
		if i > 0 {
			query += ","
		}
		query += "?"
		args = append(args, code)
	}
	query += ")"

	if insuranceCode != "" {
		query += " AND dcp.insurer = ?"
		args = append(args, insuranceCode)
	}
	query += " ORDER BY dcp.product ASC, dcp.level ASC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch default config: " + err.Error(),
		})
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id, sequence, selected, renewalCount, order int
		var level, kind, category, insurer, product string
		var paymentFrequency, policyYear sql.NullString
		var commission, point, uplineBonus, bonus, maxDiscount float64
		var createdAt string

		err := rows.Scan(&id, &level, &sequence, &kind, &category,
			&insurer, &product, &paymentFrequency, &policyYear,
			&selected, &renewalCount, &commission, &point,
			&uplineBonus, &bonus, &maxDiscount, &order, &createdAt)
		if err != nil {
			continue
		}

		result := map[string]interface{}{
			"id":               id,
			"level":            level,
			"sequence":         sequence,
			"kind":             kind,
			"category":         category,
			"insurer":          insurer,
			"product":          product,
			"product_code":     product, // For compatibility
			"commission":       commission,
			"basic_commission": commission, // For compatibility
			"point":            point,
			"upline_bonus":     uplineBonus,
			"bonus":            bonus,
			"max_discount":     maxDiscount,
			"order":            order,
			"created_at":       createdAt,
		}

		if paymentFrequency.Valid {
			result["payment_frequency"] = paymentFrequency.String
		} else {
			result["payment_frequency"] = nil
		}

		if policyYear.Valid {
			result["policy_year"] = policyYear.String
		} else {
			result["policy_year"] = nil
		}

		results = append(results, result)
	}

	return c.JSON(http.StatusOK, results)
}


