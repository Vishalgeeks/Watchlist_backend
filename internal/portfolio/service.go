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
	if qtyDelta > 0 {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO portfolios (user_id, stock_id, quantity, avg_price, updated_at)
			 VALUES ($1, $2, $3, $4, NOW())
			 ON CONFLICT (user_id, stock_id) DO UPDATE SET
			   quantity = portfolios.quantity + $3,
			   avg_price = (portfolios.quantity * portfolios.avg_price + $3 * $4) / (portfolios.quantity + $3),
			   updated_at = NOW()`,
			userID, stockID, qtyDelta, price,
		)
		return err
	}

	_, err := tx.ExecContext(ctx,
		`UPDATE portfolios SET quantity = quantity + $3, updated_at = NOW()
		 WHERE user_id = $1 AND stock_id = $2 AND quantity >= ABS($3)`,
		userID, stockID, qtyDelta,
	)
	return err
}
