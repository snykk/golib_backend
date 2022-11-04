package requests

import (
	"github.com/snykk/golib_backend/usecases/books"
)

type BookRequest struct {
	Title       string `json:"title" validate:"required"`
	Author      string `json:"author" validate:"required"`
	Description string `json:"description" validate:"required"`
	Publisher   string `json:"publisher" validate:"required"`
	ISBN        string `json:"isbn" validate:"required"`
}

func (bookRequest *BookRequest) ToDomain() *books.Domain {
	return &books.Domain{
		Title:       bookRequest.Title,
		Description: bookRequest.Description,
		Author:      bookRequest.Author,
		Publisher:   bookRequest.Publisher,
		ISBN:        bookRequest.ISBN,
	}
}
