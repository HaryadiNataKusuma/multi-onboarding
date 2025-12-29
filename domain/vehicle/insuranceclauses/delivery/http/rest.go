package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/insuranceclauses"
	"multi-onboarding/utils"
)

type handler struct {
	usecase insuranceclauses.Usecase
}

// NewInsuranceClausesHandler creates a new insurance clauses handler
func NewInsuranceClausesHandler(usecase insuranceclauses.Usecase) *handler {
	return &handler{
		usecase: usecase,
	}
}

// AddInsuranceClausesHandler registers insurance clauses routes to Echo router
func AddInsuranceClausesHandler(e *echo.Echo, usecase insuranceclauses.Usecase) {
	handler := &handler{
		usecase: usecase,
	}

	// Onboarding routes
	e.GET("/api/insurance-clauses/download-template", handler.DownloadClausesTemplate)
	e.POST("/api/insurance-clauses/upload", handler.UploadClausesFile)
	e.GET("/api/insurance-clauses/draft", handler.GetClausesDraft)
	e.PUT("/api/insurance-clauses/draft/:id", handler.UpdateClausesDraft)
	e.DELETE("/api/insurance-clauses/draft/:id", handler.DeleteClausesDraft)
	e.POST("/api/insurance-clauses/draft/clear", handler.ClearClausesDraft)
	e.POST("/api/insurance-clauses/draft/confirm", handler.ConfirmClauses)
	e.GET("/api/insurance-clauses", handler.GetClauses)
	e.PUT("/api/insurance-clauses/:id", handler.UpdateClauses)
}

// DownloadClausesTemplate handles GET /api/insurance-clauses/download-template
func (h *handler) DownloadClausesTemplate(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	insuranceCode := c.QueryParam("insurance_code")
	if insuranceCode == "" {
		return ac.CustomResponse("Failed", nil, "insurance_code is required", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	fileBytes, fileName, err := h.usecase.DownloadClausesTemplate(c, insuranceCode)
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

// UploadClausesFile handles POST /api/insurance-clauses/upload
func (h *handler) UploadClausesFile(c echo.Context) error {
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

	rowCount, err := h.usecase.UploadClausesFile(c, src, file.Filename, insuranceCode)
	if err != nil {
		if err.Error() == "No sheet found in Excel file" || err.Error() == "Failed to read file" || strings.Contains(err.Error(), "Invalid headers. Expected") {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to process file: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"message":    fmt.Sprintf("Successfully uploaded %d clauses", rowCount),
		"rows_count": rowCount,
	}, "File uploaded successfully", http.StatusOK, "", nil)
}

// GetClausesDraft handles GET /api/insurance-clauses/draft
func (h *handler) GetClausesDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	insuranceCode := c.QueryParam("insurance_code")

	rulesList, err := h.usecase.GetClausesDraft(c, insuranceCode)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to load draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", rulesList, "Success to Get Clauses Draft", http.StatusOK, "", nil)
}

// UpdateClausesDraft handles PUT /api/insurance-clauses/draft/:id
func (h *handler) UpdateClausesDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid ID", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	var updateData struct {
		InsuranceCode string  `json:"insurance_code"`
		ProductType   string  `json:"product_type"`
		Code          string  `json:"code"`
		Title         string  `json:"title"`
		Content       string  `json:"content"`
		Description   *string `json:"description"`
		IsMandatory   *bool   `json:"is_mandatory"`
		IsActive      *bool   `json:"is_active"`
	}

	if err := ac.Bind(&updateData); err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid request: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	req := insuranceclauses.UpdateClauseDraftRequest{
		ID:            id,
		InsuranceCode: updateData.InsuranceCode,
		ProductType:   updateData.ProductType,
		Code:          updateData.Code,
		Title:         updateData.Title,
		Content:       updateData.Content,
		Description:   updateData.Description,
		IsMandatory:   updateData.IsMandatory,
		IsActive:      updateData.IsActive,
	}

	err = h.usecase.UpdateClauseDraft(c, req)
	if err != nil {
		if err.Error() == "Clauses draft not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to save draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Clauses updated successfully", http.StatusOK, "", nil)
}

// DeleteClausesDraft handles DELETE /api/insurance-clauses/draft/:id
func (h *handler) DeleteClausesDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid ID", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	err = h.usecase.DeleteClauseDraft(c, id)
	if err != nil {
		if err.Error() == "Clauses draft not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to save draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Clauses deleted successfully", http.StatusOK, "", nil)
}

// ClearClausesDraft handles POST /api/insurance-clauses/draft/clear
func (h *handler) ClearClausesDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	err := h.usecase.ClearClausesDraft(c)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to save draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Clauses draft data cleared successfully", http.StatusOK, "", nil)
}

// ConfirmClauses handles POST /api/insurance-clauses/draft/confirm
func (h *handler) ConfirmClauses(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	result, err := h.usecase.ConfirmClauses(c)
	if err != nil {
		if err.Error() == "No draft data to confirm" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to confirm clauses: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"message":        fmt.Sprintf("Successfully confirmed %d clauses", result.InsertedCount),
		"inserted_count": result.InsertedCount,
	}, "Clauses confirmed successfully", http.StatusOK, "", nil)
}

// GetClauses handles GET /api/insurance-clauses
func (h *handler) GetClauses(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	insuranceCode := c.QueryParam("insurance_code")

	rulesList, err := h.usecase.GetClauses(c, insuranceCode)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to query clauses", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", rulesList, "Success to Get Clauses", http.StatusOK, "", nil)
}

// UpdateClauses handles PUT /api/insurance-clauses/:id
func (h *handler) UpdateClauses(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid ID", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	var updateData struct {
		ProductType string `json:"product_type"`
		Code        string `json:"code"`
		Title       string `json:"title"`
		Content     string `json:"content"`
		Description string `json:"description"`
		IsMandatory int    `json:"is_mandatory"`
		IsActive    int    `json:"is_active"`
	}

	if err := ac.Bind(&updateData); err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid request: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	// Get user name from header
	userName := utils.GetUserNameFromRequest(c)

	req := insuranceclauses.UpdateClauseRequest{
		UpdateData: map[string]interface{}{
			"product_type": updateData.ProductType,
			"code":         updateData.Code,
			"title":        updateData.Title,
			"content":      updateData.Content,
			"description":  updateData.Description,
			"is_mandatory": updateData.IsMandatory,
			"is_active":    updateData.IsActive,
		},
	}

	err = h.usecase.UpdateClause(c, id, req, userName)
	if err != nil {
		if err.Error() == "Clauses not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to update clauses", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Clauses updated successfully", http.StatusOK, "", nil)
}

