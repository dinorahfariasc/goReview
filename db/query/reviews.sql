-- name: CreateReview :one
INSERT INTO reviews (movie_id, user_id, reviewer_name, rating, content)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, movie_id, user_id, reviewer_name, rating, content, created_at, updated_at;

-- name: GetReview :one
SELECT id, movie_id, user_id, reviewer_name, rating, content, created_at, updated_at
FROM reviews
WHERE id = $1;

-- name: ListReviewsByMovie :many
SELECT id, movie_id, user_id, reviewer_name, rating, content, created_at, updated_at
FROM reviews
WHERE movie_id = $1
ORDER BY id;

-- name: UpdateReview :one
UPDATE reviews
SET reviewer_name = $3,
    rating = $4,
    content = $5,
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING id, movie_id, user_id, reviewer_name, rating, content, created_at, updated_at;

-- name: DeleteReview :exec
DELETE FROM reviews
WHERE id = $1 AND user_id = $2;