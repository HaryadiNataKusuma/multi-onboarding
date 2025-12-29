package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"multi-onboarding/model/Travel"
)

// MySQLProductRepository implements ProductRepository using MySQL database
type MySQLProductRepository struct {
	db         *sql.DB
	historyRepo HistoryRepository
	draftData []model.ProductDraft
	draftMu   sync.RWMutex
}

// NewMySQLProductRepository creates a new MySQL product repository
func NewMySQLProductRepository(db *sql.DB) *MySQLProductRepository {
	return &MySQLProductRepository{
		db:        db,
		draftData: make([]model.ProductDraft, 0),
	}
}

// SetHistoryRepository sets the history repository for saving update histories
func (r *MySQLProductRepository) SetHistoryRepository(historyRepo HistoryRepository) {
	r.historyRepo = historyRepo
}

// Draft operations (in-memory)

// AddDraft adds a draft product
func (r *MySQLProductRepository) AddDraft(draft *model.ProductDraft) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	if draft.ID == 0 {
		// Auto-generate ID if not provided
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

// GetAllDrafts gets all draft products
func (r *MySQLProductRepository) GetAllDrafts() ([]model.ProductDraft, error) {
	r.draftMu.RLock()
	defer r.draftMu.RUnlock()

	result := make([]model.ProductDraft, len(r.draftData))
	copy(result, r.draftData)
	return result, nil
}

// ClearDrafts clears all draft products
func (r *MySQLProductRepository) ClearDrafts() error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	r.draftData = make([]model.ProductDraft, 0)
	return nil
}

// DeleteDraft deletes a draft product by ID
func (r *MySQLProductRepository) DeleteDraft(id int) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	for i, draft := range r.draftData {
		if draft.ID == id {
			r.draftData = append(r.draftData[:i], r.draftData[i+1:]...)
			return nil
		}
	}
	return ErrProductNotFound
}

// ConfirmDrafts confirms all draft products to database
func (r *MySQLProductRepository) ConfirmDrafts(createdBy int64) error {
	r.draftMu.Lock()
	drafts := make([]model.ProductDraft, len(r.draftData))
	copy(drafts, r.draftData)
	r.draftMu.Unlock()

	// Get all regions to map region name to region_id
	regionsQuery := `SELECT id, name FROM regions`
	rows, err := r.db.Query(regionsQuery)
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

	// Convert each draft to TravelProduct and insert into database
	for _, draft := range drafts {
		// Use region_id from draft if available, otherwise use regionMap
		regionID := draft.RegionID
		if regionID == 0 && draft.Region != "" {
			regionID = regionMap[draft.Region]
		}

		product := model.TravelProduct{
			Code:             draft.Code,
			InsuranceCode:   draft.InsuranceCode,
			Name:            draft.Name,
			Logo:            draft.Logo, // Use logo from draft
			Type:            draft.Type,
			IsActive:        convertStringToInt(draft.IsActive),
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

		if err := r.Create(&product); err != nil {
			return fmt.Errorf("failed to create product from draft: %w", err)
		}
	}

	// Clear drafts after successful confirmation
	r.draftMu.Lock()
	r.draftData = make([]model.ProductDraft, 0)
	r.draftMu.Unlock()

	return nil
}

// Helper function to convert string to int
func convertStringToInt(s string) int {
	if s == "1" || s == "true" {
		return 1
	}
	return 0
}

// GetAll retrieves all products
func (r *MySQLProductRepository) GetAll() ([]model.TravelProduct, error) {
	query := `
		SELECT id, code, insurance_code, name, logo, type, is_active, 
		       insurance_type, admin_fee, region_id, schengen_eligible, 
		       min_adult, max_adult, min_child, max_child, 
		       summary, insurance_detail, protection_detail, how_to_claim,
		       created_by, created_at, updated_by, updated_at, partner
		FROM travel_service_development.products
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.TravelProduct
	for rows.Next() {
		var p model.TravelProduct
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
			&p.UpdatedBy,
			&p.UpdatedAt,
			&p.Partner,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// GetByID retrieves a product by ID
func (r *MySQLProductRepository) GetByID(id int) (*model.TravelProduct, error) {
	query := `
		SELECT id, code, insurance_code, name, logo, type, is_active, 
		       insurance_type, admin_fee, region_id, schengen_eligible, 
		       min_adult, max_adult, min_child, max_child, 
		       summary, insurance_detail, protection_detail, how_to_claim,
		       created_by, created_at, updated_by, updated_at, partner
		FROM travel_service_development.products
		WHERE id = ?
	`

	var p model.TravelProduct
	err := r.db.QueryRow(query, id).Scan(
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
		&p.UpdatedBy,
		&p.UpdatedAt,
		&p.Partner,
	)

	if err == sql.ErrNoRows {
		return nil, ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// Create creates a new product
func (r *MySQLProductRepository) Create(product *model.TravelProduct) error {
	query := `
		INSERT INTO products 
		(code, insurance_code, name, logo, type, is_active, insurance_type, admin_fee, 
		 region_id, schengen_eligible, min_adult, max_adult, min_child, max_child, 
		 summary, insurance_detail, protection_detail, how_to_claim, created_by, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
	`
	
	result, err := r.db.Exec(query,
		product.Code,
		product.InsuranceCode,
		product.Name,
		product.Logo,
		product.Type,
		product.IsActive,
		product.InsuranceType,
		product.AdminFee,
		product.RegionID,
		product.SchengenEligible,
		product.MinAdult,
		product.MaxAdult,
		product.MinChild,
		product.MaxChild,
		product.Summary,
		product.InsuranceDetail,
		product.ProtectionDetail,
		product.HowToClaim,
		product.CreatedBy,
	)

	if err != nil {
		// Return the actual database error - log it for debugging
		// The error message will be passed through to the handler
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	product.ID = id
	product.CreatedAt = time.Now()
	now := time.Now()
	product.UpdatedAt = &now

	return nil
}

// Update updates an existing product
func (r *MySQLProductRepository) Update(product *model.TravelProduct) error {
	// Get data before update for comparison and history
	beforeProduct, err := r.GetByID(int(product.ID))
	if err != nil {
		return err
	}

	// Build changed fields map (only fields that are different)
	changedBefore := make(map[string]interface{})
	changedAfter := make(map[string]interface{})
	updateFields := []string{}
	updateValues := []interface{}{}

	// Compare each field and build update query dynamically
	if beforeProduct.Code != product.Code {
		changedBefore["code"] = beforeProduct.Code
		changedAfter["code"] = product.Code
		updateFields = append(updateFields, "code = ?")
		updateValues = append(updateValues, product.Code)
	}
	if beforeProduct.InsuranceCode != product.InsuranceCode {
		changedBefore["insurance_code"] = beforeProduct.InsuranceCode
		changedAfter["insurance_code"] = product.InsuranceCode
		updateFields = append(updateFields, "insurance_code = ?")
		updateValues = append(updateValues, product.InsuranceCode)
	}
	if beforeProduct.Name != product.Name {
		changedBefore["name"] = beforeProduct.Name
		changedAfter["name"] = product.Name
		updateFields = append(updateFields, "name = ?")
		updateValues = append(updateValues, product.Name)
	}
	if beforeProduct.Logo != product.Logo {
		changedBefore["logo"] = beforeProduct.Logo
		changedAfter["logo"] = product.Logo
		updateFields = append(updateFields, "logo = ?")
		updateValues = append(updateValues, product.Logo)
	}
	if beforeProduct.Type != product.Type {
		changedBefore["type"] = beforeProduct.Type
		changedAfter["type"] = product.Type
		updateFields = append(updateFields, "type = ?")
		updateValues = append(updateValues, product.Type)
	}
	if beforeProduct.IsActive != product.IsActive {
		changedBefore["is_active"] = beforeProduct.IsActive
		changedAfter["is_active"] = product.IsActive
		updateFields = append(updateFields, "is_active = ?")
		updateValues = append(updateValues, product.IsActive)
	}
	if beforeProduct.InsuranceType != product.InsuranceType {
		changedBefore["insurance_type"] = beforeProduct.InsuranceType
		changedAfter["insurance_type"] = product.InsuranceType
		updateFields = append(updateFields, "insurance_type = ?")
		updateValues = append(updateValues, product.InsuranceType)
	}
	if beforeProduct.AdminFee != product.AdminFee {
		changedBefore["admin_fee"] = beforeProduct.AdminFee
		changedAfter["admin_fee"] = product.AdminFee
		updateFields = append(updateFields, "admin_fee = ?")
		updateValues = append(updateValues, product.AdminFee)
	}
	if beforeProduct.RegionID != product.RegionID {
		changedBefore["region_id"] = beforeProduct.RegionID
		changedAfter["region_id"] = product.RegionID
		updateFields = append(updateFields, "region_id = ?")
		updateValues = append(updateValues, product.RegionID)
	}
	if beforeProduct.SchengenEligible != product.SchengenEligible {
		changedBefore["schengen_eligible"] = beforeProduct.SchengenEligible
		changedAfter["schengen_eligible"] = product.SchengenEligible
		updateFields = append(updateFields, "schengen_eligible = ?")
		updateValues = append(updateValues, product.SchengenEligible)
	}
	if beforeProduct.MinAdult != product.MinAdult {
		changedBefore["min_adult"] = beforeProduct.MinAdult
		changedAfter["min_adult"] = product.MinAdult
		updateFields = append(updateFields, "min_adult = ?")
		updateValues = append(updateValues, product.MinAdult)
	}
	if beforeProduct.MaxAdult != product.MaxAdult {
		changedBefore["max_adult"] = beforeProduct.MaxAdult
		changedAfter["max_adult"] = product.MaxAdult
		updateFields = append(updateFields, "max_adult = ?")
		updateValues = append(updateValues, product.MaxAdult)
	}
	if beforeProduct.MinChild != product.MinChild {
		changedBefore["min_child"] = beforeProduct.MinChild
		changedAfter["min_child"] = product.MinChild
		updateFields = append(updateFields, "min_child = ?")
		updateValues = append(updateValues, product.MinChild)
	}
	if beforeProduct.MaxChild != product.MaxChild {
		changedBefore["max_child"] = beforeProduct.MaxChild
		changedAfter["max_child"] = product.MaxChild
		updateFields = append(updateFields, "max_child = ?")
		updateValues = append(updateValues, product.MaxChild)
	}
	if beforeProduct.Summary != product.Summary {
		changedBefore["summary"] = beforeProduct.Summary
		changedAfter["summary"] = product.Summary
		updateFields = append(updateFields, "summary = ?")
		updateValues = append(updateValues, product.Summary)
	}
	if beforeProduct.InsuranceDetail != product.InsuranceDetail {
		changedBefore["insurance_detail"] = beforeProduct.InsuranceDetail
		changedAfter["insurance_detail"] = product.InsuranceDetail
		updateFields = append(updateFields, "insurance_detail = ?")
		updateValues = append(updateValues, product.InsuranceDetail)
	}
	if beforeProduct.ProtectionDetail != product.ProtectionDetail {
		changedBefore["protection_detail"] = beforeProduct.ProtectionDetail
		changedAfter["protection_detail"] = product.ProtectionDetail
		updateFields = append(updateFields, "protection_detail = ?")
		updateValues = append(updateValues, product.ProtectionDetail)
	}
	if beforeProduct.HowToClaim != product.HowToClaim {
		changedBefore["how_to_claim"] = beforeProduct.HowToClaim
		changedAfter["how_to_claim"] = product.HowToClaim
		updateFields = append(updateFields, "how_to_claim = ?")
		updateValues = append(updateValues, product.HowToClaim)
	}

	// If no fields changed, return early (no update needed)
	if len(updateFields) == 0 {
		return nil
	}

	// Always update updated_by and updated_at
	updateFields = append(updateFields, "updated_by = ?", "updated_at = NOW()")
	updateValues = append(updateValues, product.UpdatedBy)
	updateValues = append(updateValues, product.ID)

	// Build dynamic UPDATE query using strings.Join
	query := fmt.Sprintf(`
		UPDATE travel_service_development.products 
		SET %s
		WHERE id = ?
	`, strings.Join(updateFields, ", "))

	result, err := r.db.Exec(query, updateValues...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrProductNotFound
	}

	now := time.Now()
	product.UpdatedAt = &now

	// Save history only if there are changes and historyRepo is available
	if r.historyRepo != nil && len(changedBefore) > 0 {
		beforeJSON, err := json.Marshal(changedBefore)
		if err != nil {
			fmt.Printf("Warning: failed to marshal changed before data: %v\n", err)
		} else {
			afterJSON, err := json.Marshal(changedAfter)
			if err != nil {
				fmt.Printf("Warning: failed to marshal changed after data: %v\n", err)
			} else {
				history := &model.History{
					UserName:   "System",
					Section:    "products",
					Action:     "UPDATE",
					RecordID:   int(product.ID),
					RecordType: "product",
					DataBefore: string(beforeJSON),
					DataAfter:  string(afterJSON),
				}
				if err := r.historyRepo.Create(history); err != nil {
					// Log error but don't fail the update
					fmt.Printf("ERROR: failed to save history: %v\n", err)
				} else {
					fmt.Printf("DEBUG: History saved successfully for product ID %d with %d changed field(s)\n", product.ID, len(changedBefore))
				}
			}
		}
	}

	return nil
}

// Delete deletes a product by ID
func (r *MySQLProductRepository) Delete(id int) error {
	query := `DELETE FROM products WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrProductNotFound
	}

	return nil
}

