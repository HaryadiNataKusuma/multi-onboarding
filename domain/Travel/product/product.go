package product

import "time"

// Product represents a travel product entity
type Product struct {
	ID               int64      `json:"id"`
	Code             string     `json:"code"`
	InsuranceCode    string     `json:"insurance_code"`
	Name             string     `json:"name"`
	Logo             string     `json:"logo"`
	Type             string     `json:"type"`      // COUPLE, INDIVIDUAL, FAMILY
	IsActive         int        `json:"is_active"` // 1 = Active, 0 = Inactive
	InsuranceType    string     `json:"insurance_type"`
	AdminFee         float64    `json:"admin_fee"`
	RegionID         int        `json:"region_id"`
	SchengenEligible int        `json:"schengen_eligible"`
	MinAdult         int        `json:"min_adult"`
	MaxAdult         int        `json:"max_adult"`
	MinChild         int        `json:"min_child"`
	MaxChild         int        `json:"max_child"`
	Summary          string     `json:"summary"`
	InsuranceDetail  string     `json:"insurance_detail"`
	ProtectionDetail string     `json:"protection_detail"`
	HowToClaim       string     `json:"how_to_claim"`
	CreatedBy        int64      `json:"created_by"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedBy        *int64     `json:"updated_by"`
	UpdatedAt        *time.Time `json:"updated_at"`
	Partner          string     `json:"partner"`
}

// ProductDraft represents a draft product (temporary storage)
type ProductDraft struct {
	ID              int      `json:"id"`
	Code            string   `json:"code"`
	InsuranceCode   string   `json:"insurance_code"`
	Region          string   `json:"region"`
	RegionID        int      `json:"region_id"`
	Type            string   `json:"type"`
	InsuranceType   string   `json:"insurance_type"`
	ProtectionType  string   `json:"protection_type"`
	Name            string   `json:"name"`
	Logo            string   `json:"logo"`
	IsActive        string   `json:"is_active"`
	MinAdult        int      `json:"min_adult"`
	MaxAdult        int      `json:"max_adult"`
	MinChild        int      `json:"min_child"`
	MaxChild        int      `json:"max_child"`
	Summary         string   `json:"summary"`
	InsuranceDetail string   `json:"insurance_detail"`
	ProtectionDetail string  `json:"protection_detail"`
	HowToClaim      string   `json:"how_to_claim"`
	Addons          []string `json:"addons"`
}

// CreateProductRequest represents the request to create a product
type CreateProductRequest struct {
	Code             string  `json:"code" validate:"required"`
	InsuranceCode    string  `json:"insurance_code"`
	Name             string  `json:"name" validate:"required"`
	Type             string  `json:"type"`
	IsActive         int     `json:"is_active"`
	InsuranceType    string  `json:"insurance_type"`
	AdminFee         float64 `json:"admin_fee"`
	RegionID         int     `json:"region_id"`
	SchengenEligible int     `json:"schengen_eligible"`
	MinAdult         int     `json:"min_adult"`
	MaxAdult         int     `json:"max_adult"`
	MinChild         int     `json:"min_child"`
	MaxChild         int     `json:"max_child"`
	Summary          string  `json:"summary"`
	InsuranceDetail  string  `json:"insurance_detail"`
	ProtectionDetail string  `json:"protection_detail"`
	HowToClaim       string  `json:"how_to_claim"`
	CreatedBy        int64   `json:"created_by"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	Code             string  `json:"code" validate:"required"`
	InsuranceCode    string  `json:"insurance_code"`
	Name             string  `json:"name" validate:"required"`
	Type             string  `json:"type"`
	IsActive         int     `json:"is_active"`
	InsuranceType    string  `json:"insurance_type"`
	AdminFee         float64 `json:"admin_fee"`
	RegionID         int     `json:"region_id"`
	SchengenEligible int     `json:"schengen_eligible"`
	MinAdult         int     `json:"min_adult"`
	MaxAdult         int     `json:"max_adult"`
	MinChild         int     `json:"min_child"`
	MaxChild         int     `json:"max_child"`
	Summary          string  `json:"summary"`
	InsuranceDetail  string  `json:"insurance_detail"`
	ProtectionDetail string  `json:"protection_detail"`
	HowToClaim       string  `json:"how_to_claim"`
	UpdatedBy        *int64  `json:"updated_by"`
}

// ProductBeforeUpdate represents product data before update (for history)
type ProductBeforeUpdate struct {
	Code             string
	InsuranceCode    string
	Name             string
	Type             string
	IsActive         int
	InsuranceType    string
	AdminFee         float64
	RegionID         int
	SchengenEligible int
	MinAdult         int
	MaxAdult         int
	MinChild         int
	MaxChild         int
	Summary          string
	InsuranceDetail  string
	ProtectionDetail string
	HowToClaim       string
}

// Usecase interface defines business logic operations
type Usecase interface {
	// Draft operations
	AddDraft(draft *ProductDraft) error
	GetAllDrafts() ([]ProductDraft, error)
	GetDraftByID(id int) (*ProductDraft, error)
	UpdateDraft(id int, draft *ProductDraft) error
	DeleteDraft(id int) error
	ClearDrafts() error
	ConfirmDrafts(createdBy int64) error

	// Database operations
	GetAll(insuranceCode string) ([]Product, error)
	GetByID(id int) (*Product, error)
	Create(req CreateProductRequest) error
	Update(id int, req UpdateProductRequest) error
	Delete(id int) error
	GenerateTemplateDraftsForAllProducts(createdBy int64) (int, error)
	GenerateTemplateDraftsForProducts(productCodes []string, createdBy int64) (int, error)
}

// Repository interface defines data access operations
type Repository interface {
	// Draft operations (in-memory)
	AddDraft(draft *ProductDraft) error
	GetAllDrafts() ([]ProductDraft, error)
	GetDraftByID(id int) (*ProductDraft, error)
	UpdateDraft(id int, draft *ProductDraft) error
	DeleteDraft(id int) error
	ClearDrafts() error
	ConfirmDrafts(createdBy int64) error

	// Database operations
	GetAll(insuranceCode string) ([]Product, error)
	GetByID(id int) (*Product, error)
	Create(product *Product) error
	Update(id int, product *Product) error
	Delete(id int) error
	GetProductBeforeUpdate(id int) (*ProductBeforeUpdate, error)
	InsertHistory(userName, section, action string, recordID int, recordType, dataBefore, dataAfter string) error
	GenerateTemplateDraftsForAllProducts(createdBy int64) (int, error)
	GenerateTemplateDraftsForProducts(productCodes []string, createdBy int64) (int, error)
}

