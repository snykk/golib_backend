package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/utils/token"
)

func AuthorizeJWT(jwtService token.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if _, err := jwtService.ParseToken(authHeader); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		}
	}
}
