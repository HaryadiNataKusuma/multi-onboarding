package http

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	model "multi-onboarding/model/Travel"
	usecase "multi-onboarding/usecase/Travel"
	"multi-onboarding/utils"

	"github.com/gorilla/mux"
)

// TemplateHandler handles HTTP requests for templates
type TemplateHandler struct {
	usecase *usecase.TemplateUsecase
}

// NewTemplateHandler creates a new template handler
func NewTemplateHandler(usecase *usecase.TemplateUsecase) *TemplateHandler {
	return &TemplateHandler{usecase: usecase}
}

// AddDraft handles POST /api/templates/draft/add
func (h *TemplateHandler) AddDraft(w http.ResponseWriter, r *http.Request) {
	var draft model.TemplateDraft
	if err := json.NewDecoder(r.Body).Decode(&draft); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.usecase.AddDraft(&draft); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, draft)
}

// GetAllDrafts handles GET /api/templates/draft
func (h *TemplateHandler) GetAllDrafts(w http.ResponseWriter, r *http.Request) {
	insuranceCode := r.URL.Query().Get("insurance_code")

	drafts, err := h.usecase.GetAllDrafts(insuranceCode)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status": "Success",
		"data":   drafts,
	})
}

// GetDraftByID handles GET /api/templates/draft/{id}
func (h *TemplateHandler) GetDraftByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "id" // Default locale
	}

	draft, err := h.usecase.GetDraftByID(id, locale)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, draft)
}

// UpdateDraft handles PUT /api/templates/draft/{id}
func (h *TemplateHandler) UpdateDraft(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var draft model.TemplateDraft
	if err := json.NewDecoder(r.Body).Decode(&draft); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Use locale from request body, default to "id"
	locale := draft.Locale
	if locale == "" {
		locale = "id" // Default locale
	}

	if err := h.usecase.UpdateDraft(id, locale, &draft); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, draft)
}

// DeleteDraft handles DELETE /api/templates/draft/{id}
func (h *TemplateHandler) DeleteDraft(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "id" // Default locale
	}

	if err := h.usecase.DeleteDraft(id, locale); err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Draft deleted successfully"})
}

// ClearDrafts handles POST /api/templates/draft/clear
func (h *TemplateHandler) ClearDrafts(w http.ResponseWriter, r *http.Request) {
	if err := h.usecase.ClearDrafts(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Drafts cleared successfully"})
}

// ConfirmDrafts handles POST /api/templates/draft/confirm
func (h *TemplateHandler) ConfirmDrafts(w http.ResponseWriter, r *http.Request) {
	if err := h.usecase.ConfirmDrafts(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Templates confirmed successfully"})
}

// GetAll handles GET /api/templates
func (h *TemplateHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	insuranceCode := r.URL.Query().Get("insurance_code")

	// Query templates with JOIN to products to get insurance_code and product_name
	var query string
	var args []interface{}

	if insuranceCode != "" {
		query = `
			SELECT DISTINCT 
				t.locale, t.id, t.value, t.created_by, t.created_at, t.updated_by, t.updated_at,
				p.insurance_code, p.name as product_name
			FROM travel_service_development.templates t
			INNER JOIN travel_service_development.products p ON (t.id = p.summary OR t.id = p.insurance_detail)
			WHERE p.insurance_code = ?
			ORDER BY t.created_at DESC
		`
		args = []interface{}{insuranceCode}
	} else {
		query = `
			SELECT DISTINCT 
				t.locale, t.id, t.value, t.created_by, t.created_at, t.updated_by, t.updated_at,
				p.insurance_code, p.name as product_name
			FROM travel_service_development.templates t
			INNER JOIN travel_service_development.products p ON (t.id = p.summary OR t.id = p.insurance_detail)
			ORDER BY t.created_at DESC
		`
		args = []interface{}{}
	}

	// Get database connection
	db := utils.GetDB()
	rows, err := db.Query(query, args...)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to query templates: "+err.Error())
		return
	}
	defer rows.Close()

	enrichedTemplates := make([]map[string]interface{}, 0)
	for rows.Next() {
		var locale, id, value string
		var createdBy int64
		var createdAt sql.NullTime
		var updatedBy sql.NullInt64
		var updatedAt sql.NullTime
		var insCode sql.NullString
		var prodName sql.NullString

		if err := rows.Scan(&locale, &id, &value, &createdBy, &createdAt, &updatedBy, &updatedAt, &insCode, &prodName); err != nil {
			continue
		}

		enriched := map[string]interface{}{
			"locale":     locale,
			"id":         id,
			"value":      value,
			"created_by": createdBy,
		}

		if createdAt.Valid {
			enriched["created_at"] = createdAt.Time
		}
		if updatedBy.Valid {
			enriched["updated_by"] = updatedBy.Int64
		}
		if updatedAt.Valid {
			enriched["updated_at"] = updatedAt.Time
		}

		if insCode.Valid {
			enriched["insurance_code"] = insCode.String
		} else {
			enriched["insurance_code"] = ""
		}

		if prodName.Valid {
			enriched["product_name"] = prodName.String
		} else {
			enriched["product_name"] = ""
		}

		enrichedTemplates = append(enrichedTemplates, enriched)
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status": "Success",
		"data":   enrichedTemplates,
	})
}

// GetByID handles GET /api/templates/{id}
func (h *TemplateHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "id" // Default locale
	}

	template, err := h.usecase.GetByID(id, locale)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, template)
}

// Update handles PUT /api/templates/{id}
func (h *TemplateHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var template model.Template
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Use locale from request body, default to "id"
	locale := template.Locale
	if locale == "" {
		locale = "id" // Default locale
	}

	if err := h.usecase.Update(id, locale, &template); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, template)
}

// Delete handles DELETE /api/templates/{id}
func (h *TemplateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "id" // Default locale
	}

	if err := h.usecase.Delete(id, locale); err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Template deleted successfully"})
}

// UploadExcel handles POST /api/templates/upload
func (h *TemplateHandler) UploadExcel(w http.ResponseWriter, r *http.Request) {
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

	// Get insurance_code from form if provided (from selectedInsuranceCode in frontend)
	selectedInsuranceCode := r.FormValue("insurance_code")

	// Read file content
	fileBytes := make([]byte, header.Size)
	_, err = file.Read(fileBytes)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to read file: "+err.Error())
		return
	}

	// Parse CSV content
	drafts, err := parseTemplateCSV(fileBytes, selectedInsuranceCode)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse file: "+err.Error())
		return
	}

	// Clear existing drafts
	if err := h.usecase.ClearDrafts(); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to clear drafts: "+err.Error())
		return
	}

	// Add all drafts
	successCount := 0
	failedDrafts := []string{}
	for _, draft := range drafts {
		if err := h.usecase.AddDraft(&draft); err != nil {
			failedDrafts = append(failedDrafts, fmt.Sprintf("%s: %v", draft.ID, err))
			continue // Continue with other drafts instead of failing completely
		}
		successCount++
	}

	if successCount == 0 && len(failedDrafts) > 0 {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to add all drafts: %v", failedDrafts))
		return
	}

	// Return success even if some drafts failed (partial success)
	responseMessage := fmt.Sprintf("Successfully uploaded %d template(s)", successCount)
	if len(failedDrafts) > 0 {
		responseMessage += fmt.Sprintf(". Failed: %d", len(failedDrafts))
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message":       responseMessage,
		"data":          drafts,
		"success_count": successCount,
		"failed_count":  len(failedDrafts),
	})
}

// parseTemplateCSV parses CSV content and returns template drafts
// selectedInsuranceCode is used as fallback if CSV doesn't have insurance_code
func parseTemplateCSV(content []byte, selectedInsuranceCode string) ([]model.TemplateDraft, error) {
	reader := csv.NewReader(strings.NewReader(string(content)))
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Required headers only: id, locale, value
	requiredHeaders := []string{"id", "locale", "value"}
	// Optional headers (not validated): product_code, insurance_code

	headerMap := make(map[string]int)
	for i, h := range header {
		normalized := strings.TrimSpace(strings.ToLower(h))
		headerMap[normalized] = i
	}

	// Validate required headers only (case-insensitive)
	missingHeaders := []string{}
	for _, expected := range requiredHeaders {
		if _, ok := headerMap[expected]; !ok {
			missingHeaders = append(missingHeaders, expected)
		}
	}
	if len(missingHeaders) > 0 {
		return nil, fmt.Errorf("missing required header(s): %s. Found headers: %v", strings.Join(missingHeaders, ", "), header)
	}

	// Optional headers are not validated - they can be missing

	// Helper function to get value from CSV row
	getValue := func(row []string, headerMap map[string]int, key string) string {
		if idx, ok := headerMap[key]; ok && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	// Parse data rows
	var drafts []model.TemplateDraft
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // Skip invalid rows
		}

		// Skip rows that don't have at least the required fields
		if len(row) < len(requiredHeaders) {
			continue // Skip incomplete rows
		}

		draft := model.TemplateDraft{
			ID:            getValue(row, headerMap, "id"),
			Locale:        getValue(row, headerMap, "locale"),
			Value:         getValue(row, headerMap, "value"),
			ProductCode:   getValue(row, headerMap, "product_code"),
			InsuranceCode: getValue(row, headerMap, "insurance_code"),
			CreatedBy:     1, // Default created_by
		}

		// Set defaults
		if draft.Locale == "" {
			draft.Locale = "id" // Default locale
		}
		if draft.ID == "" {
			continue // Skip rows without ID
		}

		// If insurance_code is empty in CSV, use selectedInsuranceCode from form (if provided)
		if draft.InsuranceCode == "" && selectedInsuranceCode != "" {
			draft.InsuranceCode = selectedInsuranceCode
		}

		drafts = append(drafts, draft)
	}

	return drafts, nil
}
