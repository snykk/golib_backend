package request

import users "github.com/snykk/golib_backend/domains/users"

type UserChangePassRequest struct {
	Password    string `json:"password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (u *UserChangePassRequest) ToDomain() *users.Domain {
	return &users.Domain{
		Password: u.Password,
	}
}
