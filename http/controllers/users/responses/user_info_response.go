package responses

import (
	"time"

	"github.com/snykk/golib_backend/domains/users"
)

type UserInfoResponse struct {
	Id        int       `json:"id"`
	FullName  string    `json:"fullname"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Gender    string    `json:"gender"`
	Reviews   int       `json:"reviews"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *UserInfoResponse) ToDomain() users.Domain {
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

func FromDomainToUserInfo(u users.Domain) UserInfoResponse {
	return UserInfoResponse{
		Id:        u.ID,
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

func ToResponseUserinfoList(domains []users.Domain) []UserInfoResponse {
	var result []UserInfoResponse

	for _, val := range domains {
		result = append(result, FromDomainToUserInfo(val))
	}

	return result
}
