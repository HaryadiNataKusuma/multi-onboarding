package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"multi-onboarding/domain/vehicle/product"
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

	// Convert each draft to Product and insert into database
	for _, draft := range drafts {
		p := &product.Product{
			Code:                 draft.Code,
			InsuranceCode:        draft.InsuranceCode,
			InsuranceProductCode: draft.InsuranceProductCode,
			VehicleType:          draft.VehicleType,
			Name:                 draft.Name,
			Logo:                 draft.Logo,
			IsActive:             draft.IsActive,
			IsInsurerDriven:      draft.IsInsurerDriven,
			IsPersonalUsage:      draft.IsPersonalUsage,
			IsElectricVehicle:    draft.IsElectricVehicle,
			Summary:              draft.Summary,
			Config:               draft.Config,
			InsuranceDetail:      draft.InsuranceDetail,
			ProtectionDetail:     draft.ProtectionDetail,
			HowToClaim:           draft.HowToClaim,
			BaseProduct:          draft.BaseProduct,
			CreatedBy:            createdBy,
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

	if insuranceCode != "" {
		query = `
			SELECT id, code, insurance_code, insurance_product_code, vehicle_type,
			       is_insurance_api, name, logo, is_active, is_insurer_driven,
			       is_personal_usage, is_electric_vehicle, summary, config,
			       insurance_detail, protection_detail, how_to_claim, base_product,
			       created_by, created_at, updated_by, updated_at
			FROM vehicle_service_development.products
			WHERE insurance_code = ?
			ORDER BY created_at DESC
		`
		args = []interface{}{insuranceCode}
	} else {
		query = `
			SELECT id, code, insurance_code, insurance_product_code, vehicle_type,
			       is_insurance_api, name, logo, is_active, is_insurer_driven,
			       is_personal_usage, is_electric_vehicle, summary, config,
			       insurance_detail, protection_detail, how_to_claim, base_product,
			       created_by, created_at, updated_by, updated_at
			FROM vehicle_service_development.products
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
			&p.InsuranceProductCode,
			&p.VehicleType,
			&p.IsInsuranceAPI,
			&p.Name,
			&p.Logo,
			&p.IsActive,
			&p.IsInsurerDriven,
			&p.IsPersonalUsage,
			&p.IsElectricVehicle,
			&p.Summary,
			&p.Config,
			&p.InsuranceDetail,
			&p.ProtectionDetail,
			&p.HowToClaim,
			&p.BaseProduct,
		&p.CreatedBy,
		&p.CreatedAt,
		&updatedBy,
		&updatedAt,
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
		SELECT id, code, insurance_code, insurance_product_code, vehicle_type,
		       is_insurance_api, name, logo, is_active, is_insurer_driven,
		       is_personal_usage, is_electric_vehicle, summary, config,
		       insurance_detail, protection_detail, how_to_claim, base_product,
		       created_by, created_at, updated_by, updated_at
		FROM vehicle_service_development.products
		WHERE id = ?
	`

	var p product.Product
	var updatedBy sql.NullInt64
	var updatedAt sql.NullTime

	err := r.mysqlSess.QueryRow(query, id).Scan(
		&p.ID,
		&p.Code,
		&p.InsuranceCode,
		&p.InsuranceProductCode,
		&p.VehicleType,
		&p.IsInsuranceAPI,
		&p.Name,
		&p.Logo,
		&p.IsActive,
		&p.IsInsurerDriven,
		&p.IsPersonalUsage,
		&p.IsElectricVehicle,
		&p.Summary,
		&p.Config,
		&p.InsuranceDetail,
		&p.ProtectionDetail,
		&p.HowToClaim,
		&p.BaseProduct,
		&p.CreatedBy,
		&p.CreatedAt,
		&updatedBy,
		&updatedAt,
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
		INSERT INTO vehicle_service_development.products 
		(code, insurance_code, insurance_product_code, vehicle_type, is_insurance_api,
		 name, logo, is_active, is_insurer_driven, is_personal_usage, is_electric_vehicle,
		 summary, config, insurance_detail, protection_detail, how_to_claim, base_product,
		 created_by, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
	`

	result, err := r.mysqlSess.Exec(query,
		p.Code,
		p.InsuranceCode,
		p.InsuranceProductCode,
		p.VehicleType,
		p.IsInsuranceAPI,
		p.Name,
		p.Logo,
		p.IsActive,
		p.IsInsurerDriven,
		p.IsPersonalUsage,
		p.IsElectricVehicle,
		p.Summary,
		p.Config,
		p.InsuranceDetail,
		p.ProtectionDetail,
		p.HowToClaim,
		p.BaseProduct,
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
			Code:                 p.Code,
			InsuranceCode:        p.InsuranceCode,
			InsuranceProductCode: p.InsuranceProductCode,
			VehicleType:          p.VehicleType,
			Name:                 p.Name,
			IsActive:             p.IsActive,
			IsInsurerDriven:      p.IsInsurerDriven,
			IsPersonalUsage:      p.IsPersonalUsage,
			IsElectricVehicle:    p.IsElectricVehicle,
			BaseProduct:          p.BaseProduct,
			Summary:              p.Summary,
			Config:               p.Config,
			InsuranceDetail:      p.InsuranceDetail,
			ProtectionDetail:     p.ProtectionDetail,
			HowToClaim:           p.HowToClaim,
			Logo:                 p.Logo,
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
	if beforeData.InsuranceProductCode != p.InsuranceProductCode {
		changedBefore["insurance_product_code"] = beforeData.InsuranceProductCode
		changedAfter["insurance_product_code"] = p.InsuranceProductCode
		updateFields = append(updateFields, "insurance_product_code = ?")
		updateValues = append(updateValues, p.InsuranceProductCode)
	}
	if beforeData.VehicleType != p.VehicleType {
		changedBefore["vehicle_type"] = beforeData.VehicleType
		changedAfter["vehicle_type"] = p.VehicleType
		updateFields = append(updateFields, "vehicle_type = ?")
		updateValues = append(updateValues, p.VehicleType)
	}
	if beforeData.IsActive != p.IsActive {
		changedBefore["is_active"] = beforeData.IsActive
		changedAfter["is_active"] = p.IsActive
		updateFields = append(updateFields, "is_active = ?")
		updateValues = append(updateValues, p.IsActive)
	}
	if beforeData.IsInsurerDriven != p.IsInsurerDriven {
		changedBefore["is_insurer_driven"] = beforeData.IsInsurerDriven
		changedAfter["is_insurer_driven"] = p.IsInsurerDriven
		updateFields = append(updateFields, "is_insurer_driven = ?")
		updateValues = append(updateValues, p.IsInsurerDriven)
	}
	if beforeData.IsPersonalUsage != p.IsPersonalUsage {
		changedBefore["is_personal_usage"] = beforeData.IsPersonalUsage
		changedAfter["is_personal_usage"] = p.IsPersonalUsage
		updateFields = append(updateFields, "is_personal_usage = ?")
		updateValues = append(updateValues, p.IsPersonalUsage)
	}
	if beforeData.IsElectricVehicle != p.IsElectricVehicle {
		changedBefore["is_electric_vehicle"] = beforeData.IsElectricVehicle
		changedAfter["is_electric_vehicle"] = p.IsElectricVehicle
		updateFields = append(updateFields, "is_electric_vehicle = ?")
		updateValues = append(updateValues, p.IsElectricVehicle)
	}
	if beforeData.BaseProduct != p.BaseProduct {
		changedBefore["base_product"] = beforeData.BaseProduct
		changedAfter["base_product"] = p.BaseProduct
		updateFields = append(updateFields, "base_product = ?")
		updateValues = append(updateValues, p.BaseProduct)
	}
	if beforeData.Summary != p.Summary {
		changedBefore["summary"] = beforeData.Summary
		changedAfter["summary"] = p.Summary
		updateFields = append(updateFields, "summary = ?")
		updateValues = append(updateValues, p.Summary)
	}
	if beforeData.Config != p.Config {
		changedBefore["config"] = beforeData.Config
		changedAfter["config"] = p.Config
		updateFields = append(updateFields, "config = ?")
		updateValues = append(updateValues, p.Config)
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
	if beforeData.Logo != p.Logo {
		changedBefore["logo"] = beforeData.Logo
		changedAfter["logo"] = p.Logo
		updateFields = append(updateFields, "logo = ?")
		updateValues = append(updateValues, p.Logo)
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
		UPDATE vehicle_service_development.products
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
	query := `DELETE FROM vehicle_service_development.products WHERE id = ?`

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
		SELECT code, insurance_code, insurance_product_code, vehicle_type, name, 
		       is_active, is_insurer_driven, is_personal_usage, is_electric_vehicle, base_product,
		       summary, config, insurance_detail, protection_detail, how_to_claim, logo
		FROM vehicle_service_development.products
		WHERE id = ?
	`

	var before product.ProductBeforeUpdate
	err := r.mysqlSess.QueryRow(query, id).Scan(
		&before.Code,
		&before.InsuranceCode,
		&before.InsuranceProductCode,
		&before.VehicleType,
		&before.Name,
		&before.IsActive,
		&before.IsInsurerDriven,
		&before.IsPersonalUsage,
		&before.IsElectricVehicle,
		&before.BaseProduct,
		&before.Summary,
		&before.Config,
		&before.InsuranceDetail,
		&before.ProtectionDetail,
		&before.HowToClaim,
		&before.Logo,
	)
	if err != nil {
		return nil, err
	}

	return &before, nil
}

func (r *repository) InsertHistory(userName, section, action string, recordID int, recordType, dataBefore, dataAfter string) error {
	historyQuery := `INSERT INTO vehicle_service_development.histories (user_name, section, action, record_id, record_type, data_before, data_after, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	_, err := r.mysqlSess.Exec(historyQuery, userName, section, action, recordID, recordType, dataBefore, dataAfter)
	if err != nil {
		return fmt.Errorf("failed to insert history: %w", err)
	}

	return nil
}

