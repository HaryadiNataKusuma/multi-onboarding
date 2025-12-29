package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"multi-onboarding/domain/Travel/addon"
)

type handlerAddon struct {
	usecase addon.Usecase
}

// NewAddonHandler creates a new addon handler
func NewAddonHandler(usecase addon.Usecase) *handlerAddon {
	return &handlerAddon{
		usecase: usecase,
	}
}

// GetAddons handles GET /api/addons
func (h *handlerAddon) GetAddons(c echo.Context) error {
	addons, err := h.usecase.GetAllAddons()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if addons == nil {
		addons = []addon.Addon{}
	}

	return c.JSON(http.StatusOK, addons)
}

// GetAddon handles GET /api/addons/:id
func (h *handlerAddon) GetAddon(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid addon ID",
		})
	}

	a, err := h.usecase.GetAddonByID(id)
	if err != nil {
		if err.Error() == "addon not found" || err.Error() == "invalid addon ID" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Addon not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, a)
}

// CreateAddon handles POST /api/addons
func (h *handlerAddon) CreateAddon(c echo.Context) error {
	var req addon.CreateAddonRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	if err := h.usecase.CreateAddon(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Addon created successfully",
	})
}

// UpdateAddon handles PUT /api/addons/:id
func (h *handlerAddon) UpdateAddon(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid addon ID",
		})
	}

	var req addon.UpdateAddonRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	if err := h.usecase.UpdateAddon(id, req); err != nil {
		if err.Error() == "addon not found" || err.Error() == "invalid addon ID" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Addon not found",
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	updatedAddon, _ := h.usecase.GetAddonByID(id)
	return c.JSON(http.StatusOK, updatedAddon)
}

// DeleteAddon handles DELETE /api/addons/:id
func (h *handlerAddon) DeleteAddon(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid addon ID",
		})
	}

	err = h.usecase.DeleteAddon(id)
	if err != nil {
		if err.Error() == "addon not found" || err.Error() == "invalid addon ID" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Addon not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}
