package handler

import (
	"net/http"

	"github.com/PurePras/shri-ram-service/backend/internal/dto"
	"github.com/PurePras/shri-ram-service/backend/internal/service"
)

type DeliverySlotHandler struct {
	svc *service.DeliverySlotService
}

func NewDeliverySlotHandler(svc *service.DeliverySlotService) *DeliverySlotHandler {
	return &DeliverySlotHandler{svc: svc}
}

func (h *DeliverySlotHandler) RegisterRoutes(mux *http.ServeMux, requireAdmin RouteGuard) {
	mux.HandleFunc("GET /delivery-slots", h.List)
	mux.HandleFunc("GET /delivery-slots/{id}", h.Get)
	mux.HandleFunc("POST /delivery-slots", requireAdmin(h.Create))
	mux.HandleFunc("PUT /delivery-slots/{id}", requireAdmin(h.Update))
	mux.HandleFunc("DELETE /delivery-slots/{id}", requireAdmin(h.Delete))
}

func (h *DeliverySlotHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)
	slots, err := h.svc.List(r.Context(), queryBoolPtr(r, "is_active"), limit, offset)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, slots)
}

func (h *DeliverySlotHandler) Get(w http.ResponseWriter, r *http.Request) {
	slot, err := h.svc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, slot)
}

func (h *DeliverySlotHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in dto.DeliverySlotInput
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	slot, err := h.svc.Create(r.Context(), in)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, slot)
}

func (h *DeliverySlotHandler) Update(w http.ResponseWriter, r *http.Request) {
	var in dto.DeliverySlotInput
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	slot, err := h.svc.Update(r.Context(), r.PathValue("id"), in)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, slot)
}

func (h *DeliverySlotHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Delete(r.Context(), r.PathValue("id")); err != nil {
		handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
