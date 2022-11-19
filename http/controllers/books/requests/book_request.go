package requests

import (
	"github.com/snykk/golib_backend/domains/books"
)

type BookRequest struct {
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author" binding:"required"`
	Description string `json:"description" binding:"required"`
	Publisher   string `json:"publisher" binding:"required"`
	ISBN        string `json:"isbn" binding:"required"`
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
