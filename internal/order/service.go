package order

import (
	"context"
	"errors"
	"watchlist-backend/pkg/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateOrder(ctx context.Context, userID int, req *models.CreateOrderRequest) (*models.Order, error) {
	if req.Side != "BUY" && req.Side != "SELL" {
		return nil, errors.New("invalid side: must be BUY or SELL")
	}
	if req.OrderType != "MARKET" && req.OrderType != "LIMIT" {
		return nil, errors.New("invalid order_type: must be MARKET or LIMIT")
	}
	if req.Quantity < 1 {
		return nil, errors.New("quantity must be at least 1")
	}

	if req.OrderType == "LIMIT" {
		if req.Price == nil || *req.Price <= 0 {
			return nil, errors.New("price is required and must be greater than 0 for LIMIT orders")
		}
	}

	exists, err := s.repo.StockExists(ctx, req.StockID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("stock not found")
	}

	order := &models.Order{
		UserID:    userID,
		StockID:   req.StockID,
		Side:      req.Side,
		OrderType: req.OrderType,
		Quantity:  req.Quantity,
		Price:     req.Price,
		Status:    "PENDING",
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, order.ID, userID)
}

func (s *Service) GetOrder(ctx context.Context, userID, orderID int) (*models.Order, error) {
	return s.repo.GetByID(ctx, orderID, userID)
}

func (s *Service) GetAllOrders(ctx context.Context, userID int) ([]models.Order, error) {
	orders, err := s.repo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if orders == nil {
		return []models.Order{}, nil
	}
	return orders, nil
}

func (s *Service) CancelOrder(ctx context.Context, userID, orderID int) error {
	cancelled, err := s.repo.CancelOrder(ctx, orderID, userID)
	if err != nil {
		return err
	}
	if !cancelled {
		return errors.New("order cannot be cancelled")
	}
	return nil
}

func (s *Service) FillOrder(ctx context.Context, orderID, userID int) error {
	return s.repo.UpdateStatus(ctx, orderID, userID, "FILLED")
}