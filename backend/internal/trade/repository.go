package trade

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

func (r *Repository) Create(ctx context.Context, trade *models.Trade) error {
	query := `
		INSERT INTO trades (order_id, user_id, stock_id, side, quantity, execution_price, total_amount, executed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		RETURNING id, executed_at
	`
	return r.db.QueryRowContext(ctx, query,
		trade.OrderID, trade.UserID, trade.StockID, trade.Side, trade.Quantity,
		trade.ExecutionPrice, trade.TotalAmount,
	).Scan(&trade.ID, &trade.ExecutedAt)
}

func (r *Repository) CreateWithinTx(ctx context.Context, tx *sql.Tx, trade *models.Trade) error {
	query := `
		INSERT INTO trades (order_id, user_id, stock_id, side, quantity, execution_price, total_amount, executed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		RETURNING id, executed_at
	`
	return tx.QueryRowContext(ctx, query,
		trade.OrderID, trade.UserID, trade.StockID, trade.Side, trade.Quantity,
		trade.ExecutionPrice, trade.TotalAmount,
	).Scan(&trade.ID, &trade.ExecutedAt)
}

func (r *Repository) GetByOrderID(ctx context.Context, orderID, userID int) (*models.Trade, error) {
	query := `
		SELECT t.id, t.order_id, t.user_id, t.stock_id, t.side, t.quantity,
		       t.execution_price, t.total_amount, t.executed_at,
		       s.id, s.exchange_instrument_id, s.segment, s.instrument_type, s.symbol, s.display_name,
		       s.company_name, s.isin, s.series, s.exchange, s.contract_expiration, s.strike,
		       s.option_type, s.underlying_symbol_id, s.underlying_symbol, s.lot_size, s.tick_size,
		       s.upper_circuit, s.lower_circuit, s.freeze_qty, s.description, s.ltp, s.open, s.high,
		       s.low, s.close, s.vol, s.oi, s.bid, s.ask, s.bid_qty, s.ask_qty,
		       s.cautionary_message_info, s.last_updated
		FROM trades t
		JOIN stocks s ON s.id = t.stock_id
		WHERE t.order_id = $1 AND t.user_id = $2
	`
	var t models.Trade
	var s models.Stock
	err := r.db.QueryRowContext(ctx, query, orderID, userID).Scan(
		&t.ID, &t.OrderID, &t.UserID, &t.StockID, &t.Side, &t.Quantity,
		&t.ExecutionPrice, &t.TotalAmount, &t.ExecutedAt,
		&s.ID, &s.ExchangeInstrumentID, &s.Segment, &s.InstrumentType, &s.Symbol, &s.DisplayName,
		&s.CompanyName, &s.ISIN, &s.Series, &s.Exchange, &s.ContractExpiration, &s.Strike,
		&s.OptionType, &s.UnderlyingSymbolID, &s.UnderlyingSymbol, &s.LotSize, &s.TickSize,
		&s.UpperCircuit, &s.LowerCircuit, &s.FreezeQty, &s.Description,
		&s.LTP, &s.Open, &s.High, &s.Low, &s.Close, &s.Vol, &s.OI, &s.Bid, &s.Ask, &s.BidQty, &s.AskQty,
		&s.CautionaryMessageInfo, &s.LastUpdated,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("trade not found")
	}
	if err != nil {
		return nil, err
	}
	t.Stock = &s
	return &t, nil
}

func (r *Repository) GetByID(ctx context.Context, tradeID, userID int) (*models.Trade, error) {
	query := `
		SELECT t.id, t.order_id, t.user_id, t.stock_id, t.side, t.quantity,
		       t.execution_price, t.total_amount, t.executed_at,
		       s.id, s.exchange_instrument_id, s.segment, s.instrument_type, s.symbol, s.display_name,
		       s.company_name, s.isin, s.series, s.exchange, s.contract_expiration, s.strike,
		       s.option_type, s.underlying_symbol_id, s.underlying_symbol, s.lot_size, s.tick_size,
		       s.upper_circuit, s.lower_circuit, s.freeze_qty, s.description, s.ltp, s.open, s.high,
		       s.low, s.close, s.vol, s.oi, s.bid, s.ask, s.bid_qty, s.ask_qty,
		       s.cautionary_message_info, s.last_updated
		FROM trades t
		JOIN stocks s ON s.id = t.stock_id
		WHERE t.id = $1 AND t.user_id = $2
	`
	var t models.Trade
	var s models.Stock
	err := r.db.QueryRowContext(ctx, query, tradeID, userID).Scan(
		&t.ID, &t.OrderID, &t.UserID, &t.StockID, &t.Side, &t.Quantity,
		&t.ExecutionPrice, &t.TotalAmount, &t.ExecutedAt,
		&s.ID, &s.ExchangeInstrumentID, &s.Segment, &s.InstrumentType, &s.Symbol, &s.DisplayName,
		&s.CompanyName, &s.ISIN, &s.Series, &s.Exchange, &s.ContractExpiration, &s.Strike,
		&s.OptionType, &s.UnderlyingSymbolID, &s.UnderlyingSymbol, &s.LotSize, &s.TickSize,
		&s.UpperCircuit, &s.LowerCircuit, &s.FreezeQty, &s.Description,
		&s.LTP, &s.Open, &s.High, &s.Low, &s.Close, &s.Vol, &s.OI, &s.Bid, &s.Ask, &s.BidQty, &s.AskQty,
		&s.CautionaryMessageInfo, &s.LastUpdated,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("trade not found")
	}
	if err != nil {
		return nil, err
	}
	t.Stock = &s
	return &t, nil
}

func (r *Repository) GetAllByUserID(ctx context.Context, userID int) ([]models.Trade, error) {
	query := `
		SELECT t.id, t.order_id, t.user_id, t.stock_id, t.side, t.quantity,
		       t.execution_price, t.total_amount, t.executed_at,
		       s.id, s.exchange_instrument_id, s.segment, s.instrument_type, s.symbol, s.display_name,
		       s.company_name, s.isin, s.series, s.exchange, s.contract_expiration, s.strike,
		       s.option_type, s.underlying_symbol_id, s.underlying_symbol, s.lot_size, s.tick_size,
		       s.upper_circuit, s.lower_circuit, s.freeze_qty, s.description, s.ltp, s.open, s.high,
		       s.low, s.close, s.vol, s.oi, s.bid, s.ask, s.bid_qty, s.ask_qty,
		       s.cautionary_message_info, s.last_updated
		FROM trades t
		JOIN stocks s ON s.id = t.stock_id
		WHERE t.user_id = $1
		ORDER BY t.executed_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []models.Trade
	for rows.Next() {
		var t models.Trade
		var s models.Stock
		err := rows.Scan(
			&t.ID, &t.OrderID, &t.UserID, &t.StockID, &t.Side, &t.Quantity,
			&t.ExecutionPrice, &t.TotalAmount, &t.ExecutedAt,
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
		t.Stock = &s
		trades = append(trades, t)
	}
	return trades, nil
}
