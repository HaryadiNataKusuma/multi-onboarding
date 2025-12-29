package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"multi-onboarding/model/Travel"
)

var (
	ErrHistoryNotFound = errors.New("history not found")
)

// HistoryRepository defines the interface for history operations
type HistoryRepository interface {
	Create(history *model.History) error
	GetAll() ([]model.History, error)
	GetByTableName(tableName string) ([]model.History, error)
	GetByRecordID(tableName string, recordID int64) ([]model.History, error)
	GetByID(id int64) (*model.History, error)
}

// MySQLHistoryRepository implements HistoryRepository using MySQL database
type MySQLHistoryRepository struct {
	db *sql.DB
}

// NewMySQLHistoryRepository creates a new MySQL history repository
func NewMySQLHistoryRepository(db *sql.DB) *MySQLHistoryRepository {
	return &MySQLHistoryRepository{
		db: db,
	}
}

// Create creates a new history record
func (r *MySQLHistoryRepository) Create(history *model.History) error {
	query := `
		INSERT INTO travel_service_development.histories 
		(user_name, section, action, record_id, record_type, data_before, data_after, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())
	`

	result, err := r.db.Exec(query,
		history.UserName,
		history.Section,
		history.Action,
		history.RecordID,
		history.RecordType,
		history.DataBefore,
		history.DataAfter,
	)

	if err != nil {
		return fmt.Errorf("failed to create history: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	history.ID = id
	return nil
}

// GetAll retrieves all history records
func (r *MySQLHistoryRepository) GetAll() ([]model.History, error) {
	query := `
		SELECT id, user_name, section, action, record_id, record_type, data_before, data_after, created_at
		FROM travel_service_development.histories
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all histories: %w", err)
	}
	defer rows.Close()

	var histories []model.History
	for rows.Next() {
		var history model.History
		err := rows.Scan(
			&history.ID,
			&history.UserName,
			&history.Section,
			&history.Action,
			&history.RecordID,
			&history.RecordType,
			&history.DataBefore,
			&history.DataAfter,
			&history.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history: %w", err)
		}
		histories = append(histories, history)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return histories, nil
}

// GetByTableName retrieves all history records for a specific table (using section field)
func (r *MySQLHistoryRepository) GetByTableName(tableName string) ([]model.History, error) {
	query := `
		SELECT id, user_name, section, action, record_id, record_type, data_before, data_after, created_at
		FROM travel_service_development.histories
		WHERE section = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get histories by table name: %w", err)
	}
	defer rows.Close()

	var histories []model.History
	for rows.Next() {
		var history model.History
		err := rows.Scan(
			&history.ID,
			&history.UserName,
			&history.Section,
			&history.Action,
			&history.RecordID,
			&history.RecordType,
			&history.DataBefore,
			&history.DataAfter,
			&history.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history: %w", err)
		}
		histories = append(histories, history)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return histories, nil
}

// GetByRecordID retrieves all history records for a specific record
func (r *MySQLHistoryRepository) GetByRecordID(tableName string, recordID int64) ([]model.History, error) {
	query := `
		SELECT id, user_name, section, action, record_id, record_type, data_before, data_after, created_at
		FROM travel_service_development.histories
		WHERE section = ? AND record_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, tableName, recordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get histories by record id: %w", err)
	}
	defer rows.Close()

	var histories []model.History
	for rows.Next() {
		var history model.History
		err := rows.Scan(
			&history.ID,
			&history.UserName,
			&history.Section,
			&history.Action,
			&history.RecordID,
			&history.RecordType,
			&history.DataBefore,
			&history.DataAfter,
			&history.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history: %w", err)
		}
		histories = append(histories, history)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return histories, nil
}

// GetByID retrieves a history record by ID
func (r *MySQLHistoryRepository) GetByID(id int64) (*model.History, error) {
	query := `
		SELECT id, user_name, section, action, record_id, record_type, data_before, data_after, created_at
		FROM travel_service_development.histories
		WHERE id = ?
	`

	var history model.History
	err := r.db.QueryRow(query, id).Scan(
		&history.ID,
		&history.UserName,
		&history.Section,
		&history.Action,
		&history.RecordID,
		&history.RecordType,
		&history.DataBefore,
		&history.DataAfter,
		&history.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrHistoryNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get history by ID: %w", err)
	}

	return &history, nil
}

