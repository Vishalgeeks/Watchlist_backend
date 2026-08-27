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
		SELECT id, COALESCE(exchange_instrument_id,''), COALESCE(segment,''), COALESCE(instrument_type,''), symbol,
		       COALESCE(display_name,''), COALESCE(company_name,''), COALESCE(isin,''), COALESCE(series,''),
		       COALESCE(exchange,''), COALESCE(contract_expiration,''), COALESCE(strike,0), COALESCE(option_type,''),
		       COALESCE(underlying_symbol_id,''), COALESCE(underlying_symbol,''), COALESCE(lot_size,0),
		       COALESCE(tick_size,0), COALESCE(upper_circuit,0), COALESCE(lower_circuit,0), COALESCE(freeze_qty,0),
		       COALESCE(description,''), ltp, open, high, low, COALESCE(close,0), vol, oi, bid, ask,
		       COALESCE(bid_qty,0), COALESCE(ask_qty,0), COALESCE(cautionary_message_info,''), last_updated
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
			&s.ID, &s.ExchangeInstrumentID, &s.Segment, &s.InstrumentType, &s.Symbol, &s.DisplayName,
			&s.CompanyName, &s.ISIN, &s.Series, &s.Exchange, &s.ContractExpiration, &s.Strike, &s.OptionType,
			&s.UnderlyingSymbolID, &s.UnderlyingSymbol, &s.LotSize, &s.TickSize, &s.UpperCircuit,
			&s.LowerCircuit, &s.FreezeQty, &s.Description, &s.LTP, &s.Open, &s.High, &s.Low, &s.Close,
			&s.Vol, &s.OI, &s.Bid, &s.Ask, &s.BidQty, &s.AskQty,
			&s.CautionaryMessageInfo, &s.LastUpdated,
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
		SELECT id, COALESCE(exchange_instrument_id,''), COALESCE(segment,''), COALESCE(instrument_type,''), symbol,
		       COALESCE(display_name,''), COALESCE(company_name,''), COALESCE(isin,''), COALESCE(series,''),
		       COALESCE(exchange,''), COALESCE(contract_expiration,''), COALESCE(strike,0), COALESCE(option_type,''),
		       COALESCE(underlying_symbol_id,''), COALESCE(underlying_symbol,''), COALESCE(lot_size,0),
		       COALESCE(tick_size,0), COALESCE(upper_circuit,0), COALESCE(lower_circuit,0), COALESCE(freeze_qty,0),
		       COALESCE(description,''), ltp, open, high, low, COALESCE(close,0), vol, oi, bid, ask,
		       COALESCE(bid_qty,0), COALESCE(ask_qty,0), COALESCE(cautionary_message_info,''), last_updated
		FROM stocks
		WHERE id = $1
	`
	var stock models.Stock
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&stock.ID, &stock.ExchangeInstrumentID, &stock.Segment, &stock.InstrumentType, &stock.Symbol, &stock.DisplayName,
		&stock.CompanyName, &stock.ISIN, &stock.Series, &stock.Exchange, &stock.ContractExpiration, &stock.Strike, &stock.OptionType,
		&stock.UnderlyingSymbolID, &stock.UnderlyingSymbol, &stock.LotSize, &stock.TickSize, &stock.UpperCircuit,
		&stock.LowerCircuit, &stock.FreezeQty, &stock.Description, &stock.LTP, &stock.Open, &stock.High, &stock.Low, &stock.Close,
		&stock.Vol, &stock.OI, &stock.Bid, &stock.Ask, &stock.BidQty, &stock.AskQty,
		&stock.CautionaryMessageInfo, &stock.LastUpdated,
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
