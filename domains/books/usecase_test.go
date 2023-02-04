package books_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	bookMocks "github.com/snykk/golib_backend/datasources/databases/books/mocks"
	"github.com/snykk/golib_backend/domains/books"
	"github.com/snykk/golib_backend/http/controllers/books/requests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	bookRepository  *bookMocks.Repository
	bookUsecase     books.Usecase
	booksDataFromDB []books.Domain
	bookDataFromDB  books.Domain
)

func setup(t *testing.T) {
	bookRepository = bookMocks.NewRepository(t)
	bookUsecase = books.NewBookUsecase(bookRepository)
	booksDataFromDB = []books.Domain{
		{
			ID:          1,
			Title:       "Atomic Habits",
			Description: "lorem ipsum doler sit amet",
			Author:      "James Clear",
			Publisher:   "Gramedia",
			ISBN:        "1111111111111",
			Rating:      new(float64),
			CreatedAt:   time.Now(),
		},
		{
			ID:          2,
			Title:       "Selena",
			Description: "lorem ipsum doler sit amet",
			Author:      "Tere Liye",
			Publisher:   "Gramedia",
			ISBN:        "1111111111111",
			Rating:      new(float64),
			CreatedAt:   time.Now(),
		},
	}

	bookDataFromDB = booksDataFromDB[0]
}

func TestStore(t *testing.T) {
	setup(t)
	req := requests.BookRequest{
		Title:       "Atomic Habits",
		Description: "lorem ipsum doler sit amet",
		Author:      "James Clear",
		Publisher:   "Gramedia",
		ISBN:        "1111111111111",
	}
	t.Run("When Success Store Book Data", func(t *testing.T) {
		booksFromDB := books.Domain{
			ID:          1,
			Title:       "Atomic Habits",
			Description: "lorem ipsum doler sit amet",
			Author:      "James Clear",
			Publisher:   "Gramedia",
			ISBN:        "1111111111111",
			CreatedAt:   time.Now(),
		}

		bookRepository.Mock.On("Store", mock.Anything, mock.AnythingOfType("*books.Domain")).Return(booksFromDB, nil).Once()
		result, statusCode, err := bookUsecase.Store(context.Background(), req.ToDomain())

		assert.Nil(t, err)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, http.StatusCreated, statusCode)
		assert.Equal(t, "Atomic Habits", result.Title)
		assert.Equal(t, "lorem ipsum doler sit amet", result.Description)
		assert.Equal(t, "James Clear", result.Author)
		assert.Equal(t, "Gramedia", result.Publisher)
		assert.Equal(t, "1111111111111", result.ISBN)
		assert.NotNil(t, result.CreatedAt)
	})

	t.Run("When Failure", func(t *testing.T) {
		bookRepository.Mock.On("Store", mock.Anything, mock.AnythingOfType("*books.Domain")).Return(books.Domain{}, errors.New("create book failed")).Once()
		result, statusCode, err := bookUsecase.Store(context.Background(), req.ToDomain())

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
		assert.Equal(t, 0, result.ID)
	})

}

func TestGetAll(t *testing.T) {
	setup(t)
	t.Run("When Success Get Books Data", func(t *testing.T) {
		bookRepository.Mock.On("GetAll", mock.Anything).Return(booksDataFromDB, nil).Once()
		result, statusCode, err := bookUsecase.GetAll(context.Background())

		t.Run("Check Book 1", func(t *testing.T) {
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, statusCode)
			assert.Equal(t, 1, result[0].ID)
			assert.Equal(t, booksDataFromDB[0].Title, result[0].Title)
			assert.Equal(t, booksDataFromDB[0].Description, result[0].Description)
			assert.Equal(t, booksDataFromDB[0].Author, result[0].Author)
			assert.Equal(t, booksDataFromDB[0].Publisher, result[0].Publisher)
			assert.Equal(t, booksDataFromDB[0].ISBN, result[0].ISBN)
			assert.NotNil(t, result[0].CreatedAt)
		})

		t.Run("Check Book 2", func(t *testing.T) {
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, statusCode)
			assert.Equal(t, 2, result[1].ID)
			assert.Equal(t, booksDataFromDB[1].Title, result[1].Title)
			assert.Equal(t, booksDataFromDB[1].Description, result[1].Description)
			assert.Equal(t, booksDataFromDB[1].Author, result[1].Author)
			assert.Equal(t, booksDataFromDB[1].Publisher, result[1].Publisher)
			assert.Equal(t, booksDataFromDB[1].ISBN, result[1].ISBN)
			assert.NotNil(t, result[1].CreatedAt)
		})
	})

	t.Run("When Failure Get Books Data", func(t *testing.T) {
		bookRepository.Mock.On("GetAll", mock.Anything).Return([]books.Domain{}, errors.New("get all books failed")).Once()
		result, statusCode, err := bookUsecase.GetAll(context.Background())

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
		assert.Equal(t, []books.Domain{}, result)
	})
}

func TestGetById(t *testing.T) {
	setup(t)
	t.Run("When Success Get Book Data", func(t *testing.T) {
		bookRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(bookDataFromDB, nil).Once()

		result, statusCode, err := bookUsecase.GetById(context.Background(), bookDataFromDB.ID)

		assert.Equal(t, bookDataFromDB, result)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Nil(t, err)
	})

	t.Run("When Failure Book doesn't exist", func(t *testing.T) {
		bookRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(books.Domain{}, errors.New("book not found")).Once()

		result, statusCode, err := bookUsecase.GetById(context.Background(), bookDataFromDB.ID)

		assert.Equal(t, books.Domain{}, result)
		assert.Equal(t, http.StatusNotFound, statusCode)
		assert.Equal(t, errors.New("book not found"), err)
	})
}

func TestDelete(t *testing.T) {
	setup(t)
	t.Run("When Success Delete Book Data", func(t *testing.T) {
		bookRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(bookDataFromDB, nil).Once()
		bookRepository.Mock.On("Delete", mock.Anything, mock.AnythingOfType("int")).Return(nil).Once()

		statusCode, err := bookUsecase.Delete(context.Background(), bookDataFromDB.ID)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, statusCode)
	})
	t.Run("When Failure Delete Book Data", func(t *testing.T) {

		t.Run("User doesn't exist", func(t *testing.T) {
			bookRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(books.Domain{}, errors.New("book not found")).Once()

			statusCode, err := bookUsecase.Delete(context.Background(), 1)

			assert.Equal(t, errors.New("book not found"), err)
			assert.Equal(t, http.StatusNotFound, statusCode)
		})

		t.Run("Failed Delete Book", func(t *testing.T) {
			bookRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(bookDataFromDB, nil).Once()
			bookRepository.Mock.On("Delete", mock.Anything, mock.AnythingOfType("int")).Return(errors.New("failed")).Once()

			statusCode, err := bookUsecase.Delete(context.Background(), 1)

			assert.Equal(t, http.StatusInternalServerError, statusCode)
			assert.Equal(t, errors.New("failed"), err)
		})
	})
}

func TestUpdate(t *testing.T) {
	setup(t)
	t.Run("When Success Update Book", func(t *testing.T) {
		updatedBookFromDB := bookDataFromDB
		updatedBookFromDB.UpdatedAt = time.Now()
		bookRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(updatedBookFromDB, nil).Once()
		bookRepository.Mock.On("Update", mock.Anything, mock.AnythingOfType("*books.Domain")).Return(nil).Once()

		result, statusCode, err := bookUsecase.Update(context.Background(), &bookDataFromDB, bookDataFromDB.ID)

		assert.Equal(t, updatedBookFromDB, result)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Nil(t, err)
		assert.NotNil(t, result.UpdatedAt)
	})
}
