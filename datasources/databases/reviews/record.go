package reviews

import (
	"time"

	"github.com/snykk/golib_backend/datasources/databases/books"
	"github.com/snykk/golib_backend/datasources/databases/users"
	"github.com/snykk/golib_backend/domains/reviews"
	"gorm.io/gorm"
)

type Review struct {
	Id        int    `gorm:"primaryKey;autoIncrement"`
	Text      string `gorm:"type:text; not null"`
	Rating    int    `gorm:"type:integer; not null"`
	BookId    int    `gorm:"not null"`
	Book      books.Book
	UserId    int `gorm:"not null"`
	User      users.User
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *Review) ToDomain() reviews.Domain {
	return reviews.Domain{
		ID:        u.Id,
		Text:      u.Text,
		Rating:    u.Rating,
		BookId:    u.BookId,
		Book:      u.Book.ToDomain(),
		UserId:    u.UserId,
		User:      u.User.ToDomain(),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func FromDomain(domain *reviews.Domain) Review {
	return Review{
		Id:        domain.ID,
		Text:      domain.Text,
		Rating:    domain.Rating,
		BookId:    domain.BookId,
		UserId:    domain.UserId,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}
}

func ToArrayOfDomain(reviewss *[]Review) []reviews.Domain {
	var result []reviews.Domain

	for _, review := range *reviewss {
		result = append(result, review.ToDomain())
	}

	return result
}
