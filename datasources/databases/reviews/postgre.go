package reviews

import (
	"context"

	bookRepo "github.com/snykk/golib_backend/datasources/databases/books"
	"github.com/snykk/golib_backend/domains/reviews"
	"gorm.io/gorm"
)

type postgreReviewRepository struct {
	conn *gorm.DB
}

func NewPostgreReviewRepository(conn *gorm.DB) reviews.Repository {
	return &postgreReviewRepository{
		conn: conn,
	}
}

func (r *postgreReviewRepository) Store(ctx context.Context, domain *reviews.Domain) (reviews.Domain, error) {
	review := FromDomain(domain)

	err := r.conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&review).Error; err != nil {
			return err
		}

		var rating float64
		tx.Raw(`SELECT AVG("reviews".rating) FROM "reviews" WHERE book_id = ?`, review.BookId).Scan(&rating)

		if err := tx.Model(bookRepo.Book{}).Where(bookRepo.Book{Id: review.BookId}).Updates(bookRepo.Book{Rating: rating}).Error; err != nil {
			return err
		}

		if err := tx.Preload("User.Role").Preload("User.Gender").Preload("Book").First(&review, review.Id).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return reviews.Domain{}, err
	}

	return review.ToDomain(), nil
}

func (r *postgreReviewRepository) GetAll() ([]reviews.Domain, error) {
	var reviewRecords []Review
	if err := r.conn.Preload("User.Role").Preload("User.Gender").Preload("Book").Find(&reviewRecords).Error; err != nil {
		return []reviews.Domain{}, err
	}

	reviewDomains := ToArrayOfDomain(&reviewRecords)

	return reviewDomains, nil
}

func (r *postgreReviewRepository) GetById(ctx context.Context, id int) (reviews.Domain, error) {
	var review Review
	if err := r.conn.Preload("User.Role").Preload("User.Gender").Preload("Book").Where(Review{Id: id}).First(&review).Error; err != nil {
		return reviews.Domain{}, err
	}

	return review.ToDomain(), nil
}

func (r *postgreReviewRepository) Update(ctx context.Context, b *reviews.Domain) (err error) {
	review := FromDomain(b)

	err = r.conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Review{}).Where("id = ?", review.Id).Updates(&review).Error; err != nil {
			return err
		}

		var rating float64
		tx.Raw(`SELECT AVG("reviews".rating) FROM "reviews" WHERE book_id = ?`, review.BookId).Scan(&rating)

		if err := tx.Model(bookRepo.Book{}).Where(bookRepo.Book{Id: review.BookId}).Updates(bookRepo.Book{Rating: rating}).Error; err != nil {
			return err
		}

		return nil
	})

	return

}

func (r *postgreReviewRepository) Delete(ctx context.Context, id int) (err error) {
	err = r.conn.Delete(&Review{}, id).Error
	return
}
