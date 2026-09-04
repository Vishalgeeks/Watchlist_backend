package stock

import (
	"encoding/json"
	"net/http"

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

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/stocks", h.CreateStock).Methods("POST")
	r.HandleFunc("/stocks", h.GetAllStocks).Methods("GET")
	r.HandleFunc("/stocks/{id}", h.GetStockByID).Methods("GET")
	r.HandleFunc("/stocks/{id}", h.UpdateStock).Methods("PUT")
	r.HandleFunc("/stocks/{id}", h.DeleteStock).Methods("DELETE")
}

func (h *Handler) CreateStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var stock models.Stock
	if err := json.NewDecoder(r.Body).Decode(&stock); err != nil {
		writeJSON(w, http.StatusBadRequest, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if err := h.service.CreateStock(ctx, &stock); err != nil {
		writeJSON(w, http.StatusBadRequest, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, models.Response{
		Success: true,
		Message: "stock created successfully",
		Data:    stock,
	})
}

func (h *Handler) GetAllStocks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	stocks, err := h.service.GetAllStocks(ctx, page, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "stocks fetched successfully",
		Data:    stocks,
	})
}

func (h *Handler) GetStockByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	stock, err := h.service.GetStockByID(ctx, id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, models.Response{
			Success: false,
			Message: "stock not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "stock fetched successfully",
		Data:    stock,
	})
}

func (h *Handler) UpdateStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	var stock models.Stock
	if err := json.NewDecoder(r.Body).Decode(&stock); err != nil {
		writeJSON(w, http.StatusBadRequest, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if err := h.service.UpdateStock(ctx, id, &stock); err != nil {
		writeJSON(w, http.StatusBadRequest, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "stock updated successfully",
		Data:    stock,
	})
}

func (h *Handler) DeleteStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	if err := h.service.DeleteStock(ctx, id); err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, models.Response{
		Success: true,
		Message: "stock deleted successfully",
	})
}
