package usecase

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/insurance"
	"multi-onboarding/utils"
)

type usecase struct {
	repository insurance.Repository
}

// NewInsuranceUsecase creates a new insurance usecase
func NewInsuranceUsecase(repository insurance.Repository) insurance.Usecase {
	return &usecase{
		repository: repository,
	}
}

// SaveInsuranceDraft saves an insurance draft to JSON file
func (u *usecase) SaveInsuranceDraft(c echo.Context, req insurance.SaveInsuranceDraftRequest) (*utils.InsuranceDraft, error) {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return nil, err
	}

	newID := utils.GetNextID(draftData)

	newInsurance := utils.InsuranceDraft{
		ID:        newID,
		Code:      strings.ToUpper(req.Code),
		Name:      req.Name,
		Status:    req.Status,
		Logo:      req.Logo,
		Timestamp: utils.GetTimestamp(),
	}

	draftData.Insurances = append(draftData.Insurances, newInsurance)

	if err := u.repository.SaveDraftData(c, draftData); err != nil {
		return nil, err
	}

	return &newInsurance, nil
}

// GetInsurancesDraft retrieves all insurance drafts
func (u *usecase) GetInsurancesDraft(c echo.Context, insuranceCode string) ([]utils.InsuranceDraft, error) {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return nil, err
	}

	if insuranceCode != "" {
		var filtered []utils.InsuranceDraft
		for _, ins := range draftData.Insurances {
			if strings.EqualFold(ins.Code, insuranceCode) {
				filtered = append(filtered, ins)
			}
		}
		return filtered, nil
	}

	return draftData.Insurances, nil
}

// UpdateInsuranceDraft updates an insurance draft
func (u *usecase) UpdateInsuranceDraft(c echo.Context, req insurance.UpdateInsuranceDraftRequest) error {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return err
	}

	found := false
	for i := range draftData.Insurances {
		if draftData.Insurances[i].ID == req.ID {
			draftData.Insurances[i].Code = strings.ToUpper(req.Code)
			draftData.Insurances[i].Name = req.Name
			draftData.Insurances[i].Status = req.Status
			draftData.Insurances[i].Logo = req.Logo
			draftData.Insurances[i].Timestamp = utils.GetTimestamp()
			found = true
			break
		}
	}

	if !found {
		return &notFoundError{Message: "Draft data not found"}
	}

	return u.repository.SaveDraftData(c, draftData)
}

// DeleteInsuranceDraft deletes an insurance draft
func (u *usecase) DeleteInsuranceDraft(c echo.Context, id int) error {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return err
	}

	found := false
	for i := range draftData.Insurances {
		if draftData.Insurances[i].ID == id {
			draftData.Insurances = append(draftData.Insurances[:i], draftData.Insurances[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return &notFoundError{Message: "Draft data not found"}
	}

	return u.repository.SaveDraftData(c, draftData)
}

// ClearInsurancesDraft clears all insurance drafts
func (u *usecase) ClearInsurancesDraft(c echo.Context) error {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return err
	}

	draftData.Insurances = []utils.InsuranceDraft{}
	return u.repository.SaveDraftData(c, draftData)
}

// ConfirmInsurances confirms insurances from draft to database
func (u *usecase) ConfirmInsurances(c echo.Context) (*insurance.ConfirmInsurancesResult, error) {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		log.Printf("ERROR loading draft data: %v", err)
		return nil, err
	}

	log.Printf("Loaded draft data: %d insurances found", len(draftData.Insurances))
	if len(draftData.Insurances) == 0 {
		log.Printf("No insurances in draft to confirm")
		return nil, &validationError{Message: "No data available in Draft to confirm"}
	}

	successCount := 0

	for _, ins := range draftData.Insurances {
		log.Printf("Processing insurance draft: code=%s, name=%s, status=%d", ins.Code, ins.Name, ins.Status)
		
		// Insert into vehicle_service_development.insurances
		if err := u.repository.InsertInsurance(c, ins.Code, ins.Name, ins.Status, ins.Logo); err != nil {
			log.Printf("ERROR inserting insurance %s: %v", ins.Code, err)
			continue // Skip this insurance if insert fails
		}
		
		log.Printf("Successfully inserted insurance %s", ins.Code)

		// Check and insert into quotation_service_local if not exists
		exists, err := u.repository.CheckInsuranceExistsInQuotation(c, ins.Code)
		if err != nil {
			log.Printf("Error checking vehicle_service_development.insurances: %v", err)
			// Continue anyway
		} else if !exists {
			isActive := 0
			if ins.Status > 0 {
				isActive = 1
			}
			if err := u.repository.InsertInsuranceToQuotation(c, ins.Code, ins.Name, isActive, ins.Logo); err != nil {
				log.Printf("Error inserting into vehicle_service_development.insurances: %v", err)
				// Continue anyway, don't fail the whole operation
			}
		}

		successCount++
		log.Printf("Insurance %s confirmed successfully. Total success: %d", ins.Code, successCount)
	}

	// Clear draft data
	draftData.Insurances = []utils.InsuranceDraft{}
	u.repository.SaveDraftData(c, draftData)

	return &insurance.ConfirmInsurancesResult{
		SuccessCount: successCount,
	}, nil
}

// GetInsurances retrieves all insurances from database
func (u *usecase) GetInsurances(c echo.Context) ([]map[string]interface{}, error) {
	return u.repository.GetAllInsurances(c)
}

// UpdateInsurance updates a confirmed insurance in the database
func (u *usecase) UpdateInsurance(c echo.Context, id int, req insurance.UpdateInsuranceRequest, userName string) error {
	beforeData, errBefore := u.repository.GetInsuranceBeforeUpdate(c, id)
	if errBefore != nil {
		log.Printf("Error getting data before update: %v", errBefore)
		beforeData = &insurance.InsuranceBeforeUpdate{
			Code:   "",
			Name:   "",
			Status: 0,
		}
	}

	beforeLogoStr := ""
	if beforeData.Logo.Valid {
		beforeLogoStr = beforeData.Logo.String
	}

	rowsAffected, err := u.repository.UpdateInsurance(c, id, req.Code, req.Name, req.Status, req.Logo)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &notFoundError{Message: "Insurance not found"}
	}

	// Compare and only log changed fields
	changedFields := make(map[string]interface{})

	if beforeData.Code != strings.ToUpper(req.Code) {
		changedFields["code"] = map[string]interface{}{
			"before": beforeData.Code,
			"after":  strings.ToUpper(req.Code),
		}
	}
	if beforeData.Name != req.Name {
		changedFields["name"] = map[string]interface{}{
			"before": beforeData.Name,
			"after":  req.Name,
		}
	}
	if beforeData.Status != req.Status {
		changedFields["status"] = map[string]interface{}{
			"before": beforeData.Status,
			"after":  req.Status,
		}
	}
	if beforeLogoStr != req.Logo {
		changedFields["logo"] = map[string]interface{}{
			"before": beforeLogoStr,
			"after":  req.Logo,
		}
	}

	// Always insert history if there are changes, or if we couldn't get before data
	if len(changedFields) > 0 || errBefore != nil {
		beforeChanged := make(map[string]interface{})
		afterChanged := make(map[string]interface{})

		// Always include code (insurance_code) for filtering
		beforeChanged["insurance_code"] = beforeData.Code
		afterChanged["insurance_code"] = strings.ToUpper(req.Code)

		// If we couldn't get before data, include all fields in history
		if errBefore != nil {
			beforeChanged["code"] = ""
			beforeChanged["name"] = ""
			beforeChanged["status"] = 0
			beforeChanged["logo"] = ""
			afterChanged["code"] = strings.ToUpper(req.Code)
			afterChanged["name"] = req.Name
			afterChanged["status"] = req.Status
			afterChanged["logo"] = req.Logo
		} else {
			// Add only changed fields
			for field, changeData := range changedFields {
				changeMap := changeData.(map[string]interface{})
				beforeChanged[field] = changeMap["before"]
				afterChanged[field] = changeMap["after"]
			}
		}

		beforeJSONBytes, _ := json.Marshal(beforeChanged)
		afterJSONBytes, _ := json.Marshal(afterChanged)

		// Insert history
		errHistory := u.repository.InsertHistory(
			c,
			userName,
			"Insurances",
			"UPDATE",
			id,
			"insurance",
			string(beforeJSONBytes),
			string(afterJSONBytes),
		)
		if errHistory != nil {
			log.Printf("Error inserting history: %v", errHistory)
		}
	}

	return nil
}

// Error types
type notFoundError struct {
	Message string
}

func (e *notFoundError) Error() string {
	return e.Message
}

type validationError struct {
	Message string
}

func (e *validationError) Error() string {
	return e.Message
}

