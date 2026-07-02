package usecase

import (
	"context"
	"errors"
	"testing"

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
	return []domain.Movie{}, nil
}

func (m *memoryRepository) GetMovie(_ context.Context, id int64) (domain.Movie, error) {
	movie, ok := m.movies[id]
	if !ok {
		return domain.Movie{}, ErrMovieNotFound
	}
	return movie, nil
}

func (m *memoryRepository) CreateMovie(_ context.Context, input domain.CreateMovieInput) (domain.Movie, error) {
	movie := domain.Movie{ID: m.nextMovieID, Title: input.Title, Synopsis: input.Synopsis, ReleaseYear: input.ReleaseYear}
	m.movies[movie.ID] = movie
	m.nextMovieID++
	return movie, nil
}

func (m *memoryRepository) UpdateMovie(context.Context, domain.Movie) (domain.Movie, error) {
	return domain.Movie{}, nil
}

func (m *memoryRepository) DeleteMovie(context.Context, int64) error {
	return nil
}

func (m *memoryRepository) ListReviewsByMovie(context.Context, int64) ([]domain.Review, error) {
	return []domain.Review{}, nil
}

func (m *memoryRepository) GetReview(_ context.Context, id int64) (domain.Review, error) {
	review, ok := m.reviews[id]
	if !ok {
		return domain.Review{}, ErrReviewNotFound
	}
	return review, nil
}

func (m *memoryRepository) CreateReview(_ context.Context, input domain.CreateReviewInput) (domain.Review, error) {
	review := domain.Review{ID: m.nextReviewID, MovieID: input.MovieID, UserID: input.UserID, ReviewerName: input.ReviewerName, Rating: input.Rating, Content: input.Content}
	m.reviews[review.ID] = review
	m.nextReviewID++
	return review, nil
}

func (m *memoryRepository) UpdateReview(context.Context, domain.Review) (domain.Review, error) {
	return domain.Review{}, nil
}

func (m *memoryRepository) DeleteReview(_ context.Context, id int64, userID int64) error {
	return nil
}

// Teste 5
func TestServiceMovieAndReviewFlow(t *testing.T) {
	repo := newMemoryRepository()
	service := NewService(repo, repo)

	movie, _ := service.CreateMovie(context.Background(), domain.CreateMovieInput{Title: "A", Synopsis: "B", ReleaseYear: 2000})
	review, _ := service.CreateReview(context.Background(), domain.CreateReviewInput{MovieID: movie.ID, UserID: 1, ReviewerName: "Me", Rating: 5, Content: "Good"})

	if review.UserID != 1 {
		t.Fatalf("unexpected user ID")
	}
}

// Teste 6
func TestServiceValidation(t *testing.T) {
	repo := newMemoryRepository()
	service := NewService(repo, repo)

	_, err := service.CreateMovie(context.Background(), domain.CreateMovieInput{})
	if err == nil || !errors.Is(err, ErrValidation) {
		t.Fatal("expected validation error")
	}
}

// Teste 7 (OWASP BOLA Protection Test)
func TestServiceReviewBOLA(t *testing.T) {
	repo := newMemoryRepository()
	service := NewService(repo, repo)

	movie, _ := service.CreateMovie(context.Background(), domain.CreateMovieInput{Title: "A", Synopsis: "B", ReleaseYear: 2000})

	review, _ := service.CreateReview(context.Background(), domain.CreateReviewInput{MovieID: movie.ID, UserID: 1, ReviewerName: "Dinorah", Rating: 5, Content: "Top"})

	content := "Hacked"
	_, err := service.UpdateReview(context.Background(), review.ID, domain.UpdateReviewInput{UserID: 2, Content: &content})
	if err == nil {
		t.Fatal("expected error blocking user 2 from editing user 1's review")
	}
}
