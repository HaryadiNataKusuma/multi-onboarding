package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/addon"
	"multi-onboarding/utils"
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
	ac := utils.GetCustomApplicationContextStub(c)

	addons, err := h.usecase.GetAddons(c)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to query database", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", addons, "Success to Get Addons", http.StatusOK, "", nil)
}

// UpdateAddon handles PUT /api/addons/:id
func (h *handlerAddon) UpdateAddon(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid ID", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	var input struct {
		Code     string `json:"code" validate:"required"`
		Name     string `json:"name" validate:"required"`
		IsActive bool   `json:"is_active"`
	}

	if err := ac.Bind(&input); err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid request: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	if err := ac.Validate(input); err != nil {
		return ac.CustomResponse("Failed", nil, "Validation error: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
	}

	// Get user name from header
	userName := utils.GetUserNameFromRequest(c)

	req := addon.UpdateAddonRequest{
		Code:     input.Code,
		Name:     input.Name,
		IsActive: input.IsActive,
	}

	err = h.usecase.UpdateAddon(c, id, req, userName)
	if err != nil {
		// Check error message for not found
		if err.Error() == "Addon not found" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusNotFound, "COMM-SYSTEM-01", nil)
		}
		return ac.CustomResponse("Failed", nil, "Failed to update addon", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", nil, "Addon updated successfully", http.StatusOK, "", nil)
}

// ConfirmAddons handles POST /api/addons/confirm
func (h *handlerAddon) ConfirmAddons(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	var request struct {
		Addons []struct {
			Code     string `json:"code" validate:"required"`
			Name     string `json:"name" validate:"required"`
			IsActive bool   `json:"is_active"`
		} `json:"addons" validate:"required"`
	}

	if err := ac.Bind(&request); err != nil {
		return ac.CustomResponse("Failed", nil, "Invalid request: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	if err := ac.Validate(request); err != nil {
		return ac.CustomResponse("Failed", nil, "Validation error: "+err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
	}

	// Convert request.Addons to addon.ConfirmAddonsRequest format
	addons := make([]struct {
		Code     string
		Name     string
		IsActive bool
	}, len(request.Addons))

	for i, a := range request.Addons {
		addons[i] = struct {
			Code     string
			Name     string
			IsActive bool
		}{
			Code:     a.Code,
			Name:     a.Name,
			IsActive: a.IsActive,
		}
	}

	req := addon.ConfirmAddonsRequest{
		Addons: addons,
	}

	insertedCount, err := h.usecase.ConfirmAddons(c, req)
	if err != nil {
		// Check error message for validation error
		if err.Error() == "No addons to confirm" {
			return ac.CustomResponse("Failed", nil, err.Error(), http.StatusBadRequest, "COMM-REQUEST-02", nil)
		}
		return ac.CustomResponse("Failed", map[string]interface{}{
			"failed_at_index": insertedCount,
			"details":         err.Error(),
		}, "Failed to insert addons", http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	return ac.CustomResponse("Success", map[string]interface{}{
		"inserted_count": insertedCount,
	}, "Addons confirmed successfully", http.StatusOK, "", nil)
}


