package wallet

import (
	"encoding/json"
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

// GET /api/wallet
func (h *Handler) GetWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(int)

	wallet, err := h.service.GetWallet(ctx, userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "failed to fetch wallet",
		})
		return
	}

	writeJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "wallet fetched successfully",
		Data:    wallet,
	})
}

// POST /api/wallet/deposit
func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(int)

	var req models.DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.Response{
			Success: false,
			Message: "invalid request body",
		})
		return
	}

	if req.Amount <= 0 {
		writeJSON(w, http.StatusBadRequest, models.Response{
			Success: false,
			Message: "amount must be greater than 0",
		})
		return
	}

	tx, err := h.service.Deposit(ctx, userID, req.Amount, req.Description)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "failed to deposit",
		})
		return
	}

	writeJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "deposit successful",
		Data:    tx,
	})
}

// POST /api/wallet/withdraw
func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(int)

	var req models.WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.Response{
			Success: false,
			Message: "invalid request body",
		})
		return
	}

	if req.Amount <= 0 {
		writeJSON(w, http.StatusBadRequest, models.Response{
			Success: false,
			Message: "amount must be greater than 0",
		})
		return
	}

	tx, err := h.service.Withdraw(ctx, userID, req.Amount, req.Description)
	if err != nil {
		msg := err.Error()
		if msg == "insufficient balance" {
			writeJSON(w, http.StatusBadRequest, models.Response{
				Success: false,
				Message: msg,
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "failed to withdraw",
		})
		return
	}

	writeJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "withdrawal successful",
		Data:    tx,
	})
}

// GET /api/wallet/transactions
func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(int)

	limit := 50
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	txs, err := h.service.GetTransactions(ctx, userID, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "failed to fetch transactions",
		})
		return
	}

	if txs == nil {
		txs = []models.WalletTransaction{}
	}

	writeJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "transactions fetched successfully",
		Data:    txs,
	})
}

// GET /api/wallet/transactions/{id}
func (h *Handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(int)
	vars := mux.Vars(r)

	transactionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.Response{
			Success: false,
			Message: "invalid transaction id",
		})
		return
	}

	tx, err := h.service.GetTransaction(ctx, userID, transactionID)
	if err != nil {
		msg := err.Error()
		if msg == "transaction not found" {
			writeJSON(w, http.StatusNotFound, models.Response{
				Success: false,
				Message: msg,
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "failed to fetch transaction",
		})
		return
	}

	writeJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "transaction fetched successfully",
		Data:    tx,
	})
}