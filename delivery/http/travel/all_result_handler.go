package http

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/xuri/excelize/v2"
)

// AllResultHandler handles HTTP requests for downloading all result data
type AllResultHandler struct {
	db *sql.DB
}

// NewAllResultHandler creates a new all result handler
func NewAllResultHandler(db *sql.DB) *AllResultHandler {
	return &AllResultHandler{
		db: db,
	}
}

// DownloadAllData handles GET /api/all-result/download
func (h *AllResultHandler) DownloadAllData(w http.ResponseWriter, r *http.Request) {
	// Create new Excel file
	f := excelize.NewFile()
	defer f.Close()

	// Download data from Product Rules and Addon Rules sections only
	sections := []struct {
		name     string
		query    string
		headers  []string
		getRow   func(*sql.Rows) ([]interface{}, error)
	}{
		{
			name:    "Product Rules",
			query:   "SELECT id, product_code, premium_type, start_days, end_days, base_premium_value, rules, created_at, updated_at FROM travel_service_development.product_rules ORDER BY id",
			headers: []string{"ID", "Product Code", "Premium Type", "Start Days", "End Days", "Base Premium Value", "Rules", "Created At", "Updated At"},
			getRow: func(rows *sql.Rows) ([]interface{}, error) {
				var id, startDays, endDays int
				var productCode, premiumType, rules string
				var basePremiumValue sql.NullFloat64
				var createdAt, updatedAt sql.NullTime
				err := rows.Scan(&id, &productCode, &premiumType, &startDays, &endDays, &basePremiumValue, &rules, &createdAt, &updatedAt)
				if err != nil {
					return nil, err
				}
				basePremiumVal := ""
				if basePremiumValue.Valid {
					basePremiumVal = fmt.Sprintf("%.2f", basePremiumValue.Float64)
				}
				return []interface{}{id, productCode, premiumType, startDays, endDays, basePremiumVal, rules,
					formatTime(createdAt), formatTime(updatedAt)}, nil
			},
		},
		{
			name:    "Addon Rule Details",
			query:   "SELECT id, addon_rule_id, start_condition, end_condition, value_type, value, duration_rule_type, min_adult, max_adult, max_age, created_by, created_at, updated_by, updated_at FROM travel_service_development.addon_rule_details ORDER BY id",
			headers: []string{"ID", "Addon Rule ID", "Start Condition", "End Condition", "Value Type", "Value", "Duration Rule Type", "Min Adult", "Max Adult", "Max Age", "Created By", "Created At", "Updated By", "Updated At"},
			getRow: func(rows *sql.Rows) ([]interface{}, error) {
				var id, addonRuleID, minAdult, maxAdult, maxAge, createdBy int64
				var startCondition, endCondition, valueType, durationRuleType string
				var value float64
				var createdAt sql.NullTime
				var updatedBy sql.NullInt64
				var updatedAt sql.NullTime
				err := rows.Scan(&id, &addonRuleID, &startCondition, &endCondition, &valueType, &value, &durationRuleType, &minAdult, &maxAdult, &maxAge, &createdBy, &createdAt, &updatedBy, &updatedAt)
				if err != nil {
					return nil, err
				}
				updatedByVal := ""
				if updatedBy.Valid {
					updatedByVal = fmt.Sprintf("%d", updatedBy.Int64)
				}
				return []interface{}{id, addonRuleID, startCondition, endCondition, valueType, fmt.Sprintf("%.2f", value), durationRuleType, minAdult, maxAdult, maxAge, createdBy,
					formatTime(createdAt), updatedByVal, formatTime(updatedAt)}, nil
			},
		},
		{
			name:    "Templates",
			query:   "SELECT locale, id, value, created_by, created_at, updated_by, updated_at FROM travel_service_development.templates ORDER BY created_at DESC",
			headers: []string{"Locale", "ID", "Value", "Created By", "Created At", "Updated By", "Updated At"},
			getRow: func(rows *sql.Rows) ([]interface{}, error) {
				var locale, id, value string
				var createdBy int64
				var createdAt sql.NullTime
				var updatedBy sql.NullInt64
				var updatedAt sql.NullTime
				err := rows.Scan(&locale, &id, &value, &createdBy, &createdAt, &updatedBy, &updatedAt)
				if err != nil {
					return nil, err
				}
				// Truncate value if too long for Excel (limit to first 1000 characters)
				if len(value) > 1000 {
					value = value[:1000] + "..."
				}
				updatedByVal := ""
				if updatedBy.Valid {
					updatedByVal = fmt.Sprintf("%d", updatedBy.Int64)
				}
				return []interface{}{locale, id, value, createdBy,
					formatTime(createdAt), updatedByVal, formatTime(updatedAt)}, nil
			},
		},
	}

	// Create a sheet for each section
	for idx, section := range sections {
		sheetName := section.name
		if idx == 0 {
			// Use default sheet for first section
			f.SetSheetName("Sheet1", sheetName)
		} else {
			// Create new sheet for other sections
			f.NewSheet(sheetName)
		}

		// Write headers
		for colIdx, header := range section.headers {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
			f.SetCellValue(sheetName, cell, header)
		}

		// Style headers (bold, background color)
		styleID, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
		})
		firstCell, _ := excelize.CoordinatesToCellName(1, 1)
		lastCell, _ := excelize.CoordinatesToCellName(len(section.headers), 1)
		f.SetCellStyle(sheetName, firstCell, lastCell, styleID)

		// Fetch and write data
		rows, err := h.db.Query(section.query)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch %s: %v", section.name, err))
			return
		}
		defer rows.Close()

		rowNum := 2
		for rows.Next() {
			rowData, err := section.getRow(rows)
			if err != nil {
				continue
			}
			for colIdx, val := range rowData {
				cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowNum)
				f.SetCellValue(sheetName, cell, val)
			}
			rowNum++
		}

		// Auto-fit columns (approximate)
		for colIdx := range section.headers {
			colName, _ := excelize.ColumnNumberToName(colIdx + 1)
			f.SetColWidth(sheetName, colName, colName, 15)
		}
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=all_result_data_%s.xlsx", time.Now().Format("2006-01-02")))
	w.WriteHeader(http.StatusOK)

	// Write Excel file to response
	if err := f.Write(w); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to write Excel file: "+err.Error())
		return
	}
}

// formatTime formats sql.NullTime to string
func formatTime(t sql.NullTime) string {
	if t.Valid {
		return t.Time.Format("2006-01-02 15:04:05")
	}
	return ""
}

