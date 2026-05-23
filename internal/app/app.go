package app

import (
	"context"
	"fmt"
	"net/http"

	db "goreview/db/sqlc"
	httpadapter "goreview/internal/adapter/http"
	"goreview/internal/adapter/postgres"
	"goreview/internal/config"
	"goreview/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	queries := db.New(pool)
	repository := postgres.NewRepository(queries)
	service := usecase.NewService(repository, repository)
	handler := httpadapter.NewHandler(service)

	return http.ListenAndServe(":"+cfg.Port, httpadapter.NewRouter(handler))
}
