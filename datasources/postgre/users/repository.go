package users

import (
	"context"

	"github.com/snykk/golib_backend/domains/users"

	"gorm.io/gorm"
)

type userRepository struct {
	conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) users.Repository {
	return &userRepository{
		conn: conn,
	}
}

func (userR *userRepository) Store(ctx context.Context, domain *users.Domain) (users.Domain, error) {
	var user = FromDomain(domain)

	err := userR.conn.Transaction(func(tx *gorm.DB) error {
		if err := userR.conn.Create(&user).Error; err != nil {
			return err
		}

		if err := userR.conn.Preload("Role").Preload("Gender").Find(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return users.Domain{}, err
	}

	return user.ToDomain(), nil
}

func (userR *userRepository) GetAll() ([]users.Domain, error) {
	var usersFromDB []User

	if err := userR.conn.Preload("Role").Preload("Gender").Where("is_activated = true").Find(&usersFromDB).Error; err != nil {
		return []users.Domain{}, err
	}

	return ToArrayOfDomain(&usersFromDB), nil
}

func (userR *userRepository) GetById(ctx context.Context, id int) (users.Domain, error) {
	var user User
	if err := userR.conn.Preload("Role").Preload("Gender").Where("is_activated = true").First(&user, id).Error; err != nil {
		return users.Domain{}, err
	}

	return user.ToDomain(), nil
}

func (userR *userRepository) Update(ctx context.Context, domain *users.Domain) (err error) {
	user := FromDomain(domain)
	err = userR.conn.Model(&User{}).Model(&user).Updates(&user).Error
	return
}

func (userR *userRepository) UpdateEmail(ctx context.Context, domain *users.Domain) (err error) {
	user := FromDomain(domain)
	err = userR.conn.Model(&User{}).Model(&user).Updates(map[string]interface{}{"email": domain.Email, "is_activated": false}).Error
	return
}

func (userR *userRepository) Delete(ctx context.Context, id int) (err error) {
	err = userR.conn.Delete(&User{}, id).Error
	return
}

func (userR *userRepository) GetByEmail(ctx context.Context, domain *users.Domain) (users.Domain, error) {
	var result User
	if err := userR.conn.Preload("Role").Preload("Gender").First(&result, "email = ?", domain.Email).Error; err != nil {
		return users.Domain{}, err
	}

	return result.ToDomain(), nil
}
