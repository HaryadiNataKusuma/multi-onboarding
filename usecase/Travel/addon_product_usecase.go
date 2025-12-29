package usecase

import (
	"errors"
	"fmt"
	"multi-onboarding/model/Travel"
	"multi-onboarding/repository/Travel"
)

// AddonProductUsecase handles business logic for saving addons for products
type AddonProductUsecase struct {
	addonRepo                    repository.AddonRepository
	addonRuleRepo                repository.AddonRuleRepository
	insuranceProductAddonMappingRepo repository.InsuranceProductAddonMappingRepository
}

// NewAddonProductUsecase creates a new addon product usecase
func NewAddonProductUsecase(
	addonRepo repository.AddonRepository,
	addonRuleRepo repository.AddonRuleRepository,
	insuranceProductAddonMappingRepo repository.InsuranceProductAddonMappingRepository,
) *AddonProductUsecase {
	return &AddonProductUsecase{
		addonRepo:                    addonRepo,
		addonRuleRepo:                addonRuleRepo,
		insuranceProductAddonMappingRepo: insuranceProductAddonMappingRepo,
	}
}

// SaveProductAddons saves addons for a product to all 3 tables
func (u *AddonProductUsecase) SaveProductAddons(req *model.AddonRequest) error {
	createdBy := int64(1) // Default created_by
	
	// Remove duplicates from req.Addons before processing
	// Create a map to track unique addon codes
	uniqueAddons := make(map[string]*model.AddonWithNames)
	for _, addon := range req.Addons {
		if addon.Code != "" {
			// If we haven't seen this addon code before, add it
			if _, exists := uniqueAddons[addon.Code]; !exists {
				uniqueAddons[addon.Code] = &addon
			}
		}
	}
	
	// Convert map back to slice
	req.Addons = make([]model.AddonWithNames, 0, len(uniqueAddons))
	for _, addon := range uniqueAddons {
		req.Addons = append(req.Addons, *addon)
	}
	
	// Delete existing addon_rules and insurance_product_addon_mappings for this product_code first
	// This prevents duplicates when editing
	// Use transaction to ensure atomicity
	if err := u.addonRuleRepo.DeleteByProductCode(req.ProductCode); err != nil {
		return fmt.Errorf("failed to delete existing addon rules: %w", err)
	}
	
	if err := u.insuranceProductAddonMappingRepo.DeleteByProductCode(req.ProductCode); err != nil {
		return fmt.Errorf("failed to delete existing insurance product addon mappings: %w", err)
	}
	
	// Log for debugging
	fmt.Printf("DEBUG: Deleted existing addon rules for product_code: %s\n", req.ProductCode)
	fmt.Printf("DEBUG: About to insert %d unique addon rules for product_code: %s\n", len(req.Addons), req.ProductCode)
	
	// Save to addons table (only if not exists)
	for _, addonWithNames := range req.Addons {
		// Check if addon already exists
		existingAddon, err := u.addonRepo.GetByCode(addonWithNames.Code)
		if err != nil {
			// If error is not "not found", return the error
			if !errors.Is(err, repository.ErrAddonNotFound) {
				return err
			}
			// If addon not found, proceed to create it
		}
		
		// Only create if addon doesn't exist
		if existingAddon == nil {
			addon := &model.Addon{
				Code:      addonWithNames.Code,
				Name:      addonWithNames.Name,
				NameMy:    addonWithNames.NameMy,
				NameEn:    addonWithNames.NameEn,
				IsActive:  1,
				CreatedBy: createdBy,
			}
			
			if err := u.addonRepo.Create(addon); err != nil {
				return err
			}
		}
		// If addon already exists, skip insertion
	}
	
	// Save to addon_rules table
	addonRules := make([]model.AddonRule, 0, len(req.Addons))
	for _, addonWithNames := range req.Addons {
		insuranceCode := req.InsuranceCode
		rulesJSON := "{}"
		addonRule := model.AddonRule{
			ProductCode:   req.ProductCode,
			AddonCode:     addonWithNames.Code,
			InsuranceCode: &insuranceCode,
			Rules:         &rulesJSON,
			CreatedBy:     createdBy,
		}
		addonRules = append(addonRules, addonRule)
	}
	
	if len(addonRules) > 0 {
		if err := u.addonRuleRepo.CreateBatch(addonRules); err != nil {
			return err
		}
	}
	
	// Save to insurance_product_addon_mappings table
	mappings := make([]model.InsuranceProductAddonMapping, 0, len(req.Addons))
	for _, addonWithNames := range req.Addons {
		mapping := model.InsuranceProductAddonMapping{
			InsuranceCode: req.InsuranceCode,
			ProductCode:   req.ProductCode,
			AddonCode:     addonWithNames.Code,
			CreatedBy:     createdBy,
		}
		mappings = append(mappings, mapping)
	}
	
	if len(mappings) > 0 {
		if err := u.insuranceProductAddonMappingRepo.CreateBatch(mappings); err != nil {
			return err
		}
	}
	
	return nil
}

// CleanupDuplicateAddonRules removes duplicate addon_rules entries (keeping only the first one)
func (u *AddonProductUsecase) CleanupDuplicateAddonRules() error {
	// This will be implemented in repository
	return u.addonRuleRepo.CleanupDuplicates()
}

