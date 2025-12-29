package http

import (
	"encoding/json"
	"net/http"

	"multi-onboarding/model/Travel"
	"multi-onboarding/usecase/Travel"
)

// AddonProductHandler handles HTTP requests for saving addons for products
type AddonProductHandler struct {
	usecase *usecase.AddonProductUsecase
}

// NewAddonProductHandler creates a new addon product handler
func NewAddonProductHandler(uc *usecase.AddonProductUsecase) *AddonProductHandler {
	return &AddonProductHandler{
		usecase: uc,
	}
}

// SaveProductAddons handles POST /api/products/addons
func (h *AddonProductHandler) SaveProductAddons(w http.ResponseWriter, r *http.Request) {
	var req model.AddonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.ProductCode == "" || req.InsuranceCode == "" {
		respondWithError(w, http.StatusBadRequest, "product_code and insurance_code are required")
		return
	}

	if err := h.usecase.SaveProductAddons(&req); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Addons saved successfully",
	})
}

// CleanupDuplicates handles POST /api/products/addons/cleanup-duplicates
// This endpoint removes duplicate addon_rules entries (keeping only the first one)
func (h *AddonProductHandler) CleanupDuplicates(w http.ResponseWriter, r *http.Request) {
	if err := h.usecase.CleanupDuplicateAddonRules(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Duplicate addon rules cleaned up successfully"})
}

