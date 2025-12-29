package http

import (
	"net/http"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/history"
	"multi-onboarding/utils"
)

type handlerHistory struct {
	usecase history.Usecase
}

// NewHistoryHandler creates a new history handler
func NewHistoryHandler(usecase history.Usecase) *handlerHistory {
	return &handlerHistory{
		usecase: usecase,
	}
}

// GetHistories handles GET /api/histories
func (h *handlerHistory) GetHistories(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)
	
	insuranceCode := c.QueryParam("insurance_code")
	sectionFilter := c.QueryParam("section")

	histories, err := h.usecase.GetHistories(c, insuranceCode, sectionFilter)
	if err != nil {
		return ac.CustomResponse("Failed", nil, "Failed to query database: "+err.Error(), http.StatusInternalServerError, "HIST-SYSTEM-01", nil)
	}

	return ac.CustomResponse("Success", histories, "Success to Get Histories", http.StatusOK, "", nil)
}


