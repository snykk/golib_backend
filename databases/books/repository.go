package books

import (
	"context"

	"github.com/snykk/golib_backend/usecase/books"
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

func (repo *BookRepository) Store(ctx context.Context, b *books.Domain) (books.Domain, error) {
	var result = FromDomain(b)
	if err := repo.Conn.Save(&result).Error; err != nil {
		return books.Domain{}, err
	}

	return result.ToDomain(), nil
}

func (repo *BookRepository) GetAll() ([]books.Domain, error) {
	var booksFromDB []Book
	result := repo.Conn.Find(&booksFromDB)

	if result.Error != nil {
		return []books.Domain{}, result.Error
	}

	var convertedBook []books.Domain

	for _, val := range booksFromDB {
		convertedBook = append(convertedBook, val.ToDomain())
	}

	return convertedBook, nil
}

func (repo *BookRepository) GetById(ctx context.Context, id int) (books.Domain, error) {
	var b Book

	if err := repo.Conn.First(&b, id).Error; err != nil {
		return books.Domain{}, err
	}

	return b.ToDomain(), nil
}

func (repo *BookRepository) Update(ctx context.Context, b *books.Domain) (books.Domain, error) {
	bookFromDB := FromDomain(b)
	if err := repo.Conn.Save(bookFromDB).Error; err != nil {
		return books.Domain{}, err
	}

	return bookFromDB.ToDomain(), nil

}

func (repo *BookRepository) Delete(ctx context.Context, id int) error {
	err := repo.Conn.Delete(&Book{}, id).Error
	if err != nil {
		return err
	}

	return nil
}
