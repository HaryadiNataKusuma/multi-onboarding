package usecase

import (
	"fmt"

	"multi-onboarding/domain/Travel/insurance"
)

type usecase struct {
	repository insurance.Repository
}

// NewInsuranceUsecase creates a new insurance usecase
func NewInsuranceUsecase(repository insurance.Repository) insurance.Usecase {
	return &usecase{
		repository: repository,
	}
}

// Draft operations

func (u *usecase) AddDraft(draft *insurance.InsuranceDraft) error {
	if err := u.validateInsuranceDraft(draft); err != nil {
		return err
	}
	return u.repository.AddDraft(draft)
}

func (u *usecase) GetAllDrafts() ([]insurance.InsuranceDraft, error) {
	return u.repository.GetAllDrafts()
}

func (u *usecase) GetDraftByID(id int) (*insurance.InsuranceDraft, error) {
	return u.repository.GetDraftByID(id)
}

func (u *usecase) UpdateDraft(id int, draft *insurance.InsuranceDraft) error {
	if err := u.validateInsuranceDraft(draft); err != nil {
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

func (u *usecase) ConfirmDrafts() error {
	return u.repository.ConfirmDrafts()
}

// Database operations

func (u *usecase) Create(req insurance.CreateInsuranceRequest) error {
	if err := u.validateCreateRequest(req); err != nil {
		return err
	}

	ins := &insurance.Insurance{
		Code:   req.Code,
		Name:   req.Name,
		Status: req.Status,
		Logo:   req.Logo,
	}

	return u.repository.Create(ins)
}

func (u *usecase) GetByID(id int) (*insurance.Insurance, error) {
	return u.repository.GetByID(id)
}

func (u *usecase) GetAll() ([]insurance.Insurance, error) {
	return u.repository.GetAll()
}

func (u *usecase) Update(id int, req insurance.UpdateInsuranceRequest) error {
	if err := u.validateUpdateRequest(req); err != nil {
		return err
	}

	ins := &insurance.Insurance{
		ID:     id,
		Code:   req.Code,
		Name:   req.Name,
		Status: req.Status,
		Logo:   req.Logo,
	}

	return u.repository.Update(id, ins)
}

func (u *usecase) Delete(id int) error {
	return u.repository.Delete(id)
}

// Validation functions

func (u *usecase) validateInsuranceDraft(draft *insurance.InsuranceDraft) error {
	if draft.Code == "" {
		return fmt.Errorf("code is required")
	}
	if draft.Name == "" {
		return fmt.Errorf("name is required")
	}
	if draft.Status != 0 && draft.Status != 1 {
		return fmt.Errorf("status must be 0 or 1")
	}
	return nil
}

func (u *usecase) validateCreateRequest(req insurance.CreateInsuranceRequest) error {
	if req.Code == "" {
		return fmt.Errorf("code is required")
	}
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.Status != 0 && req.Status != 1 {
		return fmt.Errorf("status must be 0 or 1")
	}
	return nil
}

func (u *usecase) validateUpdateRequest(req insurance.UpdateInsuranceRequest) error {
	if req.Code == "" {
		return fmt.Errorf("code is required")
	}
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.Status != 0 && req.Status != 1 {
		return fmt.Errorf("status must be 0 or 1")
	}
	return nil
}

