package search

import (
	"context"

	"watchlist-backend/pkg/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SearchStocks(
	ctx context.Context,
	query string,
) ([]models.Stock, error) {

	return s.repo.SearchStocks(ctx, query)
}
