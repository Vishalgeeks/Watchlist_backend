package stock

import (
	"context"
	"errors"

	"watchlist-backend/pkg/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateStock(ctx context.Context, stock *models.Stock) error {
	if stock.Symbol == "" {
		return errors.New("symbol is required")
	}
	if stock.CompanyName == "" {
		return errors.New("company name is required")
	}
	return s.repo.Create(ctx, stock)
}

func (s *Service) GetAllStocks(ctx context.Context) ([]models.Stock, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetStockByID(ctx context.Context, id string) (*models.Stock, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) UpdateStock(ctx context.Context, id string, stock *models.Stock) error {
	if stock.Symbol == "" {
		return errors.New("symbol is required")
	}
	if stock.CompanyName == "" {
		return errors.New("company name is required")
	}
	return s.repo.Update(ctx, id, stock)
}

func (s *Service) DeleteStock(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
