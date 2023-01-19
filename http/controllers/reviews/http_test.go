package reviews_test

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
	reviewMocks "github.com/snykk/golib_backend/datasources/databases/reviews/mocks"
	"github.com/snykk/golib_backend/domains/books"
	"github.com/snykk/golib_backend/domains/reviews"
	"github.com/snykk/golib_backend/domains/users"
	"github.com/snykk/golib_backend/helpers"
	controllers "github.com/snykk/golib_backend/http/controllers/reviews"
	"github.com/snykk/golib_backend/http/controllers/reviews/requests"
	"github.com/snykk/golib_backend/http/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	reviewRepository  *reviewMocks.Repository
	reviewUsecase     reviews.Usecase
	reviewController  controllers.ReviewController
	ristrettoMock     *cacheMocks.RistrettoCache
	s                 *gin.Engine
	reviewsDataFromDB []reviews.Domain
	reviewDataFromDB  reviews.Domain
	bookFromDB        books.Domain
	userFromDB        users.Domain
)

func setup(t *testing.T) {
	ristrettoMock = cacheMocks.NewRistrettoCache(t)
	reviewRepository = reviewMocks.NewRepository(t)
	reviewUsecase = reviews.NewReviewUsecase(reviewRepository)
	reviewController = controllers.NewReviewController(reviewUsecase, ristrettoMock)

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

	// Create gin engine
	s = gin.Default()
	s.Use(lazyAuth)
}

func lazyAuth(ctx *gin.Context) {
	// hash
	pass, _ := helpers.GenerateHash(userFromDB.Password)
	// prepare claims
	jwtClaims := token.JwtCustomClaim{
		UserID:   userFromDB.ID,
		IsAdmin:  true,
		Email:    userFromDB.Email,
		Password: pass,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(config.AppConfig.JWTExpired)).Unix(),
			Issuer:    userFromDB.Username,
			IssuedAt:  time.Now().Unix(),
		},
	}
	ctx.Set(constants.CtxAuthenticatedUserKey, jwtClaims)
}

func TestStore(t *testing.T) {
	setup(t)
	// Define route
	s.POST("/reviews", reviewController.Store)
	t.Run("When Success Create Review", func(t *testing.T) {
		req := requests.ReviewRequest{
			Text:   "gege bet yagesya bintang 9",
			Rating: 9,
			BookId: bookFromDB.ID,
		}
		reqBody, _ := json.Marshal(req)

		reviewRepository.Mock.On("GetUserReview", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(reviews.Domain{}, errors.New("reviews not found")).Once() // when user does'nt have review yet
		reviewRepository.Mock.On("Store", mock.Anything, mock.AnythingOfType("*reviews.Domain")).Return(reviewDataFromDB, nil).Once()
		ristrettoMock.Mock.On("Del", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader(reqBody))

		r.Header.Set("Content-Type", "application/json")

		// Perform requests
		s.ServeHTTP(w, r)

		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "review created successfully")
	})
	t.Run("When Failure", func(t *testing.T) {
		t.Run("When Request is Empty", func(t *testing.T) {
			req := requests.ReviewRequest{}
			reqBody, _ := json.Marshal(req)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader(reqBody))

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
		t.Run("When Invalid Rating", func(t *testing.T) {
			req := requests.ReviewRequest{
				Text:   "gege bet yagesya bintang 15",
				Rating: 15,
				BookId: bookFromDB.ID,
			}
			reqBody, _ := json.Marshal(req)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader(reqBody))

			r.Header.Set("Content-Type", "application/json")

			// Perform requests
			s.ServeHTTP(w, r)

			body := w.Body.String()

			// Assertions
			// Assert status code
			assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
			assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
			assert.Contains(t, body, "the rating must be in the range 1 - 10")
		})
		t.Run("When Unexpexted Error", func(t *testing.T) {
			req := requests.ReviewRequest{
				Text:   "gege bet yagesya bintang 9",
				Rating: 9,
				BookId: bookFromDB.ID,
			}
			reqBody, _ := json.Marshal(req)

			reviewRepository.Mock.On("GetUserReview", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(reviews.Domain{}, errors.New("reviews not found")).Once() // when user does'nt have review yet
			reviewRepository.Mock.On("Store", mock.Anything, mock.AnythingOfType("*reviews.Domain")).Return(reviews.Domain{}, constants.ErrUnexpected).Once()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader(reqBody))

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
	s.GET("/reviews", reviewController.GetAll)
	t.Run("When Success", func(t *testing.T) {
		t.Run("Fetched review Data", func(t *testing.T) {
			reviewRepository.Mock.On("GetAll", mock.Anything).Return(reviewsDataFromDB, nil).Once()
			ristrettoMock.Mock.On("Get", "reviews").Return(nil).Once()
			ristrettoMock.Mock.On("Set", "reviews", mock.Anything).Once()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/reviews", nil)

			r.Header.Set("Content-Type", "application/json")

			// Perform requests
			s.ServeHTTP(w, r)

			// parsing json to raw text
			body := w.Body.String()

			// Assertions
			// Assert status code
			assert.Equal(t, http.StatusOK, w.Result().StatusCode)
			assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
			assert.Contains(t, body, "review data fetched successfully")
		})
		t.Run("When Empty Data", func(t *testing.T) {
			reviewRepository.Mock.On("GetAll", mock.Anything).Return([]reviews.Domain{}, nil).Once()
			ristrettoMock.Mock.On("Get", "reviews").Return(nil).Once()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/reviews", nil)

			r.Header.Set("Content-Type", "application/json")

			// Perform requests
			s.ServeHTTP(w, r)

			body := w.Body.String()

			// Assertions
			// Assert status code
			assert.Equal(t, http.StatusOK, w.Result().StatusCode)
			assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
			assert.Contains(t, body, "review data is empty")
		})
	})
	t.Run("When Failure", func(t *testing.T) {
		reviewRepository.Mock.On("GetAll", mock.Anything).Return([]reviews.Domain{}, constants.ErrUnexpected).Once()
		ristrettoMock.Mock.On("Get", "reviews").Return(nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/reviews", nil)

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
	s.GET("/reviews/:id", reviewController.GetById)

	id := 1
	t.Run("When Success Fetched review Data By Id", func(t *testing.T) {
		reviewRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(reviewDataFromDB, nil).Once()
		ristrettoMock.Mock.On("Get", fmt.Sprintf("review/%d", id)).Return(nil).Once()
		ristrettoMock.Mock.On("Set", fmt.Sprintf("review/%d", id), mock.Anything).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/reviews/%d", id), nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform requests
		s.ServeHTTP(w, r)

		// parsing json to raw text
		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, fmt.Sprintf("review data with id %d fetched successfully", id))
	})
	t.Run("When Failure Fetched reviews Data", func(t *testing.T) {
		ristrettoMock.Mock.On("Get", fmt.Sprintf("review/%d", id)).Return(nil).Once()
		reviewRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(reviews.Domain{}, constants.ErrUnexpected).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/reviews/%d", id), nil)

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
	s.PUT("/reviews/:id", reviewController.Update)
	t.Run("When Success Create Review", func(t *testing.T) {
		updatedReview := reviewDataFromDB
		req := requests.ReviewRequest{
			Text:   "gege bet yagesya bintang 9",
			Rating: 9,
			BookId: bookFromDB.ID,
		}
		updatedReview.Text = req.Text
		updatedReview.Rating = 9
		updatedReview.UpdatedAt = time.Now()

		reqBody, _ := json.Marshal(req)

		reviewRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(reviewDataFromDB, nil).Once() // when user does'nt have review yet
		reviewRepository.Mock.On("Update", mock.Anything, mock.AnythingOfType("*reviews.Domain")).Return(nil).Once()
		reviewRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(updatedReview, nil).Once() // when user does'nt have review yet
		ristrettoMock.Mock.On("Del", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/reviews/%d", reviewDataFromDB.ID), bytes.NewReader(reqBody))

		r.Header.Set("Content-Type", "application/json")

		// Perform requests
		s.ServeHTTP(w, r)

		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "review updated successfully")
	})
	t.Run("When Failure", func(t *testing.T) {
		t.Run("When Request is Empty", func(t *testing.T) {
			req := requests.ReviewRequest{}
			reqBody, _ := json.Marshal(req)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/reviews/%d", reviewDataFromDB.ID), bytes.NewReader(reqBody))

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
		t.Run("When Invalid Rating", func(t *testing.T) {
			req := requests.ReviewRequest{
				Text:   "gege bet yagesya bintang 15",
				Rating: 15,
				BookId: bookFromDB.ID,
			}
			reqBody, _ := json.Marshal(req)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/reviews/%d", reviewDataFromDB.ID), bytes.NewReader(reqBody))

			r.Header.Set("Content-Type", "application/json")

			// Perform requests
			s.ServeHTTP(w, r)

			body := w.Body.String()

			// Assertions
			// Assert status code
			assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
			assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
			assert.Contains(t, body, "the rating must be in the range 1 - 10")
		})
		t.Run("When Unexpexted Error", func(t *testing.T) {
			req := requests.ReviewRequest{
				Text:   "gege bet yagesya bintang 9",
				Rating: 9,
				BookId: bookFromDB.ID,
			}
			reqBody, _ := json.Marshal(req)

			reviewRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(reviewDataFromDB, nil).Once()
			reviewRepository.Mock.On("Update", mock.Anything, mock.AnythingOfType("*reviews.Domain")).Return(constants.ErrUnexpected).Once()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/reviews/%d", reviewDataFromDB.ID), bytes.NewReader(reqBody))

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
	s.DELETE("/reviews/:id", reviewController.Delete)
	t.Run("When Success Delete review Data", func(t *testing.T) {
		reviewRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(reviewDataFromDB, nil).Once()
		reviewRepository.Mock.On("Delete", mock.Anything, mock.AnythingOfType("*reviews.Domain")).Return(reviewDataFromDB.BookId, nil).Once()
		ristrettoMock.Mock.On("Del", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/reviews/%d", reviewDataFromDB.ID), nil)

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
	t.Run("When Failure review Not Found", func(t *testing.T) {
		reviewRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(reviews.Domain{}, errors.New("review not found")).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/reviews/%d", reviewDataFromDB.ID), nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform requests
		s.ServeHTTP(w, r)

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
	})
}
