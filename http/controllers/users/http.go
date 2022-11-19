package users

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/datasources/cache"
	"github.com/snykk/golib_backend/domains/users"
	"github.com/snykk/golib_backend/http/controllers"
	"github.com/snykk/golib_backend/http/controllers/users/request"
	"github.com/snykk/golib_backend/http/controllers/users/responses"
	"github.com/snykk/golib_backend/utils/otp"
)

type UserController interface {
	Regis(ctx *gin.Context)
	Login(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetById(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	SendOTP(ctx *gin.Context)
	VerifOTP(ctx *gin.Context)
}

type userController struct {
	UserUsecase    users.Usecase
	RedisCache     cache.RedisCache
	RistrettoCache cache.RistrettoCache
}

func NewUserController(usecase users.Usecase, redisCache cache.RedisCache, ristrettoCache cache.RistrettoCache) UserController {
	return &userController{
		UserUsecase:    usecase,
		RedisCache:     redisCache,
		RistrettoCache: ristrettoCache,
	}
}

func (c userController) Regis(ctx *gin.Context) {
	var UserRegisRequest request.UserRegisRequest
	if err := ctx.ShouldBindJSON(&UserRegisRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	userDomain := UserRegisRequest.ToDomain()
	userDomainn, err := c.UserUsecase.Store(ctxx, &userDomain)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, "registration user success", map[string]interface{}{
		"user": responses.FromDomain(userDomainn),
	})
}

func (c userController) Login(ctx *gin.Context) {
	var UserLoginRequest request.UserLoginRequest
	if err := ctx.ShouldBindJSON(&UserLoginRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	userDomain, err := c.UserUsecase.Login(ctxx, UserLoginRequest.ToDomain())
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, "login success", responses.FromDomainLogin(userDomain))
}

func (c userController) GetAll(ctx *gin.Context) {
	if val := c.RistrettoCache.Get("users"); val != nil {
		controllers.NewSuccessResponse(ctx, "user data fetched successfully", map[string]interface{}{
			"users": val,
		})
		return
	}

	usersFromUseCase, err := c.UserUsecase.GetAll()

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

	go c.RistrettoCache.Set("users", userResponses)

	controllers.NewSuccessResponse(ctx, "user data fetched successfully", map[string]interface{}{
		"users": userResponses,
	})
}

func (c userController) GetById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if val := c.RistrettoCache.Get(fmt.Sprintf("user/%d", id)); val != nil {
		controllers.NewSuccessResponse(ctx, fmt.Sprintf("user data with id %d fetched successfully", id), map[string]interface{}{
			"user": val,
		})
		return
	}

	ctxx := ctx.Request.Context()
	authHeader := ctx.GetHeader("Authorization")
	userFromUsecase, err := c.UserUsecase.GetById(ctxx, id, authHeader)

	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	userResponse := responses.FromDomain(userFromUsecase)

	go c.RistrettoCache.Set(fmt.Sprintf("user/%d", id), userResponse)

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("user data with id %d fetched successfully", id), map[string]interface{}{
		"user": userResponse,
	})
}

func (c userController) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var userRequest request.UserRequest

	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userDomain := userRequest.ToDomain()
	ctxx := ctx.Request.Context()
	userDomainn, err := c.UserUsecase.Update(ctxx, &userDomain, id)

	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	go c.RistrettoCache.Del("users", fmt.Sprintf("user/%d", id))

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("user data with id %d updated successfully", id), responses.FromDomain(userDomainn))
}

func (c userController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	ctxx := ctx.Request.Context()

	err := c.UserUsecase.Delete(ctxx, id)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	go c.RistrettoCache.Del("users", fmt.Sprintf("user/%d", id))

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("user data with id %d deleted successfully", id), nil)
}

func (c userController) SendOTP(ctx *gin.Context) {
	var userOTP request.UserSendOTP

	if err := ctx.ShouldBindJSON(&userOTP); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	userDom, err := c.UserUsecase.GetByEmail(ctxx, userOTP.Email)
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
	go c.RedisCache.Set(otpKey, code)

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("otp code has been send to %s", userOTP.Email), nil)
}

func (c userController) VerifOTP(ctx *gin.Context) {
	var userOTP request.UserVerifOTP

	if err := ctx.ShouldBindJSON(&userOTP); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	userDom, err := c.UserUsecase.GetByEmail(ctxx, userOTP.Email)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if userDom.IsActive {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, "account is already activated")
		return
	}

	otpKey := fmt.Sprintf("user_otp:%s", userOTP.Email)
	otpCode := c.RedisCache.Get(otpKey)

	if otpCode != userOTP.Code {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, "invalid otp code")
		return
	}

	if err := c.UserUsecase.ActivateUser(ctxx, userOTP.Email); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	go c.RedisCache.Del(otpKey)

	controllers.NewSuccessResponse(ctx, "otp verification success", nil)
}
