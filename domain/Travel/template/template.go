package template

import "time"

// Template represents a template entity
type Template struct {
	Locale    string     `json:"locale"` // Locale (e.g., "id")
	ID        string     `json:"id"`     // Template ID (e.g., "INSURANCE_DETAIL_")
	Value     string     `json:"value"`  // HTML content
	CreatedBy int64      `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedBy *int64     `json:"updated_by,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// TemplateDraft represents a draft template (temporary storage)
type TemplateDraft struct {
	ID            string    `json:"id"`
	Locale        string    `json:"locale"`
	Value         string    `json:"value"`
	ProductCode   string    `json:"product_code"`
	Insurance     string    `json:"insurance"`      // Insurance code (can be array in JSON)
	InsuranceCode string    `json:"insurance_code"` // Alternative field name
	CreatedBy     int64     `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateTemplateRequest represents the request to create a template
type CreateTemplateRequest struct {
	Locale    string `json:"locale" validate:"required"`
	ID        string `json:"id" validate:"required"`
	Value     string `json:"value" validate:"required"`
	CreatedBy int64  `json:"created_by"`
}

// UpdateTemplateRequest represents the request to update a template
type UpdateTemplateRequest struct {
	Locale    string `json:"locale" validate:"required"`
	ID        string `json:"id" validate:"required"`
	Value     string `json:"value" validate:"required"`
	UpdatedBy *int64 `json:"updated_by"`
}

// TemplateBeforeUpdate represents template data before update (for history)
type TemplateBeforeUpdate struct {
	Locale string
	ID     string
	Value  string
}

// Usecase interface defines business logic operations
type Usecase interface {
	// Draft operations
	AddDraft(draft *TemplateDraft) error
	GetAllDrafts(insuranceCode string) ([]TemplateDraft, error)
	GetDraftByID(id string) (*TemplateDraft, error)
	UpdateDraft(id string, draft *TemplateDraft) error
	DeleteDraft(id string) error
	ClearDrafts() error
	ConfirmDrafts(insuranceCodeFilter string) error

	// Database operations
	GetAll(insuranceCode string) ([]Template, error)
	GetByID(locale, id string) (*Template, error)
	Create(req CreateTemplateRequest) error
	Update(locale, id string, req UpdateTemplateRequest) error
	Delete(locale, id string) error
}

// Repository interface defines data access operations
type Repository interface {
	// Draft operations (in-memory)
	AddDraft(draft *TemplateDraft) error
	GetAllDrafts(insuranceCode string) ([]TemplateDraft, error)
	GetDraftByID(id string) (*TemplateDraft, error)
	UpdateDraft(id string, draft *TemplateDraft) error
	DeleteDraft(id string) error
	ClearDrafts() error
	ConfirmDrafts(insuranceCodeFilter string) error

	// Database operations
	GetAll(insuranceCode string) ([]Template, error)
	GetByID(locale, id string) (*Template, error)
	Create(template *Template) error
	Update(locale, id string, template *Template) error
	Delete(locale, id string) error
	GetTemplateBeforeUpdate(locale, id string) (*TemplateBeforeUpdate, error)
	InsertHistory(userName, section, action string, recordID string, recordType, dataBefore, dataAfter string) error
}
