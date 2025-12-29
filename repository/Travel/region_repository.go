package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"multi-onboarding/model/Travel"
)

// RegionRepository interface defines the contract for region data operations
type RegionRepository interface {
	GetAll() ([]model.Region, error)
	GetByID(id int) (*model.Region, error)
	Create(region *model.Region) error
	Update(region *model.Region) error
	Delete(id int) error
}

// MySQLRegionRepository implements RegionRepository using MySQL database
type MySQLRegionRepository struct {
	db         *sql.DB
	historyRepo HistoryRepository
}

// NewMySQLRegionRepository creates a new MySQL region repository
func NewMySQLRegionRepository(db *sql.DB) *MySQLRegionRepository {
	return &MySQLRegionRepository{
		db: db,
	}
}

// SetHistoryRepository sets the history repository for saving update histories
func (r *MySQLRegionRepository) SetHistoryRepository(historyRepo HistoryRepository) {
	r.historyRepo = historyRepo
}

// GetAll retrieves all regions
func (r *MySQLRegionRepository) GetAll() ([]model.Region, error) {
	query := `
		SELECT id, name, type, country_ids, created_at, updated_at
		FROM travel_service_development.regions
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var regions []model.Region
	for rows.Next() {
		var region model.Region
		var countryIDsJSON []byte
		err := rows.Scan(
			&region.ID,
			&region.Name,
			&region.Type,
			&countryIDsJSON,
			&region.CreatedAt,
			&region.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal JSON array to []int
		if err := json.Unmarshal(countryIDsJSON, &region.CountryIDs); err != nil {
			return nil, err
		}

		regions = append(regions, region)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return regions, nil
}

// GetByID retrieves a region by ID
func (r *MySQLRegionRepository) GetByID(id int) (*model.Region, error) {
	query := `
		SELECT id, name, type, country_ids, created_at, updated_at
		FROM travel_service_development.regions
		WHERE id = ?
	`

	var region model.Region
	var countryIDsJSON []byte
	err := r.db.QueryRow(query, id).Scan(
		&region.ID,
		&region.Name,
		&region.Type,
		&countryIDsJSON,
		&region.CreatedAt,
		&region.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrRegionNotFound
	}
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON array to []int
	if err := json.Unmarshal(countryIDsJSON, &region.CountryIDs); err != nil {
		return nil, err
	}

	return &region, nil
}

// Create creates a new region
func (r *MySQLRegionRepository) Create(region *model.Region) error {
	// Marshal []int to JSON
	countryIDsJSON, err := json.Marshal(region.CountryIDs)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO travel_service_development.regions (name, type, country_ids, created_by, created_at, updated_at)
		VALUES (?, ?, ?, 1, NOW(), NOW())
	`

	result, err := r.db.Exec(query,
		region.Name,
		region.Type,
		countryIDsJSON,
	)

	if err != nil {
		return fmt.Errorf("failed to create region: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	region.ID = int(id)
	region.CreatedAt = time.Now()
	region.UpdatedAt = time.Now()

	return nil
}

// Update updates an existing region
func (r *MySQLRegionRepository) Update(region *model.Region) error {
	// Get data before update for comparison and history
	beforeRegion, err := r.GetByID(region.ID)
	if err != nil {
		return fmt.Errorf("failed to get region before update: %w", err)
	}

	// Build changed fields map (only fields that are different)
	changedBefore := make(map[string]interface{})
	changedAfter := make(map[string]interface{})
	updateFields := []string{}
	updateValues := []interface{}{}

	// Compare each field and build update query dynamically
	if beforeRegion.Name != region.Name {
		changedBefore["name"] = beforeRegion.Name
		changedAfter["name"] = region.Name
		updateFields = append(updateFields, "name = ?")
		updateValues = append(updateValues, region.Name)
	}
	if beforeRegion.Type != region.Type {
		changedBefore["type"] = beforeRegion.Type
		changedAfter["type"] = region.Type
		updateFields = append(updateFields, "type = ?")
		updateValues = append(updateValues, region.Type)
	}
	
	// Compare country_ids (need to compare as JSON strings)
	beforeCountryIDsJSON, _ := json.Marshal(beforeRegion.CountryIDs)
	afterCountryIDsJSON, err := json.Marshal(region.CountryIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal country_ids: %w", err)
	}
	if string(beforeCountryIDsJSON) != string(afterCountryIDsJSON) {
		changedBefore["country_ids"] = beforeRegion.CountryIDs
		changedAfter["country_ids"] = region.CountryIDs
		updateFields = append(updateFields, "country_ids = ?")
		updateValues = append(updateValues, afterCountryIDsJSON)
	}

	// If no fields changed, return early (no update needed)
	if len(updateFields) == 0 {
		return nil
	}

	// Always update updated_at
	updateFields = append(updateFields, "updated_at = NOW()")
	updateValues = append(updateValues, region.ID)

	// Build dynamic UPDATE query
	query := fmt.Sprintf(`
		UPDATE travel_service_development.regions
		SET %s
		WHERE id = ?
	`, strings.Join(updateFields, ", "))

	result, err := r.db.Exec(query, updateValues...)
	if err != nil {
		return fmt.Errorf("failed to update region: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrRegionNotFound
	}

	region.UpdatedAt = time.Now()

	// Save history only if there are changes and historyRepo is available
	if r.historyRepo != nil && len(changedBefore) > 0 {
		beforeJSON, err := json.Marshal(changedBefore)
		if err != nil {
			fmt.Printf("Warning: failed to marshal changed before data: %v\n", err)
		} else {
			afterJSON, err := json.Marshal(changedAfter)
			if err != nil {
				fmt.Printf("Warning: failed to marshal changed after data: %v\n", err)
			} else {
				history := &model.History{
					UserName:   "System",
					Section:    "regions",
					Action:     "UPDATE",
					RecordID:   region.ID,
					RecordType: "region",
					DataBefore: string(beforeJSON),
					DataAfter:  string(afterJSON),
				}
				if err := r.historyRepo.Create(history); err != nil {
					// Log error but don't fail the update
					fmt.Printf("Warning: failed to save history: %v\n", err)
				} else {
					fmt.Printf("DEBUG: History saved successfully for region ID %d with %d changed field(s)\n", region.ID, len(changedBefore))
				}
			}
		}
	}

	return nil
}

// Delete deletes a region by ID
func (r *MySQLRegionRepository) Delete(id int) error {
	query := `DELETE FROM travel_service_development.regions WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRegionNotFound
	}

	return nil
}
