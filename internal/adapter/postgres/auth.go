package postgres

import (
	"context"
	"errors"
	"time"

	db "goreview/db/sqlc"
	"goreview/internal/domain"
	"goreview/internal/usecase"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r Repository) CreateUser(ctx context.Context, email, passwordHash string) (domain.User, error) {
	row, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:           row.ID,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt.Time.UTC(),
		UpdatedAt:    row.UpdatedAt.Time.UTC(),
	}, nil
}

func (r Repository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, usecase.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:           row.ID,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt.Time.UTC(),
		UpdatedAt:    row.UpdatedAt.Time.UTC(),
	}, nil
}

func (r Repository) CreateRefreshToken(ctx context.Context, token string, userID int64, expiresAt time.Time) error {
	_, err := r.queries.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		Token:     token,
		UserID:    userID,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	return err
}

func (r Repository) GetRefreshTokenUserID(ctx context.Context, token string) (int64, error) {
	row, err := r.queries.GetRefreshToken(ctx, token)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, errors.New("refresh token not found")
	}
	if err != nil {
		return 0, err
	}
	if row.ExpiresAt.Time.Before(time.Now()) {
		return 0, errors.New("refresh token expired")
	}
	return row.UserID, nil
}

func (r Repository) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.queries.DeleteRefreshToken(ctx, token)
}
