package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/insurance"
	"multi-onboarding/utils"
)

type handlerInsurance struct {
	usecase insurance.Usecase
}

// NewInsuranceHandler creates a new insurance handler
func NewInsuranceHandler(usecase insurance.Usecase) *handlerInsurance {
	return &handlerInsurance{
		usecase: usecase,
	}
}

// SaveInsuranceDraft handles POST /api/insurances/draft/add
func (h *handlerInsurance) SaveInsuranceDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	var input struct {
		Code   string `json:"code" validate:"required"`
		Name   string `json:"name" validate:"required"`
		Status int    `json:"status" validate:"required"`
		Logo   string `json:"logo"`
	}

	if err := ac.Bind(&input); err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid request: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	if err := ac.Validate(input); err != nil {
		return ac.CustomResponse("Failed", nil, "Validation error: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
	}

	req := insurance.SaveInsuranceDraftRequest{
		Code:   input.Code,
		Name:   input.Name,
		Status: input.Status,
		Logo:   input.Logo,
	}

	newInsurance, err := h.usecase.SaveInsuranceDraft(c, req)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to save draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"id":       newInsurance.ID,
		"row_data": newInsurance,
	}, "Data saved to draft successfully", http.StatusOK, "", nil)
}

// GetInsurancesDraft handles GET /api/insurances/draft
func (h *handlerInsurance) GetInsurancesDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	insuranceCode := c.QueryParam("insurance_code")

	insurances, err := h.usecase.GetInsurancesDraft(c, insuranceCode)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to load draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	// Ensure data is always an array, not null
	if insurances == nil {
		insurances = []utils.InsuranceDraft{}
	}

	return ac.CustomResponse("Success", insurances, "Success to Get Insurances Draft", http.StatusOK, "", nil)
}

// UpdateInsuranceDraft handles PUT /api/insurances/draft/:id
func (h *handlerInsurance) UpdateInsuranceDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid ID", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	var input struct {
		Code   string `json:"code" validate:"required"`
		Name   string `json:"name" validate:"required"`
		Status int    `json:"status" validate:"required"`
		Logo   string `json:"logo"`
	}

	if err := ac.Bind(&input); err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid request: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	if err := ac.Validate(input); err != nil {
		return ac.CustomResponse("Failed", nil, "Validation error: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
	}

	req := insurance.UpdateInsuranceDraftRequest{
		ID:     id,
		Code:   input.Code,
		Name:   input.Name,
		Status: input.Status,
		Logo:   input.Logo,
	}

	err = h.usecase.UpdateInsuranceDraft(c, req)
	if err != nil {
		if err.Error() == "Draft data not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to save draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"id": id,
	}, "Draft data updated successfully", http.StatusOK, "", nil)
}

// DeleteInsuranceDraft handles DELETE /api/insurances/draft/:id
func (h *handlerInsurance) DeleteInsuranceDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid ID", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	err = h.usecase.DeleteInsuranceDraft(c, id)
	if err != nil {
		if err.Error() == "Draft data not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to save draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Draft data deleted successfully", http.StatusOK, "", nil)
}

// ClearInsurancesDraft handles POST /api/insurances/draft/clear
func (h *handlerInsurance) ClearInsurancesDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	err := h.usecase.ClearInsurancesDraft(c)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to save draft data", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Insurances draft data cleared successfully", http.StatusOK, "", nil)
}

// ConfirmInsurances handles POST /api/insurances/draft/confirm
func (h *handlerInsurance) ConfirmInsurances(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	result, err := h.usecase.ConfirmInsurances(c)
	if err != nil {
		if err.Error() == "No data available in Draft to confirm" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to confirm insurances", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"success_count": result.SuccessCount,
	}, "Data confirmed successfully", http.StatusOK, "", nil)
}

// GetInsurances handles GET /api/insurances
func (h *handlerInsurance) GetInsurances(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	insurances, err := h.usecase.GetInsurances(c)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to query database", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	// Ensure data is always an array, not null
	if insurances == nil {
		insurances = []map[string]interface{}{}
	}

	return ac.CustomResponse("Success", insurances, "Success to Get Insurances", http.StatusOK, "", nil)
}

// UpdateInsurance handles PUT /api/insurances/:id
func (h *handlerInsurance) UpdateInsurance(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid ID", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	var input struct {
		Code   string `json:"code" validate:"required"`
		Name   string `json:"name" validate:"required"`
		Status int    `json:"status" validate:"required"`
		Logo   string `json:"logo"`
	}

	if err := ac.Bind(&input); err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid request: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	if err := ac.Validate(input); err != nil {
		return ac.CustomResponse("Failed", nil, "Validation error: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
	}

	// Get user name from header
	userName := utils.GetUserNameFromRequest(c)

	req := insurance.UpdateInsuranceRequest{
		Code:   input.Code,
		Name:   input.Name,
		Status: input.Status,
		Logo:   input.Logo,
	}

	err = h.usecase.UpdateInsurance(c, id, req, userName)
	if err != nil {
		if err.Error() == "Insurance not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to update insurance", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"id": id,
	}, "Insurance updated successfully", http.StatusOK, "", nil)
}


