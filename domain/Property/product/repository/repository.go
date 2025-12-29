package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"multi-onboarding/domain/Property/product"
	"multi-onboarding/utils"
)

type repository struct {
	mysqlSess *sql.DB
	draftData []product.ProductDraft
	draftMu   sync.RWMutex
}

// NewProductRepository creates a new product repository
func NewProductRepository() product.Repository {
	db := utils.GetDB()

	// Ensure template_drafts table exists
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS property_service_development.template_drafts (
		id VARCHAR(255) NOT NULL,
		locale VARCHAR(10) NOT NULL,
		value TEXT NOT NULL,
		product_code VARCHAR(255) DEFAULT NULL,
		insurance JSON DEFAULT NULL,
		created_by BIGINT NOT NULL DEFAULT 1,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id, locale),
		INDEX idx_product_code (product_code)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		fmt.Printf("Warning: Failed to create template_drafts table: %v\n", err)
	}

	return &repository{
		mysqlSess: db,
		draftData: make([]product.ProductDraft, 0),
	}
}

// Draft operations (in-memory)

func (r *repository) AddDraft(draft *product.ProductDraft) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	if draft.ID == 0 {
		if len(r.draftData) == 0 {
			draft.ID = 1
		} else {
			maxID := 0
			for _, d := range r.draftData {
				if d.ID > maxID {
					maxID = d.ID
				}
			}
			draft.ID = maxID + 1
		}
	}

	r.draftData = append(r.draftData, *draft)
	return nil
}

func (r *repository) GetAllDrafts() ([]product.ProductDraft, error) {
	r.draftMu.RLock()
	defer r.draftMu.RUnlock()

	result := make([]product.ProductDraft, len(r.draftData))
	copy(result, r.draftData)
	return result, nil
}

func (r *repository) GetDraftByID(id int) (*product.ProductDraft, error) {
	r.draftMu.RLock()
	defer r.draftMu.RUnlock()

	for i := range r.draftData {
		if r.draftData[i].ID == id {
			result := r.draftData[i]
			return &result, nil
		}
	}

	return nil, fmt.Errorf("draft not found")
}

func (r *repository) UpdateDraft(id int, draft *product.ProductDraft) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	for i := range r.draftData {
		if r.draftData[i].ID == id {
			draft.ID = id
			r.draftData[i] = *draft
			return nil
		}
	}

	return fmt.Errorf("draft not found")
}

func (r *repository) DeleteDraft(id int) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	for i, draft := range r.draftData {
		if draft.ID == id {
			r.draftData = append(r.draftData[:i], r.draftData[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("draft not found")
}

func (r *repository) ClearDrafts() error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	r.draftData = make([]product.ProductDraft, 0)
	return nil
}

func (r *repository) ConfirmDrafts(createdBy int64) error {
	r.draftMu.Lock()
	drafts := make([]product.ProductDraft, len(r.draftData))
	copy(drafts, r.draftData)
	r.draftMu.Unlock()

	// Convert each draft to Product and insert into database
	for _, draft := range drafts {
		isActive := 0
		if draft.IsActive == "1" || draft.IsActive == "true" {
			isActive = 1
		}

		p := &product.Product{
			Code:             draft.Code,
			InsuranceCode:    draft.InsuranceCode,
			Name:             draft.Name,
			Logo:             draft.Logo,
			IsActive:         isActive,
			InsuranceType:    draft.InsuranceType,
			AdminFee:         0, // AdminFee column doesn't exist in table
			PropertyCategory: draft.PropertyType,
			CoverageType:     draft.CoverageType,
			MinCoverage:      draft.MinCoverage,
			MaxCoverage:      draft.MaxCoverage,
			Summary:          draft.Summary,
			InsuranceDetail:  draft.InsuranceDetail,
			ProtectionDetail: draft.ProtectionDetail,
			HowToClaim:       draft.HowToClaim,
			CreatedBy:        createdBy,
		}

		if err := r.Create(p); err != nil {
			return fmt.Errorf("failed to create product from draft: %w", err)
		}
	}

	// Clear drafts after successful confirmation
	r.draftMu.Lock()
	r.draftData = make([]product.ProductDraft, 0)
	r.draftMu.Unlock()

	return nil
}

// Database operations

func (r *repository) GetAll(insuranceCode string) ([]product.Product, error) {
	var query string
	var args []interface{}

	// Query directly from hardcopy_admin_fee column
	if insuranceCode != "" {
		query = `
			SELECT id, code, insurance_code, name, logo, is_active, 
			       product_category, coverage_type,
			       hardcopy_admin_fee as admin_fee,
			       summary, insurance_detail, 
			       protection_detail, how_to_claim,
			       created_by, created_at, updated_by, updated_at, partner
			FROM property_service_development.products
			WHERE insurance_code = ?
			ORDER BY created_at DESC
		`
		args = []interface{}{insuranceCode}
	} else {
		query = `
			SELECT id, code, insurance_code, name, logo, is_active, 
			       product_category, coverage_type,
			       hardcopy_admin_fee as admin_fee,
			       summary, insurance_detail, 
			       protection_detail, how_to_claim,
			       created_by, created_at, updated_by, updated_at, partner
			FROM property_service_development.products
			ORDER BY created_at DESC
		`
		args = []interface{}{}
	}

	rows, err := r.mysqlSess.Query(query, args...)
	// Debug: log query error if any
	if err != nil {
		fmt.Printf("DEBUG GetAll: Query error: %v\n", err)
	}
	// If query fails due to missing hardcopy_admin_fee column, try with fallback to admin_fee
	if err != nil && strings.Contains(err.Error(), "Unknown column 'hardcopy_admin_fee'") {
		if insuranceCode != "" {
			query = `
				SELECT id, code, insurance_code, name, logo, is_active, 
				       product_category, coverage_type,
				       IFNULL(admin_fee, 0) as admin_fee,
				       summary, insurance_detail, 
				       protection_detail, how_to_claim,
				       created_by, created_at, updated_by, updated_at, partner
				FROM property_service_development.products
				WHERE insurance_code = ?
				ORDER BY created_at DESC
			`
			args = []interface{}{insuranceCode}
		} else {
			query = `
				SELECT id, code, insurance_code, name, logo, is_active, 
				       product_category, coverage_type,
				       IFNULL(admin_fee, 0) as admin_fee,
				       summary, insurance_detail, 
				       protection_detail, how_to_claim,
				       created_by, created_at, updated_by, updated_at, partner
				FROM property_service_development.products
				ORDER BY created_at DESC
			`
			args = []interface{}{}
		}
		rows, err = r.mysqlSess.Query(query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []product.Product
	for rows.Next() {
		var p product.Product
		var updatedBy sql.NullInt64
		var updatedAt sql.NullTime
		var partner sql.NullString
		var logo sql.NullString

		err := rows.Scan(
			&p.ID,
			&p.Code,
			&p.InsuranceCode,
			&p.Name,
			&logo,
			&p.IsActive,
			&p.PropertyCategory,
			&p.CoverageType,
			&p.AdminFee,
			&p.Summary,
			&p.InsuranceDetail,
			&p.ProtectionDetail,
			&p.HowToClaim,
			&p.CreatedBy,
			&p.CreatedAt,
			&updatedBy,
			&updatedAt,
			&partner,
		)
		if logo.Valid {
			p.Logo = logo.String
		}
		if err != nil {
			return nil, err
		}
		// Debug: log AdminFee value
		fmt.Printf("DEBUG GetAll: Product ID=%d, Code=%s, AdminFee=%d\n", p.ID, p.Code, p.AdminFee)
		// Set default values for columns that don't exist
		p.MinCoverage = 0
		p.MaxCoverage = 0
		// Admin_fee is already loaded from query (either from hardcopy_admin_fee, admin_fee, or 0)
		// Generate deduction_and_clause from code
		codeWithUnderscore := strings.ReplaceAll(p.Code, "-", "_")
		p.DeductibleAndClause = "DEDUCTION_AND_CLAUSE_" + codeWithUnderscore

		if updatedBy.Valid {
			p.UpdatedBy = &updatedBy.Int64
		}
		if updatedAt.Valid {
			p.UpdatedAt = &updatedAt.Time
		}
		if partner.Valid {
			p.Partner = partner.String
		}

		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *repository) GetByID(id int) (*product.Product, error) {
	query := `
		SELECT id, code, insurance_code, name, logo, is_active, 
		       product_category, coverage_type,
		       hardcopy_admin_fee as admin_fee,
		       summary, insurance_detail, 
		       protection_detail, how_to_claim,
		       created_by, created_at, updated_by, updated_at, partner
		FROM property_service_development.products
		WHERE id = ?
	`

	var p product.Product
	var updatedBy sql.NullInt64
	var updatedAt sql.NullTime
	var partner sql.NullString
	var logo sql.NullString

	err := r.mysqlSess.QueryRow(query, id).Scan(
		&p.ID,
		&p.Code,
		&p.InsuranceCode,
		&p.Name,
		&logo,
		&p.IsActive,
		&p.PropertyCategory,
		&p.CoverageType,
		&p.AdminFee,
		&p.Summary,
		&p.InsuranceDetail,
		&p.ProtectionDetail,
		&p.HowToClaim,
		&p.CreatedBy,
		&p.CreatedAt,
		&updatedBy,
		&updatedAt,
		&partner,
	)
	// Debug: log AdminFee value
	if err == nil {
		fmt.Printf("DEBUG GetByID: Product ID=%d, Code=%s, AdminFee=%d\n", p.ID, p.Code, p.AdminFee)
	}
	// If query fails due to missing hardcopy_admin_fee column, try with fallback
	if err != nil && strings.Contains(err.Error(), "Unknown column 'hardcopy_admin_fee'") {
		query = `
			SELECT id, code, insurance_code, name, logo, is_active, 
			       product_category, coverage_type,
			       IFNULL(admin_fee, 0) as admin_fee,
			       summary, insurance_detail, 
			       protection_detail, how_to_claim,
			       created_by, created_at, updated_by, updated_at, partner
			FROM property_service_development.products
			WHERE id = ?
		`
		err = r.mysqlSess.QueryRow(query, id).Scan(
			&p.ID,
			&p.Code,
			&p.InsuranceCode,
			&p.Name,
			&logo,
			&p.IsActive,
			&p.PropertyCategory,
			&p.CoverageType,
			&p.AdminFee,
			&p.Summary,
			&p.InsuranceDetail,
			&p.ProtectionDetail,
			&p.HowToClaim,
			&p.CreatedBy,
			&p.CreatedAt,
			&updatedBy,
			&updatedAt,
			&partner,
		)
		p.AdminFee = 0
	}
	if logo.Valid {
		p.Logo = logo.String
	}
	// Set default values for columns that don't exist
	p.MinCoverage = 0
	p.MaxCoverage = 0
	// Generate deduction_and_clause from code
	codeWithUnderscore := strings.ReplaceAll(p.Code, "-", "_")
	p.DeductibleAndClause = "DEDUCTION_AND_CLAUSE_" + codeWithUnderscore

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}
	if err != nil {
		return nil, err
	}

	if updatedBy.Valid {
		p.UpdatedBy = &updatedBy.Int64
	}
	if updatedAt.Valid {
		p.UpdatedAt = &updatedAt.Time
	}
	if partner.Valid {
		p.Partner = partner.String
	}

	return &p, nil
}

func (r *repository) Create(p *product.Product) error {
	// Generate deduction_and_clause from code: DEDUCTION_AND_CLAUSE_ + code (replace "-" with "_")
	codeWithUnderscore := strings.ReplaceAll(p.Code, "-", "_")
	deductionAndClause := "DEDUCTION_AND_CLAUSE_" + codeWithUnderscore

	// Ensure ID starts from latest insurance ID + 1 if products table is empty
	// Get latest ID from insurances table
	var latestInsuranceID int64
	err := r.mysqlSess.QueryRow("SELECT COALESCE(MAX(id), 0) FROM property_service_development.insurances").Scan(&latestInsuranceID)
	if err != nil {
		latestInsuranceID = 0
	}

	// Get latest ID from products table
	var latestProductID int64
	err = r.mysqlSess.QueryRow("SELECT COALESCE(MAX(id), 0) FROM property_service_development.products").Scan(&latestProductID)
	if err != nil {
		latestProductID = 0
	}

	// If products table is empty, set AUTO_INCREMENT to start from latestInsuranceID + 1
	if latestProductID == 0 && latestInsuranceID > 0 {
		// Set AUTO_INCREMENT to start from latestInsuranceID + 1
		setAutoIncrementQuery := fmt.Sprintf("ALTER TABLE property_service_development.products AUTO_INCREMENT = %d", latestInsuranceID+1)
		_, err = r.mysqlSess.Exec(setAutoIncrementQuery)
		if err != nil {
			// Log error but continue - MySQL will use default auto-increment
			fmt.Printf("Warning: Failed to set AUTO_INCREMENT: %v\n", err)
		}
	}

	// Try to insert with hardcopy_admin_fee first, fallback to admin_fee if column doesn't exist
	query := `
		INSERT INTO property_service_development.products 
		(code, insurance_code, name, logo, is_active, 
		 product_category, coverage_type,
		 hardcopy_admin_fee,
		 summary, insurance_detail, protection_detail, how_to_claim,
		 no_past_claim, construction_class_id, deduction_and_clause,
		 partner, created_by, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1, 1, ?, 'QOALAPLUS', 1, NOW())
	`

	result, err := r.mysqlSess.Exec(query,
		p.Code,
		p.InsuranceCode,
		p.Name,
		p.Logo,
		p.IsActive,
		p.PropertyCategory,
		p.CoverageType,
		p.AdminFee,
		p.Summary,
		p.InsuranceDetail,
		p.ProtectionDetail,
		p.HowToClaim,
		deductionAndClause,
	)

	// If hardcopy_admin_fee column doesn't exist, try with admin_fee
	if err != nil && strings.Contains(err.Error(), "Unknown column 'hardcopy_admin_fee'") {
		query = `
			INSERT INTO property_service_development.products 
			(code, insurance_code, name, logo, is_active, 
			 product_category, coverage_type,
			 admin_fee,
			 summary, insurance_detail, protection_detail, how_to_claim,
			 no_past_claim, construction_class_id, deduction_and_clause,
			 partner, created_by, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1, 1, ?, 'QOALAPLUS', 1, NOW())
		`
		result, err = r.mysqlSess.Exec(query,
			p.Code,
			p.InsuranceCode,
			p.Name,
			p.Logo,
			p.IsActive,
			p.PropertyCategory,
			p.CoverageType,
			p.AdminFee,
			p.Summary,
			p.InsuranceDetail,
			p.ProtectionDetail,
			p.HowToClaim,
			deductionAndClause,
		)
	}

	// If both columns don't exist, log warning but continue without admin_fee
	// Admin_fee will be set to 0 in the application layer
	if err != nil && (strings.Contains(err.Error(), "Unknown column 'hardcopy_admin_fee'") || strings.Contains(err.Error(), "Unknown column 'admin_fee'")) {
		fmt.Printf("Warning: Admin fee columns not found, inserting without admin_fee. Value: %d\n", p.AdminFee)
		query = `
			INSERT INTO property_service_development.products 
			(code, insurance_code, name, logo, is_active, 
			 product_category, coverage_type,
			 summary, insurance_detail, protection_detail, how_to_claim,
			 no_past_claim, construction_class_id, deduction_and_clause,
			 partner, created_by, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1, 1, ?, 'QOALAPLUS', 1, NOW())
		`
		result, err = r.mysqlSess.Exec(query,
			p.Code,
			p.InsuranceCode,
			p.Name,
			p.Logo,
			p.IsActive,
			p.PropertyCategory,
			p.CoverageType,
			p.Summary,
			p.InsuranceDetail,
			p.ProtectionDetail,
			p.HowToClaim,
			deductionAndClause,
		)
		// Set AdminFee to 0 since column doesn't exist
		p.AdminFee = 0
	}

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	p.ID = id
	p.CreatedAt = time.Now()

	return nil
}

func (r *repository) Update(id int, p *product.Product) error {
	// Get data before update
	beforeData, err := r.GetProductBeforeUpdate(id)
	if err != nil {
		beforeData = &product.ProductBeforeUpdate{
			Code:             p.Code,
			InsuranceCode:    p.InsuranceCode,
			Name:             p.Name,
			Logo:             p.Logo,
			IsActive:         p.IsActive,
			InsuranceType:    p.InsuranceType,
			AdminFee:         p.AdminFee,
			PropertyCategory: p.PropertyCategory,
			CoverageType:     p.CoverageType,
			MinCoverage:      p.MinCoverage,
			MaxCoverage:      p.MaxCoverage,
			Summary:          p.Summary,
			InsuranceDetail:  p.InsuranceDetail,
			ProtectionDetail: p.ProtectionDetail,
			HowToClaim:       p.HowToClaim,
		}
	}

	// Build changed fields
	changedBefore := make(map[string]interface{})
	changedAfter := make(map[string]interface{})
	updateFields := []string{}
	updateValues := []interface{}{}

	if beforeData.Code != p.Code {
		changedBefore["code"] = beforeData.Code
		changedAfter["code"] = p.Code
		updateFields = append(updateFields, "code = ?")
		updateValues = append(updateValues, p.Code)
	}
	if beforeData.Name != p.Name {
		changedBefore["name"] = beforeData.Name
		changedAfter["name"] = p.Name
		updateFields = append(updateFields, "name = ?")
		updateValues = append(updateValues, p.Name)
	}
	if beforeData.InsuranceCode != p.InsuranceCode {
		changedBefore["insurance_code"] = beforeData.InsuranceCode
		changedAfter["insurance_code"] = p.InsuranceCode
		updateFields = append(updateFields, "insurance_code = ?")
		updateValues = append(updateValues, p.InsuranceCode)
	}
	if beforeData.Logo != p.Logo {
		changedBefore["logo"] = beforeData.Logo
		changedAfter["logo"] = p.Logo
		updateFields = append(updateFields, "logo = ?")
		updateValues = append(updateValues, p.Logo)
	}
	// Type field removed - not in database schema
	if beforeData.IsActive != p.IsActive {
		changedBefore["is_active"] = beforeData.IsActive
		changedAfter["is_active"] = p.IsActive
		updateFields = append(updateFields, "is_active = ?")
		updateValues = append(updateValues, p.IsActive)
	}
	// Update admin_fee (try hardcopy_admin_fee first, fallback to admin_fee)
	if beforeData.AdminFee != p.AdminFee {
		changedBefore["admin_fee"] = beforeData.AdminFee
		changedAfter["admin_fee"] = p.AdminFee
		// Try hardcopy_admin_fee first
		updateFields = append(updateFields, "hardcopy_admin_fee = ?")
		updateValues = append(updateValues, p.AdminFee)
	}
	if beforeData.PropertyCategory != p.PropertyCategory {
		changedBefore["product_category"] = beforeData.PropertyCategory
		changedAfter["product_category"] = p.PropertyCategory
		updateFields = append(updateFields, "product_category = ?")
		updateValues = append(updateValues, p.PropertyCategory)
	}
	if beforeData.CoverageType != p.CoverageType {
		changedBefore["coverage_type"] = beforeData.CoverageType
		changedAfter["coverage_type"] = p.CoverageType
		updateFields = append(updateFields, "coverage_type = ?")
		updateValues = append(updateValues, p.CoverageType)
	}
	// MinCoverage and MaxCoverage columns don't exist in database, skip update
	if beforeData.Summary != p.Summary {
		changedBefore["summary"] = beforeData.Summary
		changedAfter["summary"] = p.Summary
		updateFields = append(updateFields, "summary = ?")
		updateValues = append(updateValues, p.Summary)
	}
	if beforeData.InsuranceDetail != p.InsuranceDetail {
		changedBefore["insurance_detail"] = beforeData.InsuranceDetail
		changedAfter["insurance_detail"] = p.InsuranceDetail
		updateFields = append(updateFields, "insurance_detail = ?")
		updateValues = append(updateValues, p.InsuranceDetail)
	}
	if beforeData.ProtectionDetail != p.ProtectionDetail {
		changedBefore["protection_detail"] = beforeData.ProtectionDetail
		changedAfter["protection_detail"] = p.ProtectionDetail
		updateFields = append(updateFields, "protection_detail = ?")
		updateValues = append(updateValues, p.ProtectionDetail)
	}
	if beforeData.HowToClaim != p.HowToClaim {
		changedBefore["how_to_claim"] = beforeData.HowToClaim
		changedAfter["how_to_claim"] = p.HowToClaim
		updateFields = append(updateFields, "how_to_claim = ?")
		updateValues = append(updateValues, p.HowToClaim)
	}

	if len(updateFields) == 0 {
		return nil // No changes
	}

	var updatedBy interface{}
	if p.UpdatedBy != nil {
		updatedBy = *p.UpdatedBy
		updateFields = append(updateFields, "updated_by = ?")
		updateValues = append(updateValues, updatedBy)
	}

	updateFields = append(updateFields, "updated_at = NOW()")
	updateValues = append(updateValues, id)

	query := fmt.Sprintf(`
		UPDATE property_service_development.products
		SET %s
		WHERE id = ?
	`, strings.Join(updateFields, ", "))

	result, err := r.mysqlSess.Exec(query, updateValues...)
	if err != nil {
		// If hardcopy_admin_fee column doesn't exist, try with admin_fee
		if strings.Contains(err.Error(), "Unknown column 'hardcopy_admin_fee'") {
			// Replace hardcopy_admin_fee with admin_fee in updateFields
			for i, field := range updateFields {
				if strings.Contains(field, "hardcopy_admin_fee") {
					updateFields[i] = strings.Replace(field, "hardcopy_admin_fee", "admin_fee", 1)
				}
			}
			query = fmt.Sprintf(`
				UPDATE property_service_development.products
				SET %s
				WHERE id = ?
			`, strings.Join(updateFields, ", "))
			result, err = r.mysqlSess.Exec(query, updateValues...)
			if err != nil {
				// If admin_fee also doesn't exist, remove it from update
				if strings.Contains(err.Error(), "Unknown column 'admin_fee'") {
					// Remove admin_fee from updateFields
					newUpdateFields := []string{}
					for _, field := range updateFields {
						if !strings.Contains(field, "admin_fee") {
							newUpdateFields = append(newUpdateFields, field)
						}
					}
					// Also remove corresponding value from updateValues
					adminFeeIndex := -1
					for i, field := range updateFields {
						if strings.Contains(field, "admin_fee") {
							adminFeeIndex = i
							break
						}
					}
					if adminFeeIndex >= 0 {
						newUpdateValues := make([]interface{}, 0, len(updateValues)-1)
						for i, val := range updateValues {
							if i != adminFeeIndex {
								newUpdateValues = append(newUpdateValues, val)
							}
						}
						updateValues = newUpdateValues
					}
					updateFields = newUpdateFields
					query = fmt.Sprintf(`
						UPDATE property_service_development.products
						SET %s
						WHERE id = ?
					`, strings.Join(updateFields, ", "))
					result, err = r.mysqlSess.Exec(query, updateValues...)
				}
			}
		}
		if err != nil {
			return fmt.Errorf("failed to update product: %w", err)
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	p.ID = int64(id)
	now := time.Now()
	p.UpdatedAt = &now

	// Save history if there are changes
	if len(changedBefore) > 0 {
		beforeJSON, _ := json.Marshal(changedBefore)
		afterJSON, _ := json.Marshal(changedAfter)
		_ = r.InsertHistory("system", "property_products", "UPDATE", id, "product", string(beforeJSON), string(afterJSON))
	}

	return nil
}

func (r *repository) Delete(id int) error {
	query := `DELETE FROM property_service_development.products WHERE id = ?`

	result, err := r.mysqlSess.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (r *repository) GetProductBeforeUpdate(id int) (*product.ProductBeforeUpdate, error) {
	query := `
		SELECT code, insurance_code, name, logo, is_active, 
		       product_category, coverage_type,
		       hardcopy_admin_fee as admin_fee,
		       summary, insurance_detail, protection_detail, how_to_claim
		FROM property_service_development.products
		WHERE id = ?
	`

	var before product.ProductBeforeUpdate
	var logo sql.NullString
	err := r.mysqlSess.QueryRow(query, id).Scan(
		&before.Code,
		&before.InsuranceCode,
		&before.Name,
		&logo,
		&before.IsActive,
		&before.PropertyCategory,
		&before.CoverageType,
		&before.AdminFee,
		&before.Summary,
		&before.InsuranceDetail,
		&before.ProtectionDetail,
		&before.HowToClaim,
	)
	// If query fails due to missing hardcopy_admin_fee column, try with fallback
	if err != nil && strings.Contains(err.Error(), "Unknown column 'hardcopy_admin_fee'") {
		query = `
			SELECT code, insurance_code, name, logo, is_active, 
			       product_category, coverage_type,
			       IFNULL(admin_fee, 0) as admin_fee,
			       summary, insurance_detail, protection_detail, how_to_claim
			FROM property_service_development.products
			WHERE id = ?
		`
		err = r.mysqlSess.QueryRow(query, id).Scan(
			&before.Code,
			&before.InsuranceCode,
			&before.Name,
			&logo,
			&before.IsActive,
			&before.PropertyCategory,
			&before.CoverageType,
			&before.AdminFee,
			&before.Summary,
			&before.InsuranceDetail,
			&before.ProtectionDetail,
			&before.HowToClaim,
		)
		before.AdminFee = 0
	}
	// Set default values for columns that don't exist
	before.MinCoverage = 0
	before.MaxCoverage = 0
	if err != nil {
		return nil, err
	}
	// Logo
	if logo.Valid {
		before.Logo = logo.String
	}

	return &before, nil
}

func (r *repository) InsertHistory(userName, section, action string, recordID int, recordType, dataBefore, dataAfter string) error {
	historyQuery := `INSERT INTO histories (user_name, section, action, record_id, record_type, data_before, data_after, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	_, err := r.mysqlSess.Exec(historyQuery, userName, section, action, recordID, recordType, dataBefore, dataAfter)
	if err != nil {
		return fmt.Errorf("failed to insert history: %w", err)
	}

	return nil
}

// GenerateTemplateDraftsForProducts generates template drafts for specific products
func (r *repository) GenerateTemplateDraftsForProducts(productCodes []string, createdBy int64) (int, error) {
	if len(productCodes) == 0 {
		return 0, fmt.Errorf("no product codes provided")
	}

	// Build query with IN clause
	placeholders := strings.Repeat("?,", len(productCodes))
	placeholders = placeholders[:len(placeholders)-1] // Remove trailing comma

	query := fmt.Sprintf(`
		SELECT id, code, insurance_code, name, logo, is_active, 
		       product_category, coverage_type,
		       IFNULL(hardcopy_admin_fee, 0) as admin_fee,
		       summary, insurance_detail, 
		       protection_detail, how_to_claim,
		       created_by, created_at, updated_by, updated_at, partner
		FROM property_service_development.products
		WHERE code IN (%s)
		ORDER BY created_at DESC
	`, placeholders)

	args := make([]interface{}, len(productCodes))
	for i, code := range productCodes {
		args[i] = code
	}

	rows, err := r.mysqlSess.Query(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	successCount := 0
	errorCount := 0
	for rows.Next() {
		var p product.Product
		var updatedBy sql.NullInt64
		var updatedAt sql.NullTime
		var partner sql.NullString
		var logo sql.NullString

		err := rows.Scan(
			&p.ID, &p.Code, &p.InsuranceCode, &p.Name, &logo,
			&p.IsActive, &p.PropertyCategory, &p.CoverageType,
			&p.AdminFee, &p.Summary, &p.InsuranceDetail, &p.ProtectionDetail,
			&p.HowToClaim, &p.CreatedBy, &p.CreatedAt, &updatedBy, &updatedAt, &partner,
		)
		if err != nil {
			fmt.Printf("Error scanning product: %v\n", err)
			errorCount++
			continue
		}

		if logo.Valid {
			p.Logo = logo.String
		}
		if updatedBy.Valid {
			p.UpdatedBy = &updatedBy.Int64
		}
		if partner.Valid {
			p.Partner = partner.String
		}

		// Generate template drafts for this product
		if err := r.createTemplateDraftsForProduct(
			p.Code,
			p.InsuranceCode,
			p.Summary,
			p.InsuranceDetail,
			createdBy,
		); err != nil {
			fmt.Printf("Error generating template drafts for product %s: %v\n", p.Code, err)
			errorCount++
			continue
		}
		successCount++
		fmt.Printf("Success: Generated template drafts for product %s\n", p.Code)
	}

	fmt.Printf("Template generation complete: %d successful, %d failed out of %d products\n", successCount, errorCount, len(productCodes))
	return successCount, nil
}

// createTemplateDraftsForProduct creates SUMMARY and INSURANCE_DETAIL template drafts for a product
func (r *repository) createTemplateDraftsForProduct(productCode, insuranceCode, summaryTemplateID, insuranceDetailTemplateID string, createdBy int64) error {
	if productCode == "" || insuranceCode == "" {
		return fmt.Errorf("product code and insurance code are required")
	}

	// Create SUMMARY template draft
	if summaryTemplateID != "" {
		summaryValue := `["SUMMARY1","SUMMARY2","SUMMARY3","SUMMARY4","SUMMARY5"]`
		summaryQuery := `
			INSERT INTO property_service_development.template_drafts (id, locale, value, product_code, insurance, created_by, created_at)
			VALUES (?, 'id', ?, ?, ?, ?, NOW())
			ON DUPLICATE KEY UPDATE
				value = VALUES(value),
				product_code = VALUES(product_code),
				insurance = VALUES(insurance),
				created_by = VALUES(created_by)
		`

		insuranceJSON, err := json.Marshal([]string{insuranceCode})
		if err != nil {
			return fmt.Errorf("failed to marshal insurance code: %w", err)
		}

		_, err = r.mysqlSess.Exec(summaryQuery, summaryTemplateID, summaryValue, productCode, insuranceJSON, createdBy)
		if err != nil {
			return fmt.Errorf("failed to create SUMMARY template draft: %w", err)
		}
	}

	// Create INSURANCE_DETAIL template draft
	if insuranceDetailTemplateID != "" {
		insuranceDetailValue := "" // Empty value, user will fill it later
		insuranceDetailQuery := `
			INSERT INTO property_service_development.template_drafts (id, locale, value, product_code, insurance, created_by, created_at)
			VALUES (?, 'id', ?, ?, ?, ?, NOW())
			ON DUPLICATE KEY UPDATE
				value = VALUES(value),
				product_code = VALUES(product_code),
				insurance = VALUES(insurance),
				created_by = VALUES(created_by)
		`

		insuranceJSON, err := json.Marshal([]string{insuranceCode})
		if err != nil {
			return fmt.Errorf("failed to marshal insurance code: %w", err)
		}

		_, err = r.mysqlSess.Exec(insuranceDetailQuery, insuranceDetailTemplateID, insuranceDetailValue, productCode, insuranceJSON, createdBy)
		if err != nil {
			return fmt.Errorf("failed to create INSURANCE_DETAIL template draft: %w", err)
		}
	}

	return nil
}
