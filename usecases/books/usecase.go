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

func (bookUC *BookUsecase) GetAll() ([]Domain, error) {
	books, err := bookUC.Repo.GetAll()

	if err != nil {
		return []Domain{}, err
	}

	return books, nil
}

func (bookUC *BookUsecase) Store(ctx context.Context, book *Domain) (Domain, error) {
	result, err := bookUC.Repo.Store(ctx, book)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (bookUC *BookUsecase) GetById(ctx context.Context, id int) (Domain, error) {
	result, err := bookUC.Repo.GetById(ctx, id)

	if err != nil {
		return Domain{}, err
	}

	return result, nil
}

func (bookUC *BookUsecase) Update(ctx context.Context, book *Domain, id int) (Domain, error) {
	book.ID = id
	if err := bookUC.Repo.Update(ctx, book); err != nil {
		return Domain{}, err
	}

	newBook, err := bookUC.Repo.GetById(ctx, id)

	return newBook, err
}

func (bookUC *BookUsecase) Delete(ctx context.Context, id int) error {
	_, err := bookUC.Repo.GetById(ctx, id)
	if err != nil { // check wheter data is exists or not
		return err
	}
	err = bookUC.Repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
