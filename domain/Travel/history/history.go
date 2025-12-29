package history

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

// GetHistoriesRequest represents the request to get histories
type GetHistoriesRequest struct {
	TableName string `json:"table_name"`
	RecordID  int    `json:"record_id"`
}

// Usecase interface defines business logic operations
type Usecase interface {
	GetAll() ([]History, error)
	GetByID(id int64) (*History, error)
	GetByTableName(tableName string) ([]History, error)
	GetByRecordID(tableName string, recordID int) ([]History, error)
}

// Repository interface defines data access operations
type Repository interface {
	GetAll() ([]History, error)
	GetByID(id int64) (*History, error)
	GetByTableName(tableName string) ([]History, error)
	GetByRecordID(tableName string, recordID int) ([]History, error)
	Create(history *History) error
}

