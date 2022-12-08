package books

import (
	"time"

	books "github.com/snykk/golib_backend/domains/books"
	"gorm.io/gorm"
)

type Book struct {
	Id          int     `gorm:"primaryKey;autoIncrement"`
	Title       string  `gorm:"type:varchar(100); not null"`
	Description string  `gorm:"type:text; not null"`
	Author      string  `gorm:"type:varchar(30); not null"`
	Publisher   string  `gorm:"type:varchar(30); not null"`
	ISBN        string  `gorm:"type:char(13); not null"`
	Rating      float64 `gorm:"type:NUMERIC(2,1); not null"`
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
		Rating:      book.Rating,
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
		Rating:      book.Rating,
		CreatedAt:   book.CreatedAt,
		UpdatedAt:   book.UpdatedAt,
	}
}
