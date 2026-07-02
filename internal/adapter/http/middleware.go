package httpadapter

import (
	"context"
	"net/http"
	"strings"

	"goreview/internal/util"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "missing authorization header"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "invalid authorization format"})
			return
		}

		tokenString := parts[1]
		userID, err := util.ValidateToken(tokenString)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "invalid or expired token"})
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserID(r *http.Request) int64 {
	userID, ok := r.Context().Value(UserIDKey).(int64)
	if !ok {
		return 0
	}
	return userID
}
