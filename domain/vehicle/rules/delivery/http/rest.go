package http

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/rules"
	"multi-onboarding/utils"
)

// Helper function to convert echo.Context to CustomApplicationContextStub
func toCustomContext(c echo.Context) *utils.CustomApplicationContextStub {
	return utils.GetCustomApplicationContextStub(c)
}

type handler struct {
	usecase rules.Usecase
}

// NewRulesHandler creates a new rules handler
func NewRulesHandler(usecase rules.Usecase) *handler {
	return &handler{
		usecase: usecase,
	}
}

// AddRulesHandler registers rules routes to Echo router
func AddRulesHandler(e *echo.Echo, usecase rules.Usecase) {
	handler := &handler{
		usecase: usecase,
	}

	// Onboarding routes
	e.GET("/api/product-rules/download-template", handler.DownloadRulesTemplate)
	e.POST("/api/product-rules/upload", handler.UploadRulesFile)
	e.GET("/api/product-rules/draft", handler.GetRulesDraft)
	e.PUT("/api/product-rules/draft/:id", handler.UpdateRuleDraft)
	e.DELETE("/api/product-rules/draft/:id", handler.DeleteRuleDraft)
	e.POST("/api/product-rules/draft/clear", handler.ClearRulesDraft)
	e.POST("/api/product-rules/draft/confirm", handler.ConfirmRules)
	e.PUT("/api/product-rules/draft/bulk", handler.BulkUpdateRulesDraft)
	e.GET("/api/product-rules", handler.GetRules)
	e.PUT("/api/product-rules/:id", handler.UpdateRule)
	e.DELETE("/api/product-rules/:id", handler.DeleteRule)
	e.PUT("/api/product-rules/bulk", handler.BulkUpdateRules)
}

// DownloadRulesTemplate handles GET /api/product-rules/download-template
func (h *handler) DownloadRulesTemplate(c echo.Context) error {
	ac := toCustomContext(c)

	insuranceCode := c.QueryParam("insurance_code")
	log.Printf("DownloadRulesTemplate called with insurance_code: %s", insuranceCode)

	if insuranceCode == "" {
		log.Printf("ERROR: insurance_code is required")
		return ac.CustomResponse("Failed", nil, "insurance_code is required", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	fileBytes, fileName, err := h.usecase.DownloadRulesTemplate(c, insuranceCode)
	if err != nil {
		log.Printf("ERROR: Failed to generate template for insurance_code %s: %v", insuranceCode, err)
		if err.Error() == "No active products found for this insurance" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to generate template: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	log.Printf("SUCCESS: Generated template file %s (%d bytes) for insurance_code %s", fileName, len(fileBytes), insuranceCode)

	// Set headers and send file using Blob method
	c.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	return c.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}

// UploadRulesFile handles POST /api/product-rules/upload
func (h *handler) UploadRulesFile(c echo.Context) error {
	ac := toCustomContext(c)
	log.Printf("UploadRulesFile called")

	// Get file from form
	file, err := c.FormFile("file")

	if err != nil {
		log.Printf("ERROR: Failed to get file from form: %v", err)
		return ac.CustomResponse("Failed", nil, "File is required", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	src, err := file.Open()
	if err != nil {
		log.Printf("ERROR: Failed to open file: %v", err)
		return ac.CustomResponse("Failed", nil, "Failed to open file", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}
	defer src.Close()

	// Get insurance_code from form
	insuranceCode := c.FormValue("insurance_code")

	log.Printf("UploadRulesFile: file=%s, insurance_code=%s", file.Filename, insuranceCode)

	newRulesCount, err := h.usecase.UploadRulesFile(c, src, file.Filename, insuranceCode)
	if err != nil {
		if err.Error() == "File must be Excel (.xlsx) or CSV format" || err.Error() == "File must contain at least header and one data row" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to process file: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"message":         fmt.Sprintf("Successfully uploaded %d rules", newRulesCount),
		"new_rules_count": newRulesCount,
	}, "File uploaded successfully", http.StatusOK, "", nil)
}

// GetRulesDraft handles GET /api/product-rules/draft
func (h *handler) GetRulesDraft(c echo.Context) error {
	ac := toCustomContext(c)
	insuranceCode := c.QueryParam("insurance_code")

	rulesList, err := h.usecase.GetRulesDraft(c, insuranceCode)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to load draft data: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	// Ensure data is always an array, not null
	if rulesList == nil {
		rulesList = []utils.RuleDraft{}
	}

	return ac.CustomResponse("Success", rulesList, "Success to Get Rules Draft", http.StatusOK, "", nil)
}

// UpdateRuleDraft handles PUT /api/product-rules/draft/:id
func (h *handler) UpdateRuleDraft(c echo.Context) error {
	ac := toCustomContext(c)

	ruleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid rule ID", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	var rule struct {
		ProductCode             string   `json:"product_code"`
		InsuranceCode           string   `json:"insurance_code"`
		VehicleType             string   `json:"vehicle_type"`
		VehicleCategory         string   `json:"vehicle_category"`
		StartVehicleValue       *float64 `json:"start_vehicle_value"`
		EndVehicleValue         *float64 `json:"end_vehicle_value"`
		LimitVehicleAge         *int     `json:"limit_vehicle_age"`
		AdminFee                *float64 `json:"admin_fee"`
		HardcopyAdminFee        *float64 `json:"hardcopy_admin_fee"`
		RegionID                *int     `json:"region_id"`
		BasePremiumValue        *float64 `json:"base_premium_value"`
		LoadingFeePremiumValue  *float64 `json:"loading_fee_premium_value"`
		CommercialUsageValue    *float64 `json:"commercial_usage_value"`
		StartLoadingAge         *int     `json:"start_loading_age"`
		AdditionalPremium       *float64 `json:"additional_premium"`
		TypeAdditionalPremium   string   `json:"type_additional_premium"`
		Rules                   string   `json:"rules"`
		BasePremiumType         string   `json:"base_premium_type"`
		BaseLoadingPremiumValue *float64 `json:"base_loading_premium_value"`
	}

	if err := ac.Bind(&rule); err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid request: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	req := rules.UpdateRuleDraftRequest{
		ID:                      ruleID,
		ProductCode:             rule.ProductCode,
		InsuranceCode:           rule.InsuranceCode,
		VehicleType:             rule.VehicleType,
		VehicleCategory:         rule.VehicleCategory,
		StartVehicleValue:       rule.StartVehicleValue,
		EndVehicleValue:         rule.EndVehicleValue,
		LimitVehicleAge:         rule.LimitVehicleAge,
		AdminFee:                rule.AdminFee,
		HardcopyAdminFee:        rule.HardcopyAdminFee,
		RegionID:                rule.RegionID,
		BasePremiumValue:        rule.BasePremiumValue,
		LoadingFeePremiumValue:  rule.LoadingFeePremiumValue,
		CommercialUsageValue:    rule.CommercialUsageValue,
		StartLoadingAge:         rule.StartLoadingAge,
		AdditionalPremium:       rule.AdditionalPremium,
		TypeAdditionalPremium:   rule.TypeAdditionalPremium,
		Rules:                   rule.Rules,
		BasePremiumType:         rule.BasePremiumType,
		BaseLoadingPremiumValue: rule.BaseLoadingPremiumValue,
	}

	err = h.usecase.UpdateRuleDraft(c, req)
	if err != nil {
		if err.Error() == "Rule draft not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to save draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Rule updated successfully", http.StatusOK, "", nil)
}

// DeleteRuleDraft handles DELETE /api/product-rules/draft/:id
func (h *handler) DeleteRuleDraft(c echo.Context) error {
	ac := toCustomContext(c)

	ruleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid rule ID", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	err = h.usecase.DeleteRuleDraft(c, ruleID)
	if err != nil {
		if err.Error() == "Rule draft not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to save draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Rule deleted successfully", http.StatusOK, "", nil)
}

// ClearRulesDraft handles POST /api/product-rules/draft/clear
func (h *handler) ClearRulesDraft(c echo.Context) error {
	ac := toCustomContext(c)

	err := h.usecase.ClearRulesDraft(c)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to save draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Rule drafts cleared successfully", http.StatusOK, "", nil)
}

// ConfirmRules handles POST /api/product-rules/draft/confirm
func (h *handler) ConfirmRules(c echo.Context) error {
	ac := toCustomContext(c)

	result, err := h.usecase.ConfirmRules(c)
	if err != nil {
		if err.Error() == "No draft rules to confirm" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to confirm rules: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"message":        fmt.Sprintf("Successfully confirmed %d rules", result.InsertedCount),
		"inserted_count": result.InsertedCount,
	}, "Rules confirmed successfully", http.StatusOK, "", nil)
}

// BulkUpdateRulesDraft handles PUT /api/product-rules/draft/bulk
func (h *handler) BulkUpdateRulesDraft(c echo.Context) error {
	ac := toCustomContext(c)

	var req rules.BulkUpdateRulesDraftRequest
	if err := ac.Bind(&req); err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid request: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	result, err := h.usecase.BulkUpdateRulesDraft(c, req)
	if err != nil {
		if err.Error() == "No rules found matching the filter criteria" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to update rules: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"message":       fmt.Sprintf("Successfully updated %d rules", result.UpdatedCount),
		"updated_count": result.UpdatedCount,
	}, "Rules updated successfully", http.StatusOK, "", nil)
}

// GetRules handles GET /api/product-rules
func (h *handler) GetRules(c echo.Context) error {
	ac := toCustomContext(c)

	insuranceCode := c.QueryParam("insurance_code")

	rulesList, err := h.usecase.GetRules(c, insuranceCode)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to query database", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", rulesList, "Success to Get Rules", http.StatusOK, "", nil)
}

// UpdateRule handles PUT /api/product-rules/:id
func (h *handler) UpdateRule(c echo.Context) error {
	ac := toCustomContext(c)

	ruleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid rule ID", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	var updateData map[string]interface{}
	if err := ac.Bind(&updateData); err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid request: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	// Get user name from request
	userName := utils.GetUserNameFromRequest(c)

	req := rules.UpdateRuleRequest{
		UpdateData: updateData,
	}

	err = h.usecase.UpdateRule(c, ruleID, req, userName)
	if err != nil {
		if err.Error() == "Rule not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to update rule", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Rule updated successfully", http.StatusOK, "", nil)
}

// DeleteRule handles DELETE /api/product-rules/:id
func (h *handler) DeleteRule(c echo.Context) error {
	ac := toCustomContext(c)

	ruleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid rule ID", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	err = h.usecase.DeleteRule(c, ruleID)
	if err != nil {
		if err.Error() == "Rule not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to delete rule", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Rule deleted successfully", http.StatusOK, "", nil)
}

// BulkUpdateRules handles PUT /api/product-rules/bulk
func (h *handler) BulkUpdateRules(c echo.Context) error {
	ac := toCustomContext(c)

	var req rules.BulkUpdateRulesRequest
	if err := ac.Bind(&req); err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid request: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	result, err := h.usecase.BulkUpdateRules(c, req)
	if err != nil {
		if err.Error() == "At least one filter condition is required" || err.Error() == "At least one update field is required" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to update rules: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"message":       fmt.Sprintf("Successfully updated %d rules", result.UpdatedCount),
		"updated_count": result.UpdatedCount,
	}, "Rules updated successfully", http.StatusOK, "", nil)
}

