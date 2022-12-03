package request

import users "github.com/snykk/golib_backend/domains/users"

type UserChangeEmailRequest struct {
	NewEmail string `json:"new_email" binding:"required"`
}

func (u *UserChangeEmailRequest) ToDomain() *users.Domain {
	return &users.Domain{
		Email: u.NewEmail,
	}
}
