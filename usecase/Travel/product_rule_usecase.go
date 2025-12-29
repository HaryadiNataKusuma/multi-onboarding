package usecase

import (
	"errors"

	"multi-onboarding/model/Travel"
	"multi-onboarding/repository/Travel"
)

var (
	ErrProductCodeRequired      = errors.New("product code is required")
	ErrPremiumTypeRequired     = errors.New("premium type is required")
	ErrStartDaysRequired        = errors.New("start days is required")
	ErrEndDaysRequired          = errors.New("end days is required")
	ErrBasePremiumValueRequired = errors.New("base premium value is required")
	ErrRulesRequired            = errors.New("rules is required")
)

// ProductRuleUsecase handles product rule business logic
type ProductRuleUsecase struct {
	repo repository.ProductRuleRepository
}

// NewProductRuleUsecase creates a new product rule usecase
func NewProductRuleUsecase(repo repository.ProductRuleRepository) *ProductRuleUsecase {
	return &ProductRuleUsecase{
		repo: repo,
	}
}

// AddDraft adds a draft product rule
func (u *ProductRuleUsecase) AddDraft(draft *model.ProductRuleDraft) error {
	if err := u.validateDraft(draft); err != nil {
		return err
	}

	return u.repo.AddDraft(draft)
}

// GetAllDrafts gets all draft product rules
func (u *ProductRuleUsecase) GetAllDrafts() ([]model.ProductRuleDraft, error) {
	return u.repo.GetAllDrafts()
}

// UpdateDraft updates a draft product rule by index
func (u *ProductRuleUsecase) UpdateDraft(index int, draft *model.ProductRuleDraft) error {
	if err := u.validateDraft(draft); err != nil {
		return err
	}

	return u.repo.UpdateDraft(index, draft)
}

// DeleteDraft deletes a draft product rule by index
func (u *ProductRuleUsecase) DeleteDraft(index int) error {
	return u.repo.DeleteDraft(index)
}

// ClearDrafts clears all draft product rules
func (u *ProductRuleUsecase) ClearDrafts() error {
	return u.repo.ClearDrafts()
}

// ConfirmDrafts confirms all draft product rules to database
func (u *ProductRuleUsecase) ConfirmDrafts(createdBy int64) error {
	return u.repo.ConfirmDrafts(createdBy)
}

// Create creates a new product rule
func (u *ProductRuleUsecase) Create(rule *model.ProductRule) error {
	if err := u.validateRule(rule); err != nil {
		return err
	}

	return u.repo.Create(rule)
}

// GetAll gets all product rules
func (u *ProductRuleUsecase) GetAll(productCode string, insuranceCode string) ([]model.ProductRule, error) {
	return u.repo.GetAll(productCode, insuranceCode)
}

// GetByID gets a product rule by ID
func (u *ProductRuleUsecase) GetByID(id int64) (*model.ProductRule, error) {
	return u.repo.GetByID(id)
}

// Update updates a product rule
func (u *ProductRuleUsecase) Update(id int64, rule *model.ProductRule) error {
	if err := u.validateRule(rule); err != nil {
		return err
	}

	return u.repo.Update(id, rule)
}

// Delete deletes a product rule
func (u *ProductRuleUsecase) Delete(id int64) error {
	return u.repo.Delete(id)
}

// validateDraft validates a draft product rule
func (u *ProductRuleUsecase) validateDraft(draft *model.ProductRuleDraft) error {
	if draft.ProductCode == "" {
		return ErrProductCodeRequired
	}
	if draft.PremiumType == "" {
		return ErrPremiumTypeRequired
	}
	if draft.StartDays < 0 {
		return ErrStartDaysRequired
	}
	if draft.EndDays < 0 {
		return ErrEndDaysRequired
	}
	if draft.BasePremiumValue < 0 {
		return ErrBasePremiumValueRequired
	}
	if draft.Rules == "" {
		return ErrRulesRequired
	}
	return nil
}

// validateRule validates a product rule
func (u *ProductRuleUsecase) validateRule(rule *model.ProductRule) error {
	if rule.ProductCode == "" {
		return ErrProductCodeRequired
	}
	if rule.PremiumType == "" {
		return ErrPremiumTypeRequired
	}
	if rule.StartDays < 0 {
		return ErrStartDaysRequired
	}
	if rule.EndDays < 0 {
		return ErrEndDaysRequired
	}
	if rule.BasePremiumValue < 0 {
		return ErrBasePremiumValueRequired
	}
	if rule.Rules == "" {
		return ErrRulesRequired
	}
	return nil
}

