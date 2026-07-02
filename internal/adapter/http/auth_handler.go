package httpadapter

import (
	"context"
	"net/http"

	"goreview/internal/domain"
)

type AuthService interface {
	Register(ctx context.Context, input domain.RegisterInput) (domain.User, error)
	Login(ctx context.Context, input domain.LoginInput) (domain.TokenResponse, error)
	RefreshToken(ctx context.Context, oldToken string) (domain.TokenResponse, error)
}

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) AuthHandler {
	return AuthHandler{service: service}
}

func (h AuthHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var input domain.RegisterInput
	if err := decodeJSON(r, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid body"})
		return
	}

	user, err := h.service.Register(r.Context(), input)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

func (h AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var input domain.LoginInput
	if err := decodeJSON(r, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid body"})
		return
	}

	tokens, err := h.service.Login(r.Context(), input)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, tokens)
}

func (h AuthHandler) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid body"})
		return
	}

	tokens, err := h.service.RefreshToken(r.Context(), input.RefreshToken)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, tokens)
}
