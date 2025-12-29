package http

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	model "multi-onboarding/model/Travel"
	repository "multi-onboarding/repository/Travel"
	usecase "multi-onboarding/usecase/Travel"

	"github.com/gorilla/mux"
	"github.com/labstack/echo"
)

// AddonRuleDetailHandler handles HTTP requests for addon rule details
type AddonRuleDetailHandler struct {
	usecase     *usecase.AddonRuleDetailUsecase
	productRepo *sql.DB
}

// NewAddonRuleDetailHandler creates a new addon rule detail handler
func NewAddonRuleDetailHandler(usecase *usecase.AddonRuleDetailUsecase, db *sql.DB) *AddonRuleDetailHandler {
	return &AddonRuleDetailHandler{
		usecase:     usecase,
		productRepo: db,
	}
}

// DownloadTemplate handles GET /api/addon-rule-details/template
func (h *AddonRuleDetailHandler) DownloadTemplate(w http.ResponseWriter, r *http.Request) {
	// Fetch only products that exist in addon_rules (products that are already confirmed in section products)
	productQuery := `
		SELECT DISTINCT p.code, p.name 
		FROM travel_service_development.products p
		INNER JOIN travel_service_development.addon_rules ar ON p.code = ar.product_code
		ORDER BY p.code ASC
	`
	productRows, err := h.productRepo.Query(productQuery)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch products: "+err.Error())
		return
	}
	defer productRows.Close()

	type ProductInfo struct {
		Code string
		Name string
	}
	var products []ProductInfo
	productMap := make(map[string]ProductInfo) // Map product_code -> ProductInfo
	for productRows.Next() {
		var p ProductInfo
		if err := productRows.Scan(&p.Code, &p.Name); err != nil {
			continue
		}
		products = append(products, p)
		productMap[p.Code] = p
	}

	// Fetch premium_type from product_rules table (join with products to get product_code)
	premiumTypeQuery := `
		SELECT DISTINCT pr.product_code, pr.premium_type 
		FROM travel_service_development.product_rules pr
		WHERE pr.product_code IN (
			SELECT DISTINCT product_code FROM travel_service_development.addon_rules
		)
		ORDER BY pr.product_code ASC
	`
	premiumTypeRows, err := h.productRepo.Query(premiumTypeQuery)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch premium types: "+err.Error())
		return
	}
	defer premiumTypeRows.Close()

	// Map product_code -> premium_type
	premiumTypeMap := make(map[string]string)
	for premiumTypeRows.Next() {
		var productCode, premiumType string
		if err := premiumTypeRows.Scan(&productCode, &premiumType); err != nil {
			continue
		}
		// Use the first premium_type found for each product_code
		// If multiple premium_types exist, we'll use the first one
		if _, exists := premiumTypeMap[productCode]; !exists {
			premiumTypeMap[productCode] = premiumType
		}
	}

	// Fetch all addons with names from database
	addonQuery := `SELECT code, name FROM travel_service_development.addons ORDER BY code ASC`
	addonRows, err := h.productRepo.Query(addonQuery)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch addons: "+err.Error())
		return
	}
	defer addonRows.Close()

	type AddonInfo struct {
		Code string
		Name string
	}
	var addons []AddonInfo
	for addonRows.Next() {
		var a AddonInfo
		if err := addonRows.Scan(&a.Code, &a.Name); err != nil {
			continue
		}
		addons = append(addons, a)
	}

	// Fetch addon_rules to get addon_rule_id mapping
	addonRuleQuery := `SELECT id, product_code, addon_code FROM travel_service_development.addon_rules`
	addonRuleRows, err := h.productRepo.Query(addonRuleQuery)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch addon rules: "+err.Error())
		return
	}
	defer addonRuleRows.Close()

	// Create map: product_code+addon_code -> addon_rule_id
	type AddonRuleKey struct {
		ProductCode string
		AddonCode   string
	}
	addonRuleMap := make(map[AddonRuleKey]int64)
	type AddonRuleInfo struct {
		ID          int64
		ProductCode string
		AddonCode   string
	}
	var addonRules []AddonRuleInfo
	for addonRuleRows.Next() {
		var ar AddonRuleInfo
		if err := addonRuleRows.Scan(&ar.ID, &ar.ProductCode, &ar.AddonCode); err != nil {
			continue
		}
		key := AddonRuleKey{ProductCode: ar.ProductCode, AddonCode: ar.AddonCode}
		addonRuleMap[key] = ar.ID
		addonRules = append(addonRules, ar)
	}

	// Build CSV template with all required columns
	var template strings.Builder
	template.WriteString("product_name,addon_name,product_code,addon_code,addon_rule_id,start_condition,end_condition,value_type,value,duration_rule_type,min_adult,max_adult,max_age\n")

	// Only include combinations that exist in addon_rules table
	if len(addonRules) == 0 {
		// If no addon_rules, return empty template with just header
	} else {
		// Generate template only for product+addon combinations that exist in addon_rules
		for _, ar := range addonRules {
			// Find product info
			product, productExists := productMap[ar.ProductCode]
			var productName string
			if productExists {
				productName = product.Name
			} else {
				productName = ar.ProductCode // Fallback to code if name not found
			}

			// Get premium_type from premiumTypeMap
			premiumType, premiumTypeExists := premiumTypeMap[ar.ProductCode]
			if !premiumTypeExists {
				premiumType = "DAILY" // Default to DAILY if premium_type not found
			}

			// Find addon name
			var addonName string
			for _, addon := range addons {
				if addon.Code == ar.AddonCode {
					addonName = addon.Name
					break
				}
			}
			if addonName == "" {
				addonName = ar.AddonCode // Fallback to code if name not found
			}

			// Generate rows based on product premium_type
			// If premium_type is ANNUAL, only generate ANNUAL duration_rule_type (1 row per addon code)
			// If premium_type is DAILY, only generate DAILY and PER_EXTRA_SEVEN_DAYS duration_rule_type (1 row per addon code for PER_EXTRA_SEVEN_DAYS)
			if premiumType == "ANNUAL" {
				// ANNUAL: 1 row only
				template.WriteString(fmt.Sprintf("%s,%s,%s,%s,%d,%d,%d,FIXED,0,ANNUAL,0,0,0\n",
					productName,
					addonName,
					ar.ProductCode,
					ar.AddonCode,
					ar.ID,
					1,   // start_condition
					365, // end_condition
				))
			} else {
				// DAILY: 8 rows
				dailyStartDays := []int{1, 5, 7, 9, 11, 16, 21, 26}
				dailyEndDays := []int{4, 6, 8, 10, 15, 20, 25, 31}
				for i := 0; i < len(dailyStartDays); i++ {
					template.WriteString(fmt.Sprintf("%s,%s,%s,%s,%d,%d,%d,FIXED,0,DAILY,0,0,0\n",
						productName,
						addonName,
						ar.ProductCode,
						ar.AddonCode,
						ar.ID,
						dailyStartDays[i],
						dailyEndDays[i],
					))
				}

				// PER_EXTRA_SEVEN_DAYS: 1 row only
				template.WriteString(fmt.Sprintf("%s,%s,%s,%s,%d,%d,%d,FIXED,0,PER_EXTRA_SEVEN_DAYS,0,0,0\n",
					productName,
					addonName,
					ar.ProductCode,
					ar.AddonCode,
					ar.ID,
					32,  // start_condition
					180, // end_condition
				))
			}
		}
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=addon_rule_details_template.csv")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(template.String()))
}

// UploadExcel handles POST /api/addon-rule-details/upload
func (h *AddonRuleDetailHandler) UploadExcel(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse form: "+err.Error())
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "File is required")
		return
	}
	defer file.Close()

	// Read file content
	fileBytes := make([]byte, header.Size)
	_, err = file.Read(fileBytes)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to read file: "+err.Error())
		return
	}

	// Parse CSV content
	drafts, err := parseAddonRuleDetailCSV(fileBytes)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse file: "+err.Error())
		return
	}

	// Fetch addon_rules to get addon_rule_id mapping
	addonRuleQuery := `SELECT id, product_code, addon_code FROM travel_service_development.addon_rules`
	addonRuleRows, err := h.productRepo.Query(addonRuleQuery)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch addon rules: "+err.Error())
		return
	}
	defer addonRuleRows.Close()

	// Create map: product_code+addon_code -> addon_rule_id
	type AddonRuleKey struct {
		ProductCode string
		AddonCode   string
	}
	addonRuleMap := make(map[AddonRuleKey]int64)
	for addonRuleRows.Next() {
		var id int64
		var productCode, addonCode string
		if err := addonRuleRows.Scan(&id, &productCode, &addonCode); err != nil {
			continue
		}
		key := AddonRuleKey{ProductCode: productCode, AddonCode: addonCode}
		addonRuleMap[key] = id
	}

	// Fill addon_rule_id for each draft (if not already set from CSV)
	for i := range drafts {
		if drafts[i].AddonRuleID == "" {
			key := AddonRuleKey{ProductCode: drafts[i].ProductCode, AddonCode: drafts[i].AddonCode}
			if id, exists := addonRuleMap[key]; exists {
				drafts[i].AddonRuleID = fmt.Sprintf("%d", id)
			}
		}
	}

	// Clear existing drafts
	if err := h.usecase.ClearDrafts(); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to clear drafts: "+err.Error())
		return
	}

	// Add all drafts
	for _, draft := range drafts {
		if err := h.usecase.AddDraft(&draft); err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to add draft: %v", err))
			return
		}
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status": "Success",
		"data":   drafts,
		"count":  len(drafts),
	})
}

// GetAllDrafts handles GET /api/addon-rule-details/draft
func (h *AddonRuleDetailHandler) GetAllDrafts(w http.ResponseWriter, r *http.Request) {
	drafts, err := h.usecase.GetAllDrafts()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, drafts)
}

// ClearDrafts handles POST /api/addon-rule-details/draft/clear
func (h *AddonRuleDetailHandler) ClearDrafts(w http.ResponseWriter, r *http.Request) {
	if err := h.usecase.ClearDrafts(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Drafts cleared successfully"})
}

// ConfirmDrafts handles POST /api/addon-rule-details/draft/confirm
func (h *AddonRuleDetailHandler) ConfirmDrafts(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CreatedBy int64 `json:"created_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.usecase.ConfirmDrafts(req.CreatedBy); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Addon rule details confirmed successfully"})
}

// GetAll handles GET /api/addon-rule-details
func (h *AddonRuleDetailHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	details, err := h.usecase.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": details,
	})
}

// GetByID handles GET /api/addon-rule-details/{id}
func (h *AddonRuleDetailHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Try to get id from Echo context first (if available)
	var idStr string
	if echoCtx := r.Context().Value("echo_context"); echoCtx != nil {
		if c, ok := echoCtx.(echo.Context); ok {
			idStr = c.Param("id")
		}
	}

	// Fallback: Extract id from URL path
	if idStr == "" {
		// Path format: /api/addon-rule-details/{id}
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		for i, part := range pathParts {
			if part == "addon-rule-details" && i+1 < len(pathParts) {
				idStr = pathParts[i+1]
				break
			}
		}
	}

	// Try mux.Vars as last fallback
	if idStr == "" {
		vars := mux.Vars(r)
		if val, ok := vars["id"]; ok {
			idStr = val
		}
	}

	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid ID: id parameter not found")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID: "+err.Error())
		return
	}

	detail, err := h.usecase.GetByID(id)
	if err != nil {
		if err == repository.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "Addon rule detail not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, detail)
}

// Update handles PUT /api/addon-rule-details/{id}
func (h *AddonRuleDetailHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Try to get id from Echo context first (if available)
	var idStr string
	if echoCtx := r.Context().Value("echo_context"); echoCtx != nil {
		if c, ok := echoCtx.(echo.Context); ok {
			idStr = c.Param("id")
		}
	}

	// Fallback: Extract id from URL path
	if idStr == "" {
		// Path format: /api/addon-rule-details/{id}
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		for i, part := range pathParts {
			if part == "addon-rule-details" && i+1 < len(pathParts) {
				idStr = pathParts[i+1]
				break
			}
		}
	}

	// Try mux.Vars as last fallback
	if idStr == "" {
		vars := mux.Vars(r)
		if val, ok := vars["id"]; ok {
			idStr = val
		}
	}

	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid ID: id parameter not found")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID: "+err.Error())
		return
	}

	var detail model.AddonRuleDetail
	if err := json.NewDecoder(r.Body).Decode(&detail); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.usecase.Update(id, &detail); err != nil {
		if err == repository.ErrNotFound {
			respondWithError(w, http.StatusNotFound, "Addon rule detail not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Addon rule detail updated successfully"})
}

// parseAddonRuleDetailCSV parses CSV content and returns draft addon rule details
func parseAddonRuleDetailCSV(content []byte) ([]model.AddonRuleDetailDraft, error) {
	reader := csv.NewReader(strings.NewReader(string(content)))
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	expectedHeaders := []string{"product_code", "addon_code", "start_condition", "end_condition", "value_type", "value", "duration_rule_type", "min_adult", "max_adult", "max_age"}
	headerMap := make(map[string]int)
	for i, h := range header {
		headerMap[strings.TrimSpace(h)] = i
	}

	// Validate required headers (product_name, addon_name, and addon_rule_id are optional, for display/reference only)
	for _, expected := range expectedHeaders {
		if _, ok := headerMap[expected]; !ok {
			return nil, fmt.Errorf("missing required header: %s", expected)
		}
	}

	// Helper function to get value from CSV row
	getValue := func(row []string, headerMap map[string]int, key string) string {
		if idx, ok := headerMap[key]; ok && idx < len(row) {
			return row[idx]
		}
		return ""
	}

	// Helper function to parse float64
	parseFloat := func(s string) float64 {
		val, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
		if err != nil {
			return 0
		}
		return val
	}

	// Helper function to parse int
	parseInt := func(s string) int {
		val, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			return 0
		}
		return val
	}

	// Parse data rows
	var drafts []model.AddonRuleDetailDraft
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // Skip invalid rows
		}

		if len(row) < len(expectedHeaders) {
			continue // Skip incomplete rows
		}

		// Get product_name, addon_name, and addon_rule_id (for display, but not saved to DB)
		productName := strings.TrimSpace(getValue(row, headerMap, "product_name"))
		addonName := strings.TrimSpace(getValue(row, headerMap, "addon_name"))
		addonRuleID := strings.TrimSpace(getValue(row, headerMap, "addon_rule_id"))

		draft := model.AddonRuleDetailDraft{
			ProductName:      productName,
			AddonName:        addonName,
			AddonRuleID:      addonRuleID,
			ProductCode:      strings.TrimSpace(getValue(row, headerMap, "product_code")),
			AddonCode:        strings.TrimSpace(getValue(row, headerMap, "addon_code")),
			StartCondition:   strings.TrimSpace(getValue(row, headerMap, "start_condition")),
			EndCondition:     strings.TrimSpace(getValue(row, headerMap, "end_condition")),
			ValueType:        strings.TrimSpace(getValue(row, headerMap, "value_type")),
			Value:            parseFloat(getValue(row, headerMap, "value")),
			DurationRuleType: strings.TrimSpace(getValue(row, headerMap, "duration_rule_type")),
			MinAdult:         parseInt(getValue(row, headerMap, "min_adult")),
			MaxAdult:         parseInt(getValue(row, headerMap, "max_adult")),
			MaxAge:           parseInt(getValue(row, headerMap, "max_age")),
		}

		// Set defaults if empty
		if draft.StartCondition == "" {
			draft.StartCondition = "-1"
		}
		if draft.EndCondition == "" {
			draft.EndCondition = "-1"
		}
		if draft.ValueType == "" {
			draft.ValueType = "FIXED"
		}
		if draft.DurationRuleType == "" {
			draft.DurationRuleType = "DAILY"
		}

		drafts = append(drafts, draft)
	}

	return drafts, nil
}
