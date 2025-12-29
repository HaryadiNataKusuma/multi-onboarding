package model

// Country represents a country entity
type Country struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"` // ISO country code (e.g., "ID", "US")
}


