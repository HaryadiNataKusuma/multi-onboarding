package repository

import (
	"database/sql"
	"fmt"
	"time"

	"multi-onboarding/domain/Travel/addon"
	"multi-onboarding/utils"
)

type repository struct {
	mysqlSess *sql.DB
}

// NewAddonRepository creates a new addon repository
func NewAddonRepository() addon.Repository {
	return &repository{
		mysqlSess: utils.GetDB(),
	}
}

// GetAll retrieves all addons from database
func (r *repository) GetAll() ([]addon.Addon, error) {
	query := `
		SELECT id, code, name, name_my, name_en, is_active, created_by, created_at, updated_by, updated_at
		FROM addons
		ORDER BY created_at DESC
	`

	rows, err := r.mysqlSess.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addons []addon.Addon
	for rows.Next() {
		var a addon.Addon
		var updatedBy sql.NullInt64
		var updatedAt sql.NullTime

		err := rows.Scan(
			&a.ID,
			&a.Code,
			&a.Name,
			&a.NameMy,
			&a.NameEn,
			&a.IsActive,
			&a.CreatedBy,
			&a.CreatedAt,
			&updatedBy,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if updatedBy.Valid {
			a.UpdatedBy = &updatedBy.Int64
		}
		if updatedAt.Valid {
			a.UpdatedAt = &updatedAt.Time
		}

		addons = append(addons, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return addons, nil
}

// GetByID retrieves an addon by ID
func (r *repository) GetByID(id int64) (*addon.Addon, error) {
	query := `
		SELECT id, code, name, name_my, name_en, is_active, created_by, created_at, updated_by, updated_at
		FROM addons
		WHERE id = ?
	`

	var a addon.Addon
	var updatedBy sql.NullInt64
	var updatedAt sql.NullTime

	err := r.mysqlSess.QueryRow(query, id).Scan(
		&a.ID,
		&a.Code,
		&a.Name,
		&a.NameMy,
		&a.NameEn,
		&a.IsActive,
		&a.CreatedBy,
		&a.CreatedAt,
		&updatedBy,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("addon not found")
	}
	if err != nil {
		return nil, err
	}

	if updatedBy.Valid {
		a.UpdatedBy = &updatedBy.Int64
	}
	if updatedAt.Valid {
		a.UpdatedAt = &updatedAt.Time
	}

	return &a, nil
}

// GetByCode retrieves an addon by code
func (r *repository) GetByCode(code string) (*addon.Addon, error) {
	query := `
		SELECT id, code, name, name_my, name_en, is_active, created_by, created_at, updated_by, updated_at
		FROM addons
		WHERE code = ?
	`

	var a addon.Addon
	var updatedBy sql.NullInt64
	var updatedAt sql.NullTime

	err := r.mysqlSess.QueryRow(query, code).Scan(
		&a.ID,
		&a.Code,
		&a.Name,
		&a.NameMy,
		&a.NameEn,
		&a.IsActive,
		&a.CreatedBy,
		&a.CreatedAt,
		&updatedBy,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("addon not found")
	}
	if err != nil {
		return nil, err
	}

	if updatedBy.Valid {
		a.UpdatedBy = &updatedBy.Int64
	}
	if updatedAt.Valid {
		a.UpdatedAt = &updatedAt.Time
	}

	return &a, nil
}

// Create creates a new addon
func (r *repository) Create(a *addon.Addon) error {
	query := `
		INSERT INTO addons (code, name, name_my, name_en, is_active, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.mysqlSess.Exec(query,
		a.Code,
		a.Name,
		a.NameMy,
		a.NameEn,
		a.IsActive,
		a.CreatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create addon: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	a.ID = id
	a.CreatedAt = time.Now()
	now := time.Now()
	a.UpdatedAt = &now

	return nil
}

// Update updates an existing addon
func (r *repository) Update(a *addon.Addon) error {
	query := `
		UPDATE addons
		SET code = ?, name = ?, name_my = ?, name_en = ?, is_active = ?, updated_by = ?, updated_at = NOW()
		WHERE id = ?
	`

	var updatedBy interface{}
	if a.UpdatedBy != nil {
		updatedBy = *a.UpdatedBy
	} else {
		updatedBy = nil
	}

	result, err := r.mysqlSess.Exec(query,
		a.Code,
		a.Name,
		a.NameMy,
		a.NameEn,
		a.IsActive,
		updatedBy,
		a.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("addon not found")
	}

	now := time.Now()
	a.UpdatedAt = &now

	return nil
}

// Delete deletes an addon by ID
func (r *repository) Delete(id int64) error {
	query := `DELETE FROM addons WHERE id = ?`

	result, err := r.mysqlSess.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("addon not found")
	}

	return nil
}

// GetAddonBeforeUpdate retrieves addon data before update for history
func (r *repository) GetAddonBeforeUpdate(id int64) (*addon.AddonBeforeUpdate, error) {
	query := `SELECT code, name, name_my, name_en, is_active FROM addons WHERE id = ?`

	var before addon.AddonBeforeUpdate
	err := r.mysqlSess.QueryRow(query, id).Scan(
		&before.Code,
		&before.Name,
		&before.NameMy,
		&before.NameEn,
		&before.IsActive,
	)
	if err != nil {
		return nil, err
	}

	return &before, nil
}

// InsertHistory inserts a history record
func (r *repository) InsertHistory(userName, section, action string, recordID int64, recordType, dataBefore, dataAfter string) error {
	historyQuery := `INSERT INTO histories (user_name, section, action, record_id, record_type, data_before, data_after, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	_, err := r.mysqlSess.Exec(historyQuery, userName, section, action, recordID, recordType, dataBefore, dataAfter)
	if err != nil {
		return fmt.Errorf("failed to insert history: %w", err)
	}

	return nil
}

