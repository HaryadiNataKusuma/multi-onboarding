package addon

import (
	"database/sql"

	"github.com/labstack/echo"
)

// Addon represents an addon entity
type Addon struct {
	ID       int
	Code     string
	Name     string
	IsActive bool
}

// UpdateAddonRequest represents the request to update an addon
type UpdateAddonRequest struct {
	Code     string
	Name     string
	IsActive bool
}

// ConfirmAddonsRequest represents the request to confirm addons
type ConfirmAddonsRequest struct {
	Addons []struct {
		Code     string
		Name     string
		IsActive bool
	}
}

// AddonBeforeUpdate represents addon data before update (for history)
type AddonBeforeUpdate struct {
	Code     string
	Name     string
	IsActive sql.NullBool
}

// Usecase interface defines business logic operations
type Usecase interface {
	GetAddons(c echo.Context) ([]map[string]interface{}, error)
	UpdateAddon(c echo.Context, id int, req UpdateAddonRequest, userName string) error
	ConfirmAddons(c echo.Context, req ConfirmAddonsRequest) (int, error)
}

// Repository interface defines data access operations
type Repository interface {
	GetAllAddons(c echo.Context) ([]map[string]interface{}, error)
	GetAddonBeforeUpdate(c echo.Context, id int) (*AddonBeforeUpdate, error)
	UpdateAddon(c echo.Context, id int, code, name string, isActive bool) (int64, error)
	UpdateAddonCodeInAddonRules(c echo.Context, newCode, oldCode string) error
	UpdateAddonCodeInMappings(c echo.Context, newCode, oldCode string) error
	InsertHistory(c echo.Context, userName, section, action string, recordID int, recordType, dataBefore, dataAfter string) error
	InsertAddons(c echo.Context, addons []struct {
		Code     string
		Name     string
		IsActive bool
	}) (int, error)
}


