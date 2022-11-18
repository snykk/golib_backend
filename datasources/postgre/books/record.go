package books

import (
	"time"

	books "github.com/snykk/golib_backend/usecases/books"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Id          int
	Title       string
	Description string
	Author      string
	Publisher   string
	ISBN        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
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
