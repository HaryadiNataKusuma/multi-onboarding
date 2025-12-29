package usecase

import "errors"

var (
	ErrInvalidID        = errors.New("invalid product ID")
	ErrNameRequired     = errors.New("product name is required")
	ErrInvalidPrice     = errors.New("price must be greater than or equal to 0")
	ErrInvalidDuration  = errors.New("duration must be greater than 0")
	ErrCategoryRequired = errors.New("category is required")
	ErrInvalidStatus    = errors.New("invalid status")
	ErrProductNotFound  = errors.New("product not found")

	// Region errors
	ErrInvalidRegionID    = errors.New("invalid region ID")
	ErrRegionNameRequired = errors.New("region name is required")
	ErrInvalidRegionType  = errors.New("region type must be WHITELIST or BLACKLIST")
	ErrCountryIDRequired  = errors.New("country ID is required")
	ErrRegionNotFound     = errors.New("region not found")

	// Addon errors
	ErrAddonNotFound       = errors.New("addon not found")
	ErrInvalidAddonID      = errors.New("invalid addon ID")
	ErrAddonCodeRequired   = errors.New("addon code is required")
	ErrAddonNameRequired   = errors.New("addon name is required")
	ErrAddonNameMyRequired = errors.New("addon name_my is required")
	ErrAddonNameEnRequired = errors.New("addon name_en is required")

	// AddonRule errors
	ErrAddonRuleNotFound            = errors.New("addon rule not found")
	ErrInvalidAddonRuleID           = errors.New("invalid addon rule ID")
	ErrAddonRuleProductCodeRequired = errors.New("addon rule product code is required")
	ErrAddonRuleAddonCodeRequired   = errors.New("addon rule addon code is required")
	ErrEmptyAddonRulesBatch         = errors.New("empty addon rules batch")
)
