package repository

import (
	"database/sql"

	"multi-onboarding/model/Travel"
)

// CountryRepository interface defines the contract for country data operations
type CountryRepository interface {
	GetAll() ([]model.Country, error)
	GetByID(id int) (*model.Country, error)
}

// MySQLCountryRepository implements CountryRepository using MySQL database
type MySQLCountryRepository struct {
	db *sql.DB
}

// NewMySQLCountryRepository creates a new MySQL country repository
func NewMySQLCountryRepository(db *sql.DB) *MySQLCountryRepository {
	return &MySQLCountryRepository{
		db: db,
	}
}

// GetAll retrieves all countries
func (r *MySQLCountryRepository) GetAll() ([]model.Country, error) {
	query := `
		SELECT id, name
		FROM countries
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []model.Country
	for rows.Next() {
		var country model.Country
		err := rows.Scan(
			&country.ID,
			&country.Name,
		)
		if err != nil {
			return nil, err
		}
		countries = append(countries, country)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return countries, nil
}

// GetByID retrieves a country by ID
func (r *MySQLCountryRepository) GetByID(id int) (*model.Country, error) {
	query := `
		SELECT id, name
		FROM countries
		WHERE id = ?
	`

	var country model.Country
	err := r.db.QueryRow(query, id).Scan(
		&country.ID,
		&country.Name,
	)

	if err == sql.ErrNoRows {
		return nil, ErrCountryNotFound
	}
	if err != nil {
		return nil, err
	}

	return &country, nil
}
