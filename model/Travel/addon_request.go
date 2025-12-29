package model

// AddonRequest represents a request to save addons for a product
type AddonRequest struct {
	ProductCode   string           `json:"product_code"`
	InsuranceCode string           `json:"insurance_code"`
	Addons        []AddonWithNames `json:"addons"`
}

// AddonWithNames represents an addon with its names
type AddonWithNames struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	NameMy string `json:"name_my"`
	NameEn string `json:"name_en"`
}


