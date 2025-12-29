package region

import "time"

// Region represents a region entity
type Region struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"` // WHITELIST or BLACKLIST
	CountryIDs []int     `json:"country_ids"` // Store as JSON array of integers
	CreatedBy  int64     `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedBy  *int64    `json:"updated_by"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateRegionRequest represents the request to create a region
type CreateRegionRequest struct {
	Name       string `json:"name" validate:"required"`
	Type       string `json:"type" validate:"required"` // WHITELIST or BLACKLIST
	CountryIDs []int  `json:"country_ids" validate:"required"`
	CreatedBy  int64  `json:"created_by"`
}

// UpdateRegionRequest represents the request to update a region
type UpdateRegionRequest struct {
	Name       string  `json:"name" validate:"required"`
	Type       string  `json:"type" validate:"required"` // WHITELIST or BLACKLIST
	CountryIDs []int   `json:"country_ids" validate:"required"`
	UpdatedBy  *int64  `json:"updated_by"`
}

// RegionBeforeUpdate represents region data before update (for history)
type RegionBeforeUpdate struct {
	Name       string
	Type       string
	CountryIDs []int
}

// Usecase interface defines business logic operations
type Usecase interface {
	GetAll() ([]Region, error)
	GetByID(id int) (*Region, error)
	Create(req CreateRegionRequest) error
	Update(id int, req UpdateRegionRequest) error
	Delete(id int) error
}

// Repository interface defines data access operations
type Repository interface {
	GetAll() ([]Region, error)
	GetByID(id int) (*Region, error)
	Create(region *Region) error
	Update(id int, region *Region) error
	Delete(id int) error
	GetRegionBeforeUpdate(id int) (*RegionBeforeUpdate, error)
	InsertHistory(userName, section, action string, recordID int, recordType, dataBefore, dataAfter string) error
}

