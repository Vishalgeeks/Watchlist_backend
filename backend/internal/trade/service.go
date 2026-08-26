package trade

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"
	"watchlist-backend/pkg/models"
	"watchlist-backend/internal/order"
	"watchlist-backend/internal/portfolio"
	"watchlist-backend/internal/stock"
	"watchlist-backend/internal/wallet"
)

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

	executionPrice := execPrice
	if executionPrice <= 0 {
		if o.OrderType == "LIMIT" {
			if o.Price == nil || *o.Price <= 0 {
				return nil, errors.New("invalid execution price for limit order")
			}
			executionPrice = *o.Price
		} else {
stock, err := s.stockRepo.GetByID(ctx, strconv.Itoa(o.StockID))
		if err != nil {
			return nil, err
		}
		executionPrice = stock.LTP
		}
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
		OrderID:         orderID,
		UserID:          userID,
		StockID:         o.StockID,
		Side:            o.Side,
		Quantity:        o.Quantity,
		ExecutionPrice:  executionPrice,
		TotalAmount:     totalAmount,
	}

	if err := s.repo.CreateWithinTx(ctx, tx, trade); err != nil {
		return nil, err
	}

	if err := s.orderRepo.UpdateStatusWithinTx(ctx, tx, orderID, userID, "FILLED"); err != nil {
		return nil, err
	}

delta := -totalAmount
	if o.Side == "SELL" {
		delta = totalAmount
	}
	if err := s.walletSvc.AdjustBalanceWithinTx(ctx, tx, userID, delta); err != nil {
		return nil, err
	}

	// Get updated balance after adjustment
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

	stock, _ := s.stockRepo.GetByID(ctx, strconv.Itoa(o.StockID))
	if stock != nil {
		trade.Stock = &models.Stock{
			ID:         stock.ID,
			Symbol:     stock.Symbol,
			CompanyName: stock.CompanyName,
			Exchange:   stock.Exchange,
			LTP:        stock.LTP,
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
