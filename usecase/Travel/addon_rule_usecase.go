package usecase

import (
	"errors"
	"multi-onboarding/model/Travel"
	"multi-onboarding/repository/Travel"
)

// AddonRuleUsecase interface defines business logic operations for addon rules
type AddonRuleUsecase interface {
	GetAllAddonRules() ([]model.AddonRule, error)
	GetAddonRuleByID(id int64) (*model.AddonRule, error)
	GetAddonRulesByProductCode(productCode string) ([]model.AddonRule, error)
	GetAddonRulesByAddonCode(addonCode string) ([]model.AddonRule, error)
	CreateAddonRule(addonRule *model.AddonRule) error
	CreateAddonRulesBatch(addonRules []model.AddonRule) error
	UpdateAddonRule(id int64, addonRule *model.AddonRule) error
	DeleteAddonRule(id int64) error
}

// addonRuleUsecase implements AddonRuleUsecase
type addonRuleUsecase struct {
	repo repository.AddonRuleRepository
}

// NewAddonRuleUsecase creates a new addon rule usecase
func NewAddonRuleUsecase(repo repository.AddonRuleRepository) AddonRuleUsecase {
	return &addonRuleUsecase{
		repo: repo,
	}
}

// GetAllAddonRules retrieves all addon rules
func (uc *addonRuleUsecase) GetAllAddonRules() ([]model.AddonRule, error) {
	return uc.repo.GetAll()
}

// GetAddonRuleByID retrieves an addon rule by ID
func (uc *addonRuleUsecase) GetAddonRuleByID(id int64) (*model.AddonRule, error) {
	if id <= 0 {
		return nil, ErrInvalidAddonRuleID
	}
	addonRule, err := uc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrAddonRuleNotFound) {
			return nil, ErrAddonRuleNotFound
		}
		return nil, err
	}
	return addonRule, nil
}

// GetAddonRulesByProductCode retrieves addon rules by product code
func (uc *addonRuleUsecase) GetAddonRulesByProductCode(productCode string) ([]model.AddonRule, error) {
	return uc.repo.GetByProductCode(productCode)
}

// GetAddonRulesByAddonCode retrieves addon rules by addon code
func (uc *addonRuleUsecase) GetAddonRulesByAddonCode(addonCode string) ([]model.AddonRule, error) {
	return uc.repo.GetByAddonCode(addonCode)
}

// CreateAddonRule creates a new addon rule with validation
func (uc *addonRuleUsecase) CreateAddonRule(addonRule *model.AddonRule) error {
	if err := uc.validateAddonRule(addonRule); err != nil {
		return err
	}
	// Set default created_by if not provided
	if addonRule.CreatedBy == 0 {
		addonRule.CreatedBy = 1
	}
	return uc.repo.Create(addonRule)
}

// CreateAddonRulesBatch creates multiple addon rules in batch
func (uc *addonRuleUsecase) CreateAddonRulesBatch(addonRules []model.AddonRule) error {
	if len(addonRules) == 0 {
		return ErrEmptyAddonRulesBatch
	}
	
	// Validate all addon rules
	for i := range addonRules {
		if err := uc.validateAddonRule(&addonRules[i]); err != nil {
			return err
		}
		// Set default created_by if not provided
		if addonRules[i].CreatedBy == 0 {
			addonRules[i].CreatedBy = 1
		}
	}
	
	return uc.repo.CreateBatch(addonRules)
}

// UpdateAddonRule updates an existing addon rule
func (uc *addonRuleUsecase) UpdateAddonRule(id int64, addonRule *model.AddonRule) error {
	if id <= 0 {
		return ErrInvalidAddonRuleID
	}
	if err := uc.validateAddonRule(addonRule); err != nil {
		return err
	}
	err := uc.repo.Update(id, addonRule)
	if err != nil {
		if errors.Is(err, repository.ErrAddonRuleNotFound) {
			return ErrAddonRuleNotFound
		}
		return err
	}
	return nil
}

// DeleteAddonRule deletes an addon rule by ID
func (uc *addonRuleUsecase) DeleteAddonRule(id int64) error {
	if id <= 0 {
		return ErrInvalidAddonRuleID
	}
	err := uc.repo.Delete(id)
	if err != nil {
		if errors.Is(err, repository.ErrAddonRuleNotFound) {
			return ErrAddonRuleNotFound
		}
		return err
	}
	return nil
}

// validateAddonRule validates addon rule data
func (uc *addonRuleUsecase) validateAddonRule(addonRule *model.AddonRule) error {
	if addonRule.ProductCode == "" {
		return ErrAddonRuleProductCodeRequired
	}
	if addonRule.AddonCode == "" {
		return ErrAddonRuleAddonCodeRequired
	}
	return nil
}


