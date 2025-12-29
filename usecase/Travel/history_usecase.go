package usecase

import (
	"fmt"
	"multi-onboarding/model/Travel"
	"multi-onboarding/repository/Travel"
)

// HistoryUsecase handles business logic for histories
type HistoryUsecase struct {
	repo repository.HistoryRepository
}

// NewHistoryUsecase creates a new history usecase
func NewHistoryUsecase(repo repository.HistoryRepository) *HistoryUsecase {
	return &HistoryUsecase{
		repo: repo,
	}
}

// Create creates a new history record
func (u *HistoryUsecase) Create(history *model.History) error {
	if history.Section == "" || history.Action == "" {
		return fmt.Errorf("section and action cannot be empty")
	}
	return u.repo.Create(history)
}

// GetAll retrieves all history records
func (u *HistoryUsecase) GetAll() ([]model.History, error) {
	return u.repo.GetAll()
}

// GetByTableName retrieves all history records for a specific table
func (u *HistoryUsecase) GetByTableName(tableName string) ([]model.History, error) {
	return u.repo.GetByTableName(tableName)
}

// GetByRecordID retrieves all history records for a specific record
func (u *HistoryUsecase) GetByRecordID(tableName string, recordID int64) ([]model.History, error) {
	return u.repo.GetByRecordID(tableName, recordID)
}

// GetByID retrieves a history record by ID
func (u *HistoryUsecase) GetByID(id int64) (*model.History, error) {
	return u.repo.GetByID(id)
}

