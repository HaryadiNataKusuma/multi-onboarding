package usecase

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/addon"
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

// GetAddons retrieves all addons
func (u *usecase) GetAddons(c echo.Context) ([]map[string]interface{}, error) {
	return u.repository.GetAllAddons(c)
}

// UpdateAddon updates an addon with history tracking
func (u *usecase) UpdateAddon(c echo.Context, id int, req addon.UpdateAddonRequest, userName string) error {
	// Get data before update for history
	beforeData, errBefore := u.repository.GetAddonBeforeUpdate(c, id)
	if errBefore != nil {
		log.Printf("Error getting data before update: %v", errBefore)
		// Continue with empty before data
		beforeData = &addon.AddonBeforeUpdate{
			Code:     "",
			Name:     "",
			IsActive: sql.NullBool{Valid: false, Bool: false},
		}
	}

	// Convert beforeIsActive to bool
	beforeIsActiveBool := false
	if beforeData.IsActive.Valid {
		beforeIsActiveBool = beforeData.IsActive.Bool
	}

	// Update addon
	rowsAffected, err := u.repository.UpdateAddon(c, id, req.Code, req.Name, req.IsActive)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &notFoundError{Message: "Addon not found"}
	}

	// If code changed, update addon_code in related tables
	if beforeData.Code != req.Code {
		// Update addon_rules
		if err := u.repository.UpdateAddonCodeInAddonRules(c, req.Code, beforeData.Code); err != nil {
			log.Printf("Error updating addon_code in addon_rules: %v", err)
			// Don't fail the main update
		}

		// Update mappings
		if err := u.repository.UpdateAddonCodeInMappings(c, req.Code, beforeData.Code); err != nil {
			log.Printf("Error updating addon_code in insurance_product_addon_mappings: %v", err)
			// Don't fail the main update
		}
	}

	// Compare and only log changed fields
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
	if beforeIsActiveBool != req.IsActive {
		changedFields["is_active"] = map[string]interface{}{
			"before": beforeIsActiveBool,
			"after":  req.IsActive,
		}
	}

	// Always insert history if there are changes, or if we couldn't get before data
	if len(changedFields) > 0 || errBefore != nil {
		// Build before and after JSON with only changed fields
		beforeChanged := make(map[string]interface{})
		afterChanged := make(map[string]interface{})

		// If we couldn't get before data, include all fields in history
		if errBefore != nil {
			beforeChanged["code"] = ""
			beforeChanged["name"] = ""
			beforeChanged["is_active"] = false
			afterChanged["code"] = req.Code
			afterChanged["name"] = req.Name
			afterChanged["is_active"] = req.IsActive
		} else {
			// Add only changed fields
			for field, changeData := range changedFields {
				changeMap := changeData.(map[string]interface{})
				beforeChanged[field] = changeMap["before"]
				afterChanged[field] = changeMap["after"]
			}
		}

		beforeJSONBytes, errMarshal := json.Marshal(beforeChanged)
		if errMarshal != nil {
			log.Printf("Error marshaling before data: %v", errMarshal)
			beforeJSONBytes = []byte("{}")
		}

		afterJSONBytes, errMarshal := json.Marshal(afterChanged)
		if errMarshal != nil {
			log.Printf("Error marshaling after data: %v", errMarshal)
			afterJSONBytes = []byte("{}")
		}

		// Insert history
		errHistory := u.repository.InsertHistory(
			c,
			userName,
			"Addons",
			"UPDATE",
			id,
			"addon",
			string(beforeJSONBytes),
			string(afterJSONBytes),
		)
		if errHistory != nil {
			log.Printf("Error inserting history: %v", errHistory)
			// Don't fail the update if history logging fails
		}
	}

	return nil
}

// ConfirmAddons confirms addons from request to database
func (u *usecase) ConfirmAddons(c echo.Context, req addon.ConfirmAddonsRequest) (int, error) {
	if len(req.Addons) == 0 {
		return 0, &validationError{Message: "No addons to confirm"}
	}

	insertedCount, err := u.repository.InsertAddons(c, req.Addons)
	if err != nil {
		return 0, err
	}

	return insertedCount, nil
}

// notFoundError represents a not found error
type notFoundError struct {
	Message string
}

func (e *notFoundError) Error() string {
	return e.Message
}

// validationError represents a validation error
type validationError struct {
	Message string
}

func (e *validationError) Error() string {
	return e.Message
}


