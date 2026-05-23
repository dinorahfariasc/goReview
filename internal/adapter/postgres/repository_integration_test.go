package postgres

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	db "goreview/db/sqlc"
	"goreview/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestRepositoryIntegration(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer pool.Close()

	applyMigrations(t, ctx, pool, "../../../db/migrations")

	repo := NewRepository(db.New(pool))

	movie, err := repo.CreateMovie(ctx, domain.CreateMovieInput{
		Title:       "Arrival",
		Synopsis:    "Sci-fi first contact",
		ReleaseYear: 2016,
	})
	if err != nil {
		t.Fatalf("create movie failed: %v", err)
	}

	review, err := repo.CreateReview(ctx, domain.CreateReviewInput{
		MovieID:      movie.ID,
		ReviewerName: "Dinorah",
		Rating:       5,
		Content:      "Excelente atmosfera",
	})
	if err != nil {
		t.Fatalf("create review failed: %v", err)
	}

	details, err := repo.ListReviewsByMovie(ctx, movie.ID)
	if err != nil {
		t.Fatalf("list reviews failed: %v", err)
	}
	if len(details) != 1 || details[0].ID != review.ID {
		t.Fatalf("unexpected reviews: %+v", details)
	}
}

func applyMigrations(t *testing.T, ctx context.Context, pool *pgxpool.Pool, dir string) {
	t.Helper()

	if _, err := pool.Exec(ctx, `
		DROP TABLE IF EXISTS reviews;
		DROP TABLE IF EXISTS movies;
	`); err != nil {
		t.Fatalf("reset schema failed: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read migrations failed: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		content, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			t.Fatalf("read migration failed: %v", err)
		}

		if _, err := pool.Exec(ctx, string(content)); err != nil {
			t.Fatalf("apply migration failed: %v", err)
		}
	}
}
