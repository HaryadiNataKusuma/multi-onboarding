package addonrules

import (
	"database/sql"
	"io"

	"github.com/labstack/echo"
	"multi-onboarding/utils"
)

// UpdateAddonRuleRequest represents request to update an addon rule
type UpdateAddonRuleRequest struct {
	UpdateData map[string]interface{}
}

// UpdateAddonRuleDraftRequest represents request to update a draft addon rule
type UpdateAddonRuleDraftRequest struct {
	ID            int
	AddonCode     string
	ProductCode   string
	InsuranceCode string
	Rules         string
	CreatedBy     int
}

// ConfirmAddonRulesResult represents result of confirming addon rules
type ConfirmAddonRulesResult struct {
	InsertedCount int
}

// Usecase interface defines business logic operations
type Usecase interface {
	// Template & Upload
	DownloadAddonRulesTemplate(c echo.Context, insuranceCode string) ([]byte, string, error)
	UploadAddonRulesFile(c echo.Context, file io.Reader, fileName string, insuranceCode string) (int, error)

	// Draft operations
	GetAddonRulesDraft(c echo.Context, insuranceCode string) ([]utils.AddonRuleDraft, error)
	UpdateAddonRuleDraft(c echo.Context, req UpdateAddonRuleDraftRequest) error
	DeleteAddonRuleDraft(c echo.Context, id int) error
	ClearAddonRulesDraft(c echo.Context) error
	ConfirmAddonRules(c echo.Context) (*ConfirmAddonRulesResult, error)

	// Confirmed operations
	GetAddonRules(c echo.Context, insuranceCode string) ([]map[string]interface{}, error)
	UpdateAddonRule(c echo.Context, id int, req UpdateAddonRuleRequest, userName string) error
}

// Repository interface defines data access operations
type Repository interface {
	// Template & Upload
	GetProductsByInsuranceCode(c echo.Context, insuranceCode string) ([]string, error)
	GetActiveAddons(c echo.Context) ([]string, error)
	GetLatestAddons(c echo.Context) ([]string, error)
	LoadDraftData(c echo.Context) (*utils.DraftData, error)
	SaveDraftData(c echo.Context, data *utils.DraftData) error

	// Draft operations
	GetAllAddonRulesDraft(c echo.Context, insuranceCode string) ([]utils.AddonRuleDraft, error)
	UpdateAddonRuleDraftInMemory(c echo.Context, draftData *utils.DraftData, req UpdateAddonRuleDraftRequest) error
	DeleteAddonRuleDraftFromMemory(c echo.Context, draftData *utils.DraftData, id int) error
	ClearAddonRulesDraftInMemory(c echo.Context, draftData *utils.DraftData)

	// Confirmed operations
	GetAllAddonRules(c echo.Context, insuranceCode string) ([]map[string]interface{}, error)
	GetAddonRuleBeforeUpdate(c echo.Context, id int) (*AddonRuleBeforeUpdate, error)
	UpdateAddonRule(c echo.Context, id int, updateData map[string]interface{}) (int64, error)
	InsertAddonRule(c echo.Context, addonCode, productCode, insuranceCode, rules string, createdBy int) error
	InsertHistory(c echo.Context, userName, section, action string, recordID int, recordType, dataBefore, dataAfter string) error
}

// AddonRuleBeforeUpdate represents addon rule data before update (for history)
type AddonRuleBeforeUpdate struct {
	AddonCode     sql.NullString
	ProductCode   sql.NullString
	InsuranceCode sql.NullString
	Rules         sql.NullString
	CreatedBy     sql.NullInt64
}


