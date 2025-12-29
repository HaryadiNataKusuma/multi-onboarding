package usecase

import (
	"errors"

	"multi-onboarding/model/Travel"
	"multi-onboarding/repository/Travel"
)

var (
	ErrTemplateIDRequired = errors.New("template ID is required")
	ErrLocaleRequired     = errors.New("locale is required")
	ErrValueRequired      = errors.New("template value is required")
	ErrTemplateNotFound   = errors.New("template not found")
)

// TemplateUsecase handles template business logic
type TemplateUsecase struct {
	repo repository.TemplateRepository
}

// NewTemplateUsecase creates a new template usecase
func NewTemplateUsecase(repo repository.TemplateRepository) *TemplateUsecase {
	return &TemplateUsecase{repo: repo}
}

// AddDraft adds a new template draft
func (u *TemplateUsecase) AddDraft(draft *model.TemplateDraft) error {
	if err := u.validateTemplateDraft(draft); err != nil {
		return err
	}

	return u.repo.AddDraft(draft)
}

// GetAllDrafts retrieves all template drafts
func (u *TemplateUsecase) GetAllDrafts(insuranceCode string) ([]model.TemplateDraft, error) {
	return u.repo.GetAllDrafts(insuranceCode)
}

// GetDraftByID retrieves a template draft by ID and locale
func (u *TemplateUsecase) GetDraftByID(id, locale string) (*model.TemplateDraft, error) {
	return u.repo.GetDraftByID(id, locale)
}

// UpdateDraft updates a template draft
func (u *TemplateUsecase) UpdateDraft(id, locale string, draft *model.TemplateDraft) error {
	if err := u.validateTemplateDraft(draft); err != nil {
		return err
	}

	return u.repo.UpdateDraft(id, locale, draft)
}

// DeleteDraft deletes a template draft
func (u *TemplateUsecase) DeleteDraft(id, locale string) error {
	return u.repo.DeleteDraft(id, locale)
}

// ClearDrafts clears all template drafts
func (u *TemplateUsecase) ClearDrafts() error {
	return u.repo.ClearDrafts()
}

// ConfirmDrafts confirms all drafts and saves them to database
func (u *TemplateUsecase) ConfirmDrafts() error {
	return u.repo.ConfirmDrafts()
}

// Create creates a new template in database
func (u *TemplateUsecase) Create(template *model.Template) error {
	if err := u.validateTemplate(template); err != nil {
		return err
	}

	return u.repo.Create(template)
}

// GetByID retrieves a template by ID and locale
func (u *TemplateUsecase) GetByID(id, locale string) (*model.Template, error) {
	return u.repo.GetByID(id, locale)
}

// GetAll retrieves all templates
func (u *TemplateUsecase) GetAll(insuranceCode string) ([]model.Template, error) {
	return u.repo.GetAll(insuranceCode)
}

// Update updates a template
func (u *TemplateUsecase) Update(id, locale string, template *model.Template) error {
	if err := u.validateTemplate(template); err != nil {
		return err
	}

	return u.repo.Update(id, locale, template)
}

// Delete deletes a template
func (u *TemplateUsecase) Delete(id, locale string) error {
	return u.repo.Delete(id, locale)
}

// validateTemplateDraft validates template draft data
// Note: value is optional for draft (can be empty, user will fill later)
func (u *TemplateUsecase) validateTemplateDraft(draft *model.TemplateDraft) error {
	if draft.ID == "" {
		return ErrTemplateIDRequired
	}
	if draft.Locale == "" {
		return ErrLocaleRequired
	}
	// Value is optional for draft - can be empty when auto-created from products
	// User will fill the value later
	return nil
}

// validateTemplate validates template data
func (u *TemplateUsecase) validateTemplate(template *model.Template) error {
	if template.ID == "" {
		return ErrTemplateIDRequired
	}
	if template.Locale == "" {
		return ErrLocaleRequired
	}
	if template.Value == "" {
		return ErrValueRequired
	}
	return nil
}


