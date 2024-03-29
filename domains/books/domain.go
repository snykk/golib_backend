package books

import (
	"context"
	"time"
)

type Domain struct {
	ID          int
	Title       string
	Description string
	Author      string
	Publisher   string
	ISBN        string
	Rating      *float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Usecase interface {
	GetAll(ctx context.Context) (domains []Domain, statusCode int, err error)
	Store(ctx context.Context, book *Domain) (domain Domain, statusCode int, err error)
	GetById(ctx context.Context, id int) (domain Domain, statusCode int, err error)
	Update(ctx context.Context, book *Domain, id int) (domain Domain, statusCode int, err error)
	Delete(ctx context.Context, id int) (statusCode int, err error)
}

type Repository interface {
	GetAll(ctx context.Context) ([]Domain, error)
	Store(ctx context.Context, book *Domain) (Domain, error)
	GetById(ctx context.Context, id int) (Domain, error)
	Update(ctx context.Context, book *Domain) (err error)
	Delete(ctx context.Context, id int) error
}
