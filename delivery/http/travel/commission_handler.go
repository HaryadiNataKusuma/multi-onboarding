package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"multi-onboarding/model/Travel"
	"multi-onboarding/usecase/Travel"

	"github.com/gorilla/mux"
)

// CommissionHandler handles HTTP requests for commissions
type CommissionHandler struct {
	usecase *usecase.CommissionUsecase
}

// NewCommissionHandler creates a new commission handler
func NewCommissionHandler(usecase *usecase.CommissionUsecase) *CommissionHandler {
	return &CommissionHandler{
		usecase: usecase,
	}
}

// GetAll handles GET /api/commissions
func (h *CommissionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	commissions, err := h.usecase.GetAllCommissions()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, commissions)
}

// GetByID handles GET /api/commissions/{id}
func (h *CommissionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid commission ID")
		return
	}

	commission, err := h.usecase.GetCommissionByID(id)
	if err != nil {
		if err == usecase.ErrCommissionNotFound {
			respondWithError(w, http.StatusNotFound, "Commission not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, commission)
}

// GetByProductCode handles GET /api/commissions/product/{productCode}
func (h *CommissionHandler) GetByProductCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productCode := vars["productCode"]

	commissions, err := h.usecase.GetCommissionsByProductCode(productCode)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, commissions)
}

// Create handles POST /api/commissions
func (h *CommissionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var commission model.Commission
	if err := json.NewDecoder(r.Body).Decode(&commission); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if commission.CreatedBy == 0 {
		commission.CreatedBy = 1 // Default to 1 if not provided
	}

	if err := h.usecase.CreateCommission(&commission); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, commission)
}

// Update handles PUT /api/commissions/{id}
func (h *CommissionHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid commission ID")
		return
	}

	var commission model.Commission
	if err := json.NewDecoder(r.Body).Decode(&commission); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if commission.UpdatedBy == nil {
		updatedBy := int64(1) // Default to 1 if not provided
		commission.UpdatedBy = &updatedBy
	}

	if err := h.usecase.UpdateCommission(id, &commission); err != nil {
		if err == usecase.ErrCommissionNotFound {
			respondWithError(w, http.StatusNotFound, "Commission not found")
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, commission)
}

// Delete handles DELETE /api/commissions/{id}
func (h *CommissionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid commission ID")
		return
	}

	if err := h.usecase.DeleteCommission(id); err != nil {
		if err == usecase.ErrCommissionNotFound {
			respondWithError(w, http.StatusNotFound, "Commission not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Commission deleted successfully"})
}

// GetAllDrafts handles GET /api/commissions/draft
func (h *CommissionHandler) GetAllDrafts(w http.ResponseWriter, r *http.Request) {
	drafts, err := h.usecase.GetAllDrafts()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, drafts)
}

// CreateDraft handles POST /api/commissions/draft
func (h *CommissionHandler) CreateDraft(w http.ResponseWriter, r *http.Request) {
	var draft model.CommissionDraft
	if err := json.NewDecoder(r.Body).Decode(&draft); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.usecase.CreateDraft(&draft); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, draft)
}

// UpdateDraft handles PUT /api/commissions/draft/{id}
func (h *CommissionHandler) UpdateDraft(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid draft ID")
		return
	}

	var draft model.CommissionDraft
	if err := json.NewDecoder(r.Body).Decode(&draft); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.usecase.UpdateDraft(id, &draft); err != nil {
		if err == usecase.ErrCommissionNotFound {
			respondWithError(w, http.StatusNotFound, "Draft not found")
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, draft)
}

// DeleteDraft handles DELETE /api/commissions/draft/{id}
func (h *CommissionHandler) DeleteDraft(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid draft ID")
		return
	}

	if err := h.usecase.DeleteDraft(id); err != nil {
		if err == usecase.ErrCommissionNotFound {
			respondWithError(w, http.StatusNotFound, "Draft not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Draft deleted successfully"})
}

// ConfirmDrafts handles POST /api/commissions/draft/confirm
func (h *CommissionHandler) ConfirmDrafts(w http.ResponseWriter, r *http.Request) {
	if err := h.usecase.ConfirmDrafts(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Drafts confirmed successfully"})
}

// GenerateCommissions handles POST /api/commissions/generate
// This endpoint generates commission data to 3 tables:
// 1. commission_service_development.commissions
// 2. agent_service_development.default_config_products
// 3. quotation_service_development.plan_commissions
func (h *CommissionHandler) GenerateCommissions(w http.ResponseWriter, r *http.Request) {
	var request usecase.GenerateCommissionRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.usecase.GenerateCommissionsToMultipleTables(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Commissions generated successfully"})
}

// GetCommissionsFromTable handles GET /api/commissions/table/commissions
func (h *CommissionHandler) GetCommissionsFromTable(w http.ResponseWriter, r *http.Request) {
	insuranceCode := r.URL.Query().Get("insurance_code")
	data, err := h.usecase.GetCommissionsFromTable(insuranceCode)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, data)
}

// GetPlanCommissionsFromTable handles GET /api/commissions/table/plan_commissions
func (h *CommissionHandler) GetPlanCommissionsFromTable(w http.ResponseWriter, r *http.Request) {
	insuranceCode := r.URL.Query().Get("insurance_code")
	data, err := h.usecase.GetPlanCommissionsFromTable(insuranceCode)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, data)
}

// GetDefaultConfigFromTable handles GET /api/commissions/table/default_config_products
func (h *CommissionHandler) GetDefaultConfigFromTable(w http.ResponseWriter, r *http.Request) {
	insuranceCode := r.URL.Query().Get("insurance_code")
	data, err := h.usecase.GetDefaultConfigFromTable(insuranceCode)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, data)
}

// CheckDuplicateProductCodes handles POST /api/commissions/check-duplicates
func (h *CommissionHandler) CheckDuplicateProductCodes(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ProductCodes []string `json:"product_codes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	duplicates, err := h.usecase.CheckDuplicateProductCodes(request.ProductCodes)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, duplicates)
}
