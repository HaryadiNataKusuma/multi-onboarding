package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/insuranceownrisks"
	"multi-onboarding/utils"
)

type handler struct {
	usecase insuranceownrisks.Usecase
}

// NewInsuranceOwnRisksHandler creates a new insurance own risks handler
func NewInsuranceOwnRisksHandler(usecase insuranceownrisks.Usecase) *handler {
	return &handler{
		usecase: usecase,
	}
}

// AddInsuranceOwnRisksHandler registers insurance own risks routes to Echo router
func AddInsuranceOwnRisksHandler(e *echo.Echo, usecase insuranceownrisks.Usecase) {
	handler := &handler{
		usecase: usecase,
	}

	// Onboarding routes
	e.GET("/api/insurance-own-risks/download-template", handler.DownloadOwnRisksTemplate)
	e.POST("/api/insurance-own-risks/upload", handler.UploadOwnRisksFile)
	e.GET("/api/insurance-own-risks/draft", handler.GetOwnRisksDraft)
	e.PUT("/api/insurance-own-risks/draft/:id", handler.UpdateOwnRisksDraft)
	e.DELETE("/api/insurance-own-risks/draft/:id", handler.DeleteOwnRisksDraft)
	e.POST("/api/insurance-own-risks/draft/clear", handler.ClearOwnRisksDraft)
	e.POST("/api/insurance-own-risks/draft/confirm", handler.ConfirmOwnRisks)
	e.GET("/api/insurance-own-risks", handler.GetOwnRisks)
	e.PUT("/api/insurance-own-risks/:id", handler.UpdateOwnRisks)
}

// DownloadOwnRisksTemplate handles GET /api/insurance-own-risks/download-template
func (h *handler) DownloadOwnRisksTemplate(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	insuranceCode := c.QueryParam("insurance_code")
	if insuranceCode == "" {
		return ac.CustomResponse("Failed", nil, "insurance_code is required", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	fileBytes, fileName, err := h.usecase.DownloadOwnRisksTemplate(c, insuranceCode)
	if err != nil {
		if err.Error() == "No active products found for this insurance" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to generate template", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	c.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	return c.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}

// UploadOwnRisksFile handles POST /api/insurance-own-risks/upload
func (h *handler) UploadOwnRisksFile(c echo.Context) error {
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

	rowCount, err := h.usecase.UploadOwnRisksFile(c, src, file.Filename, insuranceCode)
	if err != nil {
		if err.Error() == "No sheet found in Excel file" || err.Error() == "Failed to read file" || strings.Contains(err.Error(), "Invalid headers. Expected") {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to process file: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"message":    fmt.Sprintf("File uploaded successfully. %d rows processed.", rowCount),
		"rows_count": rowCount,
	}, "File uploaded successfully", http.StatusOK, "", nil)
}

// GetOwnRisksDraft handles GET /api/insurance-own-risks/draft
func (h *handler) GetOwnRisksDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	insuranceCode := c.QueryParam("insurance_code")

	rulesList, err := h.usecase.GetOwnRisksDraft(c, insuranceCode)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to load draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", rulesList, "Success to Get Own Risks Draft", http.StatusOK, "", nil)
}

// UpdateOwnRisksDraft handles PUT /api/insurance-own-risks/draft/:id
func (h *handler) UpdateOwnRisksDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid ID", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	var updateData struct {
		InsuranceCode     string   `json:"insurance_code"`
		ProductType       string   `json:"product_type"`
		Code              string   `json:"code"`
		Title             string   `json:"title"`
		Value             *float64 `json:"value"`
		ValueType         string   `json:"value_type"`
		Description       *string  `json:"description"`
		IsMandatory       *bool    `json:"is_mandatory"`
		IsActive          *bool    `json:"is_active"`
		IsElectricVehicle *int     `json:"is_electric_vehicle"`
	}

	if err := ac.Bind(&updateData); err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid request: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	req := insuranceownrisks.UpdateOwnRiskDraftRequest{
		ID:                id,
		InsuranceCode:     updateData.InsuranceCode,
		ProductType:       updateData.ProductType,
		Code:              updateData.Code,
		Title:             updateData.Title,
		Value:             updateData.Value,
		ValueType:         updateData.ValueType,
		Description:       updateData.Description,
		IsMandatory:       updateData.IsMandatory,
		IsActive:          updateData.IsActive,
		IsElectricVehicle: updateData.IsElectricVehicle,
	}

	err = h.usecase.UpdateOwnRiskDraft(c, req)
	if err != nil {
		if err.Error() == "Own risk draft not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to save draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Own risk updated successfully", http.StatusOK, "", nil)
}

// DeleteOwnRisksDraft handles DELETE /api/insurance-own-risks/draft/:id
func (h *handler) DeleteOwnRisksDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid ID", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	err = h.usecase.DeleteOwnRiskDraft(c, id)
	if err != nil {
		if err.Error() == "Own risk draft not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to save draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Own risk deleted successfully", http.StatusOK, "", nil)
}

// ClearOwnRisksDraft handles POST /api/insurance-own-risks/draft/clear
func (h *handler) ClearOwnRisksDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	err := h.usecase.ClearOwnRisksDraft(c)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to save draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Own risks draft data cleared successfully", http.StatusOK, "", nil)
}

// ConfirmOwnRisks handles POST /api/insurance-own-risks/draft/confirm
func (h *handler) ConfirmOwnRisks(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	result, err := h.usecase.ConfirmOwnRisks(c)
	if err != nil {
		if err.Error() == "No draft data to confirm" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to confirm own risks: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"message":        fmt.Sprintf("Successfully confirmed %d own risks", result.InsertedCount),
		"inserted_count": result.InsertedCount,
	}, "Own risks confirmed successfully", http.StatusOK, "", nil)
}

// GetOwnRisks handles GET /api/insurance-own-risks
func (h *handler) GetOwnRisks(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	insuranceCode := c.QueryParam("insurance_code")

	rulesList, err := h.usecase.GetOwnRisks(c, insuranceCode)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to query own risks", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", rulesList, "Success to Get Own Risks", http.StatusOK, "", nil)
}

// UpdateOwnRisks handles PUT /api/insurance-own-risks/:id
func (h *handler) UpdateOwnRisks(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid ID", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	var updateData struct {
		ProductType string `json:"product_type"`
		Code        string `json:"code"`
		Title       string `json:"title"`
		Value       string `json:"value"`
		ValueType   string `json:"value_type"`
		Description string `json:"description"`
		IsMandatory int    `json:"is_mandatory"`
		IsActive    int    `json:"is_active"`
	}

	if err := ac.Bind(&updateData); err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid request: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	// Get user name from header
	userName := utils.GetUserNameFromRequest(c)

	req := insuranceownrisks.UpdateOwnRiskRequest{
		UpdateData: map[string]interface{}{
			"product_type": updateData.ProductType,
			"code":         updateData.Code,
			"title":        updateData.Title,
			"value":        updateData.Value,
			"value_type":   updateData.ValueType,
			"description":  updateData.Description,
			"is_mandatory": updateData.IsMandatory,
			"is_active":    updateData.IsActive,
		},
	}

	err = h.usecase.UpdateOwnRisk(c, id, req, userName)
	if err != nil {
		if err.Error() == "Own risks not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to update own risks", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Own risks updated successfully", http.StatusOK, "", nil)
}


