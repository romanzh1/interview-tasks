package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"route256/cart/internal/models"

	"github.com/go-playground/validator/v10"
)

type cartService interface {
	AddItemToUserCart(ctx context.Context, cart models.CartRequest) error
	DeleteItemFromUserCart(ctx context.Context, userID, skuID int64) error
	ClearUserCart(ctx context.Context, userID int64) error
	ListUserCart(ctx context.Context, userID int64) ([]models.CartItem, uint32, error)
}

type Handler struct {
	service  cartService
	validate *validator.Validate
}

func NewHandler(service cartService) *Handler {
	return &Handler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /user/{userID}/cart/{skuID}", h.AddItemToUserCart)
	mux.HandleFunc("DELETE /user/{userID}/cart/{skuID}", h.DeleteItemFromUserCart)
	mux.HandleFunc("DELETE /user/{userID}/cart", h.ClearUserCart)
	mux.HandleFunc("GET /user/{userID}/cart", h.ListUserCart)
}

func (h *Handler) extractUserID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	userID, err := strconv.ParseInt(r.PathValue("userID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return 0, false
	}

	return userID, true
}

func (h *Handler) extractPathParams(w http.ResponseWriter, r *http.Request) (int64, int64, bool) {
	userID, err := strconv.ParseInt(r.PathValue("userID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return 0, 0, false
	}

	skuID, err := strconv.ParseInt(r.PathValue("skuID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid SKU ID", http.StatusBadRequest)
		return 0, 0, false
	}

	return userID, skuID, true
}

func (h *Handler) AddItemToUserCart(w http.ResponseWriter, r *http.Request) {
	userID, skuID, ok := h.extractPathParams(w, r)
	if !ok {
		return
	}

	req := models.CartRequest{
		UserID: userID,
		SkuID:  skuID,
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if err := h.service.AddItemToUserCart(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteItemFromUserCart(w http.ResponseWriter, r *http.Request) {
	userID, skuID, ok := h.extractPathParams(w, r)
	if !ok {
		return
	}

	req := models.CartRequest{
		UserID: userID,
		SkuID:  skuID,
	}

	if err := h.validate.StructPartial(req, "UserID", "SkuID"); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteItemFromUserCart(r.Context(), req.UserID, req.SkuID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ClearUserCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.extractUserID(w, r)
	if !ok {
		return
	}

	req := models.CartRequest{
		UserID: userID,
	}

	if err := h.validate.StructPartial(req, "UserID"); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if err := h.service.ClearUserCart(r.Context(), req.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListUserCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.extractUserID(w, r)
	if !ok {
		return
	}

	req := models.CartRequest{
		UserID: userID,
	}

	if err := h.validate.StructPartial(req, "UserID"); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	items, totalPrice, err := h.service.ListUserCart(r.Context(), req.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(items) == 0 {
		http.Error(w, "cart not found", http.StatusNotFound)
		return
	}

	response := struct {
		Items      []models.CartItem `json:"items"`
		TotalPrice uint32            `json:"total_price"`
	}{
		Items:      items,
		TotalPrice: totalPrice,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
