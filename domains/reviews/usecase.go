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

func (uc *reviewUsecase) GetAll() ([]Domain, error) {
	domains, err := uc.repo.GetAll()

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

func (uc *reviewUsecase) Delete(ctx context.Context, userId, reviewId int) (err error) {
	beforeUpdate, err := uc.repo.GetById(ctx, reviewId)
	if err != nil {
		return errors.New("review not found")
	}
	if beforeUpdate.UserId != userId {
		return errors.New("you don't have access to update this review")
	}

	return uc.repo.Delete(ctx, reviewId)
}
