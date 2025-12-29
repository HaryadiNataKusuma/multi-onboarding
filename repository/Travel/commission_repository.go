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

// CommissionRepository interface defines the contract for commission data operations
type CommissionRepository interface {
	GetAll() ([]model.Commission, error)
	GetByID(id int64) (*model.Commission, error)
	GetByProductCode(productCode string) ([]model.Commission, error)
	Create(commission *model.Commission) error
	Update(commission *model.Commission) error
	Delete(id int64) error
	// Draft operations
	GetAllDrafts() ([]model.CommissionDraft, error)
	CreateDraft(draft *model.CommissionDraft) error
	UpdateDraft(id int, draft *model.CommissionDraft) error
	DeleteDraft(id int) error
	DeleteAllDrafts() error
	ConfirmDrafts() error
	// Generate to multiple tables
	GenerateCommissionsToMultipleTables(insuranceCode string, productCodes []string, commissionPercentage, afPercentage, adminFee float64, commissionVATType, afVATType string) error
	// Get data from specific tables
	GetCommissionsFromTable(insuranceCode string) ([]map[string]interface{}, error)
	GetPlanCommissionsFromTable(insuranceCode string) ([]map[string]interface{}, error)
	GetDefaultConfigFromTable(insuranceCode string) ([]map[string]interface{}, error)
	// Check for duplicate product codes
	CheckDuplicateProductCodes(productCodes []string) (map[string][]string, error)
}

// MySQLCommissionRepository implements CommissionRepository using MySQL database
type MySQLCommissionRepository struct {
	db         *sql.DB
	historyRepo HistoryRepository
	draftData []model.CommissionDraft
	draftMu   sync.RWMutex
	nextID    int
}

// NewMySQLCommissionRepository creates a new MySQL commission repository
func NewMySQLCommissionRepository(db *sql.DB) *MySQLCommissionRepository {
	return &MySQLCommissionRepository{
		db:        db,
		draftData: make([]model.CommissionDraft, 0),
		nextID:    1,
	}
}

// SetHistoryRepository sets the history repository for saving update histories
func (r *MySQLCommissionRepository) SetHistoryRepository(historyRepo HistoryRepository) {
	r.historyRepo = historyRepo
}

// GetAll retrieves all commissions
func (r *MySQLCommissionRepository) GetAll() ([]model.Commission, error) {
	query := `
		SELECT id, product_code, commission_percentage, commission_vat_type, 
		       af_percentage, af_vat_type, admin_fee, created_by, created_at, updated_by, updated_at
		FROM travel_service_development.commissions
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commissions []model.Commission
	for rows.Next() {
		var commission model.Commission
		var updatedBy sql.NullInt64
		var updatedAt sql.NullTime

		err := rows.Scan(
			&commission.ID,
			&commission.ProductCode,
			&commission.CommissionPercentage,
			&commission.CommissionVATType,
			&commission.AFPercentage,
			&commission.AFVATType,
			&commission.AdminFee,
			&commission.CreatedBy,
			&commission.CreatedAt,
			&updatedBy,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if updatedBy.Valid {
			updatedByVal := updatedBy.Int64
			commission.UpdatedBy = &updatedByVal
		}
		if updatedAt.Valid {
			updatedAtVal := updatedAt.Time
			commission.UpdatedAt = &updatedAtVal
		}

		commissions = append(commissions, commission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return commissions, nil
}

// GetByID retrieves a commission by ID
func (r *MySQLCommissionRepository) GetByID(id int64) (*model.Commission, error) {
	query := `
		SELECT id, product_code, commission_percentage, commission_vat_type, 
		       af_percentage, af_vat_type, admin_fee, created_by, created_at, updated_by, updated_at
		FROM travel_service_development.commissions
		WHERE id = ?
	`

	var commission model.Commission
	var updatedBy sql.NullInt64
	var updatedAt sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&commission.ID,
		&commission.ProductCode,
		&commission.CommissionPercentage,
		&commission.CommissionVATType,
		&commission.AFPercentage,
		&commission.AFVATType,
		&commission.AdminFee,
		&commission.CreatedBy,
		&commission.CreatedAt,
		&updatedBy,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrCommissionNotFound
	}
	if err != nil {
		return nil, err
	}

	if updatedBy.Valid {
		updatedByVal := updatedBy.Int64
		commission.UpdatedBy = &updatedByVal
	}
	if updatedAt.Valid {
		updatedAtVal := updatedAt.Time
		commission.UpdatedAt = &updatedAtVal
	}

	return &commission, nil
}

// GetByProductCode retrieves commissions by product code
func (r *MySQLCommissionRepository) GetByProductCode(productCode string) ([]model.Commission, error) {
	query := `
		SELECT id, product_code, commission_percentage, commission_vat_type, 
		       af_percentage, af_vat_type, admin_fee, created_by, created_at, updated_by, updated_at
		FROM travel_service_development.commissions
		WHERE product_code = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, productCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commissions []model.Commission
	for rows.Next() {
		var commission model.Commission
		var updatedBy sql.NullInt64
		var updatedAt sql.NullTime

		err := rows.Scan(
			&commission.ID,
			&commission.ProductCode,
			&commission.CommissionPercentage,
			&commission.CommissionVATType,
			&commission.AFPercentage,
			&commission.AFVATType,
			&commission.AdminFee,
			&commission.CreatedBy,
			&commission.CreatedAt,
			&updatedBy,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if updatedBy.Valid {
			updatedByVal := updatedBy.Int64
			commission.UpdatedBy = &updatedByVal
		}
		if updatedAt.Valid {
			updatedAtVal := updatedAt.Time
			commission.UpdatedAt = &updatedAtVal
		}

		commissions = append(commissions, commission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return commissions, nil
}

// Create creates a new commission
func (r *MySQLCommissionRepository) Create(commission *model.Commission) error {
	query := `
		INSERT INTO travel_service_development.commissions 
		(product_code, commission_percentage, commission_vat_type, af_percentage, af_vat_type, admin_fee, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.Exec(query,
		commission.ProductCode,
		commission.CommissionPercentage,
		commission.CommissionVATType,
		commission.AFPercentage,
		commission.AFVATType,
		commission.AdminFee,
		commission.CreatedBy,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	commission.ID = id
	commission.CreatedAt = time.Now()
	now := time.Now()
	commission.UpdatedAt = &now

	return nil
}

// Update updates an existing commission
func (r *MySQLCommissionRepository) Update(commission *model.Commission) error {
	// Get data before update for comparison and history
	beforeCommission, err := r.GetByID(commission.ID)
	if err != nil {
		return fmt.Errorf("failed to get commission before update: %w", err)
	}

	// Build changed fields map (only fields that are different)
	changedBefore := make(map[string]interface{})
	changedAfter := make(map[string]interface{})
	updateFields := []string{}
	updateValues := []interface{}{}

	// Compare each field and build update query dynamically
	if beforeCommission.ProductCode != commission.ProductCode {
		changedBefore["product_code"] = beforeCommission.ProductCode
		changedAfter["product_code"] = commission.ProductCode
		updateFields = append(updateFields, "product_code = ?")
		updateValues = append(updateValues, commission.ProductCode)
	}
	if beforeCommission.CommissionPercentage != commission.CommissionPercentage {
		changedBefore["commission_percentage"] = beforeCommission.CommissionPercentage
		changedAfter["commission_percentage"] = commission.CommissionPercentage
		updateFields = append(updateFields, "commission_percentage = ?")
		updateValues = append(updateValues, commission.CommissionPercentage)
	}
	if beforeCommission.CommissionVATType != commission.CommissionVATType {
		changedBefore["commission_vat_type"] = beforeCommission.CommissionVATType
		changedAfter["commission_vat_type"] = commission.CommissionVATType
		updateFields = append(updateFields, "commission_vat_type = ?")
		updateValues = append(updateValues, commission.CommissionVATType)
	}
	if beforeCommission.AFPercentage != commission.AFPercentage {
		changedBefore["af_percentage"] = beforeCommission.AFPercentage
		changedAfter["af_percentage"] = commission.AFPercentage
		updateFields = append(updateFields, "af_percentage = ?")
		updateValues = append(updateValues, commission.AFPercentage)
	}
	if beforeCommission.AFVATType != commission.AFVATType {
		changedBefore["af_vat_type"] = beforeCommission.AFVATType
		changedAfter["af_vat_type"] = commission.AFVATType
		updateFields = append(updateFields, "af_vat_type = ?")
		updateValues = append(updateValues, commission.AFVATType)
	}
	if beforeCommission.AdminFee != commission.AdminFee {
		changedBefore["admin_fee"] = beforeCommission.AdminFee
		changedAfter["admin_fee"] = commission.AdminFee
		updateFields = append(updateFields, "admin_fee = ?")
		updateValues = append(updateValues, commission.AdminFee)
	}

	// If no fields changed, return early (no update needed)
	if len(updateFields) == 0 {
		return nil
	}

	// Always update updated_by and updated_at
	updateFields = append(updateFields, "updated_by = ?", "updated_at = NOW()")
	updateValues = append(updateValues, commission.UpdatedBy)
	updateValues = append(updateValues, commission.ID)

	// Build dynamic UPDATE query
	query := fmt.Sprintf(`
		UPDATE travel_service_development.commissions
		SET %s
		WHERE id = ?
	`, strings.Join(updateFields, ", "))

	result, err := r.db.Exec(query, updateValues...)
	if err != nil {
		return fmt.Errorf("failed to update commission: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrCommissionNotFound
	}

	now := time.Now()
	commission.UpdatedAt = &now

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
					Section:    "commissions",
					Action:     "UPDATE",
					RecordID:   int(commission.ID),
					RecordType: "commission",
					DataBefore: string(beforeJSON),
					DataAfter:  string(afterJSON),
				}
				if err := r.historyRepo.Create(history); err != nil {
					// Log error but don't fail the update
					fmt.Printf("Warning: failed to save history: %v\n", err)
				} else {
					fmt.Printf("DEBUG: History saved successfully for commission ID %d with %d changed field(s)\n", commission.ID, len(changedBefore))
				}
			}
		}
	}

	return nil
}

// Delete deletes a commission by ID
func (r *MySQLCommissionRepository) Delete(id int64) error {
	query := `DELETE FROM travel_service_development.commissions WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrCommissionNotFound
	}

	return nil
}

// GetAllDrafts retrieves all commission drafts (in-memory)
func (r *MySQLCommissionRepository) GetAllDrafts() ([]model.CommissionDraft, error) {
	r.draftMu.RLock()
	defer r.draftMu.RUnlock()

	result := make([]model.CommissionDraft, len(r.draftData))
	copy(result, r.draftData)
	return result, nil
}

// CreateDraft creates a new commission draft (in-memory)
func (r *MySQLCommissionRepository) CreateDraft(draft *model.CommissionDraft) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	draft.ID = r.nextID
	r.nextID++
	draft.CreatedAt = time.Now()
	r.draftData = append(r.draftData, *draft)
	return nil
}

// UpdateDraft updates an existing commission draft (in-memory)
func (r *MySQLCommissionRepository) UpdateDraft(id int, draft *model.CommissionDraft) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	for i := range r.draftData {
		if r.draftData[i].ID == id {
			draft.ID = id
			r.draftData[i] = *draft
			return nil
		}
	}

	return ErrCommissionNotFound
}

// DeleteDraft deletes a commission draft by ID (in-memory)
func (r *MySQLCommissionRepository) DeleteDraft(id int) error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	for i := range r.draftData {
		if r.draftData[i].ID == id {
			r.draftData = append(r.draftData[:i], r.draftData[i+1:]...)
			return nil
		}
	}

	return ErrCommissionNotFound
}

// DeleteAllDrafts deletes all commission drafts (in-memory)
func (r *MySQLCommissionRepository) DeleteAllDrafts() error {
	r.draftMu.Lock()
	defer r.draftMu.Unlock()

	r.draftData = make([]model.CommissionDraft, 0)
	r.nextID = 1
	return nil
}

// ConfirmDrafts confirms all drafts to the main commissions table
func (r *MySQLCommissionRepository) ConfirmDrafts() error {
	r.draftMu.RLock()
	drafts := make([]model.CommissionDraft, len(r.draftData))
	copy(drafts, r.draftData)
	r.draftMu.RUnlock()

	if len(drafts) == 0 {
		return nil
	}

	query := `
		INSERT INTO travel_service_development.commissions 
		(product_code, commission_percentage, commission_vat_type, af_percentage, af_vat_type, admin_fee, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, 1, NOW(), NOW())
	`

	for _, draft := range drafts {
		_, err := r.db.Exec(query,
			draft.ProductCode,
			draft.CommissionPercentage,
			draft.CommissionVATType,
			draft.AFPercentage,
			draft.AFVATType,
			draft.AdminFee,
		)
		if err != nil {
			return err
		}
	}

	// Clear drafts after confirmation
	return r.DeleteAllDrafts()
}

// GenerateCommissionsToMultipleTables generates commission data to 3 tables:
// 1. commission_service_development.commissions
// 2. agent_service_development.default_config_products
// 3. quotation_service_development.plan_commissions
func (r *MySQLCommissionRepository) GenerateCommissionsToMultipleTables(insuranceCode string, productCodes []string, commissionPercentage, afPercentage, adminFee float64, commissionVATType, afVATType string) error {
	// Agent levels
	agentLevels := []string{
		"GREEN", "SILVER", "GOLD", "DIAMOND", "PLATINUM", "CORPORATE", "CORPORATE2",
		"CORPORATE3", "CORPORATE4", "TIED", "SKB", "SHOPDRIVE", "MVPARTNERSHIP1", "DIRECTPROPERTY",
	}

	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// For each product code
	for _, productCode := range productCodes {
		// For each agent level
		for _, agentLevel := range agentLevels {
			// 1. Insert into commission_service_development.commissions
			// Set corporate_id: 1 for CORPORATE, CORPORATE2, CORPORATE3, CORPORATE4, otherwise 0
			corporateLevels := []string{"CORPORATE", "CORPORATE2", "CORPORATE3", "CORPORATE4"}
			corporateId := 0
			for _, cl := range corporateLevels {
				if agentLevel == cl {
					corporateId = 1
					break
				}
			}

			query1 := `
				INSERT INTO commission_service_development.commissions 
				(product_code, insurance_code, agent_level, corporate_id, basic_commission, bonus_point, note, start_date, end_date, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, 0, ?, NOW(), NOW(), NOW(), NOW())
				ON DUPLICATE KEY UPDATE 
					basic_commission = VALUES(basic_commission),
					corporate_id = VALUES(corporate_id),
					updated_at = NOW()
			`
			note := "Commission for " + productCode + " - " + agentLevel
			_, err = tx.Exec(query1, productCode, insuranceCode, agentLevel, corporateId, commissionPercentage, note)
			if err != nil {
				return err
			}

			// 2. Insert into agent_service_development.default_config_products
			query2 := `
				INSERT INTO agent_service_development.default_config_products 
				(level, product, commission, max_discount, insurer, kind, category, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, 'PL', 'TV', NOW(), NOW())
				ON DUPLICATE KEY UPDATE 
					commission = VALUES(commission),
					max_discount = VALUES(max_discount),
					kind = VALUES(kind),
					category = VALUES(category),
					updated_at = NOW()
			`
			_, err = tx.Exec(query2, agentLevel, productCode, commissionPercentage, commissionPercentage, insuranceCode)
			if err != nil {
				return err
			}
		}

		// 3. Insert into quotation_service_development.plan_commissions (one row per product, not per agent level)
		query3 := `
			INSERT INTO quotation_service_development.plan_commissions 
			(plan_code, product_id, commission_percentage, commission_vat_type, af_percentage, af_vat_type, admin_fee, hardcopy_fee, is_active, version, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, 0, 1, 1, NOW())
			ON DUPLICATE KEY UPDATE 
				commission_percentage = VALUES(commission_percentage),
				commission_vat_type = VALUES(commission_vat_type),
				af_percentage = VALUES(af_percentage),
				af_vat_type = VALUES(af_vat_type),
				admin_fee = VALUES(admin_fee),
				hardcopy_fee = VALUES(hardcopy_fee)
		`
		_, err = tx.Exec(query3, productCode, 5, commissionPercentage, commissionVATType, afPercentage, afVATType, adminFee)
		if err != nil {
			return err
		}
	}

	// Commit transaction
	return tx.Commit()
}

// GetCommissionsFromTable retrieves all data from commission_service_development.commissions
// Only returns data for product codes that exist in travel_service_development.products
func (r *MySQLCommissionRepository) GetCommissionsFromTable(insuranceCode string) ([]map[string]interface{}, error) {
	var query string
	var args []interface{}
	
	if insuranceCode != "" {
		query = `
			SELECT c.id, c.product_code, c.insurance_code, c.agent_level, c.corporate_id, 
			       c.basic_commission, c.bonus_point, c.note, c.start_date, c.end_date, 
			       c.created_at, c.updated_at
			FROM commission_service_development.commissions c
			INNER JOIN travel_service_development.products p ON c.product_code = p.code
			WHERE c.insurance_code = ? AND p.code IS NOT NULL
			ORDER BY c.product_code ASC, c.agent_level ASC
		`
		args = []interface{}{insuranceCode}
	} else {
		query = `
			SELECT c.id, c.product_code, c.insurance_code, c.agent_level, c.corporate_id, 
			       c.basic_commission, c.bonus_point, c.note, c.start_date, c.end_date, 
			       c.created_at, c.updated_at
			FROM commission_service_development.commissions c
			INNER JOIN travel_service_development.products p ON c.product_code = p.code
			WHERE p.code IS NOT NULL
			ORDER BY c.product_code ASC, c.agent_level ASC
		`
		args = []interface{}{}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				entry[col] = string(b)
			} else {
				entry[col] = val
			}
		}
		results = append(results, entry)
	}

	return results, nil
}

// GetPlanCommissionsFromTable retrieves all data from quotation_service_development.plan_commissions
// Only returns data for product codes that exist in travel_service_development.products
func (r *MySQLCommissionRepository) GetPlanCommissionsFromTable(insuranceCode string) ([]map[string]interface{}, error) {
	var query string
	var args []interface{}
	
	if insuranceCode != "" {
		query = `
			SELECT pc.id, pc.plan_code, pc.product_id, pc.insurer_id, pc.commission_percentage, 
			       pc.commission_vat_type, pc.af_percentage, pc.af_vat_type, pc.admin_fee, 
			       pc.hardcopy_fee, pc.is_active, pc.version, pc.created_at
			FROM quotation_service_development.plan_commissions pc
			INNER JOIN travel_service_development.products p ON pc.plan_code = p.code
			WHERE p.insurance_code = ? AND p.code IS NOT NULL
			ORDER BY pc.plan_code ASC
		`
		args = []interface{}{insuranceCode}
	} else {
		query = `
			SELECT pc.id, pc.plan_code, pc.product_id, pc.insurer_id, pc.commission_percentage, 
			       pc.commission_vat_type, pc.af_percentage, pc.af_vat_type, pc.admin_fee, 
			       pc.hardcopy_fee, pc.is_active, pc.version, pc.created_at
			FROM quotation_service_development.plan_commissions pc
			INNER JOIN travel_service_development.products p ON pc.plan_code = p.code
			WHERE p.code IS NOT NULL
			ORDER BY pc.plan_code ASC
		`
		args = []interface{}{}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				entry[col] = string(b)
			} else {
				entry[col] = val
			}
		}
		results = append(results, entry)
	}

	return results, nil
}

// GetDefaultConfigFromTable retrieves all data from agent_service_development.default_config_products
// Only returns data for product codes that exist in travel_service_development.products
func (r *MySQLCommissionRepository) GetDefaultConfigFromTable(insuranceCode string) ([]map[string]interface{}, error) {
	var query string
	var args []interface{}
	
	if insuranceCode != "" {
		query = `
			SELECT dcp.id, dcp.level, dcp.sequence, dcp.kind, dcp.category, dcp.insurer, dcp.product, 
			       dcp.payment_frequency, dcp.policy_year, dcp.selected, dcp.renewal_count, 
			       dcp.commission, dcp.point, dcp.upline_bonus, dcp.bonus, dcp.max_discount, 
			       dcp.` + "`order`" + `, dcp.created_at, dcp.updated_at
			FROM agent_service_development.default_config_products dcp
			INNER JOIN travel_service_development.products p ON dcp.product = p.code
			WHERE dcp.insurer = ? AND p.code IS NOT NULL
			ORDER BY dcp.product ASC, dcp.level ASC
		`
		args = []interface{}{insuranceCode}
	} else {
		query = `
			SELECT dcp.id, dcp.level, dcp.sequence, dcp.kind, dcp.category, dcp.insurer, dcp.product, 
			       dcp.payment_frequency, dcp.policy_year, dcp.selected, dcp.renewal_count, 
			       dcp.commission, dcp.point, dcp.upline_bonus, dcp.bonus, dcp.max_discount, 
			       dcp.` + "`order`" + `, dcp.created_at, dcp.updated_at
			FROM agent_service_development.default_config_products dcp
			INNER JOIN travel_service_development.products p ON dcp.product = p.code
			WHERE p.code IS NOT NULL
			ORDER BY dcp.product ASC, dcp.level ASC
		`
		args = []interface{}{}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				entry[col] = string(b)
			} else {
				entry[col] = val
			}
		}
		results = append(results, entry)
	}

	return results, nil
}

// CheckDuplicateProductCodes checks if product codes already exist in the 3 tables
// Returns a map with table names as keys and arrays of duplicate product codes as values
func (r *MySQLCommissionRepository) CheckDuplicateProductCodes(productCodes []string) (map[string][]string, error) {
	if len(productCodes) == 0 {
		return make(map[string][]string), nil
	}

	duplicates := make(map[string][]string)
	
	// Create placeholders for IN clause
	placeholders := ""
	args := make([]interface{}, len(productCodes))
	for i, code := range productCodes {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args[i] = code
	}

	// 1. Check commission_service_development.commissions
	query1 := `
		SELECT DISTINCT product_code 
		FROM commission_service_development.commissions 
		WHERE product_code IN (` + placeholders + `)
	`
	rows1, err := r.db.Query(query1, args...)
	if err != nil {
		return nil, err
	}
	defer rows1.Close()
	
	var duplicates1 []string
	for rows1.Next() {
		var productCode string
		if err := rows1.Scan(&productCode); err != nil {
			return nil, err
		}
		duplicates1 = append(duplicates1, productCode)
	}
	if len(duplicates1) > 0 {
		duplicates["commission_service_development.commissions"] = duplicates1
	}

	// 2. Check quotation_service_development.plan_commissions
	query2 := `
		SELECT DISTINCT plan_code 
		FROM quotation_service_development.plan_commissions 
		WHERE plan_code IN (` + placeholders + `)
	`
	rows2, err := r.db.Query(query2, args...)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()
	
	var duplicates2 []string
	for rows2.Next() {
		var planCode string
		if err := rows2.Scan(&planCode); err != nil {
			return nil, err
		}
		duplicates2 = append(duplicates2, planCode)
	}
	if len(duplicates2) > 0 {
		duplicates["quotation_service_development.plan_commissions"] = duplicates2
	}

	// 3. Check agent_service_development.default_config_products
	query3 := `
		SELECT DISTINCT product 
		FROM agent_service_development.default_config_products 
		WHERE product IN (` + placeholders + `)
	`
	rows3, err := r.db.Query(query3, args...)
	if err != nil {
		return nil, err
	}
	defer rows3.Close()
	
	var duplicates3 []string
	for rows3.Next() {
		var product string
		if err := rows3.Scan(&product); err != nil {
			return nil, err
		}
		duplicates3 = append(duplicates3, product)
	}
	if len(duplicates3) > 0 {
		duplicates["agent_service_development.default_config_products"] = duplicates3
	}

	return duplicates, nil
}
