package request

import users "github.com/snykk/golib_backend/usecase/users"

type UserRegisRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	IsAdmin  bool   `json:"is_admin"`
}

func (user UserRegisRequest) ToDomain() users.Domain {
	return users.Domain{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		IsAdmin:  user.IsAdmin,
	}
}
