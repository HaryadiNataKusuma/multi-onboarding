package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"multi-onboarding/domain/Travel/commission"
)

type handlerCommission struct {
	usecase commission.Usecase
}

// NewCommissionHandler creates a new commission handler
func NewCommissionHandler(usecase commission.Usecase) *handlerCommission {
	return &handlerCommission{
		usecase: usecase,
	}
}

// GetAll handles GET /api/commissions
func (h *handlerCommission) GetAll(c echo.Context) error {
	commissions, err := h.usecase.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if commissions == nil {
		commissions = []commission.Commission{}
	}

	return c.JSON(http.StatusOK, commissions)
}

// GetByID handles GET /api/commissions/:id
func (h *handlerCommission) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid commission ID",
		})
	}

	comm, err := h.usecase.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, comm)
}

// GetByProductCode handles GET /api/commissions/product/:productCode
func (h *handlerCommission) GetByProductCode(c echo.Context) error {
	productCode := c.Param("product_code")

	comm, err := h.usecase.GetByProductCode(productCode)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, comm)
}

// Create handles POST /api/commissions
func (h *handlerCommission) Create(c echo.Context) error {
	var req commission.CreateCommissionRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	if err := h.usecase.Create(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Commission created successfully",
	})
}

// Update handles PUT /api/commissions/:id
func (h *handlerCommission) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid commission ID",
		})
	}

	var req commission.UpdateCommissionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	if err := h.usecase.Update(id, req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	updatedCommission, _ := h.usecase.GetByID(id)
	return c.JSON(http.StatusOK, updatedCommission)
}

// Delete handles DELETE /api/commissions/:id
func (h *handlerCommission) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid commission ID",
		})
	}

	err = h.usecase.Delete(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// Draft operations

// AddDraft handles POST /api/commissions/draft
func (h *handlerCommission) AddDraft(c echo.Context) error {
	var draft commission.CommissionDraft

	if err := c.Bind(&draft); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := h.usecase.AddDraft(&draft); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, draft)
}

// GetAllDrafts handles GET /api/commissions/draft
func (h *handlerCommission) GetAllDrafts(c echo.Context) error {
	drafts, err := h.usecase.GetAllDrafts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if drafts == nil {
		drafts = []commission.CommissionDraft{}
	}

	return c.JSON(http.StatusOK, drafts)
}

// GetDraftByID handles GET /api/commissions/draft/:id
func (h *handlerCommission) GetDraftByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID",
		})
	}

	draft, err := h.usecase.GetDraftByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, draft)
}

// UpdateDraft handles PUT /api/commissions/draft/:id
func (h *handlerCommission) UpdateDraft(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID",
		})
	}

	var draft commission.CommissionDraft
	if err := c.Bind(&draft); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := h.usecase.UpdateDraft(id, &draft); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, draft)
}

// DeleteDraft handles DELETE /api/commissions/draft/:id
func (h *handlerCommission) DeleteDraft(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID",
		})
	}

	if err := h.usecase.DeleteDraft(id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// ClearDrafts handles POST /api/commissions/draft/clear
func (h *handlerCommission) ClearDrafts(c echo.Context) error {
	if err := h.usecase.ClearDrafts(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// ConfirmDrafts handles POST /api/commissions/draft/confirm
func (h *handlerCommission) ConfirmDrafts(c echo.Context) error {
	var requestData struct {
		CreatedBy int64 `json:"created_by"`
	}

	if err := c.Bind(&requestData); err != nil {
		requestData.CreatedBy = 1
	}
	if requestData.CreatedBy == 0 {
		requestData.CreatedBy = 1
	}

	if err := h.usecase.ConfirmDrafts(requestData.CreatedBy); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Commissions confirmed successfully",
	})
}
