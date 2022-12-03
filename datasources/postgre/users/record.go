package users

import (
	"time"

	"github.com/snykk/golib_backend/constants"
	"github.com/snykk/golib_backend/domains/users"
	"gorm.io/gorm"
)

type Role struct {
	Id   int    `gorm:"PrimaryKey "`
	Name string `gorm:"type:varchar(15) not null"`
}

type Gender struct {
	Id   int    `gorm:"PrimaryKey"`
	Name string `gorm:"type:varchar(15) not null"`
}

type User struct {
	Id          int    `gorm:"PrimaryKey"`
	FullName    string `gorm:"type:varchar(30) not null"`
	Username    string `gorm:"uniqueIndex:idx_username; type:varchar(30) not null"`
	Email       string `gorm:"uniqueIndex:idx_email; type:varchar(50) not null"`
	Password    string `gorm:"type:varchar(255) not null"`
	IsActivated bool   `gorm:"not null"`
	RoleId      int    `gorm:"not null"`
	Role        Role
	GenderId    int `gorm:"not null"`
	Gender      Gender
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (u *User) ToDomain() users.Domain {
	return users.Domain{
		ID:          u.Id,
		FullName:    u.FullName,
		Username:    u.Username,
		Email:       u.Email,
		Password:    u.Password,
		Role:        u.Role.Name,
		Gender:      u.Gender.Name,
		IsActivated: u.IsActivated,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func FromDomain(u *users.Domain) User {
	return User{
		Id:          u.ID,
		FullName:    u.FullName,
		Username:    u.Username,
		Email:       u.Email,
		Password:    u.Password,
		RoleId:      constants.MapperRoleToId[u.Role],
		GenderId:    constants.MapperGenderToId[u.Gender],
		IsActivated: u.IsActivated,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func ToArrayOfDomain(u *[]User) []users.Domain {
	var result []users.Domain

	for _, val := range *u {
		result = append(result, val.ToDomain())
	}

	return result
}
