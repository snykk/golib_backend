package books

import (
	"context"

	"github.com/snykk/golib_backend/usecases/books"
	"gorm.io/gorm"
)

type BookRepository struct {
	Conn *gorm.DB
}

func NewBookRepository(conn *gorm.DB) books.Repository {
	return &BookRepository{
		Conn: conn,
	}
}

func (bookRepo *BookRepository) Store(ctx context.Context, b *books.Domain) (books.Domain, error) {
	var result = FromDomain(b)
	if err := bookRepo.Conn.Save(&result).Error; err != nil {
		return books.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (bookRepo *BookRepository) GetAll() ([]books.Domain, error) {
	var booksFromDB []Book
	err := bookRepo.Conn.Find(&booksFromDB).Error

	if err != nil {
		return []books.Domain{}, err
	}

	var convertedBook []books.Domain

	for _, val := range booksFromDB {
		convertedBook = append(convertedBook, val.ToDomain())
	}

	return convertedBook, nil
}

func (bookRepo *BookRepository) GetById(ctx context.Context, id int) (books.Domain, error) {
	var book Book

	if err := bookRepo.Conn.First(&book, id).Error; err != nil {
		return books.Domain{}, err
	}

	return book.ToDomain(), nil
}

func (bookRepo *BookRepository) Update(ctx context.Context, b *books.Domain) (err error) {
	bookFromDB := FromDomain(b)
	err = bookRepo.Conn.Model(&Book{}).Model(&bookFromDB).Updates(&bookFromDB).Error
	return

}

func (bookRepo *BookRepository) Delete(ctx context.Context, id int) (err error) {
	err = bookRepo.Conn.Delete(&Book{}, id).Error
	return
}
