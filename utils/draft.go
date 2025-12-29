package utils

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo"
)

// DraftData represents the structure of draft_data.json
type DraftData struct {
	Insurances  []InsuranceDraft  `json:"insurances"`
	Products    []ProductDraft    `json:"products"`
	Rules       []RuleDraft       `json:"rules"`
	Templates   []TemplateDraft   `json:"templates"`
	AddonRules  []AddonRuleDraft  `json:"addon_rules"`
	OwnRisks    []OwnRisksDraft   `json:"own_risks"`
	Clauses     []ClausesDraft    `json:"clauses"`
	Commissions []CommissionDraft `json:"commissions"`
}

// InsuranceDraft represents insurance draft item
type InsuranceDraft struct {
	ID        int    `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Status    int    `json:"status"`
	Logo      string `json:"logo"`
	Timestamp string `json:"timestamp"`
}

// ProductDraft represents product draft item
type ProductDraft struct {
	ID                   int    `json:"id"`
	Code                 string `json:"code"`
	InsuranceCode        string `json:"insurance_code"`
	InsuranceProductCode string `json:"insurance_product_code"`
	VehicleType          string `json:"vehicle_type"`
	IsInsuranceAPI       int    `json:"is_insurance_api"`
	Name                 string `json:"name"`
	Logo                 string `json:"logo"`
	IsActive             int    `json:"is_active"`
	IsInsurerDriven      int    `json:"is_insurer_driven"`
	IsPersonalUsage      int    `json:"is_personal_usage"`
	IsElectricVehicle    int    `json:"is_electric_vehicle"`
	Summary              string `json:"summary"`
	Config               string `json:"config"`
	InsuranceDetail      string `json:"insurance_detail"`
	ProtectionDetail     string `json:"protection_detail"`
	HowToClaim           string `json:"how_to_claim"`
	CreatedBy            int    `json:"created_by"`
	BaseProduct          string `json:"base_product"`
	Timestamp            string `json:"timestamp"`
}

// RuleDraft represents rule draft item
type RuleDraft struct {
	ID                      int      `json:"id"`
	ProductCode             string   `json:"product_code"`
	InsuranceCode           string   `json:"insurance_code"`
	VehicleType             string   `json:"vehicle_type"`
	VehicleCategory         string   `json:"vehicle_category"`
	StartVehicleValue       *float64 `json:"start_vehicle_value,omitempty"`
	EndVehicleValue         *float64 `json:"end_vehicle_value,omitempty"`
	LimitVehicleAge         *int     `json:"limit_vehicle_age,omitempty"`
	AdminFee                *float64 `json:"admin_fee,omitempty"`
	HardcopyAdminFee        *float64 `json:"hardcopy_admin_fee,omitempty"`
	RegionID                *int     `json:"region_id,omitempty"`
	BasePremiumValue        *float64 `json:"base_premium_value,omitempty"`
	LoadingFeePremiumValue  *float64 `json:"loading_fee_premium_value,omitempty"`
	CommercialUsageValue    *float64 `json:"commercial_usage_value,omitempty"`
	StartLoadingAge         *int     `json:"start_loading_age,omitempty"`
	AdditionalPremium       *float64 `json:"additional_premium,omitempty"`
	TypeAdditionalPremium   string   `json:"type_additional_premium,omitempty"`
	Rules                   string   `json:"rules"`
	CreatedBy               int      `json:"created_by"`
	BasePremiumType         string   `json:"base_premium_type,omitempty"`
	BaseLoadingPremiumValue *float64 `json:"base_loading_premium_value,omitempty"`
	Timestamp               string   `json:"timestamp"`
}

// TemplateDraft represents template draft item
type TemplateDraft struct {
	Locale        string   `json:"locale"`
	TemplateID    string   `json:"template_id"`
	Value         []string `json:"value"` // Array of strings, default: ["","",""]
	InsuranceCode string   `json:"insurance_code,omitempty"`
	Insurance     string   `json:"insurance,omitempty"` // Alternative field name
	CreatedBy     int      `json:"created_by"`
	CreatedAt     string   `json:"created_at"` // CURRENT_TIMESTAMP
}

// AddonRuleDraft represents addon rule draft item
type AddonRuleDraft struct {
	ID            int    `json:"id"`
	AddonCode     string `json:"addon_code"`
	ProductCode   string `json:"product_code"`
	InsuranceCode string `json:"insurance_code"`
	Rules         string `json:"rules"`
	CreatedBy     int    `json:"created_by"`
	Timestamp     string `json:"timestamp"`
}

// OwnRisksDraft represents insurance own risks draft item
type OwnRisksDraft struct {
	ID                int    `json:"id"`
	InsuranceCode     string `json:"insurance_code"`
	ProductType       string `json:"product_type"`
	Code              string `json:"code"`
	Title             string `json:"title"`
	Value             string `json:"value"`
	ValueType         string `json:"value_type,omitempty"`
	Description       string `json:"description,omitempty"`
	IsMandatory       int    `json:"is_mandatory,omitempty"`
	IsActive          int    `json:"is_active,omitempty"`
	IsElectricVehicle int    `json:"is_electric_vehicle,omitempty"`
	CreatedBy         int    `json:"created_by"`
	Timestamp         string `json:"timestamp"`
}

// ClausesDraft represents insurance clauses draft item
type ClausesDraft struct {
	ID                int    `json:"id"`
	InsuranceCode     string `json:"insurance_code"`
	ProductType       string `json:"product_type"`
	Code              string `json:"code"`
	Title             string `json:"title"`
	Content           string `json:"content"`
	Description       string `json:"description,omitempty"`
	IsMandatory       int    `json:"is_mandatory,omitempty"`
	IsActive          int    `json:"is_active,omitempty"`
	IsElectricVehicle int    `json:"is_electric_vehicle,omitempty"`
	CreatedBy         int    `json:"created_by"`
	Timestamp         string `json:"timestamp"`
}

// CommissionDraft represents commission draft item
type CommissionDraft struct {
	ID              int                     `json:"id"`
	InsuranceCode   string                  `json:"insurance_code"`
	BasicCommission float64                 `json:"basic_commission"`
	PlanCommissions PlanCommissionDraftData `json:"plan_commissions"`
	CreatedBy       int                     `json:"created_by"`
	Timestamp       string                  `json:"timestamp"`
}

// PlanCommissionDraftData represents plan commission data in draft
type PlanCommissionDraftData struct {
	CommissionVatType string  `json:"commission_vat_type"`
	AfPercentage      float64 `json:"af_percentage"`
	AfVatType         string  `json:"af_vat_type"`
	AdminFee          float64 `json:"admin_fee"`
}

// GetDraftFilePath returns the absolute path to draft_data.json
func GetDraftFilePath() string {
	dir, err := os.Getwd()
	if err != nil {
		dir = "."
	}
	paths := []string{
		filepath.Join(dir, "draft_data.json"),
		filepath.Join(dir, "..", "draft_data.json"),
		filepath.Join("/tmp", "draft_data.json"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return filepath.Join(dir, "draft_data.json")
}

// LoadDraftData loads draft data from JSON file
func LoadDraftData(c echo.Context) (*DraftData, error) {
	data := &DraftData{
		Insurances:  []InsuranceDraft{},
		Products:    []ProductDraft{},
		Rules:       []RuleDraft{},
		Templates:   []TemplateDraft{},
		AddonRules:  []AddonRuleDraft{},
		OwnRisks:    []OwnRisksDraft{},
		Clauses:     []ClausesDraft{},
		Commissions: []CommissionDraft{},
	}

	draftFilePath := GetDraftFilePath()
	file, err := os.Open(draftFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(data); err != nil {
		log.Printf("WARNING: Failed to decode draft_data.json: %v", err)
		return data, nil
	}

	// Ensure all slices are initialized
	if data.Products == nil {
		data.Products = []ProductDraft{}
	}
	if data.Templates == nil {
		data.Templates = []TemplateDraft{}
	}
	if data.Insurances == nil {
		data.Insurances = []InsuranceDraft{}
	}
	if data.Rules == nil {
		data.Rules = []RuleDraft{}
	}
	if data.AddonRules == nil {
		data.AddonRules = []AddonRuleDraft{}
	}
	if data.OwnRisks == nil {
		data.OwnRisks = []OwnRisksDraft{}
	}
	if data.Clauses == nil {
		data.Clauses = []ClausesDraft{}
	}
	if data.Commissions == nil {
		data.Commissions = []CommissionDraft{}
	}

	return data, nil
}

// SaveDraftData saves draft data to JSON file
func SaveDraftData(c echo.Context, data *DraftData) error {
	draftFilePath := GetDraftFilePath()
	file, err := os.Create(draftFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// GetTimestamp returns current timestamp in format "YYYY-MM-DD HH:MM:SS"
func GetTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// GetInsuranceLogo gets logo from insurances table
func GetInsuranceLogo(c echo.Context, insuranceCode string) string {
	db := GetDB()
	if db == nil {
		return ""
	}

	var logo sql.NullString
	query := "SELECT logo FROM vehicle_service_development.insurances WHERE code = ?"
	err := db.QueryRow(query, insuranceCode).Scan(&logo)
	if err != nil {
		return ""
	}

	if logo.Valid {
		return logo.String
	}
	return ""
}

// GetNextID gets the next ID for insurance draft
func GetNextID(data *DraftData) int {
	maxID := 0
	for _, item := range data.Insurances {
		if item.ID > maxID {
			maxID = item.ID
		}
	}
	return maxID + 1
}

// GetNextOwnRisksID gets the next ID for own risks draft
func GetNextOwnRisksID(data *DraftData) int {
	maxID := 0
	for _, item := range data.OwnRisks {
		if item.ID > maxID {
			maxID = item.ID
		}
	}
	return maxID + 1
}

// GetNextRuleID gets the next ID for rule draft
func GetNextRuleID(data *DraftData) int {
	maxID := 0
	for _, item := range data.Rules {
		if item.ID > maxID {
			maxID = item.ID
		}
	}
	return maxID + 1
}

// GetNextClausesID gets the next ID for clauses draft
func GetNextClausesID(data *DraftData) int {
	maxID := 0
	for _, item := range data.Clauses {
		if item.ID > maxID {
			maxID = item.ID
		}
	}
	return maxID + 1
}

// GetNextAddonRuleID gets the next ID for addon rule draft
func GetNextAddonRuleID(data *DraftData) int {
	maxID := 0
	for _, item := range data.AddonRules {
		if item.ID > maxID {
			maxID = item.ID
		}
	}
	return maxID + 1
}

