package users

import (
	"context"

	"github.com/snykk/golib_backend/usecase/users"

	"gorm.io/gorm"
)

type UserRepository struct {
	Conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) users.Repository {
	return &UserRepository{
		conn,
	}
}

func (userRepo UserRepository) Store(ctx context.Context, domain *users.Domain) (users.Domain, error) {
	var user = FromDomain(domain)

	if err := userRepo.Conn.Create(&user).Error; err != nil {
		return users.Domain{}, err
	}

	return user.ToDomain(), nil
}

func (userRepo UserRepository) GetAll() ([]users.Domain, error) {
	var usersFromDB []User

	if err := userRepo.Conn.Find(&usersFromDB).Error; err != nil {
		return []users.Domain{}, err
	}

	return ToArrayOfDomain(&usersFromDB), nil
}

func (userRepo UserRepository) GetById(ctx context.Context, id int) (users.Domain, error) {
	var user User
	if err := userRepo.Conn.First(&user, id).Error; err != nil {
		return users.Domain{}, err
	}

	return user.ToDomain(), nil
}

func (userRepo UserRepository) Update(ctx context.Context, domain *users.Domain) (err error) {
	user := FromDomain(domain)
	err = userRepo.Conn.Model(&User{}).Model(&user).Updates(&user).Error
	return
}

func (userRepo UserRepository) Delete(ctx context.Context, id int) (err error) {
	err = userRepo.Conn.Delete(&User{}, id).Error
	return
}

func (userRepo UserRepository) GetByEmail(ctx context.Context, domain *users.Domain) (users.Domain, error) {
	var result User
	if err := userRepo.Conn.First(&result, "email = ?", domain.Email).Error; err != nil {
		return users.Domain{}, err
	}

	return result.ToDomain(), nil
}
