package usecase

import (
	"context"
	"errors"
	"strings"

	"goreview/internal/domain"
)

type MovieRepository interface {
	ListMovies(ctx context.Context) ([]domain.Movie, error)
	GetMovie(ctx context.Context, id int64) (domain.Movie, error)
	CreateMovie(ctx context.Context, input domain.CreateMovieInput) (domain.Movie, error)
	UpdateMovie(ctx context.Context, movie domain.Movie) (domain.Movie, error)
	DeleteMovie(ctx context.Context, id int64) error
}

type ReviewRepository interface {
	ListReviewsByMovie(ctx context.Context, movieID int64) ([]domain.Review, error)
	GetReview(ctx context.Context, id int64) (domain.Review, error)
	CreateReview(ctx context.Context, input domain.CreateReviewInput) (domain.Review, error)
	UpdateReview(ctx context.Context, review domain.Review) (domain.Review, error)
	DeleteReview(ctx context.Context, id int64) error
}

type Service struct {
	movies  MovieRepository
	reviews ReviewRepository
}

func NewService(movieRepo MovieRepository, reviewRepo ReviewRepository) Service {
	return Service{
		movies:  movieRepo,
		reviews: reviewRepo,
	}
}

func (s Service) ListMovies(ctx context.Context) ([]domain.Movie, error) {
	return s.movies.ListMovies(ctx)
}

func (s Service) GetMovie(ctx context.Context, id int64) (domain.Movie, error) {
	movie, err := s.movies.GetMovie(ctx, id)
	if err != nil {
		return domain.Movie{}, mapMovieError(err)
	}
	return movie, nil
}

func (s Service) CreateMovie(ctx context.Context, input domain.CreateMovieInput) (domain.Movie, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Synopsis = strings.TrimSpace(input.Synopsis)

	if err := validateMovie(input.Title, input.Synopsis, input.ReleaseYear); err != nil {
		return domain.Movie{}, err
	}

	return s.movies.CreateMovie(ctx, input)
}

func (s Service) UpdateMovie(ctx context.Context, id int64, input domain.UpdateMovieInput) (domain.Movie, error) {
	current, err := s.movies.GetMovie(ctx, id)
	if err != nil {
		return domain.Movie{}, mapMovieError(err)
	}

	if input.Title != nil {
		current.Title = strings.TrimSpace(*input.Title)
	}
	if input.Synopsis != nil {
		current.Synopsis = strings.TrimSpace(*input.Synopsis)
	}
	if input.ReleaseYear != nil {
		current.ReleaseYear = *input.ReleaseYear
	}

	if err := validateMovie(current.Title, current.Synopsis, current.ReleaseYear); err != nil {
		return domain.Movie{}, err
	}

	return s.movies.UpdateMovie(ctx, current)
}

func (s Service) DeleteMovie(ctx context.Context, id int64) error {
	if _, err := s.movies.GetMovie(ctx, id); err != nil {
		return mapMovieError(err)
	}
	return s.movies.DeleteMovie(ctx, id)
}

func (s Service) GetMovieDetails(ctx context.Context, id int64) (domain.MovieDetails, error) {
	movie, err := s.movies.GetMovie(ctx, id)
	if err != nil {
		return domain.MovieDetails{}, mapMovieError(err)
	}

	reviews, err := s.reviews.ListReviewsByMovie(ctx, id)
	if err != nil {
		return domain.MovieDetails{}, err
	}

	return domain.MovieDetails{
		Movie:   movie,
		Reviews: reviews,
	}, nil
}

func (s Service) ListReviewsByMovie(ctx context.Context, movieID int64) ([]domain.Review, error) {
	if _, err := s.movies.GetMovie(ctx, movieID); err != nil {
		return nil, mapMovieError(err)
	}
	return s.reviews.ListReviewsByMovie(ctx, movieID)
}

func (s Service) GetReview(ctx context.Context, id int64) (domain.Review, error) {
	review, err := s.reviews.GetReview(ctx, id)
	if err != nil {
		return domain.Review{}, mapReviewError(err)
	}
	return review, nil
}

func (s Service) CreateReview(ctx context.Context, input domain.CreateReviewInput) (domain.Review, error) {
	if _, err := s.movies.GetMovie(ctx, input.MovieID); err != nil {
		return domain.Review{}, mapMovieError(err)
	}

	input.ReviewerName = strings.TrimSpace(input.ReviewerName)
	input.Content = strings.TrimSpace(input.Content)
	if err := validateReview(input.ReviewerName, input.Rating, input.Content); err != nil {
		return domain.Review{}, err
	}

	return s.reviews.CreateReview(ctx, input)
}

func (s Service) UpdateReview(ctx context.Context, id int64, input domain.UpdateReviewInput) (domain.Review, error) {
	current, err := s.reviews.GetReview(ctx, id)
	if err != nil {
		return domain.Review{}, mapReviewError(err)
	}

	if input.ReviewerName != nil {
		current.ReviewerName = strings.TrimSpace(*input.ReviewerName)
	}
	if input.Rating != nil {
		current.Rating = *input.Rating
	}
	if input.Content != nil {
		current.Content = strings.TrimSpace(*input.Content)
	}

	if err := validateReview(current.ReviewerName, current.Rating, current.Content); err != nil {
		return domain.Review{}, err
	}

	return s.reviews.UpdateReview(ctx, current)
}

func (s Service) DeleteReview(ctx context.Context, id int64) error {
	if _, err := s.reviews.GetReview(ctx, id); err != nil {
		return mapReviewError(err)
	}
	return s.reviews.DeleteReview(ctx, id)
}

func validateMovie(title, synopsis string, releaseYear int32) error {
	if title == "" {
		return ValidationError{Message: "title is required"}
	}
	if synopsis == "" {
		return ValidationError{Message: "synopsis is required"}
	}
	if releaseYear < 1888 {
		return ValidationError{Message: "release_year must be 1888 or later"}
	}
	return nil
}

func validateReview(reviewerName string, rating int32, content string) error {
	if reviewerName == "" {
		return ValidationError{Message: "reviewer_name is required"}
	}
	if content == "" {
		return ValidationError{Message: "content is required"}
	}
	if rating < 1 || rating > 5 {
		return ValidationError{Message: "rating must be between 1 and 5"}
	}
	return nil
}

func mapMovieError(err error) error {
	if errors.Is(err, ErrMovieNotFound) {
		return ErrMovieNotFound
	}
	return err
}

func mapReviewError(err error) error {
	if errors.Is(err, ErrReviewNotFound) {
		return ErrReviewNotFound
	}
	return err
}
