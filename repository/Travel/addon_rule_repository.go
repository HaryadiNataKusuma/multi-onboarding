package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"multi-onboarding/model/Travel"
)

// AddonRuleRepository interface defines the contract for addon rule data operations
type AddonRuleRepository interface {
	GetAll() ([]model.AddonRule, error)
	GetByID(id int64) (*model.AddonRule, error)
	GetByProductCode(productCode string) ([]model.AddonRule, error)
	GetByAddonCode(addonCode string) ([]model.AddonRule, error)
	Create(addonRule *model.AddonRule) error
	CreateBatch(addonRules []model.AddonRule) error
	Update(id int64, addonRule *model.AddonRule) error
	Delete(id int64) error
	DeleteByProductCodeAndAddonCode(productCode, addonCode string) error
	DeleteByProductCode(productCode string) error
	CleanupDuplicates() error
}

// MySQLAddonRuleRepository implements AddonRuleRepository using MySQL database
type MySQLAddonRuleRepository struct {
	db         *sql.DB
	historyRepo HistoryRepository
}

// NewMySQLAddonRuleRepository creates a new MySQL addon rule repository
func NewMySQLAddonRuleRepository(db *sql.DB) *MySQLAddonRuleRepository {
	return &MySQLAddonRuleRepository{
		db: db,
	}
}

// SetHistoryRepository sets the history repository for saving update histories
func (r *MySQLAddonRuleRepository) SetHistoryRepository(historyRepo HistoryRepository) {
	r.historyRepo = historyRepo
}

// GetAll retrieves all addon rules
func (r *MySQLAddonRuleRepository) GetAll() ([]model.AddonRule, error) {
	query := `
		SELECT id, product_code, addon_code, insurance_code, rules, created_by, created_at, updated_by, updated_at
		FROM travel_service_development.addon_rules
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addonRules []model.AddonRule
	for rows.Next() {
		var addonRule model.AddonRule
		var updatedBy sql.NullInt64
		var updatedAt sql.NullTime
		var insuranceCode sql.NullString
		var rules sql.NullString
		
		err := rows.Scan(
			&addonRule.ID,
			&addonRule.ProductCode,
			&addonRule.AddonCode,
			&insuranceCode,
			&rules,
			&addonRule.CreatedBy,
			&addonRule.CreatedAt,
			&updatedBy,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if insuranceCode.Valid {
			addonRule.InsuranceCode = &insuranceCode.String
		}
		if rules.Valid {
			addonRule.Rules = &rules.String
		}
		if updatedBy.Valid {
			addonRule.UpdatedBy = &updatedBy.Int64
		}
		if updatedAt.Valid {
			addonRule.UpdatedAt = &updatedAt.Time
		}

		addonRules = append(addonRules, addonRule)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return addonRules, nil
}

// GetByID retrieves an addon rule by ID
func (r *MySQLAddonRuleRepository) GetByID(id int64) (*model.AddonRule, error) {
	query := `
		SELECT id, product_code, addon_code, insurance_code, rules, created_by, created_at, updated_by, updated_at
		FROM travel_service_development.addon_rules
		WHERE id = ?
	`

	var addonRule model.AddonRule
	var updatedBy sql.NullInt64
	var updatedAt sql.NullTime
	var insuranceCode sql.NullString
	var rules sql.NullString
	
	err := r.db.QueryRow(query, id).Scan(
		&addonRule.ID,
		&addonRule.ProductCode,
		&addonRule.AddonCode,
		&insuranceCode,
		&rules,
		&addonRule.CreatedBy,
		&addonRule.CreatedAt,
		&updatedBy,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrAddonRuleNotFound
	}
	if err != nil {
		return nil, err
	}

	if insuranceCode.Valid {
		addonRule.InsuranceCode = &insuranceCode.String
	}
	if rules.Valid {
		addonRule.Rules = &rules.String
	}
	if updatedBy.Valid {
		addonRule.UpdatedBy = &updatedBy.Int64
	}
	if updatedAt.Valid {
		addonRule.UpdatedAt = &updatedAt.Time
	}

	return &addonRule, nil
}

// GetByProductCode retrieves addon rules by product code (including inactive ones with INACTIVE- prefix)
func (r *MySQLAddonRuleRepository) GetByProductCode(productCode string) ([]model.AddonRule, error) {
	query := `
		SELECT id, product_code, addon_code, created_by, created_at, updated_by, updated_at
		FROM travel_service_development.addon_rules
		WHERE product_code = ? OR product_code = ?
		ORDER BY created_at DESC
	`
	inactiveProductCode := "INACTIVE-" + productCode

	rows, err := r.db.Query(query, productCode, inactiveProductCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addonRules []model.AddonRule
	for rows.Next() {
		var addonRule model.AddonRule
		var updatedBy sql.NullInt64
		var updatedAt sql.NullTime
		
		err := rows.Scan(
			&addonRule.ID,
			&addonRule.ProductCode,
			&addonRule.AddonCode,
			&addonRule.CreatedBy,
			&addonRule.CreatedAt,
			&updatedBy,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if updatedBy.Valid {
			addonRule.UpdatedBy = &updatedBy.Int64
		}
		if updatedAt.Valid {
			addonRule.UpdatedAt = &updatedAt.Time
		}

		addonRules = append(addonRules, addonRule)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return addonRules, nil
}

// GetByAddonCode retrieves addon rules by addon code
func (r *MySQLAddonRuleRepository) GetByAddonCode(addonCode string) ([]model.AddonRule, error) {
	query := `
		SELECT id, product_code, addon_code, insurance_code, rules, created_by, created_at, updated_by, updated_at
		FROM travel_service_development.addon_rules
		WHERE addon_code = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, addonCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addonRules []model.AddonRule
	for rows.Next() {
		var addonRule model.AddonRule
		var updatedBy sql.NullInt64
		var updatedAt sql.NullTime
		var insuranceCode sql.NullString
		var rules sql.NullString
		
		err := rows.Scan(
			&addonRule.ID,
			&addonRule.ProductCode,
			&addonRule.AddonCode,
			&insuranceCode,
			&rules,
			&addonRule.CreatedBy,
			&addonRule.CreatedAt,
			&updatedBy,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if insuranceCode.Valid {
			addonRule.InsuranceCode = &insuranceCode.String
		}
		if rules.Valid {
			addonRule.Rules = &rules.String
		}
		if updatedBy.Valid {
			addonRule.UpdatedBy = &updatedBy.Int64
		}
		if updatedAt.Valid {
			addonRule.UpdatedAt = &updatedAt.Time
		}

		addonRules = append(addonRules, addonRule)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return addonRules, nil
}

// extractInsuranceCode extracts insurance code from product code (e.g., "TV-DAMAI-ASN-IND-01" -> "DAMAI")
func extractInsuranceCode(productCode string) string {
	parts := strings.Split(productCode, "-")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

// Create creates a new addon rule
func (r *MySQLAddonRuleRepository) Create(addonRule *model.AddonRule) error {
	insuranceCode := extractInsuranceCode(addonRule.ProductCode)
	
	query := `
		INSERT INTO travel_service_development.addon_rules (product_code, addon_code, insurance_code, rules, created_by, created_at, updated_at)
		VALUES (?, ?, ?, '{}', ?, NOW(), NOW())
	`

	result, err := r.db.Exec(query,
		addonRule.ProductCode,
		addonRule.AddonCode,
		insuranceCode,
		addonRule.CreatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create addon rule: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	addonRule.ID = id
	addonRule.CreatedAt = time.Now()
	now := time.Now()
	addonRule.UpdatedAt = &now

	return nil
}

// CreateBatch creates multiple addon rules in a single transaction
func (r *MySQLAddonRuleRepository) CreateBatch(addonRules []model.AddonRule) error {
	if len(addonRules) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Remove duplicates from addonRules slice before inserting
	// Create a map to track unique product_code + addon_code combinations
	type AddonRuleKey struct {
		ProductCode string
		AddonCode   string
	}
	uniqueRules := make(map[AddonRuleKey]model.AddonRule)
	for _, rule := range addonRules {
		key := AddonRuleKey{
			ProductCode: rule.ProductCode,
			AddonCode:   rule.AddonCode,
		}
		// Only keep the first occurrence of each unique combination
		if _, exists := uniqueRules[key]; !exists {
			uniqueRules[key] = rule
		}
	}

	// Convert map back to slice
	uniqueAddonRules := make([]model.AddonRule, 0, len(uniqueRules))
	for _, rule := range uniqueRules {
		uniqueAddonRules = append(uniqueAddonRules, rule)
	}

	query := `
		INSERT INTO travel_service_development.addon_rules (product_code, addon_code, insurance_code, rules, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for i := range uniqueAddonRules {
		// Use InsuranceCode from addonRule if available, otherwise extract from ProductCode
		var insuranceCode string
		if uniqueAddonRules[i].InsuranceCode != nil && *uniqueAddonRules[i].InsuranceCode != "" {
			insuranceCode = *uniqueAddonRules[i].InsuranceCode
		} else {
			insuranceCode = extractInsuranceCode(uniqueAddonRules[i].ProductCode)
		}
		
		// Use Rules from addonRule if available, otherwise use '{}'
		var rulesValue string
		if uniqueAddonRules[i].Rules != nil && *uniqueAddonRules[i].Rules != "" {
			rulesValue = *uniqueAddonRules[i].Rules
		} else {
			rulesValue = "{}"
		}
		
		result, err := stmt.Exec(
			uniqueAddonRules[i].ProductCode,
			uniqueAddonRules[i].AddonCode,
			insuranceCode,
			rulesValue,
			uniqueAddonRules[i].CreatedBy,
		)
		if err != nil {
			return fmt.Errorf("failed to insert addon rule: %w", err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id: %w", err)
		}

		uniqueAddonRules[i].ID = id
		uniqueAddonRules[i].CreatedAt = time.Now()
		now := time.Now()
		uniqueAddonRules[i].UpdatedAt = &now
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Update updates an existing addon rule
func (r *MySQLAddonRuleRepository) Update(id int64, addonRule *model.AddonRule) error {
	// Get data before update for comparison and history
	beforeAddonRule, err := r.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get addon rule before update: %w", err)
	}

	// Build changed fields map (only fields that are different)
	changedBefore := make(map[string]interface{})
	changedAfter := make(map[string]interface{})
	updateFields := []string{}
	updateValues := []interface{}{}

	// Compare each field and build update query dynamically
	if beforeAddonRule.ProductCode != addonRule.ProductCode {
		changedBefore["product_code"] = beforeAddonRule.ProductCode
		changedAfter["product_code"] = addonRule.ProductCode
		updateFields = append(updateFields, "product_code = ?")
		updateValues = append(updateValues, addonRule.ProductCode)
	}
	if beforeAddonRule.AddonCode != addonRule.AddonCode {
		changedBefore["addon_code"] = beforeAddonRule.AddonCode
		changedAfter["addon_code"] = addonRule.AddonCode
		updateFields = append(updateFields, "addon_code = ?")
		updateValues = append(updateValues, addonRule.AddonCode)
	}
	
	// Compare rules (need to handle nil pointers)
	var beforeRules, afterRules string
	if beforeAddonRule.Rules != nil {
		beforeRules = *beforeAddonRule.Rules
	}
	if addonRule.Rules != nil {
		afterRules = *addonRule.Rules
	}
	if beforeRules != afterRules {
		changedBefore["rules"] = beforeRules
		changedAfter["rules"] = afterRules
		updateFields = append(updateFields, "rules = ?")
		if addonRule.Rules != nil {
			updateValues = append(updateValues, *addonRule.Rules)
		} else {
			updateValues = append(updateValues, nil)
		}
	}

	// If no fields changed, return early (no update needed)
	if len(updateFields) == 0 {
		return nil
	}

	// Always update updated_by and updated_at
	var updatedBy interface{}
	if addonRule.UpdatedBy != nil {
		updatedBy = *addonRule.UpdatedBy
	} else {
		updatedBy = nil
	}
	updateFields = append(updateFields, "updated_by = ?", "updated_at = NOW()")
	updateValues = append(updateValues, updatedBy)
	updateValues = append(updateValues, id)

	// Build dynamic UPDATE query
	query := fmt.Sprintf(`
		UPDATE travel_service_development.addon_rules
		SET %s
		WHERE id = ?
	`, strings.Join(updateFields, ", "))

	result, err := r.db.Exec(query, updateValues...)
	if err != nil {
		return fmt.Errorf("failed to update addon rule: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrAddonRuleNotFound
	}

	now := time.Now()
	addonRule.UpdatedAt = &now

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
					Section:    "addon_rules",
					Action:     "UPDATE",
					RecordID:   int(id),
					RecordType: "addon_rule",
					DataBefore: string(beforeJSON),
					DataAfter:  string(afterJSON),
				}
				if err := r.historyRepo.Create(history); err != nil {
					// Log error but don't fail the update
					fmt.Printf("Warning: failed to save history: %v\n", err)
				} else {
					fmt.Printf("DEBUG: History saved successfully for addon rule ID %d with %d changed field(s)\n", id, len(changedBefore))
				}
			}
		}
	}

	return nil
}

// Delete deletes an addon rule by ID
func (r *MySQLAddonRuleRepository) Delete(id int64) error {
	query := `DELETE FROM travel_service_development.addon_rules WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrAddonRuleNotFound
	}

	return nil
}

// DeleteByProductCodeAndAddonCode deletes addon rule by product code and addon code
func (r *MySQLAddonRuleRepository) DeleteByProductCodeAndAddonCode(productCode, addonCode string) error {
	query := `DELETE FROM travel_service_development.addon_rules WHERE product_code = ? AND addon_code = ?`

	result, err := r.db.Exec(query, productCode, addonCode)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrAddonRuleNotFound
	}

	return nil
}

// DeleteByProductCode deletes all addon rules for a product code
func (r *MySQLAddonRuleRepository) DeleteByProductCode(productCode string) error {
	query := `DELETE FROM travel_service_development.addon_rules WHERE product_code = ?`

	_, err := r.db.Exec(query, productCode)
	if err != nil {
		return fmt.Errorf("failed to delete addon rules by product code: %w", err)
	}

	return nil
}

// CleanupDuplicates removes duplicate addon_rules entries (keeping only the one with the minimum id for each (product_code, addon_code) combination)
func (r *MySQLAddonRuleRepository) CleanupDuplicates() error {
	// Delete duplicates, keeping only the one with the minimum id for each (product_code, addon_code) combination
	query := `
		DELETE ar1 FROM travel_service_development.addon_rules ar1
		INNER JOIN travel_service_development.addon_rules ar2
		WHERE ar1.product_code = ar2.product_code
		AND ar1.addon_code = ar2.addon_code
		AND ar1.id > ar2.id
	`

	result, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to cleanup duplicate addon rules: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	fmt.Printf("DEBUG: Cleaned up %d duplicate addon rule entries\n", rowsAffected)

	return nil
}

