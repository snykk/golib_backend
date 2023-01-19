package books

import (
	"context"
)

type bookUsecase struct {
	repo Repository
}

func NewBookUsecase(repo Repository) Usecase {
	return &bookUsecase{
		repo,
	}
}

func (uc *bookUsecase) GetAll(ctx context.Context) ([]Domain, error) {
	books, err := uc.repo.GetAll(ctx)

	if err != nil {
		return []Domain{}, err
	}

	return books, nil
}

func (uc *bookUsecase) Store(ctx context.Context, book *Domain) (Domain, error) {
	result, err := uc.repo.Store(ctx, book)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (uc *bookUsecase) GetById(ctx context.Context, id int) (Domain, error) {
	result, err := uc.repo.GetById(ctx, id)

	if err != nil {
		return Domain{}, err
	}

	return result, nil
}

func (uc *bookUsecase) Update(ctx context.Context, book *Domain, id int) (Domain, error) {
	book.ID = id
	if err := uc.repo.Update(ctx, book); err != nil {
		return Domain{}, err
	}

	newBook, err := uc.repo.GetById(ctx, id)

	return newBook, err
}

func (uc *bookUsecase) Delete(ctx context.Context, id int) error {
	_, err := uc.repo.GetById(ctx, id)
	if err != nil { // check wheter data is exists or not
		return err
	}
	err = uc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
