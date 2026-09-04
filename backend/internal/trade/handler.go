package trade

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"watchlist-backend/pkg/models"

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

// POST /api/orders/{id}/execute
func (h *Handler) ExecuteOrder(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		ExecutionPrice float64 `json:"execution_price"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	trade, err := h.service.ExecuteOrder(ctx, userID, orderID, req.ExecutionPrice)
	if err != nil {
		log.Printf("execute order failed: user_id=%d order_id=%d err=%v", userID, orderID, err)
		msg := err.Error()
		switch msg {
		case "order not found":
			writeJSON(w, http.StatusNotFound, models.Response{Success: false, Message: msg})
		case "order already executed", "order not executable", "insufficient balance", "insufficient holdings", "invalid execution price for limit order":
			writeJSON(w, http.StatusBadRequest, models.Response{Success: false, Message: msg})
		default:
			writeJSON(w, http.StatusInternalServerError, models.Response{Success: false, Message: "failed to execute order"})
		}
		return
	}

	writeJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "order executed successfully",
		Data:    trade,
	})
}

// GET /api/trades
func (h *Handler) GetAllTrades(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(int)

	trades, err := h.service.GetAllTrades(ctx, userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "failed to fetch trades",
		})
		return
	}

	writeJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "trades fetched successfully",
		Data:    trades,
	})
}

// GET /api/trades/{id}
func (h *Handler) GetTrade(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(int)
	vars := mux.Vars(r)

	tradeID, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.Response{
			Success: false,
			Message: "invalid trade id",
		})
		return
	}

	trade, err := h.service.GetTrade(ctx, userID, tradeID)
	if err != nil {
		msg := err.Error()
		if msg == "trade not found" {
			writeJSON(w, http.StatusNotFound, models.Response{
				Success: false,
				Message: msg,
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "failed to fetch trade",
		})
		return
	}

	writeJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "trade fetched successfully",
		Data:    trade,
	})
}
