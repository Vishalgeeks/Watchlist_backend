package order

import (
	"encoding/json"
	"net/http"
	"strconv"
	"watchlist-backend/pkg/models"
	"watchlist-backend/pkg/validator"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// POST /api/orders
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(int)

	var req models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.Response{
			Success: false,
			Message: "invalid request body",
		})
		return
	}

	if errs := validator.Validate(req); len(errs) > 0 {
		writeJSON(w, http.StatusBadRequest, models.Response{
			Success: false,
			Message: "validation failed",
			Data:    errs,
		})
		return
	}

	order, err := h.service.CreateOrder(ctx, userID, &req)
	if err != nil {
		msg := err.Error()
		if msg == "stock not found" {
			writeJSON(w, http.StatusBadRequest, models.Response{
				Success: false,
				Message: msg,
			})
			return
		}
		writeJSON(w, http.StatusBadRequest, models.Response{
			Success: false,
			Message: msg,
		})
		return
	}

	writeJSON(w, http.StatusCreated, models.Response{
		Success: true,
		Message: "order placed successfully",
		Data:    order,
	})
}

// GET /api/orders
func (h *Handler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(int)

	orders, err := h.service.GetAllOrders(ctx, userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "failed to fetch orders",
		})
		return
	}

	writeJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "orders fetched successfully",
		Data:    orders,
	})
}

// GET /api/orders/{id}
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(int)
	vars := mux.Vars(r)

	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.Response{
			Success: false,
			Message: "invalid order id",
		})
		return
	}

	order, err := h.service.GetOrder(ctx, userID, orderID)
	if err != nil {
		msg := err.Error()
		if msg == "order not found" {
			writeJSON(w, http.StatusNotFound, models.Response{
				Success: false,
				Message: msg,
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "failed to fetch order",
		})
		return
	}

	writeJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "order fetched successfully",
		Data:    order,
	})
}

// PUT /api/orders/{id}/cancel
func (h *Handler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(int)
	vars := mux.Vars(r)

	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.Response{
			Success: false,
			Message: "invalid order id",
		})
		return
	}

	if err := h.service.CancelOrder(ctx, userID, orderID); err != nil {
		msg := err.Error()
		if msg == "order cannot be cancelled" {
			writeJSON(w, http.StatusBadRequest, models.Response{
				Success: false,
				Message: msg,
			})
			return
		}
		if msg == "order not found" {
			writeJSON(w, http.StatusNotFound, models.Response{
				Success: false,
				Message: msg,
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "failed to cancel order",
		})
		return
	}

	writeJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "order cancelled successfully",
	})
}