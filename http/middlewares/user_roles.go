package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/utils/token"
)

func IsAdmin(jwtService token.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		claims, err := jwtService.ParseToken(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		if claims.IsAdmin {
			return
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": "you're not admin"})
	}
}

func IsValidUser(jwtService token.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		idParam, _ := strconv.Atoi(c.Param("id"))

		claims, err := jwtService.ParseToken(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
			return
		}

		if claims.IsAdmin || idParam == claims.UserID {
			return
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": "you don't have access for this action"})
	}
}
