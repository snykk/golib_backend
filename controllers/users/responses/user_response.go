package responses

import (
	"time"

	"github.com/snykk/golib_backend/usecase/users"
)

type UserResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	Password  string    `json:"password,omitempty"`
	Token     string    `json:"token,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *UserResponse) ToDomain() users.Domain {
	return users.Domain{
		Id:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func FromDomain(u users.Domain) UserResponse {
	return UserResponse{
		Id:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		IsAdmin:   u.IsAdmin,
		Password:  u.Password,
		Token:     u.Token,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func ToResponseList(domains *[]users.Domain) []UserResponse {
	var result []UserResponse

	for _, val := range *domains {
		result = append(result, FromDomain(val))
	}

	return result
}
