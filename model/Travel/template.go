package model

import "time"

// Template represents a template entity
type Template struct {
	Locale    string     `json:"locale"` // Locale (e.g., "id")
	ID        string     `json:"id"`     // Template ID (e.g., "INSURANCE_DETAIL_")
	Value     string     `json:"value"`  // HTML content
	CreatedBy int64      `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedBy *int64     `json:"updated_by,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// TemplateDraft represents a draft template (temporary storage)
type TemplateDraft struct {
	ID            string    `json:"id"`
	Locale        string    `json:"locale"`
	Value         string    `json:"value"`
	ProductCode   string    `json:"product_code"`
	Insurance     string    `json:"insurance"`      // Insurance code (can be array in JSON)
	InsuranceCode string    `json:"insurance_code"` // Alternative field name
	CreatedBy     int64     `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
}


