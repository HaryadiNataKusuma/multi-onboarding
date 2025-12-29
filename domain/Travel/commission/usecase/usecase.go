package usecase

import (
	"errors"
	"fmt"

	"multi-onboarding/domain/Travel/commission"
)

type usecase struct {
	repository commission.Repository
}

// NewCommissionUsecase creates a new commission usecase
func NewCommissionUsecase(repository commission.Repository) commission.Usecase {
	return &usecase{
		repository: repository,
	}
}

// Draft operations

func (u *usecase) AddDraft(draft *commission.CommissionDraft) error {
	if err := u.validateCommissionDraft(draft); err != nil {
		return err
	}
	return u.repository.AddDraft(draft)
}

func (u *usecase) GetAllDrafts() ([]commission.CommissionDraft, error) {
	return u.repository.GetAllDrafts()
}

func (u *usecase) GetDraftByID(id int) (*commission.CommissionDraft, error) {
	return u.repository.GetDraftByID(id)
}

func (u *usecase) UpdateDraft(id int, draft *commission.CommissionDraft) error {
	if err := u.validateCommissionDraft(draft); err != nil {
		return err
	}
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

func (u *usecase) GetAll() ([]commission.Commission, error) {
	return u.repository.GetAll()
}

func (u *usecase) GetByID(id int64) (*commission.Commission, error) {
	return u.repository.GetByID(id)
}

func (u *usecase) GetByProductCode(productCode string) (*commission.Commission, error) {
	return u.repository.GetByProductCode(productCode)
}

func (u *usecase) Create(req commission.CreateCommissionRequest) error {
	if err := u.validateCreateRequest(req); err != nil {
		return err
	}

	c := &commission.Commission{
		ProductCode:          req.ProductCode,
		CommissionPercentage: req.CommissionPercentage,
		CommissionVATType:    req.CommissionVATType,
		AFPercentage:         req.AFPercentage,
		AFVATType:            req.AFVATType,
		AdminFee:             req.AdminFee,
		CreatedBy:            req.CreatedBy,
	}

	if c.CreatedBy == 0 {
		c.CreatedBy = 1
	}

	return u.repository.Create(c)
}

func (u *usecase) Update(id int64, req commission.UpdateCommissionRequest) error {
	if err := u.validateUpdateRequest(req); err != nil {
		return err
	}

	existing, err := u.repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("commission not found")
	}

	existing.ProductCode = req.ProductCode
	existing.CommissionPercentage = req.CommissionPercentage
	existing.CommissionVATType = req.CommissionVATType
	existing.AFPercentage = req.AFPercentage
	existing.AFVATType = req.AFVATType
	existing.AdminFee = req.AdminFee
	if req.UpdatedBy != nil {
		existing.UpdatedBy = req.UpdatedBy
	}

	return u.repository.Update(id, existing)
}

func (u *usecase) Delete(id int64) error {
	return u.repository.Delete(id)
}

// Validation functions

func (u *usecase) validateCommissionDraft(draft *commission.CommissionDraft) error {
	if draft.ProductCode == "" {
		return errors.New("product code is required")
	}
	if draft.CommissionVATType != "EXCLUSIVE" && draft.CommissionVATType != "INCLUSIVE" {
		return errors.New("commission VAT type must be EXCLUSIVE or INCLUSIVE")
	}
	if draft.AFVATType != "EXCLUSIVE" && draft.AFVATType != "INCLUSIVE" {
		return errors.New("AF VAT type must be EXCLUSIVE or INCLUSIVE")
	}
	if draft.CommissionPercentage < 0 {
		return errors.New("commission percentage must be non-negative")
	}
	if draft.AFPercentage < 0 {
		return errors.New("AF percentage must be non-negative")
	}
	if draft.AdminFee < 0 {
		return errors.New("admin fee must be non-negative")
	}
	return nil
}

func (u *usecase) validateCreateRequest(req commission.CreateCommissionRequest) error {
	if req.ProductCode == "" {
		return errors.New("product code is required")
	}
	if req.CommissionVATType != "EXCLUSIVE" && req.CommissionVATType != "INCLUSIVE" {
		return errors.New("commission VAT type must be EXCLUSIVE or INCLUSIVE")
	}
	if req.AFVATType != "EXCLUSIVE" && req.AFVATType != "INCLUSIVE" {
		return errors.New("AF VAT type must be EXCLUSIVE or INCLUSIVE")
	}
	if req.CommissionPercentage < 0 {
		return errors.New("commission percentage must be non-negative")
	}
	if req.AFPercentage < 0 {
		return errors.New("AF percentage must be non-negative")
	}
	if req.AdminFee < 0 {
		return errors.New("admin fee must be non-negative")
	}
	return nil
}

func (u *usecase) validateUpdateRequest(req commission.UpdateCommissionRequest) error {
	if req.ProductCode == "" {
		return errors.New("product code is required")
	}
	if req.CommissionVATType != "EXCLUSIVE" && req.CommissionVATType != "INCLUSIVE" {
		return errors.New("commission VAT type must be EXCLUSIVE or INCLUSIVE")
	}
	if req.AFVATType != "EXCLUSIVE" && req.AFVATType != "INCLUSIVE" {
		return errors.New("AF VAT type must be EXCLUSIVE or INCLUSIVE")
	}
	if req.CommissionPercentage < 0 {
		return errors.New("commission percentage must be non-negative")
	}
	if req.AFPercentage < 0 {
		return errors.New("AF percentage must be non-negative")
	}
	if req.AdminFee < 0 {
		return errors.New("admin fee must be non-negative")
	}
	return nil
}

