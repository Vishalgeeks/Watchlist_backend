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
	if err == nil {
		return w, nil
	}
	if err.Error() == "wallet not found" {
		if createErr := s.CreateWalletForUser(ctx, userID); createErr == nil {
			return s.repo.GetByUserID(ctx, userID)
		}
	}
	return nil, err
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

func (s *Service) ensureWalletInTx(ctx context.Context, tx *sql.Tx, userID int) (*models.Wallet, error) {
	wallet, err := s.repo.GetByUserIDWithinTx(ctx, tx, userID)
	if err == nil {
		return wallet, nil
	}
	if err.Error() != "wallet not found" {
		return nil, err
	}
	if err := s.repo.CreateWalletWithinTx(ctx, tx, userID, 100000.00); err != nil {
		return nil, err
	}
	return s.repo.GetByUserIDWithinTx(ctx, tx, userID)
}

func (s *Service) Deposit(ctx context.Context, userID int, amount float64, description *string) (*models.WalletTransaction, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	wallet, err := s.ensureWalletInTx(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateBalanceWithinTx(ctx, tx, userID, amount); err != nil {
		return nil, err
	}

	wallet, err = s.repo.GetByUserIDWithinTx(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	t := &models.WalletTransaction{
		WalletID:     wallet.ID,
		UserID:       userID,
		Type:         "DEPOSIT",
		Amount:       amount,
		BalanceAfter: wallet.Balance,
		Description:  description,
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

	wallet, err := s.ensureWalletInTx(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	if wallet.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	if err := s.repo.UpdateBalanceWithinTx(ctx, tx, userID, -amount); err != nil {
		return nil, err
	}

	wallet, err = s.repo.GetByUserIDWithinTx(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	if wallet.Balance < 0 {
		return nil, errors.New("insufficient balance")
	}

	t := &models.WalletTransaction{
		WalletID:     wallet.ID,
		UserID:       userID,
		Type:         "WITHDRAWAL",
		Amount:       amount,
		BalanceAfter: wallet.Balance,
		Description:  description,
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
	wallet, err := s.ensureWalletInTx(ctx, tx, userID)
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
