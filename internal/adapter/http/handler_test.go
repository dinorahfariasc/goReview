package httpadapter

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"goreview/internal/domain"
)

type serviceStub struct {
	movies  map[int64]domain.Movie
	reviews map[int64]domain.Review
}

func newServiceStub() *serviceStub {
	now := time.Now().UTC()
	return &serviceStub{
		movies: map[int64]domain.Movie{
			1: {ID: 1, Title: "Arrival", Synopsis: "Sci-fi", ReleaseYear: 2016, CreatedAt: now, UpdatedAt: now},
		},
		reviews: map[int64]domain.Review{},
	}
}

func (s *serviceStub) ListMovies(context.Context) ([]domain.Movie, error)    { return nil, nil }
func (s *serviceStub) GetMovie(context.Context, int64) (domain.Movie, error) { return s.movies[1], nil }
func (s *serviceStub) CreateMovie(context.Context, domain.CreateMovieInput) (domain.Movie, error) {
	return s.movies[1], nil
}
func (s *serviceStub) UpdateMovie(_ context.Context, id int64, input domain.UpdateMovieInput) (domain.Movie, error) {
	return s.movies[id], nil
}
func (s *serviceStub) DeleteMovie(context.Context, int64) error { return nil }

func (s *serviceStub) GetMovieDetails(context.Context, int64) (domain.MovieDetails, error) {
	reviews := []domain.Review{}
	for _, review := range s.reviews {
		reviews = append(reviews, review)
	}
	return domain.MovieDetails{Movie: s.movies[1], Reviews: reviews}, nil
}

func (s *serviceStub) ListReviewsByMovie(context.Context, int64) ([]domain.Review, error) {
	return nil, nil
}
func (s *serviceStub) GetReview(_ context.Context, id int64) (domain.Review, error) {
	return s.reviews[id], nil
}

func (s *serviceStub) CreateReview(_ context.Context, input domain.CreateReviewInput) (domain.Review, error) {
	review := domain.Review{ID: 1, MovieID: input.MovieID, UserID: input.UserID, ReviewerName: input.ReviewerName, Rating: input.Rating, Content: input.Content}
	s.reviews[review.ID] = review
	return review, nil
}

func (s *serviceStub) UpdateReview(context.Context, int64, domain.UpdateReviewInput) (domain.Review, error) {
	return s.reviews[1], nil
}
func (s *serviceStub) DeleteReview(context.Context, int64, int64) error { return nil }

type authServiceStub struct{}

func (a authServiceStub) Register(context.Context, domain.RegisterInput) (domain.User, error) {
	return domain.User{}, nil
}
func (a authServiceStub) Login(context.Context, domain.LoginInput) (domain.TokenResponse, error) {
	return domain.TokenResponse{}, nil
}
func (a authServiceStub) RefreshToken(context.Context, string) (domain.TokenResponse, error) {
	return domain.TokenResponse{}, nil
}

// Teste 8
func TestHandlerMovieDetailsFlow(t *testing.T) {
	handler := NewHandler(newServiceStub())
	router := NewRouter(handler, NewAuthHandler(authServiceStub{}))

	// Rota GET é pública, deve funcionar sem token
	detailsReq := httptest.NewRequest(http.MethodGet, "/movies/1/details", nil)
	detailsRes := httptest.NewRecorder()
	router.ServeHTTP(detailsRes, detailsReq)

	if detailsRes.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", detailsRes.Code)
	}
}

// Teste 9
func TestHandlerMiddlewareNoToken(t *testing.T) {
	handler := NewHandler(newServiceStub())
	router := NewRouter(handler, NewAuthHandler(authServiceStub{}))

	// Rota POST é privada, deve ser bloqueada
	req := httptest.NewRequest(http.MethodPost, "/movies", bytes.NewBufferString(`{}`))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized got %d", res.Code)
	}
}

// Teste 10
func TestHandlerAuthRoute(t *testing.T) {
	authHandler := NewAuthHandler(authServiceStub{})
	router := NewRouter(NewHandler(newServiceStub()), authHandler)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{"email":"a@b.com", "password":"123"}`))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200 OK got %d", res.Code)
	}
}
