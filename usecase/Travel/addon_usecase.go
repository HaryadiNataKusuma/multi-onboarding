package usecase

import (
	"errors"
	"multi-onboarding/model/Travel"
	"multi-onboarding/repository/Travel"
)

// AddonUsecase interface defines business logic operations for addons
type AddonUsecase interface {
	GetAllAddons() ([]model.Addon, error)
	GetAddonByID(id int64) (*model.Addon, error)
	GetAddonByCode(code string) (*model.Addon, error)
	CreateAddon(addon *model.Addon) error
	UpdateAddon(id int64, addon *model.Addon) error
	DeleteAddon(id int64) error
}

// addonUsecase implements AddonUsecase
type addonUsecase struct {
	repo repository.AddonRepository
}

// NewAddonUsecase creates a new addon usecase
func NewAddonUsecase(repo repository.AddonRepository) AddonUsecase {
	return &addonUsecase{
		repo: repo,
	}
}

// GetAllAddons retrieves all addons
func (uc *addonUsecase) GetAllAddons() ([]model.Addon, error) {
	return uc.repo.GetAll()
}

// GetAddonByID retrieves an addon by ID
func (uc *addonUsecase) GetAddonByID(id int64) (*model.Addon, error) {
	if id <= 0 {
		return nil, ErrInvalidAddonID
	}
	addon, err := uc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrAddonNotFound) {
			return nil, ErrAddonNotFound
		}
		return nil, err
	}
	return addon, nil
}

// GetAddonByCode retrieves an addon by code
func (uc *addonUsecase) GetAddonByCode(code string) (*model.Addon, error) {
	if code == "" {
		return nil, ErrAddonCodeRequired
	}
	addon, err := uc.repo.GetByCode(code)
	if err != nil {
		if errors.Is(err, repository.ErrAddonNotFound) {
			return nil, ErrAddonNotFound
		}
		return nil, err
	}
	return addon, nil
}

// CreateAddon creates a new addon with validation
func (uc *addonUsecase) CreateAddon(addon *model.Addon) error {
	if err := uc.validateAddon(addon); err != nil {
		return err
	}
	// Set default created_by if not provided
	if addon.CreatedBy == 0 {
		addon.CreatedBy = 1
	}
	// Set default is_active if not provided
	if addon.IsActive == 0 {
		addon.IsActive = 0 // Default to inactive
	}
	return uc.repo.Create(addon)
}

// UpdateAddon updates an existing addon
func (uc *addonUsecase) UpdateAddon(id int64, addon *model.Addon) error {
	if id <= 0 {
		return ErrInvalidAddonID
	}
	if err := uc.validateAddon(addon); err != nil {
		return err
	}
	addon.ID = id
	err := uc.repo.Update(addon)
	if err != nil {
		if errors.Is(err, repository.ErrAddonNotFound) {
			return ErrAddonNotFound
		}
		return err
	}
	return nil
}

// DeleteAddon deletes an addon by ID
func (uc *addonUsecase) DeleteAddon(id int64) error {
	if id <= 0 {
		return ErrInvalidAddonID
	}
	err := uc.repo.Delete(id)
	if err != nil {
		if errors.Is(err, repository.ErrAddonNotFound) {
			return ErrAddonNotFound
		}
		return err
	}
	return nil
}

// validateAddon validates addon data
func (uc *addonUsecase) validateAddon(addon *model.Addon) error {
	if addon.Code == "" {
		return ErrAddonCodeRequired
	}
	if addon.Name == "" {
		return ErrAddonNameRequired
	}
	if addon.NameMy == "" {
		return ErrAddonNameMyRequired
	}
	if addon.NameEn == "" {
		return ErrAddonNameEnRequired
	}
	return nil
}
