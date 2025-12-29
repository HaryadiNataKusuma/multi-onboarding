package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"multi-onboarding/domain/Travel/addon"
)

type usecase struct {
	repository addon.Repository
}

// NewAddonUsecase creates a new addon usecase
func NewAddonUsecase(repository addon.Repository) addon.Usecase {
	return &usecase{
		repository: repository,
	}
}

// GetAllAddons retrieves all addons
func (u *usecase) GetAllAddons() ([]addon.Addon, error) {
	return u.repository.GetAll()
}

// GetAddonByID retrieves an addon by ID
func (u *usecase) GetAddonByID(id int64) (*addon.Addon, error) {
	if id <= 0 {
		return nil, errors.New("invalid addon ID")
	}
	return u.repository.GetByID(id)
}

// GetAddonByCode retrieves an addon by code
func (u *usecase) GetAddonByCode(code string) (*addon.Addon, error) {
	if code == "" {
		return nil, errors.New("addon code is required")
	}
	return u.repository.GetByCode(code)
}

// CreateAddon creates a new addon with validation
func (u *usecase) CreateAddon(req addon.CreateAddonRequest) error {
	if err := u.validateAddonRequest(req); err != nil {
		return err
	}

	a := &addon.Addon{
		Code:      req.Code,
		Name:      req.Name,
		NameMy:    req.NameMy,
		NameEn:    req.NameEn,
		IsActive:  req.IsActive,
		CreatedBy: req.CreatedBy,
	}

	// Set default created_by if not provided
	if a.CreatedBy == 0 {
		a.CreatedBy = 1
	}

	return u.repository.Create(a)
}

// UpdateAddon updates an existing addon
func (u *usecase) UpdateAddon(id int64, req addon.UpdateAddonRequest) error {
	if id <= 0 {
		return errors.New("invalid addon ID")
	}
	if err := u.validateUpdateRequest(req); err != nil {
		return err
	}

	// Get existing addon
	existing, err := u.repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("addon not found")
	}

	// Get data before update for history
	beforeData, errBefore := u.repository.GetAddonBeforeUpdate(id)
	if errBefore != nil {
		log.Printf("Error getting data before update: %v", errBefore)
		beforeData = &addon.AddonBeforeUpdate{
			Code:     existing.Code,
			Name:     existing.Name,
			NameMy:   existing.NameMy,
			NameEn:   existing.NameEn,
			IsActive: existing.IsActive,
		}
	}

	// Update addon
	existing.Code = req.Code
	existing.Name = req.Name
	existing.NameMy = req.NameMy
	existing.NameEn = req.NameEn
	existing.IsActive = req.IsActive
	if req.UpdatedBy != nil {
		existing.UpdatedBy = req.UpdatedBy
	}

	err = u.repository.Update(existing)
	if err != nil {
		return err
	}

	// Compare and log changed fields for history
	changedFields := make(map[string]interface{})
	if beforeData.Code != req.Code {
		changedFields["code"] = map[string]interface{}{
			"before": beforeData.Code,
			"after":  req.Code,
		}
	}
	if beforeData.Name != req.Name {
		changedFields["name"] = map[string]interface{}{
			"before": beforeData.Name,
			"after":  req.Name,
		}
	}
	if beforeData.NameMy != req.NameMy {
		changedFields["name_my"] = map[string]interface{}{
			"before": beforeData.NameMy,
			"after":  req.NameMy,
		}
	}
	if beforeData.NameEn != req.NameEn {
		changedFields["name_en"] = map[string]interface{}{
			"before": beforeData.NameEn,
			"after":  req.NameEn,
		}
	}
	if beforeData.IsActive != req.IsActive {
		changedFields["is_active"] = map[string]interface{}{
			"before": beforeData.IsActive,
			"after":  req.IsActive,
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

		// Insert history (userName would come from context in real implementation)
		_ = u.repository.InsertHistory(
			"system",
			"Addons",
			"UPDATE",
			id,
			"addon",
			string(beforeJSONBytes),
			string(afterJSONBytes),
		)
	}

	return nil
}

// DeleteAddon deletes an addon by ID
func (u *usecase) DeleteAddon(id int64) error {
	if id <= 0 {
		return errors.New("invalid addon ID")
	}
	return u.repository.Delete(id)
}

// validateAddonRequest validates addon creation request
func (u *usecase) validateAddonRequest(req addon.CreateAddonRequest) error {
	if req.Code == "" {
		return errors.New("addon code is required")
	}
	if req.Name == "" {
		return errors.New("addon name is required")
	}
	if req.NameMy == "" {
		return errors.New("addon name_my is required")
	}
	if req.NameEn == "" {
		return errors.New("addon name_en is required")
	}
	return nil
}

// validateUpdateRequest validates addon update request
func (u *usecase) validateUpdateRequest(req addon.UpdateAddonRequest) error {
	if req.Code == "" {
		return errors.New("addon code is required")
	}
	if req.Name == "" {
		return errors.New("addon name is required")
	}
	if req.NameMy == "" {
		return errors.New("addon name_my is required")
	}
	if req.NameEn == "" {
		return errors.New("addon name_en is required")
	}
	return nil
}

