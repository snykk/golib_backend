package middlewares

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/utils/token"
)

func IsAdmin(jwtService token.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		token, err := token.GetToken(authHeader, jwtService)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if claims["IsAdmin"] == true {
			return
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": "you're not admin"})
	}
}
