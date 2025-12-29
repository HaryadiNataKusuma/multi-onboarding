package model

import "time"

// Insurance represents an insurance entity
type Insurance struct {
	ID        int       `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Status    int       `json:"status"` // 1 = Active, 0 = Inactive
	Logo      string    `json:"logo"`   // URL to logo image
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// InsuranceDraft represents a draft insurance (temporary storage)
type InsuranceDraft struct {
	ID        int       `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Status    int       `json:"status"`
	Logo      string    `json:"logo"`
	Timestamp time.Time `json:"timestamp"`
}


