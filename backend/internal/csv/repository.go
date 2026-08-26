package csv

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

func (r *Repository) LoadFromURL(
	ctx context.Context,
	url string,
) error {

	stocks, err := ParseCSV(ctx, url)
	if err != nil {
		return err
	}

	for _, stock := range stocks {
		if err := r.UpsertStock(ctx, &stock); err != nil {
			return err
		}
	}

	return nil
}

// UPSERT
func (r *Repository) UpsertStock(
	ctx context.Context,
	s *models.Stock,
) error {

	query := `
		INSERT INTO stocks (
			exchange_instrument_id, segment, instrument_type, symbol,
			display_name, company_name, isin, series, exchange,
			contract_expiration, strike, option_type, underlying_symbol_id,
			underlying_symbol, lot_size, tick_size, upper_circuit, lower_circuit,
			freeze_qty, description, ltp, open, high, low, close,
			vol, oi, bid, ask, bid_qty, ask_qty, cautionary_message_info,
			last_updated
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,
			$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,NOW()
		)
		ON CONFLICT (symbol) DO UPDATE SET
			exchange_instrument_id = EXCLUDED.exchange_instrument_id,
			segment = EXCLUDED.segment,
			instrument_type = EXCLUDED.instrument_type,
			display_name = EXCLUDED.display_name,
			company_name = EXCLUDED.company_name,
			isin = EXCLUDED.isin,
			series = EXCLUDED.series,
			exchange = EXCLUDED.exchange,
			contract_expiration = EXCLUDED.contract_expiration,
			strike = EXCLUDED.strike,
			option_type = EXCLUDED.option_type,
			underlying_symbol_id = EXCLUDED.underlying_symbol_id,
			underlying_symbol = EXCLUDED.underlying_symbol,
			lot_size = EXCLUDED.lot_size,
			tick_size = EXCLUDED.tick_size,
			upper_circuit = EXCLUDED.upper_circuit,
			lower_circuit = EXCLUDED.lower_circuit,
			freeze_qty = EXCLUDED.freeze_qty,
			description = EXCLUDED.description,
			ltp = EXCLUDED.ltp,
			open = EXCLUDED.open,
			high = EXCLUDED.high,
			low = EXCLUDED.low,
			close = EXCLUDED.close,
			vol = EXCLUDED.vol,
			oi = EXCLUDED.oi,
			bid = EXCLUDED.bid,
			ask = EXCLUDED.ask,
			bid_qty = EXCLUDED.bid_qty,
			ask_qty = EXCLUDED.ask_qty,
			cautionary_message_info = EXCLUDED.cautionary_message_info,
			last_updated = NOW()
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		s.ExchangeInstrumentID,
		s.Segment,
		s.InstrumentType,
		s.Symbol,
		s.DisplayName,
		s.CompanyName,
		s.ISIN,
		s.Series,
		s.Exchange,
		s.ContractExpiration,
		s.Strike,
		s.OptionType,
		s.UnderlyingSymbolID,
		s.UnderlyingSymbol,
		s.LotSize,
		s.TickSize,
		s.UpperCircuit,
		s.LowerCircuit,
		s.FreezeQty,
		s.Description,
		s.LTP,
		s.Open,
		s.High,
		s.Low,
		s.Close,
		s.Vol,
		s.OI,
		s.Bid,
		s.Ask,
		s.BidQty,
		s.AskQty,
		s.CautionaryMessageInfo,
	)

	return err
}

func (r *Repository) GetStats(
	ctx context.Context,
) (int, int, error) {

	var total int

	err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM stocks`,
	).Scan(&total)

	return total, 0, err
}
