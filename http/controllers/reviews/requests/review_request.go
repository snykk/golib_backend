package requests

import "github.com/snykk/golib_backend/domains/reviews"

type ReviewRequest struct {
	Text   string `json:"text" binding:"required"`
	Rating int    `json:"rating" binding:"required"`
	BookId int    `json:"book_id" binding:"required"`
}

func (r *ReviewRequest) ToDomain() *reviews.Domain {
	return &reviews.Domain{
		Text:   r.Text,
		Rating: r.Rating,
		BookId: r.BookId,
	}
}
