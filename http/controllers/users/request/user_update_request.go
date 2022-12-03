package request

import (
	users "github.com/snykk/golib_backend/domains/users"
)

type UserUpdateRequest struct {
	FullName string `json:"fullname" binding:"required"`
	Username string `json:"username" binding:"required"`
	Gender   string `json:"gender" binding:"required"`
}

func (user UserUpdateRequest) ToDomain() *users.Domain {
	return &users.Domain{
		FullName: user.FullName,
		Username: user.Username,
		Gender:   user.Gender,
	}
}
