package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"multi-onboarding/domain/Travel/product"
	"multi-onboarding/utils"
)

type repository struct {
	mysqlSess *sql.DB
	draftData []product.ProductDraft
	draftMu   sync.RWMutex
}

// NewProductRepository creates a new product repository
func NewProductRepository() product.Repository {
	return &repository{
		mysqlSess: utils.GetDB(),
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

	// Get all regions to map region name to region_id
	regionsQuery := `SELECT id, name FROM regions`
	rows, err := r.mysqlSess.Query(regionsQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	regionMap := make(map[string]int)
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			continue
		}
		regionMap[name] = id
	}

	// Convert each draft to Product and insert into database
	for _, draft := range drafts {
		regionID := draft.RegionID
		if regionID == 0 && draft.Region != "" {
			regionID = regionMap[draft.Region]
		}

		isActive := 0
		if draft.IsActive == "1" || draft.IsActive == "true" {
			isActive = 1
		}

		p := &product.Product{
			Code:             draft.Code,
			InsuranceCode:   draft.InsuranceCode,
			Name:            draft.Name,
			Logo:            draft.Logo,
			Type:            draft.Type,
			IsActive:        isActive,
			InsuranceType:   draft.InsuranceType,
			AdminFee:        0,
			RegionID:        regionID,
			SchengenEligible: 0,
			MinAdult:        draft.MinAdult,
			MaxAdult:        draft.MaxAdult,
			MinChild:        draft.MinChild,
			MaxChild:        draft.MaxChild,
			Summary:         draft.Summary,
			InsuranceDetail: draft.InsuranceDetail,
			ProtectionDetail: draft.ProtectionDetail,
			HowToClaim:      draft.HowToClaim,
			CreatedBy:       createdBy,
		}

		if err := r.Create(p); err != nil {
			return fmt.Errorf("failed to create product from draft: %w", err)
		}

		// Auto-create template drafts for SUMMARY and INSURANCE_DETAIL
		if err := r.createTemplateDraftsForProduct(draft.Code, draft.InsuranceCode, draft.Summary, draft.InsuranceDetail, createdBy); err != nil {
			// Log error but don't fail the confirmation
			// This allows products to be confirmed even if template creation fails
			fmt.Printf("Warning: Failed to create template drafts for product %s: %v\n", draft.Code, err)
		}
	}

	// Clear drafts after successful confirmation
	r.draftMu.Lock()
	r.draftData = make([]product.ProductDraft, 0)
	r.draftMu.Unlock()

	return nil
}

// createTemplateDraftsForProduct creates SUMMARY and INSURANCE_DETAIL template drafts
func (r *repository) createTemplateDraftsForProduct(productCode, insuranceCode, summaryID, insuranceDetailID string, createdBy int64) error {
	if productCode == "" || insuranceCode == "" {
		return fmt.Errorf("product code and insurance code are required")
	}

	// Create SUMMARY template draft
	if summaryID != "" {
		summaryValue := `["SUMMARY1","SUMMARY2","SUMMARY3","SUMMARY4","SUMMARY5"]`
		summaryQuery := `
			INSERT INTO travel_service_development.template_drafts (id, locale, value, product_code, insurance, created_by, created_at)
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
		
		_, err = r.mysqlSess.Exec(summaryQuery, summaryID, summaryValue, productCode, insuranceJSON, createdBy)
		if err != nil {
			return fmt.Errorf("failed to create SUMMARY template draft: %w", err)
		}
	}

	// Create INSURANCE_DETAIL template draft
	if insuranceDetailID != "" {
		insuranceDetailValue := "" // Empty value, user will fill it later
		insuranceDetailQuery := `
			INSERT INTO travel_service_development.template_drafts (id, locale, value, product_code, insurance, created_by, created_at)
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
		
		_, err = r.mysqlSess.Exec(insuranceDetailQuery, insuranceDetailID, insuranceDetailValue, productCode, insuranceJSON, createdBy)
		if err != nil {
			return fmt.Errorf("failed to create INSURANCE_DETAIL template draft: %w", err)
		}
	}

	return nil
}

// GenerateTemplateDraftsForAllProducts generates template drafts for all existing products
func (r *repository) GenerateTemplateDraftsForAllProducts(createdBy int64) (int, error) {
	// Get all products from database
	products, err := r.GetAll("")
	if err != nil {
		return 0, fmt.Errorf("failed to get products: %w", err)
	}

	successCount := 0
	errorCount := 0
	for _, product := range products {
		// Check if product has required fields
		if product.Code == "" || product.InsuranceCode == "" {
			fmt.Printf("Warning: Product missing code or insurance_code, skipping\n")
			errorCount++
			continue
		}
		
		if product.Summary == "" && product.InsuranceDetail == "" {
			fmt.Printf("Warning: Product %s has no summary or insurance_detail, skipping template generation\n", product.Code)
			errorCount++
			continue
		}
		
		// Generate template drafts for each product
		if err := r.createTemplateDraftsForProduct(
			product.Code,
			product.InsuranceCode,
			product.Summary,
			product.InsuranceDetail,
			createdBy,
		); err != nil {
			// Log error but continue with other products
			fmt.Printf("Error: Failed to generate template drafts for product %s (summary: %s, insurance_detail: %s): %v\n", 
				product.Code, product.Summary, product.InsuranceDetail, err)
			errorCount++
			continue
		}
		successCount++
		fmt.Printf("Success: Generated template drafts for product %s\n", product.Code)
	}

	fmt.Printf("Template generation complete: %d successful, %d failed out of %d products\n", successCount, errorCount, len(products))
	return successCount, nil
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
		SELECT id, code, insurance_code, name, logo, type, is_active, 
		       insurance_type, admin_fee, region_id, schengen_eligible, 
		       min_adult, max_adult, min_child, max_child, 
		       summary, insurance_detail, protection_detail, how_to_claim,
		       created_by, created_at, updated_by, updated_at, partner
		FROM products
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

		err := rows.Scan(
			&p.ID, &p.Code, &p.InsuranceCode, &p.Name, &p.Logo, &p.Type,
			&p.IsActive, &p.InsuranceType, &p.AdminFee, &p.RegionID,
			&p.SchengenEligible, &p.MinAdult, &p.MaxAdult, &p.MinChild,
			&p.MaxChild, &p.Summary, &p.InsuranceDetail, &p.ProtectionDetail,
			&p.HowToClaim, &p.CreatedBy, &p.CreatedAt, &updatedBy, &updatedAt, &p.Partner,
		)
		if err != nil {
			errorCount++
			continue
		}

		if updatedBy.Valid {
			p.UpdatedBy = &updatedBy.Int64
		}
		if updatedAt.Valid {
			p.UpdatedAt = &updatedAt.Time
		}

		// Check if product has required fields
		if p.Code == "" || p.InsuranceCode == "" {
			fmt.Printf("Warning: Product missing code or insurance_code, skipping\n")
			errorCount++
			continue
		}
		
		if p.Summary == "" && p.InsuranceDetail == "" {
			fmt.Printf("Warning: Product %s has no summary or insurance_detail, skipping template generation\n", p.Code)
			errorCount++
			continue
		}

		// Generate template drafts for each product
		if err := r.createTemplateDraftsForProduct(
			p.Code,
			p.InsuranceCode,
			p.Summary,
			p.InsuranceDetail,
			createdBy,
		); err != nil {
			fmt.Printf("Error: Failed to generate template drafts for product %s: %v\n", p.Code, err)
			errorCount++
			continue
		}
		successCount++
		fmt.Printf("Success: Generated template drafts for product %s\n", p.Code)
	}

	fmt.Printf("Template generation complete: %d successful, %d failed out of %d requested products\n", successCount, errorCount, len(productCodes))
	return successCount, nil
}

// Database operations

func (r *repository) GetAll(insuranceCode string) ([]product.Product, error) {
	var query string
	var args []interface{}

	if insuranceCode != "" {
		query = `
			SELECT id, code, insurance_code, name, logo, type, is_active, 
			       insurance_type, admin_fee, region_id, schengen_eligible, 
			       min_adult, max_adult, min_child, max_child, 
			       summary, insurance_detail, protection_detail, how_to_claim,
			       created_by, created_at, updated_by, updated_at, partner
			FROM products
			WHERE insurance_code = ?
			ORDER BY created_at DESC
		`
		args = []interface{}{insuranceCode}
	} else {
		query = `
			SELECT id, code, insurance_code, name, logo, type, is_active, 
			       insurance_type, admin_fee, region_id, schengen_eligible, 
			       min_adult, max_adult, min_child, max_child, 
			       summary, insurance_detail, protection_detail, how_to_claim,
			       created_by, created_at, updated_by, updated_at, partner
			FROM products
			ORDER BY created_at DESC
		`
		args = []interface{}{}
	}

	rows, err := r.mysqlSess.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []product.Product
	for rows.Next() {
		var p product.Product
		var updatedBy sql.NullInt64
		var updatedAt sql.NullTime

		err := rows.Scan(
			&p.ID,
			&p.Code,
			&p.InsuranceCode,
			&p.Name,
			&p.Logo,
			&p.Type,
			&p.IsActive,
			&p.InsuranceType,
			&p.AdminFee,
			&p.RegionID,
			&p.SchengenEligible,
			&p.MinAdult,
			&p.MaxAdult,
			&p.MinChild,
			&p.MaxChild,
			&p.Summary,
			&p.InsuranceDetail,
			&p.ProtectionDetail,
			&p.HowToClaim,
			&p.CreatedBy,
			&p.CreatedAt,
			&updatedBy,
			&updatedAt,
			&p.Partner,
		)
		if err != nil {
			return nil, err
		}

		if updatedBy.Valid {
			p.UpdatedBy = &updatedBy.Int64
		}
		if updatedAt.Valid {
			p.UpdatedAt = &updatedAt.Time
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
		SELECT id, code, insurance_code, name, logo, type, is_active, 
		       insurance_type, admin_fee, region_id, schengen_eligible, 
		       min_adult, max_adult, min_child, max_child, 
		       summary, insurance_detail, protection_detail, how_to_claim,
		       created_by, created_at, updated_by, updated_at, partner
		FROM products
		WHERE id = ?
	`

	var p product.Product
	var updatedBy sql.NullInt64
	var updatedAt sql.NullTime

	err := r.mysqlSess.QueryRow(query, id).Scan(
		&p.ID,
		&p.Code,
		&p.InsuranceCode,
		&p.Name,
		&p.Logo,
		&p.Type,
		&p.IsActive,
		&p.InsuranceType,
		&p.AdminFee,
		&p.RegionID,
		&p.SchengenEligible,
		&p.MinAdult,
		&p.MaxAdult,
		&p.MinChild,
		&p.MaxChild,
		&p.Summary,
		&p.InsuranceDetail,
		&p.ProtectionDetail,
		&p.HowToClaim,
		&p.CreatedBy,
		&p.CreatedAt,
		&updatedBy,
		&updatedAt,
		&p.Partner,
	)

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

	return &p, nil
}

func (r *repository) Create(p *product.Product) error {
	query := `
		INSERT INTO products 
		(code, insurance_code, name, logo, type, is_active, insurance_type, admin_fee, 
		 region_id, schengen_eligible, min_adult, max_adult, min_child, max_child, 
		 summary, insurance_detail, protection_detail, how_to_claim, created_by, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
	`

	result, err := r.mysqlSess.Exec(query,
		p.Code,
		p.InsuranceCode,
		p.Name,
		p.Logo,
		p.Type,
		p.IsActive,
		p.InsuranceType,
		p.AdminFee,
		p.RegionID,
		p.SchengenEligible,
		p.MinAdult,
		p.MaxAdult,
		p.MinChild,
		p.MaxChild,
		p.Summary,
		p.InsuranceDetail,
		p.ProtectionDetail,
		p.HowToClaim,
		p.CreatedBy,
	)

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
			Type:             p.Type,
			IsActive:         p.IsActive,
			InsuranceType:    p.InsuranceType,
			AdminFee:         p.AdminFee,
			RegionID:         p.RegionID,
			SchengenEligible: p.SchengenEligible,
			MinAdult:         p.MinAdult,
			MaxAdult:         p.MaxAdult,
			MinChild:         p.MinChild,
			MaxChild:         p.MaxChild,
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
	if beforeData.Type != p.Type {
		changedBefore["type"] = beforeData.Type
		changedAfter["type"] = p.Type
		updateFields = append(updateFields, "type = ?")
		updateValues = append(updateValues, p.Type)
	}
	if beforeData.IsActive != p.IsActive {
		changedBefore["is_active"] = beforeData.IsActive
		changedAfter["is_active"] = p.IsActive
		updateFields = append(updateFields, "is_active = ?")
		updateValues = append(updateValues, p.IsActive)
	}
	if beforeData.InsuranceType != p.InsuranceType {
		changedBefore["insurance_type"] = beforeData.InsuranceType
		changedAfter["insurance_type"] = p.InsuranceType
		updateFields = append(updateFields, "insurance_type = ?")
		updateValues = append(updateValues, p.InsuranceType)
	}
	if beforeData.RegionID != p.RegionID {
		changedBefore["region_id"] = beforeData.RegionID
		changedAfter["region_id"] = p.RegionID
		updateFields = append(updateFields, "region_id = ?")
		updateValues = append(updateValues, p.RegionID)
	}
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
		UPDATE products
		SET %s
		WHERE id = ?
	`, strings.Join(updateFields, ", "))

	result, err := r.mysqlSess.Exec(query, updateValues...)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
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
		_ = r.InsertHistory("system", "products", "UPDATE", id, "product", string(beforeJSON), string(afterJSON))
	}

	return nil
}

func (r *repository) Delete(id int) error {
	query := `DELETE FROM products WHERE id = ?`

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
		SELECT code, insurance_code, name, type, is_active, insurance_type, admin_fee, 
		       region_id, schengen_eligible, min_adult, max_adult, min_child, max_child,
		       summary, insurance_detail, protection_detail, how_to_claim
		FROM products
		WHERE id = ?
	`

	var before product.ProductBeforeUpdate
	err := r.mysqlSess.QueryRow(query, id).Scan(
		&before.Code,
		&before.InsuranceCode,
		&before.Name,
		&before.Type,
		&before.IsActive,
		&before.InsuranceType,
		&before.AdminFee,
		&before.RegionID,
		&before.SchengenEligible,
		&before.MinAdult,
		&before.MaxAdult,
		&before.MinChild,
		&before.MaxChild,
		&before.Summary,
		&before.InsuranceDetail,
		&before.ProtectionDetail,
		&before.HowToClaim,
	)
	if err != nil {
		return nil, err
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

