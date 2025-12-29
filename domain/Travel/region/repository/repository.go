package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"multi-onboarding/domain/Travel/region"
	"multi-onboarding/utils"
)

type repository struct {
	mysqlSess *sql.DB
}

// NewRegionRepository creates a new region repository
func NewRegionRepository() region.Repository {
	return &repository{
		mysqlSess: utils.GetDB(),
	}
}

// GetAll retrieves all regions
func (r *repository) GetAll() ([]region.Region, error) {
	query := `
		SELECT id, name, type, country_ids, created_by, created_at, updated_by, updated_at
		FROM regions
		ORDER BY created_at DESC
	`

	rows, err := r.mysqlSess.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query regions: %w", err)
	}
	defer rows.Close()

	var regions []region.Region
	for rows.Next() {
		var reg region.Region
		var countryIDsJSON []byte
		var updatedBy sql.NullInt64

		err := rows.Scan(
			&reg.ID,
			&reg.Name,
			&reg.Type,
			&countryIDsJSON,
			&reg.CreatedBy,
			&reg.CreatedAt,
			&updatedBy,
			&reg.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan region: %w", err)
		}

		if updatedBy.Valid {
			reg.UpdatedBy = &updatedBy.Int64
		}

		// Unmarshal JSON array to []int
		if err := json.Unmarshal(countryIDsJSON, &reg.CountryIDs); err != nil {
			// If unmarshal fails, set empty array instead of failing
			reg.CountryIDs = []int{}
		}

		regions = append(regions, reg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Always return an array, even if empty
	if regions == nil {
		regions = []region.Region{}
	}

	return regions, nil
}

// GetByID retrieves a region by ID
func (r *repository) GetByID(id int) (*region.Region, error) {
	query := `
		SELECT id, name, type, country_ids, created_by, created_at, updated_by, updated_at
		FROM regions
		WHERE id = ?
	`

	var reg region.Region
	var countryIDsJSON []byte
	var updatedBy sql.NullInt64

	err := r.mysqlSess.QueryRow(query, id).Scan(
		&reg.ID,
		&reg.Name,
		&reg.Type,
		&countryIDsJSON,
		&reg.CreatedBy,
		&reg.CreatedAt,
		&updatedBy,
		&reg.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("region not found")
	}
	if err != nil {
		return nil, err
	}

	if updatedBy.Valid {
		reg.UpdatedBy = &updatedBy.Int64
	}

	// Unmarshal JSON array to []int
	if err := json.Unmarshal(countryIDsJSON, &reg.CountryIDs); err != nil {
		return nil, err
	}

	return &reg, nil
}

// Create creates a new region
func (r *repository) Create(reg *region.Region) error {
	// Marshal []int to JSON
	countryIDsJSON, err := json.Marshal(reg.CountryIDs)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO regions (name, type, country_ids, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.mysqlSess.Exec(query,
		reg.Name,
		reg.Type,
		countryIDsJSON,
		reg.CreatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create region: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	reg.ID = int(id)
	reg.CreatedAt = time.Now()
	reg.UpdatedAt = time.Now()

	return nil
}

// Update updates an existing region
func (r *repository) Update(id int, reg *region.Region) error {
	// Get data before update
	beforeData, err := r.GetRegionBeforeUpdate(id)
	if err != nil {
		beforeData = &region.RegionBeforeUpdate{
			Name:       reg.Name,
			Type:       reg.Type,
			CountryIDs: reg.CountryIDs,
		}
	}

	// Build changed fields
	changedBefore := make(map[string]interface{})
	changedAfter := make(map[string]interface{})
	updateFields := []string{}
	updateValues := []interface{}{}

	if beforeData.Name != reg.Name {
		changedBefore["name"] = beforeData.Name
		changedAfter["name"] = reg.Name
		updateFields = append(updateFields, "name = ?")
		updateValues = append(updateValues, reg.Name)
	}
	if beforeData.Type != reg.Type {
		changedBefore["type"] = beforeData.Type
		changedAfter["type"] = reg.Type
		updateFields = append(updateFields, "type = ?")
		updateValues = append(updateValues, reg.Type)
	}

	// Compare country_ids
	beforeCountryIDsJSON, _ := json.Marshal(beforeData.CountryIDs)
	afterCountryIDsJSON, _ := json.Marshal(reg.CountryIDs)
	if string(beforeCountryIDsJSON) != string(afterCountryIDsJSON) {
		changedBefore["country_ids"] = beforeData.CountryIDs
		changedAfter["country_ids"] = reg.CountryIDs
		updateFields = append(updateFields, "country_ids = ?")
		updateValues = append(updateValues, afterCountryIDsJSON)
	}

	if len(updateFields) == 0 {
		return nil // No changes
	}

	var updatedBy interface{}
	if reg.UpdatedBy != nil {
		updatedBy = *reg.UpdatedBy
		updateFields = append(updateFields, "updated_by = ?")
		updateValues = append(updateValues, updatedBy)
	}

	updateFields = append(updateFields, "updated_at = NOW()")
	updateValues = append(updateValues, id)

	query := fmt.Sprintf(`
		UPDATE regions
		SET %s
		WHERE id = ?
	`, strings.Join(updateFields, ", "))

	result, err := r.mysqlSess.Exec(query, updateValues...)
	if err != nil {
		return fmt.Errorf("failed to update region: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("region not found")
	}

	reg.ID = id
	reg.UpdatedAt = time.Now()

	// Save history if there are changes
	if len(changedBefore) > 0 {
		beforeJSON, _ := json.Marshal(changedBefore)
		afterJSON, _ := json.Marshal(changedAfter)
		_ = r.InsertHistory("system", "regions", "UPDATE", id, "region", string(beforeJSON), string(afterJSON))
	}

	return nil
}

// Delete deletes a region by ID
func (r *repository) Delete(id int) error {
	query := `DELETE FROM regions WHERE id = ?`

	result, err := r.mysqlSess.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("region not found")
	}

	return nil
}

// GetRegionBeforeUpdate retrieves region data before update for history
func (r *repository) GetRegionBeforeUpdate(id int) (*region.RegionBeforeUpdate, error) {
	query := `SELECT name, type, country_ids FROM regions WHERE id = ?`

	var before region.RegionBeforeUpdate
	var countryIDsJSON []byte

	err := r.mysqlSess.QueryRow(query, id).Scan(
		&before.Name,
		&before.Type,
		&countryIDsJSON,
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON array to []int
	if err := json.Unmarshal(countryIDsJSON, &before.CountryIDs); err != nil {
		return nil, err
	}

	return &before, nil
}

// InsertHistory inserts a history record
func (r *repository) InsertHistory(userName, section, action string, recordID int, recordType, dataBefore, dataAfter string) error {
	historyQuery := `INSERT INTO histories (user_name, section, action, record_id, record_type, data_before, data_after, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	_, err := r.mysqlSess.Exec(historyQuery, userName, section, action, recordID, recordType, dataBefore, dataAfter)
	if err != nil {
		return fmt.Errorf("failed to insert history: %w", err)
	}

	return nil
}

