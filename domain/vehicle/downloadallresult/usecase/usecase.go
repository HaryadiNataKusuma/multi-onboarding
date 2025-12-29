package usecase

import (
	"bytes"
	"fmt"
	"log"

	"github.com/labstack/echo"
	"github.com/xuri/excelize/v2"
	"multi-onboarding/domain/vehicle/downloadallresult"
)

type usecase struct {
	repository downloadallresult.Repository
}

// NewDownloadAllResultUsecase creates a new download all result usecase
func NewDownloadAllResultUsecase(repository downloadallresult.Repository) downloadallresult.Usecase {
	return &usecase{
		repository: repository,
	}
}

// DownloadAllResult generates an Excel file with multiple sheets containing all RESULT data
func (u *usecase) DownloadAllResult(c echo.Context, insuranceCode string) ([]byte, string, error) {
	// Create new Excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing Excel file: %v", err)
		}
	}()

	// Delete default Sheet1
	f.DeleteSheet("Sheet1")

	// 1. Insurances (active only)
	if err := u.addInsurancesSheet(c, f, insuranceCode); err != nil {
		log.Printf("Error adding insurances sheet: %v", err)
	}

	// 2. Products (active only)
	if err := u.addProductsSheet(c, f, insuranceCode); err != nil {
		log.Printf("Error adding products sheet: %v", err)
	}

	// 3. Templates
	if err := u.addTemplatesSheet(c, f, insuranceCode); err != nil {
		log.Printf("Error adding templates sheet: %v", err)
	}

	// 4. Commissions (3 tables: commissions, plan_commissions, default_config_products)
	if err := u.addCommissionsSheet(c, f, insuranceCode); err != nil {
		log.Printf("Error adding commissions sheet: %v", err)
	}
	if err := u.addPlanCommissionsSheet(c, f, insuranceCode); err != nil {
		log.Printf("Error adding plan commissions sheet: %v", err)
	}
	if err := u.addDefaultConfigProductsSheet(c, f, insuranceCode); err != nil {
		log.Printf("Error adding default config products sheet: %v", err)
	}

	// 5. Product Rules
	if err := u.addProductRulesSheet(c, f, insuranceCode); err != nil {
		log.Printf("Error adding product rules sheet: %v", err)
	}

	// 6. Addons (all addons, latest input by user - we'll get all active addons)
	if err := u.addAddonsSheet(c, f); err != nil {
		log.Printf("Error adding addons sheet: %v", err)
	}

	// 7. Addon Rules
	if err := u.addAddonRulesSheet(c, f, insuranceCode); err != nil {
		log.Printf("Error adding addon rules sheet: %v", err)
	}

	// 8. Insurance Own Risks (active only)
	if err := u.addInsuranceOwnRisksSheet(c, f, insuranceCode); err != nil {
		log.Printf("Error adding insurance own risks sheet: %v", err)
	}

	// 9. Insurance Clauses (active only)
	if err := u.addInsuranceClausesSheet(c, f, insuranceCode); err != nil {
		log.Printf("Error adding insurance clauses sheet: %v", err)
	}

	// Write Excel file to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", fmt.Errorf("failed to write Excel file: %w", err)
	}

	fileName := fmt.Sprintf("%s_All_RESULT.xlsx", insuranceCode)
	return buf.Bytes(), fileName, nil
}

// Helper function to add Insurances sheet
func (u *usecase) addInsurancesSheet(c echo.Context, f *excelize.File, insuranceCode string) error {
	rows, err := u.repository.GetInsurances(c, insuranceCode)
	if err != nil {
		return err
	}

	sheetName := "Insurances"
	f.NewSheet(sheetName)

	// Write headers
	headers := []string{"ID", "Code", "Name", "Status", "Logo", "Created At", "Updated At"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data
	for i, row := range rows {
		excelRow := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", excelRow), row.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", excelRow), row.Code)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", excelRow), row.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", excelRow), row.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", excelRow), row.Logo)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", excelRow), row.CreatedAt)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", excelRow), row.UpdatedAt)
	}

	return nil
}

// Helper function to add Products sheet
func (u *usecase) addProductsSheet(c echo.Context, f *excelize.File, insuranceCode string) error {
	rows, err := u.repository.GetProducts(c, insuranceCode)
	if err != nil {
		return err
	}

	sheetName := "Products"
	f.NewSheet(sheetName)

	// Write headers
	headers := []string{"ID", "Code", "Insurance Code", "Insurance Product Code", "Vehicle Type",
		"Is Insurance API", "Name", "Logo", "Is Active", "Is Insurer Driven",
		"Is Personal Usage", "Is Electric Vehicle", "Summary", "Config",
		"Insurance Detail", "Protection Detail", "How To Claim", "Created By", "Base Product",
		"Created At", "Updated At"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data
	for i, row := range rows {
		excelRow := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", excelRow), row.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", excelRow), row.Code)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", excelRow), row.InsuranceCode)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", excelRow), row.InsuranceProductCode)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", excelRow), row.VehicleType)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", excelRow), row.IsInsuranceAPI)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", excelRow), row.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", excelRow), row.Logo)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", excelRow), row.IsActive)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", excelRow), row.IsInsurerDriven)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", excelRow), row.IsPersonalUsage)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", excelRow), row.IsElectricVehicle)
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", excelRow), row.Summary)
		f.SetCellValue(sheetName, fmt.Sprintf("N%d", excelRow), row.Config)
		f.SetCellValue(sheetName, fmt.Sprintf("O%d", excelRow), row.InsuranceDetail)
		f.SetCellValue(sheetName, fmt.Sprintf("P%d", excelRow), row.ProtectionDetail)
		f.SetCellValue(sheetName, fmt.Sprintf("Q%d", excelRow), row.HowToClaim)
		f.SetCellValue(sheetName, fmt.Sprintf("R%d", excelRow), row.CreatedBy)
		f.SetCellValue(sheetName, fmt.Sprintf("S%d", excelRow), row.BaseProduct)
		f.SetCellValue(sheetName, fmt.Sprintf("T%d", excelRow), row.CreatedAt)
		f.SetCellValue(sheetName, fmt.Sprintf("U%d", excelRow), row.UpdatedAt)
	}

	return nil
}

// Helper function to add Templates sheet
func (u *usecase) addTemplatesSheet(c echo.Context, f *excelize.File, insuranceCode string) error {
	rows, err := u.repository.GetTemplates(c, insuranceCode)
	if err != nil {
		return err
	}

	sheetName := "Templates"
	f.NewSheet(sheetName)

	// Write headers
	headers := []string{"ID", "Locale", "Template ID", "Value", "Insurance Code", "Created By", "Created At"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data
	for i, row := range rows {
		excelRow := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", excelRow), row.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", excelRow), row.Locale)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", excelRow), row.TemplateID)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", excelRow), row.Value)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", excelRow), row.InsuranceCode)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", excelRow), row.CreatedBy)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", excelRow), row.CreatedAt)
	}

	return nil
}

// Helper function to add Commissions sheet
func (u *usecase) addCommissionsSheet(c echo.Context, f *excelize.File, insuranceCode string) error {
	rows, err := u.repository.GetCommissions(c, insuranceCode)
	if err != nil {
		return err
	}

	sheetName := "Commissions"
	f.NewSheet(sheetName)

	// Write headers
	headers := []string{"ID", "Product Code", "Insurance Code", "Agent Level", "Corporate ID",
		"Basic Commission", "Bonus Point", "Note", "Start Date", "End Date", "Created At", "Updated At"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data
	for i, row := range rows {
		excelRow := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", excelRow), row.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", excelRow), row.ProductCode)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", excelRow), row.InsuranceCode)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", excelRow), row.AgentLevel)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", excelRow), row.CorporateID)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", excelRow), row.BasicCommission)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", excelRow), row.BonusPoint)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", excelRow), row.Note)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", excelRow), row.StartDate)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", excelRow), row.EndDate)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", excelRow), row.CreatedAt)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", excelRow), row.UpdatedAt)
	}

	return nil
}

// Helper function to add Plan Commissions sheet
func (u *usecase) addPlanCommissionsSheet(c echo.Context, f *excelize.File, insuranceCode string) error {
	rows, err := u.repository.GetPlanCommissions(c, insuranceCode)
	if err != nil {
		return err
	}

	sheetName := "Plan Commissions"
	f.NewSheet(sheetName)

	// Write headers
	headers := []string{"ID", "Plan Code", "Product ID", "Insurer ID", "Commission Percentage", "Commission VAT Type",
		"AF Percentage", "AF VAT Type", "Admin Fee", "Hardcopy Fee", "Is Active", "Version", "Created At"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data
	for i, row := range rows {
		excelRow := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", excelRow), row.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", excelRow), row.PlanCode)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", excelRow), row.ProductID)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", excelRow), row.InsurerID)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", excelRow), row.CommissionPercentage)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", excelRow), row.CommissionVatType)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", excelRow), row.AFPercentage)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", excelRow), row.AFVatType)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", excelRow), row.AdminFee)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", excelRow), row.HardcopyFee)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", excelRow), row.IsActive)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", excelRow), row.Version)
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", excelRow), row.CreatedAt)
	}

	return nil
}

// Helper function to add Default Config Products sheet
func (u *usecase) addDefaultConfigProductsSheet(c echo.Context, f *excelize.File, insuranceCode string) error {
	rows, err := u.repository.GetDefaultConfigProducts(c, insuranceCode)
	if err != nil {
		return err
	}

	sheetName := "Default Config Prods"
	f.NewSheet(sheetName)

	// Write headers
	headers := []string{"ID", "Level", "Sequence", "Kind", "Category", "Insurer", "Product", "Selected",
		"Commission", "Point", "Upline Bonus", "Bonus", "Max Discount", "Order", "Created At", "Updated At"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data
	for i, row := range rows {
		excelRow := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", excelRow), row.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", excelRow), row.Level)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", excelRow), row.Sequence)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", excelRow), row.Kind)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", excelRow), row.Category)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", excelRow), row.Insurer)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", excelRow), row.Product)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", excelRow), row.Selected)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", excelRow), row.Commission)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", excelRow), row.Point)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", excelRow), row.UplineBonus)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", excelRow), row.Bonus)
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", excelRow), row.MaxDiscount)
		f.SetCellValue(sheetName, fmt.Sprintf("N%d", excelRow), row.Order)
		f.SetCellValue(sheetName, fmt.Sprintf("O%d", excelRow), row.CreatedAt)
		f.SetCellValue(sheetName, fmt.Sprintf("P%d", excelRow), row.UpdatedAt)
	}

	return nil
}

// Helper function to add Product Rules sheet
func (u *usecase) addProductRulesSheet(c echo.Context, f *excelize.File, insuranceCode string) error {
	rows, err := u.repository.GetProductRules(c, insuranceCode)
	if err != nil {
		return err
	}

	sheetName := "Product Rules"
	f.NewSheet(sheetName)

	// Write headers - match upload template format
	headers := []string{"product_code", "insurance_code", "vehicle_type", "vehicle_category",
		"start_vehicle_value", "end_vehicle_value", "limit_vehicle_age", "admin_fee", "hardcopy_admin_fee",
		"region_id", "base_premium_value", "loading_fee_premium_value", "commercial_usage_value",
		"start_loading_age", "additional_premium", "type_additional_premium", "rules", "created_by",
		"created_at", "base_premium_type", "base_loading_premium_value"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data
	for i, row := range rows {
		excelRow := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", excelRow), row.ProductCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", excelRow), row.InsuranceCode)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", excelRow), row.VehicleType)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", excelRow), row.VehicleCategory)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", excelRow), row.StartVehicleValue)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", excelRow), row.EndVehicleValue)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", excelRow), row.LimitVehicleAge)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", excelRow), row.AdminFee)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", excelRow), row.HardcopyAdminFee)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", excelRow), row.RegionID)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", excelRow), row.BasePremiumValue)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", excelRow), row.LoadingFeePremiumValue)
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", excelRow), row.CommercialUsageValue)
		f.SetCellValue(sheetName, fmt.Sprintf("N%d", excelRow), row.StartLoadingAge)
		f.SetCellValue(sheetName, fmt.Sprintf("O%d", excelRow), row.AdditionalPremium)
		f.SetCellValue(sheetName, fmt.Sprintf("P%d", excelRow), row.TypeAdditionalPremium)
		f.SetCellValue(sheetName, fmt.Sprintf("Q%d", excelRow), row.Rules)
		f.SetCellValue(sheetName, fmt.Sprintf("R%d", excelRow), row.CreatedBy)
		f.SetCellValue(sheetName, fmt.Sprintf("S%d", excelRow), row.CreatedAt)
		f.SetCellValue(sheetName, fmt.Sprintf("T%d", excelRow), row.BasePremiumType)
		f.SetCellValue(sheetName, fmt.Sprintf("U%d", excelRow), row.BaseLoadingPremiumValue)
	}

	return nil
}

// Helper function to add Addons sheet
func (u *usecase) addAddonsSheet(c echo.Context, f *excelize.File) error {
	rows, err := u.repository.GetAddons(c)
	if err != nil {
		return err
	}

	sheetName := "Addons"
	f.NewSheet(sheetName)

	// Write headers
	headers := []string{"ID", "Code", "Name", "Is Active", "Created At"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data
	for i, row := range rows {
		excelRow := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", excelRow), row.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", excelRow), row.Code)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", excelRow), row.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", excelRow), row.IsActive)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", excelRow), row.CreatedAt)
	}

	return nil
}

// Helper function to add Addon Rules sheet
func (u *usecase) addAddonRulesSheet(c echo.Context, f *excelize.File, insuranceCode string) error {
	rows, err := u.repository.GetAddonRules(c, insuranceCode)
	if err != nil {
		return err
	}

	sheetName := "Addon Rules"
	f.NewSheet(sheetName)

	// Write headers - match upload template format
	headers := []string{"addon_code", "product_code", "insurance_code", "rules", "created_by"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data
	for i, row := range rows {
		excelRow := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", excelRow), row.AddonCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", excelRow), row.ProductCode)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", excelRow), row.InsuranceCode)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", excelRow), row.Rules)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", excelRow), row.CreatedBy)
	}

	return nil
}

// Helper function to add Insurance Own Risks sheet
func (u *usecase) addInsuranceOwnRisksSheet(c echo.Context, f *excelize.File, insuranceCode string) error {
	rows, err := u.repository.GetInsuranceOwnRisks(c, insuranceCode)
	if err != nil {
		return err
	}

	sheetName := "Insurance Own Risks"
	f.NewSheet(sheetName)

	// Write headers - match upload template format
	headers := []string{"insurance_code", "product_type", "code", "title", "value", "value_type",
		"description", "is_mandatory", "is_active", "is_electric_vehicle", "created_by"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data
	for i, row := range rows {
		excelRow := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", excelRow), row.InsuranceCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", excelRow), row.ProductType)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", excelRow), row.Code)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", excelRow), row.Title)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", excelRow), row.Value)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", excelRow), row.ValueType)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", excelRow), row.Description)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", excelRow), row.IsMandatory)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", excelRow), row.IsActive)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", excelRow), row.IsElectricVehicle)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", excelRow), row.CreatedBy)
	}

	return nil
}

// Helper function to add Insurance Clauses sheet
func (u *usecase) addInsuranceClausesSheet(c echo.Context, f *excelize.File, insuranceCode string) error {
	rows, err := u.repository.GetInsuranceClauses(c, insuranceCode)
	if err != nil {
		return err
	}

	sheetName := "Insurance Clauses"
	f.NewSheet(sheetName)

	// Write headers - match upload template format
	headers := []string{"insurance_code", "product_type", "code", "title", "content",
		"description", "is_mandatory", "is_active", "created_by"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data
	for i, row := range rows {
		excelRow := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", excelRow), row.InsuranceCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", excelRow), row.ProductType)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", excelRow), row.Code)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", excelRow), row.Title)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", excelRow), row.Content)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", excelRow), row.Description)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", excelRow), row.IsMandatory)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", excelRow), row.IsActive)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", excelRow), row.CreatedBy)
	}

	return nil
}


