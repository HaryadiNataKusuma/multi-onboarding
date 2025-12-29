package usecase

import (
	"encoding/json"
	"errors"
	"fmt"

	"multi-onboarding/domain/Travel/region"
)

type usecase struct {
	repository region.Repository
}

// NewRegionUsecase creates a new region usecase
func NewRegionUsecase(repository region.Repository) region.Usecase {
	return &usecase{
		repository: repository,
	}
}

// GetAll retrieves all regions
func (u *usecase) GetAll() ([]region.Region, error) {
	return u.repository.GetAll()
}

// GetByID retrieves a region by ID
func (u *usecase) GetByID(id int) (*region.Region, error) {
	if id <= 0 {
		return nil, errors.New("invalid region ID")
	}
	return u.repository.GetByID(id)
}

// Create creates a new region with validation
func (u *usecase) Create(req region.CreateRegionRequest) error {
	if err := u.validateCreateRequest(req); err != nil {
		return err
	}

	reg := &region.Region{
		Name:       req.Name,
		Type:       req.Type,
		CountryIDs: req.CountryIDs,
		CreatedBy:  req.CreatedBy,
	}

	// Set default created_by if not provided
	if reg.CreatedBy == 0 {
		reg.CreatedBy = 1
	}

	return u.repository.Create(reg)
}

// Update updates an existing region
func (u *usecase) Update(id int, req region.UpdateRegionRequest) error {
	if id <= 0 {
		return errors.New("invalid region ID")
	}
	if err := u.validateUpdateRequest(req); err != nil {
		return err
	}

	// Get existing region
	existing, err := u.repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("region not found")
	}

	// Get data before update for history
	beforeData, errBefore := u.repository.GetRegionBeforeUpdate(id)
	if errBefore != nil {
		beforeData = &region.RegionBeforeUpdate{
			Name:       existing.Name,
			Type:       existing.Type,
			CountryIDs: existing.CountryIDs,
		}
	}

	// Update region
	existing.Name = req.Name
	existing.Type = req.Type
	existing.CountryIDs = req.CountryIDs
	if req.UpdatedBy != nil {
		existing.UpdatedBy = req.UpdatedBy
	}

	err = u.repository.Update(id, existing)
	if err != nil {
		return err
	}

	// Compare and log changed fields for history
	changedFields := make(map[string]interface{})
	if beforeData.Name != req.Name {
		changedFields["name"] = map[string]interface{}{
			"before": beforeData.Name,
			"after":  req.Name,
		}
	}
	if beforeData.Type != req.Type {
		changedFields["type"] = map[string]interface{}{
			"before": beforeData.Type,
			"after":  req.Type,
		}
	}

	// Compare country_ids
	beforeCountryIDsJSON, _ := json.Marshal(beforeData.CountryIDs)
	afterCountryIDsJSON, _ := json.Marshal(req.CountryIDs)
	if string(beforeCountryIDsJSON) != string(afterCountryIDsJSON) {
		changedFields["country_ids"] = map[string]interface{}{
			"before": beforeData.CountryIDs,
			"after":  req.CountryIDs,
		}
	}

	// Insert history if there are changes
	if len(changedFields) > 0 {
		beforeChanged := make(map[string]interface{})
		afterChanged := make(map[string]interface{})

		for field, changeData := range changedFields {
			changeMap := changeData.(map[string]interface{})
			beforeChanged[field] = changeMap["before"]
			afterChanged[field] = changeMap["after"]
		}

		beforeJSONBytes, _ := json.Marshal(beforeChanged)
		afterJSONBytes, _ := json.Marshal(afterChanged)

		_ = u.repository.InsertHistory(
			"system",
			"Regions",
			"UPDATE",
			id,
			"region",
			string(beforeJSONBytes),
			string(afterJSONBytes),
		)
	}

	return nil
}

// Delete deletes a region by ID
func (u *usecase) Delete(id int) error {
	if id <= 0 {
		return errors.New("invalid region ID")
	}
	return u.repository.Delete(id)
}

// validateCreateRequest validates region creation request
func (u *usecase) validateCreateRequest(req region.CreateRegionRequest) error {
	if req.Name == "" {
		return errors.New("region name is required")
	}
	if req.Type != "WHITELIST" && req.Type != "BLACKLIST" {
		return errors.New("region type must be WHITELIST or BLACKLIST")
	}
	if len(req.CountryIDs) == 0 {
		return errors.New("at least one country ID is required")
	}
	return nil
}

// validateUpdateRequest validates region update request
func (u *usecase) validateUpdateRequest(req region.UpdateRegionRequest) error {
	if req.Name == "" {
		return errors.New("region name is required")
	}
	if req.Type != "WHITELIST" && req.Type != "BLACKLIST" {
		return errors.New("region type must be WHITELIST or BLACKLIST")
	}
	if len(req.CountryIDs) == 0 {
		return errors.New("at least one country ID is required")
	}
	return nil
}

