package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"multi-onboarding/domain/Travel/country"
)

type handlerCountry struct {
	usecase country.Usecase
}

// NewCountryHandler creates a new country handler
func NewCountryHandler(usecase country.Usecase) *handlerCountry {
	return &handlerCountry{
		usecase: usecase,
	}
}

// GetCountries handles GET /api/countries
func (h *handlerCountry) GetCountries(c echo.Context) error {
	countries, err := h.usecase.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if countries == nil {
		countries = []country.Country{}
	}

	return c.JSON(http.StatusOK, countries)
}

// GetCountry handles GET /api/countries/:id
func (h *handlerCountry) GetCountry(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid country ID",
		})
	}

	cntry, err := h.usecase.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, cntry)
}
