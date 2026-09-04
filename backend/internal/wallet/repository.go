package wallet

import (
	"context"
	"database/sql"
	"errors"
	"watchlist-backend/pkg/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *Repository) GetByUserID(ctx context.Context, userID int) (*models.Wallet, error) {
	query := `
		SELECT id, user_id, balance, created_at, updated_at
		FROM wallets
		WHERE user_id = $1
	`
	var w models.Wallet
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&w.ID, &w.UserID, &w.Balance, &w.CreatedAt, &w.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("wallet not found")
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *Repository) GetByUserIDWithinTx(ctx context.Context, tx *sql.Tx, userID int) (*models.Wallet, error) {
	query := `
		SELECT id, user_id, balance, created_at, updated_at
		FROM wallets
		WHERE user_id = $1
	`
	var w models.Wallet
	err := tx.QueryRowContext(ctx, query, userID).Scan(
		&w.ID, &w.UserID, &w.Balance, &w.CreatedAt, &w.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("wallet not found")
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *Repository) UpdateBalance(ctx context.Context, userID int, delta float64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE wallets SET balance = balance + $2, updated_at = NOW() WHERE user_id = $1`,
		userID, delta,
	)
	return err
}

func (r *Repository) UpdateBalanceWithinTx(ctx context.Context, tx *sql.Tx, userID int, delta float64) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE wallets SET balance = balance + $2, updated_at = NOW() WHERE user_id = $1`,
		userID, delta,
	)
	return err
}

func (r *Repository) CreateWallet(ctx context.Context, userID int, initialBalance float64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO wallets (user_id, balance) VALUES ($1, $2)
		ON CONFLICT (user_id) DO NOTHING`,
		userID, initialBalance,
	)
	return err
}

func (r *Repository) CreateWalletWithinTx(ctx context.Context, tx *sql.Tx, userID int, initialBalance float64) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO wallets (user_id, balance) VALUES ($1, $2)
		ON CONFLICT (user_id) DO NOTHING`,
		userID, initialBalance,
	)
	return err
}

func (r *Repository) CreateTransaction(ctx context.Context, t *models.WalletTransaction) error {
	query := `
		INSERT INTO wallet_transactions (wallet_id, user_id, type, amount, balance_after, reference_type, reference_id, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(ctx, query,
		t.WalletID, t.UserID, t.Type, t.Amount, t.BalanceAfter,
		t.ReferenceType, t.ReferenceID, t.Description,
	).Scan(&t.ID, &t.CreatedAt)
}

func (r *Repository) CreateTransactionWithinTx(ctx context.Context, sqlTx *sql.Tx, t *models.WalletTransaction) error {
	query := `
		INSERT INTO wallet_transactions (wallet_id, user_id, type, amount, balance_after, reference_type, reference_id, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		RETURNING id, created_at
	`
	return sqlTx.QueryRowContext(ctx, query,
		t.WalletID, t.UserID, t.Type, t.Amount, t.BalanceAfter,
		t.ReferenceType, t.ReferenceID, t.Description,
	).Scan(&t.ID, &t.CreatedAt)
}

func (r *Repository) GetTransactionsByUserID(ctx context.Context, userID int, limit, offset int) ([]models.WalletTransaction, error) {
	query := `
		SELECT id, wallet_id, user_id, type, amount, balance_after, reference_type, reference_id, description, created_at
		FROM wallet_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []models.WalletTransaction
	for rows.Next() {
		var t models.WalletTransaction
		err := rows.Scan(&t.ID, &t.WalletID, &t.UserID, &t.Type, &t.Amount, &t.BalanceAfter, &t.ReferenceType, &t.ReferenceID, &t.Description, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		txs = append(txs, t)
	}
	return txs, nil
}

func (r *Repository) GetTransactionByID(ctx context.Context, transactionID, userID int) (*models.WalletTransaction, error) {
	query := `
		SELECT id, wallet_id, user_id, type, amount, balance_after, reference_type, reference_id, description, created_at
		FROM wallet_transactions
		WHERE id = $1 AND user_id = $2
	`
	var t models.WalletTransaction
	err := r.db.QueryRowContext(ctx, query, transactionID, userID).Scan(
		&t.ID, &t.WalletID, &t.UserID, &t.Type, &t.Amount, &t.BalanceAfter, &t.ReferenceType, &t.ReferenceID, &t.Description, &t.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("transaction not found")
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}
