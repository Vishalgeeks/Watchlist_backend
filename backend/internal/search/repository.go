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
		SELECT id, COALESCE(exchange_instrument_id,''), COALESCE(segment,''), COALESCE(instrument_type,''), symbol,
		       COALESCE(display_name,''), COALESCE(company_name,''), COALESCE(isin,''), COALESCE(series,''),
		       COALESCE(exchange,''), COALESCE(contract_expiration,''), COALESCE(strike,0), COALESCE(option_type,''),
		       COALESCE(underlying_symbol_id,''), COALESCE(underlying_symbol,''), COALESCE(lot_size,0),
		       COALESCE(tick_size,0), COALESCE(upper_circuit,0), COALESCE(lower_circuit,0), COALESCE(freeze_qty,0),
		       COALESCE(description,''), COALESCE(ltp,0), COALESCE(open,0), COALESCE(high,0), COALESCE(low,0),
		       COALESCE(close,0), COALESCE(vol,0), COALESCE(oi,0), COALESCE(bid,0), COALESCE(ask,0),
		       COALESCE(bid_qty,0), COALESCE(ask_qty,0), COALESCE(cautionary_message_info,''), last_updated
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

		result = append(result, s)
	}

	return result, nil
}
