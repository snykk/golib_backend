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

func (uc UserUsecase) Store(ctx context.Context, domain *Domain) (Domain, error) {
	var err error
	domain.Password, err = encrpyt.GenerateHash(domain.Password)
	if err != nil {
		return Domain{}, err
	}

	user, err := uc.Repo.Store(ctx, domain)
	if err != nil {
		return Domain{}, err
	}

	return user, nil
}

func (uc UserUsecase) GetAll() ([]Domain, error) {
	usersFromRepo, err := uc.Repo.GetAll()

	if err != nil {
		return []Domain{}, err
	}

	return usersFromRepo, err
}

func (uc UserUsecase) GetById(ctx context.Context, id int) (Domain, error) {
	user, err := uc.Repo.GetById(ctx, id)
	if err != nil {
		return Domain{}, err
	}
	return user, nil
}

func (uc UserUsecase) Update(ctx context.Context, domain *Domain, id int) (Domain, error) {
	userFromDB, err := uc.Repo.GetById(ctx, id)

	if err != nil {
		return Domain{}, err
	}

	domain.Id = userFromDB.Id
	domain.CreatedAt = userFromDB.UpdatedAt

	result, err := uc.Repo.Update(ctx, domain)

	if err != nil {
		return Domain{}, err
	}

	return result, nil

}

func (uc UserUsecase) Delete(ctx context.Context, id int) error {
	_, err := uc.Repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	err = uc.Repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (uc UserUsecase) Login(ctx context.Context, domain *Domain) (Domain, error) {
	var err error

	result, err := uc.Repo.GetByEmail(ctx, domain)
	if err != nil {
		return Domain{}, err
	}

	if !encrpyt.ValidateHash(domain.Password, result.Password) {
		return Domain{}, errors.New("userame or password is not valid")
	}

	if result.IsAdmin {
		result.Token, err = uc.JwtService.GenerateToken(result.Id, true)
	} else {
		result.Token, err = uc.JwtService.GenerateToken(result.Id, false)
	}

	if err != nil {
		return Domain{}, err
	}
	return result, nil
}
