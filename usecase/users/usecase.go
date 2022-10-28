package users

import (
	"context"
	"errors"

	encrpyt "github.com/snykk/golib_backend/utils/hash"
	"github.com/snykk/golib_backend/utils/token"
)

type UserUsecase struct {
	JwtService token.JWTService
	Repo       Repository
}

func NewUserUsecase(repo Repository, jwtService token.JWTService) Usecase {
	return &UserUsecase{
		JwtService: jwtService,
		Repo:       repo,
	}
}

func (userUC UserUsecase) Store(ctx context.Context, domain *Domain) (Domain, error) {
	var err error
	domain.Password, err = encrpyt.GenerateHash(domain.Password)
	if err != nil {
		return Domain{}, err
	}

	user, err := userUC.Repo.Store(ctx, domain)
	if err != nil {
		return Domain{}, err
	}

	return user, nil
}

func (userUC UserUsecase) Login(ctx context.Context, domain *Domain) (Domain, error) {
	var err error

	userDomain, err := userUC.Repo.GetByEmail(ctx, domain)
	if err != nil {
		return Domain{}, err
	}

	if !encrpyt.ValidateHash(domain.Password, userDomain.Password) {
		return Domain{}, errors.New("invalid username or password")
	}

	if userDomain.IsAdmin {
		userDomain.Token, err = userUC.JwtService.GenerateToken(userDomain.Id, true)
	} else {
		userDomain.Token, err = userUC.JwtService.GenerateToken(userDomain.Id, false)
	}

	if err != nil {
		return Domain{}, err
	}

	return userDomain, nil
}

func (userUC UserUsecase) GetAll() ([]Domain, error) {
	usersFromRepo, err := userUC.Repo.GetAll()

	if err != nil {
		return []Domain{}, err
	}

	return usersFromRepo, err
}

func (userUC UserUsecase) GetById(ctx context.Context, id int) (Domain, error) {
	user, err := userUC.Repo.GetById(ctx, id)
	if err != nil {
		return Domain{}, err
	}
	return user, nil
}

func (userUC UserUsecase) Update(ctx context.Context, domain *Domain, id int) (Domain, error) {
	var err error
	domain.Id = id

	if domain.Password, err = encrpyt.GenerateHash(domain.Password); err != nil {
		return Domain{}, err
	}

	if err := userUC.Repo.Update(ctx, domain); err != nil {
		return Domain{}, err
	}

	newUserFromDB, err := userUC.Repo.GetById(ctx, id)

	return newUserFromDB, err
}

func (userUC UserUsecase) Delete(ctx context.Context, id int) error {
	_, err := userUC.Repo.GetById(ctx, id)
	if err != nil { // check wheter data is exists or not
		return err
	}
	err = userUC.Repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
