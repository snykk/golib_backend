package users

import (
	"context"

	"github.com/snykk/golib_backend/domains/users"

	"gorm.io/gorm"
)

type postgreUserRepository struct {
	conn *gorm.DB
}

func NewPostgreUserRepository(conn *gorm.DB) users.Repository {
	return &postgreUserRepository{
		conn: conn,
	}
}

func (r *postgreUserRepository) Store(ctx context.Context, domain *users.Domain) (users.Domain, error) {
	var user = FromDomain(domain)

	err := r.conn.Transaction(func(tx *gorm.DB) error {
		if err := r.conn.Create(&user).Error; err != nil {
			return err
		}

		if err := r.conn.Preload("Role").Preload("Gender").Find(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return users.Domain{}, err
	}

	return user.ToDomain(), nil
}

func (r *postgreUserRepository) GetAll(ctx context.Context) ([]users.Domain, error) {
	var usersFromDB []User

	if err := r.conn.Preload("Role").Preload("Gender").Where("is_activated = true").Find(&usersFromDB).Error; err != nil {
		return []users.Domain{}, err
	}

	return ToArrayOfDomain(&usersFromDB), nil
}

func (r *postgreUserRepository) GetById(ctx context.Context, id int) (users.Domain, error) {
	var user User
	if err := r.conn.Preload("Role").Preload("Gender").Where("is_activated = true").First(&user, id).Error; err != nil {
		return users.Domain{}, err
	}

	return user.ToDomain(), nil
}

func (r *postgreUserRepository) Update(ctx context.Context, domain *users.Domain) (err error) {
	user := FromDomain(domain)
	err = r.conn.Model(&User{}).Model(&user).Updates(&user).Error
	return
}

func (r *postgreUserRepository) UpdateEmail(ctx context.Context, domain *users.Domain) (err error) {
	user := FromDomain(domain)
	err = r.conn.Model(&User{}).Model(&user).Updates(map[string]interface{}{"email": domain.Email, "is_activated": false}).Error
	return
}

func (r *postgreUserRepository) Delete(ctx context.Context, id int) (err error) {
	err = r.conn.Delete(&User{}, id).Error
	return
}

func (r *postgreUserRepository) GetByEmail(ctx context.Context, domain *users.Domain) (users.Domain, error) {
	var result User
	if err := r.conn.Preload("Role").Preload("Gender").First(&result, "email = ?", domain.Email).Error; err != nil {
		return users.Domain{}, err
	}

	return result.ToDomain(), nil
}
