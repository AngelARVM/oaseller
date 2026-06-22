package merchants

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
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
			"error": "unable to read merchants",
		})
		return
	}

	writeJSON(w, http.StatusOK, merchants)
}

func (h *Handler) Merchant(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	merchantID, err := parseMerchantID(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid merchant id",
		})
		return
	}

	merchant, err := h.service.Merchant(r.Context(), merchantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "merchant not found",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "unable to read merchant",
		})
		return
	}

	writeJSON(w, http.StatusOK, merchant)
}

func (h *Handler) PatchMerchant(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	merchantID, err := parseMerchantID(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid merchant id",
		})
		return
	}

	var req UpdateMerchantRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	merchant, err := h.service.PatchMerchant(r.Context(), merchantID, req)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "merchant not found",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "unable to update merchant",
		})
		return
	}

	writeJSON(w, http.StatusOK, merchant)
}

func parseMerchantID(id string) (int64, error) {
	merchantID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, err
	}

	if merchantID <= 0 {
		return 0, errors.New("invalid merchant id")
	}

	return merchantID, nil
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(data)
}
