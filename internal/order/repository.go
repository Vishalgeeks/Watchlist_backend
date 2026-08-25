package order

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

func (r *Repository) Create(ctx context.Context, o *models.Order) error {
	query := `
		INSERT INTO orders (user_id, stock_id, side, order_type, quantity, price, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		o.UserID,
		o.StockID,
		o.Side,
		o.OrderType,
		o.Quantity,
		o.Price,
		o.Status,
	).Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)
}

func (r *Repository) GetByID(ctx context.Context, orderID, userID int) (*models.Order, error) {
	query := `
		SELECT
			o.id, o.user_id, o.stock_id, o.side, o.order_type, o.quantity, o.price, o.status, o.created_at, o.updated_at,
			s.id, s.exchange_instrument_id, s.segment, s.instrument_type, s.symbol, s.display_name, s.company_name, s.isin, s.series,
			s.exchange, s.contract_expiration, s.strike, s.option_type, s.underlying_symbol_id, s.underlying_symbol,
			s.lot_size, s.tick_size, s.upper_circuit, s.lower_circuit, s.freeze_qty, s.description,
			s.ltp, s.open, s.high, s.low, s.close, s.vol, s.oi, s.bid, s.ask, s.bid_qty, s.ask_qty,
			s.cautionary_message_info, s.last_updated
		FROM orders o
		JOIN stocks s ON s.id = o.stock_id
		WHERE o.id = $1 AND o.user_id = $2
	`
	var o models.Order
	var s models.Stock
	err := r.db.QueryRowContext(ctx, query, orderID, userID).Scan(
		&o.ID, &o.UserID, &o.StockID, &o.Side, &o.OrderType, &o.Quantity, &o.Price, &o.Status, &o.CreatedAt, &o.UpdatedAt,
		&s.ID, &s.ExchangeInstrumentID, &s.Segment, &s.InstrumentType, &s.Symbol, &s.DisplayName, &s.CompanyName, &s.ISIN, &s.Series,
		&s.Exchange, &s.ContractExpiration, &s.Strike, &s.OptionType, &s.UnderlyingSymbolID, &s.UnderlyingSymbol,
		&s.LotSize, &s.TickSize, &s.UpperCircuit, &s.LowerCircuit, &s.FreezeQty, &s.Description,
		&s.LTP, &s.Open, &s.High, &s.Low, &s.Close, &s.Vol, &s.OI, &s.Bid, &s.Ask, &s.BidQty, &s.AskQty,
		&s.CautionaryMessageInfo, &s.LastUpdated,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("order not found")
	}
	if err != nil {
		return nil, err
	}
	o.Stock = &s
	return &o, nil
}

func (r *Repository) GetAllByUserID(ctx context.Context, userID int) ([]models.Order, error) {
	query := `
		SELECT
			o.id, o.user_id, o.stock_id, o.side, o.order_type, o.quantity, o.price, o.status, o.created_at, o.updated_at,
			s.id, s.exchange_instrument_id, s.segment, s.instrument_type, s.symbol, s.display_name, s.company_name, s.isin, s.series,
			s.exchange, s.contract_expiration, s.strike, s.option_type, s.underlying_symbol_id, s.underlying_symbol,
			s.lot_size, s.tick_size, s.upper_circuit, s.lower_circuit, s.freeze_qty, s.description,
			s.ltp, s.open, s.high, s.low, s.close, s.vol, s.oi, s.bid, s.ask, s.bid_qty, s.ask_qty,
			s.cautionary_message_info, s.last_updated
		FROM orders o
		JOIN stocks s ON s.id = o.stock_id
		WHERE o.user_id = $1
		ORDER BY o.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		var s models.Stock
		err := rows.Scan(
			&o.ID, &o.UserID, &o.StockID, &o.Side, &o.OrderType, &o.Quantity, &o.Price, &o.Status, &o.CreatedAt, &o.UpdatedAt,
			&s.ID, &s.ExchangeInstrumentID, &s.Segment, &s.InstrumentType, &s.Symbol, &s.DisplayName, &s.CompanyName, &s.ISIN, &s.Series,
			&s.Exchange, &s.ContractExpiration, &s.Strike, &s.OptionType, &s.UnderlyingSymbolID, &s.UnderlyingSymbol,
			&s.LotSize, &s.TickSize, &s.UpperCircuit, &s.LowerCircuit, &s.FreezeQty, &s.Description,
			&s.LTP, &s.Open, &s.High, &s.Low, &s.Close, &s.Vol, &s.OI, &s.Bid, &s.Ask, &s.BidQty, &s.AskQty,
			&s.CautionaryMessageInfo, &s.LastUpdated,
		)
		if err != nil {
			return nil, err
		}
		o.Stock = &s
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *Repository) StockExists(ctx context.Context, stockID int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT 1 FROM stocks WHERE id = $1`, stockID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repository) CancelOrder(ctx context.Context, orderID, userID int) (bool, error) {
	result, err := r.db.ExecContext(ctx,
		`UPDATE orders SET status = 'CANCELLED', updated_at = NOW() WHERE id = $1 AND user_id = $2 AND status = 'PENDING'`,
		orderID, userID,
	)
	if err != nil {
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}