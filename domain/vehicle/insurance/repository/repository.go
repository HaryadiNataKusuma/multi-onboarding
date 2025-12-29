package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/insurance"
	"multi-onboarding/utils"
)

type repository struct{}

// NewInsuranceRepository creates a new insurance repository
func NewInsuranceRepository() insurance.Repository {
	return &repository{}
}

// getDB gets the database connection
func (r *repository) getDB(c echo.Context) *sql.DB {
	return utils.GetDB()
}

// LoadDraftData loads draft data from JSON file
func (r *repository) LoadDraftData(c echo.Context) (*utils.DraftData, error) {
	return utils.LoadDraftData(c)
}

// SaveDraftData saves draft data to JSON file
func (r *repository) SaveDraftData(c echo.Context, data *utils.DraftData) error {
	return utils.SaveDraftData(c, data)
}

// GetAllInsurances retrieves all insurances from database
func (r *repository) GetAllInsurances(c echo.Context) ([]map[string]interface{}, error) {
	db := r.getDB(c)
	rows, err := db.Query("SELECT id, code, name, status, logo, created_at, updated_at FROM vehicle_service_development.insurances")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var insurances []map[string]interface{}
	for rows.Next() {
		var id int
		var code, name string
		var status int
		var logo sql.NullString
		var createdAt, updatedAt sql.NullString

		if err := rows.Scan(&id, &code, &name, &status, &logo, &createdAt, &updatedAt); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		insurance := map[string]interface{}{
			"id":     id,
			"code":   code,
			"name":   name,
			"status": status,
		}

		if logo.Valid {
			insurance["logo"] = logo.String
		} else {
			insurance["logo"] = ""
		}

		if createdAt.Valid {
			insurance["created_at"] = createdAt.String
		} else {
			insurance["created_at"] = ""
		}

		if updatedAt.Valid {
			insurance["updated_at"] = updatedAt.String
		} else {
			insurance["updated_at"] = ""
		}

		insurances = append(insurances, insurance)
	}

	return insurances, rows.Err()
}

// GetInsuranceBeforeUpdate retrieves insurance data before update for history
func (r *repository) GetInsuranceBeforeUpdate(c echo.Context, id int) (*insurance.InsuranceBeforeUpdate, error) {
	db := r.getDB(c)
	beforeQuery := `SELECT code, name, status, logo FROM vehicle_service_development.insurances WHERE id = ?`

	var beforeCode, beforeName string
	var beforeStatus int
	var beforeLogo sql.NullString

	err := db.QueryRow(beforeQuery, id).Scan(&beforeCode, &beforeName, &beforeStatus, &beforeLogo)
	if err != nil {
		return nil, err
	}

	return &insurance.InsuranceBeforeUpdate{
		Code:   beforeCode,
		Name:   beforeName,
		Status: beforeStatus,
		Logo:   beforeLogo,
	}, nil
}

// UpdateInsurance updates an insurance in the database
func (r *repository) UpdateInsurance(c echo.Context, id int, code, name string, status int, logo string) (int64, error) {
	db := r.getDB(c)
	query := `UPDATE vehicle_service_development.insurances 
		SET code = ?, name = ?, status = ?, logo = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`

	result, err := db.Exec(query, strings.ToUpper(code), name, status, logo, id)
	if err != nil {
		log.Printf("Error updating insurance: %v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return 0, err
	}

	return rowsAffected, nil
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

	return nil
}

// InsertInsurance inserts an insurance into database
func (r *repository) InsertInsurance(c echo.Context, code, name string, status int, logo string) error {
	db := r.getDB(c)
	if db == nil {
		log.Printf("ERROR: Database connection is nil")
		return fmt.Errorf("database connection is nil")
	}

	// Convert status to is_active (1 if status > 0, 0 otherwise)
	isActive := 0
	if status > 0 {
		isActive = 1
	}

	// Include is_active and created_by (required fields)
	query := "INSERT INTO vehicle_service_development.insurances (code, name, status, logo, is_active, created_by) VALUES (?, ?, ?, ?, ?, 1)"
	log.Printf("Executing query: %s with params: code=%s, name=%s, status=%d, logo=%s, is_active=%d", query, code, name, status, logo, isActive)

	result, err := db.Exec(query, code, name, status, logo, isActive)
	if err != nil {
		log.Printf("ERROR inserting insurance %s: %v", code, err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Successfully inserted insurance %s, rows affected: %d", code, rowsAffected)
	return nil
}

// CheckInsuranceExistsInQuotation checks if insurance exists in quotation_service_local
func (r *repository) CheckInsuranceExistsInQuotation(c echo.Context, code string) (bool, error) {
	db := r.getDB(c)
	var existingID int
	checkQuery := "SELECT id FROM vehicle_service_development.insurances WHERE code = ?"
	err := db.QueryRow(checkQuery, code).Scan(&existingID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// InsertInsuranceToQuotation inserts insurance into quotation_service_local
func (r *repository) InsertInsuranceToQuotation(c echo.Context, code, name string, isActive int, logo string) error {
	db := r.getDB(c)
	insertQuotationQuery := `INSERT INTO vehicle_service_development.insurances 
		(code, name, is_active, is_pay_to_head_office, aliases, filterable, logo, created_at, updated_at) 
		VALUES (?, ?, ?, 1, ?, 1, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	_, err := db.Exec(insertQuotationQuery, code, name, isActive, name, logo)
	if err != nil {
		log.Printf("Error inserting into vehicle_service_development.insurances: %v", err)
		return err
	}
	return nil
}


