package users_test

import (
	"context"
	"errors"
	"testing"
	"time"

	repositoryMocks "github.com/snykk/golib_backend/datasources/databases/users/mocks"
	"github.com/snykk/golib_backend/domains/users"
	"github.com/snykk/golib_backend/helpers"
	"github.com/snykk/golib_backend/http/controllers/users/request"
	jwtMocks "github.com/snykk/golib_backend/http/token/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	jwtService      *jwtMocks.JWTService
	userRepository  *repositoryMocks.Repository
	userUsecase     users.Usecase
	usersDataFromDB []users.Domain
	userDataFromDB  users.Domain
)

func setup(t *testing.T) {
	jwtService = jwtMocks.NewJWTService(t)
	userRepository = repositoryMocks.NewRepository(t)
	userUsecase = users.NewUserUsecase(userRepository, jwtService)
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
			CreatedAt:   time.Now(),
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
			CreatedAt:   time.Now(),
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
		CreatedAt:   time.Now(),
	}
}

func TestStore(t *testing.T) {
	setup(t)
	req := request.UserRequest{
		FullName: "patrick star",
		Username: "itsmepatrick",
		Email:    "najibfikri13@gmail.com",
		Password: "11111",
		Gender:   "male",
	}
	t.Run("When Success Store User Data", func(t *testing.T) {
		pass, _ := helpers.GenerateHash("11111")

		userRepository.Mock.On("Store", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(userDataFromDB, nil).Once()
		result, err := userUsecase.Store(context.Background(), req.ToDomain())

		assert.Nil(t, err)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, "patrick star", result.FullName)
		assert.Equal(t, "user", result.Role)
		assert.Equal(t, "male", result.Gender)
		assert.Equal(t, true, helpers.ValidateHash("11111", pass))
		assert.NotNil(t, result.CreatedAt)
	})

	t.Run("When Failure", func(t *testing.T) {
		userRepository.Mock.On("Store", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(users.Domain{}, errors.New("registration failed")).Once()
		result, err := userUsecase.Store(context.Background(), req.ToDomain())

		assert.NotNil(t, err)
		assert.Equal(t, 0, result.ID)
	})

}

func TestGetAll(t *testing.T) {
	setup(t)
	t.Run("When Success Get Users Data", func(t *testing.T) {
		userRepository.Mock.On("GetAll", mock.Anything).Return(usersDataFromDB, nil).Once()
		result, err := userUsecase.GetAll(context.Background())

		t.Run("Check User 1", func(t *testing.T) {
			assert.Nil(t, err)
			assert.Equal(t, 1, result[0].ID)
			assert.Equal(t, "patrick star", result[0].FullName)
			assert.Equal(t, "itsmepatrick", result[0].Username)
			assert.Equal(t, "najibfikri13@gmail.com", result[0].Email)
			assert.Equal(t, "11111", result[0].Password)
		})

		t.Run("Check User 2", func(t *testing.T) {
			assert.Nil(t, err)
			assert.Equal(t, 2, result[1].ID)
			assert.Equal(t, "john doe", result[1].FullName)
			assert.Equal(t, "johny", result[1].Username)
			assert.Equal(t, "johny123@gmail.com", result[1].Email)
			assert.Equal(t, "11111", result[1].Password)
		})
	})

	t.Run("When Failure Get Users Data", func(t *testing.T) {
		userRepository.Mock.On("GetAll", mock.Anything).Return([]users.Domain{}, errors.New("get all users failed")).Once()
		result, err := userUsecase.GetAll(context.Background())

		assert.NotNil(t, err)
		assert.Equal(t, []users.Domain{}, result)
	})
}

func TestGetById(t *testing.T) {
	setup(t)
	t.Run("When Success Get User Data By Id", func(t *testing.T) {
		t.Run("With User Itself", func(t *testing.T) {
			userRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(userDataFromDB, nil).Once()

			result, err := userUsecase.GetById(context.Background(), userDataFromDB.ID, userDataFromDB.ID)

			assert.Equal(t, userDataFromDB, result)
			assert.Nil(t, err)
			assert.NotEqual(t, "", result.Password)
		})

		t.Run("With Strangers", func(t *testing.T) {
			userDataFromDB.Password = ""
			userRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(userDataFromDB, nil).Once()

			result, err := userUsecase.GetById(context.Background(), userDataFromDB.ID, userDataFromDB.ID+1)

			assert.Equal(t, userDataFromDB, result)
			assert.Nil(t, err)
			assert.Equal(t, "", result.Password)
		})
	})

	t.Run("When Failure User doesn't exist", func(t *testing.T) {
		userRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(users.Domain{}, errors.New("user doesn't exist")).Once()

		result, err := userUsecase.GetById(context.Background(), userDataFromDB.ID, userDataFromDB.ID)

		assert.Equal(t, users.Domain{}, result)
		assert.Equal(t, errors.New("user doesn't exist"), err)
	})
}

func TestDelete(t *testing.T) {
	setup(t)
	t.Run("When Success Delete User Data", func(t *testing.T) {
		userRepository.Mock.On("Delete", mock.Anything, mock.AnythingOfType("int")).Return(nil).Once()

		err := userUsecase.Delete(context.Background(), userDataFromDB.ID)

		assert.Nil(t, err)
	})
	t.Run("When Failure Delete User Data", func(t *testing.T) {
		userRepository.Mock.On("Delete", mock.Anything, mock.AnythingOfType("int")).Return(errors.New("failed")).Once()

		err := userUsecase.Delete(context.Background(), 1)

		assert.Equal(t, errors.New("failed"), err)
	})
}

func TestUpdate(t *testing.T) {
	setup(t)
	t.Run("When Success Update User", func(t *testing.T) {
		data := userDataFromDB
		data.UpdatedAt = time.Now()

		userRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(data, nil).Once()
		userRepository.Mock.On("Update", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(nil).Once()

		result, err := userUsecase.Update(context.Background(), &data, data.ID)

		assert.Equal(t, data, result)
		assert.Nil(t, err)
		assert.NotNil(t, result.UpdatedAt)
	})
}

func TestLogin(t *testing.T) {
	setup(t)
	t.Run("When Success Login", func(t *testing.T) {
		req := request.UserLoginRequest{
			Email:    "najibfikri13@gmail.com",
			Password: "11111",
		}
		userDataFromDB.IsActivated = true
		userDataFromDB.Password, _ = helpers.GenerateHash(userDataFromDB.Password)

		userRepository.Mock.On("GetByEmail", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(userDataFromDB, nil).Once()
		jwtService.Mock.On("GenerateToken", mock.AnythingOfType("int"), mock.AnythingOfType("bool"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("eyBlablablabla", nil).Once()

		result, err := userUsecase.Login(context.Background(), req.ToDomain())

		assert.NotNil(t, result)
		assert.Nil(t, err)
		assert.Contains(t, result.Token, "ey")
	})
	t.Run("When Failure Account Not Activated Yet", func(t *testing.T) {
		t.Run("Account Not Activated Yet", func(t *testing.T) {
			req := request.UserLoginRequest{
				Email:    "najibfikri13@gmail.com",
				Password: "11111",
			}
			userDataFromDB.IsActivated = false
			userDataFromDB.Password, _ = helpers.GenerateHash(userDataFromDB.Password)

			userRepository.Mock.On("GetByEmail", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(userDataFromDB, nil).Once()
			result, err := userUsecase.Login(context.Background(), req.ToDomain())

			assert.Equal(t, users.Domain{}, result)
			assert.NotNil(t, err)
			assert.Equal(t, "", result.Token)
		})
		t.Run("Invalid Credential", func(t *testing.T) {
			req := request.UserLoginRequest{
				Email:    "najibfikri13@gmail.com",
				Password: "111112",
			}
			userDataFromDB.IsActivated = true
			userDataFromDB.Password, _ = helpers.GenerateHash(userDataFromDB.Password)

			userRepository.Mock.On("GetByEmail", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(userDataFromDB, nil).Once()

			result, err := userUsecase.Login(context.Background(), req.ToDomain())

			assert.Equal(t, users.Domain{}, result)
			assert.NotNil(t, err)
			assert.Equal(t, "", result.Token)
		})
	})
}

func TestGetByEmail(t *testing.T) {
	setup(t)
	t.Run("When Success Get User Data By Email", func(t *testing.T) {
		userRepository.Mock.On("GetByEmail", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(userDataFromDB, nil).Once()

		result, err := userUsecase.GetByEmail(context.Background(), "najibfikri13@gmail.com")

		assert.Equal(t, userDataFromDB, result)
		assert.Nil(t, err)
	})

	t.Run("When Failure User doesn't exist", func(t *testing.T) {
		userRepository.Mock.On("GetByEmail", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(users.Domain{}, errors.New("user doesn't exist")).Once()

		result, err := userUsecase.GetByEmail(context.Background(), "johndoe@gmail.com")

		assert.Equal(t, users.Domain{}, result)
		assert.Equal(t, errors.New("user doesn't exist"), err)
	})
}

func TestActivate(t *testing.T) {
	setup(t)
	t.Run("When Success Activate Email", func(t *testing.T) {
		userRepository.Mock.On("GetByEmail", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(userDataFromDB, nil).Once()
		userRepository.Mock.On("Update", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(nil).Once()

		err := userUsecase.ActivateUser(context.Background(), "najibfikri13@gmail.com")

		assert.Nil(t, err)
	})

	t.Run("When Failure Activate Email", func(t *testing.T) {
		userRepository.Mock.On("GetByEmail", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(users.Domain{}, errors.New("user doesn't exist")).Once()

		result, err := userUsecase.GetByEmail(context.Background(), "johndoe@gmail.com")

		assert.Equal(t, users.Domain{}, result)
		assert.Equal(t, errors.New("user doesn't exist"), err)
	})
}

func TestChangePassword(t *testing.T) {
	setup(t)
	newPass := "newPass"
	t.Run("When Success Change Password", func(t *testing.T) {
		userDomDB := userDataFromDB
		userDomDB.Password, _ = helpers.GenerateHash(userDataFromDB.Password)

		userRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(userDomDB, nil).Once()
		userRepository.Mock.On("Update", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(nil).Once()

		err := userUsecase.ChangePassword(context.Background(), &userDataFromDB, newPass, userDataFromDB.ID)

		assert.Nil(t, err)
	})

	t.Run("When Failure Change Password", func(t *testing.T) {
		userRepository.Mock.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(users.Domain{}, errors.New("user doesn't exist")).Once()

		err := userUsecase.ChangePassword(context.Background(), &userDataFromDB, newPass, userDataFromDB.ID)

		assert.NotNil(t, err)
	})
}

func TestChangeEmail(t *testing.T) {
	setup(t)
	t.Run("When Success Change Email", func(t *testing.T) {
		userDataFromDB.Email = "newemail@gmail.com"

		userRepository.Mock.On("GetByEmail", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(users.Domain{}, errors.New("user not found")).Once()
		userRepository.Mock.On("UpdateEmail", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(nil).Once()
		err := userUsecase.ChangeEmail(context.Background(), &userDataFromDB, userDataFromDB.ID)

		assert.Nil(t, err)
	})
	t.Run("When Failure Change Email", func(t *testing.T) {
		userRepository.Mock.On("GetByEmail", mock.Anything, mock.AnythingOfType("*users.Domain")).Return(userDataFromDB, nil).Once()
		err := userUsecase.ChangeEmail(context.Background(), &userDataFromDB, userDataFromDB.ID)

		assert.NotNil(t, err)
	})
}
