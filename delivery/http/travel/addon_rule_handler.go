package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"multi-onboarding/model/Travel"
	"multi-onboarding/usecase/Travel"
)

// AddonRuleHandler handles HTTP requests for addon rules
type AddonRuleHandler struct {
	usecase usecase.AddonRuleUsecase
}

// NewAddonRuleHandler creates a new addon rule handler
func NewAddonRuleHandler(uc usecase.AddonRuleUsecase) *AddonRuleHandler {
	return &AddonRuleHandler{
		usecase: uc,
	}
}

// GetAddonRules handles GET /api/addon-rules
func (h *AddonRuleHandler) GetAddonRules(w http.ResponseWriter, r *http.Request) {
	productCode := r.URL.Query().Get("product_code")
	addonCode := r.URL.Query().Get("addon_code")
	
	var addonRules []model.AddonRule
	var err error
	
	if productCode != "" {
		addonRules, err = h.usecase.GetAddonRulesByProductCode(productCode)
	} else if addonCode != "" {
		addonRules, err = h.usecase.GetAddonRulesByAddonCode(addonCode)
	} else {
		addonRules, err = h.usecase.GetAllAddonRules()
	}
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Ensure we always return an array, not null
	if addonRules == nil {
		addonRules = []model.AddonRule{}
	}

	respondWithJSON(w, http.StatusOK, addonRules)
}

// GetAddonRule handles GET /api/addon-rules/{id}
func (h *AddonRuleHandler) GetAddonRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid addon rule ID")
		return
	}

	addonRule, err := h.usecase.GetAddonRuleByID(id)
	if err != nil {
		if err == usecase.ErrAddonRuleNotFound || err == usecase.ErrInvalidAddonRuleID {
			respondWithError(w, http.StatusNotFound, "Addon rule not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, addonRule)
}

// CreateAddonRule handles POST /api/addon-rules
func (h *AddonRuleHandler) CreateAddonRule(w http.ResponseWriter, r *http.Request) {
	var addonRule model.AddonRule

	if err := json.NewDecoder(r.Body).Decode(&addonRule); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.usecase.CreateAddonRule(&addonRule); err != nil {
		statusCode := http.StatusBadRequest
		if err == usecase.ErrAddonRuleNotFound {
			statusCode = http.StatusNotFound
		}
		respondWithError(w, statusCode, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, addonRule)
}

// CreateAddonRulesBatch handles POST /api/addon-rules/batch
func (h *AddonRuleHandler) CreateAddonRulesBatch(w http.ResponseWriter, r *http.Request) {
	var addonRules []model.AddonRule

	if err := json.NewDecoder(r.Body).Decode(&addonRules); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.usecase.CreateAddonRulesBatch(addonRules); err != nil {
		statusCode := http.StatusBadRequest
		if err == usecase.ErrEmptyAddonRulesBatch {
			statusCode = http.StatusBadRequest
		}
		respondWithError(w, statusCode, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, addonRules)
}

// DeleteAddonRule handles DELETE /api/addon-rules/{id}
func (h *AddonRuleHandler) DeleteAddonRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid addon rule ID")
		return
	}

	err = h.usecase.DeleteAddonRule(id)
	if err != nil {
		if err == usecase.ErrAddonRuleNotFound || err == usecase.ErrInvalidAddonRuleID {
			respondWithError(w, http.StatusNotFound, "Addon rule not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateAddonRule handles PUT /api/addon-rules/{id}
func (h *AddonRuleHandler) UpdateAddonRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid addon rule ID")
		return
	}

	var addonRule model.AddonRule
	if err := json.NewDecoder(r.Body).Decode(&addonRule); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = h.usecase.UpdateAddonRule(id, &addonRule)
	if err != nil {
		if err == usecase.ErrAddonRuleNotFound || err == usecase.ErrInvalidAddonRuleID {
			respondWithError(w, http.StatusNotFound, "Addon rule not found")
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Fetch updated addon rule to return
	updatedAddonRule, err := h.usecase.GetAddonRuleByID(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch updated addon rule")
		return
	}

	respondWithJSON(w, http.StatusOK, updatedAddonRule)
}

