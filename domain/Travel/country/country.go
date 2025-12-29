package country

// Country represents a country entity
type Country struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"` // ISO country code (e.g., "ID", "US")
}

// Usecase interface defines business logic operations
type Usecase interface {
	GetAll() ([]Country, error)
	GetByID(id int) (*Country, error)
	GetByCode(code string) (*Country, error)
}

// Repository interface defines data access operations
type Repository interface {
	GetAll() ([]Country, error)
	GetByID(id int) (*Country, error)
	GetByCode(code string) (*Country, error)
}

