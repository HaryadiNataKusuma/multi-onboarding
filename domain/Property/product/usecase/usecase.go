package usecase

import (
	"errors"
	"fmt"

	"multi-onboarding/domain/Property/product"
)

type usecase struct {
	repository product.Repository
}

// NewProductUsecase creates a new product usecase
func NewProductUsecase(repository product.Repository) product.Usecase {
	return &usecase{
		repository: repository,
	}
}

// Draft operations

func (u *usecase) AddDraft(draft *product.ProductDraft) error {
	return u.repository.AddDraft(draft)
}

func (u *usecase) GetAllDrafts() ([]product.ProductDraft, error) {
	return u.repository.GetAllDrafts()
}

func (u *usecase) GetDraftByID(id int) (*product.ProductDraft, error) {
	return u.repository.GetDraftByID(id)
}

func (u *usecase) UpdateDraft(id int, draft *product.ProductDraft) error {
	return u.repository.UpdateDraft(id, draft)
}

func (u *usecase) DeleteDraft(id int) error {
	return u.repository.DeleteDraft(id)
}

func (u *usecase) ClearDrafts() error {
	return u.repository.ClearDrafts()
}

func (u *usecase) ConfirmDrafts(createdBy int64) error {
	return u.repository.ConfirmDrafts(createdBy)
}

// Database operations

func (u *usecase) GetAll(insuranceCode string) ([]product.Product, error) {
	return u.repository.GetAll(insuranceCode)
}

func (u *usecase) GetByID(id int) (*product.Product, error) {
	if id <= 0 {
		return nil, errors.New("invalid product ID")
	}
	return u.repository.GetByID(id)
}

func (u *usecase) Create(req product.CreateProductRequest) error {
	if err := u.validateCreateRequest(req); err != nil {
		return err
	}

	p := &product.Product{
		Code:             req.Code,
		InsuranceCode:    req.InsuranceCode,
		Name:             req.Name,
		Logo:             req.Logo,
		Type:             req.Type,
		IsActive:         req.IsActive,
		InsuranceType:    req.InsuranceType,
		AdminFee:         req.AdminFee,
		PropertyCategory: req.PropertyCategory,
		CoverageType:     req.CoverageType,
		MinCoverage:      req.MinCoverage,
		MaxCoverage:      req.MaxCoverage,
		Summary:          req.Summary,
		InsuranceDetail:  req.InsuranceDetail,
		ProtectionDetail: req.ProtectionDetail,
		HowToClaim:       req.HowToClaim,
		CreatedBy:        req.CreatedBy,
	}

	if p.CreatedBy == 0 {
		p.CreatedBy = 1
	}

	return u.repository.Create(p)
}

func (u *usecase) Update(id int, req product.UpdateProductRequest) error {
	if id <= 0 {
		return errors.New("invalid product ID")
	}
	if err := u.validateUpdateRequest(req); err != nil {
		return err
	}

	existing, err := u.repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("product not found")
	}

	existing.Code = req.Code
	existing.InsuranceCode = req.InsuranceCode
	existing.Name = req.Name
	existing.Type = req.Type
	existing.IsActive = req.IsActive
	existing.InsuranceType = req.InsuranceType
	existing.AdminFee = req.AdminFee
	existing.PropertyCategory = req.PropertyCategory
	existing.CoverageType = req.CoverageType
	existing.MinCoverage = req.MinCoverage
	existing.MaxCoverage = req.MaxCoverage
	existing.Summary = req.Summary
	existing.InsuranceDetail = req.InsuranceDetail
	existing.ProtectionDetail = req.ProtectionDetail
	existing.HowToClaim = req.HowToClaim
	if req.UpdatedBy != nil {
		existing.UpdatedBy = req.UpdatedBy
	}

	return u.repository.Update(id, existing)
}

func (u *usecase) Delete(id int) error {
	if id <= 0 {
		return errors.New("invalid product ID")
	}
	return u.repository.Delete(id)
}

// Template generation

func (u *usecase) GenerateTemplateDraftsForProducts(productCodes []string, createdBy int64) (int, error) {
	return u.repository.GenerateTemplateDraftsForProducts(productCodes, createdBy)
}

// Validation functions

func (u *usecase) validateCreateRequest(req product.CreateProductRequest) error {
	if req.Code == "" {
		return errors.New("product code is required")
	}
	if req.Name == "" {
		return errors.New("product name is required")
	}
	if req.InsuranceCode == "" {
		return errors.New("insurance code is required")
	}
	return nil
}

func (u *usecase) validateUpdateRequest(req product.UpdateProductRequest) error {
	if req.Code == "" {
		return errors.New("product code is required")
	}
	if req.Name == "" {
		return errors.New("product name is required")
	}
	if req.InsuranceCode == "" {
		return errors.New("insurance code is required")
	}
	return nil
}
