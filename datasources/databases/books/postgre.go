package books

import (
	"context"

	"github.com/snykk/golib_backend/domains/books"
	"gorm.io/gorm"
)

type postgreBookRepository struct {
	conn *gorm.DB
}

func NewPostgreBookRepository(conn *gorm.DB) books.Repository {
	return &postgreBookRepository{
		conn: conn,
	}
}

func (r *postgreBookRepository) Store(ctx context.Context, b *books.Domain) (books.Domain, error) {
	var result = FromDomain(b)
	if err := r.conn.Save(&result).Error; err != nil {
		return books.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (r *postgreBookRepository) GetAll() ([]books.Domain, error) {
	var booksFromDB []Book
	err := r.conn.Find(&booksFromDB).Error

	if err != nil {
		return []books.Domain{}, err
	}

	var convertedBook []books.Domain

	for _, val := range booksFromDB {
		convertedBook = append(convertedBook, val.ToDomain())
	}

	return convertedBook, nil
}

func (r *postgreBookRepository) GetById(ctx context.Context, id int) (books.Domain, error) {
	var book Book

	if err := r.conn.First(&book, id).Error; err != nil {
		return books.Domain{}, err
	}

	return book.ToDomain(), nil
}

func (r *postgreBookRepository) Update(ctx context.Context, b *books.Domain) (err error) {
	bookFromDB := FromDomain(b)
	err = r.conn.Model(&Book{}).Model(&bookFromDB).Updates(&bookFromDB).Error
	return

}

func (r *postgreBookRepository) Delete(ctx context.Context, id int) (err error) {
	err = r.conn.Delete(&Book{}, id).Error
	return
}
