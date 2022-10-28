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

func validateRequest(request interface{}) error {
	validator := validator.New()
	err := validator.Struct(request)
	return err
}

func (controller UserController) Regis(ctx *gin.Context) {
	var UserRegisRequest request.UserRegisRequest
	if err := ctx.ShouldBindJSON(&UserRegisRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateRequest(UserRegisRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	userDomain := UserRegisRequest.ToDomain()
	userDomainn, err := controller.UserUsecase.Store(ctxx, &userDomain)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, "registration user success", map[string]interface{}{
		"user": responses.FromDomain(userDomainn),
	})
}

func (controller UserController) Login(ctx *gin.Context) {
	var UserLoginRequest request.UserLoginRequest
	if err := ctx.ShouldBindJSON(&UserLoginRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateRequest(UserLoginRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	userDomain, err := controller.UserUsecase.Login(ctxx, UserLoginRequest.ToDomain())
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, "login success", responses.FromDomainLogin(userDomain))
}

func (controller UserController) GetAll(ctx *gin.Context) {
	usersFromUseCase, err := controller.UserUsecase.GetAll()

	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if usersFromUseCase == nil {
		controllers.NewSuccessResponse(ctx, "user data is empty", map[string]interface{}{
			"users": []int{},
		})
		return
	}

	controllers.NewSuccessResponse(ctx, "user data fetched successfully", map[string]interface{}{
		"users": responses.ToResponseList(&usersFromUseCase),
	})
}

func (controller UserController) GetById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	ctxx := ctx.Request.Context()
	userFromUsecase, err := controller.UserUsecase.GetById(ctxx, id)

	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("user data with id %d fetched successfully", id), map[string]interface{}{
		"user": userFromUsecase,
	})
}

func (controller UserController) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var userRequest request.UserRequest

	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userDomain := userRequest.ToDomain()
	ctxx := ctx.Request.Context()
	userDomainn, err := controller.UserUsecase.Update(ctxx, &userDomain, id)

	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("user data with id %d updated successfully", id), responses.FromDomain(userDomainn))
}

func (controller UserController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	ctxx := ctx.Request.Context()

	err := controller.UserUsecase.Delete(ctxx, id)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("user data with id %d deleted successfully", id), []int{})
}
