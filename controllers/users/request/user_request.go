package request

import users "github.com/snykk/golib_backend/usecase/users"

type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

func (user UserRequest) ToDomain() users.Domain {
	return users.Domain{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		IsAdmin:  user.IsAdmin,
	}
}
