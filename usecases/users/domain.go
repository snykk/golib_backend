package users

import (
	"context"
	"time"
)

type Domain struct {
	ID        int
	Name      string
	Email     string
	Password  string
	IsAdmin   bool
	Token     string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Usecase interface {
	Store(ctx context.Context, domain *Domain) (Domain, error)
	GetAll() ([]Domain, error)
	GetById(ctx context.Context, id int, authHeader string) (Domain, error)
	Update(ctx context.Context, domain *Domain, id int) (Domain, error)
	Delete(ctx context.Context, id int) error
	Login(ctx context.Context, domain *Domain) (Domain, error)
	ActivateUser(ctx context.Context, email string) (err error)
	GetByEmail(ctx context.Context, email string) (Domain, error)
}

type Repository interface {
	Store(ctx context.Context, domain *Domain) (Domain, error)
	GetAll() ([]Domain, error)
	GetById(ctx context.Context, id int) (Domain, error)
	Update(ctx context.Context, domain *Domain) (err error)
	Delete(ctx context.Context, id int) (err error)
	GetByEmail(ctx context.Context, domain *Domain) (Domain, error)
}
