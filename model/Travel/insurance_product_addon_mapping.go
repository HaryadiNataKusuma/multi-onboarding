package model

import "time"

// InsuranceProductAddonMapping represents a mapping between insurance, product, and addon
type InsuranceProductAddonMapping struct {
	ID            int64     `json:"id"`
	InsuranceCode string    `json:"insurance_code"`
	ProductCode   string    `json:"product_code"`
	AddonCode     string    `json:"addon_code"`
	CreatedBy     int64     `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}




