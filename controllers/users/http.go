package users

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/controllers"
	"github.com/snykk/golib_backend/controllers/users/request"
	"github.com/snykk/golib_backend/controllers/users/responses"
	"github.com/snykk/golib_backend/datasources/cache"
	"github.com/snykk/golib_backend/usecases/users"
	"github.com/snykk/golib_backend/utils/otp"

	validator "github.com/go-playground/validator/v10"
)

type UserController struct {
	UserUsecase    users.Usecase
	RedisCache     cache.RedisCache
	RistrettoCache cache.RistrettoCache
}

func NewUserController(usecase users.Usecase, redisCache cache.RedisCache, ristrettoCache cache.RistrettoCache) *UserController {
	return &UserController{
		UserUsecase:    usecase,
		RedisCache:     redisCache,
		RistrettoCache: ristrettoCache,
	}
}

func validateRequest(request interface{}) error {
	validator := validator.New()
	err := validator.Struct(request)
	return err
}

func (userController UserController) Regis(ctx *gin.Context) {
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
	userDomainn, err := userController.UserUsecase.Store(ctxx, &userDomain)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, "registration user success", map[string]interface{}{
		"user": responses.FromDomain(userDomainn),
	})
}

func (userController UserController) Login(ctx *gin.Context) {
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
	userDomain, err := userController.UserUsecase.Login(ctxx, UserLoginRequest.ToDomain())
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, "login success", responses.FromDomainLogin(userDomain))
}

func (userController UserController) GetAll(ctx *gin.Context) {
	if val := userController.RistrettoCache.Get("users"); val != nil {
		controllers.NewSuccessResponse(ctx, "user data fetched successfully", map[string]interface{}{
			"users": val,
		})
		return
	}

	usersFromUseCase, err := userController.UserUsecase.GetAll()

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

	userResponses := responses.ToResponseList(usersFromUseCase)

	userController.RistrettoCache.Set("users", userResponses)

	controllers.NewSuccessResponse(ctx, "user data fetched successfully", map[string]interface{}{
		"users": userResponses,
	})
}

func (userController UserController) GetById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if val := userController.RistrettoCache.Get(fmt.Sprintf("user/%d", id)); val != nil {
		controllers.NewSuccessResponse(ctx, fmt.Sprintf("user data with id %d fetched successfully", id), map[string]interface{}{
			"user": val,
		})
		return
	}

	ctxx := ctx.Request.Context()
	authHeader := ctx.GetHeader("Authorization")
	userFromUsecase, err := userController.UserUsecase.GetById(ctxx, id, authHeader)

	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	userResponse := responses.FromDomain(userFromUsecase)

	userController.RistrettoCache.Set(fmt.Sprintf("user/%d", id), userResponse)

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("user data with id %d fetched successfully", id), map[string]interface{}{
		"user": userResponse,
	})
}

func (userController UserController) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var userRequest request.UserRequest

	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userDomain := userRequest.ToDomain()
	ctxx := ctx.Request.Context()
	userDomainn, err := userController.UserUsecase.Update(ctxx, &userDomain, id)

	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	userController.RistrettoCache.Del("users", fmt.Sprintf("user/%d", id))

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("user data with id %d updated successfully", id), responses.FromDomain(userDomainn))
}

func (userController UserController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	ctxx := ctx.Request.Context()

	err := userController.UserUsecase.Delete(ctxx, id)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	userController.RistrettoCache.Del("users", fmt.Sprintf("user/%d", id))

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("user data with id %d deleted successfully", id), nil)
}

func (userController UserController) SendOTP(ctx *gin.Context) {
	var userOTP request.UserSendOTP

	if err := ctx.ShouldBindJSON(&userOTP); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateRequest(userOTP); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	userDom, err := userController.UserUsecase.GetByEmail(ctxx, userOTP.Email)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if userDom.IsActive {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, "account alreaday activated")
		return
	}

	code, err := otp.GenerateCode(6)
	if err != nil {
		log.Println(err)
	}

	if err = otp.SendOTP(code, userOTP.Email); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	otpKey := fmt.Sprintf("user_otp:%s", userOTP.Email)
	userController.RedisCache.Set(otpKey, code)

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("otp code has been send to %s", userOTP.Email), nil)
}

func (userController UserController) VerifOTP(ctx *gin.Context) {
	var userOTP request.UserVerifOTP

	if err := ctx.ShouldBindJSON(&userOTP); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateRequest(userOTP); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	userDom, err := userController.UserUsecase.GetByEmail(ctxx, userOTP.Email)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if userDom.IsActive {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, "account is already activated")
		return
	}

	otpKey := fmt.Sprintf("user_otp:%s", userOTP.Email)
	otpCode := userController.RedisCache.Get(otpKey)

	if otpCode != userOTP.Code {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, "invalid otp code")
		return
	}

	if err := userController.UserUsecase.ActivateUser(ctxx, userOTP.Email); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	userController.RedisCache.Del(otpKey)

	controllers.NewSuccessResponse(ctx, "otp verification success", nil)
}
