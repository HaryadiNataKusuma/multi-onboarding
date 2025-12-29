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

// ProductRuleRepository defines the interface for product rule operations
type ProductRuleRepository interface {
	// Draft operations
	AddDraft(rule *model.ProductRuleDraft) error
	GetAllDrafts() ([]model.ProductRuleDraft, error)
	UpdateDraft(index int, rule *model.ProductRuleDraft) error
	DeleteDraft(index int) error
	ClearDrafts() error
	ConfirmDrafts(createdBy int64) error

	// Database operations
	Create(rule *model.ProductRule) error
	GetAll(productCode string, insuranceCode string) ([]model.ProductRule, error)
	GetByID(id int64) (*model.ProductRule, error)
	Update(id int64, rule *model.ProductRule) error
	Delete(id int64) error
}

// MySQLProductRuleRepository implements ProductRuleRepository with MySQL database storage
type MySQLProductRuleRepository struct {
	db         *sql.DB
	historyRepo HistoryRepository
	draftData []model.ProductRuleDraft
	draftMu   sync.RWMutex
}

// NewMySQLProductRuleRepository creates a new MySQL product rule repository
func NewMySQLProductRuleRepository(db *sql.DB) *MySQLProductRuleRepository {
	return &MySQLProductRuleRepository{
		db:        db,
		draftData: make([]model.ProductRuleDraft, 0),
	}
}

// SetHistoryRepository sets the history repository for saving update histories
func (r *MySQLProductRuleRepository) SetHistoryRepository(historyRepo HistoryRepository) {
	r.historyRepo = historyRepo
}

// Draft operations (in-memory)

func (r *MySQLProductRuleRepository) AddDraft(rule *model.ProductRuleDraft) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	r.draftData = append(r.draftData, *rule)
	return nil
}

func (r *MySQLProductRuleRepository) GetAllDrafts() ([]model.ProductRuleDraft, error) {
	r.draftMu.RLock()
	defer r.draftMu.RUnlock()

	result := make([]model.ProductRuleDraft, len(r.draftData))
	copy(result, r.draftData)
	return result, nil
}

func (r *MySQLProductRuleRepository) UpdateDraft(index int, rule *model.ProductRuleDraft) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	if index < 0 || index >= len(r.draftData) {
		return fmt.Errorf("draft index out of range")
	}

	r.draftData[index] = *rule
	return nil
}

func (r *MySQLProductRuleRepository) DeleteDraft(index int) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	if index < 0 || index >= len(r.draftData) {
		return fmt.Errorf("draft index out of range")
	}

	r.draftData = append(r.draftData[:index], r.draftData[index+1:]...)
	return nil
}

func (r *MySQLProductRuleRepository) ClearDrafts() error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	r.draftData = make([]model.ProductRuleDraft, 0)
	return nil
}

func (r *MySQLProductRuleRepository) ConfirmDrafts(createdBy int64) error {
	// Get all drafts
	drafts, err := r.GetAllDrafts()
	if err != nil {
		return fmt.Errorf("failed to get drafts: %w", err)
	}

	// Insert each draft into product_rules table
	for _, draft := range drafts {
		rule := &model.ProductRule{
			ProductCode:      draft.ProductCode,
			PremiumType:      draft.PremiumType,
			StartDays:        draft.StartDays,
			EndDays:          draft.EndDays,
			BasePremiumValue: draft.BasePremiumValue,
			Rules:            draft.Rules,
			CreatedBy:        createdBy,
			CreatedAt:         time.Now(),
		}

		if err := r.Create(rule); err != nil {
			return fmt.Errorf("failed to confirm draft: %w", err)
		}
	}

	// Clear drafts after confirmation
	if err := r.ClearDrafts(); err != nil {
		return fmt.Errorf("failed to clear drafts: %w", err)
	}

	return nil
}

// Database operations

func (r *MySQLProductRuleRepository) Create(rule *model.ProductRule) error {
	query := `
		INSERT INTO product_rules (product_code, premium_type, start_days, end_days, base_premium_value, rules, created_by, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())
	`

	result, err := r.db.Exec(query,
		rule.ProductCode,
		rule.PremiumType,
		rule.StartDays,
		rule.EndDays,
		rule.BasePremiumValue,
		rule.Rules,
		rule.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to create product rule: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	rule.ID = id
	rule.CreatedAt = time.Now()
	return nil
}

func (r *MySQLProductRuleRepository) GetAll(productCode string, insuranceCode string) ([]model.ProductRule, error) {
	var query string
	var args []interface{}

	// Build WHERE clause
	whereClauses := []string{}
	if productCode != "" {
		whereClauses = append(whereClauses, "pr.product_code = ?")
		args = append(args, productCode)
	}
	if insuranceCode != "" {
		whereClauses = append(whereClauses, "p.insurance_code = ?")
		args = append(args, insuranceCode)
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	query = `
		SELECT pr.id, pr.product_code, p.name as product_name, pr.premium_type, pr.start_days, pr.end_days, 
		       pr.base_premium_value, pr.rules, pr.created_by, pr.created_at, pr.updated_by, pr.updated_at
		FROM product_rules pr
		LEFT JOIN products p ON pr.product_code = p.code
		` + whereClause + `
		ORDER BY pr.created_at DESC
	`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query product rules: %w", err)
	}
	defer rows.Close()

	rules := []model.ProductRule{}
	for rows.Next() {
		var rule model.ProductRule
		var updatedBy sql.NullInt64
		var updatedAt sql.NullTime
		var productName sql.NullString
		err := rows.Scan(
			&rule.ID,
			&rule.ProductCode,
			&productName,
			&rule.PremiumType,
			&rule.StartDays,
			&rule.EndDays,
			&rule.BasePremiumValue,
			&rule.Rules,
			&rule.CreatedBy,
			&rule.CreatedAt,
			&updatedBy,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product rule: %w", err)
		}

		if updatedBy.Valid {
			rule.UpdatedBy = &updatedBy.Int64
		}
		if updatedAt.Valid {
			rule.UpdatedAt = &updatedAt.Time
		}
		if productName.Valid {
			rule.ProductName = &productName.String
		}

		rules = append(rules, rule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return rules, nil
}

func (r *MySQLProductRuleRepository) GetByID(id int64) (*model.ProductRule, error) {
	query := `
		SELECT id, product_code, premium_type, start_days, end_days, base_premium_value, rules, created_by, created_at, updated_by, updated_at
		FROM product_rules
		WHERE id = ?
	`

	rule := &model.ProductRule{}
	var updatedBy sql.NullInt64
	var updatedAt sql.NullTime
	err := r.db.QueryRow(query, id).Scan(
		&rule.ID,
		&rule.ProductCode,
		&rule.PremiumType,
		&rule.StartDays,
		&rule.EndDays,
		&rule.BasePremiumValue,
		&rule.Rules,
		&rule.CreatedBy,
		&rule.CreatedAt,
		&updatedBy,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product rule not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product rule: %w", err)
	}

	if updatedBy.Valid {
		rule.UpdatedBy = &updatedBy.Int64
	}
	if updatedAt.Valid {
		rule.UpdatedAt = &updatedAt.Time
	}

	return rule, nil
}

func (r *MySQLProductRuleRepository) Update(id int64, rule *model.ProductRule) error {
	// Get data before update for comparison and history
	beforeRule, err := r.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get product rule before update: %w", err)
	}

	// Build changed fields map (only fields that are different)
	changedBefore := make(map[string]interface{})
	changedAfter := make(map[string]interface{})
	updateFields := []string{}
	updateValues := []interface{}{}

	// Compare each field and build update query dynamically
	if beforeRule.ProductCode != rule.ProductCode {
		changedBefore["product_code"] = beforeRule.ProductCode
		changedAfter["product_code"] = rule.ProductCode
		updateFields = append(updateFields, "product_code = ?")
		updateValues = append(updateValues, rule.ProductCode)
	}
	if beforeRule.PremiumType != rule.PremiumType {
		changedBefore["premium_type"] = beforeRule.PremiumType
		changedAfter["premium_type"] = rule.PremiumType
		updateFields = append(updateFields, "premium_type = ?")
		updateValues = append(updateValues, rule.PremiumType)
	}
	if beforeRule.StartDays != rule.StartDays {
		changedBefore["start_days"] = beforeRule.StartDays
		changedAfter["start_days"] = rule.StartDays
		updateFields = append(updateFields, "start_days = ?")
		updateValues = append(updateValues, rule.StartDays)
	}
	if beforeRule.EndDays != rule.EndDays {
		changedBefore["end_days"] = beforeRule.EndDays
		changedAfter["end_days"] = rule.EndDays
		updateFields = append(updateFields, "end_days = ?")
		updateValues = append(updateValues, rule.EndDays)
	}
	if beforeRule.BasePremiumValue != rule.BasePremiumValue {
		changedBefore["base_premium_value"] = beforeRule.BasePremiumValue
		changedAfter["base_premium_value"] = rule.BasePremiumValue
		updateFields = append(updateFields, "base_premium_value = ?")
		updateValues = append(updateValues, rule.BasePremiumValue)
	}
	if beforeRule.Rules != rule.Rules {
		changedBefore["rules"] = beforeRule.Rules
		changedAfter["rules"] = rule.Rules
		updateFields = append(updateFields, "rules = ?")
		updateValues = append(updateValues, rule.Rules)
	}

	// If no fields changed, return early (no update needed)
	if len(updateFields) == 0 {
		return nil
	}

	// Always update updated_by and updated_at
	var updatedBy interface{}
	if rule.UpdatedBy != nil {
		updatedBy = *rule.UpdatedBy
	} else {
		updatedBy = nil
	}
	updateFields = append(updateFields, "updated_by = ?", "updated_at = NOW()")
	updateValues = append(updateValues, updatedBy)
	updateValues = append(updateValues, id)

	// Build dynamic UPDATE query
	query := fmt.Sprintf(`
		UPDATE travel_service_development.product_rules
		SET %s
		WHERE id = ?
	`, strings.Join(updateFields, ", "))

	result, err := r.db.Exec(query, updateValues...)
	if err != nil {
		return fmt.Errorf("failed to update product rule: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product rule not found")
	}

	rule.ID = id
	now := time.Now()
	rule.UpdatedAt = &now

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
					Section:    "product_rules",
					Action:     "UPDATE",
					RecordID:   int(id),
					RecordType: "product_rule",
					DataBefore: string(beforeJSON),
					DataAfter:  string(afterJSON),
				}
				if err := r.historyRepo.Create(history); err != nil {
					// Log error but don't fail the update
					fmt.Printf("Warning: failed to save history: %v\n", err)
				} else {
					fmt.Printf("DEBUG: History saved successfully for product rule ID %d with %d changed field(s)\n", id, len(changedBefore))
				}
			}
		}
	}

	return nil
}

func (r *MySQLProductRuleRepository) Delete(id int64) error {
	query := `DELETE FROM product_rules WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product rule: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product rule not found")
	}

	return nil
}


