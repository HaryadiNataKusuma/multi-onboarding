package http

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/template"
	"multi-onboarding/utils"
)

type handlerTemplate struct {
	usecase template.Usecase
}

// NewTemplateHandler creates a new template handler
func NewTemplateHandler(usecase template.Usecase) *handlerTemplate {
	return &handlerTemplate{
		usecase: usecase,
	}
}

// AddTemplateHandler registers template routes to Echo router
func AddTemplateHandler(e *echo.Echo, usecase template.Usecase) {
	handler := handlerTemplate{
		usecase: usecase,
	}

	// Onboarding routes
	e.GET("/api/templates", handler.GetTemplates)
	e.PUT("/api/templates/:id", handler.UpdateTemplate)
	e.GET("/api/templates/draft", handler.GetTemplatesDraft)
	e.PUT("/api/templates/draft/:id", handler.UpdateTemplateDraft)
	e.DELETE("/api/templates/draft/:id", handler.DeleteTemplateDraft)
	e.POST("/api/templates/draft/confirm", handler.ConfirmTemplates)
	e.POST("/api/templates/draft/clear", handler.ClearTemplatesDraft)
}

// GetTemplates handles GET /api/templates
func (h *handlerTemplate) GetTemplates(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	insuranceCode := c.QueryParam("insurance_code")

	templates, err := h.usecase.GetTemplates(c, insuranceCode)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to query database: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", templates, "Success to Get Templates", http.StatusOK, "", nil)
}

// GetTemplatesDraft handles GET /api/templates/draft
func (h *handlerTemplate) GetTemplatesDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	insuranceCode := c.QueryParam("insurance_code")

	templates, err := h.usecase.GetTemplatesDraft(c, insuranceCode)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to load draft data: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	// Ensure data is always an array, not null
	if templates == nil {
		templates = []utils.TemplateDraft{}
	}

	// Transform response: template_id -> id, remove insurance field
	transformedTemplates := make([]map[string]interface{}, 0)
	for _, tmpl := range templates {
		// Skip templates with empty TemplateID
		if strings.TrimSpace(tmpl.TemplateID) == "" {
			continue
		}

		// Value adalah array, default ["","",""]
		value := tmpl.Value
		if value == nil || len(value) == 0 {
			value = []string{"", "", ""}
		}

		transformedTemplates = append(transformedTemplates, map[string]interface{}{
			"id":         tmpl.TemplateID,
			"locale":     tmpl.Locale,
			"value":      value,
			"created_by": tmpl.CreatedBy,
			"created_at": tmpl.CreatedAt,
		})
	}

	return ac.CustomResponse("Success", transformedTemplates, "Success to Get Templates Draft", http.StatusOK, "", nil)
}

// UpdateTemplate handles PUT /api/templates/:id
func (h *handlerTemplate) UpdateTemplate(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	idStr := c.Param("id")
	if idStr == "" {
		return ac.CustomResponse("Failed", nil, "ID parameter is required", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	templateID := idStr

	var input struct {
		Locale    string `json:"locale" validate:"required"`
		ID        string `json:"id" validate:"required"`
		Value     string `json:"value" validate:"required"`
		CreatedBy int    `json:"created_by"`
	}

	if err := ac.Bind(&input); err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid request: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	if err := ac.Validate(input); err != nil {
		return ac.CustomResponse("Failed", nil, "Validation error: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
	}

	// Get user name from request
	userName := utils.GetUserNameFromRequest(c)

	req := template.UpdateTemplateRequest{
		Locale:    input.Locale,
		ID:        input.ID,
		Value:     input.Value,
		CreatedBy: input.CreatedBy,
	}

	if err := h.usecase.UpdateTemplate(c, templateID, req, userName); err != nil {
		if err.Error() == "Template not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to update template", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"id": templateID,
	}, "Template updated successfully", http.StatusOK, "", nil)
}

// UpdateTemplateDraft handles PUT /api/templates/draft/:id
func (h *handlerTemplate) UpdateTemplateDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	templateID := c.Param("id")
	if templateID == "" {
		return ac.CustomResponse("Failed", nil, "Template ID parameter is required", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	var input struct {
		Locale    string   `json:"locale" validate:"required"`
		Value     []string `json:"value" validate:"required"`
		CreatedBy int      `json:"created_by"`
	}

	if err := ac.Bind(&input); err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid request: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	if err := ac.Validate(input); err != nil {
		return ac.CustomResponse("Failed", nil, "Validation error: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
	}

	req := template.UpdateTemplateDraftRequest{
		ID:        templateID,
		Locale:    input.Locale,
		Value:     input.Value,
		CreatedBy: input.CreatedBy,
	}

	err := h.usecase.UpdateTemplateDraft(c, req)
	if err != nil {
		if err.Error() == "Template draft not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to save draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"id": templateID,
	}, "Template draft updated successfully", http.StatusOK, "", nil)
}

// DeleteTemplateDraft handles DELETE /api/templates/draft/:id
func (h *handlerTemplate) DeleteTemplateDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	templateID := c.Param("id")
	if templateID == "" {
		return ac.CustomResponse("Failed", nil, "Template ID parameter is required", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	log.Printf("DeleteTemplateDraft called: method=%s, path=%s, templateID=%s", c.Request().Method, c.Request().URL.Path, templateID)

	err := h.usecase.DeleteTemplateDraft(c, templateID)
	if err != nil {
		log.Printf("Error deleting template draft %s: %v", templateID, err)
		if err.Error() == "Template draft not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to delete template draft: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	log.Printf("Successfully deleted template draft: %s", templateID)
	return ac.CustomResponse("Success", nil, "Template draft deleted successfully", http.StatusOK, "", nil)
}

// ConfirmTemplates handles POST /api/templates/draft/confirm
func (h *handlerTemplate) ConfirmTemplates(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	insuranceCodeFilter := c.QueryParam("insurance_code")

	result, err := h.usecase.ConfirmTemplates(c, insuranceCodeFilter)
	if err != nil {
		if err.Error() == "No template drafts available to confirm" || err.Error() == "No templates found for the specified insurance code in draft" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
		}
		if err.Error() == "Failed to confirm any templates" {
			responseData := map[string]interface{}{
				"success_count": result.SuccessCount,
				"total_count":   result.TotalCount,
			}
			if result.ErrorCount > 0 {
				responseData["error_count"] = result.ErrorCount
				responseData["errors"] = result.Errors
			}
			return ac.CustomResponse("Failed", responseData, "Templates confirmation completed with errors", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to confirm templates", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	responseData := map[string]interface{}{
		"success_count": result.SuccessCount,
		"total_count":   result.TotalCount,
	}
	if result.ErrorCount > 0 {
		responseData["error_count"] = result.ErrorCount
		responseData["errors"] = result.Errors
	}

	return ac.CustomResponse("Success", responseData, "Templates confirmation completed", http.StatusOK, "", nil)
}

// ClearTemplatesDraft handles POST /api/templates/draft/clear
func (h *handlerTemplate) ClearTemplatesDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	err := h.usecase.ClearTemplatesDraft(c)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to save draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Template drafts cleared successfully", http.StatusOK, "", nil)
}


