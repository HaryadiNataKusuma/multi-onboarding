package usecase

import "multi-onboarding/domain/Travel/country"

type usecase struct {
	repository country.Repository
}

// NewCountryUsecase creates a new country usecase
func NewCountryUsecase(repository country.Repository) country.Usecase {
	return &usecase{
		repository: repository,
	}
}

// GetAll retrieves all countries
func (u *usecase) GetAll() ([]country.Country, error) {
	return u.repository.GetAll()
}

// GetByID retrieves a country by ID
func (u *usecase) GetByID(id int) (*country.Country, error) {
	return u.repository.GetByID(id)
}

// GetByCode retrieves a country by code
func (u *usecase) GetByCode(code string) (*country.Country, error) {
	return u.repository.GetByCode(code)
}

