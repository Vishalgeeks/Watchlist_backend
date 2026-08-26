package wallet

import (
	"context"
	"database/sql"
	"errors"
	"watchlist-backend/pkg/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetWallet(ctx context.Context, userID int) (*models.Wallet, error) {
	w, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *Service) AdjustBalance(ctx context.Context, userID int, amount float64) error {
	return s.repo.UpdateBalance(ctx, userID, amount)
}

func (s *Service) AdjustBalanceWithinTx(ctx context.Context, tx *sql.Tx, userID int, amount float64) error {
	return s.repo.UpdateBalanceWithinTx(ctx, tx, userID, amount)
}

func (s *Service) CreateWalletForUser(ctx context.Context, userID int) error {
	return s.repo.CreateWallet(ctx, userID, 100000.00)
}

func (s *Service) Deposit(ctx context.Context, userID int, amount float64, description *string) (*models.WalletTransaction, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update balance
	if err := s.repo.UpdateBalanceWithinTx(ctx, tx, userID, amount); err != nil {
		return nil, err
	}

	// Get wallet to find wallet_id and updated balance
	wallet, err := s.repo.GetByUserIDWithinTx(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	// Create transaction record
	t := &models.WalletTransaction{
		WalletID:      wallet.ID,
		UserID:        userID,
		Type:          "DEPOSIT",
		Amount:        amount,
		BalanceAfter:  wallet.Balance,
		Description:   description,
	}
	if err := s.repo.CreateTransactionWithinTx(ctx, tx, t); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) Withdraw(ctx context.Context, userID int, amount float64, description *string) (*models.WalletTransaction, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update balance (negative amount for withdrawal)
	if err := s.repo.UpdateBalanceWithinTx(ctx, tx, userID, -amount); err != nil {
		return nil, err
	}

	// Get wallet to find wallet_id and updated balance
	wallet, err := s.repo.GetByUserIDWithinTx(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	// Check balance is not negative
	if wallet.Balance < 0 {
		return nil, errors.New("insufficient balance")
	}

	// Create transaction record
	t := &models.WalletTransaction{
		WalletID:      wallet.ID,
		UserID:        userID,
		Type:          "WITHDRAWAL",
		Amount:        amount,
		BalanceAfter:  wallet.Balance,
		Description:   description,
	}
	if err := s.repo.CreateTransactionWithinTx(ctx, tx, t); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) RecordTransactionWithinTx(ctx context.Context, tx *sql.Tx, userID int, txType string, amount float64, balanceAfter float64, referenceType *string, referenceID *int, description *string) (*models.WalletTransaction, error) {
	// Get wallet to find wallet_id
	wallet, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	t := &models.WalletTransaction{
		WalletID:      wallet.ID,
		UserID:        userID,
		Type:          txType,
		Amount:        amount,
		BalanceAfter:  balanceAfter,
		ReferenceType: referenceType,
		ReferenceID:   referenceID,
		Description:   description,
	}
	if err := s.repo.CreateTransactionWithinTx(ctx, tx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) GetTransactions(ctx context.Context, userID int, limit, offset int) ([]models.WalletTransaction, error) {
	return s.repo.GetTransactionsByUserID(ctx, userID, limit, offset)
}

func (s *Service) GetTransaction(ctx context.Context, userID, transactionID int) (*models.WalletTransaction, error) {
	return s.repo.GetTransactionByID(ctx, transactionID, userID)
}
