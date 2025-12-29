package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"multi-onboarding/model/Travel"
)

var (
	ErrNotFound = errors.New("record not found")
)

// AddonRuleDetailRepository defines the interface for addon rule detail operations
type AddonRuleDetailRepository interface {
	AddDraft(draft *model.AddonRuleDetailDraft) error
	GetAllDrafts() ([]model.AddonRuleDetailDraft, error)
	ClearDrafts() error
	ConfirmDrafts(createdBy int64) error
	GetAll() ([]model.AddonRuleDetail, error)
	GetByID(id int64) (*model.AddonRuleDetail, error)
	Update(id int64, detail *model.AddonRuleDetail) error
}

// MySQLAddonRuleDetailRepository implements AddonRuleDetailRepository using MySQL database
type MySQLAddonRuleDetailRepository struct {
	db        *sql.DB
	draftData []model.AddonRuleDetailDraft
	draftMu   sync.RWMutex
	historyRepo HistoryRepository
}

// NewMySQLAddonRuleDetailRepository creates a new MySQL addon rule detail repository
func NewMySQLAddonRuleDetailRepository(db *sql.DB, historyRepo HistoryRepository) *MySQLAddonRuleDetailRepository {
	return &MySQLAddonRuleDetailRepository{
		db:         db,
		draftData:  make([]model.AddonRuleDetailDraft, 0),
		historyRepo: historyRepo,
	}
}

// AddDraft adds a draft addon rule detail (in-memory)
func (r *MySQLAddonRuleDetailRepository) AddDraft(draft *model.AddonRuleDetailDraft) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	// Add to in-memory draft data
	r.draftData = append(r.draftData, *draft)

	return nil
}

// GetAllDrafts retrieves all draft addon rule details (from memory)
func (r *MySQLAddonRuleDetailRepository) GetAllDrafts() ([]model.AddonRuleDetailDraft, error) {
	r.draftMu.RLock()
	defer r.draftMu.RUnlock()

	// Return a copy to prevent external modifications
	result := make([]model.AddonRuleDetailDraft, len(r.draftData))
	copy(result, r.draftData)
	return result, nil
}

// ClearDrafts clears all draft addon rule details (from memory)
func (r *MySQLAddonRuleDetailRepository) ClearDrafts() error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	r.draftData = make([]model.AddonRuleDetailDraft, 0)
	return nil
}

// ConfirmDrafts confirms all drafts by inserting them into the main table
func (r *MySQLAddonRuleDetailRepository) ConfirmDrafts(createdBy int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get all drafts
	drafts, err := r.GetAllDrafts()
	if err != nil {
		return err
	}

	// Insert into main table
	insertQuery := `
		INSERT INTO travel_service_development.addon_rule_details 
		(addon_rule_id, start_condition, end_condition, value_type, value, duration_rule_type, min_adult, max_adult, max_age, created_by, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1, NOW())
		ON DUPLICATE KEY UPDATE 
			start_condition = VALUES(start_condition),
			end_condition = VALUES(end_condition),
			value_type = VALUES(value_type),
			value = VALUES(value),
			duration_rule_type = VALUES(duration_rule_type),
			min_adult = VALUES(min_adult),
			max_adult = VALUES(max_adult),
			max_age = VALUES(max_age),
			updated_by = 1,
			updated_at = NOW()
	`

	stmt, err := tx.Prepare(insertQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, draft := range drafts {
		// Parse addon_rule_id from string to int64
		var addonRuleID int64
		if draft.AddonRuleID != "" {
			_, err := fmt.Sscanf(draft.AddonRuleID, "%d", &addonRuleID)
			if err != nil {
				return fmt.Errorf("invalid addon_rule_id: %s for product_code: %s, addon_code: %s", draft.AddonRuleID, draft.ProductCode, draft.AddonCode)
			}
		} else {
			return fmt.Errorf("addon_rule_id is required for product_code: %s, addon_code: %s", draft.ProductCode, draft.AddonCode)
		}

		_, err = stmt.Exec(
			addonRuleID,
			draft.StartCondition,
			draft.EndCondition,
			draft.ValueType,
			draft.Value,
			draft.DurationRuleType,
			draft.MinAdult,
			draft.MaxAdult,
			draft.MaxAge,
		)
		if err != nil {
			return fmt.Errorf("failed to insert addon rule detail: %w", err)
		}
	}

	// Clear drafts from memory
	r.draftMu.Lock()
	r.draftData = make([]model.AddonRuleDetailDraft, 0)
	r.draftMu.Unlock()

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetAll retrieves all addon rule details
func (r *MySQLAddonRuleDetailRepository) GetAll() ([]model.AddonRuleDetail, error) {
	query := `
		SELECT ard.id, ard.addon_rule_id, ard.start_condition, ard.end_condition, 
		       ard.value_type, ard.value, ard.duration_rule_type, ard.min_adult, 
		       ard.max_adult, ard.max_age, ard.created_by, ard.created_at, 
		       ard.updated_by, ard.updated_at
		FROM travel_service_development.addon_rule_details ard
		ORDER BY ard.id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []model.AddonRuleDetail
	for rows.Next() {
		var detail model.AddonRuleDetail
		var createdAt string
		var updatedBy sql.NullInt64
		var updatedAt sql.NullString

		err := rows.Scan(
			&detail.ID,
			&detail.AddonRuleID,
			&detail.StartCondition,
			&detail.EndCondition,
			&detail.ValueType,
			&detail.Value,
			&detail.DurationRuleType,
			&detail.MinAdult,
			&detail.MaxAdult,
			&detail.MaxAge,
			&detail.CreatedBy,
			&createdAt,
			&updatedBy,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		detail.CreatedAt = createdAt
		if updatedBy.Valid {
			detail.UpdatedBy = &updatedBy.Int64
		}
		if updatedAt.Valid {
			detail.UpdatedAt = &updatedAt.String
		}

		details = append(details, detail)
	}

	return details, nil
}

// GetByID retrieves an addon rule detail by ID
func (r *MySQLAddonRuleDetailRepository) GetByID(id int64) (*model.AddonRuleDetail, error) {
	query := `
		SELECT ard.id, ard.addon_rule_id, ard.start_condition, ard.end_condition, 
		       ard.value_type, ard.value, ard.duration_rule_type, ard.min_adult, 
		       ard.max_adult, ard.max_age, ard.created_by, ard.created_at, 
		       ard.updated_by, ard.updated_at
		FROM travel_service_development.addon_rule_details ard
		WHERE ard.id = ?
	`

	var detail model.AddonRuleDetail
	var createdAt string
	var updatedBy sql.NullInt64
	var updatedAt sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&detail.ID,
		&detail.AddonRuleID,
		&detail.StartCondition,
		&detail.EndCondition,
		&detail.ValueType,
		&detail.Value,
		&detail.DurationRuleType,
		&detail.MinAdult,
		&detail.MaxAdult,
		&detail.MaxAge,
		&detail.CreatedBy,
		&createdAt,
		&updatedBy,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get addon rule detail by ID: %w", err)
	}

	detail.CreatedAt = createdAt
	if updatedBy.Valid {
		detail.UpdatedBy = &updatedBy.Int64
	}
	if updatedAt.Valid {
		detail.UpdatedAt = &updatedAt.String
	}

	return &detail, nil
}

// Update updates an existing addon rule detail
func (r *MySQLAddonRuleDetailRepository) Update(id int64, detail *model.AddonRuleDetail) error {
	// Get data before update
	beforeDetail, err := r.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get data before update: %w", err)
	}

	// Build changed fields map (only fields that are different)
	changedBefore := make(map[string]interface{})
	changedAfter := make(map[string]interface{})

	if beforeDetail.StartCondition != detail.StartCondition {
		changedBefore["start_condition"] = beforeDetail.StartCondition
		changedAfter["start_condition"] = detail.StartCondition
	}
	if beforeDetail.EndCondition != detail.EndCondition {
		changedBefore["end_condition"] = beforeDetail.EndCondition
		changedAfter["end_condition"] = detail.EndCondition
	}
	if beforeDetail.ValueType != detail.ValueType {
		changedBefore["value_type"] = beforeDetail.ValueType
		changedAfter["value_type"] = detail.ValueType
	}
	if beforeDetail.Value != detail.Value {
		changedBefore["value"] = beforeDetail.Value
		changedAfter["value"] = detail.Value
	}
	if beforeDetail.DurationRuleType != detail.DurationRuleType {
		changedBefore["duration_rule_type"] = beforeDetail.DurationRuleType
		changedAfter["duration_rule_type"] = detail.DurationRuleType
	}
	if beforeDetail.MinAdult != detail.MinAdult {
		changedBefore["min_adult"] = beforeDetail.MinAdult
		changedAfter["min_adult"] = detail.MinAdult
	}
	if beforeDetail.MaxAdult != detail.MaxAdult {
		changedBefore["max_adult"] = beforeDetail.MaxAdult
		changedAfter["max_adult"] = detail.MaxAdult
	}
	if beforeDetail.MaxAge != detail.MaxAge {
		changedBefore["max_age"] = beforeDetail.MaxAge
		changedAfter["max_age"] = detail.MaxAge
	}

	// Update the record
	query := `
		UPDATE travel_service_development.addon_rule_details
		SET start_condition = ?, end_condition = ?, value_type = ?, value = ?,
		    duration_rule_type = ?, min_adult = ?, max_adult = ?, max_age = ?,
		    updated_by = 1, updated_at = NOW()
		WHERE id = ?
	`

	result, err := r.db.Exec(query,
		detail.StartCondition,
		detail.EndCondition,
		detail.ValueType,
		detail.Value,
		detail.DurationRuleType,
		detail.MinAdult,
		detail.MaxAdult,
		detail.MaxAge,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update addon rule detail: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	// Save history only if there are changes
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
					Section:    "addon_rule_details",
					Action:     "UPDATE",
					RecordID:   int(id),
					RecordType: "addon_rule_detail",
					DataBefore: string(beforeJSON),
					DataAfter:  string(afterJSON),
				}
				if err := r.historyRepo.Create(history); err != nil {
					// Log error but don't fail the update
					fmt.Printf("Warning: failed to save history: %v\n", err)
				}
			}
		}
	}

	return nil
}

