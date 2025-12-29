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
	"multi-onboarding/domain/vehicle/insuranceclauses"
	"multi-onboarding/utils"
)

type usecase struct {
	repository insuranceclauses.Repository
}

// NewInsuranceClausesUsecase creates a new insurance clauses usecase
func NewInsuranceClausesUsecase(repository insuranceclauses.Repository) insuranceclauses.Usecase {
	return &usecase{
		repository: repository,
	}
}

// DownloadClausesTemplate generates Excel template for clauses
func (u *usecase) DownloadClausesTemplate(c echo.Context, insuranceCode string) ([]byte, string, error) {
	products, err := u.repository.GetProductsByInsuranceCode(c, insuranceCode)
	if err != nil {
		return nil, "", err
	}

	if len(products) == 0 {
		return nil, "", &notFoundError{Message: "No active products found for this insurance"}
	}

	// Create Excel file
	f := excelize.NewFile()
	sheetName := "Insurance Clauses"
	f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")

	// Write header row
	headers := []string{
		"insurance_code", "product_type", "code", "title", "content", "description", "is_mandatory", "is_active", "created_by",
	}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	row := 2
	// Generate rows for each product
	for _, productCode := range products {
		// Generate rows for different product types
		productTypes := []string{"CAR", "MOTORCYCLE"}
		for _, productType := range productTypes {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), insuranceCode)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), productType)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "CLAUSE_"+productCode+"_"+productType)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), "Clause "+productType)
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), "")
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), "")
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), "0")
			f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), "1")
			f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), "1")
			row++
		}
	}

	// Write to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", err
	}

	fileName := fmt.Sprintf("%s_Insurance_Clauses_Template.xlsx", insuranceCode)
	return buf.Bytes(), fileName, nil
}

// UploadClausesFile handles file upload and parsing
func (u *usecase) UploadClausesFile(c echo.Context, file io.Reader, fileName string, insuranceCode string) (int, error) {
	// Read file into memory
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}

	expectedHeaders := []string{"insurance_code", "product_type", "code", "title", "content", "description", "is_mandatory", "is_active", "created_by"}
	expectedHeadersLower := make([]string, len(expectedHeaders))
	for i, h := range expectedHeaders {
		expectedHeadersLower[i] = strings.ToLower(h)
	}

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
		for col := 1; col <= 9; col++ {
			cell, err := f.GetCellValue(sheetName, fmt.Sprintf("%c1", 'A'+col-1))
			if err != nil {
				break
			}
			if cell == "" && col > 1 {
				break
			}
			headers = append(headers, strings.TrimSpace(strings.ToLower(cell)))
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

			nextID := utils.GetNextClausesID(draftData)
			clausesDraft := u.parseClauseRow(record, nextID, insuranceCode)
			draftData.Clauses = append(draftData.Clauses, clausesDraft)
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

		nextID := utils.GetNextClausesID(draftData)
		clausesDraft := u.parseClauseRow(record, nextID, insuranceCode)
		draftData.Clauses = append(draftData.Clauses, clausesDraft)
		rowCount++
	}

	if err := u.repository.SaveDraftData(c, draftData); err != nil {
		return 0, err
	}

	return rowCount, nil
}

// parseClauseRow parses a row from Excel/CSV into ClausesDraft
func (u *usecase) parseClauseRow(record []string, nextID int, insuranceCode string) utils.ClausesDraft {
	insuranceCodeVal := strings.TrimSpace(record[0])
	if insuranceCode != "" {
		insuranceCodeVal = insuranceCode
	}

	productType := strings.TrimSpace(record[1])
	code := strings.TrimSpace(record[2])
	title := strings.TrimSpace(record[3])
	content := strings.TrimSpace(record[4])
	description := ""
	isMandatory := 0
	isActive := 1

	if len(record) > 5 {
		description = strings.TrimSpace(record[5])
	}
	if len(record) > 6 {
		if parsed, err := strconv.Atoi(strings.TrimSpace(record[6])); err == nil {
			isMandatory = parsed
		}
	}
	if len(record) > 7 {
		if parsed, err := strconv.Atoi(strings.TrimSpace(record[7])); err == nil {
			isActive = parsed
		}
	}

	createdBy := 1
	if len(record) > 8 {
		if parsed, err := strconv.Atoi(strings.TrimSpace(record[8])); err == nil {
			createdBy = parsed
		}
	}

	return utils.ClausesDraft{
		ID:            nextID,
		InsuranceCode: insuranceCodeVal,
		ProductType:   productType,
		Code:          code,
		Title:         title,
		Content:       content,
		Description:   description,
		IsMandatory:   isMandatory,
		IsActive:      isActive,
		CreatedBy:     createdBy,
		Timestamp:     utils.GetTimestamp(),
	}
}

// GetClausesDraft gets all clause drafts
func (u *usecase) GetClausesDraft(c echo.Context, insuranceCode string) ([]utils.ClausesDraft, error) {
	return u.repository.GetAllClausesDraft(c, insuranceCode)
}

// UpdateClauseDraft updates a draft clause
func (u *usecase) UpdateClauseDraft(c echo.Context, req insuranceclauses.UpdateClauseDraftRequest) error {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return err
	}

	if err := u.repository.UpdateClauseDraftInMemory(c, draftData, req); err != nil {
		return err
	}

	return u.repository.SaveDraftData(c, draftData)
}

// DeleteClauseDraft deletes a draft clause
func (u *usecase) DeleteClauseDraft(c echo.Context, id int) error {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return err
	}

	if err := u.repository.DeleteClauseDraftFromMemory(c, draftData, id); err != nil {
		return err
	}

	return u.repository.SaveDraftData(c, draftData)
}

// ClearClausesDraft clears all draft clauses
func (u *usecase) ClearClausesDraft(c echo.Context) error {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return err
	}

	u.repository.ClearClausesDraftInMemory(c, draftData)
	return u.repository.SaveDraftData(c, draftData)
}

// ConfirmClauses confirms clauses from draft to database
func (u *usecase) ConfirmClauses(c echo.Context) (*insuranceclauses.ConfirmClausesResult, error) {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return nil, err
	}

	if len(draftData.Clauses) == 0 {
		return nil, &validationError{Message: "No draft data to confirm"}
	}

	insertedCount := 0
	for i, clause := range draftData.Clauses {
		if err := u.repository.InsertClause(c, clause); err != nil {
			log.Printf("Error inserting clause: %v", err)
			return nil, fmt.Errorf("failed to insert clause at index %d: %v", i, err)
		}
		insertedCount++
	}

	// Clear draft data
	draftData.Clauses = []utils.ClausesDraft{}
	u.repository.SaveDraftData(c, draftData)

	return &insuranceclauses.ConfirmClausesResult{
		InsertedCount: insertedCount,
	}, nil
}

// GetClauses gets all confirmed clauses
func (u *usecase) GetClauses(c echo.Context, insuranceCode string) ([]map[string]interface{}, error) {
	return u.repository.GetAllClauses(c, insuranceCode)
}

// UpdateClause updates a confirmed clause
func (u *usecase) UpdateClause(c echo.Context, id int, req insuranceclauses.UpdateClauseRequest, userName string) error {
	beforeData, errBefore := u.repository.GetClauseBeforeUpdate(c, id)
	if errBefore != nil {
		log.Printf("Error getting data before update: %v", errBefore)
		beforeData = &insuranceclauses.ClauseBeforeUpdate{}
	}

	rowsAffected, err := u.repository.UpdateClause(c, id, req.UpdateData)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &notFoundError{Message: "Clauses not found"}
	}

	// Compare and log changed fields
	changedFields := u.compareClauseFields(beforeData, req.UpdateData)

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
			beforeChanged["content"] = ""
			beforeChanged["description"] = ""
			beforeChanged["is_mandatory"] = 0
			beforeChanged["is_active"] = 0

			afterChanged["product_type"] = req.UpdateData["product_type"]
			afterChanged["code"] = req.UpdateData["code"]
			afterChanged["title"] = req.UpdateData["title"]
			afterChanged["content"] = req.UpdateData["content"]
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

		u.repository.InsertHistory(c, userName, "Insurance Clauses", "UPDATE", id, "insurance_clause",
			string(beforeJSONBytes), string(afterJSONBytes))
	}

	return nil
}

// compareClauseFields compares clause fields before and after update
func (u *usecase) compareClauseFields(before *insuranceclauses.ClauseBeforeUpdate, updateData map[string]interface{}) map[string]interface{} {
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
	if val, ok := updateData["content"]; ok {
		compareNullableString(before.Content, val, "content")
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


