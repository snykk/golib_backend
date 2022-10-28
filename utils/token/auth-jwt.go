package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/snykk/golib_backend/config"
)

type JWTService interface {
	GenerateToken(userID int, isAdmin bool) (t string, err error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtCustomClaim struct {
	UserID  int
	IsAdmin bool
	jwt.StandardClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() JWTService {
	return &jwtService{
		issuer:    config.AppConfig.JWTIssuer,
		secretKey: getSecretKey(),
	}
}

func getSecretKey() string {
	secretKey := config.AppConfig.JWTSecret
	if secretKey == "" {
		secretKey = "this-is-not-secret-anymore-mwuehehe"
	}

	return secretKey
}

func (j *jwtService) GenerateToken(UserID int, isAdmin bool) (t string, err error) {
	claims := &jwtCustomClaim{
		UserID,
		isAdmin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(1, 0, 0).Unix(),
			Issuer:    j.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err = token.SignedString([]byte(j.secretKey))
	return
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}
