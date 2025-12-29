package model

import "time"

// AddonRule represents an addon rule entity (links product_code with addon)
type AddonRule struct {
	ID            int64      `json:"id"`
	ProductCode   string     `json:"product_code"`
	AddonCode     string     `json:"addon_code"` // Reference to addon code
	InsuranceCode *string    `json:"insurance_code,omitempty"`
	Rules         *string    `json:"rules,omitempty"`
	CreatedBy     int64      `json:"created_by"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

