package requests

import (
	"github.com/snykk/golib_backend/domains/books"
)

type BookUpdateRequests struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Publisher   string `json:"publisher"`
	ISBN        string `json:"isbn"`
}

func (b *BookUpdateRequests) ToDomain() *books.Domain {
	return &books.Domain{
		Title:       b.Title,
		Description: b.Description,
		Author:      b.Author,
		Publisher:   b.Publisher,
		ISBN:        b.ISBN,
	}
}
