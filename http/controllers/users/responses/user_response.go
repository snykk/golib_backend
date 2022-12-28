package responses

import (
	"time"

	"github.com/snykk/golib_backend/domains/users"
)

type UserResponse struct {
	Id        int       `json:"id"`
	FullName  string    `json:"fullname"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Gender    string    `json:"gender"`
	Password  string    `json:"password,omitempty"`
	Reviews   int       `json:"reviews"`
	Token     string    `json:"token,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *UserResponse) ToDomain() users.Domain {
	return users.Domain{
		ID:        u.Id,
		FullName:  u.FullName,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		Gender:    u.Gender,
		Reviews:   u.Reviews,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func FromDomain(u users.Domain) UserResponse {
	return UserResponse{
		Id:        u.ID,
		FullName:  u.FullName,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		Gender:    u.Gender,
		Password:  u.Password,
		Reviews:   u.Reviews,
		Token:     u.Token,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func ToResponseList(domains []users.Domain) []UserResponse {
	var result []UserResponse

	for _, val := range domains {
		result = append(result, FromDomain(val))
	}

	return result
}
