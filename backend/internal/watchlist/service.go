package watchlist

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"watchlist-backend/cache"
	"watchlist-backend/pkg/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Cache keys — ek jagah define karo
func watchlistsKey(userID int) string {
	return fmt.Sprintf("watchlists:user:%d", userID)
}

func stocksKey(watchlistID int) string {
	return fmt.Sprintf("watchlist:stocks:%d", watchlistID)
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

	// Naya watchlist bana → user ki list purani ho gayi → invalidate
	cache.Rdb.Del(ctx, watchlistsKey(userID))

	return w, nil
}

func (s *Service) GetWatchlists(ctx context.Context, userID int) ([]models.Watchlist, error) {
	key := watchlistsKey(userID)

	// Step 1: Cache check
	val, err := cache.Rdb.Get(ctx, key).Result()
	if err == nil {
		// Cache HIT — JSON decode karke return karo
		var watchlists []models.Watchlist
		if err := json.Unmarshal([]byte(val), &watchlists); err == nil {
			return watchlists, nil
		}
	}

	// Step 2: Cache MISS — DB se lo
	watchlists, err := s.repo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Step 3: Cache mein store karo — 5 min TTL
	if data, err := json.Marshal(watchlists); err == nil {
		cache.Rdb.Set(ctx, key, data, 5*time.Minute)
	}

	return watchlists, nil
}

func (s *Service) DeleteWatchlist(ctx context.Context, userID, watchlistID int) error {
	w, err := s.repo.GetByID(ctx, watchlistID)
	if err != nil {
		return err
	}
	if w.UserID != userID {
		return errors.New("unauthorized")
	}
	if err := s.repo.Delete(ctx, watchlistID); err != nil {
		return err
	}

	// Delete hone ke baad dono cache invalidate karo
	cache.Rdb.Del(ctx, watchlistsKey(userID))
	cache.Rdb.Del(ctx, stocksKey(watchlistID))

	return nil
}

func (s *Service) GetStocks(ctx context.Context, userID, watchlistID int) ([]models.WatchlistItem, error) {
	w, err := s.repo.GetByID(ctx, watchlistID)
	if err != nil {
		return nil, err
	}
	if w.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	key := stocksKey(watchlistID)

	// Step 1: Cache check
	val, err := cache.Rdb.Get(ctx, key).Result()
	if err == nil {
		var items []models.WatchlistItem
		if err := json.Unmarshal([]byte(val), &items); err == nil {
			return items, nil
		}
	}

	// Step 2: Cache MISS — DB se lo
	items, err := s.repo.GetStocks(ctx, watchlistID)
	if err != nil {
		return nil, err
	}

	// Step 3: Cache store — 2 min TTL (stocks change hote hain)
	if data, err := json.Marshal(items); err == nil {
		cache.Rdb.Set(ctx, key, data, 2*time.Minute)
	}

	return items, nil
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
	if err := s.repo.AddStock(ctx, watchlistID, stockID); err != nil {
		return err
	}

	// Stock add hua → stocks cache purana ho gaya
	cache.Rdb.Del(ctx, stocksKey(watchlistID))

	return nil
}

func (s *Service) RemoveStock(ctx context.Context, userID, watchlistID, stockID int) error {
	w, err := s.repo.GetByID(ctx, watchlistID)
	if err != nil {
		return err
	}
	if w.UserID != userID {
		return errors.New("unauthorized")
	}
	if err := s.repo.RemoveStock(ctx, watchlistID, stockID); err != nil {
		return err
	}

	// Stock remove hua → stocks cache purana ho gaya
	cache.Rdb.Del(ctx, stocksKey(watchlistID))

	return nil
}

func ParseInt(s string) (int, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("invalid id format")
	}
	return val, nil
}
