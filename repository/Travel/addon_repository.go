package repository

import (
	"database/sql"
	"fmt"
	"time"

	"multi-onboarding/model/Travel"
)

// AddonRepository interface defines the contract for addon data operations
type AddonRepository interface {
	GetAll() ([]model.Addon, error)
	GetByID(id int64) (*model.Addon, error)
	GetByCode(code string) (*model.Addon, error)
	Create(addon *model.Addon) error
	Update(addon *model.Addon) error
	Delete(id int64) error
}

// MySQLAddonRepository implements AddonRepository using MySQL database
type MySQLAddonRepository struct {
	db *sql.DB
}

// NewMySQLAddonRepository creates a new MySQL addon repository
func NewMySQLAddonRepository(db *sql.DB) *MySQLAddonRepository {
	return &MySQLAddonRepository{
		db: db,
	}
}

// GetAll retrieves all addons
func (r *MySQLAddonRepository) GetAll() ([]model.Addon, error) {
	query := `
		SELECT id, code, name, name_my, name_en, is_active, created_by, created_at, updated_by, updated_at
		FROM travel_service_development.addons
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addons []model.Addon
	for rows.Next() {
		var addon model.Addon
		var updatedBy sql.NullInt64
		var updatedAt sql.NullTime
		
		err := rows.Scan(
			&addon.ID,
			&addon.Code,
			&addon.Name,
			&addon.NameMy,
			&addon.NameEn,
			&addon.IsActive,
			&addon.CreatedBy,
			&addon.CreatedAt,
			&updatedBy,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if updatedBy.Valid {
			addon.UpdatedBy = &updatedBy.Int64
		}
		if updatedAt.Valid {
			addon.UpdatedAt = &updatedAt.Time
		}

		addons = append(addons, addon)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return addons, nil
}

// GetByID retrieves an addon by ID
func (r *MySQLAddonRepository) GetByID(id int64) (*model.Addon, error) {
	query := `
		SELECT id, code, name, name_my, name_en, is_active, created_by, created_at, updated_by, updated_at
		FROM travel_service_development.addons
		WHERE id = ?
	`

	var addon model.Addon
	var updatedBy sql.NullInt64
	var updatedAt sql.NullTime
	
	err := r.db.QueryRow(query, id).Scan(
		&addon.ID,
		&addon.Code,
		&addon.Name,
		&addon.NameMy,
		&addon.NameEn,
		&addon.IsActive,
		&addon.CreatedBy,
		&addon.CreatedAt,
		&updatedBy,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrAddonNotFound
	}
	if err != nil {
		return nil, err
	}

	if updatedBy.Valid {
		addon.UpdatedBy = &updatedBy.Int64
	}
	if updatedAt.Valid {
		addon.UpdatedAt = &updatedAt.Time
	}

	return &addon, nil
}

// GetByCode retrieves an addon by code
func (r *MySQLAddonRepository) GetByCode(code string) (*model.Addon, error) {
	query := `
		SELECT id, code, name, name_my, name_en, is_active, created_by, created_at, updated_by, updated_at
		FROM travel_service_development.addons
		WHERE code = ?
	`

	var addon model.Addon
	var updatedBy sql.NullInt64
	var updatedAt sql.NullTime
	
	err := r.db.QueryRow(query, code).Scan(
		&addon.ID,
		&addon.Code,
		&addon.Name,
		&addon.NameMy,
		&addon.NameEn,
		&addon.IsActive,
		&addon.CreatedBy,
		&addon.CreatedAt,
		&updatedBy,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrAddonNotFound
	}
	if err != nil {
		return nil, err
	}

	if updatedBy.Valid {
		addon.UpdatedBy = &updatedBy.Int64
	}
	if updatedAt.Valid {
		addon.UpdatedAt = &updatedAt.Time
	}

	return &addon, nil
}

// Create creates a new addon (only if not exists)
func (r *MySQLAddonRepository) Create(addon *model.Addon) error {
	query := `
		INSERT INTO travel_service_development.addons (code, name, name_my, name_en, is_active, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.Exec(query,
		addon.Code,
		addon.Name,
		addon.NameMy,
		addon.NameEn,
		addon.IsActive,
		addon.CreatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create addon: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	addon.ID = id
	addon.CreatedAt = time.Now()
	now := time.Now()
	addon.UpdatedAt = &now

	return nil
}

// Update updates an existing addon
func (r *MySQLAddonRepository) Update(addon *model.Addon) error {
	query := `
		UPDATE travel_service_development.addons
		SET code = ?, name = ?, name_my = ?, name_en = ?, is_active = ?, updated_by = ?, updated_at = NOW()
		WHERE id = ?
	`

	var updatedBy interface{}
	if addon.UpdatedBy != nil {
		updatedBy = *addon.UpdatedBy
	} else {
		updatedBy = nil
	}

	result, err := r.db.Exec(query,
		addon.Code,
		addon.Name,
		addon.NameMy,
		addon.NameEn,
		addon.IsActive,
		updatedBy,
		addon.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrAddonNotFound
	}

	now := time.Now()
	addon.UpdatedAt = &now

	return nil
}

// Delete deletes an addon by ID
func (r *MySQLAddonRepository) Delete(id int64) error {
	query := `DELETE FROM travel_service_development.addons WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrAddonNotFound
	}

	return nil
}
