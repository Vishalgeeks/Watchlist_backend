package search

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

func (r *Repository) SearchStocks(
	ctx context.Context,
	query string,
) ([]models.Stock, error) {

	query = "%" + query + "%"

	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT id, symbol, company_name, exchange, ltp, last_updated
		FROM stocks
		WHERE symbol ILIKE $1 OR company_name ILIKE $1
		LIMIT 20
	`,
		query,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Stock

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

		result = append(result, s)
	}

	return result, nil
}
