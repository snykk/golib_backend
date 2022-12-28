package reviews

import (
	"context"
	"time"

	"github.com/snykk/golib_backend/domains/books"
	"github.com/snykk/golib_backend/domains/users"
)

type Domain struct {
	ID        int
	Text      string
	Rating    int
	BookId    int
	Book      books.Domain
	UserId    int
	User      users.Domain
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Usecase interface {
	Store(ctx context.Context, book *Domain, userId int) (Domain, error)
	GetAll() ([]Domain, error)
	GetById(ctx context.Context, id int) (Domain, error)
	GetByBookId(ctx context.Context, bookId int) ([]Domain, error)
	GetByUserId(ctx context.Context, userId int) ([]Domain, error)
	Update(ctx context.Context, book *Domain, userId, reviewId int) (Domain, error)
	Delete(ctx context.Context, userId, reviewId int) (bookId int, err error)
}

type Repository interface {
	Store(ctx context.Context, book *Domain) (Domain, error)
	GetAll() ([]Domain, error)
	GetById(ctx context.Context, id int) (Domain, error)
	GetByBookId(ctx context.Context, bookId int) ([]Domain, error)
	GetByUserId(ctx context.Context, userId int) ([]Domain, error)
	Update(ctx context.Context, book *Domain) error
	Delete(ctx context.Context, domain *Domain) (bookId int, err error)
}
