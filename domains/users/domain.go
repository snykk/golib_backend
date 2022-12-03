package users

import (
	"context"
	"time"
)

type Domain struct {
	ID          int
	FullName    string
	Username    string
	Email       string
	Password    string
	Token       string
	Role        string
	Gender      string
	IsActivated bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Usecase interface {
	Store(ctx context.Context, domain *Domain) (Domain, error)
	GetAll() ([]Domain, error)
	GetById(ctx context.Context, id int, idClaims int) (Domain, error)
	Update(ctx context.Context, domain *Domain, id int) (Domain, error)
	Delete(ctx context.Context, id int) error
	Login(ctx context.Context, domain *Domain) (Domain, error)
	ActivateUser(ctx context.Context, email string) (err error)
	GetByEmail(ctx context.Context, email string) (Domain, error)
	ChangePassword(ctx context.Context, domain *Domain, new_pass string, id int) (err error)
	ChangeEmail(ctx context.Context, domain *Domain, id int) (err error)
}

type Repository interface {
	Store(ctx context.Context, domain *Domain) (Domain, error)
	GetAll() ([]Domain, error)
	GetById(ctx context.Context, id int) (Domain, error)
	Update(ctx context.Context, domain *Domain) (err error)
	Delete(ctx context.Context, id int) (err error)
	GetByEmail(ctx context.Context, domain *Domain) (Domain, error)
	UpdateEmail(ctx context.Context, domain *Domain) (err error)
}
