package httpadapter

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"goreview/internal/domain"
	"goreview/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type Service interface {
	ListMovies(ctx context.Context) ([]domain.Movie, error)
	GetMovie(ctx context.Context, id int64) (domain.Movie, error)
	CreateMovie(ctx context.Context, input domain.CreateMovieInput) (domain.Movie, error)
	UpdateMovie(ctx context.Context, id int64, input domain.UpdateMovieInput) (domain.Movie, error)
	DeleteMovie(ctx context.Context, id int64) error
	GetMovieDetails(ctx context.Context, id int64) (domain.MovieDetails, error)
	ListReviewsByMovie(ctx context.Context, movieID int64) ([]domain.Review, error)
	GetReview(ctx context.Context, id int64) (domain.Review, error)
	CreateReview(ctx context.Context, input domain.CreateReviewInput) (domain.Review, error)
	UpdateReview(ctx context.Context, id int64, input domain.UpdateReviewInput) (domain.Review, error)
	DeleteReview(ctx context.Context, id int64, userID int64) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h Handler) RegisterRoutes(r chi.Router) {
	r.Get("/health", h.handleHealth)

	r.Route("/movies", func(r chi.Router) {
		r.Get("/", h.handleListMovies)
		r.Post("/", h.handleCreateMovie)
		r.Get("/{id}", h.handleGetMovie)
		r.Put("/{id}", h.handleUpdateMovie)
		r.Delete("/{id}", h.handleDeleteMovie)
		r.Get("/{id}/details", h.handleGetMovieDetails)
		r.Get("/{id}/reviews", h.handleListReviewsByMovie)
		r.Post("/{id}/reviews", h.handleCreateReview)
	})

	r.Route("/reviews", func(r chi.Router) {
		r.Get("/{id}", h.handleGetReview)
		r.Put("/{id}", h.handleUpdateReview)
		r.Delete("/{id}", h.handleDeleteReview)
	})
}

func (h Handler) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h Handler) handleListMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.service.ListMovies(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to list movies"})
		return
	}
	writeJSON(w, http.StatusOK, movies)
}

func (h Handler) handleGetMovie(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	movie, err := h.service.GetMovie(r.Context(), id)
	if err != nil {
		handleUsecaseError(w, err, "failed to fetch movie")
		return
	}
	writeJSON(w, http.StatusOK, movie)
}

func (h Handler) handleCreateMovie(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateMovieInput
	if err := decodeJSON(r, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid body"})
		return
	}

	movie, err := h.service.CreateMovie(r.Context(), input)
	if err != nil {
		handleUsecaseError(w, err, "failed to create movie")
		return
	}
	writeJSON(w, http.StatusCreated, movie)
}

func (h Handler) handleUpdateMovie(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var input domain.UpdateMovieInput
	if err := decodeJSON(r, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid body"})
		return
	}

	movie, err := h.service.UpdateMovie(r.Context(), id, input)
	if err != nil {
		handleUsecaseError(w, err, "failed to update movie")
		return
	}
	writeJSON(w, http.StatusOK, movie)
}

func (h Handler) handleDeleteMovie(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.service.DeleteMovie(r.Context(), id); err != nil {
		handleUsecaseError(w, err, "failed to delete movie")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h Handler) handleGetMovieDetails(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	details, err := h.service.GetMovieDetails(r.Context(), id)
	if err != nil {
		handleUsecaseError(w, err, "failed to fetch movie details")
		return
	}
	writeJSON(w, http.StatusOK, details)
}

func (h Handler) handleListReviewsByMovie(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	reviews, err := h.service.ListReviewsByMovie(r.Context(), id)
	if err != nil {
		handleUsecaseError(w, err, "failed to list reviews")
		return
	}
	writeJSON(w, http.StatusOK, reviews)
}

func (h Handler) handleCreateReview(w http.ResponseWriter, r *http.Request) {
	movieID, ok := parseID(w, r)
	if !ok {
		return
	}

	var input domain.CreateReviewInput
	if err := decodeJSON(r, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid body"})
		return
	}
	input.MovieID = movieID
	input.UserID = getUserID(r)

	review, err := h.service.CreateReview(r.Context(), input)
	if err != nil {
		handleUsecaseError(w, err, "failed to create review")
		return
	}
	writeJSON(w, http.StatusCreated, review)
}

func (h Handler) handleGetReview(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	review, err := h.service.GetReview(r.Context(), id)
	if err != nil {
		handleUsecaseError(w, err, "failed to fetch review")
		return
	}
	writeJSON(w, http.StatusOK, review)
}

func (h Handler) handleUpdateReview(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var input domain.UpdateReviewInput
	if err := decodeJSON(r, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid body"})
		return
	}
	input.UserID = getUserID(r)

	review, err := h.service.UpdateReview(r.Context(), id, input)
	if err != nil {
		handleUsecaseError(w, err, "failed to update review")
		return
	}
	writeJSON(w, http.StatusOK, review)
}

func (h Handler) handleDeleteReview(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.service.DeleteReview(r.Context(), id, getUserID(r)); err != nil {
		handleUsecaseError(w, err, "failed to delete review")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func handleUsecaseError(w http.ResponseWriter, err error, fallback string) {
	switch {
	case errors.Is(err, usecase.ErrMovieNotFound), errors.Is(err, usecase.ErrReviewNotFound):
		writeJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})
	case errors.Is(err, usecase.ErrValidation):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: fallback})
	}
}

func writeJSON(w http.ResponseWriter, code int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(value)
}

func decodeJSON(r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return 0, false
	}
	return id, true
}
