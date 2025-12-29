package addon

import "time"

// Addon represents an addon entity
type Addon struct {
	ID        int64      `json:"id"`
	Code      string     `json:"code"`      // Product code (unique)
	Name      string     `json:"name"`      // Addon name
	NameMy    string     `json:"name_my"`   // Name in Malaysian
	NameEn    string     `json:"name_en"`   // Name in English
	IsActive  int        `json:"is_active"` // 1 = Active, 0 = Inactive
	CreatedBy int64      `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedBy *int64     `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// CreateAddonRequest represents the request to create an addon
type CreateAddonRequest struct {
	Code      string `json:"code" validate:"required"`
	Name      string `json:"name" validate:"required"`
	NameMy    string `json:"name_my" validate:"required"`
	NameEn    string `json:"name_en" validate:"required"`
	IsActive  int    `json:"is_active"`
	CreatedBy int64  `json:"created_by"`
}

// UpdateAddonRequest represents the request to update an addon
type UpdateAddonRequest struct {
	Code      string `json:"code" validate:"required"`
	Name      string `json:"name" validate:"required"`
	NameMy    string `json:"name_my" validate:"required"`
	NameEn    string `json:"name_en" validate:"required"`
	IsActive  int    `json:"is_active"`
	UpdatedBy *int64 `json:"updated_by"`
}

// AddonBeforeUpdate represents addon data before update (for history)
type AddonBeforeUpdate struct {
	Code     string
	Name     string
	NameMy   string
	NameEn   string
	IsActive int
}

// Usecase interface defines business logic operations
type Usecase interface {
	GetAllAddons() ([]Addon, error)
	GetAddonByID(id int64) (*Addon, error)
	GetAddonByCode(code string) (*Addon, error)
	CreateAddon(req CreateAddonRequest) error
	UpdateAddon(id int64, req UpdateAddonRequest) error
	DeleteAddon(id int64) error
}

// Repository interface defines data access operations
type Repository interface {
	GetAll() ([]Addon, error)
	GetByID(id int64) (*Addon, error)
	GetByCode(code string) (*Addon, error)
	Create(addon *Addon) error
	Update(addon *Addon) error
	Delete(id int64) error
	GetAddonBeforeUpdate(id int64) (*AddonBeforeUpdate, error)
	InsertHistory(userName, section, action string, recordID int64, recordType, dataBefore, dataAfter string) error
}
