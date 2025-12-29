package downloadallresult

import "github.com/labstack/echo"

// Usecase interface defines business logic operations
type Usecase interface {
	DownloadAllResult(c echo.Context, insuranceCode string) ([]byte, string, error)
}

// Repository interface defines data access operations
type Repository interface {
	// Insurances
	GetInsurances(c echo.Context, insuranceCode string) ([]InsuranceRow, error)

	// Products
	GetProducts(c echo.Context, insuranceCode string) ([]ProductRow, error)

	// Templates
	GetTemplates(c echo.Context, insuranceCode string) ([]TemplateRow, error)

	// Commissions
	GetCommissions(c echo.Context, insuranceCode string) ([]CommissionRow, error)
	GetPlanCommissions(c echo.Context, insuranceCode string) ([]PlanCommissionRow, error)
	GetDefaultConfigProducts(c echo.Context, insuranceCode string) ([]DefaultConfigProductRow, error)

	// Product Rules
	GetProductRules(c echo.Context, insuranceCode string) ([]ProductRuleRow, error)

	// Addons
	GetAddons(c echo.Context) ([]AddonRow, error)

	// Addon Rules
	GetAddonRules(c echo.Context, insuranceCode string) ([]AddonRuleRow, error)

	// Insurance Own Risks
	GetInsuranceOwnRisks(c echo.Context, insuranceCode string) ([]InsuranceOwnRiskRow, error)

	// Insurance Clauses
	GetInsuranceClauses(c echo.Context, insuranceCode string) ([]InsuranceClauseRow, error)
}

// Data models for each sheet
type InsuranceRow struct {
	ID        int
	Code      string
	Name      string
	Status    int
	Logo      string
	CreatedAt string
	UpdatedAt string
}

type ProductRow struct {
	ID                   int
	Code                 string
	InsuranceCode        string
	InsuranceProductCode string
	VehicleType          string
	IsInsuranceAPI       int
	Name                 string
	Logo                 string
	IsActive             int
	IsInsurerDriven      int
	IsPersonalUsage      int
	IsElectricVehicle    int
	Summary              string
	Config               string
	InsuranceDetail      string
	ProtectionDetail     string
	HowToClaim           string
	CreatedBy            int
	BaseProduct          string
	CreatedAt            string
	UpdatedAt            string
}

type TemplateRow struct {
	ID            int
	Locale        string
	TemplateID    string
	Value         string
	InsuranceCode string
	CreatedBy     int
	CreatedAt     string
}

type CommissionRow struct {
	ID              int
	ProductCode     string
	InsuranceCode   string
	AgentLevel      string
	CorporateID     int
	BasicCommission float64
	BonusPoint      float64
	Note            string
	StartDate       string
	EndDate         string
	CreatedAt       string
	UpdatedAt       string
}

type PlanCommissionRow struct {
	ID                   int
	PlanCode             string
	ProductID            int
	InsurerID            int
	CommissionPercentage float64
	CommissionVatType    string
	AFPercentage         float64
	AFVatType            string
	AdminFee             float64
	HardcopyFee          float64
	IsActive             bool
	Version              int
	CreatedAt            string
}

type DefaultConfigProductRow struct {
	ID          int
	Level       string
	Sequence    int
	Kind        string
	Category    string
	Insurer     string
	Product     string
	Selected    int
	Commission  float64
	Point       float64
	UplineBonus float64
	Bonus       float64
	MaxDiscount float64
	Order       int
	CreatedAt   string
	UpdatedAt   string
}

type ProductRuleRow struct {
	ProductCode             string
	InsuranceCode           string
	VehicleType             string
	VehicleCategory         string
	StartVehicleValue       float64
	EndVehicleValue         float64
	LimitVehicleAge         int
	AdminFee                float64
	HardcopyAdminFee        float64
	RegionID                int
	BasePremiumValue        float64
	LoadingFeePremiumValue  float64
	CommercialUsageValue    float64
	StartLoadingAge         float64
	AdditionalPremium       float64
	TypeAdditionalPremium   string
	Rules                   string
	CreatedBy               int
	CreatedAt               string
	BasePremiumType         int
	BaseLoadingPremiumValue float64
}

type AddonRow struct {
	ID        int
	Code      string
	Name      string
	IsActive  bool
	CreatedAt string
}

type AddonRuleRow struct {
	AddonCode     string
	ProductCode   string
	InsuranceCode string
	Rules         string
	CreatedBy     int
}

type InsuranceOwnRiskRow struct {
	InsuranceCode     string
	ProductType       string
	Code              string
	Title             string
	Value             float64
	ValueType         string
	Description       string
	IsMandatory       int
	IsActive          int
	IsElectricVehicle int
	CreatedBy         int
}

type InsuranceClauseRow struct {
	InsuranceCode string
	ProductType   string
	Code          string
	Title         string
	Content       string
	Description   string
	IsMandatory   int
	IsActive      int
	CreatedBy     int
}


