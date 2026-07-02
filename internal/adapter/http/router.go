package httpadapter

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func NewRouter(handler Handler, authHandler AuthHandler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// OWASP: Rate Limiting (Máximo de 100 requisições por minuto por IP)
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	// Rotas Públicas
	r.Get("/health", handler.handleHealth)
	r.Post("/register", authHandler.handleRegister)
	r.Post("/login", authHandler.handleLogin)
	r.Post("/refresh", authHandler.handleRefreshToken)

	r.Get("/movies", handler.handleListMovies)
	r.Get("/movies/{id}", handler.handleGetMovie)
	r.Get("/movies/{id}/details", handler.handleGetMovieDetails)
	r.Get("/movies/{id}/reviews", handler.handleListReviewsByMovie)
	r.Get("/reviews/{id}", handler.handleGetReview)

	// Rotas Privadas (Protegidas pelo AuthMiddleware)
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)

		r.Post("/movies", handler.handleCreateMovie)
		r.Put("/movies/{id}", handler.handleUpdateMovie)
		r.Delete("/movies/{id}", handler.handleDeleteMovie)

		r.Post("/movies/{id}/reviews", handler.handleCreateReview)
		r.Put("/reviews/{id}", handler.handleUpdateReview)
		r.Delete("/reviews/{id}", handler.handleDeleteReview)
	})

	return r
}
