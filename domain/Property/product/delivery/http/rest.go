package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"multi-onboarding/domain/Property/product"

	"github.com/labstack/echo"
)

type handlerProduct struct {
	usecase product.Usecase
}

// NewProductHandler creates a new product handler
func NewProductHandler(usecase product.Usecase) *handlerProduct {
	return &handlerProduct{
		usecase: usecase,
	}
}

// GetProducts handles GET /api/property/products
func (h *handlerProduct) GetProducts(c echo.Context) error {
	insuranceCode := c.QueryParam("insurance_code")

	products, err := h.usecase.GetAll(insuranceCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if products == nil {
		products = []product.Product{}
	}

	return c.JSON(http.StatusOK, products)
}

// GetProduct handles GET /api/property/products/:id
func (h *handlerProduct) GetProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid product ID",
		})
	}

	p, err := h.usecase.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, p)
}

// CreateProduct handles POST /api/property/products
func (h *handlerProduct) CreateProduct(c echo.Context) error {
	var req product.CreateProductRequest
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
		"message": "Product created successfully",
	})
}

// UpdateProduct handles PUT /api/property/products/:id
func (h *handlerProduct) UpdateProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid product ID",
		})
	}

	var req product.UpdateProductRequest
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

	updatedProduct, _ := h.usecase.GetByID(id)
	return c.JSON(http.StatusOK, updatedProduct)
}

// DeleteProduct handles DELETE /api/property/products/:id
func (h *handlerProduct) DeleteProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid product ID",
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

// AddDraft handles POST /api/property/products/draft/add
func (h *handlerProduct) AddDraft(c echo.Context) error {
	var draft product.ProductDraft
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

// GetAllDrafts handles GET /api/property/products/draft
func (h *handlerProduct) GetAllDrafts(c echo.Context) error {
	drafts, err := h.usecase.GetAllDrafts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if drafts == nil {
		drafts = []product.ProductDraft{}
	}

	return c.JSON(http.StatusOK, drafts)
}

// GetDraftByID handles GET /api/property/products/draft/:id
func (h *handlerProduct) GetDraftByID(c echo.Context) error {
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

// UpdateDraft handles PUT /api/property/products/draft/:id
func (h *handlerProduct) UpdateDraft(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID",
		})
	}

	// Use flexible binding to handle is_active as int or string
	var payload map[string]interface{}
	if err2 := json.NewDecoder(c.Request().Body).Decode(&payload); err2 != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body: " + err2.Error(),
		})
	}

	// Convert payload to ProductDraft struct
	var draft product.ProductDraft
	if idVal, ok := payload["id"].(float64); ok {
		draft.ID = int(idVal)
	}
	if code, ok := payload["code"].(string); ok {
		draft.Code = code
	}
	if insuranceCode, ok := payload["insurance_code"].(string); ok {
		draft.InsuranceCode = insuranceCode
	}
	if typeVal, ok := payload["type"].(string); ok {
		draft.Type = typeVal
	}
	// Handle insurance_type: can be number (float64 from JSON) or string (Hardcopy Admin Fee)
	if insuranceTypeNum, ok := payload["insurance_type"].(float64); ok {
		draft.InsuranceType = strconv.FormatFloat(insuranceTypeNum, 'f', -1, 64)
	} else if insuranceType, ok := payload["insurance_type"].(string); ok {
		draft.InsuranceType = insuranceType
	}
	// Accept both property_type (from frontend) and property_category
	if propertyCategory, ok := payload["property_category"].(string); ok {
		draft.PropertyType = propertyCategory
	} else if propertyType, ok := payload["property_type"].(string); ok {
		draft.PropertyType = propertyType
	}
	if coverageType, ok := payload["coverage_type"].(string); ok {
		draft.CoverageType = coverageType
	}
	if name, ok := payload["name"].(string); ok {
		draft.Name = name
	}
	if logo, ok := payload["logo"].(string); ok {
		draft.Logo = logo
	}
	// Handle is_active: can be int (float64 from JSON) or string
	if isActive, ok := payload["is_active"].(float64); ok {
		draft.IsActive = strconv.Itoa(int(isActive))
	} else if isActiveStr, ok := payload["is_active"].(string); ok {
		draft.IsActive = isActiveStr
	}
	if minCoverage, ok := payload["min_coverage"].(float64); ok {
		draft.MinCoverage = minCoverage
	}
	if maxCoverage, ok := payload["max_coverage"].(float64); ok {
		draft.MaxCoverage = maxCoverage
	}
	if summary, ok := payload["summary"].(string); ok {
		draft.Summary = summary
	}
	if insuranceDetail, ok := payload["insurance_detail"].(string); ok {
		draft.InsuranceDetail = insuranceDetail
	}
	if protectionDetail, ok := payload["protection_detail"].(string); ok {
		draft.ProtectionDetail = protectionDetail
	}
	if howToClaim, ok := payload["how_to_claim"].(string); ok {
		draft.HowToClaim = howToClaim
	}
	if addons, ok := payload["addons"].([]interface{}); ok {
		draft.Addons = make([]string, 0, len(addons))
		for _, v := range addons {
			if str, ok := v.(string); ok {
				draft.Addons = append(draft.Addons, str)
			}
		}
	}

	if err := h.usecase.UpdateDraft(id, &draft); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, draft)
}

// DeleteDraft handles DELETE /api/property/products/draft/:id
func (h *handlerProduct) DeleteDraft(c echo.Context) error {
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

// ClearDrafts handles POST /api/property/products/draft/clear
func (h *handlerProduct) ClearDrafts(c echo.Context) error {
	if err := h.usecase.ClearDrafts(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// ConfirmDrafts handles POST /api/property/products/draft/confirm
func (h *handlerProduct) ConfirmDrafts(c echo.Context) error {
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
		"message": "Products confirmed successfully",
	})
}

// GenerateTemplateDraftsForProducts handles POST /api/property/products/generate-templates/selected
func (h *handlerProduct) GenerateTemplateDraftsForProducts(c echo.Context) error {
	var requestData struct {
		ProductCodes []string `json:"product_codes"`
		CreatedBy    int64    `json:"created_by"`
	}

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if len(requestData.ProductCodes) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No product codes provided",
		})
	}

	if requestData.CreatedBy == 0 {
		requestData.CreatedBy = 1
	}

	count, err := h.usecase.GenerateTemplateDraftsForProducts(requestData.ProductCodes, requestData.CreatedBy)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":            "Template drafts generated successfully",
		"products_processed": count,
	})
}
