package stock

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"watchlist-backend/cache"
	"watchlist-backend/pkg/models"
)

type Service struct {
	repo *Repository
}

type PaginatedStocks struct {
	Items      []models.Stock `json:"items"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	Total      int            `json:"total"`
	TotalPages int            `json:"total_pages"`
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func stocksListKey() string {
	return "stocks:list"
}

func stockKey(id string) string {
	return fmt.Sprintf("stock:%s", id)
}

func normalizePagination(pageStr, limitStr string) (int, int, error) {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 50
	}

	return page, limit, nil
}

func (s *Service) GetAllStocks(ctx context.Context, pageStr, limitStr string) (*PaginatedStocks, error) {
	page, limit, err := normalizePagination(pageStr, limitStr)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("stocks:list:%d:%d", page, limit)

	var cached PaginatedStocks
	if hit, err := cache.GetJSON(ctx, key, &cached); hit && err == nil {
		return &cached, nil
	}

	stocks, total, err := s.repo.GetAllPaginated(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	result := &PaginatedStocks{
		Items:      stocks,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: 0,
	}
	if limit > 0 {
		result.TotalPages = (total + limit - 1) / limit
	}

	cache.SetJSON(ctx, key, result, 5*time.Minute)
	return result, nil
}

func (s *Service) GetStockByID(ctx context.Context, id string) (*models.Stock, error) {
	key := stockKey(id)

	// Try cache first
	var cached models.Stock
	if hit, err := cache.GetJSON(ctx, key, &cached); hit && err == nil {
		return &cached, nil
	}

	// Cache miss - fetch from DB
	stock, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Store in cache (TTL 60s - short to bound live-data staleness)
	cache.SetJSON(ctx, key, stock, 60*time.Second)
	return stock, nil
}

func (s *Service) CreateStock(ctx context.Context, stock *models.Stock) error {
	if stock.Symbol == "" {
		return errors.New("symbol is required")
	}
	if stock.CompanyName == "" {
		return errors.New("company name is required")
	}
	if err := s.repo.Create(ctx, stock); err != nil {
		return err
	}
	// Invalidate stocks list cache
	cache.Del(ctx, stocksListKey())
	return nil
}

func (s *Service) UpdateStock(ctx context.Context, id string, stock *models.Stock) error {
	if stock.Symbol == "" {
		return errors.New("symbol is required")
	}
	if stock.CompanyName == "" {
		return errors.New("company name is required")
	}
	if err := s.repo.Update(ctx, id, stock); err != nil {
		return err
	}
	// Invalidate both list and specific stock
	cache.Del(ctx, stocksListKey(), stockKey(id))
	return nil
}

func (s *Service) DeleteStock(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	// Invalidate both list and specific stock
	cache.Del(ctx, stocksListKey(), stockKey(id))
	return nil
}
