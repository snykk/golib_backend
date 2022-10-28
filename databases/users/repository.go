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

func (m UserRepository) Store(ctx context.Context, domain *users.Domain) (users.Domain, error) {
	var user = FromDomain(domain)
	err := m.Conn.Create(&user).Error

	if err != nil {
		return users.Domain{}, err
	}
	return user.ToDomain(), nil
}

func (m UserRepository) GetAll() ([]users.Domain, error) {
	var usersFromDB []User

	err := m.Conn.Find(&usersFromDB).Error

	if err != nil {
		return []users.Domain{}, err
	}

	return ToArrayOfDomain(&usersFromDB), nil
}

func (m UserRepository) GetById(ctx context.Context, id int) (users.Domain, error) {
	var user User
	if err := m.Conn.First(&user, id).Error; err != nil {
		return users.Domain{}, err
	}

	return user.ToDomain(), nil
}

func (m UserRepository) Update(ctx context.Context, domain *users.Domain) (users.Domain, error) {
	user := FromDomain(domain)
	if err := m.Conn.Save(&user).Error; err != nil {
		return users.Domain{}, err
	}

	return user.ToDomain(), nil
}

func (m UserRepository) Delete(ctx context.Context, id int) error {
	err := m.Conn.Delete(&User{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (m UserRepository) GetByEmail(ctx context.Context, domain *users.Domain) (users.Domain, error) {
	var result User
	if err := m.Conn.First(&result, "email = ?", domain.Email).Error; err != nil {
		return users.Domain{}, err
	}
	return result.ToDomain(), nil
}
