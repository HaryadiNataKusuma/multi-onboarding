package usecase

import (
	"errors"
	"fmt"

	"multi-onboarding/domain/Travel/template"
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

// Draft operations

func (u *usecase) AddDraft(draft *template.TemplateDraft) error {
	if err := u.validateTemplateDraft(draft); err != nil {
		return err
	}
	return u.repository.AddDraft(draft)
}

func (u *usecase) GetAllDrafts(insuranceCode string) ([]template.TemplateDraft, error) {
	return u.repository.GetAllDrafts(insuranceCode)
}

func (u *usecase) GetDraftByID(id string) (*template.TemplateDraft, error) {
	return u.repository.GetDraftByID(id)
}

func (u *usecase) UpdateDraft(id string, draft *template.TemplateDraft) error {
	if err := u.validateTemplateDraft(draft); err != nil {
		return err
	}
	return u.repository.UpdateDraft(id, draft)
}

func (u *usecase) DeleteDraft(id string) error {
	return u.repository.DeleteDraft(id)
}

func (u *usecase) ClearDrafts() error {
	return u.repository.ClearDrafts()
}

func (u *usecase) ConfirmDrafts(insuranceCodeFilter string) error {
	return u.repository.ConfirmDrafts(insuranceCodeFilter)
}

// Database operations

func (u *usecase) GetAll(insuranceCode string) ([]template.Template, error) {
	return u.repository.GetAll(insuranceCode)
}

func (u *usecase) GetByID(locale, id string) (*template.Template, error) {
	return u.repository.GetByID(locale, id)
}

func (u *usecase) Create(req template.CreateTemplateRequest) error {
	if err := u.validateCreateRequest(req); err != nil {
		return err
	}

	t := &template.Template{
		Locale:    req.Locale,
		ID:        req.ID,
		Value:     req.Value,
		CreatedBy: req.CreatedBy,
	}

	if t.CreatedBy == 0 {
		t.CreatedBy = 1
	}

	return u.repository.Create(t)
}

func (u *usecase) Update(locale, id string, req template.UpdateTemplateRequest) error {
	if err := u.validateUpdateRequest(req); err != nil {
		return err
	}

	existing, err := u.repository.GetByID(locale, id)
	if err != nil {
		return fmt.Errorf("template not found")
	}

	existing.Value = req.Value
	if req.UpdatedBy != nil {
		existing.UpdatedBy = req.UpdatedBy
	}

	return u.repository.Update(locale, id, existing)
}

func (u *usecase) Delete(locale, id string) error {
	return u.repository.Delete(locale, id)
}

// Validation functions

func (u *usecase) validateTemplateDraft(draft *template.TemplateDraft) error {
	if draft.ID == "" {
		return errors.New("template ID is required")
	}
	if draft.Locale == "" {
		return errors.New("locale is required")
	}
	return nil
}

func (u *usecase) validateCreateRequest(req template.CreateTemplateRequest) error {
	if req.ID == "" {
		return errors.New("template ID is required")
	}
	if req.Locale == "" {
		return errors.New("locale is required")
	}
	if req.Value == "" {
		return errors.New("template value is required")
	}
	return nil
}

func (u *usecase) validateUpdateRequest(req template.UpdateTemplateRequest) error {
	if req.ID == "" {
		return errors.New("template ID is required")
	}
	if req.Locale == "" {
		return errors.New("locale is required")
	}
	if req.Value == "" {
		return errors.New("template value is required")
	}
	return nil
}

