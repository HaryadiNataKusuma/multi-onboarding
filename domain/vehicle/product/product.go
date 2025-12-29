package product

import "time"

// Product represents a vehicle product entity
type Product struct {
	ID                   int64      `json:"id"`
	Code                 string     `json:"code"`
	InsuranceCode        string     `json:"insurance_code"`
	InsuranceProductCode  string     `json:"insurance_product_code"`
	VehicleType          string     `json:"vehicle_type"`
	IsInsuranceAPI       int        `json:"is_insurance_api"`
	Name                 string     `json:"name"`
	Logo                 string     `json:"logo"`
	IsActive             int        `json:"is_active"`
	IsInsurerDriven      int        `json:"is_insurer_driven"`
	IsPersonalUsage      int        `json:"is_personal_usage"`
	IsElectricVehicle    int        `json:"is_electric_vehicle"`
	Summary              string     `json:"summary"`
	Config               string     `json:"config"`
	InsuranceDetail      string     `json:"insurance_detail"`
	ProtectionDetail     string     `json:"protection_detail"`
	HowToClaim           string     `json:"how_to_claim"`
	BaseProduct          string     `json:"base_product"`
	CreatedBy            int64      `json:"created_by"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedBy            *int64     `json:"updated_by"`
	UpdatedAt            *time.Time `json:"updated_at"`
}

// ProductDraft represents a draft product (temporary storage)
type ProductDraft struct {
	ID                   int    `json:"id"`
	Code                 string `json:"code"`
	InsuranceCode        string `json:"insurance_code"`
	InsuranceProductCode string `json:"insurance_product_code"`
	VehicleType          string `json:"vehicle_type"`
	Name                 string `json:"name"`
	IsInsurerDriven      int    `json:"is_insurer_driven"`
	IsPersonalUsage      int    `json:"is_personal_usage"`
	IsElectricVehicle    int    `json:"is_electric_vehicle"`
	IsActive             int    `json:"is_active"`
	Logo                 string `json:"logo"`
	Summary              string `json:"summary"`
	Config               string `json:"config"`
	InsuranceDetail      string `json:"insurance_detail"`
	ProtectionDetail     string `json:"protection_detail"`
	HowToClaim           string `json:"how_to_claim"`
	BaseProduct          string `json:"base_product"`
}

// CreateProductRequest represents the request to create a product
type CreateProductRequest struct {
	Code                 string `json:"code" validate:"required"`
	InsuranceCode        string `json:"insurance_code"`
	InsuranceProductCode string `json:"insurance_product_code"`
	VehicleType          string `json:"vehicle_type"`
	Name                 string `json:"name" validate:"required"`
	IsInsurerDriven      int    `json:"is_insurer_driven"`
	IsPersonalUsage      int    `json:"is_personal_usage"`
	IsElectricVehicle    int    `json:"is_electric_vehicle"`
	IsActive             int    `json:"is_active"`
	Logo                 string `json:"logo"`
	Summary              string `json:"summary"`
	Config               string `json:"config"`
	InsuranceDetail      string `json:"insurance_detail"`
	ProtectionDetail     string `json:"protection_detail"`
	HowToClaim           string `json:"how_to_claim"`
	BaseProduct          string `json:"base_product"`
	CreatedBy            int64  `json:"created_by"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	Code                 string `json:"code" validate:"required"`
	InsuranceCode        string `json:"insurance_code"`
	InsuranceProductCode string `json:"insurance_product_code"`
	VehicleType          string `json:"vehicle_type"`
	Name                 string `json:"name" validate:"required"`
	IsInsurerDriven      int    `json:"is_insurer_driven"`
	IsPersonalUsage      int    `json:"is_personal_usage"`
	IsElectricVehicle    int    `json:"is_electric_vehicle"`
	IsActive             int    `json:"is_active"`
	Summary              string `json:"summary"`
	Config               string `json:"config"`
	InsuranceDetail      string `json:"insurance_detail"`
	ProtectionDetail     string `json:"protection_detail"`
	HowToClaim           string `json:"how_to_claim"`
	BaseProduct          string `json:"base_product"`
	UpdatedBy            *int64 `json:"updated_by"`
}

// ProductBeforeUpdate represents product data before update (for history)
type ProductBeforeUpdate struct {
	Code                 string
	InsuranceCode        string
	InsuranceProductCode string
	VehicleType          string
	Name                 string
	IsActive             int
	IsInsurerDriven      int
	IsPersonalUsage      int
	IsElectricVehicle    int
	BaseProduct          string
	Summary              string
	Config               string
	InsuranceDetail      string
	ProtectionDetail     string
	HowToClaim           string
	Logo                 string
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
}

