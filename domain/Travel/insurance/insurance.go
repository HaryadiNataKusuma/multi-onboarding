package insurance

import "time"

// Insurance represents an insurance entity
type Insurance struct {
	ID        int       `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Status    int       `json:"status"` // 1 = Active, 0 = Inactive
	Logo      string    `json:"logo"`  // URL to logo image
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// InsuranceDraft represents a draft insurance (temporary storage)
type InsuranceDraft struct {
	ID        int       `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Status    int       `json:"status"`
	Logo      string    `json:"logo"`
	Timestamp time.Time `json:"timestamp"`
}

// CreateInsuranceRequest represents the request to create an insurance
type CreateInsuranceRequest struct {
	Code   string `json:"code" validate:"required"`
	Name   string `json:"name" validate:"required"`
	Status int    `json:"status"`
	Logo   string `json:"logo"`
}

// UpdateInsuranceRequest represents the request to update an insurance
type UpdateInsuranceRequest struct {
	Code   string `json:"code" validate:"required"`
	Name   string `json:"name" validate:"required"`
	Status int    `json:"status"`
	Logo   string `json:"logo"`
}

// InsuranceBeforeUpdate represents insurance data before update (for history)
type InsuranceBeforeUpdate struct {
	Code   string
	Name   string
	Status int
	Logo   string
}

// Usecase interface defines business logic operations
type Usecase interface {
	// Draft operations
	AddDraft(draft *InsuranceDraft) error
	GetAllDrafts() ([]InsuranceDraft, error)
	GetDraftByID(id int) (*InsuranceDraft, error)
	UpdateDraft(id int, draft *InsuranceDraft) error
	DeleteDraft(id int) error
	ClearDrafts() error
	ConfirmDrafts() error

	// Database operations
	Create(req CreateInsuranceRequest) error
	GetByID(id int) (*Insurance, error)
	GetAll() ([]Insurance, error)
	Update(id int, req UpdateInsuranceRequest) error
	Delete(id int) error
}

// Repository interface defines data access operations
type Repository interface {
	// Draft operations (in-memory)
	AddDraft(draft *InsuranceDraft) error
	GetAllDrafts() ([]InsuranceDraft, error)
	GetDraftByID(id int) (*InsuranceDraft, error)
	UpdateDraft(id int, draft *InsuranceDraft) error
	DeleteDraft(id int) error
	ClearDrafts() error
	ConfirmDrafts() error

	// Database operations
	Create(insurance *Insurance) error
	GetByID(id int) (*Insurance, error)
	GetAll() ([]Insurance, error)
	Update(id int, insurance *Insurance) error
	Delete(id int) error
	GetInsuranceBeforeUpdate(id int) (*InsuranceBeforeUpdate, error)
	InsertHistory(userName, section, action string, recordID int, recordType, dataBefore, dataAfter string) error
}

