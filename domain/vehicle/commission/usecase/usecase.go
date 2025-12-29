package usecase

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/labstack/echo"
	"multi-onboarding/domain/vehicle/commission"
)

type usecase struct {
	repository commission.Repository
}

// Default agent levels
var defaultAgentLevels = []string{
	"GREEN", "SILVER", "GOLD", "DIAMOND", "PLATINUM",
	"CORPORATE", "CORPORATE2", "CORPORATE3", "CORPORATE4",
	"TIED", "SKB", "SHOPDRIVE", "MVPARTNERSHIP1", "DIRECTPROPERTY",
}

// NewCommissionUsecase creates a new commission usecase
func NewCommissionUsecase(repository commission.Repository) commission.Usecase {
	return &usecase{
		repository: repository,
	}
}

// SaveCommissionDraft saves commission draft data directly to database draft tables
func (u *usecase) SaveCommissionDraft(c echo.Context, req commission.SaveCommissionDraftRequest) error {
	log.Printf("SaveCommissionDraft usecase called with insurance_code: %s", req.InsuranceCode)

	// Get products for insurance
	productCodes, err := u.repository.GetProductsByInsuranceCode(c, req.InsuranceCode)
	if err != nil {
		log.Printf("ERROR: Failed to get products by insurance code %s: %v", req.InsuranceCode, err)
		return err
	}

	log.Printf("Found %d products for insurance_code %s", len(productCodes), req.InsuranceCode)
	if len(productCodes) == 0 {
		log.Printf("ERROR: No active products found for insurance_code %s", req.InsuranceCode)
		return &validationError{Message: "No active products found for this insurance"}
	}

	// Clear existing draft tables for this insurance code first
	log.Printf("Clearing existing draft tables for insurance_code %s", req.InsuranceCode)
	if err := u.repository.ClearDraftTablesForInsurance(c, req.InsuranceCode); err != nil {
		log.Printf("ERROR: Failed to clear draft tables: %v", err)
		return err
	}

	// Insert commissions_draft for each product and agent level
	log.Printf("Inserting commissions_draft for %d products...", len(productCodes))
	for _, productCode := range productCodes {
		for _, level := range defaultAgentLevels {
			noteValue := "Product " + productCode
			corporateID := 0
			if level == "CORPORATE" || level == "CORPORATE2" || level == "CORPORATE3" || level == "CORPORATE4" {
				corporateID = 1
			}

			if err := u.repository.InsertCommissionDraft(c, productCode, req.InsuranceCode, level, corporateID, req.BasicCommission, noteValue); err != nil {
				log.Printf("ERROR: Failed to insert commission draft for product %s, level %s: %v", productCode, level, err)
				return fmt.Errorf("failed to insert commission draft: %v", err)
			}
		}
	}

	// Insert plan_commissions_draft for each product
	log.Printf("Inserting plan_commissions_draft for %d products...", len(productCodes))
	for _, productCode := range productCodes {
		if err := u.repository.InsertPlanCommissionDraft(c, productCode, req.PlanCommissions.CommissionVatType, req.BasicCommission,
			req.PlanCommissions.AfPercentage, req.PlanCommissions.AfVatType, req.PlanCommissions.AdminFee); err != nil {
			log.Printf("ERROR: Failed to insert plan commission draft for product %s: %v", productCode, err)
			return fmt.Errorf("failed to insert plan commission draft: %v", err)
		}
	}

	// Insert default_config_products_draft for each product and agent level
	log.Printf("Inserting default_config_products_draft for %d products...", len(productCodes))
	for _, productCode := range productCodes {
		for _, level := range defaultAgentLevels {
			if err := u.repository.InsertConfigProductDraft(c, level, req.InsuranceCode, productCode, req.BasicCommission); err != nil {
				log.Printf("ERROR: Failed to insert config product draft for product %s, level %s: %v", productCode, level, err)
				return fmt.Errorf("failed to insert config product draft: %v", err)
			}
		}
	}

	log.Printf("SUCCESS: Commission draft saved to database draft tables for insurance_code %s", req.InsuranceCode)
	return nil
}

// GetCommissionDraft gets commission draft data from database draft tables
func (u *usecase) GetCommissionDraft(c echo.Context, insuranceCode string) (*commission.GetCommissionDraftResponse, error) {
	log.Printf("GetCommissionDraft called with insurance_code: %s", insuranceCode)

	// Get draft data from database
	commissions, err := u.repository.GetCommissionDraft(c, insuranceCode)
	if err != nil {
		log.Printf("ERROR: Failed to get commission draft: %v", err)
		return nil, err
	}

	planCommissions, err := u.repository.GetPlanCommissionDraft(c, insuranceCode)
	if err != nil {
		log.Printf("ERROR: Failed to get plan commission draft: %v", err)
		return nil, err
	}

	configProducts, err := u.repository.GetConfigProductDraft(c, insuranceCode)
	if err != nil {
		log.Printf("ERROR: Failed to get config product draft: %v", err)
		return nil, err
	}

	log.Printf("Found %d commission drafts, %d plan commission drafts, %d config product drafts for insurance_code %s",
		len(commissions), len(planCommissions), len(configProducts), insuranceCode)

	return &commission.GetCommissionDraftResponse{
		Commissions:     commissions,
		PlanCommissions: planCommissions,
		ConfigProducts:  configProducts,
	}, nil
}

// ConfirmCommissions confirms commissions from database draft tables to confirmed tables
func (u *usecase) ConfirmCommissions(c echo.Context) (*commission.ConfirmCommissionsResult, error) {
	log.Printf("=== ConfirmCommissions called ===")

	// Check if there are any draft records in database
	planCommissions, err := u.repository.GetPlanCommissionDraft(c, "")
	if err != nil {
		log.Printf("ERROR: Failed to get plan commission draft: %v", err)
		return nil, err
	}

	if len(planCommissions) == 0 {
		log.Printf("WARNING: No plan commission drafts found in database")
		return nil, &validationError{Message: "No commission drafts available to confirm"}
	}

	log.Printf("Found %d plan commission draft records to confirm", len(planCommissions))

	// Move all draft data to confirmed tables
	log.Printf("Moving draft data to confirmed tables...")
	commissionsMoved, err := u.repository.MoveCommissionsDraftToConfirmed(c)
	if err != nil {
		log.Printf("ERROR: Failed to move commissions to confirmed: %v", err)
		return nil, fmt.Errorf("failed to move commissions to confirmed: %v", err)
	} else {
		log.Printf("Moved %d commission records to confirmed", commissionsMoved)
	}

	planCommissionsMoved, err := u.repository.MovePlanCommissionsDraftToConfirmed(c)
	if err != nil {
		log.Printf("ERROR: Failed to move plan commissions to confirmed: %v", err)
		return nil, fmt.Errorf("failed to move plan commissions to confirmed: %v", err)
	} else {
		log.Printf("Moved %d plan commission records to confirmed", planCommissionsMoved)
	}

	configProductsMoved, err := u.repository.MoveConfigProductsDraftToConfirmed(c)
	if err != nil {
		log.Printf("ERROR: Failed to move config products to confirmed: %v", err)
		return nil, fmt.Errorf("failed to move config products to confirmed: %v", err)
	} else {
		log.Printf("Moved %d config product records to confirmed", configProductsMoved)
	}

	// Update product rules hardcopy admin fee
	log.Printf("Updating product rules hardcopy admin fee...")
	if err := u.repository.UpdateProductRulesHardcopyAdminFee(c); err != nil {
		log.Printf("WARNING: Failed to update product rules hardcopy admin fee: %v", err)
	}

	// Clear draft tables after successful confirmation
	log.Printf("Clearing draft tables...")
	if err := u.repository.ClearDraftTables(c); err != nil {
		log.Printf("WARNING: Failed to clear draft tables: %v", err)
	}

	totalMoved := int(commissionsMoved + planCommissionsMoved + configProductsMoved)
	log.Printf("=== ConfirmCommissions completed. Total records moved: %d ===", totalMoved)
	return &commission.ConfirmCommissionsResult{
		TotalMoved: totalMoved,
	}, nil
}

// ClearCommissionDraft clears commission draft data from database draft tables
func (u *usecase) ClearCommissionDraft(c echo.Context) error {
	log.Printf("ClearCommissionDraft called")

	// Clear all draft tables
	if err := u.repository.ClearDraftTables(c); err != nil {
		log.Printf("ERROR: Failed to clear draft tables: %v", err)
		return err
	}

	log.Printf("SUCCESS: Commission drafts cleared from database draft tables")
	return nil
}

// GetConfirmedCommissions gets confirmed commissions
func (u *usecase) GetConfirmedCommissions(c echo.Context, insuranceCode, productCode string) (*commission.GetConfirmedCommissionsResponse, error) {
	commissions, err := u.repository.GetConfirmedCommissions(c, insuranceCode, productCode)
	if err != nil {
		return nil, err
	}

	planCommissions, err := u.repository.GetConfirmedPlanCommissions(c, insuranceCode, productCode)
	if err != nil {
		return nil, err
	}

	configProducts, err := u.repository.GetConfirmedConfigProducts(c, insuranceCode, productCode)
	if err != nil {
		return nil, err
	}

	return &commission.GetConfirmedCommissionsResponse{
		Commissions:     commissions,
		PlanCommissions: planCommissions,
		ConfigProducts:  configProducts,
	}, nil
}

// BulkUpdatePlanCommissions performs bulk update of plan commissions
func (u *usecase) BulkUpdatePlanCommissions(c echo.Context, req commission.BulkUpdatePlanCommissionsRequest, userName string) (int, error) {
	if len(req.IDs) == 0 {
		return 0, &validationError{Message: "At least one ID is required"}
	}

	// Get existing records
	existingRecords, err := u.repository.GetPlanCommissionByIDs(c, req.IDs)
	if err != nil {
		return 0, err
	}

	if len(existingRecords) == 0 {
		return 0, &notFoundError{Message: "No plan commissions found for the provided IDs"}
	}

	// Group by plan_code
	planCodeMap := make(map[string]commission.PlanCommissionRecord)
	for _, rec := range existingRecords {
		if existing, exists := planCodeMap[rec.PlanCode]; !exists ||
			(!existing.IsActive && rec.IsActive) ||
			(existing.IsActive == rec.IsActive && rec.Version > existing.Version) {
			planCodeMap[rec.PlanCode] = rec
		}
	}

	totalInserted := 0

	// Process each plan_code
	for planCode, rec := range planCodeMap {
		// Deactivate existing
		if err := u.repository.DeactivatePlanCommissions(c, []int{rec.ID}); err != nil {
			log.Printf("Error deactivating plan commissions for %s: %v", planCode, err)
			continue
		}

		// Get max version
		maxVersion := rec.Version
		newVersion := maxVersion + 1

		// Prepare new values
		newCommissionPercentage := rec.CommissionPercentage
		if req.Updates.CommissionPercentage != nil {
			newCommissionPercentage = *req.Updates.CommissionPercentage
		}

		newCommissionVatType := rec.CommissionVatType
		if req.Updates.CommissionVatType != nil {
			newCommissionVatType = *req.Updates.CommissionVatType
		}

		newAfPercentage := rec.AfPercentage
		if req.Updates.AfPercentage != nil {
			newAfPercentage = *req.Updates.AfPercentage
		}

		newAfVatType := rec.AfVatType
		if req.Updates.AfVatType != nil {
			newAfVatType = *req.Updates.AfVatType
		}

		newAdminFee := rec.AdminFee
		if req.Updates.AdminFee != nil {
			newAdminFee = *req.Updates.AdminFee
		}

		// Get insurer_id
		insurerID, err := u.repository.GetInsurerIDByPlanCode(c, planCode)
		if err != nil {
			insurerID = rec.InsurerID
		}

		// Insert new record
		if err := u.repository.InsertPlanCommission(c, planCode, insurerID, rec.ProductID, newCommissionPercentage,
			newCommissionVatType, newAfPercentage, newAfVatType, newAdminFee, rec.HardcopyFee, newVersion); err != nil {
			log.Printf("Error inserting plan commission for %s: %v", planCode, err)
			continue
		}

		totalInserted++

		// Update product_rules if admin_fee changed
		if req.Updates.AdminFee != nil {
			u.repository.UpdateProductRulesHardcopyAdminFeeByPlanCode(c, planCode, newAdminFee)
		}

		// Track changes for history
		changedFields := make(map[string]interface{})
		if req.Updates.CommissionPercentage != nil && *req.Updates.CommissionPercentage != rec.CommissionPercentage {
			changedFields["commission_percentage"] = map[string]interface{}{
				"before": rec.CommissionPercentage,
				"after":  *req.Updates.CommissionPercentage,
			}
		}
		if req.Updates.CommissionVatType != nil && *req.Updates.CommissionVatType != rec.CommissionVatType {
			changedFields["commission_vat_type"] = map[string]interface{}{
				"before": rec.CommissionVatType,
				"after":  *req.Updates.CommissionVatType,
			}
		}
		if req.Updates.AfPercentage != nil && *req.Updates.AfPercentage != rec.AfPercentage {
			changedFields["af_percentage"] = map[string]interface{}{
				"before": rec.AfPercentage,
				"after":  *req.Updates.AfPercentage,
			}
		}
		if req.Updates.AfVatType != nil && *req.Updates.AfVatType != rec.AfVatType {
			changedFields["af_vat_type"] = map[string]interface{}{
				"before": rec.AfVatType,
				"after":  *req.Updates.AfVatType,
			}
		}
		if req.Updates.AdminFee != nil && *req.Updates.AdminFee != rec.AdminFee {
			changedFields["admin_fee"] = map[string]interface{}{
				"before": rec.AdminFee,
				"after":  *req.Updates.AdminFee,
			}
		}

		// Insert history if there are changes
		if len(changedFields) > 0 {
			insuranceCode, _ := u.repository.GetInsuranceCodeByPlanCode(c, planCode)

			beforeChanged := make(map[string]interface{})
			afterChanged := make(map[string]interface{})

			beforeChanged["insurance_code"] = insuranceCode
			afterChanged["insurance_code"] = insuranceCode

			for field, changeData := range changedFields {
				changeMap := changeData.(map[string]interface{})
				beforeChanged[field] = changeMap["before"]
				afterChanged[field] = changeMap["after"]
			}

			beforeJSONBytes, _ := json.Marshal(beforeChanged)
			afterJSONBytes, _ := json.Marshal(afterChanged)

			u.repository.InsertHistory(c, userName, "Commissions", "UPDATE", rec.ID, "plan_commission",
				string(beforeJSONBytes), string(afterJSONBytes))
		}
	}

	return totalInserted, nil
}

// Error types
type notFoundError struct {
	Message string
}

func (e *notFoundError) Error() string {
	return e.Message
}

type validationError struct {
	Message string
}

func (e *validationError) Error() string {
	return e.Message
}


