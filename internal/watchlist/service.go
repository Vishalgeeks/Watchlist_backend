package watchlist

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"watchlist-backend/pkg/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateWatchlist(ctx context.Context, userID int, name string) (*models.Watchlist, error) {
	exists, err := s.repo.ExistsByName(ctx, userID, strings.TrimSpace(name))
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("watchlist with this name already exists")
	}
	w := &models.Watchlist{
		UserID: userID,
		Name:   strings.TrimSpace(name),
	}
	if err := s.repo.Create(ctx, w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *Service) GetWatchlists(ctx context.Context, userID int) ([]models.Watchlist, error) {
	return s.repo.GetAllByUserID(ctx, userID)
}

func (s *Service) DeleteWatchlist(ctx context.Context, userID, watchlistID int) error {
	w, err := s.repo.GetByID(ctx, watchlistID)
	if err != nil {
		return err
	}
	if w.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.repo.Delete(ctx, watchlistID)
}

func (s *Service) GetStocks(ctx context.Context, userID, watchlistID int) ([]models.WatchlistItem, error) {
	w, err := s.repo.GetByID(ctx, watchlistID)
	if err != nil {
		return nil, err
	}
	if w.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	return s.repo.GetStocks(ctx, watchlistID)
}

func (s *Service) AddStock(ctx context.Context, userID, watchlistID, stockID int) error {
	w, err := s.repo.GetByID(ctx, watchlistID)
	if err != nil {
		return err
	}
	if w.UserID != userID {
		return errors.New("unauthorized")
	}
	exists, err := s.repo.StockExists(ctx, watchlistID, stockID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("stock already in watchlist")
	}
	return s.repo.AddStock(ctx, watchlistID, stockID)
}

func (s *Service) RemoveStock(ctx context.Context, userID, watchlistID, stockID int) error {
	w, err := s.repo.GetByID(ctx, watchlistID)
	if err != nil {
		return err
	}
	if w.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.repo.RemoveStock(ctx, watchlistID, stockID)
}

// String to Int helper
func ParseInt(s string) (int, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("invalid id format")
	}
	return val, nil
}
