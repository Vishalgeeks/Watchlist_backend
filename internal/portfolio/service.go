package portfolio

import (
	"context"
	"database/sql"
	"watchlist-backend/pkg/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetHoldings(ctx context.Context, userID int) ([]models.Portfolio, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *Service) AddShares(ctx context.Context, userID, stockID int, qty int, price float64) error {
	return s.repo.AdjustHolding(ctx, userID, stockID, qty, price)
}

func (s *Service) RemoveShares(ctx context.Context, userID, stockID int, qty int) error {
	return s.repo.AdjustHolding(ctx, userID, stockID, -qty, 0)
}

func (s *Service) AdjustHoldingWithinTx(ctx context.Context, tx *sql.Tx, userID, stockID int, qtyDelta int, price float64) error {
	return s.repo.AdjustHoldingWithinTx(ctx, tx, userID, stockID, qtyDelta, price)
}

func (s *Service) GetHolding(ctx context.Context, userID, stockID int) (*models.Portfolio, error) {
	return s.repo.GetHolding(ctx, userID, stockID)
}

func (s *Service) GetEnrichedHoldings(ctx context.Context, userID int) ([]models.Portfolio, error) {
	holdings, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if holdings == nil {
		return []models.Portfolio{}, nil
	}
	for i := range holdings {
		h := &holdings[i]
		if h.Stock != nil && h.Stock.LTP > 0 {
			h.CurrentValue = float64(h.Quantity) * h.Stock.LTP
			h.UnrealizedPnL = (h.Stock.LTP - h.AvgPrice) * float64(h.Quantity)
			if h.AvgPrice > 0 {
				h.PnLPercentage = ((h.Stock.LTP - h.AvgPrice) / h.AvgPrice) * 100
			}
		}
	}
	return holdings, nil
}

func (s *Service) GetPortfolioSummary(ctx context.Context, userID int) (*models.PortfolioSummary, error) {
	holdings, err := s.GetEnrichedHoldings(ctx, userID)
	if err != nil {
		return nil, err
	}
	summary := &models.PortfolioSummary{}
	for _, h := range holdings {
		summary.TotalValue += h.CurrentValue
		summary.TotalCostBasis += h.AvgPrice * float64(h.Quantity)
		summary.HoldingsCount++
	}
	summary.TotalUnrealizedPnL = summary.TotalValue - summary.TotalCostBasis
	if summary.TotalCostBasis > 0 {
		summary.TotalPnLPercentage = (summary.TotalUnrealizedPnL / summary.TotalCostBasis) * 100
	}
	return summary, nil
}
