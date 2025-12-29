package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	// Domain imports - Clean Architecture (Travel)
	travelAddonHttp "multi-onboarding/domain/Travel/addon/delivery/http"
	travelAddonRepo "multi-onboarding/domain/Travel/addon/repository"
	travelAddonUsecase "multi-onboarding/domain/Travel/addon/usecase"
	travelCommissionHttp "multi-onboarding/domain/Travel/commission/delivery/http"
	travelCommissionRepo "multi-onboarding/domain/Travel/commission/repository"
	travelCommissionUsecase "multi-onboarding/domain/Travel/commission/usecase"
	travelCountryHttp "multi-onboarding/domain/Travel/country/delivery/http"
	travelCountryRepo "multi-onboarding/domain/Travel/country/repository"
	travelCountryUsecase "multi-onboarding/domain/Travel/country/usecase"
	travelHistoryHttp "multi-onboarding/domain/Travel/history/delivery/http"
	travelHistoryRepo "multi-onboarding/domain/Travel/history/repository"
	travelHistoryUsecase "multi-onboarding/domain/Travel/history/usecase"
	travelInsuranceHttp "multi-onboarding/domain/Travel/insurance/delivery/http"
	travelInsuranceRepo "multi-onboarding/domain/Travel/insurance/repository"
	travelInsuranceUsecase "multi-onboarding/domain/Travel/insurance/usecase"
	travelProductHttp "multi-onboarding/domain/Travel/product/delivery/http"
	travelProductRepo "multi-onboarding/domain/Travel/product/repository"
	travelProductUsecase "multi-onboarding/domain/Travel/product/usecase"
	travelRegionHttp "multi-onboarding/domain/Travel/region/delivery/http"
	travelRegionRepo "multi-onboarding/domain/Travel/region/repository"
	travelRegionUsecase "multi-onboarding/domain/Travel/region/usecase"
	travelTemplateHttp "multi-onboarding/domain/Travel/template/delivery/http"
	travelTemplateRepo "multi-onboarding/domain/Travel/template/repository"
	travelTemplateUsecase "multi-onboarding/domain/Travel/template/usecase"

	// Property domain imports
	propertyInsuranceHttp "multi-onboarding/domain/Property/insurance/delivery/http"
	propertyInsuranceRepo "multi-onboarding/domain/Property/insurance/repository"
	propertyInsuranceUsecase "multi-onboarding/domain/Property/insurance/usecase"
	propertyProductHttp "multi-onboarding/domain/Property/product/delivery/http"
	propertyProductRepo "multi-onboarding/domain/Property/product/repository"
	propertyProductUsecase "multi-onboarding/domain/Property/product/usecase"
	propertyTemplateHttp "multi-onboarding/domain/Property/template/delivery/http"
	propertyTemplateRepo "multi-onboarding/domain/Property/template/repository"
	propertyTemplateUsecase "multi-onboarding/domain/Property/template/usecase"

	// Property product rules handler
	propertyProductRuleHttp "multi-onboarding/delivery/http/property"

	// Vehicle domain imports
	vehicleAddonHttp "multi-onboarding/domain/vehicle/addon/delivery/http"
	vehicleAddonRepo "multi-onboarding/domain/vehicle/addon/repository"
	vehicleAddonUsecase "multi-onboarding/domain/vehicle/addon/usecase"
	vehicleAddonRulesHttp "multi-onboarding/domain/vehicle/addonrules/delivery/http"
	vehicleAddonRulesRepo "multi-onboarding/domain/vehicle/addonrules/repository"
	vehicleAddonRulesUsecase "multi-onboarding/domain/vehicle/addonrules/usecase"
	vehicleCommissionHttp "multi-onboarding/domain/vehicle/commission/delivery/http"
	vehicleCommissionRepo "multi-onboarding/domain/vehicle/commission/repository"
	vehicleCommissionUsecase "multi-onboarding/domain/vehicle/commission/usecase"
	vehicleDownloadAllResultHttp "multi-onboarding/domain/vehicle/downloadallresult/delivery/http"
	vehicleDownloadAllResultRepo "multi-onboarding/domain/vehicle/downloadallresult/repository"
	vehicleDownloadAllResultUsecase "multi-onboarding/domain/vehicle/downloadallresult/usecase"
	vehicleHistoryHttp "multi-onboarding/domain/vehicle/history/delivery/http"
	vehicleHistoryRepo "multi-onboarding/domain/vehicle/history/repository"
	vehicleHistoryUsecase "multi-onboarding/domain/vehicle/history/usecase"
	vehicleHTMLConverterHttp "multi-onboarding/domain/vehicle/htmlconverter/delivery/http"
	vehicleInsuranceHttp "multi-onboarding/domain/vehicle/insurance/delivery/http"
	vehicleInsuranceRepo "multi-onboarding/domain/vehicle/insurance/repository"
	vehicleInsuranceUsecase "multi-onboarding/domain/vehicle/insurance/usecase"
	vehicleInsuranceClausesHttp "multi-onboarding/domain/vehicle/insuranceclauses/delivery/http"
	vehicleInsuranceClausesRepo "multi-onboarding/domain/vehicle/insuranceclauses/repository"
	vehicleInsuranceClausesUsecase "multi-onboarding/domain/vehicle/insuranceclauses/usecase"
	vehicleInsuranceOwnRisksHttp "multi-onboarding/domain/vehicle/insuranceownrisks/delivery/http"
	vehicleInsuranceOwnRisksRepo "multi-onboarding/domain/vehicle/insuranceownrisks/repository"
	vehicleInsuranceOwnRisksUsecase "multi-onboarding/domain/vehicle/insuranceownrisks/usecase"
	vehicleProductHttp "multi-onboarding/domain/vehicle/product/delivery/http"
	vehicleProductRepo "multi-onboarding/domain/vehicle/product/repository"
	vehicleProductUsecase "multi-onboarding/domain/vehicle/product/usecase"
	vehicleRulesHttp "multi-onboarding/domain/vehicle/rules/delivery/http"
	vehicleRulesRepo "multi-onboarding/domain/vehicle/rules/repository"
	vehicleRulesUsecase "multi-onboarding/domain/vehicle/rules/usecase"
	vehicleTemplateHttp "multi-onboarding/domain/vehicle/template/delivery/http"
	vehicleTemplateRepo "multi-onboarding/domain/vehicle/template/repository"
	vehicleTemplateUsecase "multi-onboarding/domain/vehicle/template/usecase"

	// Legacy imports for domains not yet migrated
	httphandler "multi-onboarding/delivery/http/travel"
	repository "multi-onboarding/repository/Travel"
	usecase "multi-onboarding/usecase/Travel"
	"multi-onboarding/utils"
)

// ProgramConfig holds configuration for each onboarding program
type ProgramConfig struct {
	Enabled bool
	Name    string
	Prefix  string // URL prefix like "vehicle", "travel", "property", "health"
}

// ProgramsConfig holds all program configurations
type ProgramsConfig struct {
	DefaultProgram string // Default program for backward compatibility routes
	Programs       map[string]ProgramConfig
}

// getDefaultProgram returns the default program from environment or "vehicle"
func getDefaultProgram() string {
	if defaultProg := os.Getenv("DEFAULT_ONBOARDING_PROGRAM"); defaultProg != "" {
		return defaultProg
	}
	return "vehicle" // Default to vehicle for MV onboarding
}

// getProgramsConfig returns configuration for all programs
func getProgramsConfig() ProgramsConfig {
	defaultProg := getDefaultProgram()

	return ProgramsConfig{
		DefaultProgram: defaultProg,
		Programs: map[string]ProgramConfig{
			"vehicle": {
				Enabled: true,
				Name:    "Motor Vehicle",
				Prefix:  "vehicle",
			},
			"travel": {
				Enabled: true,
				Name:    "Travel",
				Prefix:  "travel",
			},
			// Future programs - can be enabled when ready
			"property": {
				Enabled: true, // Property domain is ready
				Name:    "Property",
				Prefix:  "property",
			},
			"health": {
				Enabled: false, // Set to true when health domain is ready
				Name:    "Health",
				Prefix:  "health",
			},
		},
	}
}

func main() {
	// Get programs configuration
	programsConfig := getProgramsConfig()

	// Database configuration
	dbConfig := repository.DBConfig{
		Host:     "127.0.0.1",
		User:     "root",
		Password: "Qoala123**",
		Database: "travel_service_development",
		Port:     3306,
	}

	// Initialize database connection using utils
	log.Println("🔌 Connecting to MySQL database...")
	log.Printf("   Host: %s:%d", dbConfig.Host, dbConfig.Port)
	log.Printf("   Database: %s", dbConfig.Database)
	log.Printf("   User: %s", dbConfig.User)

	utils.InitDB(dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Database, dbConfig.Port)
	defer utils.CloseDB()

	log.Println("✅ Successfully connected to MySQL database")

	// Initialize database schema
	db := utils.GetDB()
	if err := repository.InitSchema(db); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	log.Println("✅ Database schema initialized")

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:5173", // Frontend Travel
			"http://localhost:5174", // Frontend Property
			"http://localhost:3000",
			"http://localhost:3333",
			"http://localhost:8080",
			"http://127.0.0.1:5173",
			"http://127.0.0.1:5174", // Frontend Property
			"http://127.0.0.1:3333",
			"http://127.0.0.1:8080",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-User-Name", "X-User-Email", "X-Requested-With"},
		AllowCredentials: true,
	}))

	// Initialize all domain repositories, usecases, and handlers
	// This will be done for each program that is enabled

	// ============================================
	// VEHICLE PROGRAM SETUP
	// ============================================
	if programsConfig.Programs["vehicle"].Enabled {
		log.Println("🚗 Initializing Vehicle onboarding program...")
		setupVehicleProgram(e, programsConfig)
	}

	// ============================================
	// TRAVEL PROGRAM SETUP
	// ============================================
	if programsConfig.Programs["travel"].Enabled {
		log.Println("✈️  Initializing Travel onboarding program...")
		setupTravelProgram(e, programsConfig)
	}

	// ============================================
	// PROPERTY PROGRAM SETUP
	// ============================================
	if programsConfig.Programs["property"].Enabled {
		log.Println("🏠 Initializing Property onboarding program...")
		setupPropertyProgram(e, programsConfig)
	}

	// ============================================
	// HEALTH PROGRAM SETUP (Future)
	// ============================================
	if programsConfig.Programs["health"].Enabled {
		log.Println("🏥 Initializing Health onboarding program...")
		// setupHealthProgram(e, programsConfig)
		// TODO: Implement when health domain is ready
	}

	// ============================================
	// LEGACY ROUTES (Shared across programs)
	// ============================================
	setupLegacyRoutes(e)

	// ============================================
	// SERVER STARTUP
	// ============================================
	log.Println("🚀 Server starting on :8080")
	log.Println("✅ Using Echo v3.3.10")
	log.Println("")
	log.Printf("📋 Default Program: %s (for backward compatibility)", programsConfig.DefaultProgram)
	log.Println("📋 Enabled Programs:")
	for _, config := range programsConfig.Programs {
		if config.Enabled {
			log.Printf("   ✅ %s - Routes: /api/%s/*", config.Name, config.Prefix)
		}
	}
	log.Println("")
	log.Println("📡 Server is ready! Access at http://localhost:8080")
	log.Println("📋 Available endpoints:")
	log.Printf("   - GET  /api/%s/* (default program routes)", programsConfig.DefaultProgram)
	for _, config := range programsConfig.Programs {
		if config.Enabled {
			log.Printf("   - GET  /api/%s/* (%s program routes)", config.Prefix, config.Name)
		}
	}
	log.Println("")

	log.Println("⏳ Starting HTTP server on port 8080...")

	// Check if port is already in use before starting
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
		log.Fatalf("💡 Tip: Check if port 8080 is already in use with: lsof -i:8080")
		log.Fatalf("💡 Tip: Kill existing process with: lsof -ti:8080 | xargs kill -9")
	}
}

// setupVehicleProgram sets up all routes for Vehicle onboarding program
func setupVehicleProgram(e *echo.Echo, config ProgramsConfig) {
	// Initialize Vehicle domain repositories, usecases, and handlers
	vehicleAddonRepository := vehicleAddonRepo.NewAddonRepository()
	vehicleAddonUsecaseInstance := vehicleAddonUsecase.NewAddonUsecase(vehicleAddonRepository)
	vehicleAddonHandler := vehicleAddonHttp.NewAddonHandler(vehicleAddonUsecaseInstance)

	vehicleInsuranceRepository := vehicleInsuranceRepo.NewInsuranceRepository()
	vehicleInsuranceUsecaseInstance := vehicleInsuranceUsecase.NewInsuranceUsecase(vehicleInsuranceRepository)
	vehicleInsuranceHandler := vehicleInsuranceHttp.NewInsuranceHandler(vehicleInsuranceUsecaseInstance)

	vehicleProductRepository := vehicleProductRepo.NewProductRepository()
	vehicleProductUsecaseInstance := vehicleProductUsecase.NewProductUsecase(vehicleProductRepository)
	vehicleProductHandler := vehicleProductHttp.NewProductHandler(vehicleProductUsecaseInstance)

	vehicleTemplateRepository := vehicleTemplateRepo.NewTemplateRepository()
	vehicleTemplateUsecaseInstance := vehicleTemplateUsecase.NewTemplateUsecase(vehicleTemplateRepository)

	vehicleHistoryRepository := vehicleHistoryRepo.NewHistoryRepository()
	vehicleHistoryUsecaseInstance := vehicleHistoryUsecase.NewHistoryUsecase(vehicleHistoryRepository)
	vehicleHistoryHandler := vehicleHistoryHttp.NewHistoryHandler(vehicleHistoryUsecaseInstance)

	vehicleAddonRulesRepository := vehicleAddonRulesRepo.NewAddonRulesRepository()
	vehicleAddonRulesUsecaseInstance := vehicleAddonRulesUsecase.NewAddonRulesUsecase(vehicleAddonRulesRepository)

	vehicleCommissionRepository := vehicleCommissionRepo.NewCommissionRepository()
	vehicleCommissionUsecaseInstance := vehicleCommissionUsecase.NewCommissionUsecase(vehicleCommissionRepository)

	vehicleRulesRepository := vehicleRulesRepo.NewRulesRepository()
	vehicleRulesUsecaseInstance := vehicleRulesUsecase.NewRulesUsecase(vehicleRulesRepository)

	vehicleInsuranceClausesRepository := vehicleInsuranceClausesRepo.NewInsuranceClausesRepository()
	vehicleInsuranceClausesUsecaseInstance := vehicleInsuranceClausesUsecase.NewInsuranceClausesUsecase(vehicleInsuranceClausesRepository)

	vehicleInsuranceOwnRisksRepository := vehicleInsuranceOwnRisksRepo.NewInsuranceOwnRisksRepository()
	vehicleInsuranceOwnRisksUsecaseInstance := vehicleInsuranceOwnRisksUsecase.NewInsuranceOwnRisksUsecase(vehicleInsuranceOwnRisksRepository)

	vehicleDownloadAllResultRepository := vehicleDownloadAllResultRepo.NewDownloadAllResultRepository()
	vehicleDownloadAllResultUsecaseInstance := vehicleDownloadAllResultUsecase.NewDownloadAllResultUsecase(vehicleDownloadAllResultRepository)
	vehicleDownloadAllResultHandler := vehicleDownloadAllResultHttp.NewDownloadAllResultHandler(vehicleDownloadAllResultUsecaseInstance)

	vehicleHTMLConverterHandler := vehicleHTMLConverterHttp.NewHTMLConverterHandler()

	// API routes group
	api := e.Group("/api")
	vehicle := api.Group("/vehicle")

	// Vehicle Insurance routes
	vehicle.POST("/insurances/draft/add", vehicleInsuranceHandler.SaveInsuranceDraft)
	vehicle.GET("/insurances/draft", vehicleInsuranceHandler.GetInsurancesDraft)
	vehicle.PUT("/insurances/draft/:id", vehicleInsuranceHandler.UpdateInsuranceDraft)
	vehicle.DELETE("/insurances/draft/:id", vehicleInsuranceHandler.DeleteInsuranceDraft)
	vehicle.POST("/insurances/draft/clear", vehicleInsuranceHandler.ClearInsurancesDraft)
	vehicle.POST("/insurances/draft/confirm", vehicleInsuranceHandler.ConfirmInsurances)
	vehicle.GET("/insurances", vehicleInsuranceHandler.GetInsurances)
	vehicle.PUT("/insurances/:id", vehicleInsuranceHandler.UpdateInsurance)

	// Vehicle Product routes
	vehicle.GET("/products", vehicleProductHandler.GetProducts)
	vehicle.POST("/products", vehicleProductHandler.CreateProduct)
	vehicle.GET("/products/draft", vehicleProductHandler.GetAllDrafts)
	vehicle.POST("/products/draft/add", vehicleProductHandler.AddDraft)
	vehicle.POST("/products/draft/clear", vehicleProductHandler.ClearDrafts)
	vehicle.POST("/products/draft/confirm", vehicleProductHandler.ConfirmDrafts)
	vehicle.DELETE("/products/draft/:id", vehicleProductHandler.DeleteDraft)
	vehicle.GET("/products/:id", vehicleProductHandler.GetProduct)
	vehicle.PUT("/products/:id", vehicleProductHandler.UpdateProduct)
	vehicle.DELETE("/products/:id", vehicleProductHandler.DeleteProduct)

	// Vehicle Addon routes
	vehicle.GET("/addons", vehicleAddonHandler.GetAddons)
	vehicle.PUT("/addons/:id", vehicleAddonHandler.UpdateAddon)
	vehicle.POST("/addons/confirm", vehicleAddonHandler.ConfirmAddons)

	// Vehicle History routes
	vehicle.GET("/histories", vehicleHistoryHandler.GetHistories)

	// Vehicle routes using AddHandler pattern
	vehicleAddonRulesHttp.AddAddonRulesHandler(e, vehicleAddonRulesUsecaseInstance)
	vehicleCommissionHttp.AddCommissionHandler(e, vehicleCommissionUsecaseInstance)
	vehicleRulesHttp.AddRulesHandler(e, vehicleRulesUsecaseInstance)
	vehicleTemplateHttp.AddTemplateHandler(e, vehicleTemplateUsecaseInstance)
	vehicleInsuranceClausesHttp.AddInsuranceClausesHandler(e, vehicleInsuranceClausesUsecaseInstance)
	vehicleInsuranceOwnRisksHttp.AddInsuranceOwnRisksHandler(e, vehicleInsuranceOwnRisksUsecaseInstance)

	// Vehicle Download All Result route
	vehicle.GET("/download-all-result", vehicleDownloadAllResultHandler.DownloadAllResult)

	// Vehicle HTML Converter route
	vehicle.POST("/html/convert", vehicleHTMLConverterHandler.ConvertHTML)

	// Backward compatibility: If vehicle is default program, map /api/* to /api/vehicle/*
	if config.DefaultProgram == "vehicle" {
		setupVehicleBackwardCompatibility(api, vehicleInsuranceHandler, vehicleProductHandler, vehicleAddonHandler, vehicleHistoryHandler)
	}
}

// setupVehicleBackwardCompatibility sets up backward compatibility routes for vehicle
func setupVehicleBackwardCompatibility(api *echo.Group, vehicleInsuranceHandler interface{}, vehicleProductHandler interface{}, vehicleAddonHandler interface{}, vehicleHistoryHandler interface{}) {
	// Insurance routes (backward compatibility) - using echo.HandlerFunc directly
	api.POST("/insurances/draft/add", func(c echo.Context) error {
		return vehicleInsuranceHandler.(interface{ SaveInsuranceDraft(echo.Context) error }).SaveInsuranceDraft(c)
	})
	api.GET("/insurances/draft", func(c echo.Context) error {
		return vehicleInsuranceHandler.(interface{ GetInsurancesDraft(echo.Context) error }).GetInsurancesDraft(c)
	})
	api.PUT("/insurances/draft/:id", func(c echo.Context) error {
		return vehicleInsuranceHandler.(interface{ UpdateInsuranceDraft(echo.Context) error }).UpdateInsuranceDraft(c)
	})
	api.DELETE("/insurances/draft/:id", func(c echo.Context) error {
		return vehicleInsuranceHandler.(interface{ DeleteInsuranceDraft(echo.Context) error }).DeleteInsuranceDraft(c)
	})
	api.POST("/insurances/draft/clear", func(c echo.Context) error {
		return vehicleInsuranceHandler.(interface{ ClearInsurancesDraft(echo.Context) error }).ClearInsurancesDraft(c)
	})
	api.POST("/insurances/draft/confirm", func(c echo.Context) error {
		return vehicleInsuranceHandler.(interface{ ConfirmInsurances(echo.Context) error }).ConfirmInsurances(c)
	})
	api.GET("/insurances", func(c echo.Context) error {
		return vehicleInsuranceHandler.(interface{ GetInsurances(echo.Context) error }).GetInsurances(c)
	})
	api.PUT("/insurances/:id", func(c echo.Context) error {
		return vehicleInsuranceHandler.(interface{ UpdateInsurance(echo.Context) error }).UpdateInsurance(c)
	})

	// Product routes (backward compatibility)
	prodHandler := vehicleProductHandler.(interface {
		GetProducts(echo.Context) error
		CreateProduct(echo.Context) error
		GetAllDrafts(echo.Context) error
		AddDraft(echo.Context) error
		ClearDrafts(echo.Context) error
		ConfirmDrafts(echo.Context) error
		GetDraftByID(echo.Context) error
		UpdateDraft(echo.Context) error
		DeleteDraft(echo.Context) error
		GetProduct(echo.Context) error
		UpdateProduct(echo.Context) error
		DeleteProduct(echo.Context) error
	})
	api.GET("/products", prodHandler.GetProducts)
	api.POST("/products", prodHandler.CreateProduct)
	api.GET("/products/draft", prodHandler.GetAllDrafts)
	api.POST("/products/draft/add", prodHandler.AddDraft)
	api.POST("/products/draft/clear", prodHandler.ClearDrafts)
	api.POST("/products/draft/confirm", prodHandler.ConfirmDrafts)
	api.GET("/products/draft/:id", prodHandler.GetDraftByID)
	api.PUT("/products/draft/:id", prodHandler.UpdateDraft)
	api.DELETE("/products/draft/:id", prodHandler.DeleteDraft)
	api.GET("/products/:id", prodHandler.GetProduct)
	api.PUT("/products/:id", prodHandler.UpdateProduct)
	api.DELETE("/products/:id", prodHandler.DeleteProduct)

	// Addon routes (backward compatibility)
	addonHandler := vehicleAddonHandler.(interface {
		GetAddons(echo.Context) error
		UpdateAddon(echo.Context) error
		ConfirmAddons(echo.Context) error
	})
	api.GET("/addons", addonHandler.GetAddons)
	api.PUT("/addons/:id", addonHandler.UpdateAddon)
	api.POST("/addons/confirm", addonHandler.ConfirmAddons)

	// History routes (backward compatibility)
	histHandler := vehicleHistoryHandler.(interface {
		GetHistories(echo.Context) error
	})
	api.GET("/histories", histHandler.GetHistories)
}

// setupTravelProgram sets up all routes for Travel onboarding program
func setupTravelProgram(e *echo.Echo, config ProgramsConfig) {
	// Initialize Travel domain repositories, usecases, and handlers
	travelAddonRepository := travelAddonRepo.NewAddonRepository()
	travelAddonUsecaseInstance := travelAddonUsecase.NewAddonUsecase(travelAddonRepository)
	travelAddonHandler := travelAddonHttp.NewAddonHandler(travelAddonUsecaseInstance)

	travelInsuranceRepository := travelInsuranceRepo.NewInsuranceRepository()
	travelInsuranceUsecaseInstance := travelInsuranceUsecase.NewInsuranceUsecase(travelInsuranceRepository)
	travelInsuranceHandler := travelInsuranceHttp.NewInsuranceHandler(travelInsuranceUsecaseInstance)

	travelCountryRepository := travelCountryRepo.NewCountryRepository()
	travelCountryUsecaseInstance := travelCountryUsecase.NewCountryUsecase(travelCountryRepository)
	travelCountryHandler := travelCountryHttp.NewCountryHandler(travelCountryUsecaseInstance)

	travelHistoryRepository := travelHistoryRepo.NewHistoryRepository()
	travelHistoryUsecaseInstance := travelHistoryUsecase.NewHistoryUsecase(travelHistoryRepository)
	travelHistoryHandler := travelHistoryHttp.NewHistoryHandler(travelHistoryUsecaseInstance)

	travelRegionRepository := travelRegionRepo.NewRegionRepository()
	travelRegionUsecaseInstance := travelRegionUsecase.NewRegionUsecase(travelRegionRepository)
	travelRegionHandler := travelRegionHttp.NewRegionHandler(travelRegionUsecaseInstance)

	travelProductRepository := travelProductRepo.NewProductRepository()
	travelProductUsecaseInstance := travelProductUsecase.NewProductUsecase(travelProductRepository)
	travelProductHandler := travelProductHttp.NewProductHandler(travelProductUsecaseInstance)

	travelTemplateRepository := travelTemplateRepo.NewTemplateRepository()
	travelTemplateUsecaseInstance := travelTemplateUsecase.NewTemplateUsecase(travelTemplateRepository)
	travelTemplateHandler := travelTemplateHttp.NewTemplateHandler(travelTemplateUsecaseInstance)

	travelCommissionRepository := travelCommissionRepo.NewCommissionRepository()
	travelCommissionUsecaseInstance := travelCommissionUsecase.NewCommissionUsecase(travelCommissionRepository)
	travelCommissionHandler := travelCommissionHttp.NewCommissionHandler(travelCommissionUsecaseInstance)

	// Initialize Travel Product Rules handler (using legacy handler for travel database)
	db := utils.GetDB()
	travelProductRuleRepo := repository.NewMySQLProductRuleRepository(db)
	travelProductRuleUsecase := usecase.NewProductRuleUsecase(travelProductRuleRepo)
	travelProductRuleHandler := httphandler.NewProductRuleHandler(travelProductRuleUsecase, db)

	// Initialize Travel Addon Rules handler (using legacy handler for travel database)
	travelAddonRuleRepo := repository.NewMySQLAddonRuleRepository(db)
	travelAddonRuleUsecase := usecase.NewAddonRuleUsecase(travelAddonRuleRepo)
	travelAddonRuleHandler := httphandler.NewAddonRuleHandler(travelAddonRuleUsecase)

	// Helper function to adapt http.HandlerFunc to echo.HandlerFunc
	adaptHandler := func(h http.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), "echo_context", c)
			r := c.Request().WithContext(ctx)
			h(c.Response(), r)
			return nil
		}
	}

	// API routes group
	api := e.Group("/api")
	travel := api.Group("/travel")

	// Travel Product routes
	travel.GET("/products", travelProductHandler.GetProducts)
	travel.POST("/products", travelProductHandler.CreateProduct)
	travel.GET("/products/draft", travelProductHandler.GetAllDrafts)
	travel.POST("/products/draft/add", travelProductHandler.AddDraft)
	travel.POST("/products/draft/clear", travelProductHandler.ClearDrafts)
	travel.POST("/products/draft/confirm", travelProductHandler.ConfirmDrafts)
	travel.DELETE("/products/draft/:id", travelProductHandler.DeleteDraft)
	travel.POST("/products/generate-templates", travelProductHandler.GenerateTemplateDraftsForAllProducts)
	travel.POST("/products/generate-templates/selected", travelProductHandler.GenerateTemplateDraftsForProducts)
	travel.GET("/products/:id", travelProductHandler.GetProduct)
	travel.PUT("/products/:id", travelProductHandler.UpdateProduct)
	travel.DELETE("/products/:id", travelProductHandler.DeleteProduct)

	// Travel Region routes
	travel.GET("/regions", travelRegionHandler.GetRegions)
	travel.POST("/regions", travelRegionHandler.CreateRegion)
	travel.GET("/regions/:id", travelRegionHandler.GetRegion)
	travel.PUT("/regions/:id", travelRegionHandler.UpdateRegion)
	travel.DELETE("/regions/:id", travelRegionHandler.DeleteRegion)

	// Travel Country routes
	travel.GET("/countries", travelCountryHandler.GetCountries)
	travel.GET("/countries/:id", travelCountryHandler.GetCountry)

	// Travel Insurance routes
	travel.GET("/insurances", travelInsuranceHandler.GetAll)
	travel.POST("/insurances", travelInsuranceHandler.Create)
	travel.GET("/insurances/:id", travelInsuranceHandler.GetByID)
	travel.PUT("/insurances/:id", travelInsuranceHandler.Update)
	travel.DELETE("/insurances/:id", travelInsuranceHandler.Delete)
	travel.GET("/insurances/draft", travelInsuranceHandler.GetAllDrafts)
	travel.POST("/insurances/draft/add", travelInsuranceHandler.AddDraft)
	travel.GET("/insurances/draft/:id", travelInsuranceHandler.GetDraftByID)
	travel.PUT("/insurances/draft/:id", travelInsuranceHandler.UpdateDraft)
	travel.DELETE("/insurances/draft/:id", travelInsuranceHandler.DeleteDraft)
	travel.POST("/insurances/draft/clear", travelInsuranceHandler.ClearDrafts)
	travel.POST("/insurances/draft/confirm", travelInsuranceHandler.ConfirmDrafts)

	// Travel Commission routes
	travel.GET("/commissions", travelCommissionHandler.GetAll)
	travel.POST("/commissions", travelCommissionHandler.Create)
	travel.GET("/commissions/:id", travelCommissionHandler.GetByID)
	travel.PUT("/commissions/:id", travelCommissionHandler.Update)
	travel.DELETE("/commissions/:id", travelCommissionHandler.Delete)
	travel.GET("/commissions/product/:product_code", travelCommissionHandler.GetByProductCode)
	travel.GET("/commissions/draft", travelCommissionHandler.GetAllDrafts)
	travel.POST("/commissions/draft", travelCommissionHandler.AddDraft)
	travel.PUT("/commissions/draft/:id", travelCommissionHandler.UpdateDraft)
	travel.DELETE("/commissions/draft/:id", travelCommissionHandler.DeleteDraft)
	travel.POST("/commissions/draft/clear", travelCommissionHandler.ClearDrafts)
	travel.POST("/commissions/draft/confirm", travelCommissionHandler.ConfirmDrafts)

	// Travel Addon routes
	travel.GET("/addons", travelAddonHandler.GetAddons)
	travel.GET("/addons/:id", travelAddonHandler.GetAddon)
	travel.POST("/addons", travelAddonHandler.CreateAddon)
	travel.PUT("/addons/:id", travelAddonHandler.UpdateAddon)
	travel.DELETE("/addons/:id", travelAddonHandler.DeleteAddon)

	// Travel Template routes
	travel.GET("/templates", travelTemplateHandler.GetAll)
	travel.GET("/templates/:id", travelTemplateHandler.GetByID)
	travel.POST("/templates", travelTemplateHandler.Create)
	travel.PUT("/templates/:id", travelTemplateHandler.Update)
	travel.DELETE("/templates/:id", travelTemplateHandler.Delete)
	travel.GET("/templates/draft", travelTemplateHandler.GetAllDrafts)
	travel.POST("/templates/draft/add", travelTemplateHandler.AddDraft)
	travel.GET("/templates/draft/:id", travelTemplateHandler.GetDraftByID)
	travel.PUT("/templates/draft/:id", travelTemplateHandler.UpdateDraft)
	travel.DELETE("/templates/draft/:id", travelTemplateHandler.DeleteDraft)
	travel.POST("/templates/draft/clear", travelTemplateHandler.ClearDrafts)
	travel.POST("/templates/draft/confirm", travelTemplateHandler.ConfirmDrafts)

	// Travel History routes
	travel.GET("/histories", travelHistoryHandler.GetAll)
	travel.GET("/histories/:id", travelHistoryHandler.GetByID)
	travel.GET("/histories/table/:table_name", travelHistoryHandler.GetByTableName)
	travel.GET("/histories/table/:table_name/record/:record_id", travelHistoryHandler.GetByRecordID)

	// Travel Product Rules routes (using legacy handler for travel_service_development.product_rules)
	travel.GET("/product-rules/template", adaptHandler(travelProductRuleHandler.DownloadTemplate))
	travel.POST("/product-rules/upload", adaptHandler(travelProductRuleHandler.UploadExcel))
	travel.GET("/product-rules/draft", adaptHandler(travelProductRuleHandler.GetAllDrafts))
	travel.POST("/product-rules/draft/clear", adaptHandler(travelProductRuleHandler.ClearDrafts))
	travel.POST("/product-rules/draft/confirm", adaptHandler(travelProductRuleHandler.ConfirmDrafts))
	travel.PUT("/product-rules/draft/:index", adaptHandler(travelProductRuleHandler.UpdateDraft))
	travel.DELETE("/product-rules/draft/:index", adaptHandler(travelProductRuleHandler.DeleteDraft))
	travel.GET("/product-rules", adaptHandler(travelProductRuleHandler.GetAll))
	travel.GET("/product-rules/:id", adaptHandler(travelProductRuleHandler.GetByID))
	travel.PUT("/product-rules/:id", adaptHandler(travelProductRuleHandler.Update))
	travel.DELETE("/product-rules/:id", adaptHandler(travelProductRuleHandler.Delete))

	// Travel Addon Rules routes (using legacy handler for travel_service_development.addon_rules)
	travel.GET("/addon-rules", adaptHandler(travelAddonRuleHandler.GetAddonRules))
	travel.POST("/addon-rules", adaptHandler(travelAddonRuleHandler.CreateAddonRule))
	travel.POST("/addon-rules/batch", adaptHandler(travelAddonRuleHandler.CreateAddonRulesBatch))
	travel.GET("/addon-rules/:id", adaptHandler(travelAddonRuleHandler.GetAddonRule))
	travel.PUT("/addon-rules/:id", adaptHandler(travelAddonRuleHandler.UpdateAddonRule))
	travel.DELETE("/addon-rules/:id", adaptHandler(travelAddonRuleHandler.DeleteAddonRule))

	// Backward compatibility: If travel is default program, map /api/* to /api/travel/*
	if config.DefaultProgram == "travel" {
		setupTravelBackwardCompatibility(api, travelProductHandler, travelInsuranceHandler, travelAddonHandler, travelHistoryHandler, travelTemplateHandler)
	}
}

// setupPropertyProgram sets up all routes for Property onboarding program
func setupPropertyProgram(e *echo.Echo, config ProgramsConfig) {
	log.Println("🏠 Initializing Property onboarding program...")

	// Initialize Property domain repositories, usecases, and handlers
	propertyProductRepository := propertyProductRepo.NewProductRepository()
	propertyProductUsecaseInstance := propertyProductUsecase.NewProductUsecase(propertyProductRepository)
	propertyProductHandler := propertyProductHttp.NewProductHandler(propertyProductUsecaseInstance)

	propertyInsuranceRepository := propertyInsuranceRepo.NewInsuranceRepository()
	propertyInsuranceUsecaseInstance := propertyInsuranceUsecase.NewInsuranceUsecase(propertyInsuranceRepository)
	propertyInsuranceHandler := propertyInsuranceHttp.NewInsuranceHandler(propertyInsuranceUsecaseInstance)

	propertyTemplateRepository := propertyTemplateRepo.NewTemplateRepository()
	propertyTemplateUsecaseInstance := propertyTemplateUsecase.NewTemplateUsecase(propertyTemplateRepository)
	propertyTemplateHandler := propertyTemplateHttp.NewTemplateHandler(propertyTemplateUsecaseInstance)

	// Initialize Property Product Rules handler
	db := utils.GetDB()
	propertyProductRuleHandler := propertyProductRuleHttp.NewPropertyProductRuleHandler(db)

	// Initialize Property Commission handler
	propertyCommissionHandler := propertyProductRuleHttp.NewPropertyCommissionHandler(db)

	// API routes group
	api := e.Group("/api")
	property := api.Group("/property")

	// Property Insurance routes
	property.GET("/insurances", propertyInsuranceHandler.GetAll)
	property.POST("/insurances", propertyInsuranceHandler.Create)
	property.GET("/insurances/:id", propertyInsuranceHandler.GetByID)
	property.PUT("/insurances/:id", propertyInsuranceHandler.Update)
	property.DELETE("/insurances/:id", propertyInsuranceHandler.Delete)
	property.GET("/insurances/draft", propertyInsuranceHandler.GetAllDrafts)
	property.POST("/insurances/draft/add", propertyInsuranceHandler.AddDraft)
	property.GET("/insurances/draft/:id", propertyInsuranceHandler.GetDraftByID)
	property.PUT("/insurances/draft/:id", propertyInsuranceHandler.UpdateDraft)
	property.DELETE("/insurances/draft/:id", propertyInsuranceHandler.DeleteDraft)
	property.POST("/insurances/draft/clear", propertyInsuranceHandler.ClearDrafts)
	property.POST("/insurances/draft/confirm", propertyInsuranceHandler.ConfirmDrafts)

	// Property Product routes
	property.GET("/products", propertyProductHandler.GetProducts)
	property.POST("/products", propertyProductHandler.CreateProduct)
	property.GET("/products/draft", propertyProductHandler.GetAllDrafts)
	property.POST("/products/draft/add", propertyProductHandler.AddDraft)
	property.POST("/products/draft/clear", propertyProductHandler.ClearDrafts)
	property.POST("/products/draft/confirm", propertyProductHandler.ConfirmDrafts)
	property.DELETE("/products/draft/:id", propertyProductHandler.DeleteDraft)
	property.GET("/products/draft/:id", propertyProductHandler.GetDraftByID)
	property.PUT("/products/draft/:id", propertyProductHandler.UpdateDraft)
	property.POST("/products/generate-templates/selected", propertyProductHandler.GenerateTemplateDraftsForProducts)
	property.GET("/products/:id", propertyProductHandler.GetProduct)
	property.PUT("/products/:id", propertyProductHandler.UpdateProduct)
	property.DELETE("/products/:id", propertyProductHandler.DeleteProduct)

	// Property Template routes
	property.GET("/templates", propertyTemplateHandler.GetAll)
	property.POST("/templates", propertyTemplateHandler.Create)
	property.GET("/templates/draft", propertyTemplateHandler.GetAllDrafts)
	property.POST("/templates/draft/add", propertyTemplateHandler.AddDraft)
	property.GET("/templates/draft/:id", propertyTemplateHandler.GetDraftByID)
	property.PUT("/templates/draft/:id", propertyTemplateHandler.UpdateDraft)
	property.DELETE("/templates/draft/:id", propertyTemplateHandler.DeleteDraft)
	property.POST("/templates/draft/clear", propertyTemplateHandler.ClearDrafts)
	property.POST("/templates/draft/confirm", propertyTemplateHandler.ConfirmDrafts)
	property.GET("/templates/:id", propertyTemplateHandler.GetByID)
	property.PUT("/templates/:id", propertyTemplateHandler.Update)
	property.DELETE("/templates/:id", propertyTemplateHandler.Delete)

	// Property Product Rules routes
	property.GET("/product-rules/template", propertyProductRuleHandler.DownloadTemplate)
	property.POST("/product-rules/upload", propertyProductRuleHandler.Upload)
	property.GET("/product-rules/draft", propertyProductRuleHandler.GetAllDrafts)
	property.POST("/product-rules/draft/clear", propertyProductRuleHandler.ClearDrafts)
	property.POST("/product-rules/draft/confirm", propertyProductRuleHandler.ConfirmDrafts)
	property.PUT("/product-rules/draft/:index", propertyProductRuleHandler.UpdateDraft)
	property.DELETE("/product-rules/draft/:index", propertyProductRuleHandler.DeleteDraft)
	property.GET("/product-rules", propertyProductRuleHandler.GetAll)

	// Property Commission routes
	property.POST("/commissions/generate", propertyCommissionHandler.GenerateCommissions)
	property.POST("/commissions/check-duplicates", propertyCommissionHandler.CheckDuplicateProductCodes)
	property.GET("/commissions/table/commissions", propertyCommissionHandler.GetCommissionsFromTable)
	property.GET("/commissions/table/plan_commissions", propertyCommissionHandler.GetPlanCommissionsFromTable)
	property.GET("/commissions/table/default_config_products", propertyCommissionHandler.GetDefaultConfigFromTable)

	// Initialize Property Addon Rule Detail handler
	propertyAddonRuleDetailHandler := propertyProductRuleHttp.NewPropertyAddonRuleDetailHandler(db)

	// Property Addon Rule Detail routes
	property.GET("/addon-rule-details/template", propertyAddonRuleDetailHandler.DownloadTemplate)
	property.POST("/addon-rule-details/upload", propertyAddonRuleDetailHandler.Upload)
	property.GET("/addon-rule-details/draft", propertyAddonRuleDetailHandler.GetAllDrafts)
	property.POST("/addon-rule-details/draft/clear", propertyAddonRuleDetailHandler.ClearDrafts)
	property.POST("/addon-rule-details/draft/confirm", propertyAddonRuleDetailHandler.ConfirmDrafts)
	property.GET("/addon-rule-details", propertyAddonRuleDetailHandler.GetAll)
	property.PUT("/addon-rule-details/:id", propertyAddonRuleDetailHandler.Update)
	property.DELETE("/addon-rule-details/:id", propertyAddonRuleDetailHandler.Delete)

	// Backward compatibility: If property is default program, map /api/* to /api/property/*
	if config.DefaultProgram == "property" {
		setupPropertyBackwardCompatibility(api, propertyProductHandler)
	}
}

// setupPropertyBackwardCompatibility sets up backward compatibility routes for property
func setupPropertyBackwardCompatibility(api *echo.Group, propertyProductHandler interface{}) {
	// Product routes (backward compatibility)
	api.GET("/products", func(c echo.Context) error {
		return propertyProductHandler.(interface{ GetProducts(echo.Context) error }).GetProducts(c)
	})
	api.POST("/products", func(c echo.Context) error {
		return propertyProductHandler.(interface{ CreateProduct(echo.Context) error }).CreateProduct(c)
	})
	api.GET("/products/draft", func(c echo.Context) error {
		return propertyProductHandler.(interface{ GetAllDrafts(echo.Context) error }).GetAllDrafts(c)
	})
	api.POST("/products/draft/add", func(c echo.Context) error {
		return propertyProductHandler.(interface{ AddDraft(echo.Context) error }).AddDraft(c)
	})
	api.POST("/products/draft/clear", func(c echo.Context) error {
		return propertyProductHandler.(interface{ ClearDrafts(echo.Context) error }).ClearDrafts(c)
	})
	api.POST("/products/draft/confirm", func(c echo.Context) error {
		return propertyProductHandler.(interface{ ConfirmDrafts(echo.Context) error }).ConfirmDrafts(c)
	})
	api.DELETE("/products/draft/:id", func(c echo.Context) error {
		return propertyProductHandler.(interface{ DeleteDraft(echo.Context) error }).DeleteDraft(c)
	})
	api.GET("/products/:id", func(c echo.Context) error {
		return propertyProductHandler.(interface{ GetProduct(echo.Context) error }).GetProduct(c)
	})
	api.PUT("/products/:id", func(c echo.Context) error {
		return propertyProductHandler.(interface{ UpdateProduct(echo.Context) error }).UpdateProduct(c)
	})
	api.DELETE("/products/:id", func(c echo.Context) error {
		return propertyProductHandler.(interface{ DeleteProduct(echo.Context) error }).DeleteProduct(c)
	})
}

// setupTravelBackwardCompatibility sets up backward compatibility routes for travel
func setupTravelBackwardCompatibility(api *echo.Group, travelProductHandler interface{}, travelInsuranceHandler interface{}, travelAddonHandler interface{}, travelHistoryHandler interface{}, travelTemplateHandler interface{}) {
	// Product routes (backward compatibility)
	api.GET("/products", func(c echo.Context) error {
		return travelProductHandler.(interface{ GetProducts(echo.Context) error }).GetProducts(c)
	})
	api.POST("/products", func(c echo.Context) error {
		return travelProductHandler.(interface{ CreateProduct(echo.Context) error }).CreateProduct(c)
	})
	api.GET("/products/draft", func(c echo.Context) error {
		return travelProductHandler.(interface{ GetAllDrafts(echo.Context) error }).GetAllDrafts(c)
	})
	api.POST("/products/draft/add", func(c echo.Context) error {
		return travelProductHandler.(interface{ AddDraft(echo.Context) error }).AddDraft(c)
	})
	api.POST("/products/draft/clear", func(c echo.Context) error {
		return travelProductHandler.(interface{ ClearDrafts(echo.Context) error }).ClearDrafts(c)
	})
	api.POST("/products/draft/confirm", func(c echo.Context) error {
		return travelProductHandler.(interface{ ConfirmDrafts(echo.Context) error }).ConfirmDrafts(c)
	})
	api.DELETE("/products/draft/:id", func(c echo.Context) error {
		return travelProductHandler.(interface{ DeleteDraft(echo.Context) error }).DeleteDraft(c)
	})
	api.GET("/products/:id", func(c echo.Context) error {
		return travelProductHandler.(interface{ GetProduct(echo.Context) error }).GetProduct(c)
	})
	api.PUT("/products/:id", func(c echo.Context) error {
		return travelProductHandler.(interface{ UpdateProduct(echo.Context) error }).UpdateProduct(c)
	})
	api.DELETE("/products/:id", func(c echo.Context) error {
		return travelProductHandler.(interface{ DeleteProduct(echo.Context) error }).DeleteProduct(c)
	})

	// Insurance routes (backward compatibility)
	api.GET("/insurances", func(c echo.Context) error {
		return travelInsuranceHandler.(interface{ GetAll(echo.Context) error }).GetAll(c)
	})
	api.POST("/insurances", func(c echo.Context) error {
		return travelInsuranceHandler.(interface{ Create(echo.Context) error }).Create(c)
	})
	api.GET("/insurances/:id", func(c echo.Context) error {
		return travelInsuranceHandler.(interface{ GetByID(echo.Context) error }).GetByID(c)
	})
	api.PUT("/insurances/:id", func(c echo.Context) error {
		return travelInsuranceHandler.(interface{ Update(echo.Context) error }).Update(c)
	})
	api.DELETE("/insurances/:id", func(c echo.Context) error {
		return travelInsuranceHandler.(interface{ Delete(echo.Context) error }).Delete(c)
	})

	// Addon routes (backward compatibility)
	api.GET("/addons", func(c echo.Context) error {
		return travelAddonHandler.(interface{ GetAddons(echo.Context) error }).GetAddons(c)
	})
	api.GET("/addons/:id", func(c echo.Context) error {
		return travelAddonHandler.(interface{ GetAddon(echo.Context) error }).GetAddon(c)
	})
	api.POST("/addons", func(c echo.Context) error {
		return travelAddonHandler.(interface{ CreateAddon(echo.Context) error }).CreateAddon(c)
	})
	api.PUT("/addons/:id", func(c echo.Context) error {
		return travelAddonHandler.(interface{ UpdateAddon(echo.Context) error }).UpdateAddon(c)
	})
	api.DELETE("/addons/:id", func(c echo.Context) error {
		return travelAddonHandler.(interface{ DeleteAddon(echo.Context) error }).DeleteAddon(c)
	})

	// History routes (backward compatibility)
	api.GET("/histories", func(c echo.Context) error {
		return travelHistoryHandler.(interface{ GetAll(echo.Context) error }).GetAll(c)
	})
	api.GET("/histories/:id", func(c echo.Context) error {
		return travelHistoryHandler.(interface{ GetByID(echo.Context) error }).GetByID(c)
	})

	// Template routes (backward compatibility)
	api.GET("/templates", func(c echo.Context) error {
		return travelTemplateHandler.(interface{ GetAll(echo.Context) error }).GetAll(c)
	})
	api.GET("/templates/:id", func(c echo.Context) error {
		return travelTemplateHandler.(interface{ GetByID(echo.Context) error }).GetByID(c)
	})
	api.POST("/templates", func(c echo.Context) error {
		return travelTemplateHandler.(interface{ Create(echo.Context) error }).Create(c)
	})
	api.PUT("/templates/:id", func(c echo.Context) error {
		return travelTemplateHandler.(interface{ Update(echo.Context) error }).Update(c)
	})
	api.DELETE("/templates/:id", func(c echo.Context) error {
		return travelTemplateHandler.(interface{ Delete(echo.Context) error }).Delete(c)
	})
}

// setupLegacyRoutes sets up legacy routes that are shared across programs
func setupLegacyRoutes(e *echo.Echo) {
	db := utils.GetDB()

	// Legacy handlers for domains not yet migrated
	legacyHistoryRepo := repository.NewMySQLHistoryRepository(db)

	// NOTE: productRuleHandler removed - product rules now use domain-specific routes:
	// - Vehicle: /api/product-rules (via vehicleRulesHttp.AddRulesHandler)
	// - Travel: /api/travel/product-rules (in setupTravelProgram)

	// NOTE: addonRuleHandler removed - addon rules now use domain-specific routes:
	// - Vehicle: /api/addon-rules (registered via vehicleAddonRulesHttp.AddAddonRulesHandler)
	// - Travel: /api/travel/addon-rules (registered in setupTravelProgram)
	// These shared legacy routes were conflicting with vehicle domain routes

	insuranceProductAddonMappingRepo := repository.NewMySQLInsuranceProductAddonMappingRepository(db)
	insuranceProductAddonMappingUsecase := usecase.NewInsuranceProductAddonMappingUsecase(insuranceProductAddonMappingRepo)
	insuranceProductAddonMappingHandler := httphandler.NewInsuranceProductAddonMappingHandler(insuranceProductAddonMappingUsecase)

	addonRuleDetailRepo := repository.NewMySQLAddonRuleDetailRepository(db, legacyHistoryRepo)
	addonRuleDetailUsecase := usecase.NewAddonRuleDetailUsecase(addonRuleDetailRepo)
	addonRuleDetailHandler := httphandler.NewAddonRuleDetailHandler(addonRuleDetailUsecase, db)

	addonProductUsecase := usecase.NewAddonProductUsecase(
		repository.NewMySQLAddonRepository(db),
		repository.NewMySQLAddonRuleRepository(db), // Create new instance for addonProductUsecase
		insuranceProductAddonMappingRepo,
	)
	addonProductHandler := httphandler.NewAddonProductHandler(addonProductUsecase)

	allResultHandler := httphandler.NewAllResultHandler(db)
	htmlHandler := httphandler.NewHTMLHandler()

	legacyProductRepo := repository.NewMySQLProductRepository(db)
	legacyCommissionUsecase := usecase.NewCommissionUsecase(repository.NewMySQLCommissionRepository(db), legacyProductRepo)
	legacyCommissionHandler := httphandler.NewCommissionHandler(legacyCommissionUsecase)

	// Helper function to adapt http.HandlerFunc to echo.HandlerFunc
	adaptHandler := func(h http.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), "echo_context", c)
			r := c.Request().WithContext(ctx)
			h(c.Response(), r)
			return nil
		}
	}

	api := e.Group("/api")

	// NOTE: Legacy Product Rule routes removed - now using domain-specific routes:
	// - Vehicle: /api/product-rules (registered via vehicleRulesHttp.AddRulesHandler)
	// - Travel: /api/travel/product-rules (registered in setupTravelProgram)
	// These shared legacy routes were conflicting with vehicle domain routes

	// NOTE: Legacy Addon Rule routes removed - now using domain-specific routes:
	// - Vehicle: /api/addon-rules (registered via vehicleAddonRulesHttp.AddAddonRulesHandler)
	// - Travel: /api/travel/addon-rules (registered in setupTravelProgram)
	// These shared legacy routes were conflicting with vehicle domain routes

	// Legacy Insurance Product Addon Mapping routes (shared)
	api.POST("/insurance-product-addon-mappings", adaptHandler(insuranceProductAddonMappingHandler.CreateInsuranceProductAddonMapping))
	api.POST("/insurance-product-addon-mappings/batch", adaptHandler(insuranceProductAddonMappingHandler.CreateInsuranceProductAddonMappingsBatch))

	// Legacy Product Addons routes (shared)
	api.POST("/products/addons", adaptHandler(addonProductHandler.SaveProductAddons))
	api.POST("/products/addons/cleanup-duplicates", adaptHandler(addonProductHandler.CleanupDuplicates))

	// Legacy Addon Rule Detail routes (shared)
	api.GET("/addon-rule-details/template", adaptHandler(addonRuleDetailHandler.DownloadTemplate))
	api.POST("/addon-rule-details/upload", adaptHandler(addonRuleDetailHandler.UploadExcel))
	api.GET("/addon-rule-details/draft", adaptHandler(addonRuleDetailHandler.GetAllDrafts))
	api.POST("/addon-rule-details/draft/clear", adaptHandler(addonRuleDetailHandler.ClearDrafts))
	api.POST("/addon-rule-details/draft/confirm", adaptHandler(addonRuleDetailHandler.ConfirmDrafts))
	api.GET("/addon-rule-details", adaptHandler(addonRuleDetailHandler.GetAll))
	api.GET("/addon-rule-details/:id", adaptHandler(addonRuleDetailHandler.GetByID))
	api.PUT("/addon-rule-details/:id", adaptHandler(addonRuleDetailHandler.Update))

	// Legacy Commission advanced routes (shared)
	api.POST("/commissions/generate", adaptHandler(legacyCommissionHandler.GenerateCommissions))
	api.POST("/commissions/check-duplicates", adaptHandler(legacyCommissionHandler.CheckDuplicateProductCodes))
	api.GET("/commissions/table/commissions", adaptHandler(legacyCommissionHandler.GetCommissionsFromTable))
	api.GET("/commissions/table/plan_commissions", adaptHandler(legacyCommissionHandler.GetPlanCommissionsFromTable))
	api.GET("/commissions/table/default_config_products", adaptHandler(legacyCommissionHandler.GetDefaultConfigFromTable))

	// Legacy All Result routes (shared)
	api.GET("/all-result/download", adaptHandler(allResultHandler.DownloadAllData))

	// Legacy HTML converter route (shared)
	api.POST("/html/convert", adaptHandler(htmlHandler.ConvertHTML))
}
