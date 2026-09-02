package handler

import (
	"net/http"

	"github.com/PurePras/shri-ram-service/backend/internal/dto"
	"github.com/PurePras/shri-ram-service/backend/internal/service"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) RegisterRoutes(mux *http.ServeMux, requireAdmin RouteGuard) {
	mux.HandleFunc("GET /categories", h.List)
	mux.HandleFunc("GET /categories/{id}", h.Get)
	mux.HandleFunc("POST /categories", requireAdmin(h.Create))
	mux.HandleFunc("PUT /categories/{id}", requireAdmin(h.Update))
	mux.HandleFunc("DELETE /categories/{id}", requireAdmin(h.Delete))
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)
	categories, err := h.svc.List(r.Context(), queryBoolPtr(r, "is_active"), limit, offset)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, categories)
}

func (h *CategoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	category, err := h.svc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, category)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in dto.CategoryInput
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	category, err := h.svc.Create(r.Context(), in)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, category)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	var in dto.CategoryInput
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	category, err := h.svc.Update(r.Context(), r.PathValue("id"), in)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, category)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Delete(r.Context(), r.PathValue("id")); err != nil {
		handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
