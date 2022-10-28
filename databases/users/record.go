package users

import (
	"time"

	"github.com/snykk/golib_backend/usecase/users"
)

type User struct {
	Id        int
	Name      string
	Email     string
	Password  string
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) ToDomain() users.Domain {
	return users.Domain{
		Id:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func FromDomain(u *users.Domain) User {
	return User{
		Id:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func ToArrayOfDomain(u *[]User) []users.Domain {
	var result []users.Domain

	for _, val := range *u {
		result = append(result, val.ToDomain())
	}

	return result
}
