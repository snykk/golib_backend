package responses

import (
	"time"

	"github.com/snykk/golib_backend/domains/books"
)

type BookResponse struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Publisher   string    `json:"publisher"`
	ISBN        string    `json:"isbn"`
	Rating      *float64  `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func FromDomain(bookDomain books.Domain) BookResponse {
	return BookResponse{
		Id:          bookDomain.ID,
		Title:       bookDomain.Title,
		Description: bookDomain.Description,
		Author:      bookDomain.Author,
		Publisher:   bookDomain.Publisher,
		ISBN:        bookDomain.ISBN,
		Rating:      bookDomain.Rating,
		CreatedAt:   bookDomain.CreatedAt,
		UpdatedAt:   bookDomain.UpdatedAt,
	}
}

func ToResponseList(domains []books.Domain) []BookResponse {
	var result []BookResponse

	for _, val := range domains {
		result = append(result, FromDomain(val))
	}

	return result
}
