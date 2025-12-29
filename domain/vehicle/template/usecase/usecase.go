package usecase

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/template"
	"multi-onboarding/utils"
)

type usecase struct {
	repository template.Repository
}

// NewTemplateUsecase creates a new template usecase
func NewTemplateUsecase(repository template.Repository) template.Usecase {
	return &usecase{
		repository: repository,
	}
}

// GetTemplatesDraft retrieves all template drafts
func (u *usecase) GetTemplatesDraft(c echo.Context, insuranceCode string) ([]utils.TemplateDraft, error) {
	log.Printf("GetTemplatesDraft called with insurance_code: %s", insuranceCode)

	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		log.Printf("ERROR: Failed to load draft data: %v", err)
		return nil, err
	}

	log.Printf("Total templates in draft: %d", len(draftData.Templates))

	// Filter out templates with empty TemplateID
	var validTemplates []utils.TemplateDraft
	for _, tmpl := range draftData.Templates {
		if strings.TrimSpace(tmpl.TemplateID) == "" {
			log.Printf("  - WARNING: Skipping template with empty TemplateID: TemplateID=%s", tmpl.TemplateID)
			continue
		}
		validTemplates = append(validTemplates, tmpl)
	}
	log.Printf("Valid templates (with non-empty TemplateID): %d", len(validTemplates))

	if insuranceCode != "" {
		var filtered []utils.TemplateDraft
		for _, tmpl := range validTemplates {
			templateInsuranceCode := tmpl.InsuranceCode
			if templateInsuranceCode == "" {
				templateInsuranceCode = tmpl.Insurance
			}
			log.Printf("  - Checking template: TemplateID=%s, InsuranceCode=%s, Insurance=%s",
				tmpl.TemplateID, tmpl.InsuranceCode, tmpl.Insurance)
			if strings.EqualFold(templateInsuranceCode, insuranceCode) {
				filtered = append(filtered, tmpl)
				log.Printf("    -> MATCH: Added to filtered list")
			}
		}
		log.Printf("Filtered templates count: %d for insurance_code %s", len(filtered), insuranceCode)
		return filtered, nil
	}

	log.Printf("No filter specified, returning all %d valid templates", len(validTemplates))
	return validTemplates, nil
}

// UpdateTemplateDraft updates a template draft
func (u *usecase) UpdateTemplateDraft(c echo.Context, req template.UpdateTemplateDraftRequest) error {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return err
	}

	found := false
	for i := range draftData.Templates {
		if draftData.Templates[i].TemplateID == req.ID {
			draftData.Templates[i].Locale = req.Locale
			draftData.Templates[i].Value = req.Value
			draftData.Templates[i].CreatedBy = req.CreatedBy
			found = true
			break
		}
	}

	if !found {
		return &notFoundError{Message: "Template draft not found"}
	}

	return u.repository.SaveDraftData(c, draftData)
}

// DeleteTemplateDraft deletes a template draft
func (u *usecase) DeleteTemplateDraft(c echo.Context, templateID string) error {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return err
	}

	found := false
	var updatedTemplates []utils.TemplateDraft
	for _, tmpl := range draftData.Templates {
		if tmpl.TemplateID == templateID {
			found = true
			continue
		}
		updatedTemplates = append(updatedTemplates, tmpl)
	}

	if !found {
		return &notFoundError{Message: "Template draft not found"}
	}

	draftData.Templates = updatedTemplates
	return u.repository.SaveDraftData(c, draftData)
}

// ClearTemplatesDraft clears all template drafts
func (u *usecase) ClearTemplatesDraft(c echo.Context) error {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return err
	}

	draftData.Templates = []utils.TemplateDraft{}
	return u.repository.SaveDraftData(c, draftData)
}

// ConfirmTemplates confirms templates from draft to database
func (u *usecase) ConfirmTemplates(c echo.Context, insuranceCodeFilter string) (*template.ConfirmTemplatesResult, error) {
	draftData, err := u.repository.LoadDraftData(c)
	if err != nil {
		return nil, err
	}

	if len(draftData.Templates) == 0 {
		return nil, &validationError{Message: "No template drafts available to confirm"}
	}

	var templatesToConfirm []utils.TemplateDraft
	if insuranceCodeFilter != "" {
		for _, tmpl := range draftData.Templates {
			if strings.TrimSpace(tmpl.TemplateID) == "" {
				log.Printf("WARNING: Skipping template with empty TemplateID during confirmation: TemplateID=%s", tmpl.TemplateID)
				continue
			}
			templateInsuranceCode := tmpl.InsuranceCode
			if templateInsuranceCode == "" {
				templateInsuranceCode = tmpl.Insurance
			}
			if strings.EqualFold(templateInsuranceCode, insuranceCodeFilter) {
				templatesToConfirm = append(templatesToConfirm, tmpl)
			}
		}
		if len(templatesToConfirm) == 0 {
			return nil, &validationError{Message: "No templates found for the specified insurance code in draft"}
		}
	} else {
		for _, tmpl := range draftData.Templates {
			if strings.TrimSpace(tmpl.TemplateID) == "" {
				log.Printf("WARNING: Skipping template with empty TemplateID during confirmation: TemplateID=%s", tmpl.TemplateID)
				continue
			}
			templatesToConfirm = append(templatesToConfirm, tmpl)
		}
		if len(templatesToConfirm) == 0 {
			return nil, &validationError{Message: "No valid template drafts available to confirm (all have empty template_id)"}
		}
	}

	successCount := 0
	errorCount := 0
	var errors []string
	confirmedTemplateIDs := make(map[string]bool)

	for _, tmpl := range templatesToConfirm {
		valueStr := ""
		if len(tmpl.Value) > 0 {
			valueBytes, err := json.Marshal(tmpl.Value)
			if err != nil {
				log.Printf("Error marshalling value for template %s: %v", tmpl.TemplateID, err)
				valueStr = "[]"
			} else {
				valueStr = string(valueBytes)
			}
		} else {
			valueStr = `["","",""]`
		}

		if err := u.repository.InsertTemplate(c, tmpl.Locale, tmpl.TemplateID, valueStr, tmpl.CreatedBy); err != nil {
			errorCount++
			errorMsg := "Failed to insert template " + tmpl.TemplateID + ": " + err.Error()
			log.Printf("Error inserting template %s: %v", tmpl.TemplateID, err)
			errors = append(errors, errorMsg)
			continue
		}

		successCount++
		confirmedTemplateIDs[tmpl.TemplateID] = true
		log.Printf("Successfully confirmed template %s", tmpl.TemplateID)
	}

	// Remove successfully confirmed templates from draft
	if successCount > 0 {
		if insuranceCodeFilter != "" {
			var remainingTemplates []utils.TemplateDraft
			for _, draftTemplate := range draftData.Templates {
				templateInsuranceCode := draftTemplate.InsuranceCode
				if templateInsuranceCode == "" {
					templateInsuranceCode = draftTemplate.Insurance
				}
				if !strings.EqualFold(templateInsuranceCode, insuranceCodeFilter) || !confirmedTemplateIDs[draftTemplate.TemplateID] {
					remainingTemplates = append(remainingTemplates, draftTemplate)
				}
			}
			draftData.Templates = remainingTemplates
		} else {
			var remainingTemplates []utils.TemplateDraft
			for _, draftTemplate := range draftData.Templates {
				if !confirmedTemplateIDs[draftTemplate.TemplateID] {
					remainingTemplates = append(remainingTemplates, draftTemplate)
				}
			}
			draftData.Templates = remainingTemplates
		}

		u.repository.SaveDraftData(c, draftData)
	}

	result := &template.ConfirmTemplatesResult{
		SuccessCount: successCount,
		ErrorCount:   errorCount,
		TotalCount:   len(templatesToConfirm),
		Errors:       errors,
	}

	if successCount == 0 {
		return result, &validationError{Message: "Failed to confirm any templates"}
	}

	return result, nil
}

// GetTemplates retrieves all templates from database
func (u *usecase) GetTemplates(c echo.Context, insuranceCode string) ([]map[string]interface{}, error) {
	return u.repository.GetAllTemplates(c, insuranceCode)
}

// UpdateTemplate updates a confirmed template in the database
func (u *usecase) UpdateTemplate(c echo.Context, id string, req template.UpdateTemplateRequest, userName string) error {
	oldLocale := req.Locale
	beforeData, errBefore := u.repository.GetTemplateBeforeUpdate(c, oldLocale, id)
	if errBefore != nil {
		oldLocale = "id"
		beforeData, errBefore = u.repository.GetTemplateBeforeUpdate(c, oldLocale, id)
		if errBefore != nil {
			log.Printf("Error getting data before update: %v", errBefore)
			beforeData = &template.TemplateBeforeUpdate{}
			oldLocale = "id"
		}
	}

	rowsAffected, err := u.repository.UpdateTemplate(c, oldLocale, id, req.Locale, req.ID, req.Value, req.CreatedBy)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &notFoundError{Message: "Template not found"}
	}

	// Compare and only log changed fields
	changedFields := make(map[string]interface{})

	if beforeData.Locale.Valid && beforeData.Locale.String != req.Locale {
		changedFields["locale"] = map[string]interface{}{"before": beforeData.Locale.String, "after": req.Locale}
	}
	if beforeData.TemplateID.Valid && beforeData.TemplateID.String != req.ID {
		changedFields["id"] = map[string]interface{}{"before": beforeData.TemplateID.String, "after": req.ID}
	}
	if beforeData.Value.Valid && beforeData.Value.String != req.Value {
		changedFields["value"] = map[string]interface{}{"before": beforeData.Value.String, "after": req.Value}
	}
	if beforeData.CreatedBy.Valid && int(beforeData.CreatedBy.Int64) != req.CreatedBy {
		changedFields["created_by"] = map[string]interface{}{"before": beforeData.CreatedBy.Int64, "after": req.CreatedBy}
	}

	// Always insert history
	beforeChanged := make(map[string]interface{})
	afterChanged := make(map[string]interface{})

	if len(changedFields) > 0 {
		for field, changeData := range changedFields {
			changeMap := changeData.(map[string]interface{})
			beforeChanged[field] = changeMap["before"]
			afterChanged[field] = changeMap["after"]
		}
	} else {
		if beforeData.Value.Valid {
			beforeChanged["value"] = beforeData.Value.String
		}
		afterChanged["value"] = req.Value
		log.Printf("No field changes detected for template %s, but logging update anyway", id)
	}

	beforeJSONBytes, _ := json.Marshal(beforeChanged)
	afterJSONBytes, _ := json.Marshal(afterChanged)

	errHistory := u.repository.InsertHistory(
		c,
		userName,
		"Templates",
		"UPDATE",
		id,
		"template",
		string(beforeJSONBytes),
		string(afterJSONBytes),
	)
	if errHistory != nil {
		log.Printf("Error inserting history: %v", errHistory)
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


