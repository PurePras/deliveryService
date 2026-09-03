package handler

import (
	"net/http"

	"github.com/PurePras/shri-ram-service/backend/internal/dto"
	"github.com/PurePras/shri-ram-service/backend/internal/middleware"
	"github.com/PurePras/shri-ram-service/backend/internal/model"
	"github.com/PurePras/shri-ram-service/backend/internal/service"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) RegisterRoutes(mux *http.ServeMux, requireAuth, requireAdmin RouteGuard) {
	mux.HandleFunc("POST /orders", requireAuth(h.Create))
	mux.HandleFunc("GET /orders", requireAuth(h.List))
	mux.HandleFunc("GET /orders/{id}", requireAuth(h.Get))
	mux.HandleFunc("PATCH /orders/{id}/status", requireAdmin(h.UpdateStatus))
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	authUser, ok := currentUser(w, r)
	if !ok {
		return
	}

	var in dto.CreateOrderRequest
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	order, err := h.svc.Create(r.Context(), authUser.ID, in)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	authUser, ok := currentUser(w, r)
	if !ok {
		return
	}

	order, err := h.svc.Get(r.Context(), r.PathValue("id"), authUser.ID, authUser.Role == model.RoleAdmin)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	authUser, ok := currentUser(w, r)
	if !ok {
		return
	}

	limit, offset := parsePagination(r)
	params := service.OrderListParams{
		Status: queryStringPtr(r, "status"),
		Limit:  limit,
		Offset: offset,
	}
	orders, err := h.svc.List(r.Context(), authUser.ID, authUser.Role == model.RoleAdmin, params)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	var in dto.UpdateOrderStatusRequest
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	order, err := h.svc.UpdateStatus(r.Context(), r.PathValue("id"), in)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, order)
}

// currentUser reads the authenticated user injected by requireAuth/requireAdmin,
// writing a 401 and returning ok=false if it's somehow missing.
func currentUser(w http.ResponseWriter, r *http.Request) (*middleware.AuthUser, bool) {
	authUser, ok := middleware.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return nil, false
	}
	return authUser, true
}
