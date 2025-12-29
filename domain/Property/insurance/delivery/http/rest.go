package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"multi-onboarding/domain/Property/insurance"
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

// AddDraft handles POST /api/property/insurances/draft/add
func (h *handlerInsurance) AddDraft(c echo.Context) error {
	var draft insurance.InsuranceDraft
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

// GetAllDrafts handles GET /api/property/insurances/draft
func (h *handlerInsurance) GetAllDrafts(c echo.Context) error {
	drafts, err := h.usecase.GetAllDrafts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if drafts == nil {
		drafts = []insurance.InsuranceDraft{}
	}

	return c.JSON(http.StatusOK, drafts)
}

// GetDraftByID handles GET /api/property/insurances/draft/:id
func (h *handlerInsurance) GetDraftByID(c echo.Context) error {
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

// UpdateDraft handles PUT /api/property/insurances/draft/:id
func (h *handlerInsurance) UpdateDraft(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID",
		})
	}

	var draft insurance.InsuranceDraft
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

// DeleteDraft handles DELETE /api/property/insurances/draft/:id
func (h *handlerInsurance) DeleteDraft(c echo.Context) error {
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

// ClearDrafts handles POST /api/property/insurances/draft/clear
func (h *handlerInsurance) ClearDrafts(c echo.Context) error {
	if err := h.usecase.ClearDrafts(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// ConfirmDrafts handles POST /api/property/insurances/draft/confirm
func (h *handlerInsurance) ConfirmDrafts(c echo.Context) error {
	if err := h.usecase.ConfirmDrafts(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Drafts confirmed successfully",
	})
}

// GetAll handles GET /api/property/insurances
func (h *handlerInsurance) GetAll(c echo.Context) error {
	c.Logger().Info("GetAll Property insurances: Request received")
	
	insurances, err := h.usecase.GetAll()
	if err != nil {
		c.Logger().Errorf("GetAll Property insurances error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Ensure we always return an array, even if empty
	if insurances == nil {
		insurances = []insurance.Insurance{}
	}

	c.Logger().Infof("GetAll Property insurances: Returning %d insurances", len(insurances))
	if len(insurances) > 0 {
		c.Logger().Infof("GetAll Property insurances: First insurance - Code: %s, Name: %s", insurances[0].Code, insurances[0].Name)
	}
	
	return c.JSON(http.StatusOK, insurances)
}

// GetByID handles GET /api/property/insurances/:id
func (h *handlerInsurance) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID",
		})
	}

	ins, err := h.usecase.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, ins)
}

// Create handles POST /api/property/insurances
func (h *handlerInsurance) Create(c echo.Context) error {
	var req insurance.CreateInsuranceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := h.usecase.Create(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Insurance created successfully",
	})
}

// Update handles PUT /api/property/insurances/:id
func (h *handlerInsurance) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID",
		})
	}

	var req insurance.UpdateInsuranceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := h.usecase.Update(id, req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	updatedInsurance, _ := h.usecase.GetByID(id)
	return c.JSON(http.StatusOK, updatedInsurance)
}

// Delete handles DELETE /api/property/insurances/:id
func (h *handlerInsurance) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID",
		})
	}

	if err := h.usecase.Delete(id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}



