package requests

import (
	"github.com/snykk/golib_backend/usecases/books"
)

type BookUpdateRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Publisher   string `json:"publisher"`
	ISBN        string `json:"isbn"`
}

func (b *BookUpdateRequest) ToDomain() *books.Domain {
	return &books.Domain{
		Title:       b.Title,
		Description: b.Description,
		Author:      b.Author,
		Publisher:   b.Publisher,
		ISBN:        b.ISBN,
	}
}
