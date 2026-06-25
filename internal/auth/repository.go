package auth

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

func (r *Repository) CreateUser(ctx context.Context, user *models.User, passwordHash string) error {
	query := `
		INSERT INTO users (name, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query, user.Name, user.Email, passwordHash).Scan(&user.ID)
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, string, error) {
	query := `
		SELECT id, name, email, password_hash, created_at
		FROM users WHERE email = $1
	`
	user := &models.User{}
	var passwordHash string

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &passwordHash, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, "", errors.New("user not found")
	}
	if err != nil {
		return nil, "", err
	}
	return user, passwordHash, nil
}
