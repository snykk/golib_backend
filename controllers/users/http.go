package users

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/controllers"
	"github.com/snykk/golib_backend/controllers/users/request"
	"github.com/snykk/golib_backend/controllers/users/responses"
	"github.com/snykk/golib_backend/usecase/users"

	validator "github.com/go-playground/validator/v10"
)

type UserController struct {
	UserUsecase users.Usecase
}

func NewUserController(usecase users.Usecase) *UserController {
	return &UserController{
		UserUsecase: usecase,
	}
}

func ValidateRequest(request interface{}) error {
	validator := validator.New()
	err := validator.Struct(request)
	return err
}

func (controller UserController) Regis(c *gin.Context) {
	var UserRegisRequest request.UserRegisRequest

	if err := c.Bind(&UserRegisRequest); err != nil {
		controllers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := ValidateRequest(UserRegisRequest); err != nil {
		controllers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user := UserRegisRequest.ToDomain()

	ctx := c.Request.Context()

	user, err := controller.UserUsecase.Store(ctx, &user)

	if err != nil {
		controllers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(c, "registration user success", map[string]interface{}{
		"user": responses.FromDomain(user),
	})
}

func (controller UserController) Login(c *gin.Context) {
	var UserLoginRequest request.UserLoginRequest
	if err := c.Bind(&UserLoginRequest); err != nil {
		controllers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := ValidateRequest(UserLoginRequest); err != nil {
		controllers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	result, err := controller.UserUsecase.Login(ctx, UserLoginRequest.ToDomain())
	if err != nil {
		controllers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(c, "login success", map[string]interface{}{
		"token": result.Token,
	})
}

func (controller UserController) GetAll(c *gin.Context) {
	usersFromUseCase, err := controller.UserUsecase.GetAll()

	if err != nil {
		controllers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if usersFromUseCase == nil {
		controllers.NewSuccessResponse(c, "user data is empty", map[string]interface{}{
			"users": []int{},
		})
		return
	}

	controllers.NewSuccessResponse(c, "user data fetched successfully", map[string]interface{}{
		"users": responses.ToResponseList(&usersFromUseCase),
	})
}

func (controller UserController) GetById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ctx := c.Request.Context()
	userFromUsecase, err := controller.UserUsecase.GetById(ctx, id)

	if err != nil {
		controllers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(c, fmt.Sprintf("user data with id %d fetched successfully", id), map[string]interface{}{
		"user": userFromUsecase,
	})
}

func (controller UserController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ctx := c.Request.Context()
	var payload request.UserRegisRequest
	err := c.Bind(&payload)

	if err := ValidateRequest(payload); err != nil {
		controllers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err != nil {
		controllers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	domainReq := payload.ToDomain()
	result, err := controller.UserUsecase.Update(ctx, &domainReq, id)

	if err != nil {
		controllers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(c, fmt.Sprintf("user data with id %d updated successfully", id), responses.FromDomain(result))
}

func (controller UserController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ctx := c.Request.Context()

	err := controller.UserUsecase.Delete(ctx, id)
	if err != nil {
		controllers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(c, fmt.Sprintf("book data with id %d deleted successfully", id), []int{})
}
