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

	if !userDomain.IsActive {
		return Domain{}, errors.New("account is not activated")
	}

	if !encrpyt.ValidateHash(domain.Password, userDomain.Password) {
		return Domain{}, errors.New("invalid email or password")
	}

	if userDomain.IsAdmin {
		userDomain.Token, err = userUC.JwtService.GenerateToken(userDomain.ID, true)
	} else {
		userDomain.Token, err = userUC.JwtService.GenerateToken(userDomain.ID, false)
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

	// encapsulate password
	for i := 0; i < len(usersFromRepo); i++ {
		usersFromRepo[i].Password = ""
	}

	return usersFromRepo, err
}

func (userUC UserUsecase) GetById(ctx context.Context, id int, authHeader string) (Domain, error) {
	user, err := userUC.Repo.GetById(ctx, id)
	if err != nil {
		return Domain{}, err
	}

	// get claims
	claims, _ := userUC.JwtService.ParseToken(authHeader)
	if id != claims.UserID {
		user.Password = ""
	}

	return user, nil
}

func (userUC UserUsecase) Update(ctx context.Context, domain *Domain, id int) (Domain, error) {
	var err error
	domain.ID = id

	if domain.Password != "" {
		if domain.Password, err = encrpyt.GenerateHash(domain.Password); err != nil {
			return Domain{}, err
		}
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

func (userUC UserUsecase) GetByEmail(ctx context.Context, email string) (Domain, error) {
	user, err := userUC.Repo.GetByEmail(ctx, &Domain{Email: email})
	if err != nil {
		return Domain{}, err
	}

	return user, nil
}

func (userUC UserUsecase) ActivateUser(ctx context.Context, email string) (err error) {
	user, err := userUC.Repo.GetByEmail(ctx, &Domain{Email: email})
	if err != nil {
		return err
	}

	if err = userUC.Repo.Update(ctx, &Domain{ID: user.ID, IsActive: true}); err != nil {
		return err
	}

	return
}
