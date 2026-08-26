package portfolio

import (
	"context"
	"database/sql"
	"errors"
	"watchlist-backend/pkg/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByUserID(ctx context.Context, userID int) ([]models.Portfolio, error) {
	query := `
		SELECT p.id, p.user_id, p.stock_id, p.quantity, p.avg_price, p.created_at, p.updated_at,
		       s.id, s.exchange_instrument_id, s.segment, s.instrument_type, s.symbol, s.display_name,
		       s.company_name, s.isin, s.series, s.exchange, s.contract_expiration, s.strike,
		       s.option_type, s.underlying_symbol_id, s.underlying_symbol, s.lot_size, s.tick_size,
		       s.upper_circuit, s.lower_circuit, s.freeze_qty, s.description, s.ltp, s.open, s.high,
		       s.low, s.close, s.vol, s.oi, s.bid, s.ask, s.bid_qty, s.ask_qty,
		       s.cautionary_message_info, s.last_updated
		FROM portfolios p
		JOIN stocks s ON s.id = p.stock_id
		WHERE p.user_id = $1 AND p.quantity > 0
		ORDER BY p.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var holdings []models.Portfolio
	for rows.Next() {
		var p models.Portfolio
		var s models.Stock
		err := rows.Scan(
			&p.ID, &p.UserID, &p.StockID, &p.Quantity, &p.AvgPrice, &p.CreatedAt, &p.UpdatedAt,
			&s.ID, &s.ExchangeInstrumentID, &s.Segment, &s.InstrumentType, &s.Symbol, &s.DisplayName,
			&s.CompanyName, &s.ISIN, &s.Series, &s.Exchange, &s.ContractExpiration, &s.Strike,
			&s.OptionType, &s.UnderlyingSymbolID, &s.UnderlyingSymbol, &s.LotSize, &s.TickSize,
			&s.UpperCircuit, &s.LowerCircuit, &s.FreezeQty, &s.Description,
			&s.LTP, &s.Open, &s.High, &s.Low, &s.Close, &s.Vol, &s.OI, &s.Bid, &s.Ask, &s.BidQty, &s.AskQty,
			&s.CautionaryMessageInfo, &s.LastUpdated,
		)
		if err != nil {
			return nil, err
		}
		p.Stock = &s
		holdings = append(holdings, p)
	}
	return holdings, nil
}

func (r *Repository) GetHolding(ctx context.Context, userID, stockID int) (*models.Portfolio, error) {
	query := `
		SELECT p.id, p.user_id, p.stock_id, p.quantity, p.avg_price, p.created_at, p.updated_at,
		       s.id, s.exchange_instrument_id, s.segment, s.instrument_type, s.symbol, s.display_name,
		       s.company_name, s.isin, s.series, s.exchange, s.contract_expiration, s.strike,
		       s.option_type, s.underlying_symbol_id, s.underlying_symbol, s.lot_size, s.tick_size,
		       s.upper_circuit, s.lower_circuit, s.freeze_qty, s.description, s.ltp, s.open, s.high,
		       s.low, s.close, s.vol, s.oi, s.bid, s.ask, s.bid_qty, s.ask_qty,
		       s.cautionary_message_info, s.last_updated
		FROM portfolios p
		JOIN stocks s ON s.id = p.stock_id
		WHERE p.user_id = $1 AND p.stock_id = $2 AND p.quantity > 0
	`
	var p models.Portfolio
	var s models.Stock
	err := r.db.QueryRowContext(ctx, query, userID, stockID).Scan(
		&p.ID, &p.UserID, &p.StockID, &p.Quantity, &p.AvgPrice, &p.CreatedAt, &p.UpdatedAt,
		&s.ID, &s.ExchangeInstrumentID, &s.Segment, &s.InstrumentType, &s.Symbol, &s.DisplayName,
		&s.CompanyName, &s.ISIN, &s.Series, &s.Exchange, &s.ContractExpiration, &s.Strike,
		&s.OptionType, &s.UnderlyingSymbolID, &s.UnderlyingSymbol, &s.LotSize, &s.TickSize,
		&s.UpperCircuit, &s.LowerCircuit, &s.FreezeQty, &s.Description,
		&s.LTP, &s.Open, &s.High, &s.Low, &s.Close, &s.Vol, &s.OI, &s.Bid, &s.Ask, &s.BidQty, &s.AskQty,
		&s.CautionaryMessageInfo, &s.LastUpdated,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("holding not found")
	}
	if err != nil {
		return nil, err
	}
	p.Stock = &s
	return &p, nil
}

func (r *Repository) AdjustHolding(ctx context.Context, userID, stockID int, qtyDelta int, price float64) error {
	if qtyDelta > 0 {
		_, err := r.db.ExecContext(ctx,
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

	_, err := r.db.ExecContext(ctx,
		`UPDATE portfolios SET quantity = quantity + $3, updated_at = NOW()
		 WHERE user_id = $1 AND stock_id = $2 AND quantity >= ABS($3)`,
		userID, stockID, qtyDelta,
	)
	return err
}

func (r *Repository) AdjustHoldingWithinTx(ctx context.Context, tx *sql.Tx, userID, stockID int, qtyDelta int, price float64) error {
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

	var newQty int
	err := tx.QueryRowContext(ctx,
		`UPDATE portfolios SET quantity = quantity + $3, updated_at = NOW()
		 WHERE user_id = $1 AND stock_id = $2 AND quantity >= ABS($3)
		 RETURNING quantity`,
		userID, stockID, qtyDelta,
	).Scan(&newQty)
	if err == sql.ErrNoRows {
		return errors.New("insufficient holdings")
	}
	if err != nil {
		return err
	}
	if newQty == 0 {
		_, err = tx.ExecContext(ctx,
			`DELETE FROM portfolios WHERE user_id = $1 AND stock_id = $2`,
			userID, stockID,
		)
		return err
	}
	return nil
}

func (r *Repository) RemoveHolding(ctx context.Context, userID, stockID int) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM portfolios WHERE user_id = $1 AND stock_id = $2`,
		userID, stockID,
	)
	return err
}
