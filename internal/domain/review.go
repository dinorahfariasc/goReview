package domain

import "time"

type Review struct {
	ID           int64     `json:"id"`
	MovieID      int64     `json:"movie_id"`
	ReviewerName string    `json:"reviewer_name"`
	Rating       int32     `json:"rating"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateReviewInput struct {
	MovieID      int64  `json:"movie_id"`
	ReviewerName string `json:"reviewer_name"`
	Rating       int32  `json:"rating"`
	Content      string `json:"content"`
}

type UpdateReviewInput struct {
	ReviewerName *string `json:"reviewer_name"`
	Rating       *int32  `json:"rating"`
	Content      *string `json:"content"`
}
