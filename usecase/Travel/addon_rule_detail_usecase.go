package usecase

import (
	"multi-onboarding/model/Travel"
	"multi-onboarding/repository/Travel"
)

// AddonRuleDetailUsecase handles business logic for addon rule details
type AddonRuleDetailUsecase struct {
	repo repository.AddonRuleDetailRepository
}

// NewAddonRuleDetailUsecase creates a new addon rule detail usecase
func NewAddonRuleDetailUsecase(repo repository.AddonRuleDetailRepository) *AddonRuleDetailUsecase {
	return &AddonRuleDetailUsecase{
		repo: repo,
	}
}

// AddDraft adds a draft addon rule detail
func (u *AddonRuleDetailUsecase) AddDraft(draft *model.AddonRuleDetailDraft) error {
	return u.repo.AddDraft(draft)
}

// GetAllDrafts retrieves all draft addon rule details
func (u *AddonRuleDetailUsecase) GetAllDrafts() ([]model.AddonRuleDetailDraft, error) {
	return u.repo.GetAllDrafts()
}

// ClearDrafts clears all draft addon rule details
func (u *AddonRuleDetailUsecase) ClearDrafts() error {
	return u.repo.ClearDrafts()
}

// ConfirmDrafts confirms all drafts
func (u *AddonRuleDetailUsecase) ConfirmDrafts(createdBy int64) error {
	return u.repo.ConfirmDrafts(createdBy)
}

// GetAll retrieves all addon rule details
func (u *AddonRuleDetailUsecase) GetAll() ([]model.AddonRuleDetail, error) {
	return u.repo.GetAll()
}

// GetByID retrieves an addon rule detail by ID
func (u *AddonRuleDetailUsecase) GetByID(id int64) (*model.AddonRuleDetail, error) {
	return u.repo.GetByID(id)
}

// Update updates an existing addon rule detail
func (u *AddonRuleDetailUsecase) Update(id int64, detail *model.AddonRuleDetail) error {
	return u.repo.Update(id, detail)
}

