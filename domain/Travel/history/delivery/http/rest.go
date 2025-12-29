package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"multi-onboarding/domain/Travel/history"
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

// GetAll handles GET /api/histories
func (h *handlerHistory) GetAll(c echo.Context) error {
	histories, err := h.usecase.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if histories == nil {
		histories = []history.History{}
	}

	return c.JSON(http.StatusOK, histories)
}

// GetByID handles GET /api/histories/:id
func (h *handlerHistory) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid history ID",
		})
	}

	hist, err := h.usecase.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, hist)
}

// GetByTableName handles GET /api/histories/table/:table_name
func (h *handlerHistory) GetByTableName(c echo.Context) error {
	tableName := c.Param("table_name")

	histories, err := h.usecase.GetByTableName(tableName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if histories == nil {
		histories = []history.History{}
	}

	return c.JSON(http.StatusOK, histories)
}

// GetByRecordID handles GET /api/histories/table/:table_name/record/:record_id
func (h *handlerHistory) GetByRecordID(c echo.Context) error {
	tableName := c.Param("table_name")
	recordID, err := strconv.Atoi(c.Param("record_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid record ID",
		})
	}

	histories, err := h.usecase.GetByRecordID(tableName, recordID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if histories == nil {
		histories = []history.History{}
	}

	return c.JSON(http.StatusOK, histories)
}
