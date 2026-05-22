package postgres

import (
	"context"
	"errors"

	db "goreview/db/sqlc"
	"goreview/internal/domain"
	"goreview/internal/usecase"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) Repository {
	return Repository{queries: queries}
}

func (r Repository) ListMovies(ctx context.Context) ([]domain.Movie, error) {
	rows, err := r.queries.ListMovies(ctx)
	if err != nil {
		return nil, err
	}

	movies := make([]domain.Movie, 0, len(rows))
	for _, row := range rows {
		movies = append(movies, toMovie(row))
	}
	return movies, nil
}

func (r Repository) GetMovie(ctx context.Context, id int64) (domain.Movie, error) {
	row, err := r.queries.GetMovie(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Movie{}, usecase.ErrMovieNotFound
	}
	if err != nil {
		return domain.Movie{}, err
	}
	return toMovie(row), nil
}

func (r Repository) CreateMovie(ctx context.Context, input domain.CreateMovieInput) (domain.Movie, error) {
	row, err := r.queries.CreateMovie(ctx, db.CreateMovieParams{
		Title:       input.Title,
		Synopsis:    input.Synopsis,
		ReleaseYear: input.ReleaseYear,
	})
	if err != nil {
		return domain.Movie{}, err
	}
	return toMovie(row), nil
}

func (r Repository) UpdateMovie(ctx context.Context, movie domain.Movie) (domain.Movie, error) {
	row, err := r.queries.UpdateMovie(ctx, db.UpdateMovieParams{
		ID:          movie.ID,
		Title:       movie.Title,
		Synopsis:    movie.Synopsis,
		ReleaseYear: movie.ReleaseYear,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Movie{}, usecase.ErrMovieNotFound
	}
	if err != nil {
		return domain.Movie{}, err
	}
	return toMovie(row), nil
}

func (r Repository) DeleteMovie(ctx context.Context, id int64) error {
	return r.queries.DeleteMovie(ctx, id)
}

func (r Repository) ListReviewsByMovie(ctx context.Context, movieID int64) ([]domain.Review, error) {
	rows, err := r.queries.ListReviewsByMovie(ctx, movieID)
	if err != nil {
		return nil, err
	}

	reviews := make([]domain.Review, 0, len(rows))
	for _, row := range rows {
		reviews = append(reviews, toReview(row))
	}
	return reviews, nil
}

func (r Repository) GetReview(ctx context.Context, id int64) (domain.Review, error) {
	row, err := r.queries.GetReview(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Review{}, usecase.ErrReviewNotFound
	}
	if err != nil {
		return domain.Review{}, err
	}
	return toReview(row), nil
}

func (r Repository) CreateReview(ctx context.Context, input domain.CreateReviewInput) (domain.Review, error) {
	row, err := r.queries.CreateReview(ctx, db.CreateReviewParams{
		MovieID:      input.MovieID,
		ReviewerName: input.ReviewerName,
		Rating:       input.Rating,
		Content:      input.Content,
	})
	if err != nil {
		return domain.Review{}, err
	}
	return toReview(row), nil
}

func (r Repository) UpdateReview(ctx context.Context, review domain.Review) (domain.Review, error) {
	row, err := r.queries.UpdateReview(ctx, db.UpdateReviewParams{
		ID:           review.ID,
		ReviewerName: review.ReviewerName,
		Rating:       review.Rating,
		Content:      review.Content,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Review{}, usecase.ErrReviewNotFound
	}
	if err != nil {
		return domain.Review{}, err
	}
	return toReview(row), nil
}

func (r Repository) DeleteReview(ctx context.Context, id int64) error {
	return r.queries.DeleteReview(ctx, id)
}

func toMovie(movie db.Movie) domain.Movie {
	return domain.Movie{
		ID:          movie.ID,
		Title:       movie.Title,
		Synopsis:    movie.Synopsis,
		ReleaseYear: movie.ReleaseYear,
		CreatedAt:   movie.CreatedAt.Time.UTC(),
		UpdatedAt:   movie.UpdatedAt.Time.UTC(),
	}
}

func toReview(review db.Review) domain.Review {
	return domain.Review{
		ID:           review.ID,
		MovieID:      review.MovieID,
		ReviewerName: review.ReviewerName,
		Rating:       review.Rating,
		Content:      review.Content,
		CreatedAt:    review.CreatedAt.Time.UTC(),
		UpdatedAt:    review.UpdatedAt.Time.UTC(),
	}
}
