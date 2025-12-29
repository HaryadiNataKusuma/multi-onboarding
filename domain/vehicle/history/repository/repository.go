package repository

import (
	"database/sql"
	"encoding/json"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/history"
	"multi-onboarding/utils"
)

type repository struct{}

// NewHistoryRepository creates a new history repository
func NewHistoryRepository() history.Repository {
	return &repository{}
}

// getDB gets the database connection
func (r *repository) getDB(c echo.Context) *sql.DB {
	return utils.GetDB()
}

// GetAllHistories retrieves all histories from database
func (r *repository) GetAllHistories(c echo.Context, insuranceCode, section string) ([]map[string]interface{}, error) {
	db := r.getDB(c)
	var rows *sql.Rows
	var err error

	// Build query with filters
	query := `SELECT id, user_name, section, action, record_id, record_type, data_before, data_after, created_at 
		FROM vehicle_service_development.histories 
		WHERE 1=1`
	var args []interface{}

	// Filter by insurance_code
	if insuranceCode != "" {
		insuranceCodePattern := "%\"insurance_code\":\"" + insuranceCode + "\"%"
		codePattern := "%\"code\":\"" + insuranceCode + "\"%"
		query += ` AND (data_before LIKE ? OR data_after LIKE ? OR data_before LIKE ? OR data_after LIKE ?)`
		args = append(args, insuranceCodePattern, insuranceCodePattern, codePattern, codePattern)
	}

	// Filter by section
	if section != "" {
		query += ` AND section = ?`
		args = append(args, section)
	}

	query += ` ORDER BY created_at DESC`
	
	rows, err = db.Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []map[string]interface{}
	for rows.Next() {
		var id, recordID int
		var userName, section, action, recordType, createdAt string
		var dataBefore, dataAfter *string

		if err := rows.Scan(&id, &userName, &section, &action, &recordID, &recordType, &dataBefore, &dataAfter, &createdAt); err != nil {
			continue
		}

		// Extract insurance_code from JSON data
		insuranceCodeValue := ""
		extractInsuranceCode := func(data map[string]interface{}) string {
			if code, ok := data["insurance_code"].(string); ok && code != "" {
				return code
			}
			if code, ok := data["code"].(string); ok && code != "" {
				return code
			}
			return ""
		}

		// Try data_after first
		if dataAfter != nil && *dataAfter != "" {
			var afterData map[string]interface{}
			if err := json.Unmarshal([]byte(*dataAfter), &afterData); err == nil {
				insuranceCodeValue = extractInsuranceCode(afterData)
			}
		}
		// If not found, try data_before
		if insuranceCodeValue == "" && dataBefore != nil && *dataBefore != "" {
			var beforeData map[string]interface{}
			if err := json.Unmarshal([]byte(*dataBefore), &beforeData); err == nil {
				insuranceCodeValue = extractInsuranceCode(beforeData)
			}
		}

		// If still not found, try to get from database
		if insuranceCodeValue == "" && recordID > 0 {
			switch section {
			case "Products":
				var code string
				err := db.QueryRow("SELECT insurance_code FROM vehicle_service_development.products WHERE id = ?", recordID).Scan(&code)
				if err == nil && code != "" {
					insuranceCodeValue = code
				}
			case "Insurances":
				var code string
				err := db.QueryRow("SELECT code FROM vehicle_service_development.insurances WHERE id = ?", recordID).Scan(&code)
				if err == nil && code != "" {
					insuranceCodeValue = code
				}
			}
		}

		history := map[string]interface{}{
			"id":             id,
			"user_name":      userName,
			"section":        section,
			"action":         action,
			"record_id":      recordID,
			"record_type":    recordType,
			"insurance_code": insuranceCodeValue,
			"created_at":     createdAt,
		}

		if dataBefore != nil {
			history["data_before"] = *dataBefore
		}
		if dataAfter != nil {
			history["data_after"] = *dataAfter
		}

		if insuranceCodeValue == "" {
			history["insurance_code"] = nil
		}

		histories = append(histories, history)
	}

	if histories == nil {
		histories = []map[string]interface{}{}
	}

	return histories, rows.Err()
}


