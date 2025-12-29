package repository

import (
	"database/sql"
	"fmt"

	"multi-onboarding/domain/Travel/history"
	"multi-onboarding/utils"
)

type repository struct {
	mysqlSess *sql.DB
}

// NewHistoryRepository creates a new history repository
func NewHistoryRepository() history.Repository {
	return &repository{
		mysqlSess: utils.GetDB(),
	}
}

// GetAll retrieves all history records
func (r *repository) GetAll() ([]history.History, error) {
	query := `
		SELECT id, user_name, section, action, record_id, record_type, data_before, data_after, created_at
		FROM histories
		ORDER BY created_at DESC
	`

	rows, err := r.mysqlSess.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all histories: %w", err)
	}
	defer rows.Close()

	var histories []history.History
	for rows.Next() {
		var h history.History
		err := rows.Scan(
			&h.ID,
			&h.UserName,
			&h.Section,
			&h.Action,
			&h.RecordID,
			&h.RecordType,
			&h.DataBefore,
			&h.DataAfter,
			&h.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history: %w", err)
		}
		histories = append(histories, h)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return histories, nil
}

// GetByID retrieves a history record by ID
func (r *repository) GetByID(id int64) (*history.History, error) {
	query := `
		SELECT id, user_name, section, action, record_id, record_type, data_before, data_after, created_at
		FROM histories
		WHERE id = ?
	`

	var h history.History
	err := r.mysqlSess.QueryRow(query, id).Scan(
		&h.ID,
		&h.UserName,
		&h.Section,
		&h.Action,
		&h.RecordID,
		&h.RecordType,
		&h.DataBefore,
		&h.DataAfter,
		&h.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("history not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get history: %w", err)
	}

	return &h, nil
}

// GetByTableName retrieves history records by table name
func (r *repository) GetByTableName(tableName string) ([]history.History, error) {
	query := `
		SELECT id, user_name, section, action, record_id, record_type, data_before, data_after, created_at
		FROM histories
		WHERE section = ?
		ORDER BY created_at DESC
	`

	rows, err := r.mysqlSess.Query(query, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get histories by table name: %w", err)
	}
	defer rows.Close()

	var histories []history.History
	for rows.Next() {
		var h history.History
		err := rows.Scan(
			&h.ID,
			&h.UserName,
			&h.Section,
			&h.Action,
			&h.RecordID,
			&h.RecordType,
			&h.DataBefore,
			&h.DataAfter,
			&h.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history: %w", err)
		}
		histories = append(histories, h)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return histories, nil
}

// GetByRecordID retrieves history records by table name and record ID
func (r *repository) GetByRecordID(tableName string, recordID int) ([]history.History, error) {
	query := `
		SELECT id, user_name, section, action, record_id, record_type, data_before, data_after, created_at
		FROM histories
		WHERE section = ? AND record_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.mysqlSess.Query(query, tableName, recordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get histories by record ID: %w", err)
	}
	defer rows.Close()

	var histories []history.History
	for rows.Next() {
		var h history.History
		err := rows.Scan(
			&h.ID,
			&h.UserName,
			&h.Section,
			&h.Action,
			&h.RecordID,
			&h.RecordType,
			&h.DataBefore,
			&h.DataAfter,
			&h.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history: %w", err)
		}
		histories = append(histories, h)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return histories, nil
}

// Create creates a new history record
func (r *repository) Create(h *history.History) error {
	query := `
		INSERT INTO histories 
		(user_name, section, action, record_id, record_type, data_before, data_after, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())
	`

	result, err := r.mysqlSess.Exec(query,
		h.UserName,
		h.Section,
		h.Action,
		h.RecordID,
		h.RecordType,
		h.DataBefore,
		h.DataAfter,
	)

	if err != nil {
		return fmt.Errorf("failed to create history: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	h.ID = id
	return nil
}

