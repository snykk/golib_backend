package users

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/snykk/golib_backend/constants"
	"github.com/snykk/golib_backend/helpers"
	"github.com/snykk/golib_backend/http/token"
)

type userUsecase struct {
	jwtService token.JWTService
	repo       Repository
}

func NewUserUsecase(repo Repository, jwtService token.JWTService) Usecase {
	return &userUsecase{
		jwtService: jwtService,
		repo:       repo,
	}
}

func (uc *userUsecase) Store(ctx context.Context, domain *Domain) (Domain, int, error) {
	var err error
	domain.Password, err = helpers.GenerateHash(domain.Password)
	if err != nil {
		return Domain{}, http.StatusInternalServerError, err
	}

	user, err := uc.repo.Store(ctx, domain)
	if err != nil {
		return Domain{}, http.StatusInternalServerError, err
	}

	return user, http.StatusCreated, nil
}

func (uc *userUsecase) Login(ctx context.Context, domain *Domain) (Domain, int, error) {
	var err error

	userDomain, err := uc.repo.GetByEmail(ctx, domain)
	if err != nil {
		return Domain{}, http.StatusUnauthorized, errors.New("invalid email or password") // for security purpose better use generic error message
	}

	if !userDomain.IsActivated {
		return Domain{}, http.StatusForbidden, errors.New("account is not activated")
	}

	if !helpers.ValidateHash(domain.Password, userDomain.Password) {
		return Domain{}, http.StatusUnauthorized, errors.New("invalid email or password")
	}

	if userDomain.Role == constants.Admin {
		userDomain.Token, err = uc.jwtService.GenerateToken(userDomain.ID, true, userDomain.Email, userDomain.Password)
	} else {
		userDomain.Token, err = uc.jwtService.GenerateToken(userDomain.ID, false, userDomain.Email, userDomain.Password)
	}

	if err != nil {
		return Domain{}, http.StatusInternalServerError, err
	}

	return userDomain, http.StatusOK, nil
}

func (uc *userUsecase) GetAll(ctx context.Context) ([]Domain, int, error) {
	usersFromRepo, err := uc.repo.GetAll(ctx)

	if err != nil {
		return []Domain{}, http.StatusInternalServerError, err
	}

	return usersFromRepo, http.StatusOK, nil
}

func (uc *userUsecase) GetById(ctx context.Context, id int, idClaims int) (Domain, int, error) {
	user, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return Domain{}, http.StatusNotFound, errors.New("user not found")
	}

	if id != idClaims {
		user.Password = ""
	}

	return user, http.StatusOK, nil
}

func (uc *userUsecase) Update(ctx context.Context, domain *Domain, id int) (Domain, int, error) {
	var err error
	domain.ID = id

	if err := uc.repo.Update(ctx, domain); err != nil {
		return Domain{}, http.StatusInternalServerError, err
	}

	newUserFromDB, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return Domain{}, http.StatusInternalServerError, err
	}

	return newUserFromDB, http.StatusOK, nil
}

func (uc *userUsecase) Delete(ctx context.Context, id int) (int, error) {
	err := uc.repo.Delete(ctx, id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusNoContent, nil
}

func (uc *userUsecase) GetByEmail(ctx context.Context, email string) (Domain, int, error) {
	user, err := uc.repo.GetByEmail(ctx, &Domain{Email: email})
	if err != nil {
		return Domain{}, http.StatusNotFound, errors.New("email not found")
	}

	return user, http.StatusOK, nil
}

func (uc *userUsecase) ActivateUser(ctx context.Context, email string) (statusCode int, err error) {
	user, err := uc.repo.GetByEmail(ctx, &Domain{Email: email})
	if err != nil {
		return http.StatusNotFound, errors.New("email not found")
	}

	if err = uc.repo.Update(ctx, &Domain{ID: user.ID, IsActivated: true}); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (uc *userUsecase) ChangePassword(ctx context.Context, domain *Domain, new_pass string, id int) (statusCode int, err error) {
	domain.ID = id

	if domain.Password == new_pass {
		return http.StatusBadRequest, errors.New("no changed detected")
	}

	userDom, err := uc.repo.GetById(ctx, domain.ID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("user with id %d not found", domain.ID)
	}

	if !helpers.ValidateHash(domain.Password, userDom.Password) {
		return http.StatusUnauthorized, errors.New("incorrect password")
	}

	domain.Password, err = helpers.GenerateHash(new_pass)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if err = uc.repo.Update(ctx, domain); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (uc *userUsecase) ChangeEmail(ctx context.Context, domain *Domain, id int) (statusCode int, err error) {
	domain.ID = id

	userDom, _ := uc.repo.GetByEmail(ctx, domain)
	if userDom.Password != "" {
		return http.StatusConflict, errors.New("email is already in used")
	}

	if err = uc.repo.UpdateEmail(ctx, domain); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
