package auth

import (
	"context"
	"errors"
	"strings"
	"time"
	"watchlist-backend/pkg/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      *Repository
	jwtSecret string
}

func NewService(repo *Repository, jwtSecret string) *Service {
	return &Service{repo: repo, jwtSecret: jwtSecret}
}

func (s *Service) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Name = strings.TrimSpace(req.Name)

	// 1. Check if user already exists
	_, _, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}
	if err.Error() != "user not found" {
		return nil, err
	}

	// 2. Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. Create user
	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
	}
	if err = s.repo.CreateUser(ctx, user, string(hash)); err != nil {
		return nil, err
	}

	// 4. Generate token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{Token: token, User: *user}, nil
}

func (s *Service) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	user, passwordHash, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{Token: token, User: *user}, nil
}

func (s *Service) generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
