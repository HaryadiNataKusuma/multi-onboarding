package http

import (
	"database/sql"
	"net/http"

	"multi-onboarding/domain/Property/template"
	"multi-onboarding/utils"

	"github.com/labstack/echo"
)

type handlerTemplate struct {
	usecase template.Usecase
}

// NewTemplateHandler creates a new template handler
func NewTemplateHandler(usecase template.Usecase) *handlerTemplate {
	return &handlerTemplate{
		usecase: usecase,
	}
}

// GetAll handles GET /api/property/templates
func (h *handlerTemplate) GetAll(c echo.Context) error {
	insuranceCode := c.QueryParam("insurance_code")

	var query string
	var args []interface{}

	if insuranceCode != "" {
		query = `
			SELECT DISTINCT 
				t.locale, t.id, t.value, t.created_by, t.created_at, t.updated_by, t.updated_at,
				p.insurance_code, p.name as product_name
			FROM property_service_development.templates t
			INNER JOIN property_service_development.products p ON (t.id = p.summary OR t.id = p.insurance_detail)
			WHERE p.insurance_code = ?
			ORDER BY t.created_at DESC
		`
		args = []interface{}{insuranceCode}
	} else {
		query = `
			SELECT DISTINCT 
				t.locale, t.id, t.value, t.created_by, t.created_at, t.updated_by, t.updated_at,
				p.insurance_code, p.name as product_name
			FROM property_service_development.templates t
			INNER JOIN property_service_development.products p ON (t.id = p.summary OR t.id = p.insurance_detail)
			ORDER BY t.created_at DESC
		`
		args = []interface{}{}
	}

	db := utils.GetDB()
	rows, err := db.Query(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to query templates: " + err.Error(),
		})
	}
	defer rows.Close()

	enrichedTemplates := make([]map[string]interface{}, 0)
	for rows.Next() {
		var locale, id, value string
		var createdBy int64
		var createdAt sql.NullTime
		var updatedBy sql.NullInt64
		var updatedAt sql.NullTime
		var insCode sql.NullString
		var prodName sql.NullString

		if err := rows.Scan(&locale, &id, &value, &createdBy, &createdAt, &updatedBy, &updatedAt, &insCode, &prodName); err != nil {
			continue
		}

		enriched := map[string]interface{}{
			"locale":     locale,
			"id":         id,
			"value":      value,
			"created_by": createdBy,
		}

		if createdAt.Valid {
			enriched["created_at"] = createdAt.Time
		}
		if updatedBy.Valid {
			enriched["updated_by"] = updatedBy.Int64
		}
		if updatedAt.Valid {
			enriched["updated_at"] = updatedAt.Time
		}

		if insCode.Valid {
			enriched["insurance_code"] = insCode.String
		} else {
			enriched["insurance_code"] = ""
		}

		if prodName.Valid {
			enriched["product_name"] = prodName.String
		} else {
			enriched["product_name"] = ""
		}

		enrichedTemplates = append(enrichedTemplates, enriched)
	}

	return c.JSON(http.StatusOK, enrichedTemplates)
}

// GetByID handles GET /api/property/templates/:id
func (h *handlerTemplate) GetByID(c echo.Context) error {
	id := c.Param("id")
	locale := c.QueryParam("locale")
	if locale == "" {
		locale = "id"
	}

	t, err := h.usecase.GetByID(locale, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, t)
}

// Create handles POST /api/property/templates
func (h *handlerTemplate) Create(c echo.Context) error {
	var req template.CreateTemplateRequest

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
		"message": "Template created successfully",
	})
}

// Update handles PUT /api/property/templates/:id
func (h *handlerTemplate) Update(c echo.Context) error {
	id := c.Param("id")
	locale := c.QueryParam("locale")
	if locale == "" {
		locale = "id"
	}

	var req template.UpdateTemplateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	req.Locale = locale
	req.ID = id

	if err := h.usecase.Update(locale, id, req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	updatedTemplate, _ := h.usecase.GetByID(locale, id)
	return c.JSON(http.StatusOK, updatedTemplate)
}

// Delete handles DELETE /api/property/templates/:id
func (h *handlerTemplate) Delete(c echo.Context) error {
	id := c.Param("id")
	locale := c.QueryParam("locale")
	if locale == "" {
		locale = "id"
	}

	err := h.usecase.Delete(locale, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// Draft operations

// AddDraft handles POST /api/property/templates/draft/add
func (h *handlerTemplate) AddDraft(c echo.Context) error {
	var draft template.TemplateDraft

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

// GetAllDrafts handles GET /api/property/templates/draft
func (h *handlerTemplate) GetAllDrafts(c echo.Context) error {
	insuranceCode := c.QueryParam("insurance_code")

	drafts, err := h.usecase.GetAllDrafts(insuranceCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if drafts == nil {
		drafts = []template.TemplateDraft{}
	}

	return c.JSON(http.StatusOK, drafts)
}

// GetDraftByID handles GET /api/property/templates/draft/:id
func (h *handlerTemplate) GetDraftByID(c echo.Context) error {
	id := c.Param("id")

	draft, err := h.usecase.GetDraftByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, draft)
}

// UpdateDraft handles PUT /api/property/templates/draft/:id
func (h *handlerTemplate) UpdateDraft(c echo.Context) error {
	id := c.Param("id")

	var draft template.TemplateDraft
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

// DeleteDraft handles DELETE /api/property/templates/draft/:id
func (h *handlerTemplate) DeleteDraft(c echo.Context) error {
	id := c.Param("id")

	if err := h.usecase.DeleteDraft(id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// ClearDrafts handles POST /api/property/templates/draft/clear
func (h *handlerTemplate) ClearDrafts(c echo.Context) error {
	if err := h.usecase.ClearDrafts(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// ConfirmDrafts handles POST /api/property/templates/draft/confirm
func (h *handlerTemplate) ConfirmDrafts(c echo.Context) error {
	var requestData struct {
		InsuranceCode string `json:"insurance_code"`
	}

	if err := c.Bind(&requestData); err != nil {
		requestData.InsuranceCode = ""
	}

	if err := h.usecase.ConfirmDrafts(requestData.InsuranceCode); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Templates confirmed successfully",
	})
}


