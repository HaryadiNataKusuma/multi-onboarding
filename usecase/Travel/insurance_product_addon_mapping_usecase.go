package usecase

import (
	"multi-onboarding/model/Travel"
	"multi-onboarding/repository/Travel"
)

// InsuranceProductAddonMappingUsecase handles business logic for insurance product addon mappings
type InsuranceProductAddonMappingUsecase struct {
	repo repository.InsuranceProductAddonMappingRepository
}

// NewInsuranceProductAddonMappingUsecase creates a new insurance product addon mapping usecase
func NewInsuranceProductAddonMappingUsecase(repo repository.InsuranceProductAddonMappingRepository) *InsuranceProductAddonMappingUsecase {
	return &InsuranceProductAddonMappingUsecase{
		repo: repo,
	}
}

// Create creates a new insurance product addon mapping
func (u *InsuranceProductAddonMappingUsecase) Create(mapping *model.InsuranceProductAddonMapping) error {
	return u.repo.Create(mapping)
}

// CreateBatch creates multiple insurance product addon mappings
func (u *InsuranceProductAddonMappingUsecase) CreateBatch(mappings []model.InsuranceProductAddonMapping) error {
	return u.repo.CreateBatch(mappings)
}

