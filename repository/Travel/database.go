package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	User     string
	Password string
	Database string
	Port     int
}

// NewDBConnection creates a new database connection
func NewDBConnection(config DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// InitSchema initializes the database schema
func InitSchema(db *sql.DB) error {
	// Create travel_products table
	productsQuery := `
	CREATE TABLE IF NOT EXISTS travel_products (
		id INT AUTO_INCREMENT PRIMARY KEY,
		code VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		destination VARCHAR(255),
		price DECIMAL(10, 2) DEFAULT 0,
		duration INT DEFAULT 1,
		category VARCHAR(100),
		status VARCHAR(50) DEFAULT 'draft',
		insurance_code VARCHAR(100),
		type VARCHAR(50),
		region_id INT,
		is_active TINYINT(1) DEFAULT 1,
		min_adult INT DEFAULT 0,
		max_adult INT DEFAULT 0,
		min_child INT DEFAULT 0,
		max_child INT DEFAULT 0,
		logo TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_status (status),
		INDEX idx_insurance_code (insurance_code),
		INDEX idx_code (code),
		INDEX idx_region_id (region_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	_, err := db.Exec(productsQuery)
	if err != nil {
		return fmt.Errorf("failed to create travel_products table: %w", err)
	}

	// Create regions table
	regionsQuery := `
	CREATE TABLE IF NOT EXISTS regions (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		type ENUM('WHITELIST', 'BLACKLIST') NOT NULL,
		country_ids JSON NOT NULL,
		created_by BIGINT NOT NULL DEFAULT 1,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_by BIGINT DEFAULT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_type (type)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	_, err = db.Exec(regionsQuery)
	if err != nil {
		return fmt.Errorf("failed to create regions table: %w", err)
	}

	// Create insurances table
	insurancesQuery := `
	CREATE TABLE IF NOT EXISTS insurances (
		id INT AUTO_INCREMENT PRIMARY KEY,
		code VARCHAR(100) NOT NULL UNIQUE,
		name VARCHAR(255) NOT NULL,
		status TINYINT(1) NOT NULL DEFAULT 1,
		logo TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_code (code),
		INDEX idx_status (status)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	_, err = db.Exec(insurancesQuery)
	if err != nil {
		return fmt.Errorf("failed to create insurances table: %w", err)
	}

	// Note: insurance_drafts table no longer needed - using in-memory storage

	// Create product_rules table
	productRulesQuery := `
	CREATE TABLE IF NOT EXISTS product_rules (
		id BIGINT NOT NULL AUTO_INCREMENT,
		product_code VARCHAR(255) NOT NULL,
		premium_type VARCHAR(255) NOT NULL,
		start_days INT NOT NULL,
		end_days INT NOT NULL,
		base_premium_value DOUBLE NOT NULL,
		rules TEXT NOT NULL,
		created_by BIGINT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_by BIGINT DEFAULT NULL,
		updated_at TIMESTAMP NULL DEFAULT NULL,
		PRIMARY KEY (id),
		KEY idx_product_rules_product_code (product_code)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	_, err = db.Exec(productRulesQuery)
	if err != nil {
		return fmt.Errorf("failed to create product_rules table: %w", err)
	}

	// Note: product_rule_drafts table no longer needed - using in-memory storage

	// Create commissions table
	commissionsQuery := `
	CREATE TABLE IF NOT EXISTS commissions (
		id BIGINT NOT NULL AUTO_INCREMENT,
		product_code VARCHAR(255) NOT NULL,
		commission_percentage DOUBLE NOT NULL,
		commission_vat_type VARCHAR(50) NOT NULL,
		af_percentage DOUBLE NOT NULL,
		af_vat_type VARCHAR(50) NOT NULL,
		admin_fee DOUBLE NOT NULL,
		created_by BIGINT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_by BIGINT DEFAULT NULL,
		updated_at TIMESTAMP NULL DEFAULT NULL,
		PRIMARY KEY (id),
		KEY idx_product_code (product_code)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	_, err = db.Exec(commissionsQuery)
	if err != nil {
		return fmt.Errorf("failed to create commissions table: %w", err)
	}

	// Note: commission_drafts table no longer needed - using in-memory storage

	// Create addons table
	addonsQuery := `
	CREATE TABLE IF NOT EXISTS addons (
		id BIGINT NOT NULL AUTO_INCREMENT,
		code VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		name_my VARCHAR(255) NOT NULL,
		name_en VARCHAR(255) NOT NULL,
		is_active TINYINT NOT NULL DEFAULT 0,
		created_by BIGINT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_by BIGINT DEFAULT NULL,
		updated_at TIMESTAMP NULL DEFAULT NULL,
		PRIMARY KEY (id),
		UNIQUE KEY idx_addons_code (code)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
	`

	_, err = db.Exec(addonsQuery)
	if err != nil {
		return fmt.Errorf("failed to create addons table: %w", err)
	}

	// Create addon_rules table
	addonRulesQuery := `
	CREATE TABLE IF NOT EXISTS travel_service_development.addon_rules (
		id BIGINT NOT NULL AUTO_INCREMENT,
		product_code VARCHAR(255) NOT NULL,
		addon_code VARCHAR(255) NOT NULL,
		insurance_code VARCHAR(255) DEFAULT NULL,
		rules JSON DEFAULT NULL,
		created_by BIGINT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_by BIGINT DEFAULT NULL,
		updated_at TIMESTAMP NULL DEFAULT NULL,
		PRIMARY KEY (id),
		KEY idx_product_code (product_code),
		KEY idx_addon_code (addon_code),
		KEY idx_insurance_code (insurance_code),
		UNIQUE KEY idx_product_addon (product_code, addon_code)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
	`

	_, err = db.Exec(addonRulesQuery)
	if err != nil {
		return fmt.Errorf("failed to create addon_rules table: %w", err)
	}

	// Create addon_rule_details table
	addonRuleDetailsQuery := `
	CREATE TABLE IF NOT EXISTS travel_service_development.addon_rule_details (
		id BIGINT NOT NULL AUTO_INCREMENT,
		addon_rule_id BIGINT NOT NULL,
		start_condition VARCHAR(255) NOT NULL DEFAULT '-1',
		end_condition VARCHAR(255) NOT NULL DEFAULT '-1',
		value_type VARCHAR(255) NOT NULL DEFAULT 'FIXED',
		value DOUBLE NOT NULL DEFAULT 0,
		duration_rule_type VARCHAR(255) NOT NULL DEFAULT 'DAILY',
		min_adult INT NOT NULL DEFAULT 0,
		max_adult INT NOT NULL DEFAULT 0,
		max_age INT NOT NULL DEFAULT 0,
		created_by BIGINT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_by BIGINT DEFAULT NULL,
		updated_at TIMESTAMP NULL DEFAULT NULL,
		PRIMARY KEY (id),
		KEY idx_addon_rule_id (addon_rule_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
	`

	_, err = db.Exec(addonRuleDetailsQuery)
	if err != nil {
		return fmt.Errorf("failed to create addon_rule_details table: %w", err)
	}

	// Create templates table
	templatesQuery := `
	CREATE TABLE IF NOT EXISTS travel_service_development.templates (
		locale VARCHAR(255) NOT NULL,
		id VARCHAR(255) NOT NULL,
		value TEXT NOT NULL,
		created_by BIGINT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_by BIGINT DEFAULT NULL,
		updated_at TIMESTAMP NULL DEFAULT NULL,
		PRIMARY KEY (locale, id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	_, err = db.Exec(templatesQuery)
	if err != nil {
		return fmt.Errorf("failed to create templates table: %w", err)
	}

	// Create template_drafts table
	templateDraftsQuery := `
	CREATE TABLE IF NOT EXISTS travel_service_development.template_drafts (
		id VARCHAR(255) NOT NULL,
		locale VARCHAR(10) NOT NULL,
		value TEXT NOT NULL,
		product_code VARCHAR(255) DEFAULT NULL,
		insurance JSON DEFAULT NULL,
		created_by BIGINT NOT NULL DEFAULT 1,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id, locale),
		INDEX idx_product_code (product_code)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	_, err = db.Exec(templateDraftsQuery)
	if err != nil {
		return fmt.Errorf("failed to create template_drafts table: %w", err)
	}

	// Create countries table if not exists
	countriesQuery := `
	CREATE TABLE IF NOT EXISTS countries (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		code VARCHAR(10) NOT NULL,
		UNIQUE KEY uk_code (code),
		INDEX idx_name (name)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	_, err = db.Exec(countriesQuery)
	if err != nil {
		return fmt.Errorf("failed to create countries table: %w", err)
	}

	// Add code column if it doesn't exist (for existing tables)
	alterQuery := `
	ALTER TABLE countries 
	ADD COLUMN IF NOT EXISTS code VARCHAR(10) NOT NULL DEFAULT '' AFTER name;
	`
	db.Exec(alterQuery) // Ignore error if column already exists

	// Note: histories table already exists with structure:
	// id, user_name, section, action, record_id, record_type, data_before, data_after, created_at
	// No need to create it here

	return nil
}
