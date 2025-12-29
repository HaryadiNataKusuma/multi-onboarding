package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"multi-onboarding/domain/Travel/region"
)

type handlerRegion struct {
	usecase region.Usecase
}

// NewRegionHandler creates a new region handler
func NewRegionHandler(usecase region.Usecase) *handlerRegion {
	return &handlerRegion{
		usecase: usecase,
	}
}

// GetRegions handles GET /api/regions
func (h *handlerRegion) GetRegions(c echo.Context) error {
	regions, err := h.usecase.GetAll()
	if err != nil {
		c.Logger().Errorf("GetRegions error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Ensure we always return an array, even if empty
	if regions == nil {
		regions = []region.Region{}
	}

	c.Logger().Infof("GetRegions: Returning %d regions", len(regions))
	return c.JSON(http.StatusOK, regions)
}

// GetRegion handles GET /api/regions/:id
func (h *handlerRegion) GetRegion(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid region ID",
		})
	}

	reg, err := h.usecase.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, reg)
}

// CreateRegion handles POST /api/regions
func (h *handlerRegion) CreateRegion(c echo.Context) error {
	var req region.CreateRegionRequest

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
		"message": "Region created successfully",
	})
}

// UpdateRegion handles PUT /api/regions/:id
func (h *handlerRegion) UpdateRegion(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid region ID",
		})
	}

	var req region.UpdateRegionRequest
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

	updatedRegion, _ := h.usecase.GetByID(id)
	return c.JSON(http.StatusOK, updatedRegion)
}

// DeleteRegion handles DELETE /api/regions/:id
func (h *handlerRegion) DeleteRegion(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid region ID",
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
