package usecase

import "multi-onboarding/domain/Travel/history"

type usecase struct {
	repository history.Repository
}

// NewHistoryUsecase creates a new history usecase
func NewHistoryUsecase(repository history.Repository) history.Usecase {
	return &usecase{
		repository: repository,
	}
}

// GetAll retrieves all history records
func (u *usecase) GetAll() ([]history.History, error) {
	return u.repository.GetAll()
}

// GetByID retrieves a history record by ID
func (u *usecase) GetByID(id int64) (*history.History, error) {
	return u.repository.GetByID(id)
}

// GetByTableName retrieves history records by table name
func (u *usecase) GetByTableName(tableName string) ([]history.History, error) {
	return u.repository.GetByTableName(tableName)
}

// GetByRecordID retrieves history records by table name and record ID
func (u *usecase) GetByRecordID(tableName string, recordID int) ([]history.History, error) {
	return u.repository.GetByRecordID(tableName, recordID)
}

