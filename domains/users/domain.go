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
	Reviews     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Usecase interface {
	Store(ctx context.Context, user *Domain) (domain Domain, statusCode int, err error)
	GetAll(ctx context.Context) (domains []Domain, statusCode int, err error)
	GetById(ctx context.Context, id int, idClaims int) (domain Domain, statusCode int, err error)
	Update(ctx context.Context, user *Domain, id int) (domain Domain, statusCode int, err error)
	Delete(ctx context.Context, id int) (statusCode int, err error)
	Login(ctx context.Context, user *Domain) (domain Domain, statusCode int, err error)
	ActivateUser(ctx context.Context, email string) (statusCode int, err error)
	GetByEmail(ctx context.Context, email string) (domain Domain, statusCode int, err error)
	ChangePassword(ctx context.Context, domain *Domain, new_pass string, id int) (statusCode int, err error)
	ChangeEmail(ctx context.Context, domain *Domain, id int) (statusCode int, err error)
	SendOTP(ctx context.Context, email string) (otpCode string, statusCode int, err error)
	VerifOTP(ctx context.Context, email string, userOTP string, otpRedis string) (statusCode int, err error)
}

type Repository interface {
	Store(ctx context.Context, domain *Domain) (Domain, error)
	GetAll(ctx context.Context) ([]Domain, error)
	GetById(ctx context.Context, id int) (Domain, error)
	Update(ctx context.Context, domain *Domain) (err error)
	Delete(ctx context.Context, id int) (err error)
	GetByEmail(ctx context.Context, domain *Domain) (Domain, error)
	UpdateEmail(ctx context.Context, domain *Domain) (err error)
}
