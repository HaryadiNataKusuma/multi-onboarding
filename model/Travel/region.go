package model

import "time"

// Region represents a region entity
type Region struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"` // WHITELIST or BLACKLIST
	CountryIDs []int     `json:"country_ids"` // Store as JSON array of integers
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}


