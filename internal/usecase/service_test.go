package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"goreview/internal/domain"
)

type memoryRepository struct {
	movies       map[int64]domain.Movie
	reviews      map[int64]domain.Review
	nextMovieID  int64
	nextReviewID int64
}

func newMemoryRepository() *memoryRepository {
	return &memoryRepository{
		movies:       make(map[int64]domain.Movie),
		reviews:      make(map[int64]domain.Review),
		nextMovieID:  1,
		nextReviewID: 1,
	}
}

func (m *memoryRepository) ListMovies(context.Context) ([]domain.Movie, error) {
	result := make([]domain.Movie, 0, len(m.movies))
	for i := int64(1); i < m.nextMovieID; i++ {
		movie, ok := m.movies[i]
		if ok {
			result = append(result, movie)
		}
	}
	return result, nil
}

func (m *memoryRepository) GetMovie(_ context.Context, id int64) (domain.Movie, error) {
	movie, ok := m.movies[id]
	if !ok {
		return domain.Movie{}, ErrMovieNotFound
	}
	return movie, nil
}

func (m *memoryRepository) CreateMovie(_ context.Context, input domain.CreateMovieInput) (domain.Movie, error) {
	now := time.Now().UTC()
	movie := domain.Movie{
		ID:          m.nextMovieID,
		Title:       input.Title,
		Synopsis:    input.Synopsis,
		ReleaseYear: input.ReleaseYear,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	m.movies[movie.ID] = movie
	m.nextMovieID++
	return movie, nil
}

func (m *memoryRepository) UpdateMovie(_ context.Context, movie domain.Movie) (domain.Movie, error) {
	if _, ok := m.movies[movie.ID]; !ok {
		return domain.Movie{}, ErrMovieNotFound
	}
	movie.UpdatedAt = time.Now().UTC()
	m.movies[movie.ID] = movie
	return movie, nil
}

func (m *memoryRepository) DeleteMovie(_ context.Context, id int64) error {
	delete(m.movies, id)
	for reviewID, review := range m.reviews {
		if review.MovieID == id {
			delete(m.reviews, reviewID)
		}
	}
	return nil
}

func (m *memoryRepository) ListReviewsByMovie(_ context.Context, movieID int64) ([]domain.Review, error) {
	result := []domain.Review{}
	for i := int64(1); i < m.nextReviewID; i++ {
		review, ok := m.reviews[i]
		if ok && review.MovieID == movieID {
			result = append(result, review)
		}
	}
	return result, nil
}

func (m *memoryRepository) GetReview(_ context.Context, id int64) (domain.Review, error) {
	review, ok := m.reviews[id]
	if !ok {
		return domain.Review{}, ErrReviewNotFound
	}
	return review, nil
}

func (m *memoryRepository) CreateReview(_ context.Context, input domain.CreateReviewInput) (domain.Review, error) {
	now := time.Now().UTC()
	review := domain.Review{
		ID:           m.nextReviewID,
		MovieID:      input.MovieID,
		ReviewerName: input.ReviewerName,
		Rating:       input.Rating,
		Content:      input.Content,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	m.reviews[review.ID] = review
	m.nextReviewID++
	return review, nil
}

func (m *memoryRepository) UpdateReview(_ context.Context, review domain.Review) (domain.Review, error) {
	if _, ok := m.reviews[review.ID]; !ok {
		return domain.Review{}, ErrReviewNotFound
	}
	review.UpdatedAt = time.Now().UTC()
	m.reviews[review.ID] = review
	return review, nil
}

func (m *memoryRepository) DeleteReview(_ context.Context, id int64) error {
	delete(m.reviews, id)
	return nil
}

func TestServiceMovieAndReviewFlow(t *testing.T) {
	repo := newMemoryRepository()
	service := NewService(repo, repo)

	movie, err := service.CreateMovie(context.Background(), domain.CreateMovieInput{
		Title:       " Arrival ",
		Synopsis:    " Sci-fi first contact ",
		ReleaseYear: 2016,
	})
	if err != nil {
		t.Fatalf("create movie failed: %v", err)
	}
	if movie.Title != "Arrival" {
		t.Fatalf("unexpected trimmed title: %q", movie.Title)
	}

	review, err := service.CreateReview(context.Background(), domain.CreateReviewInput{
		MovieID:      movie.ID,
		ReviewerName: " Dinorah ",
		Rating:       5,
		Content:      " Excelente atmosfera ",
	})
	if err != nil {
		t.Fatalf("create review failed: %v", err)
	}
	if review.ReviewerName != "Dinorah" {
		t.Fatalf("unexpected trimmed reviewer: %q", review.ReviewerName)
	}

	details, err := service.GetMovieDetails(context.Background(), movie.ID)
	if err != nil {
		t.Fatalf("details failed: %v", err)
	}
	if len(details.Reviews) != 1 {
		t.Fatalf("expected 1 review got %d", len(details.Reviews))
	}
}

func TestServiceValidation(t *testing.T) {
	repo := newMemoryRepository()
	service := NewService(repo, repo)

	_, err := service.CreateMovie(context.Background(), domain.CreateMovieInput{})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error got %v", err)
	}
}
