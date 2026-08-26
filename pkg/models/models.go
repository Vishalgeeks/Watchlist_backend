package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Stock struct {
	ID                    int       `json:"id"`
	ExchangeInstrumentID  string    `json:"exchange_instrument_id"`
	Segment               string    `json:"segment"`
	InstrumentType        string    `json:"instrument_type"`
	Symbol                string    `json:"symbol"`
	DisplayName           string    `json:"display_name"`
	CompanyName           string    `json:"company_name"`
	ISIN                  string    `json:"isin"`
	Series                string    `json:"series"`
	Exchange              string    `json:"exchange"`
	ContractExpiration    string    `json:"contract_expiration"`
	Strike                float64   `json:"strike"`
	OptionType            string    `json:"option_type"`
	UnderlyingSymbolID    string    `json:"underlying_symbol_id"`
	UnderlyingSymbol      string    `json:"underlying_symbol"`
	LotSize               int       `json:"lot_size"`
	TickSize              float64   `json:"tick_size"`
	UpperCircuit          float64   `json:"upper_circuit"`
	LowerCircuit          float64   `json:"lower_circuit"`
	FreezeQty             int       `json:"freeze_qty"`
	Description           string    `json:"description"`
	LTP                   float64   `json:"ltp"`
	Open                  float64   `json:"open"`
	High                  float64   `json:"high"`
	Low                   float64   `json:"low"`
	Close                 float64   `json:"close"`
	Vol                   int64     `json:"vol"`
	OI                    int64     `json:"oi"`
	Bid                   float64   `json:"bid"`
	Ask                   float64   `json:"ask"`
	BidQty                int       `json:"bid_qty"`
	AskQty                int       `json:"ask_qty"`
	CautionaryMessageInfo string    `json:"cautionary_message_info"`
	LastUpdated           time.Time `json:"last_updated"`
}

type Watchlist struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	StockCount int       `json:"stock_count"`
}

type WatchlistItem struct {
	ID          int       `json:"id"`
	WatchlistID int       `json:"watchlist_id"`
	StockID     int       `json:"stock_id"`
	AddedAt     time.Time `json:"added_at"`
	Stock       *Stock    `json:"stock,omitempty"`
}

// ── Request DTOs with Validations ─────────────────────
type RegisterRequest struct {
	Name     string `json:"name"     validate:"required,min=2,max=50,no_only_spaces,valid_name"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100,strong_password"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateWatchlistRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100,no_only_spaces,valid_name"`
}

type AddStockRequest struct {
	StockID int `json:"stock_id" validate:"required,min=1"`
}

// ── Order DTOs ─────────────────────────────────────────
type CreateOrderRequest struct {
	StockID   int      `json:"stock_id"    validate:"required,min=1"`
	Side      string   `json:"side"        validate:"required,oneof=BUY SELL"`
	OrderType string   `json:"order_type"  validate:"required,oneof=MARKET LIMIT"`
	Quantity  int      `json:"quantity"    validate:"required,min=1"`
	Price     *float64 `json:"price,omitempty"`
}

type Order struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	StockID   int        `json:"stock_id"`
	Side      string     `json:"side"`
	OrderType string     `json:"order_type"`
	Quantity  int        `json:"quantity"`
	Price     *float64   `json:"price,omitempty"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Stock     *Stock     `json:"stock,omitempty"`
}

// ── Wallet / Portfolio / Trade DTOs ─────────────────────
type Wallet struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Portfolio struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	StockID   int       `json:"stock_id"`
	Quantity  int       `json:"quantity"`
	AvgPrice  float64   `json:"avg_price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Stock     *Stock    `json:"stock,omitempty"`
}

type Trade struct {
	ID              int        `json:"id"`
	OrderID         int        `json:"order_id"`
	UserID          int        `json:"user_id"`
	StockID         int        `json:"stock_id"`
	Side            string     `json:"side"`
	Quantity        int        `json:"quantity"`
	ExecutionPrice  float64    `json:"execution_price"`
	TotalAmount     float64    `json:"total_amount"`
	ExecutedAt      time.Time  `json:"executed_at"`
	Stock           *Stock     `json:"stock,omitempty"`
}

type ExecuteOrderRequest struct {
	OrderID        int     `json:"order_id" validate:"required,min=1"`
	ExecutionPrice float64 `json:"execution_price,omitempty" validate:"omitempty,gt=0"`
}

// ── Wallet Transaction DTOs ─────────────────────────────
type WalletTransaction struct {
	ID            int       `json:"id"`
	WalletID      int       `json:"wallet_id"`
	UserID        int       `json:"user_id"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	BalanceAfter  float64   `json:"balance_after"`
	ReferenceType *string   `json:"reference_type,omitempty"`
	ReferenceID   *int      `json:"reference_id,omitempty"`
	Description   *string   `json:"description,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type DepositRequest struct {
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Description *string `json:"description,omitempty"`
}

type WithdrawRequest struct {
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Description *string `json:"description,omitempty"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// ── Standard Response ─────────────────────────────────
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
