package reviews

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type reviewUsecase struct {
	repo Repository
}

func NewReviewUsecase(repo Repository) Usecase {
	return &reviewUsecase{
		repo: repo,
	}
}

func (uc *reviewUsecase) Store(ctx context.Context, domain *Domain, userId int) (Domain, int, error) {
	domain.UserId = userId
	review, err := uc.repo.Store(ctx, domain)
	if err != nil {
		return review, http.StatusInternalServerError, err
	}
	return review, http.StatusCreated, nil
}

func (uc *reviewUsecase) GetAll(ctx context.Context) ([]Domain, int, error) {
	domains, err := uc.repo.GetAll(ctx)

	if err != nil {
		return []Domain{}, http.StatusInternalServerError, err
	}

	return domains, http.StatusOK, nil
}

func (uc *reviewUsecase) GetById(ctx context.Context, id int) (Domain, int, error) {
	domain, err := uc.repo.GetById(ctx, id)

	if err != nil {
		return Domain{}, http.StatusNotFound, errors.New("review not found")
	}

	return domain, http.StatusOK, nil
}

func (uc *reviewUsecase) GetByBookId(ctx context.Context, bookId int) ([]Domain, int, error) {
	domain, err := uc.repo.GetByBookId(ctx, bookId)

	if err != nil {
		return []Domain{}, http.StatusNotFound, fmt.Errorf("review with book id %d is not exists", bookId)
	}

	return domain, http.StatusOK, nil
}

func (uc *reviewUsecase) GetByUserId(ctx context.Context, userId int) ([]Domain, int, error) {
	domain, err := uc.repo.GetByUserId(ctx, userId)

	if err != nil {
		return []Domain{}, http.StatusNotFound, fmt.Errorf("review with book id %d is not exists", userId)
	}

	return domain, http.StatusOK, nil
}

func (uc *reviewUsecase) Update(ctx context.Context, domain *Domain, userId, reviewId int) (Domain, int, error) {
	beforeUpdate, err := uc.repo.GetById(ctx, reviewId)
	if err != nil {
		return Domain{}, http.StatusNotFound, errors.New("review not found")
	}

	if beforeUpdate.UserId != userId {
		return Domain{}, http.StatusUnauthorized, errors.New("you don't have access to update this review")
	}

	domain.ID = reviewId
	if err := uc.repo.Update(ctx, domain); err != nil {
		return Domain{}, http.StatusInternalServerError, err
	}

	afterUpdate, err := uc.repo.GetById(ctx, reviewId)
	if err != nil {
		return Domain{}, http.StatusNotFound, errors.New("review not found")
	}

	return afterUpdate, http.StatusOK, err
}

func (uc *reviewUsecase) Delete(ctx context.Context, userId, reviewId int) (bookId int, statusCode int, err error) {
	beforeUpdate, err := uc.repo.GetById(ctx, reviewId)
	if err != nil {
		return 0, http.StatusNotFound, errors.New("review not found")
	}
	if beforeUpdate.UserId != userId {
		return 0, http.StatusUnauthorized, errors.New("you don't have access to delete this review")
	}
	bookId, err = uc.repo.Delete(ctx, &beforeUpdate)

	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return bookId, http.StatusNoContent, err
}

func (uc *reviewUsecase) GetUserReview(ctx context.Context, bookId, userId int) (Domain, int, error) {
	userReview, err := uc.repo.GetUserReview(ctx, bookId, userId)
	if err != nil {
		return Domain{}, http.StatusNotFound, err
	}
	return userReview, http.StatusOK, err
}
