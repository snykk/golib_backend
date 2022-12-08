package responses

import (
	"time"

	"github.com/snykk/golib_backend/domains/reviews"
	bookRes "github.com/snykk/golib_backend/http/controllers/books/responses"
	userRes "github.com/snykk/golib_backend/http/controllers/users/responses"
)

type ReviewResponse struct {
	Id        int    `json:"id"`
	Text      string `json:"text"`
	Rating    int    `json:"rating"`
	BookId    int    `json:"book_id"`
	Book      bookRes.BookResponse
	UserId    int `json:"user_id"`
	User      userRes.UserResponse
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromDomain(domain reviews.Domain) ReviewResponse {
	return ReviewResponse{
		Id:        domain.ID,
		Text:      domain.Text,
		Rating:    domain.Rating,
		BookId:    domain.BookId,
		Book:      bookRes.FromDomain(domain.Book),
		UserId:    domain.UserId,
		User:      userRes.FromDomain(domain.User),
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

func ToResponseList(domains []reviews.Domain) []ReviewResponse {
	var result []ReviewResponse

	for _, val := range domains {
		result = append(result, FromDomain(val))
	}

	return result
}
