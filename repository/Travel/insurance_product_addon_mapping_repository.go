package repository

import (
	"database/sql"
	"fmt"
	"multi-onboarding/model/Travel"
)

// InsuranceProductAddonMappingRepository defines the interface for insurance product addon mapping operations
type InsuranceProductAddonMappingRepository interface {
	Create(mapping *model.InsuranceProductAddonMapping) error
	CreateBatch(mappings []model.InsuranceProductAddonMapping) error
	DeleteByProductCode(productCode string) error
}

// MySQLInsuranceProductAddonMappingRepository implements InsuranceProductAddonMappingRepository using MySQL database
type MySQLInsuranceProductAddonMappingRepository struct {
	db *sql.DB
}

// NewMySQLInsuranceProductAddonMappingRepository creates a new MySQL insurance product addon mapping repository
func NewMySQLInsuranceProductAddonMappingRepository(db *sql.DB) *MySQLInsuranceProductAddonMappingRepository {
	return &MySQLInsuranceProductAddonMappingRepository{
		db: db,
	}
}

// Create creates a new insurance product addon mapping
func (r *MySQLInsuranceProductAddonMappingRepository) Create(mapping *model.InsuranceProductAddonMapping) error {
	query := `
		INSERT INTO travel_service_development.insurance_product_addon_mappings 
		(insurance_code, product_code, addon_code, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE updated_at = NOW()
	`

	result, err := r.db.Exec(query,
		mapping.InsuranceCode,
		mapping.ProductCode,
		mapping.AddonCode,
		mapping.CreatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create insurance product addon mapping: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	mapping.ID = id
	return nil
}

// CreateBatch creates multiple insurance product addon mappings in a single transaction
func (r *MySQLInsuranceProductAddonMappingRepository) CreateBatch(mappings []model.InsuranceProductAddonMapping) error {
	if len(mappings) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO travel_service_development.insurance_product_addon_mappings 
		(insurance_code, product_code, addon_code, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE updated_at = NOW()
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for i := range mappings {
		result, err := stmt.Exec(
			mappings[i].InsuranceCode,
			mappings[i].ProductCode,
			mappings[i].AddonCode,
			mappings[i].CreatedBy,
		)
		if err != nil {
			return fmt.Errorf("failed to insert insurance product addon mapping: %w", err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id: %w", err)
		}

		mappings[i].ID = id
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DeleteByProductCode deletes all insurance product addon mappings for a product code
func (r *MySQLInsuranceProductAddonMappingRepository) DeleteByProductCode(productCode string) error {
	query := `DELETE FROM travel_service_development.insurance_product_addon_mappings WHERE product_code = ?`

	_, err := r.db.Exec(query, productCode)
	if err != nil {
		return fmt.Errorf("failed to delete insurance product addon mappings by product code: %w", err)
	}

	return nil
}

