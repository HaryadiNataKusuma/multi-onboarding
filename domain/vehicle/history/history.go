package history

import (
	"github.com/labstack/echo"
)

// Usecase interface defines business logic operations
type Usecase interface {
	GetHistories(c echo.Context, insuranceCode, section string) ([]map[string]interface{}, error)
}

// Repository interface defines data access operations
type Repository interface {
	GetAllHistories(c echo.Context, insuranceCode, section string) ([]map[string]interface{}, error)
}


