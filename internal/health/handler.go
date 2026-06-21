package health

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

func (h *Handler) Live(w http.ResponseWriter, r *http.Request) {
	response := h.service.Live(r.Context())

	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	response := h.service.Ready(r.Context())

	statusCode := http.StatusOK
	if response.Status != "ok" {
		statusCode = http.StatusServiceUnavailable
	}
	writeJSON(w, statusCode, response)
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(data)
}
