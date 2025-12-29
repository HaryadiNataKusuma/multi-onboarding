package http

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/labstack/echo"
)

// PropertyAddonRuleDraft represents a draft addon rule for property
type PropertyAddonRuleDraft struct {
	ProductCode         string `json:"product_code"`
	AddonCode           string `json:"addon_code"`
	ConstructionClassID string `json:"construction_class_id"`
	OccupationCode      string `json:"occupation_code"`
	Rules               string `json:"rules"`
}

// PropertyAddonRuleDetailHandler handles HTTP requests for property addon rule details
type PropertyAddonRuleDetailHandler struct {
	db        *sql.DB
	draftData []PropertyAddonRuleDraft
	draftMu   sync.Mutex
}

// NewPropertyAddonRuleDetailHandler creates a new property addon rule detail handler
func NewPropertyAddonRuleDetailHandler(db *sql.DB) *PropertyAddonRuleDetailHandler {
	return &PropertyAddonRuleDetailHandler{db: db}
}

// escapeCSV escapes a string for CSV format (double quotes inside are escaped with double quotes)
func escapeCSV(s string) string {
	// Replace " with "" and wrap in quotes
	escaped := strings.ReplaceAll(s, `"`, `""`)
	return `"` + escaped + `"`
}

// extractInsuranceCode extracts insurance code from product code
// Format: PR-{INSURANCE_CODE}-... or PR-HOME-{INSURANCE_CODE}-...
func extractInsuranceCode(productCode string) string {
	parts := strings.Split(productCode, "-")
	if len(parts) >= 2 {
		// Check if format is PR-HOME-{INSURANCE_CODE}-...
		if len(parts) >= 3 && parts[1] == "HOME" {
			return parts[2]
		}
		// Format is PR-{INSURANCE_CODE}-...
		return parts[1]
	}
	return ""
}

// DownloadTemplate handles GET /api/property/addon-rule-details/template
func (h *PropertyAddonRuleDetailHandler) DownloadTemplate(c echo.Context) error {
	// Fetch product codes from property_service_development.products
	productQuery := `SELECT code FROM property_service_development.products ORDER BY code ASC`
	productRows, err := h.db.Query(productQuery)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch products: " + err.Error(),
		})
	}
	defer productRows.Close()

	var productCodes []string
	for productRows.Next() {
		var code string
		if err := productRows.Scan(&code); err != nil {
			continue
		}
		productCodes = append(productCodes, code)
	}

	// Fetch occupation codes from property_service_development.product_rules (only for EQVET)
	occupationQuery := `SELECT DISTINCT occupation_code FROM property_service_development.product_rules ORDER BY occupation_code ASC`
	occupationRows, err := h.db.Query(occupationQuery)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch occupation codes: " + err.Error(),
		})
	}
	defer occupationRows.Close()

	var occupationCodes []string
	for occupationRows.Next() {
		var code string
		if err := occupationRows.Scan(&code); err != nil {
			continue
		}
		occupationCodes = append(occupationCodes, code)
	}

	// Define zones for EQVET (1-5)
	eqvetZones := []int{1, 2, 3, 4, 5}

	// Define storey ranges
	storeyRanges := []struct {
		Start int
		End   int
	}{
		{1, 9},
		{10, 100},
	}

	// Build CSV template
	var template strings.Builder

	// Write header (added RULES column)
	template.WriteString("product-code,addon-code,construction-class-id,occupation-code,zone,province-id,addon-premi,start-storey,end-storey,max addon-premi,base-premi-type,additional-premi,additional-premi-type,Value,Auto Selected,min_tsi,RULES\n")

	// Function to build JSON rules directly (instead of Excel formula)
	buildRulesJSON := func(occupationCode string, zone int, provinceIdExclude string, provinceIdInclude string, startStorey int, endStorey int, addonPremi float64, basePremiType string, maxAddonPremi float64, autoSelected string, minTsi int) string {
		selectedPart := ""
		if autoSelected == "Y" {
			selectedPart = `,"selected":true`
		}
		rules := fmt.Sprintf(`{"occupation_code_exclude":[],"ZONE":%d,"province_id_exclude":%s,"province_id_include":%s,"start_storey":%d,"end_storey":%d,"base_addon_premium":%.2f,"base_addon_premium_type":"%s","max_addon_premium":%.2f%s,"min_tsi":%d}`,
			zone, provinceIdExclude, provinceIdInclude, startStorey, endStorey, addonPremi, strings.ToLower(basePremiType), maxAddonPremi, selectedPart, minTsi)
		return rules
	}

	rowNumber := 2 // Start from row 2 (row 1 is header)

	// Generate rows for each product
	for _, productCode := range productCodes {
		// Extract insurance code from product code
		insuranceCode := extractInsuranceCode(productCode)

		// 1. EQVET - uses occupation codes from product_rules and zones 1-5, 2 storey ranges
		for _, occupationCode := range occupationCodes {
			for _, zone := range eqvetZones {
				for _, storeyRange := range storeyRanges {
					provinceId := `"province_id_exclude":[],"province_id_include":[]`
					rulesJSON := buildRulesJSON(occupationCode, zone, "[]", "[]", storeyRange.Start, storeyRange.End, 0.1, "percentage", 0, "N", -1)
					row := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%d,%d,%s,%s,%s,%s,%s,%s,%s,%s\n",
						productCode,
						"EQVET",
						"1",
						occupationCode,
						escapeCSV(fmt.Sprintf(`"ZONE":%d`, zone)),
						escapeCSV(provinceId),
						"0.1%",
						storeyRange.Start,
						storeyRange.End,
						"0%",
						"PERCENTAGE",
						"NULL",
						"NULL",
						"NULL",
						"N",
						"-1",
						escapeCSV(rulesJSON),
					)
					template.WriteString(row)
					rowNumber++
				}
			}
		}

		// 2. FWTWD - 2 rows only (different province-id), occupation_code = -1, zone = 0, start/end storey = 0
		{
			// Row 1: province_id_include:[3,6,9]
			provinceId1 := `"province_id_exclude":[],"province_id_include":[3,6,9]`
			rulesJSON1 := buildRulesJSON("-1", 0, "[]", "[3,6,9]", 0, 0, 0.1, "percentage", 0, "N", -1)
			row1 := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%d,%d,%s,%s,%s,%s,%s,%s,%s,%s\n",
				productCode,
				"FWTWD",
				"1",
				"-1",
				escapeCSV(`"ZONE":0`),
				escapeCSV(provinceId1),
				"0.1%",
				0,
				0,
				"0%",
				"PERCENTAGE",
				"NULL",
				"NULL",
				"NULL",
				"N",
				"-1",
				escapeCSV(rulesJSON1),
			)
			template.WriteString(row1)
			rowNumber++

			// Row 2: province_id_exclude:[3,6,9]
			provinceId2 := `"province_id_exclude":[3,6,9],"province_id_include":[]`
			rulesJSON2 := buildRulesJSON("-1", 0, "[3,6,9]", "[]", 0, 0, 0.1, "percentage", 0, "N", -1)
			row2 := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%d,%d,%s,%s,%s,%s,%s,%s,%s,%s\n",
				productCode,
				"FWTWD",
				"1",
				"-1",
				escapeCSV(`"ZONE":0`),
				escapeCSV(provinceId2),
				"0.1%",
				0,
				0,
				"0%",
				"PERCENTAGE",
				"NULL",
				"NULL",
				"NULL",
				"N",
				"-1",
				escapeCSV(rulesJSON2),
			)
			template.WriteString(row2)
			rowNumber++
		}

		// 3. OTHERS - only 1 row, occupation_code = -1, zone = 0, start/end storey = 0
		{
			provinceId := `"province_id_exclude":[],"province_id_include":[]`
			rulesJSON := buildRulesJSON("-1", 0, "[]", "[]", 0, 0, 0.1, "percentage", 0, "N", -1)
			row := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%d,%d,%s,%s,%s,%s,%s,%s,%s,%s\n",
				productCode,
				"OTHERS",
				"1",
				"-1",
				escapeCSV(`"ZONE":0`),
				escapeCSV(provinceId),
				"0.1%",
				0,
				0,
				"0%",
				"PERCENTAGE",
				"NULL",
				"NULL",
				"NULL",
				"N",
				"-1",
				escapeCSV(rulesJSON),
			)
			template.WriteString(row)
			rowNumber++
		}

		// 4. PA_{insurance_code} - only 1 row, occupation_code = -1, zone = 0, Auto Selected = Y, start/end storey = 0
		{
			addonCode := fmt.Sprintf("PA_%s", insuranceCode)
			provinceId := `"province_id_exclude":[],"province_id_include":[]`
			rulesJSON := buildRulesJSON("-1", 0, "[]", "[]", 0, 0, 0.1, "percentage", 0, "Y", -1)
			row := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%d,%d,%s,%s,%s,%s,%s,%s,%s,%s\n",
				productCode,
				addonCode,
				"1",
				"-1",
				escapeCSV(`"ZONE":0`),
				escapeCSV(provinceId),
				"0.1%",
				0,
				0,
				"0%",
				"PERCENTAGE",
				"NULL",
				"NULL",
				"NULL",
				"Y",
				"-1",
				escapeCSV(rulesJSON),
			)
			template.WriteString(row)
			rowNumber++
		}

		// 5. RSMDCC - only 1 row, occupation_code = -1, zone = 0, start/end storey = 0
		{
			provinceId := `"province_id_exclude":[],"province_id_include":[]`
			rulesJSON := buildRulesJSON("-1", 0, "[]", "[]", 0, 0, 0.1, "percentage", 0, "N", -1)
			row := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%d,%d,%s,%s,%s,%s,%s,%s,%s,%s\n",
				productCode,
				"RSMDCC",
				"1",
				"-1",
				escapeCSV(`"ZONE":0`),
				escapeCSV(provinceId),
				"0.1%",
				0,
				0,
				"0%",
				"PERCENTAGE",
				"NULL",
				"NULL",
				"NULL",
				"N",
				"-1",
				escapeCSV(rulesJSON),
			)
			template.WriteString(row)
			rowNumber++
		}

		// 6. TAC_{insurance_code} - only 1 row, occupation_code = -1, zone = 0, Auto Selected = Y, start/end storey = 0
		{
			addonCode := fmt.Sprintf("TAC_%s", insuranceCode)
			provinceId := `"province_id_exclude":[],"province_id_include":[]`
			rulesJSON := buildRulesJSON("-1", 0, "[]", "[]", 0, 0, 0.1, "percentage", 0, "Y", -1)
			row := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%d,%d,%s,%s,%s,%s,%s,%s,%s,%s\n",
				productCode,
				addonCode,
				"1",
				"-1",
				escapeCSV(`"ZONE":0`),
				escapeCSV(provinceId),
				"0.1%",
				0,
				0,
				"0%",
				"PERCENTAGE",
				"NULL",
				"NULL",
				"NULL",
				"Y",
				"-1",
				escapeCSV(rulesJSON),
			)
			template.WriteString(row)
			rowNumber++
		}

		// 7. TPL_{insurance_code} - only 1 row, occupation_code = -1, zone = 0, Auto Selected = Y, start/end storey = 0
		{
			addonCode := fmt.Sprintf("TPL_%s", insuranceCode)
			provinceId := `"province_id_exclude":[],"province_id_include":[]`
			rulesJSON := buildRulesJSON("-1", 0, "[]", "[]", 0, 0, 0.1, "percentage", 0, "Y", -1)
			row := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%d,%d,%s,%s,%s,%s,%s,%s,%s,%s\n",
				productCode,
				addonCode,
				"1",
				"-1",
				escapeCSV(`"ZONE":0`),
				escapeCSV(provinceId),
				"0.1%",
				0,
				0,
				"0%",
				"PERCENTAGE",
				"NULL",
				"NULL",
				"NULL",
				"Y",
				"-1",
				escapeCSV(rulesJSON),
			)
			template.WriteString(row)
			rowNumber++
		}
	}

	// Set response headers for CSV download
	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=property_addon_rule_details_template.csv")

	return c.String(http.StatusOK, template.String())
}

// Upload handles POST /api/property/addon-rule-details/upload
func (h *PropertyAddonRuleDetailHandler) Upload(c echo.Context) error {
	// Parse multipart form
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "File is required: " + err.Error(),
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to open file: " + err.Error(),
		})
	}
	defer src.Close()

	// Parse CSV
	reader := csv.NewReader(src)
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read header
	header, err := reader.Read()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to read header: " + err.Error(),
		})
	}

	// Create header map (lowercase and trimmed)
	headerMap := make(map[string]int)
	for i, h := range header {
		key := strings.TrimSpace(strings.ToLower(h))
		headerMap[key] = i
	}

	// Helper function to get value from row
	getValue := func(row []string, key string) string {
		if idx, ok := headerMap[key]; ok && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	// Parse data rows
	var drafts []PropertyAddonRuleDraft
	rowNum := 0
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		rowNum++
		if err != nil {
			// Log the error but continue
			fmt.Printf("CSV parse error at row %d: %v\n", rowNum, err)
			continue
		}

		draft := PropertyAddonRuleDraft{
			ProductCode:         getValue(row, "product-code"),
			AddonCode:           getValue(row, "addon-code"),
			ConstructionClassID: getValue(row, "construction-class-id"),
			OccupationCode:      getValue(row, "occupation-code"),
			Rules:               getValue(row, "rules"),
		}

		// Skip empty rows
		if draft.ProductCode == "" && draft.AddonCode == "" {
			continue
		}

		drafts = append(drafts, draft)
	}

	// Store drafts
	h.draftMu.Lock()
	h.draftData = drafts
	h.draftMu.Unlock()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "Success",
		"data":   drafts,
		"count":  len(drafts),
	})
}

// GetAllDrafts handles GET /api/property/addon-rule-details/draft
func (h *PropertyAddonRuleDetailHandler) GetAllDrafts(c echo.Context) error {
	h.draftMu.Lock()
	defer h.draftMu.Unlock()

	// Return empty array instead of null if no drafts
	if h.draftData == nil {
		return c.JSON(http.StatusOK, []PropertyAddonRuleDraft{})
	}
	return c.JSON(http.StatusOK, h.draftData)
}

// ClearDrafts handles POST /api/property/addon-rule-details/draft/clear
func (h *PropertyAddonRuleDetailHandler) ClearDrafts(c echo.Context) error {
	h.draftMu.Lock()
	h.draftData = nil
	h.draftMu.Unlock()

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Drafts cleared successfully",
	})
}

// ConfirmDrafts handles POST /api/property/addon-rule-details/draft/confirm
func (h *PropertyAddonRuleDetailHandler) ConfirmDrafts(c echo.Context) error {
	h.draftMu.Lock()
	defer h.draftMu.Unlock()

	if len(h.draftData) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No draft data to confirm",
		})
	}

	// Begin transaction
	tx, err := h.db.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to begin transaction: " + err.Error(),
		})
	}

	// Insert each draft into property_service_development.addon_rules
	insertQuery := `
		INSERT INTO property_service_development.addon_rules 
		(addon_code, product_code, construction_class_id, occupation_code, rules, created_by, created_at)
		VALUES (?, ?, ?, ?, ?, 1, CURRENT_TIMESTAMP)
	`

	// Collect distinct product_code + addon_code combinations for product_addon_mappings
	distinctMappings := make(map[string]struct{})

	for _, draft := range h.draftData {
		_, err := tx.Exec(insertQuery,
			draft.AddonCode,
			draft.ProductCode,
			draft.ConstructionClassID,
			draft.OccupationCode,
			draft.Rules,
		)
		if err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to insert addon rule: " + err.Error(),
			})
		}

		// Add to distinct mappings
		key := draft.ProductCode + "|" + draft.AddonCode
		distinctMappings[key] = struct{}{}
	}

	// Insert distinct mappings into property_service_development.product_addon_mappings
	// Use INSERT IGNORE to skip duplicates that already exist in the table
	mappingInsertQuery := `
		INSERT IGNORE INTO property_service_development.product_addon_mappings 
		(product_code, addon_code, created_by, created_at)
		VALUES (?, ?, 1, CURRENT_TIMESTAMP)
	`

	for key := range distinctMappings {
		parts := strings.SplitN(key, "|", 2)
		if len(parts) == 2 {
			productCode := parts[0]
			addonCode := parts[1]
			_, err := tx.Exec(mappingInsertQuery, productCode, addonCode)
			if err != nil {
				tx.Rollback()
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to insert product addon mapping: " + err.Error(),
				})
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to commit transaction: " + err.Error(),
		})
	}

	// Clear drafts after successful insert
	insertedCount := len(h.draftData)
	mappingCount := len(distinctMappings)
	h.draftData = nil

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "Addon rules confirmed successfully",
		"count":         insertedCount,
		"mapping_count": mappingCount,
	})
}

// GetAll handles GET /api/property/addon-rule-details
func (h *PropertyAddonRuleDetailHandler) GetAll(c echo.Context) error {
	query := `
		SELECT id, addon_code, product_code, construction_class_id, occupation_code, rules, created_by, created_at
		FROM property_service_development.addon_rules
		ORDER BY id DESC
	`

	rows, err := h.db.Query(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch addon rules: " + err.Error(),
		})
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int64
		var addonCode, productCode, rules string
		var constructionClassID, occupationCode sql.NullString
		var createdBy int64
		var createdAt string

		err := rows.Scan(&id, &addonCode, &productCode, &constructionClassID, &occupationCode, &rules, &createdBy, &createdAt)
		if err != nil {
			continue
		}

		result := map[string]interface{}{
			"id":                    id,
			"addon_code":            addonCode,
			"product_code":          productCode,
			"construction_class_id": nil,
			"occupation_code":       nil,
			"rules":                 rules,
			"created_by":            createdBy,
			"created_at":            createdAt,
		}

		if constructionClassID.Valid {
			result["construction_class_id"] = constructionClassID.String
		}
		if occupationCode.Valid {
			result["occupation_code"] = occupationCode.String
		}

		results = append(results, result)
	}

	return c.JSON(http.StatusOK, results)
}

// Update handles PUT /api/property/addon-rule-details/:id
func (h *PropertyAddonRuleDetailHandler) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "ID is required",
		})
	}

	// Parse request body
	var req struct {
		AddonCode           string `json:"addon_code"`
		ProductCode         string `json:"product_code"`
		ConstructionClassID string `json:"construction_class_id"`
		OccupationCode      string `json:"occupation_code"`
		Rules               string `json:"rules"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Update in database
	query := `
		UPDATE property_service_development.addon_rules 
		SET addon_code = ?, product_code = ?, construction_class_id = ?, occupation_code = ?, rules = ?
		WHERE id = ?
	`

	result, err := h.db.Exec(query, req.AddonCode, req.ProductCode, req.ConstructionClassID, req.OccupationCode, req.Rules, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update addon rule: " + err.Error(),
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Addon rule not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Addon rule updated successfully",
	})
}

// Delete handles DELETE /api/property/addon-rule-details/:id
func (h *PropertyAddonRuleDetailHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "ID is required",
		})
	}

	// Delete from database
	query := `DELETE FROM property_service_development.addon_rules WHERE id = ?`

	result, err := h.db.Exec(query, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete addon rule: " + err.Error(),
		})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Addon rule not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Addon rule deleted successfully",
	})
}


