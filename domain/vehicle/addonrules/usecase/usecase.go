package usecase

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/xuri/excelize/v2"
	"multi-onboarding/domain/vehicle/addonrules"
	"multi-onboarding/utils"
)

type usecase struct {
	repository addonrules.Repository
}

// NewAddonRulesUsecase creates a new addon rules usecase
func NewAddonRulesUsecase(repository addonrules.Repository) addonrules.Usecase {
	return &usecase{
		repository: repository,
	}
}

// DownloadAddonRulesTemplate generates Excel template for addon rules
func (u *usecase) DownloadAddonRulesTemplate(c echo.Context, insuranceCode string) ([]byte, string, error) {
	products, err := u.repository.GetProductsByInsuranceCode(c, insuranceCode)
	if err != nil {
		return nil, "", err
	}

	if len(products) == 0 {
		return nil, "", &notFoundError{Message: "No active products found for this insurance"}
	}

	// Get active addons
	activeAddons, err := u.repository.GetActiveAddons(c)
	if err != nil {
		return nil, "", err
	}

	// Get latest addons
	latestAddons, _ := u.repository.GetLatestAddons(c)

	// Combine and deduplicate
	addonMap := make(map[string]bool)
	var addons []string
	for _, code := range activeAddons {
		if !addonMap[code] {
			addons = append(addons, code)
			addonMap[code] = true
		}
	}
	for _, code := range latestAddons {
		if !addonMap[code] {
			addons = append(addons, code)
			addonMap[code] = true
		}
	}

	if len(addons) == 0 {
		return nil, "", &notFoundError{Message: "No active addons found"}
	}

	// Create Excel file
	f := excelize.NewFile()
	sheetName := "Addon Rules"
	f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")

	// Write header row
	headers := []string{"addon_code", "product_code", "insurance_code", "rules", "created_by"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	row := 2
	writeExcelRow := func(addonCode, productCode, insuranceCode, rules, createdBy string) {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), addonCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), productCode)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), insuranceCode)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), rules)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), createdBy)
		row++
	}

	// Generate rows based on addon type
	for _, addonCode := range addons {
		switch addonCode {
		case "TERRORISM_SABOTAGE":
			for _, productCode := range products {
				writeExcelRow(addonCode, productCode, insuranceCode, `{"category":"","region_id":-1,"start_tsi":0,"end_tsi":0,"start_vehicle_age":-1,"end_vehicle_age":-1,"addon_value":{"value":0,"value_type":"PERCENTAGE"},"commercial_usage_value":{"value":0,"value_type":"PERCENTAGE"}}`, "1")
			}
		case "TYPHOON_STORM_FLOOD_HAIL_LANDSLIDE":
			for _, productCode := range products {
				writeExcelRow(addonCode, productCode, insuranceCode, `{"category":"","region_id":1,"start_tsi":-1,"end_tsi":-1,"start_vehicle_age":-1,"end_vehicle_age":-1,"addon_value":{"value":0.05,"value_type":"PERCENTAGE"},"commercial_usage_value":{"value":0,"value_type":"PERCENTAGE"}}`, "1")
				writeExcelRow(addonCode, productCode, insuranceCode, `{"category":"","region_id":2,"start_tsi":-1,"end_tsi":-1,"start_vehicle_age":-1,"end_vehicle_age":-1,"addon_value":{"value":0.075,"value_type":"PERCENTAGE"},"commercial_usage_value":{"value":0,"value_type":"PERCENTAGE"}}`, "1")
				writeExcelRow(addonCode, productCode, insuranceCode, `{"category":"","region_id":3,"start_tsi":-1,"end_tsi":-1,"start_vehicle_age":-1,"end_vehicle_age":-1,"addon_value":{"value":0.05,"value_type":"PERCENTAGE"},"commercial_usage_value":{"value":0,"value_type":"PERCENTAGE"}}`, "1")
			}
		case "EARTHQUAKE_TSUNAMI_VOLCANIC_ERUPTION":
			for _, productCode := range products {
				writeExcelRow(addonCode, productCode, insuranceCode, `{"category":"","region_id":1,"start_tsi":-1,"end_tsi":-1,"start_vehicle_age":-1,"end_vehicle_age":-1,"addon_value":{"value":0.085,"value_type":"PERCENTAGE"},"commercial_usage_value":{"value":0,"value_type":"PERCENTAGE"}}`, "1")
				writeExcelRow(addonCode, productCode, insuranceCode, `{"category":"","region_id":2,"start_tsi":-1,"end_tsi":-1,"start_vehicle_age":-1,"end_vehicle_age":-1,"addon_value":{"value":0.075,"value_type":"PERCENTAGE"},"commercial_usage_value":{"value":0,"value_type":"PERCENTAGE"}}`, "1")
				writeExcelRow(addonCode, productCode, insuranceCode, `{"category":"","region_id":3,"start_tsi":-1,"end_tsi":-1,"start_vehicle_age":-1,"end_vehicle_age":-1,"addon_value":{"value":0.05,"value_type":"PERCENTAGE"},"commercial_usage_value":{"value":0,"value_type":"PERCENTAGE"}}`, "1")
			}
		case "AUTHORIZED_WORKSHOP", "TRANSPORT_FEE", "PERSONAL_ACCIDENT_DRIVER", "PERSONAL_ACCIDENT_PASSENGER", "THEFT_BY_OWN_DRIVER", "WATER_HAMMER":
			for _, productCode := range products {
				writeExcelRow(addonCode, productCode, insuranceCode, `{"category":"","region_id":-1,"start_tsi":0,"end_tsi":0,"start_vehicle_age":-1,"end_vehicle_age":-1,"addon_value":{"value":0,"value_type":"PERCENTAGE"},"commercial_usage_value":{"value":0,"value_type":"PERCENTAGE"}}`, "1")
			}
		case "THIRD_PARTY_LIABILITY":
			for _, productCode := range products {
				writeExcelRow(addonCode, productCode, insuranceCode, `{"category":"NON_TRUCK","region_id":-1,"start_tsi":0,"end_tsi":25000000,"start_vehicle_age":-1,"end_vehicle_age":-1,"addon_value":{"value":1,"value_type":"PERCENTAGE"},"max_protection_value":0,"commercial_usage_value":{"value":0,"value_type":"PERCENTAGE"}}`, "1")
				writeExcelRow(addonCode, productCode, insuranceCode, `{"category":"NON_TRUCK","region_id":-1,"start_tsi":25000001,"end_tsi":50000000,"start_vehicle_age":-1,"end_vehicle_age":-1,"addon_value":{"value":0.5,"value_type":"PERCENTAGE"},"max_protection_value":0,"commercial_usage_value":{"value":0,"value_type":"PERCENTAGE"}}`, "1")
				writeExcelRow(addonCode, productCode, insuranceCode, `{"category":"NON_TRUCK","region_id":-1,"start_tsi":50000001,"end_tsi":100000000,"start_vehicle_age":-1,"end_vehicle_age":-1,"addon_value":{"value":0.25,"value_type":"PERCENTAGE"},"max_protection_value":0,"commercial_usage_value":{"value":0,"value_type":"PERCENTAGE"}}`, "1")
				writeExcelRow(addonCode, productCode, insuranceCode, `{"category":"TRUCK","region_id":-1,"start_tsi":0,"end_tsi":25000000,"start_vehicle_age":-1,"end_vehicle_age":-1,"addon_value":{"value":1.5,"value_type":"PERCENTAGE"},"max_protection_value":0,"commercial_usage_value":{"value":0,"value_type":"PERCENTAGE"}}`, "1")
				writeExcelRow(addonCode, productCode, insuranceCode, `{"category":"TRUCK","region_id":-1,"start_tsi":25000001,"end_tsi":50000000,"start_vehicle_age":-1,"end_vehicle_age":-1,"addon_value":{"value":0.75,"value_type":"PERCENTAGE"},"max_protection_value":0,"commercial_usage_value":{"value":0,"value_type":"PERCENTAGE"}}`, "1")
				writeExcelRow(addonCode, productCode, insuranceCode, `{"category":"TRUCK","region_id":-1,"start_tsi":50000001,"end_tsi":100000000,"start_vehicle_age":-1,"end_vehicle_age":-1,"addon_value":{"value":0.375,"value_type":"PERCENTAGE"},"max_protection_value":0,"commercial_usage_value":{"value":0,"value_type":"PERCENTAGE"}}`, "1")
			}
		default:
			// Default: 1 row per product
			for _, productCode := range products {
				writeExcelRow(addonCode, productCode, insuranceCode, `{"category":"","region_id":-1,"start_tsi":0,"end_tsi":0,"start_vehicle_age":-1,"end_vehicle_age":-1,"addon_value":{"value":0,"value_type":"PERCENTAGE"},"commercial_usage_value":{"value":0,"value_type":"PERCENTAGE"}}`, "1")
			}
		}
	}

	// Write to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", err
	}

	fileName := fmt.Sprintf("%s_Addon_Rules_Template.xlsx", insuranceCode)
	return buf.Bytes(), fileName, nil
}

// UploadAddonRulesFile handles file upload and parsing
func (u *usecase) UploadAddonRulesFile(c echo.Context, file io.Reader, fileName string, insuranceCode string) (int, error) {
	// Read file into memory
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}

	expectedHeaders := []string{"addon_code", "product_code", "insurance_code", "rules", "created_by"}

	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return 0, err
	}

	nextID := utils.GetNextAddonRuleID(draftData)

	// Try Excel first
	f, err := excelize.OpenReader(bytes.NewReader(fileBytes))
	if err == nil {
		defer f.Close()

		sheetName := f.GetSheetName(0)
		if sheetName == "" {
			return 0, &validationError{Message: "No sheet found in Excel file"}
		}

		rows, err := f.GetRows(sheetName)
		if err != nil {
			return 0, err
		}

		if len(rows) < 2 {
			return 0, &validationError{Message: "File must contain at least header and one data row"}
		}

		var newRules []utils.AddonRuleDraft
		for rowIdx := 1; rowIdx < len(rows); rowIdx++ {
			row := rows[rowIdx]
			if len(row) < len(expectedHeaders) {
				continue
			}

			rule := utils.AddonRuleDraft{
				ID:            nextID,
				AddonCode:     u.getCSVValue(row, 0, ""),
				ProductCode:   u.getCSVValue(row, 1, ""),
				InsuranceCode: u.getCSVValue(row, 2, ""),
				Rules:         u.getCSVValue(row, 3, "{}"),
				CreatedBy:     1,
				Timestamp:     utils.GetTimestamp(),
			}

			if val := u.getCSVValue(row, 4, "1"); val != "" {
				if createdBy, err := strconv.Atoi(val); err == nil {
					rule.CreatedBy = createdBy
				}
			}

			if rule.AddonCode == "" || rule.ProductCode == "" || rule.InsuranceCode == "" {
				continue
			}

			newRules = append(newRules, rule)
			nextID++
		}

		draftData.AddonRules = append(draftData.AddonRules, newRules...)
		if err := u.repository.SaveDraftData(c, draftData); err != nil {
			return 0, err
		}

		return len(newRules), nil
	}

	// Try CSV
	reader := csv.NewReader(bytes.NewReader(fileBytes))
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err != nil {
		return 0, &validationError{Message: "Failed to read file"}
	}

	normalizedHeaders := make([]string, len(headers))
	for i, h := range headers {
		normalizedHeaders[i] = strings.TrimSpace(strings.ToLower(h))
	}

	if len(normalizedHeaders) < len(expectedHeaders) {
		return 0, &validationError{Message: "Invalid CSV format: missing required columns"}
	}

	var newRules []utils.AddonRuleDraft
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		rule := utils.AddonRuleDraft{
			ID:            nextID,
			AddonCode:     u.getCSVValue(row, 0, ""),
			ProductCode:   u.getCSVValue(row, 1, ""),
			InsuranceCode: u.getCSVValue(row, 2, ""),
			Rules:         u.getCSVValue(row, 3, "{}"),
			CreatedBy:     1,
			Timestamp:     utils.GetTimestamp(),
		}

		if val := u.getCSVValue(row, 4, "1"); val != "" {
			if createdBy, err := strconv.Atoi(val); err == nil {
				rule.CreatedBy = createdBy
			}
		}

		if rule.AddonCode == "" || rule.ProductCode == "" || rule.InsuranceCode == "" {
			continue
		}

		newRules = append(newRules, rule)
		nextID++
	}

	draftData.AddonRules = append(draftData.AddonRules, newRules...)
	if err := u.repository.SaveDraftData(c, draftData); err != nil {
		return 0, err
	}

	return len(newRules), nil
}

// Helper functions
func (u *usecase) getCSVValue(row []string, index int, defaultValue string) string {
	if index < len(row) && strings.TrimSpace(row[index]) != "" {
		return strings.TrimSpace(row[index])
	}
	return defaultValue
}

// GetAddonRulesDraft gets all addon rule drafts
func (u *usecase) GetAddonRulesDraft(c echo.Context, insuranceCode string) ([]utils.AddonRuleDraft, error) {
	rules, err := u.repository.GetAllAddonRulesDraft(c, insuranceCode)
	if err != nil {
		return nil, err
	}

	// Sort by product_code, then addon_code
	sort.Slice(rules, func(i, j int) bool {
		if rules[i].ProductCode != rules[j].ProductCode {
			return rules[i].ProductCode < rules[j].ProductCode
		}
		return rules[i].AddonCode < rules[j].AddonCode
	})

	return rules, nil
}

// UpdateAddonRuleDraft updates a draft addon rule
func (u *usecase) UpdateAddonRuleDraft(c echo.Context, req addonrules.UpdateAddonRuleDraftRequest) error {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return err
	}

	if err := u.repository.UpdateAddonRuleDraftInMemory(c, draftData, req); err != nil {
		return err
	}

	return u.repository.SaveDraftData(c, draftData)
}

// DeleteAddonRuleDraft deletes a draft addon rule
func (u *usecase) DeleteAddonRuleDraft(c echo.Context, id int) error {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return err
	}

	if err := u.repository.DeleteAddonRuleDraftFromMemory(c, draftData, id); err != nil {
		return err
	}

	return u.repository.SaveDraftData(c, draftData)
}

// ClearAddonRulesDraft clears all draft addon rules
func (u *usecase) ClearAddonRulesDraft(c echo.Context) error {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return err
	}

	u.repository.ClearAddonRulesDraftInMemory(c, draftData)
	return u.repository.SaveDraftData(c, draftData)
}

// ConfirmAddonRules confirms addon rules from draft to database
func (u *usecase) ConfirmAddonRules(c echo.Context) (*addonrules.ConfirmAddonRulesResult, error) {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return nil, err
	}

	if len(draftData.AddonRules) == 0 {
		return nil, &validationError{Message: "No draft addon rules to confirm"}
	}

	insertedCount := 0
	for _, rule := range draftData.AddonRules {
		if err := u.repository.InsertAddonRule(c, rule.AddonCode, rule.ProductCode, rule.InsuranceCode, rule.Rules, rule.CreatedBy); err != nil {
			log.Printf("Error inserting addon rule: %v", err)
			return nil, err
		}
		insertedCount++
	}

	// Clear draft data
	draftData.AddonRules = []utils.AddonRuleDraft{}
	u.repository.SaveDraftData(c, draftData)

	return &addonrules.ConfirmAddonRulesResult{
		InsertedCount: insertedCount,
	}, nil
}

// GetAddonRules gets all confirmed addon rules
func (u *usecase) GetAddonRules(c echo.Context, insuranceCode string) ([]map[string]interface{}, error) {
	return u.repository.GetAllAddonRules(c, insuranceCode)
}

// UpdateAddonRule updates a confirmed addon rule
func (u *usecase) UpdateAddonRule(c echo.Context, id int, req addonrules.UpdateAddonRuleRequest, userName string) error {
	beforeData, errBefore := u.repository.GetAddonRuleBeforeUpdate(c, id)
	if errBefore != nil {
		log.Printf("Error getting data before update: %v", errBefore)
		beforeData = &addonrules.AddonRuleBeforeUpdate{}
	}

	rowsAffected, err := u.repository.UpdateAddonRule(c, id, req.UpdateData)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &notFoundError{Message: "Addon rule not found"}
	}

	// Compare and log changed fields
	changedFields := make(map[string]interface{})

	beforeProductCodeStr := ""
	if beforeData.ProductCode.Valid {
		beforeProductCodeStr = beforeData.ProductCode.String
	}
	if productCode, ok := req.UpdateData["product_code"].(string); ok && beforeProductCodeStr != productCode {
		changedFields["product_code"] = map[string]interface{}{
			"before": beforeProductCodeStr,
			"after":  productCode,
		}
	}

	beforeInsuranceCodeStr := ""
	if beforeData.InsuranceCode.Valid {
		beforeInsuranceCodeStr = beforeData.InsuranceCode.String
	}
	if insuranceCode, ok := req.UpdateData["insurance_code"].(string); ok && beforeInsuranceCodeStr != insuranceCode {
		changedFields["insurance_code"] = map[string]interface{}{
			"before": beforeInsuranceCodeStr,
			"after":  insuranceCode,
		}
	}

	beforeRulesStr := ""
	if beforeData.Rules.Valid {
		beforeRulesStr = beforeData.Rules.String
	}
	if rules, ok := req.UpdateData["rules"].(string); ok && beforeRulesStr != rules {
		changedFields["rules"] = map[string]interface{}{
			"before": beforeRulesStr,
			"after":  rules,
		}
	}

	if len(changedFields) > 0 {
		beforeChanged := make(map[string]interface{})
		afterChanged := make(map[string]interface{})

		insuranceCode := beforeInsuranceCodeStr
		if insuranceCode == "" {
			if ic, ok := req.UpdateData["insurance_code"].(string); ok {
				insuranceCode = ic
			}
		}
		beforeChanged["insurance_code"] = beforeInsuranceCodeStr
		if ic, ok := req.UpdateData["insurance_code"].(string); ok {
			afterChanged["insurance_code"] = ic
		} else {
			afterChanged["insurance_code"] = beforeInsuranceCodeStr
		}

		for field, changeData := range changedFields {
			changeMap := changeData.(map[string]interface{})
			beforeChanged[field] = changeMap["before"]
			afterChanged[field] = changeMap["after"]
		}

		beforeJSONBytes, _ := json.Marshal(beforeChanged)
		afterJSONBytes, _ := json.Marshal(afterChanged)

		u.repository.InsertHistory(c, userName, "Addon Rules", "UPDATE", id, "addon_rule",
			string(beforeJSONBytes), string(afterJSONBytes))
	}

	return nil
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

