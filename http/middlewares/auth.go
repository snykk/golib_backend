package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/constants"
	"github.com/snykk/golib_backend/http/controllers"
	"github.com/snykk/golib_backend/utils/token"
)

type AuthMiddleware struct {
	jwtService token.JWTService
	isAdmin    bool
}

func NewAuthMiddleware(jwtService token.JWTService, isAdmin bool) gin.HandlerFunc {
	return (&AuthMiddleware{
		jwtService: jwtService,
		isAdmin:    isAdmin,
	}).Handle
}

func (m *AuthMiddleware) Handle(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		controllers.NewAbortResponse(ctx, "missing authorization header")
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		controllers.NewAbortResponse(ctx, "invalid header format")
		return
	}

	if headerParts[0] != "Bearer" {
		controllers.NewAbortResponse(ctx, "token must content bearer")
		return
	}

	user, err := m.jwtService.ParseToken(headerParts[1])
	if err != nil {
		controllers.NewAbortResponse(ctx, "invalid token")
		return
	}

	if user.IsAdmin != m.isAdmin && !user.IsAdmin {
		controllers.NewAbortResponse(ctx, "you don't have access for this action")
		return
	}

	ctx.Set(constants.CtxAuthenticatedUserKey, user)
	ctx.Next()
}
