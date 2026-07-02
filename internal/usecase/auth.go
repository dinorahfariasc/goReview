package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"goreview/internal/domain"
	"goreview/internal/util"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenInvalid       = errors.New("invalid or expired refresh token")
)

type AuthRepository interface {
	CreateUser(ctx context.Context, email, passwordHash string) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	CreateRefreshToken(ctx context.Context, token string, userID int64, expiresAt time.Time) error
	GetRefreshTokenUserID(ctx context.Context, token string) (int64, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}

type AuthService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) AuthService {
	return AuthService{repo: repo}
}

func (s AuthService) Register(ctx context.Context, input domain.RegisterInput) (domain.User, error) {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	if input.Email == "" || len(input.Password) < 6 {
		return domain.User{}, ValidationError{Message: "invalid email or password (min 6 chars)"}
	}

	hash, err := util.HashPassword(input.Password)
	if err != nil {
		return domain.User{}, err
	}

	user, err := s.repo.CreateUser(ctx, input.Email, hash)
	if err != nil {
		return domain.User{}, errors.New("failed to register user (email might already exist)")
	}
	return user, nil
}

func (s AuthService) Login(ctx context.Context, input domain.LoginInput) (domain.TokenResponse, error) {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))

	user, err := s.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return domain.TokenResponse{}, ErrInvalidCredentials
	}

	if err := util.CheckPassword(input.Password, user.PasswordHash); err != nil {
		return domain.TokenResponse{}, ErrInvalidCredentials
	}

	return s.generateTokens(ctx, user.ID)
}

func (s AuthService) RefreshToken(ctx context.Context, oldToken string) (domain.TokenResponse, error) {
	userID, err := s.repo.GetRefreshTokenUserID(ctx, oldToken)
	if err != nil {
		return domain.TokenResponse{}, ErrTokenInvalid
	}

	_ = s.repo.DeleteRefreshToken(ctx, oldToken)
	return s.generateTokens(ctx, userID)
}

func (s AuthService) generateTokens(ctx context.Context, userID int64) (domain.TokenResponse, error) {
	accessToken, err := util.GenerateAccessToken(userID)
	if err != nil {
		return domain.TokenResponse{}, err
	}

	refreshToken := util.GenerateRefreshToken()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	err = s.repo.CreateRefreshToken(ctx, refreshToken, userID, expiresAt)
	if err != nil {
		return domain.TokenResponse{}, err
	}

	return domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
