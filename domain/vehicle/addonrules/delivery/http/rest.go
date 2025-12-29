package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/addonrules"
	"multi-onboarding/utils"
)

type handler struct {
	usecase addonrules.Usecase
}

// NewAddonRulesHandler creates a new addon rules handler
func NewAddonRulesHandler(usecase addonrules.Usecase) *handler {
	return &handler{
		usecase: usecase,
	}
}

// DownloadAddonRulesTemplate handles GET /api/addon-rules/download-template
func (h *handler) DownloadAddonRulesTemplate(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	insuranceCode := c.QueryParam("insurance_code")
	if insuranceCode == "" {
		return ac.CustomResponse("Failed", nil, "insurance_code is required", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	fileBytes, fileName, err := h.usecase.DownloadAddonRulesTemplate(c, insuranceCode)
	if err != nil {
		if err.Error() == "No active products found for this insurance" || err.Error() == "No active addons found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to generate template", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	c.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	return c.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}

// UploadAddonRulesFile handles POST /api/addon-rules/upload
func (h *handler) UploadAddonRulesFile(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	file, err := c.FormFile("file")
	if err != nil {
		return ac.CustomResponse("Failed", nil, "No file uploaded", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	src, err := file.Open()
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to open file", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}
	defer src.Close()

	insuranceCode := c.FormValue("insurance_code")

	newRulesCount, err := h.usecase.UploadAddonRulesFile(c, src, file.Filename, insuranceCode)
	if err != nil {
		if err.Error() == "File must be Excel (.xlsx) or CSV format" || err.Error() == "File must contain at least header and one data row" || err.Error() == "Invalid CSV format: missing required columns" || err.Error() == "No sheet found in Excel file" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to process file: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"message":         fmt.Sprintf("Successfully uploaded %d addon rules", newRulesCount),
		"new_rules_count": newRulesCount,
	}, "File uploaded successfully", http.StatusOK, "", nil)
}

// GetAddonRulesDraft handles GET /api/addon-rules/draft
func (h *handler) GetAddonRulesDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	insuranceCode := c.QueryParam("insurance_code")

	rulesList, err := h.usecase.GetAddonRulesDraft(c, insuranceCode)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to load draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", rulesList, "Success to Get Addon Rules Draft", http.StatusOK, "", nil)
}

// UpdateAddonRuleDraft handles PUT /api/addon-rules/draft/:id
func (h *handler) UpdateAddonRuleDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid ID", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	var updateData struct {
		ProductCode   string `json:"product_code"`
		InsuranceCode string `json:"insurance_code"`
		Rules         string `json:"rules"`
	}

	if err := ac.Bind(&updateData); err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid request: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	req := addonrules.UpdateAddonRuleDraftRequest{
		ID:            id,
		ProductCode:   updateData.ProductCode,
		InsuranceCode: updateData.InsuranceCode,
		Rules:         updateData.Rules,
		CreatedBy:     1,
	}

	err = h.usecase.UpdateAddonRuleDraft(c, req)
	if err != nil {
		if err.Error() == "Addon rule draft not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to save draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Addon rule updated successfully", http.StatusOK, "", nil)
}

// DeleteAddonRuleDraft handles DELETE /api/addon-rules/draft/:id
func (h *handler) DeleteAddonRuleDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid ID", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	err = h.usecase.DeleteAddonRuleDraft(c, id)
	if err != nil {
		if err.Error() == "Addon rule draft not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to save draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Addon rule deleted successfully", http.StatusOK, "", nil)
}

// ClearAddonRulesDraft handles POST /api/addon-rules/draft/clear
func (h *handler) ClearAddonRulesDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	err := h.usecase.ClearAddonRulesDraft(c)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to save draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Addon rule drafts cleared successfully", http.StatusOK, "", nil)
}

// ConfirmAddonRules handles POST /api/addon-rules/draft/confirm
func (h *handler) ConfirmAddonRules(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	result, err := h.usecase.ConfirmAddonRules(c)
	if err != nil {
		if err.Error() == "No draft addon rules to confirm" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to confirm addon rules: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"message": fmt.Sprintf("Successfully confirmed %d addon rules", result.InsertedCount),
		"count":   result.InsertedCount,
	}, "Addon rules confirmed successfully", http.StatusOK, "", nil)
}

// GetAddonRules handles GET /api/addon-rules
func (h *handler) GetAddonRules(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	insuranceCode := c.QueryParam("insurance_code")

	rulesList, err := h.usecase.GetAddonRules(c, insuranceCode)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to query database", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", rulesList, "Success to Get Addon Rules", http.StatusOK, "", nil)
}

// UpdateAddonRule handles PUT /api/addon-rules/:id
func (h *handler) UpdateAddonRule(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid ID", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	var updateData struct {
		ProductCode   string `json:"product_code"`
		InsuranceCode string `json:"insurance_code"`
		Rules         string `json:"rules"`
	}

	if err := ac.Bind(&updateData); err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid request: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	// Get user name from header
	userName := utils.GetUserNameFromRequest(c)

	req := addonrules.UpdateAddonRuleRequest{
		UpdateData: map[string]interface{}{
			"product_code":   updateData.ProductCode,
			"insurance_code": updateData.InsuranceCode,
			"rules":          updateData.Rules,
		},
	}

	err = h.usecase.UpdateAddonRule(c, id, req, userName)
	if err != nil {
		if err.Error() == "Addon rule not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to update addon rule", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Addon rule updated successfully", http.StatusOK, "", nil)
}

// AddAddonRulesHandler registers addon rules routes to Echo router
func AddAddonRulesHandler(e *echo.Echo, usecase addonrules.Usecase) {
	handler := &handler{
		usecase: usecase,
	}

	// Onboarding routes
	e.GET("/api/addon-rules/download-template", handler.DownloadAddonRulesTemplate)
	e.POST("/api/addon-rules/upload", handler.UploadAddonRulesFile)
	e.GET("/api/addon-rules/draft", handler.GetAddonRulesDraft)
	e.PUT("/api/addon-rules/draft/:id", handler.UpdateAddonRuleDraft)
	e.DELETE("/api/addon-rules/draft/:id", handler.DeleteAddonRuleDraft)
	e.POST("/api/addon-rules/draft/clear", handler.ClearAddonRulesDraft)
	e.POST("/api/addon-rules/draft/confirm", handler.ConfirmAddonRules)
	e.GET("/api/addon-rules", handler.GetAddonRules)
	e.PUT("/api/addon-rules/:id", handler.UpdateAddonRule)
}

