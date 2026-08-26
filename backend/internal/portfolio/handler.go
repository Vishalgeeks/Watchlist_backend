package portfolio

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

// GET /api/portfolio
func (h *Handler) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(int)

	holdings, err := h.service.GetEnrichedHoldings(ctx, userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "failed to fetch portfolio",
		})
		return
	}

	writeJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "portfolio fetched successfully",
		Data:    holdings,
	})
}

// GET /api/portfolio/{id}
func (h *Handler) GetHolding(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(int)
	vars := mux.Vars(r)

	stockID, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.Response{
			Success: false,
			Message: "invalid stock id",
		})
		return
	}

	holding, err := h.service.GetHolding(ctx, userID, stockID)
	if err != nil {
		msg := err.Error()
		if msg == "holding not found" {
			writeJSON(w, http.StatusNotFound, models.Response{
				Success: false,
				Message: msg,
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "failed to fetch holding",
		})
		return
	}

	if holding.Stock != nil && holding.Stock.LTP > 0 {
		holding.CurrentValue = float64(holding.Quantity) * holding.Stock.LTP
		holding.UnrealizedPnL = (holding.Stock.LTP - holding.AvgPrice) * float64(holding.Quantity)
		if holding.AvgPrice > 0 {
			holding.PnLPercentage = ((holding.Stock.LTP - holding.AvgPrice) / holding.AvgPrice) * 100
		}
	}

	writeJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "holding fetched successfully",
		Data:    holding,
	})
}

// GET /api/portfolio/summary
func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(int)

	summary, err := h.service.GetPortfolioSummary(ctx, userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "failed to fetch portfolio summary",
		})
		return
	}

	writeJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "portfolio summary fetched successfully",
		Data:    summary,
	})
}
