package books

import (
	"time"

	books "github.com/snykk/golib_backend/usecase/books"
)

type Book struct {
	Id          int    `gorm:"primaryKey"`
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
	Author      string `json:"author" form:"author"`
	Publisher   string `json:"publisher" form:"publisher"`
	ISBN        string `json:"isbn" form:"isbn"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (book *Book) ToDomain() books.Domain {
	return books.Domain{
		ID:          book.Id,
		Title:       book.Title,
		Description: book.Description,
		Author:      book.Author,
		Publisher:   book.Publisher,
		ISBN:        book.ISBN,
		CreatedAt:   book.CreatedAt,
		UpdatedAt:   book.UpdatedAt,
	}
}

func FromDomain(book *books.Domain) Book {
	return Book{
		Id:          book.ID,
		Title:       book.Title,
		Description: book.Description,
		Author:      book.Author,
		Publisher:   book.Publisher,
		ISBN:        book.ISBN,
		CreatedAt:   book.CreatedAt,
		UpdatedAt:   book.UpdatedAt,
	}
}
