package books

import (
	"context"
)

type BookUsecase struct {
	Repo Repository
}

func NewBookUsecase(repo Repository) Usecase {
	return &BookUsecase{
		repo,
	}
}

func (uc *BookUsecase) GetAll() ([]Domain, error) {
	books, err := uc.Repo.GetAll()

	if err != nil {
		return []Domain{}, err
	}

	return books, nil
}

func (uc *BookUsecase) Store(ctx context.Context, book *Domain) (Domain, error) {
	result, err := uc.Repo.Store(ctx, book)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (uc *BookUsecase) GetById(ctx context.Context, id int) (Domain, error) {
	result, err := uc.Repo.GetById(ctx, id)

	if err != nil {
		return Domain{}, err
	}

	return result, nil
}

func (uc *BookUsecase) Update(ctx context.Context, book *Domain) (Domain, error) {
	bookFromDB, err := uc.Repo.GetById(ctx, book.ID)
	if err != nil {
		return Domain{}, err
	}

	book.CreatedAt = bookFromDB.CreatedAt
	result, err := uc.Repo.Update(ctx, book)

	if err != nil {
		return Domain{}, err
	}

	return result, nil
}

func (uc *BookUsecase) Delete(ctx context.Context, id int) error {
	err := uc.Repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
