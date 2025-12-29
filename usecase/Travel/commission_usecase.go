package usecase

import (
	"fmt"
	"multi-onboarding/model/Travel"
	"multi-onboarding/repository/Travel"
)

var (
	ErrCommissionNotFound = repository.ErrCommissionNotFound
)

// CommissionUsecase handles business logic for commissions
type CommissionUsecase struct {
	repo        repository.CommissionRepository
	productRepo repository.ProductRepository
}

// NewCommissionUsecase creates a new commission usecase
func NewCommissionUsecase(repo repository.CommissionRepository, productRepo repository.ProductRepository) *CommissionUsecase {
	return &CommissionUsecase{
		repo:        repo,
		productRepo: productRepo,
	}
}

// GetAllCommissions retrieves all commissions
func (u *CommissionUsecase) GetAllCommissions() ([]model.Commission, error) {
	return u.repo.GetAll()
}

// GetCommissionByID retrieves a commission by ID
func (u *CommissionUsecase) GetCommissionByID(id int64) (*model.Commission, error) {
	return u.repo.GetByID(id)
}

// GetCommissionsByProductCode retrieves commissions by product code
func (u *CommissionUsecase) GetCommissionsByProductCode(productCode string) ([]model.Commission, error) {
	return u.repo.GetByProductCode(productCode)
}

// CreateCommission creates a new commission
func (u *CommissionUsecase) CreateCommission(commission *model.Commission) error {
	if err := u.validateCommission(commission); err != nil {
		return err
	}
	return u.repo.Create(commission)
}

// UpdateCommission updates an existing commission
func (u *CommissionUsecase) UpdateCommission(id int64, commission *model.Commission) error {
	if err := u.validateCommission(commission); err != nil {
		return err
	}
	commission.ID = id
	return u.repo.Update(commission)
}

// DeleteCommission deletes a commission by ID
func (u *CommissionUsecase) DeleteCommission(id int64) error {
	return u.repo.Delete(id)
}

// GetAllDrafts retrieves all commission drafts
func (u *CommissionUsecase) GetAllDrafts() ([]model.CommissionDraft, error) {
	return u.repo.GetAllDrafts()
}

// CreateDraft creates a new commission draft
func (u *CommissionUsecase) CreateDraft(draft *model.CommissionDraft) error {
	if err := u.validateCommissionDraft(draft); err != nil {
		return err
	}
	return u.repo.CreateDraft(draft)
}

// UpdateDraft updates an existing commission draft
func (u *CommissionUsecase) UpdateDraft(id int, draft *model.CommissionDraft) error {
	if err := u.validateCommissionDraft(draft); err != nil {
		return err
	}
	draft.ID = id
	return u.repo.UpdateDraft(id, draft)
}

// DeleteDraft deletes a commission draft by ID
func (u *CommissionUsecase) DeleteDraft(id int) error {
	return u.repo.DeleteDraft(id)
}

// ConfirmDrafts confirms all drafts to the main commissions table
func (u *CommissionUsecase) ConfirmDrafts() error {
	return u.repo.ConfirmDrafts()
}

// GenerateCommissionsToMultipleTables generates commission data to 3 tables
func (u *CommissionUsecase) GenerateCommissionsToMultipleTables(request *GenerateCommissionRequest) error {
	if request.InsuranceCode == "" {
		return fmt.Errorf("insurance code is required")
	}
	if len(request.ProductCodes) == 0 {
		return fmt.Errorf("at least one product code is required")
	}
	if request.CommissionVATType != "EXCLUSIVE" && request.CommissionVATType != "INCLUSIVE" {
		return fmt.Errorf("commission VAT type must be EXCLUSIVE or INCLUSIVE")
	}
	if request.AFVATType != "EXCLUSIVE" && request.AFVATType != "INCLUSIVE" {
		return fmt.Errorf("AF VAT type must be EXCLUSIVE or INCLUSIVE")
	}

	// Validate that all product codes exist in travel_service_development.products table
	allProducts, err := u.productRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to validate product codes: %w", err)
	}

	// Create a map of valid product codes
	validProductCodes := make(map[string]bool)
	for _, product := range allProducts {
		validProductCodes[product.Code] = true
	}

	// Check each product code
	var invalidCodes []string
	for _, productCode := range request.ProductCodes {
		if !validProductCodes[productCode] {
			invalidCodes = append(invalidCodes, productCode)
		}
	}

	if len(invalidCodes) > 0 {
		return fmt.Errorf("invalid product codes (not found in travel_service_development.products): %v", invalidCodes)
	}

	return u.repo.GenerateCommissionsToMultipleTables(
		request.InsuranceCode,
		request.ProductCodes,
		request.CommissionPercentage,
		request.AFPercentage,
		request.AdminFee,
		request.CommissionVATType,
		request.AFVATType,
	)
}

// GetCommissionsFromTable retrieves data from commission_service_development.commissions
func (u *CommissionUsecase) GetCommissionsFromTable(insuranceCode string) ([]map[string]interface{}, error) {
	return u.repo.GetCommissionsFromTable(insuranceCode)
}

// GetPlanCommissionsFromTable retrieves data from quotation_service_development.plan_commissions
func (u *CommissionUsecase) GetPlanCommissionsFromTable(insuranceCode string) ([]map[string]interface{}, error) {
	return u.repo.GetPlanCommissionsFromTable(insuranceCode)
}

// GetDefaultConfigFromTable retrieves data from agent_service_development.default_config_products
func (u *CommissionUsecase) GetDefaultConfigFromTable(insuranceCode string) ([]map[string]interface{}, error) {
	return u.repo.GetDefaultConfigFromTable(insuranceCode)
}

// CheckDuplicateProductCodes checks if product codes already exist in the 3 tables
func (u *CommissionUsecase) CheckDuplicateProductCodes(productCodes []string) (map[string][]string, error) {
	return u.repo.CheckDuplicateProductCodes(productCodes)
}

// GenerateCommissionRequest represents the request to generate commissions
type GenerateCommissionRequest struct {
	InsuranceCode        string   `json:"insurance_code"`
	ProductCodes         []string `json:"product_codes"`
	CommissionPercentage float64  `json:"commission_percentage"`
	CommissionVATType    string   `json:"commission_vat_type"`
	AFPercentage         float64  `json:"af_percentage"`
	AFVATType            string   `json:"af_vat_type"`
	AdminFee             float64  `json:"admin_fee"`
}

// validateCommission validates commission data
func (u *CommissionUsecase) validateCommission(commission *model.Commission) error {
	if commission.ProductCode == "" {
		return fmt.Errorf("product code is required")
	}
	if commission.CommissionVATType == "" {
		return fmt.Errorf("commission VAT type is required")
	}
	if commission.CommissionVATType != "EXCLUSIVE" && commission.CommissionVATType != "INCLUSIVE" {
		return fmt.Errorf("commission VAT type must be EXCLUSIVE or INCLUSIVE")
	}
	if commission.AFVATType == "" {
		return fmt.Errorf("AF VAT type is required")
	}
	if commission.AFVATType != "EXCLUSIVE" && commission.AFVATType != "INCLUSIVE" {
		return fmt.Errorf("AF VAT type must be EXCLUSIVE or INCLUSIVE")
	}
	if commission.CommissionPercentage < 0 {
		return fmt.Errorf("commission percentage must be non-negative")
	}
	if commission.AFPercentage < 0 {
		return fmt.Errorf("AF percentage must be non-negative")
	}
	if commission.AdminFee < 0 {
		return fmt.Errorf("admin fee must be non-negative")
	}
	return nil
}

// validateCommissionDraft validates commission draft data
func (u *CommissionUsecase) validateCommissionDraft(draft *model.CommissionDraft) error {
	if draft.ProductCode == "" {
		return fmt.Errorf("product code is required")
	}
	if draft.CommissionVATType == "" {
		return fmt.Errorf("commission VAT type is required")
	}
	if draft.CommissionVATType != "EXCLUSIVE" && draft.CommissionVATType != "INCLUSIVE" {
		return fmt.Errorf("commission VAT type must be EXCLUSIVE or INCLUSIVE")
	}
	if draft.AFVATType == "" {
		return fmt.Errorf("AF VAT type is required")
	}
	if draft.AFVATType != "EXCLUSIVE" && draft.AFVATType != "INCLUSIVE" {
		return fmt.Errorf("AF VAT type must be EXCLUSIVE or INCLUSIVE")
	}
	if draft.CommissionPercentage < 0 {
		return fmt.Errorf("commission percentage must be non-negative")
	}
	if draft.AFPercentage < 0 {
		return fmt.Errorf("AF percentage must be non-negative")
	}
	if draft.AdminFee < 0 {
		return fmt.Errorf("admin fee must be non-negative")
	}
	return nil
}
