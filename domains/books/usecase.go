package books

import (
	"context"
	"errors"
	"net/http"
)

type bookUsecase struct {
	repo Repository
}

func NewBookUsecase(repo Repository) Usecase {
	return &bookUsecase{
		repo,
	}
}

func (uc *bookUsecase) GetAll(ctx context.Context) ([]Domain, int, error) {
	books, err := uc.repo.GetAll(ctx)

	if err != nil {
		return []Domain{}, http.StatusInternalServerError, err
	}

	return books, http.StatusOK, nil
}

func (uc *bookUsecase) Store(ctx context.Context, book *Domain) (Domain, int, error) {
	result, err := uc.repo.Store(ctx, book)
	if err != nil {
		return result, http.StatusInternalServerError, err
	}
	return result, http.StatusCreated, nil
}

func (uc *bookUsecase) GetById(ctx context.Context, id int) (Domain, int, error) {
	result, err := uc.repo.GetById(ctx, id)

	if err != nil {
		return Domain{}, http.StatusNotFound, errors.New("book not found")
	}

	return result, http.StatusOK, nil
}

func (uc *bookUsecase) Update(ctx context.Context, book *Domain, id int) (Domain, int, error) {
	book.ID = id
	if err := uc.repo.Update(ctx, book); err != nil {
		return Domain{}, http.StatusInternalServerError, err
	}

	newBook, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return Domain{}, http.StatusNotFound, err
	}

	return newBook, http.StatusOK, err
}

func (uc *bookUsecase) Delete(ctx context.Context, id int) (int, error) {
	_, err := uc.repo.GetById(ctx, id)
	if err != nil { // check wheter data is exists or not
		return http.StatusNotFound, errors.New("book not found")
	}
	err = uc.repo.Delete(ctx, id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
