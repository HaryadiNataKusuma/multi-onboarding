package model

import "time"

// ProductDraft represents a draft product (temporary storage before confirmation)
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

// TravelProduct represents a travel product entity
type TravelProduct struct {
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
