package domain

import "time"

type Movie struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Synopsis    string    `json:"synopsis"`
	ReleaseYear int32     `json:"release_year"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MovieDetails struct {
	Movie
	Reviews []Review `json:"reviews"`
}

type CreateMovieInput struct {
	Title       string `json:"title"`
	Synopsis    string `json:"synopsis"`
	ReleaseYear int32  `json:"release_year"`
}

type UpdateMovieInput struct {
	Title       *string `json:"title"`
	Synopsis    *string `json:"synopsis"`
	ReleaseYear *int32  `json:"release_year"`
}
