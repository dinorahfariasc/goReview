package usecase

import "errors"

var (
	ErrMovieNotFound  = errors.New("movie not found")
	ErrReviewNotFound = errors.New("review not found")
	ErrValidation     = errors.New("validation error")
)

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func (e ValidationError) Unwrap() error {
	return ErrValidation
}
