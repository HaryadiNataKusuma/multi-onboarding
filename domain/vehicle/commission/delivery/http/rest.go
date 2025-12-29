package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/commission"
	"multi-onboarding/utils"
)

type handler struct {
	usecase commission.Usecase
}

// NewCommissionHandler creates a new commission handler
func NewCommissionHandler(usecase commission.Usecase) *handler {
	return &handler{
		usecase: usecase,
	}
}

// AddCommissionHandler registers commission routes to Echo router
func AddCommissionHandler(e *echo.Echo, usecase commission.Usecase) {
	handler := &handler{
		usecase: usecase,
	}

	// Onboarding routes
	e.POST("/api/commission/save", handler.SaveCommissionDraft)
	e.GET("/api/commission/draft/all", handler.GetCommissionDraft)
	e.POST("/api/commission/confirm", handler.ConfirmCommissions)
	e.POST("/api/commission/draft/clear", handler.ClearCommissionDraft)
	e.GET("/api/commission/confirmed/all", handler.GetConfirmedCommissions)
	e.PUT("/api/plan-commissions/bulk", handler.BulkUpdatePlanCommissions)
}

// SaveCommissionDraft handles POST /api/commission/save
func (h *handler) SaveCommissionDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	log.Printf("=== SaveCommissionDraft called ===")

	var req struct {
		BasicCommission float64 `json:"basic_commission" validate:"required"`
		PlanCommissions struct {
			CommissionVatType string  `json:"commission_vat_type" validate:"required"`
			AfPercentage      float64 `json:"af_percentage" validate:"required"`
			AfVatType         string  `json:"af_vat_type" validate:"required"`
			AdminFee          float64 `json:"admin_fee" validate:"required"`
		} `json:"plan_commissions" validate:"required"`
	}

	if err := ac.Bind(&req); err != nil {
		log.Printf("ERROR: Failed to bind request: %v", err)
		return ac.CustomResponse("Failed", nil, "Invalid request: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	log.Printf("Request data: BasicCommission=%f, AdminFee=%f, CommissionVatType=%s, AfVatType=%s",
		req.BasicCommission, req.PlanCommissions.AdminFee, req.PlanCommissions.CommissionVatType, req.PlanCommissions.AfVatType)

	if err := ac.Validate(req); err != nil {
		log.Printf("ERROR: Validation failed: %v", err)
		return ac.CustomResponse("Failed", nil, "Validation error: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
	}

	// Validate
	if req.PlanCommissions.CommissionVatType == "" {
		log.Printf("ERROR: Commission VAT Type is empty")
		return ac.CustomResponse("Failed", nil, "Commission VAT Type is required", http.StatusBadRequest, "COMM-REQUEST-02", nil)
	}
	if req.PlanCommissions.AfVatType == "" {
		log.Printf("ERROR: AF VAT Type is empty")
		return ac.CustomResponse("Failed", nil, "AF VAT Type is required", http.StatusBadRequest, "COMM-REQUEST-02", nil)
	}
	if req.PlanCommissions.AdminFee < 0 {
		log.Printf("ERROR: Admin Fee is negative: %f", req.PlanCommissions.AdminFee)
		return ac.CustomResponse("Failed", nil, "Admin Fee must be greater than or equal to 0", http.StatusBadRequest, "COMM-REQUEST-02", nil)
	}

	insuranceCode := c.QueryParam("insurance_code")
	log.Printf("Insurance code from query: %s", insuranceCode)
	if insuranceCode == "" {
		log.Printf("ERROR: insurance_code query parameter is missing")
		return ac.CustomResponse("Failed", nil, "insurance_code query parameter is required", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	request := commission.SaveCommissionDraftRequest{
		BasicCommission: req.BasicCommission,
		PlanCommissions: commission.PlanCommissionData{
			CommissionVatType: req.PlanCommissions.CommissionVatType,
			AfPercentage:      req.PlanCommissions.AfPercentage,
			AfVatType:         req.PlanCommissions.AfVatType,
			AdminFee:          req.PlanCommissions.AdminFee,
		},
		InsuranceCode: insuranceCode,
	}

	log.Printf("Calling usecase.SaveCommissionDraft with insurance_code: %s", insuranceCode)
	err := h.usecase.SaveCommissionDraft(c, request)
	if err != nil {
		errMsg := err.Error()
		log.Printf("ERROR: SaveCommissionDraft failed: %s", errMsg)

		if errMsg == "No active products found for this insurance" {
			return ac.CustomResponse("Failed", nil, errMsg, http.StatusBadRequest, "COMM-REQUEST-02", nil)
		}

		return ac.CustomResponse("Failed", nil, "Failed to save commission draft: "+errMsg, http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	log.Printf("SUCCESS: Commission draft saved successfully")
	return ac.CustomResponse("Success", nil, "Commission data saved to draft successfully", http.StatusOK, "", nil)
}

// GetCommissionDraft handles GET /api/commission/draft/all
func (h *handler) GetCommissionDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	insuranceCode := c.QueryParam("insurance_code")
	if insuranceCode == "" {
		return ac.CustomResponse("Failed", nil, "insurance_code query parameter is required", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	result, err := h.usecase.GetCommissionDraft(c, insuranceCode)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to fetch commission draft: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"commissions":      result.Commissions,
		"plan_commissions": result.PlanCommissions,
		"config_products":  result.ConfigProducts,
	}, "Success to Get Commission Draft", http.StatusOK, "", nil)
}

// ConfirmCommissions handles POST /api/commission/confirm
func (h *handler) ConfirmCommissions(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	result, err := h.usecase.ConfirmCommissions(c)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to confirm commissions: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"message":           "Commission data berhasil dikonfirmasi! " + strconv.Itoa(result.TotalMoved) + " records moved to confirmed tables.",
		"records_confirmed": result.TotalMoved,
	}, "Commission data confirmed successfully", http.StatusOK, "", nil)
}

// ClearCommissionDraft handles POST /api/commission/draft/clear
func (h *handler) ClearCommissionDraft(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	err := h.usecase.ClearCommissionDraft(c)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to clear draft tables: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Commission draft data berhasil dihapus!", http.StatusOK, "", nil)
}

// GetConfirmedCommissions handles GET /api/commission/confirmed/all
func (h *handler) GetConfirmedCommissions(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	insuranceCode := c.QueryParam("insurance_code")
	productCode := c.QueryParam("product_code")

	result, err := h.usecase.GetConfirmedCommissions(c, insuranceCode, productCode)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to fetch confirmed commissions: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"commissions":      result.Commissions,
		"plan_commissions": result.PlanCommissions,
		"config_products":  result.ConfigProducts,
	}, "Success to Get Confirmed Commissions", http.StatusOK, "", nil)
}

// BulkUpdatePlanCommissions handles PUT /api/plan-commissions/bulk
func (h *handler) BulkUpdatePlanCommissions(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	var req commission.BulkUpdatePlanCommissionsRequest
	if err := ac.Bind(&req); err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid request: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	if len(req.IDs) == 0 {
		return ac.CustomResponse("Failed", nil, "At least one ID is required", http.StatusBadRequest, "COMM-REQUEST-02", nil)
	}

	if req.Updates.CommissionPercentage == nil && req.Updates.CommissionVatType == nil &&
		req.Updates.AfPercentage == nil && req.Updates.AfVatType == nil && req.Updates.AdminFee == nil {
		return ac.CustomResponse("Failed", nil, "At least one field to update is required", http.StatusBadRequest, "COMM-REQUEST-02", nil)
	}

	// Get user name from header
	userName := utils.GetUserNameFromRequest(c)

	totalInserted, err := h.usecase.BulkUpdatePlanCommissions(c, req, userName)
	if err != nil {
		if err.Error() == "At least one ID is required" || err.Error() == "No plan commissions found for the provided IDs" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to update plan commissions: "+err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"message":       "Successfully updated " + strconv.Itoa(totalInserted) + " record(s). Existing records deactivated, new records inserted with version + 1.",
		"rows_affected": totalInserted,
	}, "Plan commissions updated successfully", http.StatusOK, "", nil)
}


