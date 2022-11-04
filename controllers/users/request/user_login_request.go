package request

import users "github.com/snykk/golib_backend/usecases/users"

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (u *UserLoginRequest) ToDomain() *users.Domain {
	return &users.Domain{
		Email:    u.Email,
		Password: u.Password,
	}
}
