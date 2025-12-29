package rules

import (
	"database/sql"
	"io"

	"github.com/labstack/echo"
	"multi-onboarding/utils"
)

// UpdateRuleRequest represents request to update a rule
type UpdateRuleRequest struct {
	UpdateData map[string]interface{}
}

// BulkUpdateRulesDraftRequest represents request for bulk update draft rules
type BulkUpdateRulesDraftRequest struct {
	Filter struct {
		ProductCode     string `json:"product_code,omitempty"`
		RegionID        *int   `json:"region_id,omitempty"`
		VehicleCategory string `json:"vehicle_category,omitempty"`
	} `json:"filter"`
	UpdateData struct {
		Rules           *string `json:"rules,omitempty"`
		LimitVehicleAge *int    `json:"limit_vehicle_age,omitempty"`
		ProductCode     string  `json:"product_code,omitempty"`
	} `json:"update_data"`
}

// BulkUpdateRulesRequest represents request for bulk update confirmed rules
type BulkUpdateRulesRequest struct {
	Filter struct {
		ProductCode     string `json:"product_code,omitempty"`
		RegionID        *int   `json:"region_id,omitempty"`
		VehicleCategory string `json:"vehicle_category,omitempty"`
	} `json:"filter"`
	UpdateData struct {
		Rules           *string `json:"rules,omitempty"`
		LimitVehicleAge *int    `json:"limit_vehicle_age,omitempty"`
		ProductCode     string  `json:"product_code,omitempty"`
	} `json:"update_data"`
}

// UpdateRuleDraftRequest represents request to update a draft rule
type UpdateRuleDraftRequest struct {
	ID                      int
	ProductCode             string
	InsuranceCode           string
	VehicleType             string
	VehicleCategory         string
	StartVehicleValue       *float64
	EndVehicleValue         *float64
	LimitVehicleAge         *int
	AdminFee                *float64
	HardcopyAdminFee        *float64
	RegionID                *int
	BasePremiumValue        *float64
	LoadingFeePremiumValue  *float64
	CommercialUsageValue    *float64
	StartLoadingAge         *int
	AdditionalPremium       *float64
	TypeAdditionalPremium   string
	Rules                   string
	BasePremiumType         string
	BaseLoadingPremiumValue *float64
}

// ConfirmRulesResult represents result of confirming rules
type ConfirmRulesResult struct {
	InsertedCount int
}

// BulkUpdateResult represents result of bulk update
type BulkUpdateResult struct {
	UpdatedCount int
}

// Usecase interface defines business logic operations
type Usecase interface {
	// Template & Upload
	DownloadRulesTemplate(c echo.Context, insuranceCode string) ([]byte, string, error)
	UploadRulesFile(c echo.Context, file io.Reader, fileName string, insuranceCode string) (int, error)

	// Draft operations
	GetRulesDraft(c echo.Context, insuranceCode string) ([]utils.RuleDraft, error)
	UpdateRuleDraft(c echo.Context, req UpdateRuleDraftRequest) error
	DeleteRuleDraft(c echo.Context, id int) error
	ClearRulesDraft(c echo.Context) error
	ConfirmRules(c echo.Context) (*ConfirmRulesResult, error)
	BulkUpdateRulesDraft(c echo.Context, req BulkUpdateRulesDraftRequest) (*BulkUpdateResult, error)

	// Confirmed operations
	GetRules(c echo.Context, insuranceCode string) ([]map[string]interface{}, error)
	UpdateRule(c echo.Context, id int, req UpdateRuleRequest, userName string) error
	DeleteRule(c echo.Context, id int) error
	BulkUpdateRules(c echo.Context, req BulkUpdateRulesRequest) (*BulkUpdateResult, error)
}

// Repository interface defines data access operations
type Repository interface {
	// Template & Upload
	GetProductsByInsuranceCode(c echo.Context, insuranceCode string) ([]ProductInfo, error)
	GetAdminFeeByProductCode(c echo.Context, productCode string) (float64, error)
	LoadDraftData(c echo.Context) (*utils.DraftData, error)
	SaveDraftData(c echo.Context, data *utils.DraftData) error

	// Draft operations
	GetAllRulesDraft(c echo.Context, insuranceCode string) ([]utils.RuleDraft, error)
	UpdateRuleDraftInMemory(c echo.Context, draftData *utils.DraftData, req UpdateRuleDraftRequest) error
	DeleteRuleDraftFromMemory(c echo.Context, draftData *utils.DraftData, id int) error
	ClearRulesDraftInMemory(c echo.Context, draftData *utils.DraftData)

	// Confirmed operations
	GetAllRules(c echo.Context, insuranceCode string) ([]map[string]interface{}, error)
	GetRuleBeforeUpdate(c echo.Context, id int) (*RuleBeforeUpdate, error)
	UpdateRule(c echo.Context, id int, updateData map[string]interface{}) (int64, error)
	DeleteRule(c echo.Context, id int) (int64, error)
	InsertRule(c echo.Context, rule utils.RuleDraft) error
	BulkUpdateRules(c echo.Context, filter map[string]interface{}, updateData map[string]interface{}) (int64, error)
	InsertHistory(c echo.Context, userName, section, action string, recordID int, recordType, dataBefore, dataAfter string) error
}

// ProductInfo represents product information
type ProductInfo struct {
	Code          string
	InsuranceCode string
}

// RuleBeforeUpdate represents rule data before update (for history)
type RuleBeforeUpdate struct {
	ProductCode             sql.NullString
	InsuranceCode           sql.NullString
	VehicleType             sql.NullString
	VehicleCategory         sql.NullString
	StartVehicleValue       sql.NullFloat64
	EndVehicleValue         sql.NullFloat64
	LimitVehicleAge         sql.NullInt64
	AdminFee                sql.NullFloat64
	HardcopyAdminFee        sql.NullFloat64
	RegionID                sql.NullInt64
	BasePremiumValue        sql.NullFloat64
	LoadingFeePremiumValue  sql.NullFloat64
	CommercialUsageValue    sql.NullFloat64
	StartLoadingAge         sql.NullInt64
	AdditionalPremium       sql.NullFloat64
	TypeAdditionalPremium   sql.NullString
	Rules                   sql.NullString
	BasePremiumType         sql.NullString
	BaseLoadingPremiumValue sql.NullFloat64
}


