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

	"github.com/gorilla/mux"
	"github.com/labstack/echo"
	"multi-onboarding/model/Travel"
	"multi-onboarding/usecase/Travel"
)

// ProductRuleHandler handles HTTP requests for product rules
type ProductRuleHandler struct {
	usecase      *usecase.ProductRuleUsecase
	productRepo  *sql.DB // We'll use direct DB access to get products
}

// NewProductRuleHandler creates a new product rule handler
func NewProductRuleHandler(usecase *usecase.ProductRuleUsecase, db *sql.DB) *ProductRuleHandler {
	return &ProductRuleHandler{
		usecase:     usecase,
		productRepo: db,
	}
}

// DownloadTemplate handles GET /api/product-rules/template
func (h *ProductRuleHandler) DownloadTemplate(w http.ResponseWriter, r *http.Request) {
	// Fetch all products from database with name to identify 90 Days vs 180 Days
	query := `SELECT code, name FROM travel_service_development.products ORDER BY code ASC`
	rows, err := h.productRepo.Query(query)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch products: "+err.Error())
		return
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

	rulesJSON := `{"min_child_age": 0 ,"max_child_age": 18,"min_adult_age": 17,"max_adult_age": 69}`
	
	// Build CSV template
	var template strings.Builder
	template.WriteString("product_code,premium_type,start_days,end_days,base_premium_value,rules\n")

	// If no products, use default template
	if len(products) == 0 {
		template.WriteString("TV-DAMAI-DOM-IND-01,DAILY,1,3,0,\"" + rulesJSON + "\"\n")
		template.WriteString("TV-DAMAI-DOM-IND-01,DAILY,4,7,0,\"" + rulesJSON + "\"\n")
		template.WriteString("TV-DAMAI-DOM-IND-01,DAILY,8,10,0,\"" + rulesJSON + "\"\n")
		template.WriteString("TV-DAMAI-DOM-IND-01,DAILY,11,15,0,\"" + rulesJSON + "\"\n")
		template.WriteString("TV-DAMAI-DOM-IND-01,DAILY,16,20,0,\"" + rulesJSON + "\"\n")
		template.WriteString("TV-DAMAI-DOM-IND-01,DAILY,21,25,0,\"" + rulesJSON + "\"\n")
		template.WriteString("TV-DAMAI-DOM-IND-01,DAILY,26,31,0,\"" + rulesJSON + "\"\n")
		template.WriteString("TV-DAMAI-DOM-IND-01,PER_EXTRA_SEVEN_DAYS,32,180,0,\"" + rulesJSON + "\"\n")
	} else {
		// Generate template for each product code
		for _, product := range products {
			productCode := product.Code
			productName := product.Name
			
			// Check if product code contains "-ANN-" (ANNUAL)
			if strings.Contains(productCode, "-ANN-") {
				// Check if product name contains "(90 Days)" or "(180 Days)"
				if strings.Contains(productName, "(90 Days)") {
					// For ANNUAL 90 Days: start_days = 1, end_days = 90
					template.WriteString(productCode + ",ANNUAL,1,90,0,\"" + rulesJSON + "\"\n")
				} else if strings.Contains(productName, "(180 Days)") {
					// For ANNUAL 180 Days: start_days = 1, end_days = 180
					template.WriteString(productCode + ",ANNUAL,1,180,0,\"" + rulesJSON + "\"\n")
				} else {
					// Fallback: if name doesn't contain days info, use default 365
					template.WriteString(productCode + ",ANNUAL,1,365,0,\"" + rulesJSON + "\"\n")
				}
			} else {
				// For DAILY: 7 rows + 1 PER_EXTRA_SEVEN_DAYS row
				template.WriteString(productCode + ",DAILY,1,3,0,\"" + rulesJSON + "\"\n")
				template.WriteString(productCode + ",DAILY,4,7,0,\"" + rulesJSON + "\"\n")
				template.WriteString(productCode + ",DAILY,8,10,0,\"" + rulesJSON + "\"\n")
				template.WriteString(productCode + ",DAILY,11,15,0,\"" + rulesJSON + "\"\n")
				template.WriteString(productCode + ",DAILY,16,20,0,\"" + rulesJSON + "\"\n")
				template.WriteString(productCode + ",DAILY,21,25,0,\"" + rulesJSON + "\"\n")
				template.WriteString(productCode + ",DAILY,26,31,0,\"" + rulesJSON + "\"\n")
				template.WriteString(productCode + ",PER_EXTRA_SEVEN_DAYS,32,180,0,\"" + rulesJSON + "\"\n")
			}
		}
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=product_rules_template.csv")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(template.String()))
}

// UploadExcel handles POST /api/product-rules/upload
func (h *ProductRuleHandler) UploadExcel(w http.ResponseWriter, r *http.Request) {
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

	// Parse CSV/Excel content
	drafts, err := parseCSVContent(fileBytes)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse file: "+err.Error())
		return
	}

	// Auto-correct end_days for ANNUAL products based on product name
	// Fetch all products to get product names
	productQuery := `SELECT code, name FROM travel_service_development.products`
	productRows, err := h.productRepo.Query(productQuery)
	if err == nil {
		defer productRows.Close()
		
		// Create map of product code to product name
		productNameMap := make(map[string]string)
		for productRows.Next() {
			var code, name string
			if err := productRows.Scan(&code, &name); err == nil {
				productNameMap[code] = name
			}
		}
		
		// Correct end_days for ANNUAL products
		for i := range drafts {
			if strings.Contains(drafts[i].ProductCode, "-ANN-") && drafts[i].PremiumType == "ANNUAL" {
				productName := productNameMap[drafts[i].ProductCode]
				if strings.Contains(productName, "(90 Days)") {
					drafts[i].EndDays = 90
				} else if strings.Contains(productName, "(180 Days)") {
					drafts[i].EndDays = 180
				}
				// If product name doesn't contain days info, keep the original end_days from CSV
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

// GetAllDrafts handles GET /api/product-rules/draft
func (h *ProductRuleHandler) GetAllDrafts(w http.ResponseWriter, r *http.Request) {
	drafts, err := h.usecase.GetAllDrafts()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, drafts)
}

// UpdateDraft handles PUT /api/product-rules/draft/{index}
func (h *ProductRuleHandler) UpdateDraft(w http.ResponseWriter, r *http.Request) {
	// Try to get index from Echo context first (if available)
	var indexStr string
	if echoCtx := r.Context().Value("echo_context"); echoCtx != nil {
		if c, ok := echoCtx.(echo.Context); ok {
			indexStr = c.Param("index")
		}
	}
	
	// Fallback: Extract index from URL path
	if indexStr == "" {
		// Path format: /api/product-rules/draft/{index}
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		for i, part := range pathParts {
			if part == "draft" && i+1 < len(pathParts) {
				indexStr = pathParts[i+1]
				break
			}
		}
	}
	
	// Try mux.Vars as last fallback
	if indexStr == "" {
		vars := mux.Vars(r)
		if val, ok := vars["index"]; ok {
			indexStr = val
		}
	}
	
	if indexStr == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid draft index: index parameter not found")
		return
	}
	
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid draft index: "+err.Error())
		return
	}

	var draft model.ProductRuleDraft
	if err := json.NewDecoder(r.Body).Decode(&draft); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.usecase.UpdateDraft(index, &draft); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, draft)
}

// DeleteDraft handles DELETE /api/product-rules/draft/{index}
func (h *ProductRuleHandler) DeleteDraft(w http.ResponseWriter, r *http.Request) {
	// Try to get index from Echo context first (if available)
	var indexStr string
	if echoCtx := r.Context().Value("echo_context"); echoCtx != nil {
		if c, ok := echoCtx.(echo.Context); ok {
			indexStr = c.Param("index")
		}
	}
	
	// Fallback: Extract index from URL path
	if indexStr == "" {
		// Path format: /api/product-rules/draft/{index}
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		for i, part := range pathParts {
			if part == "draft" && i+1 < len(pathParts) {
				indexStr = pathParts[i+1]
				break
			}
		}
	}
	
	// Try mux.Vars as last fallback
	if indexStr == "" {
		vars := mux.Vars(r)
		if val, ok := vars["index"]; ok {
			indexStr = val
		}
	}
	
	if indexStr == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid draft index: index parameter not found")
		return
	}
	
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid draft index: "+err.Error())
		return
	}

	if err := h.usecase.DeleteDraft(index); err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Draft deleted successfully"})
}

// ClearDrafts handles POST /api/product-rules/draft/clear
func (h *ProductRuleHandler) ClearDrafts(w http.ResponseWriter, r *http.Request) {
	if err := h.usecase.ClearDrafts(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Drafts cleared successfully"})
}

// ConfirmDrafts handles POST /api/product-rules/draft/confirm
func (h *ProductRuleHandler) ConfirmDrafts(w http.ResponseWriter, r *http.Request) {
	// Get created_by from request body or default to 1
	var requestData struct {
		CreatedBy int64 `json:"created_by"`
	}
	
	// Try to decode request body, if empty use default
	json.NewDecoder(r.Body).Decode(&requestData)
	if requestData.CreatedBy == 0 {
		requestData.CreatedBy = 1 // Default value
	}

	if err := h.usecase.ConfirmDrafts(requestData.CreatedBy); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Product rules confirmed successfully"})
}

// GetAll handles GET /api/product-rules
func (h *ProductRuleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	productCode := r.URL.Query().Get("product_code")
	insuranceCode := r.URL.Query().Get("insurance_code")

	rules, err := h.usecase.GetAll(productCode, insuranceCode)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status": "Success",
		"data":   rules,
	})
}

// GetByID handles GET /api/product-rules/{id}
func (h *ProductRuleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Try to get id from Echo context first (if available)
	var idStr string
	if echoCtx := r.Context().Value("echo_context"); echoCtx != nil {
		if c, ok := echoCtx.(echo.Context); ok {
			idStr = c.Param("id")
		}
	}
	
	// Fallback: Extract id from URL path
	if idStr == "" {
		// Path format: /api/product-rules/{id}
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		for i, part := range pathParts {
			if part == "product-rules" && i+1 < len(pathParts) {
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
		respondWithError(w, http.StatusBadRequest, "Invalid product rule ID: id parameter not found")
		return
	}
	
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product rule ID: "+err.Error())
		return
	}

	rule, err := h.usecase.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, rule)
}

// Update handles PUT /api/product-rules/{id}
func (h *ProductRuleHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Try to get id from Echo context first (if available)
	var idStr string
	if echoCtx := r.Context().Value("echo_context"); echoCtx != nil {
		if c, ok := echoCtx.(echo.Context); ok {
			idStr = c.Param("id")
		}
	}
	
	// Fallback: Extract id from URL path
	if idStr == "" {
		// Path format: /api/product-rules/{id}
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		for i, part := range pathParts {
			if part == "product-rules" && i+1 < len(pathParts) {
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
		respondWithError(w, http.StatusBadRequest, "Invalid product rule ID: id parameter not found")
		return
	}
	
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product rule ID: "+err.Error())
		return
	}

	var rule model.ProductRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.usecase.Update(id, &rule); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, rule)
}

// Delete handles DELETE /api/product-rules/{id}
func (h *ProductRuleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Try to get id from Echo context first (if available)
	var idStr string
	if echoCtx := r.Context().Value("echo_context"); echoCtx != nil {
		if c, ok := echoCtx.(echo.Context); ok {
			idStr = c.Param("id")
		}
	}
	
	// Fallback: Extract id from URL path
	if idStr == "" {
		// Path format: /api/product-rules/{id}
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		for i, part := range pathParts {
			if part == "product-rules" && i+1 < len(pathParts) {
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
		respondWithError(w, http.StatusBadRequest, "Invalid product rule ID: id parameter not found")
		return
	}
	
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product rule ID: "+err.Error())
		return
	}

	if err := h.usecase.Delete(id); err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Product rule deleted successfully"})
}

// parseCSVContent parses CSV content and returns draft product rules
func parseCSVContent(content []byte) ([]model.ProductRuleDraft, error) {
	reader := csv.NewReader(strings.NewReader(string(content)))
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	expectedHeaders := []string{"product_code", "premium_type", "start_days", "end_days", "base_premium_value", "rules"}
	headerMap := make(map[string]int)
	for i, h := range header {
		headerMap[strings.TrimSpace(h)] = i
	}

	// Validate headers
	for _, expected := range expectedHeaders {
		if _, ok := headerMap[expected]; !ok {
			return nil, fmt.Errorf("missing required header: %s", expected)
		}
	}

	// Parse data rows
	var drafts []model.ProductRuleDraft
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

		startDays, _ := strconv.Atoi(strings.TrimSpace(getValue(row, headerMap, "start_days")))
		endDays, _ := strconv.Atoi(strings.TrimSpace(getValue(row, headerMap, "end_days")))
		basePremiumValue, _ := strconv.ParseFloat(strings.TrimSpace(getValue(row, headerMap, "base_premium_value")), 64)

		rulesValue := strings.TrimSpace(getValue(row, headerMap, "rules"))
		// Remove surrounding quotes if present
		if len(rulesValue) >= 2 && rulesValue[0] == '"' && rulesValue[len(rulesValue)-1] == '"' {
			rulesValue = rulesValue[1 : len(rulesValue)-1]
		}

		draft := model.ProductRuleDraft{
			ProductCode:      strings.TrimSpace(getValue(row, headerMap, "product_code")),
			PremiumType:      strings.TrimSpace(getValue(row, headerMap, "premium_type")),
			StartDays:        startDays,
			EndDays:          endDays,
			BasePremiumValue: basePremiumValue,
			Rules:            rulesValue,
		}

		drafts = append(drafts, draft)
	}

	return drafts, nil
}

// Helper functions for CSV parsing
func splitLines(s string) []string {
	var lines []string
	var current string
	for _, char := range s {
		if char == '\n' {
			lines = append(lines, current)
			current = ""
		} else if char != '\r' {
			current += string(char)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func parseCSVLine(line string) []string {
	var fields []string
	var current string
	inQuotes := false

	for i, char := range line {
		if char == '"' {
			inQuotes = !inQuotes
			// Don't add quote to current, but track state
		} else if char == ',' && !inQuotes {
			// Remove surrounding quotes if present
			field := current
			if len(field) >= 2 && field[0] == '"' && field[len(field)-1] == '"' {
				field = field[1 : len(field)-1]
			}
			fields = append(fields, field)
			current = ""
		} else {
			current += string(char)
		}
		if i == len(line)-1 {
			// Remove surrounding quotes if present
			field := current
			if len(field) >= 2 && field[0] == '"' && field[len(field)-1] == '"' {
				field = field[1 : len(field)-1]
			}
			fields = append(fields, field)
		}
	}

	return fields
}

func getValue(row []string, headerMap map[string]int, key string) string {
	if idx, ok := headerMap[key]; ok && idx < len(row) {
		return row[idx]
	}
	return ""
}

