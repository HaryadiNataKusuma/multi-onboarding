package usecase

import (
	"errors"
	"multi-onboarding/model/Travel"
	"multi-onboarding/repository/Travel"
)

// RegionUsecase interface defines business logic operations for regions
type RegionUsecase interface {
	GetAllRegions() ([]model.Region, error)
	GetRegionByID(id int) (*model.Region, error)
	CreateRegion(region *model.Region) error
	UpdateRegion(id int, region *model.Region) error
	DeleteRegion(id int) error
}

// regionUsecase implements RegionUsecase
type regionUsecase struct {
	repo repository.RegionRepository
}

// NewRegionUsecase creates a new region usecase
func NewRegionUsecase(repo repository.RegionRepository) RegionUsecase {
	return &regionUsecase{
		repo: repo,
	}
}

// GetAllRegions retrieves all regions
func (uc *regionUsecase) GetAllRegions() ([]model.Region, error) {
	return uc.repo.GetAll()
}

// GetRegionByID retrieves a region by ID
func (uc *regionUsecase) GetRegionByID(id int) (*model.Region, error) {
	if id <= 0 {
		return nil, ErrInvalidRegionID
	}
	region, err := uc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrRegionNotFound) {
			return nil, ErrRegionNotFound
		}
		return nil, err
	}
	return region, nil
}

// CreateRegion creates a new region with validation
func (uc *regionUsecase) CreateRegion(region *model.Region) error {
	if err := uc.validateRegion(region); err != nil {
		return err
	}
	return uc.repo.Create(region)
}

// UpdateRegion updates an existing region
func (uc *regionUsecase) UpdateRegion(id int, region *model.Region) error {
	if id <= 0 {
		return ErrInvalidRegionID
	}
	if err := uc.validateRegion(region); err != nil {
		return err
	}
	region.ID = id
	err := uc.repo.Update(region)
	if err != nil {
		if errors.Is(err, repository.ErrRegionNotFound) {
			return ErrRegionNotFound
		}
		return err
	}
	return nil
}

// DeleteRegion deletes a region by ID
func (uc *regionUsecase) DeleteRegion(id int) error {
	if id <= 0 {
		return ErrInvalidRegionID
	}
	err := uc.repo.Delete(id)
	if err != nil {
		if errors.Is(err, repository.ErrRegionNotFound) {
			return ErrRegionNotFound
		}
		return err
	}
	return nil
}

// validateRegion validates region data
func (uc *regionUsecase) validateRegion(region *model.Region) error {
	if region.Name == "" {
		return ErrRegionNameRequired
	}
	if region.Type != "WHITELIST" && region.Type != "BLACKLIST" {
		return ErrInvalidRegionType
	}
	if len(region.CountryIDs) == 0 {
		return ErrCountryIDRequired
	}
	// Validate all country IDs are positive
	for _, id := range region.CountryIDs {
		if id <= 0 {
			return ErrCountryIDRequired
		}
	}
	return nil
}

