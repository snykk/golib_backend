package users

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/constants"
	"github.com/snykk/golib_backend/datasources/cache"
	"github.com/snykk/golib_backend/domains/users"
	"github.com/snykk/golib_backend/helpers"
	"github.com/snykk/golib_backend/http/controllers"
	"github.com/snykk/golib_backend/http/controllers/users/request"
	"github.com/snykk/golib_backend/http/controllers/users/responses"
	"github.com/snykk/golib_backend/http/token"
)

type UserController struct {
	usecase        users.Usecase
	redisCache     cache.RedisCache
	ristrettoCache cache.RistrettoCache
}

func NewUserController(usecase users.Usecase, redisCache cache.RedisCache, ristrettoCache cache.RistrettoCache) UserController {
	return UserController{
		usecase:        usecase,
		redisCache:     redisCache,
		ristrettoCache: ristrettoCache,
	}
}

func (c *UserController) Regis(ctx *gin.Context) {
	var UserRegisRequest request.UserRequest
	if err := ctx.ShouldBindJSON(&UserRegisRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := helpers.IsGenderValid(UserRegisRequest.Gender); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	userDomain := UserRegisRequest.ToDomain()
	userDomainn, statusCode, err := c.usecase.Store(ctxx, userDomain)
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, statusCode, "registration user success", map[string]interface{}{
		"user": responses.FromDomain(userDomainn),
	})
}

func (c *UserController) Login(ctx *gin.Context) {
	var UserLoginRequest request.UserLoginRequest
	if err := ctx.ShouldBindJSON(&UserLoginRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	userDomain, statusCode, err := c.usecase.Login(ctxx, UserLoginRequest.ToDomain())
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, statusCode, "login success", responses.FromDomain(userDomain))
}

func (c *UserController) GetAll(ctx *gin.Context) {
	if val := c.ristrettoCache.Get("users"); val != nil {
		controllers.NewSuccessResponse(ctx, http.StatusOK, "user data fetched successfully", map[string]interface{}{
			"users": val,
		})
		return
	}

	ctxx := ctx.Request.Context()
	usersFromUseCase, statusCode, err := c.usecase.GetAll(ctxx)

	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	if len(usersFromUseCase) == 0 {
		controllers.NewSuccessResponse(ctx, http.StatusOK, "user data is empty", map[string]interface{}{
			"users": []int{},
		})
		return
	}

	userResponses := responses.ToResponseUserinfoList(usersFromUseCase)

	go c.ristrettoCache.Set("users", userResponses)

	controllers.NewSuccessResponse(ctx, statusCode, "user data fetched successfully", map[string]interface{}{
		"users": userResponses,
	})
}

func (c *UserController) GetById(ctx *gin.Context) {
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(token.JwtCustomClaim)
	id, _ := strconv.Atoi(ctx.Param("id"))
	if val := c.ristrettoCache.Get(fmt.Sprintf("user/%d", id)); val != nil {
		controllers.NewSuccessResponse(ctx, http.StatusOK, fmt.Sprintf("user data with id %d fetched successfully", id), map[string]interface{}{
			"user": val,
		})
		return
	}

	ctxx := ctx.Request.Context()
	userFromUsecase, statusCode, err := c.usecase.GetById(ctxx, id, userClaims.UserID)

	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	userResponse := responses.FromDomain(userFromUsecase)

	go c.ristrettoCache.Set(fmt.Sprintf("user/%d", id), userResponse)

	controllers.NewSuccessResponse(ctx, statusCode, fmt.Sprintf("user data with id %d fetched successfully", id), map[string]interface{}{
		"user": userResponse,
	})
}

func (c *UserController) GetUserData(ctx *gin.Context) {
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(token.JwtCustomClaim)
	if val := c.ristrettoCache.Get(fmt.Sprintf("user/%s", userClaims.Email)); val != nil {
		controllers.NewSuccessResponse(ctx, http.StatusOK, "user data fetched successfully", map[string]interface{}{
			"user": val,
		})
		return
	}

	ctxx := ctx.Request.Context()
	userDom, statusCode, err := c.usecase.GetByEmail(ctxx, userClaims.Email)
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	userResponse := responses.FromDomain(userDom)

	go c.ristrettoCache.Set(fmt.Sprintf("user/%s", userClaims.Email), userResponse)

	controllers.NewSuccessResponse(ctx, statusCode, "user data fetched successfully", map[string]interface{}{
		"user": userResponse,
	})

}

func (c *UserController) Update(ctx *gin.Context) {
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(token.JwtCustomClaim)
	var userRequest request.UserUpdateRequest

	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := helpers.IsGenderValid(userRequest.Gender); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userDomain := userRequest.ToDomain()
	ctxx := ctx.Request.Context()
	userDomainn, statusCode, err := c.usecase.Update(ctxx, userDomain, userClaims.UserID)

	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	go c.ristrettoCache.Del("users", fmt.Sprintf("user/%d", userClaims.UserID))

	controllers.NewSuccessResponse(ctx, statusCode, fmt.Sprintf("user data with id %d updated successfully", userClaims.UserID), responses.FromDomain(userDomainn))
}

func (c *UserController) Delete(ctx *gin.Context) {
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(token.JwtCustomClaim)
	ctxx := ctx.Request.Context()

	statusCode, err := c.usecase.Delete(ctxx, userClaims.UserID)
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	go c.ristrettoCache.Del("users", fmt.Sprintf("user/%d", userClaims.UserID))

	controllers.NewSuccessResponse(ctx, statusCode, fmt.Sprintf("user data with id %d deleted successfully", userClaims.UserID), nil)
}

func (c *UserController) SendOTP(ctx *gin.Context) {
	var userOTP request.UserSendOTP

	if err := ctx.ShouldBindJSON(&userOTP); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	otpCode, statusCode, err := c.usecase.SendOTP(ctxx, userOTP.Email)
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	otpKey := fmt.Sprintf("user_otp:%s", userOTP.Email)
	go c.redisCache.Set(otpKey, otpCode)

	controllers.NewSuccessResponse(ctx, statusCode, fmt.Sprintf("otp code has been send to %s", userOTP.Email), nil)
}

func (c *UserController) VerifOTP(ctx *gin.Context) {
	var userOTP request.UserVerifOTP

	if err := ctx.ShouldBindJSON(&userOTP); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	otpKey := fmt.Sprintf("user_otp:%s", userOTP.Email)
	otpRedis := c.redisCache.Get(otpKey)

	statusCode, err := c.usecase.VerifOTP(ctxx, userOTP.Email, userOTP.Code, otpRedis)
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	statusCode, err = c.usecase.ActivateUser(ctxx, userOTP.Email)
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	go c.redisCache.Del(otpKey)
	go c.ristrettoCache.Del("users")

	controllers.NewSuccessResponse(ctx, statusCode, "otp verification success", nil)
}

func (c *UserController) ChangePassword(ctx *gin.Context) {
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(token.JwtCustomClaim)
	var userRequest request.UserChangePassRequest

	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userDomain := userRequest.ToDomain()
	ctxx := ctx.Request.Context()
	statusCode, err := c.usecase.ChangePassword(ctxx, userDomain, userRequest.NewPassword, userClaims.UserID)

	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	go c.ristrettoCache.Del("users", fmt.Sprintf("user/%d", userClaims.UserID))

	controllers.NewSuccessResponse(ctx, statusCode, "password has been changed", nil)
}

func (c *UserController) ChangeEmail(ctx *gin.Context) {
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(token.JwtCustomClaim)
	var userRequest request.UserChangeEmailRequest

	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userDomain := userRequest.ToDomain()
	ctxx := ctx.Request.Context()
	statusCode, err := c.usecase.ChangeEmail(ctxx, userDomain, userClaims.UserID)

	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	go c.ristrettoCache.Del("users", fmt.Sprintf("user/%d", userClaims.UserID))

	otpCode, statusCode, err := c.usecase.SendOTP(ctxx, userRequest.NewEmail)
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	otpKey := fmt.Sprintf("user_otp:%s", userDomain.Email)
	go c.redisCache.Set(otpKey, otpCode)

	controllers.NewSuccessResponse(ctx, statusCode, fmt.Sprintf("email has been changed and otp code has been send to %s please verify soon", userDomain.Email), nil)
}
