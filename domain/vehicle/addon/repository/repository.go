package repository

import (
	"database/sql"
	"log"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/addon"
	"multi-onboarding/utils"
)

type repository struct{}

// NewAddonRepository creates a new addon repository
func NewAddonRepository() addon.Repository {
	return &repository{}
}

// getDB gets the database connection
func (r *repository) getDB(c echo.Context) *sql.DB {
	return utils.GetDB()
}

// GetAllAddons retrieves all addons from database
func (r *repository) GetAllAddons(c echo.Context) ([]map[string]interface{}, error) {
	db := r.getDB(c)
	query := `SELECT id, code, name, is_active 
		FROM vehicle_service_development.addons 
		ORDER BY id DESC`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying addons: %v", err)
		return nil, err
	}
	defer rows.Close()

	var addons []map[string]interface{}
	for rows.Next() {
		var id int
		var code, name string
		var isActive sql.NullBool

		err := rows.Scan(&id, &code, &name, &isActive)
		if err != nil {
			log.Printf("Error scanning addon row: %v", err)
			continue
		}

		addon := map[string]interface{}{
			"id":        id,
			"code":      code,
			"name":      name,
			"is_active": false,
		}

		if isActive.Valid {
			addon["is_active"] = isActive.Bool
		}

		addons = append(addons, addon)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating addon rows: %v", err)
		return nil, err
	}

	return addons, nil
}

// GetAddonBeforeUpdate retrieves addon data before update for history
func (r *repository) GetAddonBeforeUpdate(c echo.Context, id int) (*addon.AddonBeforeUpdate, error) {
	db := r.getDB(c)
	beforeQuery := `SELECT code, name, is_active FROM vehicle_service_development.addons WHERE id = ?`

	var beforeCode, beforeName string
	var beforeIsActive sql.NullBool

	err := db.QueryRow(beforeQuery, id).Scan(&beforeCode, &beforeName, &beforeIsActive)
	if err != nil {
		return nil, err
	}

	return &addon.AddonBeforeUpdate{
		Code:     beforeCode,
		Name:     beforeName,
		IsActive: beforeIsActive,
	}, nil
}

// UpdateAddon updates an addon in the database
func (r *repository) UpdateAddon(c echo.Context, id int, code, name string, isActive bool) (int64, error) {
	db := r.getDB(c)
	query := `UPDATE vehicle_service_development.addons 
		SET code = ?, name = ?, is_active = ? 
		WHERE id = ?`

	result, err := db.Exec(query, code, name, isActive, id)
	if err != nil {
		log.Printf("Error updating addon: %v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return 0, err
	}

	return rowsAffected, nil
}

// UpdateAddonCodeInAddonRules updates addon_code in addon_rules table
func (r *repository) UpdateAddonCodeInAddonRules(c echo.Context, newCode, oldCode string) error {
	db := r.getDB(c)
	updateAddonRulesQuery := `UPDATE vehicle_service_development.addon_rules 
		SET addon_code = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE addon_code = ?`

	_, err := db.Exec(updateAddonRulesQuery, newCode, oldCode)
	if err != nil {
		log.Printf("Error updating addon_code in addon_rules: %v", err)
		return err
	}

	return nil
}

// UpdateAddonCodeInMappings updates addon_code in insurance_product_addon_mappings table
func (r *repository) UpdateAddonCodeInMappings(c echo.Context, newCode, oldCode string) error {
	db := r.getDB(c)
	updateMappingsQuery := `UPDATE vehicle_service_development.insurance_product_addon_mappings 
		SET addon_code = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE addon_code = ?`

	_, err := db.Exec(updateMappingsQuery, newCode, oldCode)
	if err != nil {
		log.Printf("Error updating addon_code in insurance_product_addon_mappings: %v", err)
		return err
	}

	return nil
}

// InsertHistory inserts a history record
func (r *repository) InsertHistory(c echo.Context, userName, section, action string, recordID int, recordType, dataBefore, dataAfter string) error {
	db := r.getDB(c)
	historyQuery := `INSERT INTO vehicle_service_development.histories (user_name, section, action, record_id, record_type, data_before, data_after, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	_, err := db.Exec(historyQuery, userName, section, action, recordID, recordType, dataBefore, dataAfter)
	if err != nil {
		log.Printf("Error inserting history: %v", err)
		return err
	}

	log.Printf("History inserted successfully for %s ID %d", recordType, recordID)
	return nil
}

// InsertAddons inserts multiple addons in a transaction
func (r *repository) InsertAddons(c echo.Context, addons []struct {
	Code     string
	Name     string
	IsActive bool
}) (int, error) {
	db := r.getDB(c)
	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return 0, err
	}
	defer tx.Rollback()

	insertedCount := 0
	query := `INSERT INTO vehicle_service_development.addons (code, name, is_active, created_by, created_at) 
		VALUES (?, ?, ?, 1, CURRENT_TIMESTAMP)`

	for _, addon := range addons {
		_, err := tx.Exec(query, addon.Code, addon.Name, addon.IsActive)
		if err != nil {
			log.Printf("Error inserting addon: %v", err)
			return 0, err
		}
		insertedCount++
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return 0, err
	}

	return insertedCount, nil
}


