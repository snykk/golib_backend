package request

import (
	"github.com/snykk/golib_backend/constants"
	users "github.com/snykk/golib_backend/domains/users"
)

type UserRequest struct {
	FullName string `json:"fullname" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Gender   string `json:"gender" binding:"required"`
}

func (user UserRequest) ToDomain() *users.Domain {
	return &users.Domain{
		FullName: user.FullName,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Gender:   user.Gender,
		Role:     constants.User,
	}
}
