package commission

import (
	"github.com/labstack/echo"
	"multi-onboarding/utils"
)

// SaveCommissionDraftRequest represents request to save commission draft
type SaveCommissionDraftRequest struct {
	BasicCommission float64
	PlanCommissions PlanCommissionData
	InsuranceCode   string
}

// PlanCommissionData represents plan commission data
type PlanCommissionData struct {
	CommissionVatType string
	AfPercentage      float64
	AfVatType         string
	AdminFee          float64
}

// BulkUpdatePlanCommissionsRequest represents request for bulk update
type BulkUpdatePlanCommissionsRequest struct {
	IDs     []int                        `json:"ids"`
	Updates BulkUpdatePlanCommissionData `json:"updates"`
}

// BulkUpdatePlanCommissionData represents bulk update data
type BulkUpdatePlanCommissionData struct {
	CommissionPercentage *float64 `json:"commission_percentage,omitempty"`
	CommissionVatType    *string  `json:"commission_vat_type,omitempty"`
	AfPercentage         *float64 `json:"af_percentage,omitempty"`
	AfVatType            *string  `json:"af_vat_type,omitempty"`
	AdminFee             *float64 `json:"admin_fee,omitempty"`
}

// GetCommissionDraftResponse represents response for get commission draft
type GetCommissionDraftResponse struct {
	Commissions     []map[string]interface{}
	PlanCommissions []map[string]interface{}
	ConfigProducts  []map[string]interface{}
}

// GetConfirmedCommissionsResponse represents response for get confirmed commissions
type GetConfirmedCommissionsResponse struct {
	Commissions     []map[string]interface{}
	PlanCommissions []map[string]interface{}
	ConfigProducts  []map[string]interface{}
}

// ConfirmCommissionsResult represents result of confirming commissions
type ConfirmCommissionsResult struct {
	TotalMoved int
}

// Usecase interface defines business logic operations
type Usecase interface {
	SaveCommissionDraft(c echo.Context, req SaveCommissionDraftRequest) error
	GetCommissionDraft(c echo.Context, insuranceCode string) (*GetCommissionDraftResponse, error)
	ConfirmCommissions(c echo.Context) (*ConfirmCommissionsResult, error)
	ClearCommissionDraft(c echo.Context) error
	GetConfirmedCommissions(c echo.Context, insuranceCode, productCode string) (*GetConfirmedCommissionsResponse, error)
	BulkUpdatePlanCommissions(c echo.Context, req BulkUpdatePlanCommissionsRequest, userName string) (int, error)
}

// Repository interface defines data access operations
type Repository interface {
	GetProductsByInsuranceCode(c echo.Context, insuranceCode string) ([]string, error)
	LoadDraftData(c echo.Context) (*utils.DraftData, error)
	SaveDraftData(c echo.Context, data *utils.DraftData) error
	// Database operations for confirming drafts
	ClearDraftTables(c echo.Context) error
	ClearDraftTablesForInsurance(c echo.Context, insuranceCode string) error
	InsertCommissionDraft(c echo.Context, productCode, insuranceCode, agentLevel string, corporateID int, basicCommission float64, note string) error
	InsertPlanCommissionDraft(c echo.Context, productCode, commissionVatType string, basicCommission, afPercentage float64, afVatType string, adminFee float64) error
	InsertConfigProductDraft(c echo.Context, level, insurer, product string, basicCommission float64) error
	GetCommissionDraft(c echo.Context, insuranceCode string) ([]map[string]interface{}, error)
	GetPlanCommissionDraft(c echo.Context, insuranceCode string) ([]map[string]interface{}, error)
	GetConfigProductDraft(c echo.Context, insuranceCode string) ([]map[string]interface{}, error)
	MoveCommissionsDraftToConfirmed(c echo.Context) (int64, error)
	MovePlanCommissionsDraftToConfirmed(c echo.Context) (int64, error)
	MoveConfigProductsDraftToConfirmed(c echo.Context) (int64, error)
	UpdateProductRulesHardcopyAdminFee(c echo.Context) error
	GetConfirmedCommissions(c echo.Context, insuranceCode, productCode string) ([]map[string]interface{}, error)
	GetConfirmedPlanCommissions(c echo.Context, insuranceCode, productCode string) ([]map[string]interface{}, error)
	GetConfirmedConfigProducts(c echo.Context, insuranceCode, productCode string) ([]map[string]interface{}, error)
	GetPlanCommissionByIDs(c echo.Context, ids []int) ([]PlanCommissionRecord, error)
	GetInsurerIDByPlanCode(c echo.Context, planCode string) (int, error)
	GetInsuranceCodeByPlanCode(c echo.Context, planCode string) (string, error)
	InsertPlanCommission(c echo.Context, planCode string, insurerID, productID int, commissionPercentage float64, commissionVatType string, afPercentage float64, afVatType string, adminFee, hardcopyFee float64, version int) error
	DeactivatePlanCommissions(c echo.Context, ids []int) error
	UpdateProductRulesHardcopyAdminFeeByPlanCode(c echo.Context, planCode string, adminFee float64) error
	InsertHistory(c echo.Context, userName, section, action string, recordID int, recordType, dataBefore, dataAfter string) error
}

// PlanCommissionRecord represents a plan commission record
type PlanCommissionRecord struct {
	ID                   int
	PlanCode             string
	ProductID            int
	InsurerID            int
	CommissionPercentage float64
	CommissionVatType    string
	AfPercentage         float64
	AfVatType            string
	AdminFee             float64
	HardcopyFee          float64
	IsActive             bool
	Version              int
}


