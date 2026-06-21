package merchants

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateMerchant(w http.ResponseWriter, r *http.Request) {
	var req CreateMerchantRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	merchant, err := h.service.CreateMerchant(r.Context(), req)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, merchant)
}

func (h *Handler) ListMerchants(w http.ResponseWriter, r *http.Request) {
	merchants, err := h.service.ListMerchants(r.Context())

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Unable to read merchants",
		})

		return
	}

	writeJSON(w, http.StatusOK, merchants)
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(data)
}
