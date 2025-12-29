package usecase

import (
	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/history"
)

type usecase struct {
	repository history.Repository
}

// NewHistoryUsecase creates a new history usecase
func NewHistoryUsecase(repository history.Repository) history.Usecase {
	return &usecase{
		repository: repository,
	}
}

// GetHistories retrieves all histories
func (u *usecase) GetHistories(c echo.Context, insuranceCode, section string) ([]map[string]interface{}, error) {
	return u.repository.GetAllHistories(c, insuranceCode, section)
}


