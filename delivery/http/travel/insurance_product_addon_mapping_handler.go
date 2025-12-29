package http

import (
	"encoding/json"
	"net/http"
	"multi-onboarding/model/Travel"
	"multi-onboarding/usecase/Travel"
)

// InsuranceProductAddonMappingHandler handles HTTP requests for insurance product addon mappings
type InsuranceProductAddonMappingHandler struct {
	usecase *usecase.InsuranceProductAddonMappingUsecase
}

// NewInsuranceProductAddonMappingHandler creates a new insurance product addon mapping handler
func NewInsuranceProductAddonMappingHandler(usecase *usecase.InsuranceProductAddonMappingUsecase) *InsuranceProductAddonMappingHandler {
	return &InsuranceProductAddonMappingHandler{
		usecase: usecase,
	}
}

// CreateInsuranceProductAddonMapping handles POST /api/insurance-product-addon-mappings
func (h *InsuranceProductAddonMappingHandler) CreateInsuranceProductAddonMapping(w http.ResponseWriter, r *http.Request) {
	var mapping model.InsuranceProductAddonMapping
	if err := json.NewDecoder(r.Body).Decode(&mapping); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.usecase.Create(&mapping); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapping)
}

// CreateInsuranceProductAddonMappingsBatch handles POST /api/insurance-product-addon-mappings/batch
func (h *InsuranceProductAddonMappingHandler) CreateInsuranceProductAddonMappingsBatch(w http.ResponseWriter, r *http.Request) {
	var mappings []model.InsuranceProductAddonMapping
	if err := json.NewDecoder(r.Body).Decode(&mappings); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.usecase.CreateBatch(mappings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Insurance product addon mappings created successfully",
		"count":   len(mappings),
	})
}

