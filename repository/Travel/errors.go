package repository

import "errors"

var (
	ErrProductNotFound    = errors.New("product not found")
	ErrRegionNotFound     = errors.New("region not found")
	ErrCountryNotFound    = errors.New("country not found")
	ErrCommissionNotFound = errors.New("commission not found")
	ErrAddonNotFound      = errors.New("addon not found")
	ErrAddonRuleNotFound  = errors.New("addon rule not found")
)
