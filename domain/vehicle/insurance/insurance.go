package insurance

import (
	"database/sql"

	"github.com/labstack/echo"
	"multi-onboarding/utils"
)

// Insurance represents an insurance entity
type Insurance struct {
	ID        int
	Code      string
	Name      string
	Status    int
	Logo      string
	CreatedAt sql.NullString
	UpdatedAt sql.NullString
}

// SaveInsuranceDraftRequest represents request to save insurance draft
type SaveInsuranceDraftRequest struct {
	Code   string
	Name   string
	Status int
	Logo   string
}

// UpdateInsuranceDraftRequest represents request to update insurance draft
type UpdateInsuranceDraftRequest struct {
	ID     int
	Code   string
	Name   string
	Status int
	Logo   string
}

// UpdateInsuranceRequest represents request to update confirmed insurance
type UpdateInsuranceRequest struct {
	Code   string
	Name   string
	Status int
	Logo   string
}

// InsuranceBeforeUpdate represents insurance data before update (for history)
type InsuranceBeforeUpdate struct {
	Code   string
	Name   string
	Status int
	Logo   sql.NullString
}

// ConfirmInsurancesResult represents result of confirming insurances
type ConfirmInsurancesResult struct {
	SuccessCount int
}

// Usecase interface defines business logic operations
type Usecase interface {
	// Draft operations
	SaveInsuranceDraft(c echo.Context, req SaveInsuranceDraftRequest) (*utils.InsuranceDraft, error)
	GetInsurancesDraft(c echo.Context, insuranceCode string) ([]utils.InsuranceDraft, error)
	UpdateInsuranceDraft(c echo.Context, req UpdateInsuranceDraftRequest) error
	DeleteInsuranceDraft(c echo.Context, id int) error
	ClearInsurancesDraft(c echo.Context) error
	ConfirmInsurances(c echo.Context) (*ConfirmInsurancesResult, error)

	// Confirmed insurance operations
	GetInsurances(c echo.Context) ([]map[string]interface{}, error)
	UpdateInsurance(c echo.Context, id int, req UpdateInsuranceRequest, userName string) error
}

// Repository interface defines data access operations
type Repository interface {
	// Draft operations (JSON file)
	LoadDraftData(c echo.Context) (*utils.DraftData, error)
	SaveDraftData(c echo.Context, data *utils.DraftData) error

	// Confirmed insurance operations (Database)
	GetAllInsurances(c echo.Context) ([]map[string]interface{}, error)
	GetInsuranceBeforeUpdate(c echo.Context, id int) (*InsuranceBeforeUpdate, error)
	UpdateInsurance(c echo.Context, id int, code, name string, status int, logo string) (int64, error)
	InsertHistory(c echo.Context, userName, section, action string, recordID int, recordType, dataBefore, dataAfter string) error
	InsertInsurance(c echo.Context, code, name string, status int, logo string) error
	CheckInsuranceExistsInQuotation(c echo.Context, code string) (bool, error)
	InsertInsuranceToQuotation(c echo.Context, code, name string, isActive int, logo string) error
}


