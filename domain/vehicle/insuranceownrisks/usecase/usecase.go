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
	"multi-onboarding/domain/vehicle/insuranceownrisks"
	"multi-onboarding/utils"
)

type usecase struct {
	repository insuranceownrisks.Repository
}

// NewInsuranceOwnRisksUsecase creates a new insurance own risks usecase
func NewInsuranceOwnRisksUsecase(repository insuranceownrisks.Repository) insuranceownrisks.Usecase {
	return &usecase{
		repository: repository,
	}
}

// DownloadOwnRisksTemplate generates Excel template for own risks
func (u *usecase) DownloadOwnRisksTemplate(c echo.Context, insuranceCode string) ([]byte, string, error) {
	products, err := u.repository.GetProductsByInsuranceCode(c, insuranceCode)
	if err != nil {
		return nil, "", err
	}

	if len(products) == 0 {
		return nil, "", &notFoundError{Message: "No active products found for this insurance"}
	}

	// Create Excel file
	f := excelize.NewFile()
	sheetName := "Insurance Own Risks"
	f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")

	// Write header row
	headers := []string{
		"insurance_code", "product_type", "code", "title", "value", "value_type",
		"description", "is_mandatory", "is_active", "is_electric_vehicle", "created_by",
	}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Generate exactly 10 rows with product_type: COMPREHENSIVE, TOTAL_LOSS, ALL
	product := products[0]
	row := 2

	// Generate 4 rows for COMPREHENSIVE
	for i := 0; i < 4; i++ {
		codeValue := fmt.Sprintf("OWN_RISK_%s_COMPREHENSIVE_%d", product.Code, i+1)
		titleValue := fmt.Sprintf("Own Risk COMPREHENSIVE %d", i+1)

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), insuranceCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "COMPREHENSIVE")
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), codeValue)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), titleValue)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), "0")
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), "PERCENTAGE")
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), "")
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), "0")
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), "1")
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), product.IsElectricVehicle)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), "1")
		row++
	}

	// Generate 3 rows for TOTAL_LOSS
	for i := 0; i < 3; i++ {
		codeValue := fmt.Sprintf("OWN_RISK_%s_TOTAL_LOSS_%d", product.Code, i+1)
		titleValue := fmt.Sprintf("Own Risk TOTAL_LOSS %d", i+1)

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), insuranceCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "TOTAL_LOSS")
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), codeValue)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), titleValue)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), "0")
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), "PERCENTAGE")
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), "")
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), "0")
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), "1")
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), product.IsElectricVehicle)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), "1")
		row++
	}

	// Generate 3 rows for ALL
	for i := 0; i < 3; i++ {
		codeValue := fmt.Sprintf("OWN_RISK_%s_ALL_%d", product.Code, i+1)
		titleValue := fmt.Sprintf("Own Risk ALL %d", i+1)

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), insuranceCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "ALL")
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), codeValue)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), titleValue)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), "0")
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), "PERCENTAGE")
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), "")
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), "0")
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), "1")
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), product.IsElectricVehicle)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), "1")
		row++
	}

	// Write to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", err
	}

	fileName := fmt.Sprintf("%s_Insurance_Own_Risks_Template.xlsx", insuranceCode)
	return buf.Bytes(), fileName, nil
}

// UploadOwnRisksFile handles file upload and parsing
func (u *usecase) UploadOwnRisksFile(c echo.Context, file io.Reader, fileName string, insuranceCode string) (int, error) {
	// Read file into memory
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}

	expectedHeaders := []string{"insurance_code", "product_type", "code", "title", "value", "value_type", "description", "is_mandatory", "is_active", "is_electric_vehicle", "created_by"}

	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return 0, err
	}

	// Try Excel first
	f, err := excelize.OpenReader(bytes.NewReader(fileBytes))
	if err == nil {
		defer f.Close()

		sheetName := f.GetSheetName(0)
		if sheetName == "" {
			return 0, &validationError{Message: "No sheet found in Excel file"}
		}

		// Read header row
		headers := make([]string, 0)
		for col := 1; col <= 11; col++ {
			cell, err := f.GetCellValue(sheetName, fmt.Sprintf("%c1", 'A'+col-1))
			if err != nil {
				break
			}
			if cell == "" && col > 1 {
				break
			}
			headers = append(headers, strings.TrimSpace(strings.ToLower(cell)))
		}

		expectedHeadersLower := make([]string, len(expectedHeaders))
		for i, h := range expectedHeaders {
			expectedHeadersLower[i] = strings.ToLower(h)
		}

		// Validate headers
		if len(headers) != len(expectedHeadersLower) {
			return 0, &validationError{Message: fmt.Sprintf("Invalid headers. Expected %d columns, got %d", len(expectedHeaders), len(headers))}
		}

		for i, expected := range expectedHeadersLower {
			if headers[i] != expected {
				return 0, &validationError{Message: fmt.Sprintf("Invalid header at position %d: expected '%s', got '%s'", i+1, expectedHeaders[i], headers[i])}
			}
		}

		// Read and parse data rows
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return 0, err
		}

		rowCount := 0
		for rowIdx := 1; rowIdx < len(rows); rowIdx++ {
			record := rows[rowIdx]
			if len(record) < len(expectedHeaders) {
				continue
			}

			nextID := utils.GetNextOwnRisksID(draftData)
			ownRisksDraft := u.parseOwnRiskRow(record, nextID, insuranceCode)
			draftData.OwnRisks = append(draftData.OwnRisks, ownRisksDraft)
			rowCount++
		}

		if err := u.repository.SaveDraftData(c, draftData); err != nil {
			return 0, err
		}

		return rowCount, nil
	}

	// Try CSV
	reader := csv.NewReader(bytes.NewReader(fileBytes))
	reader.Comma = ','
	reader.LazyQuotes = true

	// Read header
	headers, err := reader.Read()
	if err != nil {
		return 0, &validationError{Message: "Failed to read file"}
	}

	normalizedHeaders := make([]string, len(headers))
	for i, h := range headers {
		normalizedHeaders[i] = strings.TrimSpace(strings.ToLower(h))
	}

	expectedHeadersLower := make([]string, len(expectedHeaders))
	for i, h := range expectedHeaders {
		expectedHeadersLower[i] = strings.ToLower(h)
	}

	// Validate headers
	if len(normalizedHeaders) != len(expectedHeadersLower) {
		return 0, &validationError{Message: fmt.Sprintf("Invalid headers. Expected %d columns, got %d", len(expectedHeaders), len(normalizedHeaders))}
	}

	for i, expected := range expectedHeadersLower {
		if normalizedHeaders[i] != expected {
			return 0, &validationError{Message: fmt.Sprintf("Invalid header at position %d: expected '%s', got '%s'", i+1, expectedHeaders[i], normalizedHeaders[i])}
		}
	}

	rowCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		if len(record) < len(expectedHeaders) {
			continue
		}

		nextID := utils.GetNextOwnRisksID(draftData)
		ownRisksDraft := u.parseOwnRiskRow(record, nextID, insuranceCode)
		draftData.OwnRisks = append(draftData.OwnRisks, ownRisksDraft)
		rowCount++
	}

	if err := u.repository.SaveDraftData(c, draftData); err != nil {
		return 0, err
	}

	return rowCount, nil
}

// parseOwnRiskRow parses a row from Excel/CSV into OwnRisksDraft
func (u *usecase) parseOwnRiskRow(record []string, nextID int, insuranceCode string) utils.OwnRisksDraft {
	insuranceCodeVal := strings.TrimSpace(record[0])
	if insuranceCode != "" {
		insuranceCodeVal = insuranceCode
	}

	productType := strings.TrimSpace(record[1])
	code := strings.TrimSpace(record[2])
	title := strings.TrimSpace(record[3])
	valueStr := strings.TrimSpace(record[4])
	valueType := ""
	description := ""
	isMandatory := 0
	isActive := 1

	if len(record) > 5 {
		valueType = strings.TrimSpace(record[5])
	}
	if len(record) > 6 {
		description = strings.TrimSpace(record[6])
	}
	if len(record) > 7 {
		if parsed, err := strconv.Atoi(strings.TrimSpace(record[7])); err == nil {
			isMandatory = parsed
		}
	}
	if len(record) > 8 {
		if parsed, err := strconv.Atoi(strings.TrimSpace(record[8])); err == nil {
			isActive = parsed
		}
	}

	isElectricVehicle := 0
	if len(record) > 9 {
		if parsed, err := strconv.Atoi(strings.TrimSpace(record[9])); err == nil {
			isElectricVehicle = parsed
		}
	}

	createdBy := 1
	if len(record) > 10 {
		if parsed, err := strconv.Atoi(strings.TrimSpace(record[10])); err == nil {
			createdBy = parsed
		}
	}

	return utils.OwnRisksDraft{
		ID:                nextID,
		InsuranceCode:     insuranceCodeVal,
		ProductType:       productType,
		Code:              code,
		Title:             title,
		Value:             valueStr,
		ValueType:         valueType,
		Description:       description,
		IsMandatory:       isMandatory,
		IsActive:          isActive,
		IsElectricVehicle: isElectricVehicle,
		CreatedBy:         createdBy,
		Timestamp:         utils.GetTimestamp(),
	}
}

// GetOwnRisksDraft gets all own risk drafts
func (u *usecase) GetOwnRisksDraft(c echo.Context, insuranceCode string) ([]utils.OwnRisksDraft, error) {
	return u.repository.GetAllOwnRisksDraft(c, insuranceCode)
}

// UpdateOwnRiskDraft updates a draft own risk
func (u *usecase) UpdateOwnRiskDraft(c echo.Context, req insuranceownrisks.UpdateOwnRiskDraftRequest) error {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return err
	}

	if err := u.repository.UpdateOwnRiskDraftInMemory(c, draftData, req); err != nil {
		return err
	}

	return u.repository.SaveDraftData(c, draftData)
}

// DeleteOwnRiskDraft deletes a draft own risk
func (u *usecase) DeleteOwnRiskDraft(c echo.Context, id int) error {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return err
	}

	if err := u.repository.DeleteOwnRiskDraftFromMemory(c, draftData, id); err != nil {
		return err
	}

	return u.repository.SaveDraftData(c, draftData)
}

// ClearOwnRisksDraft clears all draft own risks
func (u *usecase) ClearOwnRisksDraft(c echo.Context) error {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return err
	}

	u.repository.ClearOwnRisksDraftInMemory(c, draftData)
	return u.repository.SaveDraftData(c, draftData)
}

// ConfirmOwnRisks confirms own risks from draft to database
func (u *usecase) ConfirmOwnRisks(c echo.Context) (*insuranceownrisks.ConfirmOwnRisksResult, error) {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return nil, err
	}

	if len(draftData.OwnRisks) == 0 {
		return nil, &validationError{Message: "No draft data to confirm"}
	}

	insertedCount := 0
	for i, risk := range draftData.OwnRisks {
		if err := u.repository.InsertOwnRisk(c, risk); err != nil {
			log.Printf("Error inserting own risk: %v", err)
			return nil, fmt.Errorf("failed to insert own risk at index %d: %v", i, err)
		}
		insertedCount++
	}

	// Clear draft data
	draftData.OwnRisks = []utils.OwnRisksDraft{}
	u.repository.SaveDraftData(c, draftData)

	return &insuranceownrisks.ConfirmOwnRisksResult{
		InsertedCount: insertedCount,
	}, nil
}

// GetOwnRisks gets all confirmed own risks
func (u *usecase) GetOwnRisks(c echo.Context, insuranceCode string) ([]map[string]interface{}, error) {
	return u.repository.GetAllOwnRisks(c, insuranceCode)
}

// UpdateOwnRisk updates a confirmed own risk
func (u *usecase) UpdateOwnRisk(c echo.Context, id int, req insuranceownrisks.UpdateOwnRiskRequest, userName string) error {
	beforeData, errBefore := u.repository.GetOwnRiskBeforeUpdate(c, id)
	if errBefore != nil {
		log.Printf("Error getting data before update: %v", errBefore)
		beforeData = &insuranceownrisks.OwnRiskBeforeUpdate{}
	}

	rowsAffected, err := u.repository.UpdateOwnRisk(c, id, req.UpdateData)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &notFoundError{Message: "Own risks not found"}
	}

	// Compare and log changed fields
	changedFields := u.compareOwnRiskFields(beforeData, req.UpdateData)

	if len(changedFields) > 0 || errBefore != nil {
		beforeChanged := make(map[string]interface{})
		afterChanged := make(map[string]interface{})

		insuranceCode := ""
		if beforeData.InsuranceCode.Valid {
			insuranceCode = beforeData.InsuranceCode.String
		}
		beforeChanged["insurance_code"] = insuranceCode
		afterChanged["insurance_code"] = insuranceCode

		if errBefore != nil {
			beforeChanged["product_type"] = ""
			beforeChanged["code"] = ""
			beforeChanged["title"] = ""
			beforeChanged["value"] = ""
			beforeChanged["value_type"] = ""
			beforeChanged["description"] = ""
			beforeChanged["is_mandatory"] = 0
			beforeChanged["is_active"] = 0

			afterChanged["product_type"] = req.UpdateData["product_type"]
			afterChanged["code"] = req.UpdateData["code"]
			afterChanged["title"] = req.UpdateData["title"]
			afterChanged["value"] = req.UpdateData["value"]
			afterChanged["value_type"] = req.UpdateData["value_type"]
			afterChanged["description"] = req.UpdateData["description"]
			afterChanged["is_mandatory"] = req.UpdateData["is_mandatory"]
			afterChanged["is_active"] = req.UpdateData["is_active"]
		} else {
			for field, changeData := range changedFields {
				changeMap := changeData.(map[string]interface{})
				beforeChanged[field] = changeMap["before"]
				afterChanged[field] = changeMap["after"]
			}
		}

		beforeJSONBytes, _ := json.Marshal(beforeChanged)
		afterJSONBytes, _ := json.Marshal(afterChanged)

		u.repository.InsertHistory(c, userName, "Insurance Own Risks", "UPDATE", id, "insurance_own_risk",
			string(beforeJSONBytes), string(afterJSONBytes))
	}

	return nil
}

// compareOwnRiskFields compares own risk fields before and after update
func (u *usecase) compareOwnRiskFields(before *insuranceownrisks.OwnRiskBeforeUpdate, updateData map[string]interface{}) map[string]interface{} {
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

	compareNullableFloat64 := func(before sql.NullFloat64, after interface{}, fieldName string) {
		afterFloat := float64(0)
		if after != nil {
			if f, ok := after.(float64); ok {
				afterFloat = f
			} else if i, ok := after.(int); ok {
				afterFloat = float64(i)
			} else if i, ok := after.(int64); ok {
				afterFloat = float64(i)
			} else if str, ok := after.(string); ok {
				if parsed, err := strconv.ParseFloat(str, 64); err == nil {
					afterFloat = parsed
				}
			}
		}
		beforeFloat := float64(0)
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

	compareNullableBool := func(before sql.NullBool, after interface{}, fieldName string) {
		afterBool := false
		if after != nil {
			if b, ok := after.(bool); ok {
				afterBool = b
			} else if i, ok := after.(int); ok {
				afterBool = i != 0
			} else if i, ok := after.(int64); ok {
				afterBool = i != 0
			} else if f, ok := after.(float64); ok {
				afterBool = f != 0
			}
		}
		beforeBool := false
		if before.Valid {
			beforeBool = before.Bool
		}
		if beforeBool != afterBool {
			changedFields[fieldName] = map[string]interface{}{
				"before": beforeBool,
				"after":  afterBool,
			}
		}
	}

	// Compare each field
	if val, ok := updateData["product_type"]; ok {
		compareNullableString(before.ProductType, val, "product_type")
	}
	if val, ok := updateData["code"]; ok {
		compareNullableString(before.Code, val, "code")
	}
	if val, ok := updateData["title"]; ok {
		compareNullableString(before.Title, val, "title")
	}
	if val, ok := updateData["value"]; ok {
		compareNullableFloat64(before.Value, val, "value")
	}
	if val, ok := updateData["value_type"]; ok {
		compareNullableString(before.ValueType, val, "value_type")
	}
	if val, ok := updateData["description"]; ok {
		compareNullableString(before.Description, val, "description")
	}
	if val, ok := updateData["is_mandatory"]; ok {
		compareNullableBool(before.IsMandatory, val, "is_mandatory")
	}
	if val, ok := updateData["is_active"]; ok {
		compareNullableBool(before.IsActive, val, "is_active")
	}

	return changedFields
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


