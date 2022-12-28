package reviews

import (
	"context"
	"fmt"

	bookRepo "github.com/snykk/golib_backend/datasources/databases/books"
	userRepo "github.com/snykk/golib_backend/datasources/databases/users"
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

		// get rating of certain book
		var rating float64
		if err := tx.Raw(`SELECT AVG("reviews".rating) FROM "reviews" WHERE book_id = ? AND "deleted_at" IS NULL`, review.BookId).Scan(&rating).Error; err != nil {
			return err
		}

		// update book rating
		if err := tx.Model(bookRepo.Book{}).Where(bookRepo.Book{Id: review.BookId}).Updates(bookRepo.Book{Rating: &rating}).Error; err != nil {
			return err
		}

		// get user
		var users userRepo.User
		if err := tx.First(&users, review.UserId).Error; err != nil {
			return err
		}

		// update user rating
		if err := tx.Model(&users).Update("reviews", gorm.Expr("reviews + ?", 1)).Error; err != nil {
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

func (r *postgreReviewRepository) GetByBookId(ctx context.Context, bookId int) ([]reviews.Domain, error) {
	var review []Review
	if err := r.conn.Preload("User.Role").Preload("User.Gender").Preload("Book").Where(Review{BookId: bookId}).Find(&review).Error; err != nil {
		return []reviews.Domain{}, err
	}

	return ToArrayOfDomain(&review), nil
}

func (r *postgreReviewRepository) GetByUserId(ctx context.Context, userId int) ([]reviews.Domain, error) {
	var review []Review
	if err := r.conn.Preload("User.Role").Preload("User.Gender").Preload("Book").Where(Review{UserId: userId}).Find(&review).Error; err != nil {
		return []reviews.Domain{}, err
	}

	return ToArrayOfDomain(&review), nil
}

func (r *postgreReviewRepository) Update(ctx context.Context, b *reviews.Domain) (err error) {
	review := FromDomain(b)

	err = r.conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Review{}).Where("id = ?", review.Id).Updates(&review).Error; err != nil {
			return err
		}

		// get rating of certain book
		var rating float64
		if err := tx.Raw(`SELECT AVG("reviews".rating) FROM "reviews" WHERE book_id = ? AND "deleted_at" IS NULL`, review.BookId).Scan(&rating).Error; err != nil {
			return err
		}

		if err := tx.Model(bookRepo.Book{}).Where(bookRepo.Book{Id: review.BookId}).Updates(bookRepo.Book{Rating: &rating}).Error; err != nil {
			return err
		}

		return nil
	})

	return

}

func (r *postgreReviewRepository) Delete(ctx context.Context, domain *reviews.Domain) (bookId int, err error) {
	// err = r.conn.Delete(&Review{}, id).Error
	// return
	review := FromDomain(domain)

	err = r.conn.Transaction(func(tx *gorm.DB) error {
		fmt.Println("ini id review pake do", review.Id)
		if err = tx.Delete(&Review{}, review.Id).Error; err != nil {
			return err
		}

		// get rating of certain book
		var rating float64
		if err := tx.Raw(`SELECT COALESCE(AVG("reviews".rating), 0) FROM "reviews" WHERE book_id = ? AND "deleted_at" IS NULL`, review.BookId).Scan(&rating).Error; err != nil {
			return err
		}

		// update book rating
		if err := tx.Model(bookRepo.Book{}).Where(bookRepo.Book{Id: review.BookId}).Updates(bookRepo.Book{Rating: &rating}).Error; err != nil {
			return err
		}

		// get user
		var users userRepo.User
		if err := tx.First(&users, review.UserId).Error; err != nil {
			return err
		}

		// update user rating
		if err := tx.Model(&users).Update("reviews", gorm.Expr("reviews - ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})

	return review.BookId, err
}
