package token_test

import (
	"testing"
	"time"

	"github.com/snykk/golib_backend/config"
	"github.com/snykk/golib_backend/http/token"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	jwtService := token.NewJWTService()
	token, err := jwtService.GenerateToken(1, false, "john.doe@example.com", "password")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestParseToken(t *testing.T) {
	t.Run("With Valid Token", func(t *testing.T) {
		jwtService := token.NewJWTService()
		config.AppConfig.JWTExpired = 5

		token, _ := jwtService.GenerateToken(1, false, "john.doe@example.com", "password")

		claims, err := jwtService.ParseToken(token)
		assert.NoError(t, err)
		assert.Equal(t, 1, claims.UserID)
		assert.Equal(t, false, claims.IsAdmin)
		assert.Equal(t, "john.doe@example.com", claims.Email)
		assert.Equal(t, "password", claims.Password)
		assert.True(t, claims.StandardClaims.ExpiresAt > time.Now().Unix())
		assert.Equal(t, "john-doe", claims.StandardClaims.Issuer)
		assert.True(t, claims.StandardClaims.IssuedAt <= time.Now().Unix())
	})
	t.Run("With Invalid Token", func(t *testing.T) {
		jwtService := token.NewJWTService()

		_, err := jwtService.ParseToken("invalid_token")
		assert.Error(t, err)
		assert.Equal(t, "token is not valid", err.Error())
	})
}
