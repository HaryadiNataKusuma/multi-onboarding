package usecase

import (
	"fmt"

	"multi-onboarding/model/Travel"
	"multi-onboarding/repository/Travel"
)

// InsuranceUsecase handles insurance business logic
type InsuranceUsecase struct {
	repo repository.InsuranceRepository
}

// NewInsuranceUsecase creates a new insurance usecase
func NewInsuranceUsecase(repo repository.InsuranceRepository) *InsuranceUsecase {
	return &InsuranceUsecase{repo: repo}
}

// AddDraft adds a new insurance draft
func (u *InsuranceUsecase) AddDraft(draft *model.InsuranceDraft) error {
	if err := u.validateInsuranceDraft(draft); err != nil {
		return err
	}

	return u.repo.AddDraft(draft)
}

// GetAllDrafts retrieves all insurance drafts
func (u *InsuranceUsecase) GetAllDrafts() ([]model.InsuranceDraft, error) {
	return u.repo.GetAllDrafts()
}

// GetDraftByID retrieves an insurance draft by ID
func (u *InsuranceUsecase) GetDraftByID(id int) (*model.InsuranceDraft, error) {
	return u.repo.GetDraftByID(id)
}

// UpdateDraft updates an insurance draft
func (u *InsuranceUsecase) UpdateDraft(id int, draft *model.InsuranceDraft) error {
	if err := u.validateInsuranceDraft(draft); err != nil {
		return err
	}

	return u.repo.UpdateDraft(id, draft)
}

// DeleteDraft deletes an insurance draft
func (u *InsuranceUsecase) DeleteDraft(id int) error {
	return u.repo.DeleteDraft(id)
}

// ClearDrafts clears all insurance drafts
func (u *InsuranceUsecase) ClearDrafts() error {
	return u.repo.ClearDrafts()
}

// ConfirmDrafts confirms all drafts and saves them to database
func (u *InsuranceUsecase) ConfirmDrafts() error {
	return u.repo.ConfirmDrafts()
}

// Create creates a new insurance in database
func (u *InsuranceUsecase) Create(insurance *model.Insurance) error {
	if err := u.validateInsurance(insurance); err != nil {
		return err
	}

	return u.repo.Create(insurance)
}

// GetByID retrieves an insurance by ID
func (u *InsuranceUsecase) GetByID(id int) (*model.Insurance, error) {
	return u.repo.GetByID(id)
}

// GetAll retrieves all insurances
func (u *InsuranceUsecase) GetAll() ([]model.Insurance, error) {
	return u.repo.GetAll()
}

// Update updates an insurance
func (u *InsuranceUsecase) Update(id int, insurance *model.Insurance) error {
	if err := u.validateInsurance(insurance); err != nil {
		return err
	}

	return u.repo.Update(id, insurance)
}

// Delete deletes an insurance
func (u *InsuranceUsecase) Delete(id int) error {
	return u.repo.Delete(id)
}

func (u *InsuranceUsecase) validateInsuranceDraft(draft *model.InsuranceDraft) error {
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

func (u *InsuranceUsecase) validateInsurance(insurance *model.Insurance) error {
	if insurance.Code == "" {
		return fmt.Errorf("code is required")
	}
	if insurance.Name == "" {
		return fmt.Errorf("name is required")
	}
	if insurance.Status != 0 && insurance.Status != 1 {
		return fmt.Errorf("status must be 0 or 1")
	}
	return nil
}

