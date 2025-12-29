package model

import "time"

// Addon represents an addon entity
type Addon struct {
	ID        int64      `json:"id"`
	Code      string     `json:"code"`       // Product code (unique)
	Name      string     `json:"name"`       // Addon name
	NameMy    string     `json:"name_my"`    // Name in Malaysian
	NameEn    string     `json:"name_en"`    // Name in English
	IsActive  int        `json:"is_active"`  // 1 = Active, 0 = Inactive
	CreatedBy int64      `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedBy *int64     `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
}
