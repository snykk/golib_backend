package users

import (
	"time"

	"github.com/snykk/golib_backend/usecases/users"
	"gorm.io/gorm"
)

type User struct {
	Id        int
	Name      string
	Email     string
	Password  string
	IsAdmin   bool
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) ToDomain() users.Domain {
	return users.Domain{
		ID:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		IsAdmin:   u.IsAdmin,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func FromDomain(u *users.Domain) User {
	return User{
		Id:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		IsAdmin:   u.IsAdmin,
		IsActive:  u.IsActive,
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
