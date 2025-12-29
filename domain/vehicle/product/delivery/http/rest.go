package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/product"
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

// GetProducts handles GET /api/products/vehicle
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

// GetProduct handles GET /api/products/vehicle/:id
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

// CreateProduct handles POST /api/products/vehicle
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

// UpdateProduct handles PUT /api/products/vehicle/:id
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

// DeleteProduct handles DELETE /api/products/vehicle/:id
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

// AddDraft handles POST /api/products/vehicle/draft/add
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

// GetAllDrafts handles GET /api/products/vehicle/draft
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

// GetDraftByID handles GET /api/products/vehicle/draft/:id
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

// UpdateDraft handles PUT /api/products/vehicle/draft/:id
func (h *handlerProduct) UpdateDraft(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID",
		})
	}

	var draft product.ProductDraft
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

// DeleteDraft handles DELETE /api/products/vehicle/draft/:id
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

// ClearDrafts handles POST /api/products/vehicle/draft/clear
func (h *handlerProduct) ClearDrafts(c echo.Context) error {
	if err := h.usecase.ClearDrafts(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// ConfirmDrafts handles POST /api/products/vehicle/draft/confirm
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
