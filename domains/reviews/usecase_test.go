package reviews_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	reviewMocks "github.com/snykk/golib_backend/datasources/databases/reviews/mocks"
	"github.com/snykk/golib_backend/domains/books"
	"github.com/snykk/golib_backend/domains/reviews"
	"github.com/snykk/golib_backend/domains/users"
	"github.com/snykk/golib_backend/http/controllers/reviews/requests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	reviewRepository  *reviewMocks.Repository
	reviewUsecase     reviews.Usecase
	reviewsDataFromDB []reviews.Domain
	reviewDataFromDB  reviews.Domain
	bookFromDB        books.Domain
	userFromDB        users.Domain
)

func setup(t *testing.T) {
	reviewRepository = reviewMocks.NewRepository(t)
	reviewUsecase = reviews.NewReviewUsecase(reviewRepository)
	bookFromDB = books.Domain{
		ID:          1,
		Title:       "Atomic Habits",
		Description: "lorem ipsum doler sit amet",
		Author:      "James Clear",
		Publisher:   "Gramedia",
		ISBN:        "1111111111111",
		Rating:      new(float64),
		CreatedAt:   time.Now(),
	}
	userFromDB = users.Domain{
		ID:          1,
		FullName:    "patrick star",
		Username:    "itsmepatrick",
		Email:       "najibfikri13@gmail.com",
		Password:    "11111",
		Role:        "user",
		Gender:      "male",
		Reviews:     0,
		IsActivated: true,
	}

	reviewsDataFromDB = []reviews.Domain{
		{
			ID:        1,
			Text:      "keren bet yagesya bintang 10",
			Rating:    10,
			BookId:    bookFromDB.ID,
			Book:      bookFromDB,
			UserId:    userFromDB.ID,
			User:      userFromDB,
			CreatedAt: time.Now(),
		},
	}

	reviewDataFromDB = reviewsDataFromDB[0]
}

func TestStore(t *testing.T) {
	setup(t)
	req := requests.ReviewRequest{
		Text:   "keren bet yagesya bintang 10",
		Rating: 10,
		BookId: 1,
	}
	t.Run("When Success Store Review Data", func(t *testing.T) {
		reviewRepository.Mock.On("Store", mock.Anything, mock.AnythingOfType("*reviews.Domain")).Return(reviewDataFromDB, nil).Once()
		result, statusCode, err := reviewUsecase.Store(context.Background(), req.ToDomain(), userFromDB.ID)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, statusCode)
		assert.Equal(t, 1, result.ID)
		assert.NotEqual(t, books.Domain{}, result.Book)
		assert.NotEqual(t, users.Domain{}, result.User)
		assert.NotNil(t, result.CreatedAt)
	})

	t.Run("When Failure Store Review Data", func(t *testing.T) {
		reviewRepository.Mock.On("Store", mock.Anything, mock.AnythingOfType("*reviews.Domain")).Return(reviews.Domain{}, errors.New("create review failed")).Once()
		result, statusCode, err := reviewUsecase.Store(context.Background(), req.ToDomain(), userFromDB.ID)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
		assert.Equal(t, 0, result.ID)
	})

}

func TestGetAll(t *testing.T) {
	setup(t)
	t.Run("When Success Get reviews Data", func(t *testing.T) {
		reviewRepository.Mock.On("GetAll", mock.Anything).Return(reviewsDataFromDB, nil).Once()
		result, statusCode, err := reviewUsecase.GetAll(context.Background())

		assert.Nil(t, err)
		assert.NotEqual(t, 0, result[0].ID)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, reviewsDataFromDB, result)
	})

	t.Run("When Failure Get reviews Data", func(t *testing.T) {
		reviewRepository.Mock.On("GetAll", mock.Anything).Return([]reviews.Domain{}, errors.New("get all reviews failed")).Once()
		result, statusCode, err := reviewUsecase.GetAll(context.Background())

		assert.NotNil(t, err)
		assert.Equal(t, []reviews.Domain{}, result)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}

func TestGetById(t *testing.T) {
	setup(t)
	t.Run("When Success Get review Data", func(t *testing.T) {
		reviewRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(reviewDataFromDB, nil).Once()

		result, statusCode, err := reviewUsecase.GetById(context.Background(), reviewDataFromDB.ID)

		assert.Equal(t, reviewDataFromDB, result)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Nil(t, err)
	})

	t.Run("When Failure Review doesn't exist", func(t *testing.T) {
		reviewRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(reviews.Domain{}, errors.New("review not found")).Once()

		result, statusCode, err := reviewUsecase.GetById(context.Background(), reviewDataFromDB.ID)

		assert.Equal(t, reviews.Domain{}, result)
		assert.Equal(t, http.StatusNotFound, statusCode)
		assert.Equal(t, errors.New("review not found"), err)
	})
}

func TestDelete(t *testing.T) {
	setup(t)
	t.Run("When Success Delete review Data", func(t *testing.T) {
		reviewRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(reviewDataFromDB, nil).Once()
		reviewRepository.Mock.On("Delete", mock.Anything, mock.AnythingOfType("*reviews.Domain")).Return(reviewDataFromDB.BookId, nil).Once()

		bookId, statusCode, err := reviewUsecase.Delete(context.Background(), reviewDataFromDB.UserId, reviewDataFromDB.ID)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, statusCode)
		assert.Equal(t, bookId, reviewDataFromDB.UserId)
	})
	t.Run("When Failure Delete review Data", func(t *testing.T) {
		t.Run("Reviews doesn't exist", func(t *testing.T) {
			reviewRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(reviews.Domain{}, errors.New("review doesn't exist")).Once()

			bookId, statusCode, err := reviewUsecase.Delete(context.Background(), reviewDataFromDB.UserId, reviewDataFromDB.ID)

			assert.Equal(t, errors.New("review not found"), err)
			assert.Equal(t, http.StatusNotFound, statusCode)
			assert.Equal(t, 0, bookId)
		})

		t.Run("Failed Delete Review", func(t *testing.T) {
			reviewRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(reviewDataFromDB, nil).Once()
			reviewRepository.Mock.On("Delete", mock.Anything, mock.AnythingOfType("*reviews.Domain")).Return(0, errors.New("failed")).Once()

			bookId, statusCode, err := reviewUsecase.Delete(context.Background(), reviewDataFromDB.UserId, reviewDataFromDB.ID)

			assert.Equal(t, errors.New("failed"), err)
			assert.Equal(t, http.StatusInternalServerError, statusCode)
			assert.Equal(t, 0, bookId)
		})
	})
}

func TestUpdate(t *testing.T) {
	setup(t)
	t.Run("When Success Update Review Data", func(t *testing.T) {
		t.Run("When Success Update review", func(t *testing.T) {
			reviewRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(reviewDataFromDB, nil).Once()
			reviewRepository.Mock.On("Update", mock.Anything, mock.AnythingOfType("*reviews.Domain")).Return(nil).Once()
			reviewRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(reviewDataFromDB, nil).Once()

			result, statusCode, err := reviewUsecase.Update(context.Background(), &reviewDataFromDB, reviewDataFromDB.UserId, reviewDataFromDB.ID)

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, statusCode)
			assert.Equal(t, reviewDataFromDB, result)
		})
	})
	t.Run("When Failure Update Review Data", func(t *testing.T) {
		t.Run("Review Doesn't Exists", func(t *testing.T) {
			reviewRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(reviews.Domain{}, errors.New("review doesn't exist")).Once()

			result, statusCode, err := reviewUsecase.Update(context.Background(), &reviewDataFromDB, reviewDataFromDB.UserId, reviewDataFromDB.ID)

			assert.Equal(t, errors.New("review not found"), err)
			assert.Equal(t, reviews.Domain{}, result)
			assert.Equal(t, http.StatusNotFound, statusCode)
		})
		t.Run("User Don't Have permissions", func(t *testing.T) {
			reviewRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(reviewDataFromDB, nil).Once()

			result, statusCode, err := reviewUsecase.Update(context.Background(), &reviewDataFromDB, reviewDataFromDB.UserId+3, reviewDataFromDB.ID)

			assert.Equal(t, errors.New("you don't have access to update this review"), err)
			assert.Equal(t, reviews.Domain{}, result)
			assert.Equal(t, http.StatusUnauthorized, statusCode)
		})
	})
}

func TestGetByBookId(t *testing.T) {
	setup(t)
	t.Run("When Success Get review Data", func(t *testing.T) {
		reviewRepository.Mock.On("GetByBookId", mock.Anything, mock.AnythingOfType("int")).Return([]reviews.Domain{reviewDataFromDB}, nil).Once()

		result, statusCode, err := reviewUsecase.GetByBookId(context.Background(), reviewDataFromDB.BookId)

		assert.Equal(t, []reviews.Domain{reviewDataFromDB}, result)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Nil(t, err)
	})

	t.Run("When Failure Review doesn't exist", func(t *testing.T) {
		reviewRepository.Mock.On("GetByBookId", mock.Anything, mock.AnythingOfType("int")).Return([]reviews.Domain{}, errors.New("review doesn't exist")).Once()

		result, statusCode, err := reviewUsecase.GetByBookId(context.Background(), reviewDataFromDB.BookId)

		assert.Equal(t, []reviews.Domain{}, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, statusCode)
	})
}

func TestGetByUserId(t *testing.T) {
	setup(t)
	t.Run("When Success Get review Data", func(t *testing.T) {
		reviewRepository.Mock.On("GetByUserId", mock.Anything, mock.AnythingOfType("int")).Return([]reviews.Domain{reviewDataFromDB}, nil).Once()

		result, statusCode, err := reviewUsecase.GetByUserId(context.Background(), reviewDataFromDB.UserId)

		assert.Equal(t, []reviews.Domain{reviewDataFromDB}, result)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Nil(t, err)
	})

	t.Run("When Failure Review doesn't exist", func(t *testing.T) {
		reviewRepository.Mock.On("GetByUserId", mock.Anything, mock.AnythingOfType("int")).Return([]reviews.Domain{}, errors.New("review doesn't exist")).Once()

		result, statusCode, err := reviewUsecase.GetByUserId(context.Background(), reviewDataFromDB.UserId)

		assert.NotNil(t, err)
		assert.Equal(t, []reviews.Domain{}, result)
		assert.Equal(t, http.StatusNotFound, statusCode)
	})
}

func TestGetUserReview(t *testing.T) {
	setup(t)
	t.Run("When Success Get User Review", func(t *testing.T) {
		reviewRepository.Mock.On("GetUserReview", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(reviewDataFromDB, nil).Once()

		result, statusCode, err := reviewUsecase.GetUserReview(context.Background(), reviewDataFromDB.BookId, reviewDataFromDB.UserId)

		assert.Equal(t, reviewDataFromDB, result)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Nil(t, err)
	})

	t.Run("When Failure User Review Review doesn't exist", func(t *testing.T) {
		reviewRepository.Mock.On("GetUserReview", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(reviews.Domain{}, errors.New("review doesn't exist")).Once()

		result, statusCode, err := reviewUsecase.GetUserReview(context.Background(), reviewDataFromDB.BookId, reviewDataFromDB.UserId)

		assert.Equal(t, reviews.Domain{}, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, statusCode)
	})
}
