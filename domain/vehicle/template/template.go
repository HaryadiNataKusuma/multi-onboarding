package template

import (
	"database/sql"

	"github.com/labstack/echo"
	"multi-onboarding/utils"
)

// Template represents a template entity
type Template struct {
	ID            int
	Locale        string
	TemplateID    string
	Value         string
	CreatedBy     int
	CreatedAt     string
	InsuranceCode sql.NullString
	ProductCode   sql.NullString
}

// UpdateTemplateRequest represents request to update confirmed template
type UpdateTemplateRequest struct {
	Locale    string
	ID        string
	Value     string
	CreatedBy int
}

// UpdateTemplateDraftRequest represents request to update template draft
type UpdateTemplateDraftRequest struct {
	ID        string   // ID is used as identifier (from summary/insurance_detail)
	Locale    string
	Value     []string // Array of strings
	CreatedBy int
}

// TemplateBeforeUpdate represents template data before update (for history)
type TemplateBeforeUpdate struct {
	Locale     sql.NullString
	TemplateID sql.NullString
	Value      sql.NullString
	CreatedBy  sql.NullInt64
}

// ConfirmTemplatesResult represents result of confirming templates
type ConfirmTemplatesResult struct {
	SuccessCount int
	ErrorCount   int
	TotalCount   int
	Errors       []string
}

// Usecase interface defines business logic operations
type Usecase interface {
	// Draft operations
	GetTemplatesDraft(c echo.Context, insuranceCode string) ([]utils.TemplateDraft, error)
	UpdateTemplateDraft(c echo.Context, req UpdateTemplateDraftRequest) error
	DeleteTemplateDraft(c echo.Context, templateID string) error
	ClearTemplatesDraft(c echo.Context) error
	ConfirmTemplates(c echo.Context, insuranceCodeFilter string) (*ConfirmTemplatesResult, error)

	// Confirmed template operations
	GetTemplates(c echo.Context, insuranceCode string) ([]map[string]interface{}, error)
	UpdateTemplate(c echo.Context, id string, req UpdateTemplateRequest, userName string) error
}

// Repository interface defines data access operations
type Repository interface {
	// Draft operations (JSON file)
	LoadDraftData(c echo.Context) (*utils.DraftData, error)
	SaveDraftData(c echo.Context, data *utils.DraftData) error

	// Confirmed template operations (Database)
	GetAllTemplates(c echo.Context, insuranceCode string) ([]map[string]interface{}, error)
	GetTemplateBeforeUpdate(c echo.Context, locale, id string) (*TemplateBeforeUpdate, error)
	UpdateTemplate(c echo.Context, oldLocale, oldID string, newLocale, newTemplateID, value string, createdBy int) (int64, error)
	InsertHistory(c echo.Context, userName, section, action string, recordID string, recordType, dataBefore, dataAfter string) error
	InsertTemplate(c echo.Context, locale, templateID, value string, createdBy int) error
}


