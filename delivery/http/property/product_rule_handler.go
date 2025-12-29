package http

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/labstack/echo"
)

// PropertyProductRuleDraft represents a draft product rule for property
type PropertyProductRuleDraft struct {
	ProductCode         string  `json:"product_code"`
	OccupationCode      string  `json:"occupation_code"`
	BasePremiumValue    float64 `json:"base_premium_value"`
	BasePremiumType     string  `json:"base_premium_type"`
	CoverageType        *string `json:"coverage_type,omitempty"`
	ConstructionType    *string `json:"construction_type,omitempty"`
	BuildingType        *string `json:"building_type,omitempty"`
	ConstructionClassID *string `json:"construction_class_id,omitempty"`
	ConstructionClass   *string `json:"construction_class,omitempty"`
	Rules               string  `json:"rules"`
}

// PropertyProductRuleHandler handles HTTP requests for property product rules
type PropertyProductRuleHandler struct {
	db        *sql.DB
	draftData []PropertyProductRuleDraft
	draftMu   sync.Mutex
}

// NewPropertyProductRuleHandler creates a new property product rule handler
func NewPropertyProductRuleHandler(db *sql.DB) *PropertyProductRuleHandler {
	return &PropertyProductRuleHandler{
		db:        db,
		draftData: make([]PropertyProductRuleDraft, 0),
	}
}

// DownloadTemplate handles GET /api/property/product-rules/template
func (h *PropertyProductRuleHandler) DownloadTemplate(c echo.Context) error {
	// Fetch all products from property_service_development.products
	query := `SELECT code, name FROM property_service_development.products ORDER BY code ASC`
	rows, err := h.db.Query(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch products: " + err.Error(),
		})
	}
	defer rows.Close()

	type ProductInfo struct {
		Code string
		Name string
	}
	var products []ProductInfo
	for rows.Next() {
		var code, name string
		if err := rows.Scan(&code, &name); err != nil {
			continue
		}
		products = append(products, ProductInfo{Code: code, Name: name})
	}

	// List of occupation codes
	occupationCodes := []string{
		"2491", "2922", "2928", "29291", "2930", "29314", "2932", "2933", "2934", "29395",
		"29411", "29412", "29432", "2945", "2946", "2947", "2951", "2952", "2953", "2955",
		"2961", "2962", "2963", "2969", "2971", "2972", "2973", "2975", "2976", "29761", "2978",
	}

	// Default rules JSON for property
	rulesJSON := `{"start_tsi":0, "end_tsi":0}`

	// Build CSV template with columns matching property_service_development.product_rules table
	// Columns: product_code, occupation_code, base_premium_value, base_premium_type, coverage_type, construction_type, building_type, construction_class_id, construction_class, rules
	var template strings.Builder
	template.WriteString("product_code,occupation_code,base_premium_value,base_premium_type,coverage_type,construction_type,building_type,construction_class_id,construction_class,rules\n")

	// Generate template for each combination of product_code x occupation_code
	if len(products) == 0 {
		// Default example row if no products exist
		template.WriteString(fmt.Sprintf("PR-DAMAI-PAR-SME-01,%s,0,PERCENTAGE,,,,\"\",\"\",\"%s\"\n", occupationCodes[0], rulesJSON))
	} else {
		// Generate one row for each product_code x occupation_code combination
		for _, product := range products {
			productCode := product.Code
			for _, occCode := range occupationCodes {
				// Generate row with default values
				// User will fill in: base_premium_value, coverage_type, construction_type, building_type, construction_class_id, construction_class
				// Note: Empty fields for nullable columns (coverage_type, construction_type, building_type, construction_class_id, construction_class)
				// rules is pre-filled with {"start_tsi":0, "end_tsi":0}
				template.WriteString(fmt.Sprintf("%s,%s,0,PERCENTAGE,,,,\"\",\"\",\"%s\"\n", productCode, occCode, rulesJSON))
			}
		}
	}

	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=property_product_rules_template.csv")
	return c.String(http.StatusOK, template.String())
}

// Upload handles POST /api/property/product-rules/upload
func (h *PropertyProductRuleHandler) Upload(c echo.Context) error {
	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to parse form: " + err.Error(),
		})
	}
	defer form.RemoveAll()

	fileHeader := form.File["file"]
	if len(fileHeader) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "File is required",
		})
	}

	file, err := fileHeader[0].Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to open file: " + err.Error(),
		})
	}
	defer file.Close()

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to read file: " + err.Error(),
		})
	}

	// Parse CSV content
	drafts, err := h.parseCSVContent(fileBytes)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to parse file: " + err.Error(),
		})
	}

	// Clear existing drafts
	h.draftMu.Lock()
	h.draftData = make([]PropertyProductRuleDraft, 0)
	h.draftMu.Unlock()

	// Add all drafts
	h.draftMu.Lock()
	h.draftData = append(h.draftData, drafts...)
	h.draftMu.Unlock()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "Success",
		"data":   drafts,
		"count":  len(drafts),
	})
}

// GetAllDrafts handles GET /api/property/product-rules/draft
func (h *PropertyProductRuleHandler) GetAllDrafts(c echo.Context) error {
	h.draftMu.Lock()
	drafts := make([]PropertyProductRuleDraft, len(h.draftData))
	copy(drafts, h.draftData)
	h.draftMu.Unlock()

	return c.JSON(http.StatusOK, drafts)
}

// ClearDrafts handles POST /api/property/product-rules/draft/clear
func (h *PropertyProductRuleHandler) ClearDrafts(c echo.Context) error {
	h.draftMu.Lock()
	h.draftData = make([]PropertyProductRuleDraft, 0)
	h.draftMu.Unlock()

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Drafts cleared successfully",
	})
}

// ConfirmDrafts handles POST /api/property/product-rules/draft/confirm
func (h *PropertyProductRuleHandler) ConfirmDrafts(c echo.Context) error {
	// Get created_by from request body or default to 1
	var requestData struct {
		CreatedBy int64 `json:"created_by"`
	}

	if err := c.Bind(&requestData); err != nil {
		requestData.CreatedBy = 1 // Default value
	}
	if requestData.CreatedBy == 0 {
		requestData.CreatedBy = 1 // Default value
	}

	// Get all drafts
	h.draftMu.Lock()
	drafts := make([]PropertyProductRuleDraft, len(h.draftData))
	copy(drafts, h.draftData)
	h.draftMu.Unlock()

	if len(drafts) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No draft rules to confirm",
		})
	}

	// Insert each draft into database
	insertedCount := 0
	for _, draft := range drafts {
		// Prepare values for nullable fields
		var coverageType, constructionType, buildingType, constructionClassID, constructionClass sql.NullString

		if draft.CoverageType != nil && *draft.CoverageType != "" {
			coverageType = sql.NullString{String: *draft.CoverageType, Valid: true}
		}
		if draft.ConstructionType != nil && *draft.ConstructionType != "" {
			constructionType = sql.NullString{String: *draft.ConstructionType, Valid: true}
		}
		if draft.BuildingType != nil && *draft.BuildingType != "" {
			buildingType = sql.NullString{String: *draft.BuildingType, Valid: true}
		}
		if draft.ConstructionClassID != nil && *draft.ConstructionClassID != "" {
			constructionClassID = sql.NullString{String: *draft.ConstructionClassID, Valid: true}
		}
		if draft.ConstructionClass != nil && *draft.ConstructionClass != "" {
			constructionClass = sql.NullString{String: *draft.ConstructionClass, Valid: true}
		}

		query := `INSERT INTO property_service_development.product_rules 
			(product_code, occupation_code, base_premium_value, base_premium_type, 
			 coverage_type, construction_type, building_type, construction_class_id, 
			 construction_class, rules, created_by, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

		_, err := h.db.Exec(query,
			draft.ProductCode,
			draft.OccupationCode,
			draft.BasePremiumValue,
			draft.BasePremiumType,
			coverageType,
			constructionType,
			buildingType,
			constructionClassID,
			constructionClass,
			draft.Rules,
			requestData.CreatedBy,
		)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": fmt.Sprintf("Failed to insert rule: %v", err),
			})
		}
		insertedCount++
	}

	// Clear drafts after successful confirmation
	h.draftMu.Lock()
	h.draftData = make([]PropertyProductRuleDraft, 0)
	h.draftMu.Unlock()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":        fmt.Sprintf("Successfully confirmed %d rules", insertedCount),
		"inserted_count": insertedCount,
	})
}

// UpdateDraft handles PUT /api/property/product-rules/draft/:index
func (h *PropertyProductRuleHandler) UpdateDraft(c echo.Context) error {
	indexStr := c.Param("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid draft index: " + err.Error(),
		})
	}

	var draft PropertyProductRuleDraft
	if err := c.Bind(&draft); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	h.draftMu.Lock()
	defer h.draftMu.Unlock()

	if index < 0 || index >= len(h.draftData) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Draft index out of range",
		})
	}

	h.draftData[index] = draft

	return c.JSON(http.StatusOK, draft)
}

// DeleteDraft handles DELETE /api/property/product-rules/draft/:index
func (h *PropertyProductRuleHandler) DeleteDraft(c echo.Context) error {
	indexStr := c.Param("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid draft index: " + err.Error(),
		})
	}

	h.draftMu.Lock()
	defer h.draftMu.Unlock()

	if index < 0 || index >= len(h.draftData) {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Draft index out of range",
		})
	}

	h.draftData = append(h.draftData[:index], h.draftData[index+1:]...)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Draft deleted successfully",
	})
}

// GetAll handles GET /api/property/product-rules
func (h *PropertyProductRuleHandler) GetAll(c echo.Context) error {
	productCode := c.QueryParam("product_code")
	insuranceCode := c.QueryParam("insurance_code")

	// Build query
	query := `SELECT id, product_code, occupation_code, base_premium_value, base_premium_type,
		coverage_type, construction_type, building_type, construction_class_id, construction_class,
		rules, created_by, created_at
		FROM property_service_development.product_rules WHERE 1=1`

	var args []interface{}

	if productCode != "" {
		query += " AND product_code = ?"
		args = append(args, productCode)
	}

	if insuranceCode != "" {
		// Filter by insurance code - support multiple formats:
		// - PR-{INSURANCE_CODE}-... (e.g., PR-DAMAI-PAR-SME-01)
		// - PR-HOME-{INSURANCE_CODE}-... (e.g., PR-HOME-DAMAI-SR-01)
		// - PR-*-{INSURANCE_CODE}-... (any other format)
		query += " AND (product_code LIKE ? OR product_code LIKE ?)"
		args = append(args, "PR-"+insuranceCode+"-%", "PR-HOME-"+insuranceCode+"-%")
	}

	query += " ORDER BY id DESC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch product rules: " + err.Error(),
		})
	}
	defer rows.Close()

	type ProductRule struct {
		ID                  int64   `json:"id"`
		ProductCode         string  `json:"product_code"`
		OccupationCode      string  `json:"occupation_code"`
		BasePremiumValue    float64 `json:"base_premium_value"`
		BasePremiumType     string  `json:"base_premium_type"`
		CoverageType        *string `json:"coverage_type,omitempty"`
		ConstructionType    *string `json:"construction_type,omitempty"`
		BuildingType        *string `json:"building_type,omitempty"`
		ConstructionClassID *string `json:"construction_class_id,omitempty"`
		ConstructionClass   *string `json:"construction_class,omitempty"`
		Rules               string  `json:"rules"`
		CreatedBy           int64   `json:"created_by"`
		CreatedAt           string  `json:"created_at"`
	}

	var rules []ProductRule
	for rows.Next() {
		var rule ProductRule
		var coverageType, constructionType, buildingType, constructionClassID, constructionClass sql.NullString
		var createdAt sql.NullTime

		err := rows.Scan(
			&rule.ID,
			&rule.ProductCode,
			&rule.OccupationCode,
			&rule.BasePremiumValue,
			&rule.BasePremiumType,
			&coverageType,
			&constructionType,
			&buildingType,
			&constructionClassID,
			&constructionClass,
			&rule.Rules,
			&rule.CreatedBy,
			&createdAt,
		)
		if err != nil {
			continue
		}

		if coverageType.Valid {
			rule.CoverageType = &coverageType.String
		}
		if constructionType.Valid {
			rule.ConstructionType = &constructionType.String
		}
		if buildingType.Valid {
			rule.BuildingType = &buildingType.String
		}
		if constructionClassID.Valid {
			rule.ConstructionClassID = &constructionClassID.String
		}
		if constructionClass.Valid {
			rule.ConstructionClass = &constructionClass.String
		}
		if createdAt.Valid {
			rule.CreatedAt = createdAt.Time.Format("2006-01-02 15:04:05")
		}

		rules = append(rules, rule)
	}

	return c.JSON(http.StatusOK, rules)
}

// parseCSVContent parses CSV content and returns draft product rules for property
func (h *PropertyProductRuleHandler) parseCSVContent(content []byte) ([]PropertyProductRuleDraft, error) {
	reader := csv.NewReader(strings.NewReader(string(content)))
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	headerMap := make(map[string]int)
	for i, h := range header {
		headerMap[strings.TrimSpace(h)] = i
	}

	// Validate required headers
	requiredHeaders := []string{"product_code", "occupation_code", "base_premium_value", "base_premium_type", "rules"}
	for _, expected := range requiredHeaders {
		if _, ok := headerMap[expected]; !ok {
			return nil, fmt.Errorf("missing required header: %s", expected)
		}
	}

	// Helper function to get value from row
	getValue := func(row []string, key string) string {
		if idx, ok := headerMap[key]; ok && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	// Parse data rows
	var drafts []PropertyProductRuleDraft
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // Skip invalid rows
		}

		// Parse base_premium_value
		basePremiumValueStr := getValue(row, "base_premium_value")
		basePremiumValue, err := strconv.ParseFloat(basePremiumValueStr, 64)
		if err != nil {
			basePremiumValue = 0
		}

		// Parse rules (remove surrounding quotes if present)
		rulesValue := getValue(row, "rules")
		if len(rulesValue) >= 2 && rulesValue[0] == '"' && rulesValue[len(rulesValue)-1] == '"' {
			rulesValue = rulesValue[1 : len(rulesValue)-1]
		}

		// Handle nullable fields
		var coverageType, constructionType, buildingType, constructionClassID, constructionClass *string

		coverageTypeVal := getValue(row, "coverage_type")
		if coverageTypeVal != "" {
			coverageType = &coverageTypeVal
		}

		constructionTypeVal := getValue(row, "construction_type")
		if constructionTypeVal != "" {
			constructionType = &constructionTypeVal
		}

		buildingTypeVal := getValue(row, "building_type")
		if buildingTypeVal != "" {
			buildingType = &buildingTypeVal
		}

		constructionClassIDVal := getValue(row, "construction_class_id")
		if constructionClassIDVal != "" {
			constructionClassID = &constructionClassIDVal
		}

		constructionClassVal := getValue(row, "construction_class")
		if constructionClassVal != "" {
			constructionClass = &constructionClassVal
		}

		draft := PropertyProductRuleDraft{
			ProductCode:         getValue(row, "product_code"),
			OccupationCode:      getValue(row, "occupation_code"),
			BasePremiumValue:    basePremiumValue,
			BasePremiumType:     getValue(row, "base_premium_type"),
			CoverageType:        coverageType,
			ConstructionType:    constructionType,
			BuildingType:        buildingType,
			ConstructionClassID: constructionClassID,
			ConstructionClass:   constructionClass,
			Rules:               rulesValue,
		}

		drafts = append(drafts, draft)
	}

	return drafts, nil
}


