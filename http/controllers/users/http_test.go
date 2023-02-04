package users_test

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
	userMocks "github.com/snykk/golib_backend/datasources/databases/users/mocks"
	"github.com/snykk/golib_backend/domains/users"
	"github.com/snykk/golib_backend/helpers"
	controllers "github.com/snykk/golib_backend/http/controllers/users"
	"github.com/snykk/golib_backend/http/controllers/users/request"
	"github.com/snykk/golib_backend/http/token"
	jwtMocks "github.com/snykk/golib_backend/http/token/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	jwtService      *jwtMocks.JWTService
	userRepository  *userMocks.Repository
	userUsecase     users.Usecase
	userController  controllers.UserController
	usersDataFromDB []users.Domain
	userDataFromDB  users.Domain
	redisMock       *cacheMocks.RedisCache
	ristrettoMock   *cacheMocks.RistrettoCache
	s               *gin.Engine
)

func setup(t *testing.T) {
	jwtService = jwtMocks.NewJWTService(t)
	redisMock = cacheMocks.NewRedisCache(t)
	ristrettoMock = cacheMocks.NewRistrettoCache(t)
	userRepository = userMocks.NewRepository(t)
	userUsecase = users.NewUserUsecase(userRepository, jwtService)
	userController = controllers.NewUserController(userUsecase, redisMock, ristrettoMock)

	usersDataFromDB = []users.Domain{
		{
			ID:          1,
			FullName:    "patrick star",
			Username:    "itsmepatrick",
			Email:       "najibfikri13@gmail.com",
			Password:    "11111",
			Role:        "admin",
			Gender:      "male",
			Reviews:     0,
			IsActivated: true,
		},
		{
			ID:          2,
			FullName:    "john doe",
			Username:    "johny",
			Email:       "johny123@gmail.com",
			Password:    "11111",
			Role:        "user",
			Gender:      "male",
			Reviews:     0,
			IsActivated: true,
		},
	}

	userDataFromDB = users.Domain{
		ID:          1,
		FullName:    "patrick star",
		Username:    "itsmepatrick",
		Email:       "najibfikri13@gmail.com",
		Password:    "11111",
		Role:        "user",
		Gender:      "male",
		Reviews:     0,
		IsActivated: false,
	}

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
		IsAdmin:  false,
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

func TestRegis(t *testing.T) {
	setup(t)
	// Define route
	s.POST("/auth/regis", userController.Regis)
	t.Run("When Success Regis", func(t *testing.T) {
		req := request.UserRequest{
			FullName: "patrick star",
			Username: "itsmepatrick",
			Email:    "patrick@gmail.com",
			Password: "11111",
			Gender:   "male",
		}
		reqBody, _ := json.Marshal(req)

		userRepository.Mock.On("Store", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(userDataFromDB, nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/auth/regis", bytes.NewReader(reqBody))

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		s.ServeHTTP(w, r)

		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "registration user success")
	})
	t.Run("When Failure", func(t *testing.T) {
		t.Run("When Request is Empty", func(t *testing.T) {
			req := request.UserRequest{}
			reqBody, _ := json.Marshal(req)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/auth/regis", bytes.NewReader(reqBody))

			r.Header.Set("Content-Type", "application/json")

			// Perform request
			s.ServeHTTP(w, r)

			body := w.Body.String()

			// Assertions
			// Assert status code
			assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
			assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
			assert.Contains(t, body, "failed on the 'required' tag")
		})
		t.Run("Gender is Not Valid", func(t *testing.T) {
			req := request.UserRequest{
				FullName: "patrick star",
				Username: "itsmepatrick",
				Email:    "patrick@gmail.com",
				Password: "11111",
				Gender:   "malee",
			}
			reqBody, _ := json.Marshal(req)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/auth/regis", bytes.NewReader(reqBody))

			r.Header.Set("Content-Type", "application/json")

			// Perform request
			s.ServeHTTP(w, r)

			body := w.Body.String()

			// Assertions
			// Assert status code
			assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
			assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
			assert.Contains(t, body, "gender must be one of")
		})
	})
}

func TestLogin(t *testing.T) {
	setup(t)
	// Define route
	s.POST("/auth/login", userController.Login)
	t.Run("When Success Login", func(t *testing.T) {
		// hash password field
		var err error
		userDataFromDB.Password, err = helpers.GenerateHash(userDataFromDB.Password)
		if err != nil {
			t.Error(err)
		}
		// make account activated
		userDataFromDB.IsActivated = true
		req := request.UserLoginRequest{
			Email:    "patrick@gmail.com",
			Password: "11111",
		}
		reqBody, _ := json.Marshal(req)

		userRepository.Mock.On("GetByEmail", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(userDataFromDB, nil).Once()
		jwtService.Mock.On("GenerateToken", mock.AnythingOfType("int"), mock.AnythingOfType("bool"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("eyBlablablabla", nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(reqBody))

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		s.ServeHTTP(w, r)

		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "login success")
		assert.Contains(t, body, "eyBlablablabla")
	})
	t.Run("When Failure User is Not Exists", func(t *testing.T) {
		req := request.UserLoginRequest{
			Email:    "patrick312@gmail.com",
			Password: "11111",
		}
		reqBody, _ := json.Marshal(req)

		userRepository.Mock.On("GetByEmail", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(users.Domain{}, errors.New("user not found")).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(reqBody))

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		s.ServeHTTP(w, r)

		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "user not found")
	})
}

func TestGetAll(t *testing.T) {
	setup(t)
	// Define route
	s.GET("/users", userController.GetAll)
	t.Run("When Success Fetched User Data", func(t *testing.T) {
		userRepository.Mock.On("GetAll", mock.Anything).Return(usersDataFromDB, nil).Once()
		ristrettoMock.Mock.On("Get", "users").Return(nil).Once()
		ristrettoMock.Mock.On("Set", "users", mock.Anything).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/users", nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		s.ServeHTTP(w, r)

		// parsing json to raw text
		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "user data fetched successfully")
	})
	t.Run("When Failure Fetched Users Data", func(t *testing.T) {
		userRepository.Mock.On("GetAll", mock.Anything).Return([]users.Domain{}, nil).Once()
		ristrettoMock.Mock.On("Get", "users").Return(nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/users", nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		s.ServeHTTP(w, r)

		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "user data is empty")
	})
}

func TestGetById(t *testing.T) {
	setup(t)
	// Define route
	s.GET("/users/:id", userController.GetById)

	id := 1
	t.Run("When Success Fetched User Data By Id", func(t *testing.T) {
		userRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(userDataFromDB, nil).Once()
		ristrettoMock.Mock.On("Get", fmt.Sprintf("user/%d", id)).Return(nil).Once()
		ristrettoMock.Mock.On("Set", fmt.Sprintf("user/%d", id), mock.Anything).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d", id), nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		s.ServeHTTP(w, r)

		// parsing json to raw text
		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, fmt.Sprintf("user data with id %d fetched successfully", id))
	})
	t.Run("When Failure Fetched Users Data", func(t *testing.T) {
		ristrettoMock.Mock.On("Get", fmt.Sprintf("user/%d", id)).Return(nil).Once()
		userRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(users.Domain{}, errors.New("user not found")).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d", id), nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		s.ServeHTTP(w, r)

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
	})
}

func TestGetUserData(t *testing.T) {
	setup(t)
	// Define route
	s.GET("/users/me", userController.GetUserData)

	emailUserAuthenticated := userDataFromDB.Email
	t.Run("When Success Fetched User Data", func(t *testing.T) {
		userRepository.Mock.On("GetByEmail", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(userDataFromDB, nil).Once()
		ristrettoMock.Mock.On("Get", fmt.Sprintf("user/%s", emailUserAuthenticated)).Return(nil).Once()
		ristrettoMock.Mock.On("Set", fmt.Sprintf("user/%s", emailUserAuthenticated), mock.Anything).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/users/me", nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		s.ServeHTTP(w, r)

		// parsing json to raw text
		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "user data fetched successfully")
	})

	t.Run("When Failure Fetched User Data", func(t *testing.T) {
		userRepository.Mock.On("GetByEmail", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(users.Domain{}, constants.ErrUnexpected).Once()
		ristrettoMock.Mock.On("Get", fmt.Sprintf("user/%s", emailUserAuthenticated)).Return(nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/users/me", nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		s.ServeHTTP(w, r)

		// Assertions
		// Assert status code
		assert.NotEqual(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
	})
}

func TestUpdate(t *testing.T) {
	setup(t)
	// Define route
	s.PUT("/users", userController.Update)
	t.Run("When Success Update User Data", func(t *testing.T) {
		req := request.UserUpdateRequest{
			FullName: "patrick star edited",
			Username: "itsmepatrick",
			Gender:   "male",
		}
		reqBody, _ := json.Marshal(req)

		userDataFromDB.FullName = "patrick star edited"

		userRepository.Mock.On("Update", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(nil).Once()
		userRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(userDataFromDB, nil).Once()
		ristrettoMock.Mock.On("Del", mock.Anything, mock.Anything)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/users", bytes.NewReader(reqBody))

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		s.ServeHTTP(w, r)

		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "updated successfully")
	})
	t.Run("When Failure", func(t *testing.T) {
		req := request.UserUpdateRequest{
			FullName: "patrick star edited",
			Username: "itsmepatrick",
			Gender:   "male",
		}
		reqBody, _ := json.Marshal(req)

		userDataFromDB.FullName = "patrick star edited"

		userRepository.Mock.On("Update", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(constants.ErrUnexpected).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/users", bytes.NewReader(reqBody))

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		s.ServeHTTP(w, r)

		// Assertions
		// Assert status code
		assert.NotEqual(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
	})
}

func TestDelete(t *testing.T) {
	setup(t)
	// Define route
	s.DELETE("/users", userController.Delete)
	t.Run("When Success Delete User Data", func(t *testing.T) {
		userRepository.Mock.On("Delete", mock.Anything, mock.AnythingOfType("int")).Return(nil).Once()
		ristrettoMock.Mock.On("Del", mock.Anything, mock.Anything)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/users", nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		s.ServeHTTP(w, r)

		body := w.Body.String()

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "deleted successfully")
	})
	t.Run("When Failure", func(t *testing.T) {
		userRepository.Mock.On("Delete", mock.Anything, mock.AnythingOfType("int")).Return(constants.ErrUnexpected).Once()
		ristrettoMock.Mock.On("Del", mock.Anything, mock.Anything)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/users", nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		s.ServeHTTP(w, r)

		// Assertions
		// Assert status code
		assert.NotEqual(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
	})
}
