package handler

import (
	"net/http"

	"github.com/PurePras/shri-ram-service/backend/internal/dto"
	"github.com/PurePras/shri-ram-service/backend/internal/service"
)

type DeliveryAreaHandler struct {
	svc *service.DeliveryAreaService
}

func NewDeliveryAreaHandler(svc *service.DeliveryAreaService) *DeliveryAreaHandler {
	return &DeliveryAreaHandler{svc: svc}
}

func (h *DeliveryAreaHandler) RegisterRoutes(mux *http.ServeMux, requireAdmin RouteGuard) {
	mux.HandleFunc("GET /delivery-areas", h.List)
	mux.HandleFunc("GET /delivery-areas/{id}", h.Get)
	mux.HandleFunc("POST /delivery-areas", requireAdmin(h.Create))
	mux.HandleFunc("PUT /delivery-areas/{id}", requireAdmin(h.Update))
	mux.HandleFunc("DELETE /delivery-areas/{id}", requireAdmin(h.Delete))
}

func (h *DeliveryAreaHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)
	areas, err := h.svc.List(r.Context(), queryBoolPtr(r, "is_active"), limit, offset)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, areas)
}

func (h *DeliveryAreaHandler) Get(w http.ResponseWriter, r *http.Request) {
	area, err := h.svc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, area)
}

func (h *DeliveryAreaHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in dto.DeliveryAreaInput
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	area, err := h.svc.Create(r.Context(), in)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, area)
}

func (h *DeliveryAreaHandler) Update(w http.ResponseWriter, r *http.Request) {
	var in dto.DeliveryAreaInput
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	area, err := h.svc.Update(r.Context(), r.PathValue("id"), in)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, area)
}

func (h *DeliveryAreaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Delete(r.Context(), r.PathValue("id")); err != nil {
		handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
