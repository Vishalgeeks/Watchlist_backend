package stock

import (
	"context"
	"database/sql"

	"watchlist-backend/pkg/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, stock *models.Stock) error {
	query := `
		INSERT INTO stocks
		(symbol, company_name, exchange, LTP)
		VALUES ($1, $2, $3, $4)
		RETURNING id, last_updated
	`
	return r.db.QueryRowContext(ctx, query,
		stock.Symbol,
		stock.CompanyName,
		stock.Exchange,
		stock.LTP,
	).Scan(&stock.ID, &stock.LastUpdated)
}

func (r *Repository) GetAll(ctx context.Context) ([]models.Stock, error) {
	query := `
		SELECT id, symbol, company_name, exchange,
		       LTP, last_updated
		FROM stocks
		ORDER BY symbol
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var s models.Stock
		err := rows.Scan(
			&s.ID,
			&s.Symbol,
			&s.CompanyName,
			&s.Exchange,
			&s.LTP,
			&s.LastUpdated,
		)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, s)
	}
	return stocks, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*models.Stock, error) {
	query := `
		SELECT id, symbol, company_name,
		       exchange, LTP, last_updated
		FROM stocks
		WHERE id = $1
	`
	var stock models.Stock
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&stock.ID,
		&stock.Symbol,
		&stock.CompanyName,
		&stock.Exchange,
		&stock.LTP,
		&stock.LastUpdated,
	)
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *Repository) Update(ctx context.Context, id string, stock *models.Stock) error {
	query := `
		UPDATE stocks
		SET symbol = $1,
		    company_name = $2,
		    exchange = $3,
		    LTP = $4,
		    last_updated = NOW()
		WHERE id = $5
		RETURNING last_updated
	`
	return r.db.QueryRowContext(ctx, query,
		stock.Symbol,
		stock.CompanyName,
		stock.Exchange,
		stock.LTP,
		id,
	).Scan(&stock.LastUpdated)
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM stocks WHERE id = $1",
		id,
	)
	return err
}
