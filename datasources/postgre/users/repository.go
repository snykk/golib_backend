package users

import (
	"context"

	"github.com/snykk/golib_backend/domains/users"

	"gorm.io/gorm"
)

type userRepository struct {
	Conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) users.Repository {
	return &userRepository{
		conn,
	}
}

func (userR userRepository) Store(ctx context.Context, domain *users.Domain) (users.Domain, error) {
	var user = FromDomain(domain)

	if err := userR.Conn.Create(&user).Error; err != nil {
		return users.Domain{}, err
	}

	return user.ToDomain(), nil
}

func (userR userRepository) GetAll() ([]users.Domain, error) {
	var usersFromDB []User

	if err := userR.Conn.Find(&usersFromDB).Error; err != nil {
		return []users.Domain{}, err
	}

	return ToArrayOfDomain(&usersFromDB), nil
}

func (userR userRepository) GetById(ctx context.Context, id int) (users.Domain, error) {
	var user User
	if err := userR.Conn.First(&user, id).Error; err != nil {
		return users.Domain{}, err
	}

	return user.ToDomain(), nil
}

func (userR userRepository) Update(ctx context.Context, domain *users.Domain) (err error) {
	user := FromDomain(domain)
	err = userR.Conn.Model(&User{}).Model(&user).Updates(&user).Error
	return
}

func (userR userRepository) Delete(ctx context.Context, id int) (err error) {
	err = userR.Conn.Delete(&User{}, id).Error
	return
}

func (userR userRepository) GetByEmail(ctx context.Context, domain *users.Domain) (users.Domain, error) {
	var result User
	if err := userR.Conn.First(&result, "email = ?", domain.Email).Error; err != nil {
		return users.Domain{}, err
	}

	return result.ToDomain(), nil
}
