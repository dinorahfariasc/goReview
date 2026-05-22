-- name: CreateMovie :one
INSERT INTO movies (title, synopsis, release_year)
VALUES ($1, $2, $3)
RETURNING id, title, synopsis, release_year, created_at, updated_at;

-- name: GetMovie :one
SELECT id, title, synopsis, release_year, created_at, updated_at
FROM movies
WHERE id = $1;

-- name: ListMovies :many
SELECT id, title, synopsis, release_year, created_at, updated_at
FROM movies
ORDER BY id;

-- name: UpdateMovie :one
UPDATE movies
SET title = $2,
    synopsis = $3,
    release_year = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING id, title, synopsis, release_year, created_at, updated_at;

-- name: DeleteMovie :exec
DELETE FROM movies
WHERE id = $1;
