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
	Store(ctx context.Context, review *Domain, userId int) (domain Domain, statusCode int, err error)
	GetAll(ctx context.Context) (domains []Domain, statusCode int, err error)
	GetById(ctx context.Context, id int) (domain Domain, statusCode int, err error)
	GetByBookId(ctx context.Context, bookId int) (domains []Domain, statusCode int, err error)
	GetByUserId(ctx context.Context, userId int) (domains []Domain, statusCode int, err error)
	Update(ctx context.Context, review *Domain, userId, reviewId int) (domain Domain, statusCode int, err error)
	Delete(ctx context.Context, userId, reviewId int) (bookId int, statusCode int, err error)
	GetUserReview(ctx context.Context, bookId, userId int) (domain Domain, statusCode int, err error)
}

type Repository interface {
	Store(ctx context.Context, domain *Domain) (Domain, error)
	GetAll(ctx context.Context) ([]Domain, error)
	GetById(ctx context.Context, id int) (Domain, error)
	GetByBookId(ctx context.Context, bookId int) ([]Domain, error)
	GetByUserId(ctx context.Context, userId int) ([]Domain, error)
	Update(ctx context.Context, domain *Domain) error
	Delete(ctx context.Context, domain *Domain) (bookId int, err error)
	GetUserReview(ctx context.Context, bookId, userId int) (Domain, error)
}
