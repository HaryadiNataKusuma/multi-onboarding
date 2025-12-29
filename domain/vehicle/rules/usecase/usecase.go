package usecase

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/xuri/excelize/v2"
	"multi-onboarding/domain/vehicle/rules"
	"multi-onboarding/utils"
)

type usecase struct {
	repository rules.Repository
}

// NewRulesUsecase creates a new rules usecase
func NewRulesUsecase(repository rules.Repository) rules.Usecase {
	return &usecase{
		repository: repository,
	}
}

// DownloadRulesTemplate generates Excel template for rules
func (u *usecase) DownloadRulesTemplate(c echo.Context, insuranceCode string) ([]byte, string, error) {
	log.Printf("DownloadRulesTemplate usecase called with insurance_code: %s", insuranceCode)

	products, err := u.repository.GetProductsByInsuranceCode(c, insuranceCode)
	if err != nil {
		log.Printf("ERROR: Failed to get products: %v", err)
		return nil, "", fmt.Errorf("failed to get products: %v", err)
	}

	log.Printf("Retrieved %d products for insurance_code %s", len(products), insuranceCode)

	if len(products) == 0 {
		log.Printf("WARNING: No active products found for insurance_code %s", insuranceCode)
		return nil, "", &notFoundError{Message: "No active products found for this insurance"}
	}

	// Create Excel file
	f := excelize.NewFile()
	sheetName := "Product Rules"
	f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")

	// Write header row
	headers := []string{
		"product_code", "insurance_code", "vehicle_type", "vehicle_category",
		"start_vehicle_value", "end_vehicle_value", "limit_vehicle_age",
		"admin_fee", "hardcopy_admin_fee", "region_id",
		"base_premium_value", "loading_fee_premium_value", "commercial_usage_value",
		"start_loading_age", "additional_premium", "type_additional_premium",
		"rules", "created_by", "created_at",
		"base_premium_type", "base_loading_premium_value",
	}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Generate rows for each product
	row := 2
	nonTruckCombinations := [][]int{
		{0, 0, 10},
		{100000000, 200000000, 10},
		{200000000, 300000000, 10},
		{300000000, 500000000, 10},
		{500000000, 0, 10},
	}

	for _, product := range products {
		// Get admin_fee
		adminFee, _ := u.repository.GetAdminFeeByProductCode(c, product.Code)
		hardcopyAdminFee := strconv.FormatFloat(adminFee, 'f', 0, 64)
		if adminFee == 0 {
			hardcopyAdminFee = "0"
		}

		// NON_TRUCK: 15 rows
		for regionID := 1; regionID <= 3; regionID++ {
			for _, combo := range nonTruckCombinations {
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), product.Code)
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), product.InsuranceCode)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "CAR")
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), "NON_TRUCK")
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), combo[0])
				f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), combo[1])
				f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), combo[2])
				f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), "0")
				f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), hardcopyAdminFee)
				f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), regionID)
				f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), "0")
				f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), "5")
				f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), "0")
				f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), "6")
				f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), "0")
				f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), "FIXED")
				f.SetCellValue(sheetName, fmt.Sprintf("Q%d", row), "{}")
				f.SetCellValue(sheetName, fmt.Sprintf("R%d", row), "1")
				f.SetCellValue(sheetName, fmt.Sprintf("S%d", row), "")
				f.SetCellValue(sheetName, fmt.Sprintf("T%d", row), "0")
				f.SetCellValue(sheetName, fmt.Sprintf("U%d", row), "0")
				row++
			}
		}

		// TRUCK: 3 rows
		for regionID := 1; regionID <= 3; regionID++ {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), product.Code)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), product.InsuranceCode)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "CAR")
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), "TRUCK")
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), "0")
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), "0")
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), "10")
			f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), "0")
			f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), hardcopyAdminFee)
			f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), regionID)
			f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), "0")
			f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), "5")
			f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), "0")
			f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), "6")
			f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), "0")
			f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), "FIXED")
			f.SetCellValue(sheetName, fmt.Sprintf("Q%d", row), "{}")
			f.SetCellValue(sheetName, fmt.Sprintf("R%d", row), "1")
			f.SetCellValue(sheetName, fmt.Sprintf("S%d", row), "")
			f.SetCellValue(sheetName, fmt.Sprintf("T%d", row), "0")
			f.SetCellValue(sheetName, fmt.Sprintf("U%d", row), "0")
			row++
		}
	}

	// Write to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		log.Printf("ERROR: Failed to write Excel file to buffer: %v", err)
		return nil, "", fmt.Errorf("failed to write Excel file: %v", err)
	}

	fileName := fmt.Sprintf("%s_Product_Rules_Template.xlsx", insuranceCode)
	log.Printf("SUCCESS: Generated Excel template %s (%d bytes) for %d products", fileName, buf.Len(), len(products))
	return buf.Bytes(), fileName, nil
}

// UploadRulesFile handles file upload and parsing
func (u *usecase) UploadRulesFile(c echo.Context, file io.Reader, fileName string, insuranceCode string) (int, error) {
	log.Printf("UploadRulesFile usecase called: fileName=%s, insuranceCode=%s", fileName, insuranceCode)

	// Read file into memory
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("ERROR: Failed to read file: %v", err)
		return 0, fmt.Errorf("failed to read file: %v", err)
	}

	log.Printf("Read %d bytes from file %s", len(fileBytes), fileName)

	expectedHeaders := []string{
		"product_code", "insurance_code", "vehicle_type", "vehicle_category",
		"start_vehicle_value", "end_vehicle_value", "limit_vehicle_age",
		"admin_fee", "hardcopy_admin_fee", "region_id",
		"base_premium_value", "loading_fee_premium_value", "commercial_usage_value",
		"start_loading_age", "additional_premium", "type_additional_premium",
		"rules", "created_by", "created_at", "base_premium_type", "base_loading_premium_value",
	}

	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		log.Printf("ERROR: Failed to load draft data: %v", err)
		return 0, fmt.Errorf("failed to load draft data: %v", err)
	}

	log.Printf("Loaded draft data, current rules count: %d", len(draftData.Rules))

	// Try Excel first
	f, err := excelize.OpenReader(bytes.NewReader(fileBytes))
	if err == nil {
		log.Printf("Detected Excel file format")
		// Excel file
		sheetName := f.GetSheetName(0)
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return 0, err
		}

		if len(rows) < 2 {
			return 0, &validationError{Message: "File must contain at least header and one data row"}
		}

		// Parse rows (skip header)
		newRulesCount := 0
		for rowIdx := 1; rowIdx < len(rows); rowIdx++ {
			row := rows[rowIdx]
			rule := u.parseRuleRow(row, expectedHeaders, insuranceCode, draftData)
			if rule != nil {
				draftData.Rules = append(draftData.Rules, *rule)
				newRulesCount++
			}
		}

		log.Printf("Parsed %d new rules from Excel file", newRulesCount)

		if err := u.repository.SaveDraftData(c, draftData); err != nil {
			log.Printf("ERROR: Failed to save draft data: %v", err)
			return 0, fmt.Errorf("failed to save draft data: %v", err)
		}

		log.Printf("SUCCESS: Saved %d new rules to draft data", newRulesCount)
		return newRulesCount, nil
	}

	// Try CSV
	log.Printf("Trying CSV format")
	reader := csv.NewReader(bytes.NewReader(fileBytes))
	reader.TrimLeadingSpace = true
	csvRows, err := reader.ReadAll()
	if err != nil {
		log.Printf("ERROR: Failed to parse CSV: %v", err)
		return 0, &validationError{Message: "File must be Excel (.xlsx) or CSV format"}
	}

	if len(csvRows) < 2 {
		log.Printf("ERROR: CSV file has less than 2 rows")
		return 0, &validationError{Message: "File must contain at least header and one data row"}
	}

	log.Printf("Detected CSV file format with %d rows", len(csvRows))

	// Parse CSV rows (skip header)
	newRulesCount := 0
	for rowIdx := 1; rowIdx < len(csvRows); rowIdx++ {
		row := csvRows[rowIdx]
		rule := u.parseRuleRow(row, expectedHeaders, insuranceCode, draftData)
		if rule != nil {
			draftData.Rules = append(draftData.Rules, *rule)
			newRulesCount++
		}
	}

	log.Printf("Parsed %d new rules from CSV file", newRulesCount)

	if err := u.repository.SaveDraftData(c, draftData); err != nil {
		log.Printf("ERROR: Failed to save draft data: %v", err)
		return 0, fmt.Errorf("failed to save draft data: %v", err)
	}

	log.Printf("SUCCESS: Saved %d new rules to draft data", newRulesCount)
	return newRulesCount, nil
}

// parseRuleRow parses a row from Excel/CSV into RuleDraft
func (u *usecase) parseRuleRow(row []string, expectedHeaders []string, insuranceCode string, draftData *utils.DraftData) *utils.RuleDraft {
	if len(row) < len(expectedHeaders) {
		return nil
	}

	rule := utils.RuleDraft{
		ID:                    utils.GetNextRuleID(draftData),
		ProductCode:           strings.TrimSpace(row[0]),
		InsuranceCode:         strings.TrimSpace(row[1]),
		VehicleType:           strings.TrimSpace(row[2]),
		VehicleCategory:       strings.TrimSpace(row[3]),
		TypeAdditionalPremium: strings.TrimSpace(u.getCSVValue(row, 15, "FIXED")),
		Rules:                 strings.TrimSpace(u.getCSVValue(row, 16, "{}")),
		CreatedBy:             1,
		Timestamp:             utils.GetTimestamp(),
	}

	if insuranceCode != "" {
		rule.InsuranceCode = insuranceCode
	}

	// Parse numeric fields
	if val := u.parseFloat64(u.getCSVValue(row, 4, "0")); val != nil {
		rule.StartVehicleValue = val
	}
	if val := u.parseFloat64(u.getCSVValue(row, 5, "0")); val != nil {
		rule.EndVehicleValue = val
	}
	if val := u.parseInt(u.getCSVValue(row, 6, "10")); val != nil {
		rule.LimitVehicleAge = val
	}
	if val := u.parseFloat64(u.getCSVValue(row, 7, "0")); val != nil {
		rule.AdminFee = val
	}
	if val := u.parseFloat64(u.getCSVValue(row, 8, "25000")); val != nil {
		rule.HardcopyAdminFee = val
	}
	if val := u.parseInt(u.getCSVValue(row, 9, "1")); val != nil {
		rule.RegionID = val
	}
	if val := u.parseFloat64(u.getCSVValue(row, 10, "0")); val != nil {
		rule.BasePremiumValue = val
	}
	if val := u.parseFloat64(u.getCSVValue(row, 11, "0")); val != nil {
		rule.LoadingFeePremiumValue = val
	}
	if val := u.parseFloat64(u.getCSVValue(row, 12, "0")); val != nil {
		rule.CommercialUsageValue = val
	}
	if val := u.parseInt(u.getCSVValue(row, 13, "0")); val != nil {
		rule.StartLoadingAge = val
	}
	if val := u.parseFloat64(u.getCSVValue(row, 14, "0")); val != nil {
		rule.AdditionalPremium = val
	}
	if val := u.parseFloat64(u.getCSVValue(row, 19, "0")); val != nil {
		rule.BaseLoadingPremiumValue = val
	}
	if val := u.getCSVValue(row, 18, "0"); val != "" {
		rule.BasePremiumType = val
	}

	return &rule
}

// Helper functions
func (u *usecase) getCSVValue(row []string, index int, defaultValue string) string {
	if index < len(row) && strings.TrimSpace(row[index]) != "" {
		return strings.TrimSpace(row[index])
	}
	return defaultValue
}

func (u *usecase) parseFloat64(s string) *float64 {
	if s == "" {
		return nil
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &val
}

func (u *usecase) parseInt(s string) *int {
	if s == "" {
		return nil
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &val
}

// GetRulesDraft gets all rule drafts
func (u *usecase) GetRulesDraft(c echo.Context, insuranceCode string) ([]utils.RuleDraft, error) {
	return u.repository.GetAllRulesDraft(c, insuranceCode)
}

// UpdateRuleDraft updates a draft rule
func (u *usecase) UpdateRuleDraft(c echo.Context, req rules.UpdateRuleDraftRequest) error {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return err
	}

	if err := u.repository.UpdateRuleDraftInMemory(c, draftData, req); err != nil {
		return err
	}

	return u.repository.SaveDraftData(c, draftData)
}

// DeleteRuleDraft deletes a draft rule
func (u *usecase) DeleteRuleDraft(c echo.Context, id int) error {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return err
	}

	if err := u.repository.DeleteRuleDraftFromMemory(c, draftData, id); err != nil {
		return err
	}

	return u.repository.SaveDraftData(c, draftData)
}

// ClearRulesDraft clears all draft rules
func (u *usecase) ClearRulesDraft(c echo.Context) error {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return err
	}

	u.repository.ClearRulesDraftInMemory(c, draftData)
	return u.repository.SaveDraftData(c, draftData)
}

// ConfirmRules confirms rules from draft to database
func (u *usecase) ConfirmRules(c echo.Context) (*rules.ConfirmRulesResult, error) {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return nil, err
	}

	if len(draftData.Rules) == 0 {
		return nil, &validationError{Message: "No draft rules to confirm"}
	}

	insertedCount := 0
	for _, rule := range draftData.Rules {
		if err := u.repository.InsertRule(c, rule); err != nil {
			log.Printf("Error inserting rule: %v", err)
			return nil, err
		}
		insertedCount++
	}

	// Clear draft data
	draftData.Rules = []utils.RuleDraft{}
	u.repository.SaveDraftData(c, draftData)

	return &rules.ConfirmRulesResult{
		InsertedCount: insertedCount,
	}, nil
}

// BulkUpdateRulesDraft performs bulk update on draft rules
func (u *usecase) BulkUpdateRulesDraft(c echo.Context, req rules.BulkUpdateRulesDraftRequest) (*rules.BulkUpdateResult, error) {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return nil, err
	}

	updatedCount := 0
	for i := range draftData.Rules {
		rule := &draftData.Rules[i]

		match := true
		if req.Filter.ProductCode != "" && rule.ProductCode != req.Filter.ProductCode {
			match = false
		}
		if req.Filter.RegionID != nil && rule.RegionID != nil && *rule.RegionID != *req.Filter.RegionID {
			match = false
		}
		if req.Filter.VehicleCategory != "" && rule.VehicleCategory != req.Filter.VehicleCategory {
			match = false
		}

		if match {
			if req.UpdateData.Rules != nil {
				rule.Rules = *req.UpdateData.Rules
			}
			if req.UpdateData.LimitVehicleAge != nil {
				age := *req.UpdateData.LimitVehicleAge
				rule.LimitVehicleAge = &age
			}
			if req.UpdateData.ProductCode != "" {
				rule.ProductCode = req.UpdateData.ProductCode
			}
			rule.Timestamp = utils.GetTimestamp()
			updatedCount++
		}
	}

	if updatedCount == 0 {
		return nil, &notFoundError{Message: "No rules found matching the filter criteria"}
	}

	if err := u.repository.SaveDraftData(c, draftData); err != nil {
		return nil, err
	}

	return &rules.BulkUpdateResult{
		UpdatedCount: updatedCount,
	}, nil
}

// GetRules gets all confirmed rules
func (u *usecase) GetRules(c echo.Context, insuranceCode string) ([]map[string]interface{}, error) {
	return u.repository.GetAllRules(c, insuranceCode)
}

// UpdateRule updates a confirmed rule
func (u *usecase) UpdateRule(c echo.Context, id int, req rules.UpdateRuleRequest, userName string) error {
	beforeData, errBefore := u.repository.GetRuleBeforeUpdate(c, id)
	if errBefore != nil {
		log.Printf("Error getting data before update: %v", errBefore)
		beforeData = &rules.RuleBeforeUpdate{}
	}

	rowsAffected, err := u.repository.UpdateRule(c, id, req.UpdateData)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &notFoundError{Message: "Rule not found"}
	}

	// Compare and log changed fields
	changedFields := u.compareRuleFields(beforeData, req.UpdateData)

	if len(changedFields) > 0 {
		beforeChanged := make(map[string]interface{})
		afterChanged := make(map[string]interface{})

		insuranceCode := ""
		if beforeData.InsuranceCode.Valid {
			insuranceCode = beforeData.InsuranceCode.String
		}
		beforeChanged["insurance_code"] = insuranceCode
		afterChanged["insurance_code"] = insuranceCode

		for field, changeData := range changedFields {
			changeMap := changeData.(map[string]interface{})
			beforeChanged[field] = changeMap["before"]
			afterChanged[field] = changeMap["after"]
		}

		beforeJSONBytes, _ := json.Marshal(beforeChanged)
		afterJSONBytes, _ := json.Marshal(afterChanged)

		u.repository.InsertHistory(c, userName, "Product Rules", "UPDATE", id, "product_rule",
			string(beforeJSONBytes), string(afterJSONBytes))
	}

	return nil
}

// compareRuleFields compares rule fields before and after update
func (u *usecase) compareRuleFields(before *rules.RuleBeforeUpdate, updateData map[string]interface{}) map[string]interface{} {
	changedFields := make(map[string]interface{})

	compareNullableString := func(before sql.NullString, after interface{}, fieldName string) {
		afterStr := ""
		if after != nil {
			if str, ok := after.(string); ok {
				afterStr = str
			} else {
				afterStr = fmt.Sprintf("%v", after)
			}
		}
		beforeStr := ""
		if before.Valid {
			beforeStr = before.String
		}
		if beforeStr != afterStr {
			changedFields[fieldName] = map[string]interface{}{
				"before": beforeStr,
				"after":  afterStr,
			}
		}
	}

	compareNullableFloat := func(before sql.NullFloat64, after interface{}, fieldName string) {
		afterFloat := 0.0
		if after != nil {
			if f, ok := after.(float64); ok {
				afterFloat = f
			}
		}
		beforeFloat := 0.0
		if before.Valid {
			beforeFloat = before.Float64
		}
		if beforeFloat != afterFloat {
			changedFields[fieldName] = map[string]interface{}{
				"before": beforeFloat,
				"after":  afterFloat,
			}
		}
	}

	compareNullableInt := func(before sql.NullInt64, after interface{}, fieldName string) {
		afterInt := int64(0)
		if after != nil {
			if i, ok := after.(int64); ok {
				afterInt = i
			} else if i, ok := after.(int); ok {
				afterInt = int64(i)
			} else if f, ok := after.(float64); ok {
				afterInt = int64(f)
			}
		}
		beforeInt := int64(0)
		if before.Valid {
			beforeInt = before.Int64
		}
		if beforeInt != afterInt {
			changedFields[fieldName] = map[string]interface{}{
				"before": beforeInt,
				"after":  afterInt,
			}
		}
	}

	// Compare each field
	if val, ok := updateData["start_vehicle_value"]; ok {
		compareNullableFloat(before.StartVehicleValue, val, "start_vehicle_value")
	}
	if val, ok := updateData["end_vehicle_value"]; ok {
		compareNullableFloat(before.EndVehicleValue, val, "end_vehicle_value")
	}
	if val, ok := updateData["limit_vehicle_age"]; ok {
		compareNullableInt(before.LimitVehicleAge, val, "limit_vehicle_age")
	}
	if val, ok := updateData["admin_fee"]; ok {
		compareNullableFloat(before.AdminFee, val, "admin_fee")
	}
	if val, ok := updateData["hardcopy_admin_fee"]; ok {
		compareNullableFloat(before.HardcopyAdminFee, val, "hardcopy_admin_fee")
	}
	if val, ok := updateData["region_id"]; ok {
		compareNullableInt(before.RegionID, val, "region_id")
	}
	if val, ok := updateData["base_premium_value"]; ok {
		compareNullableFloat(before.BasePremiumValue, val, "base_premium_value")
	}
	if val, ok := updateData["loading_fee_premium_value"]; ok {
		compareNullableFloat(before.LoadingFeePremiumValue, val, "loading_fee_premium_value")
	}
	if val, ok := updateData["commercial_usage_value"]; ok {
		compareNullableFloat(before.CommercialUsageValue, val, "commercial_usage_value")
	}
	if val, ok := updateData["start_loading_age"]; ok {
		compareNullableInt(before.StartLoadingAge, val, "start_loading_age")
	}
	if val, ok := updateData["additional_premium"]; ok {
		compareNullableFloat(before.AdditionalPremium, val, "additional_premium")
	}
	if val, ok := updateData["type_additional_premium"]; ok {
		compareNullableString(before.TypeAdditionalPremium, val, "type_additional_premium")
	}
	if val, ok := updateData["rules"]; ok {
		compareNullableString(before.Rules, val, "rules")
	}
	if val, ok := updateData["base_premium_type"]; ok {
		compareNullableString(before.BasePremiumType, val, "base_premium_type")
	}
	if val, ok := updateData["base_loading_premium_value"]; ok {
		compareNullableFloat(before.BaseLoadingPremiumValue, val, "base_loading_premium_value")
	}

	return changedFields
}

// DeleteRule deletes a confirmed rule
func (u *usecase) DeleteRule(c echo.Context, id int) error {
	rowsAffected, err := u.repository.DeleteRule(c, id)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &notFoundError{Message: "Rule not found"}
	}

	return nil
}

// BulkUpdateRules performs bulk update on confirmed rules
func (u *usecase) BulkUpdateRules(c echo.Context, req rules.BulkUpdateRulesRequest) (*rules.BulkUpdateResult, error) {
	filter := map[string]interface{}{
		"product_code":     req.Filter.ProductCode,
		"region_id":        req.Filter.RegionID,
		"vehicle_category": req.Filter.VehicleCategory,
	}

	updateData := map[string]interface{}{
		"rules":             req.UpdateData.Rules,
		"limit_vehicle_age": req.UpdateData.LimitVehicleAge,
		"product_code":      req.UpdateData.ProductCode,
	}

	rowsAffected, err := u.repository.BulkUpdateRules(c, filter, updateData)
	if err != nil {
		return nil, err
	}

	return &rules.BulkUpdateResult{
		UpdatedCount: int(rowsAffected),
	}, nil
}

// Error types
type notFoundError struct {
	Message string
}

func (e *notFoundError) Error() string {
	return e.Message
}

type validationError struct {
	Message string
}

func (e *validationError) Error() string {
	return e.Message
}

