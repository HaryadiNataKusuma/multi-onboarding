package model

// History represents a history record
type History struct {
	ID         int64  `json:"id"`
	UserName   string `json:"user_name"`
	Section    string `json:"section"`
	Action     string `json:"action"` // e.g., "UPDATE", "INSERT", "DELETE"
	RecordID   int    `json:"record_id"`
	RecordType string `json:"record_type"`
	DataBefore string `json:"data_before"` // JSON string
	DataAfter  string `json:"data_after"`  // JSON string
	CreatedAt  string `json:"created_at"`
}

