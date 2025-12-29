package repository

import (
	"database/sql"
	"fmt"
	"strconv"

	"multi-onboarding/domain/Travel/country"
	"multi-onboarding/utils"
)

type repository struct {
	mysqlSess *sql.DB
}

// NewCountryRepository creates a new country repository
func NewCountryRepository() country.Repository {
	return &repository{
		mysqlSess: utils.GetDB(),
	}
}

// GetAll retrieves all countries
func (r *repository) GetAll() ([]country.Country, error) {
	// Use id as code if code column doesn't exist
	query := `
		SELECT id, name, CAST(id AS CHAR) as code
		FROM countries
		ORDER BY name ASC
	`

	rows, err := r.mysqlSess.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query countries: %w", err)
	}
	defer rows.Close()

	var countries []country.Country
	for rows.Next() {
		var c country.Country
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Code,
		)
		if err != nil {
			return nil, err
		}
		countries = append(countries, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return countries, nil
}

// GetByID retrieves a country by ID
func (r *repository) GetByID(id int) (*country.Country, error) {
	query := `
		SELECT id, name, CAST(id AS CHAR) as code
		FROM countries
		WHERE id = ?
	`

	var c country.Country
	err := r.mysqlSess.QueryRow(query, id).Scan(
		&c.ID,
		&c.Name,
		&c.Code,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("country not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get country: %w", err)
	}

	return &c, nil
}

// GetByCode retrieves a country by code
func (r *repository) GetByCode(code string) (*country.Country, error) {
	// Try by id if code is numeric (since code column may not exist)
	codeInt, err := strconv.Atoi(code)
	if err == nil && codeInt > 0 {
		return r.GetByID(codeInt)
	}
	
	// If code is not numeric, try to find by code column (if exists)
	query := `
		SELECT id, name, CAST(id AS CHAR) as code
		FROM countries
		WHERE code = ?
		LIMIT 1
	`

	var c country.Country
	err = r.mysqlSess.QueryRow(query, code).Scan(
		&c.ID,
		&c.Name,
		&c.Code,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("country not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get country by code: %w", err)
	}

	return &c, nil
}

