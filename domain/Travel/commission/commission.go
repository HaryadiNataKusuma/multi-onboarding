package commission

import "time"

// Commission represents a commission entity
type Commission struct {
	ID                   int64      `json:"id"`
	ProductCode          string     `json:"product_code"`
	CommissionPercentage float64    `json:"commission_percentage"`
	CommissionVATType    string     `json:"commission_vat_type"` // EXCLUSIVE, INCLUSIVE
	AFPercentage         float64    `json:"af_percentage"`
	AFVATType            string     `json:"af_vat_type"` // EXCLUSIVE, INCLUSIVE
	AdminFee             float64    `json:"admin_fee"`
	CreatedBy            int64      `json:"created_by"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedBy            *int64     `json:"updated_by"`
	UpdatedAt            *time.Time `json:"updated_at"`
}

// CommissionDraft represents a commission draft entity
type CommissionDraft struct {
	ID                   int       `json:"id"`
	ProductCode          string    `json:"product_code"`
	CommissionPercentage float64   `json:"commission_percentage"`
	CommissionVATType    string    `json:"commission_vat_type"`
	AFPercentage         float64   `json:"af_percentage"`
	AFVATType            string    `json:"af_vat_type"`
	AdminFee             float64   `json:"admin_fee"`
	CreatedAt            time.Time `json:"created_at"`
}

// CreateCommissionRequest represents the request to create a commission
type CreateCommissionRequest struct {
	ProductCode          string  `json:"product_code" validate:"required"`
	CommissionPercentage float64 `json:"commission_percentage" validate:"required"`
	CommissionVATType    string  `json:"commission_vat_type" validate:"required"`
	AFPercentage         float64 `json:"af_percentage" validate:"required"`
	AFVATType            string  `json:"af_vat_type" validate:"required"`
	AdminFee             float64 `json:"admin_fee" validate:"required"`
	CreatedBy            int64  `json:"created_by"`
}

// UpdateCommissionRequest represents the request to update a commission
type UpdateCommissionRequest struct {
	ProductCode          string  `json:"product_code" validate:"required"`
	CommissionPercentage float64 `json:"commission_percentage" validate:"required"`
	CommissionVATType    string  `json:"commission_vat_type" validate:"required"`
	AFPercentage         float64 `json:"af_percentage" validate:"required"`
	AFVATType            string  `json:"af_vat_type" validate:"required"`
	AdminFee             float64 `json:"admin_fee" validate:"required"`
	UpdatedBy            *int64  `json:"updated_by"`
}

// CommissionBeforeUpdate represents commission data before update (for history)
type CommissionBeforeUpdate struct {
	ProductCode          string
	CommissionPercentage float64
	CommissionVATType    string
	AFPercentage         float64
	AFVATType            string
	AdminFee             float64
}

// Usecase interface defines business logic operations
type Usecase interface {
	// Draft operations
	AddDraft(draft *CommissionDraft) error
	GetAllDrafts() ([]CommissionDraft, error)
	GetDraftByID(id int) (*CommissionDraft, error)
	UpdateDraft(id int, draft *CommissionDraft) error
	DeleteDraft(id int) error
	ClearDrafts() error
	ConfirmDrafts(createdBy int64) error

	// Database operations
	GetAll() ([]Commission, error)
	GetByID(id int64) (*Commission, error)
	GetByProductCode(productCode string) (*Commission, error)
	Create(req CreateCommissionRequest) error
	Update(id int64, req UpdateCommissionRequest) error
	Delete(id int64) error
}

// Repository interface defines data access operations
type Repository interface {
	// Draft operations (in-memory)
	AddDraft(draft *CommissionDraft) error
	GetAllDrafts() ([]CommissionDraft, error)
	GetDraftByID(id int) (*CommissionDraft, error)
	UpdateDraft(id int, draft *CommissionDraft) error
	DeleteDraft(id int) error
	ClearDrafts() error
	ConfirmDrafts(createdBy int64) error

	// Database operations
	GetAll() ([]Commission, error)
	GetByID(id int64) (*Commission, error)
	GetByProductCode(productCode string) (*Commission, error)
	Create(commission *Commission) error
	Update(id int64, commission *Commission) error
	Delete(id int64) error
	GetCommissionBeforeUpdate(id int64) (*CommissionBeforeUpdate, error)
	InsertHistory(userName, section, action string, recordID int64, recordType, dataBefore, dataAfter string) error
}

