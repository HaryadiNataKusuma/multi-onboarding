package insuranceownrisks

import (
	"database/sql"
	"io"

	"github.com/labstack/echo"
	"multi-onboarding/utils"
)

// UpdateOwnRiskRequest represents request to update an own risk
type UpdateOwnRiskRequest struct {
	UpdateData map[string]interface{}
}

// UpdateOwnRiskDraftRequest represents request to update a draft own risk
type UpdateOwnRiskDraftRequest struct {
	ID                int
	InsuranceCode     string
	ProductType       string
	Code              string
	Title             string
	Value             *float64
	ValueType         string
	Description       *string
	IsMandatory       *bool
	IsActive          *bool
	IsElectricVehicle *int
}

// ConfirmOwnRisksResult represents result of confirming own risks
type ConfirmOwnRisksResult struct {
	InsertedCount int
}

// ProductInfo represents product information with electric vehicle flag
type ProductInfo struct {
	Code              string
	IsElectricVehicle int
}

// Usecase interface defines business logic operations
type Usecase interface {
	// Template & Upload
	DownloadOwnRisksTemplate(c echo.Context, insuranceCode string) ([]byte, string, error)
	UploadOwnRisksFile(c echo.Context, file io.Reader, fileName string, insuranceCode string) (int, error)

	// Draft operations
	GetOwnRisksDraft(c echo.Context, insuranceCode string) ([]utils.OwnRisksDraft, error)
	UpdateOwnRiskDraft(c echo.Context, req UpdateOwnRiskDraftRequest) error
	DeleteOwnRiskDraft(c echo.Context, id int) error
	ClearOwnRisksDraft(c echo.Context) error
	ConfirmOwnRisks(c echo.Context) (*ConfirmOwnRisksResult, error)

	// Confirmed operations
	GetOwnRisks(c echo.Context, insuranceCode string) ([]map[string]interface{}, error)
	UpdateOwnRisk(c echo.Context, id int, req UpdateOwnRiskRequest, userName string) error
}

// Repository interface defines data access operations
type Repository interface {
	// Template & Upload
	GetProductsByInsuranceCode(c echo.Context, insuranceCode string) ([]ProductInfo, error)
	LoadDraftData(c echo.Context) (*utils.DraftData, error)
	SaveDraftData(c echo.Context, data *utils.DraftData) error

	// Draft operations
	GetAllOwnRisksDraft(c echo.Context, insuranceCode string) ([]utils.OwnRisksDraft, error)
	UpdateOwnRiskDraftInMemory(c echo.Context, draftData *utils.DraftData, req UpdateOwnRiskDraftRequest) error
	DeleteOwnRiskDraftFromMemory(c echo.Context, draftData *utils.DraftData, id int) error
	ClearOwnRisksDraftInMemory(c echo.Context, draftData *utils.DraftData)

	// Confirmed operations
	GetAllOwnRisks(c echo.Context, insuranceCode string) ([]map[string]interface{}, error)
	GetOwnRiskBeforeUpdate(c echo.Context, id int) (*OwnRiskBeforeUpdate, error)
	UpdateOwnRisk(c echo.Context, id int, updateData map[string]interface{}) (int64, error)
	InsertOwnRisk(c echo.Context, risk utils.OwnRisksDraft) error
	InsertHistory(c echo.Context, userName, section, action string, recordID int, recordType, dataBefore, dataAfter string) error
}

// OwnRiskBeforeUpdate represents own risk data before update (for history)
type OwnRiskBeforeUpdate struct {
	InsuranceCode     sql.NullString
	ProductType       sql.NullString
	Code              sql.NullString
	Title             sql.NullString
	Value             sql.NullFloat64
	ValueType         sql.NullString
	Description       sql.NullString
	IsMandatory       sql.NullBool
	IsActive          sql.NullBool
	IsElectricVehicle sql.NullInt64
}


