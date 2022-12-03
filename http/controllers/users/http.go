package users

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/constants"
	"github.com/snykk/golib_backend/datasources/cache"
	"github.com/snykk/golib_backend/domains/users"
	"github.com/snykk/golib_backend/http/controllers"
	"github.com/snykk/golib_backend/http/controllers/users/request"
	"github.com/snykk/golib_backend/http/controllers/users/responses"
	"github.com/snykk/golib_backend/utils/otp"
	"github.com/snykk/golib_backend/utils/token"
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

	if err := isGenderValid(UserRegisRequest.Gender); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	userDomain := UserRegisRequest.ToDomain()
	userDomainn, err := c.usecase.Store(ctxx, &userDomain)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, "registration user success", map[string]interface{}{
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
	userDomain, err := c.usecase.Login(ctxx, UserLoginRequest.ToDomain())
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, "login success", responses.FromDomain(userDomain))
}

func (c *UserController) GetAll(ctx *gin.Context) {
	if val := c.ristrettoCache.Get("users"); val != nil {
		controllers.NewSuccessResponse(ctx, "user data fetched successfully", map[string]interface{}{
			"users": val,
		})
		return
	}

	usersFromUseCase, err := c.usecase.GetAll()

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

	go c.ristrettoCache.Set("users", userResponses)

	controllers.NewSuccessResponse(ctx, "user data fetched successfully", map[string]interface{}{
		"users": userResponses,
	})
}

func (c *UserController) GetById(ctx *gin.Context) {
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(token.JwtCustomClaim)
	id, _ := strconv.Atoi(ctx.Param("id"))
	if val := c.ristrettoCache.Get(fmt.Sprintf("user/%d", id)); val != nil {
		controllers.NewSuccessResponse(ctx, fmt.Sprintf("user data with id %d fetched successfully", id), map[string]interface{}{
			"user": val,
		})
		return
	}

	ctxx := ctx.Request.Context()
	userFromUsecase, err := c.usecase.GetById(ctxx, id, userClaims.UserID)

	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	userResponse := responses.FromDomain(userFromUsecase)

	go c.ristrettoCache.Set(fmt.Sprintf("user/%d", id), userResponse)

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("user data with id %d fetched successfully", id), map[string]interface{}{
		"user": userResponse,
	})
}

func (c *UserController) GetUserData(ctx *gin.Context) {
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(token.JwtCustomClaim)
	if val := c.ristrettoCache.Get(fmt.Sprintf("user/%s", userClaims.Password)); val != nil {
		fmt.Println("pake redis")
		controllers.NewSuccessResponse(ctx, "user data fetched successfully", map[string]interface{}{
			"user": val,
		})
		return
	}

	ctxx := ctx.Request.Context()
	userDom, err := c.usecase.GetByEmail(ctxx, userClaims.Email)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	userResponse := responses.FromDomain(userDom)

	go c.ristrettoCache.Set(fmt.Sprintf("user/%s", userClaims.Password), userResponse)

	fmt.Println("dari db")
	controllers.NewSuccessResponse(ctx, "user data fetched successfully", map[string]interface{}{
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

	if err := isGenderValid(userRequest.Gender); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userDomain := userRequest.ToDomain()
	ctxx := ctx.Request.Context()
	userDomainn, err := c.usecase.Update(ctxx, userDomain, userClaims.UserID)

	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	go c.ristrettoCache.Del("users", fmt.Sprintf("user/%d", userClaims.UserID))

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("user data with id %d updated successfully", userClaims.UserID), responses.FromDomain(userDomainn))
}

func (c *UserController) Delete(ctx *gin.Context) {
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(token.JwtCustomClaim)
	ctxx := ctx.Request.Context()

	err := c.usecase.Delete(ctxx, userClaims.UserID)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	go c.ristrettoCache.Del("users", fmt.Sprintf("user/%d", userClaims.UserID))

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("user data with id %d deleted successfully", userClaims.UserID), nil)
}

func (c *UserController) SendOTP(ctx *gin.Context) {
	var userOTP request.UserSendOTP

	if err := ctx.ShouldBindJSON(&userOTP); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	userDom, err := c.usecase.GetByEmail(ctxx, userOTP.Email)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if userDom.IsActivated {
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
	go c.redisCache.Set(otpKey, code)

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("otp code has been send to %s", userOTP.Email), nil)
}

func (c *UserController) VerifOTP(ctx *gin.Context) {
	var userOTP request.UserVerifOTP

	if err := ctx.ShouldBindJSON(&userOTP); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	userDom, err := c.usecase.GetByEmail(ctxx, userOTP.Email)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if userDom.IsActivated {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, "account is already activated")
		return
	}

	otpKey := fmt.Sprintf("user_otp:%s", userOTP.Email)
	otpCode := c.redisCache.Get(otpKey)

	if otpCode != userOTP.Code {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, "invalid otp code")
		return
	}

	if err := c.usecase.ActivateUser(ctxx, userOTP.Email); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	go c.redisCache.Del(otpKey)

	controllers.NewSuccessResponse(ctx, "otp verification success", nil)
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
	err := c.usecase.ChangePassword(ctxx, userDomain, userRequest.NewPassword, userClaims.UserID)

	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	go c.ristrettoCache.Del("users", fmt.Sprintf("user/%d", userClaims.UserID))

	controllers.NewSuccessResponse(ctx, "password has been changed", nil)
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
	err := c.usecase.ChangeEmail(ctxx, userDomain, userClaims.UserID)

	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	go c.ristrettoCache.Del("users", fmt.Sprintf("user/%d", userClaims.UserID))

	code, err := otp.GenerateCode(6)
	if err != nil {
		log.Println(err)
	}

	if err = otp.SendOTP(code, userDomain.Email); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	otpKey := fmt.Sprintf("user_otp:%s", userDomain.Email)
	go c.redisCache.Set(otpKey, code)

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("email has been changed and otp code has been send to %s please verify soon", userDomain.Email), nil)
}

func isGenderValid(gender string) error {
	if !isArrayContains(constants.ListGender, gender) {
		var option string
		for index, g := range constants.ListGender {
			option += g
			if index != len(constants.ListGender)-1 {
				option += ", "
			}
		}

		return fmt.Errorf("gender must be one of [%s]", option)
	}

	return nil
}

func isArrayContains(arr []string, str string) bool {
	for _, item := range arr {
		if item == str {
			return true
		}
	}
	return false
}
