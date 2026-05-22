package httpadapter

import (
	"bytes"
	"context"
	"encoding/json"
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

func (s *serviceStub) ListMovies(context.Context) ([]domain.Movie, error) {
	return []domain.Movie{s.movies[1]}, nil
}

func (s *serviceStub) GetMovie(context.Context, int64) (domain.Movie, error) {
	return s.movies[1], nil
}

func (s *serviceStub) CreateMovie(context.Context, domain.CreateMovieInput) (domain.Movie, error) {
	return s.movies[1], nil
}

func (s *serviceStub) UpdateMovie(_ context.Context, id int64, input domain.UpdateMovieInput) (domain.Movie, error) {
	movie := s.movies[id]
	if input.Synopsis != nil {
		movie.Synopsis = *input.Synopsis
	}
	s.movies[id] = movie
	return movie, nil
}

func (s *serviceStub) DeleteMovie(context.Context, int64) error {
	return nil
}

func (s *serviceStub) GetMovieDetails(context.Context, int64) (domain.MovieDetails, error) {
	reviews := []domain.Review{}
	for _, review := range s.reviews {
		reviews = append(reviews, review)
	}
	return domain.MovieDetails{Movie: s.movies[1], Reviews: reviews}, nil
}

func (s *serviceStub) ListReviewsByMovie(context.Context, int64) ([]domain.Review, error) {
	reviews := []domain.Review{}
	for _, review := range s.reviews {
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func (s *serviceStub) GetReview(_ context.Context, id int64) (domain.Review, error) {
	return s.reviews[id], nil
}

func (s *serviceStub) CreateReview(_ context.Context, input domain.CreateReviewInput) (domain.Review, error) {
	now := time.Now().UTC()
	review := domain.Review{
		ID:           1,
		MovieID:      input.MovieID,
		ReviewerName: input.ReviewerName,
		Rating:       input.Rating,
		Content:      input.Content,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	s.reviews[review.ID] = review
	return review, nil
}

func (s *serviceStub) UpdateReview(context.Context, int64, domain.UpdateReviewInput) (domain.Review, error) {
	return s.reviews[1], nil
}

func (s *serviceStub) DeleteReview(context.Context, int64) error {
	return nil
}

func TestHandlerMovieDetailsFlow(t *testing.T) {
	handler := NewHandler(newServiceStub())
	router := NewRouter(handler)

	createReviewBody := bytes.NewBufferString(`{"reviewer_name":"Dinorah","rating":5,"content":"Excelente atmosfera"}`)
	createReviewReq := httptest.NewRequest(http.MethodPost, "/movies/1/reviews", createReviewBody)
	createReviewReq.Header.Set("Content-Type", "application/json")
	createReviewRes := httptest.NewRecorder()
	router.ServeHTTP(createReviewRes, createReviewReq)

	if createReviewRes.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", createReviewRes.Code)
	}

	detailsReq := httptest.NewRequest(http.MethodGet, "/movies/1/details", nil)
	detailsRes := httptest.NewRecorder()
	router.ServeHTTP(detailsRes, detailsReq)

	if detailsRes.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", detailsRes.Code)
	}

	var details domain.MovieDetails
	if err := json.NewDecoder(detailsRes.Body).Decode(&details); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if len(details.Reviews) != 1 {
		t.Fatalf("expected 1 review got %d", len(details.Reviews))
	}
}
