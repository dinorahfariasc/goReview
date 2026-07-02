package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"goreview/internal/domain"
)

type mockAuthRepo struct {
	users  map[string]domain.User
	tokens map[string]int64
	nextID int64
}

func newMockAuthRepo() *mockAuthRepo {
	return &mockAuthRepo{
		users:  make(map[string]domain.User),
		tokens: make(map[string]int64),
		nextID: 1,
	}
}

func (m *mockAuthRepo) CreateUser(ctx context.Context, email, hash string) (domain.User, error) {
	if _, exists := m.users[email]; exists {
		return domain.User{}, errors.New("conflict")
	}
	u := domain.User{ID: m.nextID, Email: email, PasswordHash: hash}
	m.users[email] = u
	m.nextID++
	return u, nil
}

func (m *mockAuthRepo) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	u, ok := m.users[email]
	if !ok {
		return domain.User{}, ErrUserNotFound
	}
	return u, nil
}

func (m *mockAuthRepo) CreateRefreshToken(ctx context.Context, token string, userID int64, expiresAt time.Time) error {
	m.tokens[token] = userID
	return nil
}

func (m *mockAuthRepo) GetRefreshTokenUserID(ctx context.Context, token string) (int64, error) {
	id, ok := m.tokens[token]
	if !ok {
		return 0, errors.New("not found")
	}
	return id, nil
}

func (m *mockAuthRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	delete(m.tokens, token)
	return nil
}

// Teste 1
func TestAuthRegisterSuccess(t *testing.T) {
	s := NewAuthService(newMockAuthRepo())
	_, err := s.Register(context.Background(), domain.RegisterInput{Email: "test@test.com", Password: "password123"})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

// Teste 2
func TestAuthRegisterInvalid(t *testing.T) {
	s := NewAuthService(newMockAuthRepo())
	_, err := s.Register(context.Background(), domain.RegisterInput{Email: "", Password: "123"})
	if err == nil {
		t.Fatal("expected error for short password/empty email")
	}
}

// Teste 3
func TestAuthLoginSuccess(t *testing.T) {
	s := NewAuthService(newMockAuthRepo())
	s.Register(context.Background(), domain.RegisterInput{Email: "test@test.com", Password: "password123"})

	tokens, err := s.Login(context.Background(), domain.LoginInput{Email: "test@test.com", Password: "password123"})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if tokens.AccessToken == "" {
		t.Fatal("missing access token")
	}
}

// Teste 4
func TestAuthLoginInvalidCreds(t *testing.T) {
	s := NewAuthService(newMockAuthRepo())
	s.Register(context.Background(), domain.RegisterInput{Email: "test@test.com", Password: "password123"})

	_, err := s.Login(context.Background(), domain.LoginInput{Email: "test@test.com", Password: "wrong"})
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}
