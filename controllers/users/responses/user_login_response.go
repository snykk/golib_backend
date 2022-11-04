package responses

import (
	"time"

	"github.com/snykk/golib_backend/usecases/users"
)

type UserLoginResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromDomainLogin(userDomain users.Domain) UserLoginResponse {
	return UserLoginResponse{
		Id:        userDomain.ID,
		Name:      userDomain.Name,
		Email:     userDomain.Email,
		IsAdmin:   userDomain.IsAdmin,
		Token:     userDomain.Token,
		CreatedAt: userDomain.CreatedAt,
		UpdatedAt: userDomain.UpdatedAt,
	}
}
