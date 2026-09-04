package trade

import (
	"testing"
	"watchlist-backend/pkg/models"
)

func TestResolveExecutionPriceUsesPositiveFallbacks(t *testing.T) {
	stock := &models.Stock{LTP: 0, Close: 0, Open: 0, Bid: 0, Ask: 12.5}

	got := resolveExecutionPrice(0, "MARKET", nil, stock)
	if got <= 0 {
		t.Fatalf("resolveExecutionPrice() = %v, want a positive fallback price", got)
	}
	if got != 12.5 {
		t.Fatalf("resolveExecutionPrice() = %v, want 12.5 from ask fallback", got)
	}

	limitPrice := 78.25
	got = resolveExecutionPrice(0, "LIMIT", &limitPrice, stock)
	if got != 78.25 {
		t.Fatalf("resolveExecutionPrice() = %v, want 78.25 from limit order price", got)
	}
}
