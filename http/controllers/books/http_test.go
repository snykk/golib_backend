package books_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/config"
	"github.com/snykk/golib_backend/constants"
	cacheMocks "github.com/snykk/golib_backend/datasources/cache/mocks"
	bookMocks "github.com/snykk/golib_backend/datasources/databases/books/mocks"
	"github.com/snykk/golib_backend/domains/books"
	"github.com/snykk/golib_backend/domains/users"
	"github.com/snykk/golib_backend/helpers"
	controllers "github.com/snykk/golib_backend/http/controllers/books"
	"github.com/snykk/golib_backend/http/controllers/books/requests"
	"github.com/snykk/golib_backend/http/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	bookRepository  *bookMocks.Repository
	bookUsecase     books.Usecase
	bookController  controllers.BookController
	booksDataFromDB []books.Domain
	bookDataFromDB  books.Domain
	ristrettoMock   *cacheMocks.RistrettoCache
	s               *gin.Engine
	userDataFromDB  users.Domain
)

func setup(t *testing.T) {
	ristrettoMock = cacheMocks.NewRistrettoCache(t)
	bookRepository = bookMocks.NewRepository(t)
	bookUsecase = books.NewBookUsecase(bookRepository)
	bookController = controllers.NewBookController(bookUsecase, ristrettoMock)

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

	// Create gin engine
	s = gin.Default()
	s.Use(lazyAuth)
}

func lazyAuth(ctx *gin.Context) {
	// hash
	pass, _ := helpers.GenerateHash(userDataFromDB.Password)
	// prepare claims
	jwtClaims := token.JwtCustomClaim{
		UserID:   userDataFromDB.ID,
		IsAdmin:  true,
		Email:    userDataFromDB.Email,
		Password: pass,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(config.AppConfig.JWTExpired)).Unix(),
			Issuer:    userDataFromDB.Username,
			IssuedAt:  time.Now().Unix(),
		},
	}
	ctx.Set(constants.CtxAuthenticatedUserKey, jwtClaims)
}

func TestStore(t *testing.T) {
	setup(t)
	// Define route
	s.POST("/books", bookController.Store)
	t.Run("When Success Regis", func(t *testing.T) {
		req := requests.BookRequest{
			Title:       "Atomic Habits",
			Author:      "James Clear",
			Description: "lorem ipsum doler sit amet",
			Publisher:   "Gramedia",
			ISBN:        "111111",
		}
		reqBody, _ := json.Marshal(req)

		bookRepository.Mock.On("Store", mock.Anything, mock.AnythingOfType("*books.Domain")).Return(bookDataFromDB, nil).Once()
		ristrettoMock.Mock.On("Del", "books")

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(reqBody))

		r.Header.Set("Content-Type", "application/json")

		// Perform requests
		s.ServeHTTP(w, r)

		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "book inserted successfully")
	})
	t.Run("When Failure", func(t *testing.T) {
		t.Run("When Request is Empty", func(t *testing.T) {
			req := requests.BookRequest{}
			reqBody, _ := json.Marshal(req)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(reqBody))

			r.Header.Set("Content-Type", "application/json")

			// Perform requests
			s.ServeHTTP(w, r)

			body := w.Body.String()

			// Assertions
			// Assert status code
			assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
			assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
			assert.Contains(t, body, "failed on the 'required' tag")
		})
		t.Run("When Unexpexted Error", func(t *testing.T) {
			req := requests.BookRequest{
				Title:       "Atomic Habits",
				Author:      "James Clear",
				Description: "lorem ipsum doler sit amet",
				Publisher:   "Gramedia",
				ISBN:        "111111",
			}
			reqBody, _ := json.Marshal(req)

			bookRepository.Mock.On("Store", mock.Anything, mock.AnythingOfType("*books.Domain")).Return(books.Domain{}, constants.ErrUnexpected).Once()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(reqBody))

			r.Header.Set("Content-Type", "application/json")

			// Perform requests
			s.ServeHTTP(w, r)

			// Assertions
			// Assert status code
			assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
			assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		})
	})
}

func TestGetAll(t *testing.T) {
	setup(t)
	// Define route
	s.GET("/books", bookController.GetAll)
	t.Run("When Success", func(t *testing.T) {
		t.Run("Fetched Book Data", func(t *testing.T) {
			bookRepository.Mock.On("GetAll", mock.Anything).Return(booksDataFromDB, nil).Once()
			ristrettoMock.Mock.On("Get", "books").Return(nil).Once()
			ristrettoMock.Mock.On("Set", "books", mock.Anything).Once()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/books", nil)

			r.Header.Set("Content-Type", "application/json")

			// Perform requests
			s.ServeHTTP(w, r)

			// parsing json to raw text
			body := w.Body.String()

			// Assertions
			// Assert status code
			assert.Equal(t, http.StatusOK, w.Result().StatusCode)
			assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
			assert.Contains(t, body, "book data fetched successfully")
		})
		t.Run("Empty Data", func(t *testing.T) {
			bookRepository.Mock.On("GetAll", mock.Anything).Return([]books.Domain{}, nil).Once()
			ristrettoMock.Mock.On("Get", "books").Return(nil).Once()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/books", nil)

			r.Header.Set("Content-Type", "application/json")

			// Perform requests
			s.ServeHTTP(w, r)

			body := w.Body.String()

			// Assertions
			// Assert status code
			assert.Equal(t, http.StatusOK, w.Result().StatusCode)
			assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
			assert.Contains(t, body, "book data is empty")
		})
	})
	t.Run("When Failure", func(t *testing.T) {
		bookRepository.Mock.On("GetAll", mock.Anything).Return([]books.Domain{}, constants.ErrUnexpected).Once()
		ristrettoMock.Mock.On("Get", "books").Return(nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/books", nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform requests
		s.ServeHTTP(w, r)

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
	})

}

func TestGetById(t *testing.T) {
	setup(t)
	// Define route
	s.GET("/books/:id", bookController.GetById)

	id := 1
	t.Run("When Success Fetched Book Data By Id", func(t *testing.T) {
		bookRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(bookDataFromDB, nil).Once()
		ristrettoMock.Mock.On("Get", fmt.Sprintf("book/%d", id)).Return(nil).Once()
		ristrettoMock.Mock.On("Set", fmt.Sprintf("book/%d", id), mock.Anything).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/books/%d", id), nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform requests
		s.ServeHTTP(w, r)

		// parsing json to raw text
		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, fmt.Sprintf("book data with id %d fetched successfully", id))
	})
	t.Run("When Failure Fetched books Data", func(t *testing.T) {
		ristrettoMock.Mock.On("Get", fmt.Sprintf("book/%d", id)).Return(nil).Once()
		bookRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(books.Domain{}, constants.ErrUnexpected).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/books/%d", id), nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform requests
		s.ServeHTTP(w, r)

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
	})
}

func TestUpdate(t *testing.T) {
	setup(t)
	// Define route
	s.PUT("/books/:id", bookController.Update)
	t.Run("When Success Update book Data", func(t *testing.T) {
		req := requests.BookRequest{
			Title:       "Atomic Habits edited",
			Author:      "James Clear",
			Description: "lorem ipsum doler sit amet",
			Publisher:   "Gramedia",
			ISBN:        "111111",
		}

		reqBody, _ := json.Marshal(req)

		bookDataFromDB.Title = "Atomic Habits edited"

		bookRepository.Mock.On("Update", mock.Anything, mock.AnythingOfType("*books.Domain")).Return(nil).Once()
		bookRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(bookDataFromDB, nil).Once()
		ristrettoMock.Mock.On("Del", mock.Anything, mock.Anything)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/books/%d", bookDataFromDB.ID), bytes.NewReader(reqBody))

		r.Header.Set("Content-Type", "application/json")

		// Perform requests
		s.ServeHTTP(w, r)

		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "updated successfully")
	})
	t.Run("When Failure", func(t *testing.T) {
		t.Run("When Request Empty", func(t *testing.T) {
			req := requests.BookRequest{}
			reqBody, _ := json.Marshal(req)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/books/%d", bookDataFromDB.ID), bytes.NewReader(reqBody))

			r.Header.Set("Content-Type", "application/json")

			// Perform requests
			s.ServeHTTP(w, r)

			body := w.Body.String()

			// Assertions
			// Assert status code
			assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
			assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
			assert.Contains(t, body, "failed on the 'required' tag")
		})
		t.Run("When Unexpected Error", func(t *testing.T) {
			req := requests.BookRequest{
				Title:       "Atomic Habits",
				Author:      "James Clear",
				Description: "lorem ipsum doler sit amet",
				Publisher:   "Gramedia",
				ISBN:        "111111",
			}
			reqBody, _ := json.Marshal(req)

			bookDataFromDB.Title = "Atomic Habits edited"

			bookRepository.Mock.On("Update", mock.Anything, mock.AnythingOfType("*books.Domain")).Return(constants.ErrUnexpected).Once()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/books/%d", bookDataFromDB.ID), bytes.NewReader(reqBody))

			r.Header.Set("Content-Type", "application/json")

			// Perform requests
			s.ServeHTTP(w, r)

			// Assertions
			// Assert status code
			assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
			assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		})
	})
}

func TestDelete(t *testing.T) {
	setup(t)
	// Define route
	s.DELETE("/books/:id", bookController.Delete)
	t.Run("When Success Delete book Data", func(t *testing.T) {
		bookRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(books.Domain{}, nil).Once()
		bookRepository.Mock.On("Delete", mock.Anything, mock.AnythingOfType("int")).Return(nil).Once()
		ristrettoMock.Mock.On("Del", mock.Anything, mock.Anything)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/books/%d", bookDataFromDB.ID), nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform requests
		s.ServeHTTP(w, r)

		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "deleted successfully")
	})
	t.Run("When Failure Book Not Found", func(t *testing.T) {
		bookRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(books.Domain{}, errors.New("book not found")).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/books/%d", bookDataFromDB.ID), nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform requests
		s.ServeHTTP(w, r)

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
	})
}
