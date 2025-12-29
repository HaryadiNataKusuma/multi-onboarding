package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/downloadallresult"
	"multi-onboarding/utils"
)

type handler struct {
	usecase downloadallresult.Usecase
}

// NewDownloadAllResultHandler creates a new download all result handler
func NewDownloadAllResultHandler(usecase downloadallresult.Usecase) *handler {
	return &handler{
		usecase: usecase,
	}
}

// DownloadAllResult handles GET /api/download-all-result
// Generates an Excel file with multiple sheets containing all RESULT data from all sections
func (h *handler) DownloadAllResult(c echo.Context) error {
	ac := utils.GetCustomApplicationContextStub(c)

	insuranceCode := c.QueryParam("insurance_code")
	if insuranceCode == "" {
		return ac.CustomResponse("Failed", nil, "insurance_code is required", http.StatusBadRequest, "COMM-REQUEST-01", nil)
	}

	fileBytes, fileName, err := h.usecase.DownloadAllResult(c, insuranceCode)
	if err != nil {
		return ac.CustomResponse("Failed", nil, err.Error(), http.StatusInternalServerError, "COMM-SYSTEM-02", nil)
	}

	// Set response headers
	c.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))

	// Write Excel file to response
	return c.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}


