package model

import "time"

// ProductRule represents a product rule entity
type ProductRule struct {
	ID                int64      `json:"id"`
	ProductCode       string     `json:"product_code"`
	ProductName       *string    `json:"product_name,omitempty"`
	PremiumType       string     `json:"premium_type"`
	StartDays         int        `json:"start_days"`
	EndDays           int        `json:"end_days"`
	BasePremiumValue  float64    `json:"base_premium_value"`
	Rules             string     `json:"rules"`
	CreatedBy         int64      `json:"created_by"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedBy         *int64     `json:"updated_by,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
}

// ProductRuleDraft represents a draft product rule (temporary storage)
type ProductRuleDraft struct {
	ProductCode      string  `json:"product_code"`
	PremiumType      string  `json:"premium_type"`
	StartDays        int     `json:"start_days"`
	EndDays          int     `json:"end_days"`
	BasePremiumValue float64 `json:"base_premium_value"`
	Rules            string  `json:"rules"`
}


