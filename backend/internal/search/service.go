package search

import (
	"context"
	"strings"

	"watchlist-backend/pkg/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func normalizeSearchQuery(query string) string {
	return strings.TrimSpace(query)
}

func (s *Service) SearchStocks(
	ctx context.Context,
	query string,
) ([]models.Stock, error) {
	query = normalizeSearchQuery(query)
	if query == "" {
		return []models.Stock{}, nil
	}

	return s.repo.SearchStocks(ctx, query)
}
