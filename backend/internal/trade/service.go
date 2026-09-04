package trade

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"
	"watchlist-backend/internal/order"
	"watchlist-backend/internal/portfolio"
	"watchlist-backend/internal/stock"
	"watchlist-backend/internal/wallet"
	"watchlist-backend/pkg/models"
)

func resolveExecutionPrice(provided float64, orderType string, limitPrice *float64, stock *models.Stock) float64 {
	if provided > 0 {
		return provided
	}
	if orderType == "LIMIT" {
		if limitPrice != nil && *limitPrice > 0 {
			return *limitPrice
		}
		return 0
	}
	if stock == nil {
		return 100
	}
	for _, price := range []float64{stock.LTP, stock.Close, stock.Open, stock.Bid, stock.Ask} {
		if price > 0 {
			return price
		}
	}
	return 100
}

type Service struct {
	repo         *Repository
	db           *sql.DB
	orderRepo    *order.Repository
	walletSvc    *wallet.Service
	portfolioSvc *portfolio.Service
	stockRepo    *stock.Repository
}

func NewService(repo *Repository, db *sql.DB, orderRepo *order.Repository, walletSvc *wallet.Service, portfolioSvc *portfolio.Service, stockRepo *stock.Repository) *Service {
	return &Service{
		repo:         repo,
		db:           db,
		orderRepo:    orderRepo,
		walletSvc:    walletSvc,
		portfolioSvc: portfolioSvc,
		stockRepo:    stockRepo,
	}
}

func (s *Service) ExecuteOrder(ctx context.Context, userID, orderID int, execPrice float64) (*models.Trade, error) {
	if _, err := s.walletSvc.GetWallet(ctx, userID); err != nil {
		return nil, err
	}

	o, err := s.orderRepo.GetByID(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}

	if o.Status != "PENDING" {
		if o.Status == "FILLED" {
			return nil, errors.New("order already executed")
		}
		return nil, errors.New("order not executable")
	}

	stock, err := s.stockRepo.GetByID(ctx, strconv.Itoa(o.StockID))
	if err != nil {
		return nil, err
	}

	executionPrice := resolveExecutionPrice(execPrice, o.OrderType, o.Price, stock)
	if executionPrice <= 0 {
		return nil, errors.New("invalid execution price")
	}

	totalAmount := float64(o.Quantity) * executionPrice

	if o.Side == "BUY" {
		w, err := s.walletSvc.GetWallet(ctx, userID)
		if err != nil {
			return nil, err
		}
		if w.Balance < totalAmount {
			return nil, errors.New("insufficient balance")
		}
	} else {
		holding, err := s.portfolioSvc.GetHoldings(ctx, userID)
		if err != nil {
			return nil, err
		}
		hasEnough := false
		for _, h := range holding {
			if h.StockID == o.StockID && h.Quantity >= o.Quantity {
				hasEnough = true
				break
			}
		}
		if !hasEnough {
			return nil, errors.New("insufficient holdings")
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	trade := &models.Trade{
		OrderID:        orderID,
		UserID:         userID,
		StockID:        o.StockID,
		Side:           o.Side,
		Quantity:       o.Quantity,
		ExecutionPrice: executionPrice,
		TotalAmount:    totalAmount,
	}

	if err := s.repo.CreateWithinTx(ctx, tx, trade); err != nil {
		return nil, err
	}

	if err := s.orderRepo.UpdateStatusWithinTx(ctx, tx, orderID, userID, "FILLED"); err != nil {
		return nil, err
	}

	if _, err := s.walletSvc.GetWallet(ctx, userID); err != nil {
		return nil, err
	}

	delta := -totalAmount
	if o.Side == "SELL" {
		delta = totalAmount
	}
	if err := s.walletSvc.AdjustBalanceWithinTx(ctx, tx, userID, delta); err != nil {
		return nil, err
	}

	wallet, err := s.walletSvc.GetWallet(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Record wallet transaction
	txType := "BUY_DEBIT"
	if o.Side == "SELL" {
		txType = "SELL_CREDIT"
	}
	refType := "trade"
	desc := "order executed"
	_, err = s.walletSvc.RecordTransactionWithinTx(ctx, tx, userID, txType, totalAmount, wallet.Balance, &refType, &trade.ID, &desc)
	if err != nil {
		return nil, err
	}

	qtyDelta := o.Quantity
	if o.Side == "SELL" {
		qtyDelta = -o.Quantity
	}
	if err := s.portfolioSvc.AdjustHoldingWithinTx(ctx, tx, userID, o.StockID, qtyDelta, executionPrice); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Invalidate portfolio summary cache since holdings have changed
	s.portfolioSvc.InvalidateSummary(ctx, userID)

	updatedStock, _ := s.stockRepo.GetByID(ctx, strconv.Itoa(o.StockID))
	if updatedStock != nil {
		trade.Stock = &models.Stock{
			ID:          updatedStock.ID,
			Symbol:      updatedStock.Symbol,
			CompanyName: updatedStock.CompanyName,
			Exchange:    updatedStock.Exchange,
			LTP:         updatedStock.LTP,
		}
	}

	trade.ExecutedAt = time.Now()
	return trade, nil
}

func (s *Service) GetTrade(ctx context.Context, userID, tradeID int) (*models.Trade, error) {
	return s.repo.GetByID(ctx, tradeID, userID)
}

func (s *Service) GetAllTrades(ctx context.Context, userID int) ([]models.Trade, error) {
	trades, err := s.repo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if trades == nil {
		return []models.Trade{}, nil
	}
	return trades, nil
}
