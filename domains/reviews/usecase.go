package reviews

import (
	"context"
	"errors"
)

type reviewUsecase struct {
	repo Repository
}

func NewReviewUsecase(repo Repository) Usecase {
	return &reviewUsecase{
		repo: repo,
	}
}

func (uc *reviewUsecase) Store(ctx context.Context, domain *Domain, userId int) (Domain, error) {
	domain.UserId = userId
	review, err := uc.repo.Store(ctx, domain)
	if err != nil {
		return review, err
	}
	return review, nil
}

func (uc *reviewUsecase) GetAll(ctx context.Context) ([]Domain, error) {
	domains, err := uc.repo.GetAll(ctx)

	if err != nil {
		return []Domain{}, err
	}

	return domains, nil
}

func (uc *reviewUsecase) GetById(ctx context.Context, id int) (Domain, error) {
	domain, err := uc.repo.GetById(ctx, id)

	if err != nil {
		return Domain{}, err
	}

	return domain, nil
}

func (uc *reviewUsecase) GetByBookId(ctx context.Context, bookId int) ([]Domain, error) {
	domain, err := uc.repo.GetByBookId(ctx, bookId)

	if err != nil {
		return []Domain{}, err
	}

	return domain, nil
}

func (uc *reviewUsecase) GetByUserId(ctx context.Context, userId int) ([]Domain, error) {
	domain, err := uc.repo.GetByUserId(ctx, userId)

	if err != nil {
		return []Domain{}, err
	}

	return domain, nil
}

func (uc *reviewUsecase) Update(ctx context.Context, domain *Domain, userId, reviewId int) (Domain, error) {
	beforeUpdate, err := uc.repo.GetById(ctx, reviewId)
	if err != nil {
		return Domain{}, errors.New("review not found")
	}

	if beforeUpdate.UserId != userId {
		return Domain{}, errors.New("you don't have access to update this review")
	}

	domain.ID = reviewId
	if err := uc.repo.Update(ctx, domain); err != nil {
		return Domain{}, err
	}

	afterUpdate, err := uc.repo.GetById(ctx, reviewId)
	if err != nil {
		return Domain{}, errors.New("internal server error")
	}

	return afterUpdate, err
}

func (uc *reviewUsecase) Delete(ctx context.Context, userId, reviewId int) (bookId int, err error) {
	beforeUpdate, err := uc.repo.GetById(ctx, reviewId)
	if err != nil {
		return 0, errors.New("review not found")
	}
	if beforeUpdate.UserId != userId {
		return 0, errors.New("you don't have access to delete this review")
	}

	return uc.repo.Delete(ctx, &beforeUpdate)
}

func (uc *reviewUsecase) GetUserReview(ctx context.Context, bookId, userId int) (Domain, error) {
	userReview, err := uc.repo.GetUserReview(ctx, bookId, userId)
	return userReview, err
}
