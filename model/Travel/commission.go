package model

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
