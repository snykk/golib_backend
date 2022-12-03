package users

import (
	"context"
	"errors"

	"github.com/snykk/golib_backend/constants"
	"github.com/snykk/golib_backend/helpers"
	"github.com/snykk/golib_backend/packages/token"
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

func (uc *userUsecase) Store(ctx context.Context, domain *Domain) (Domain, error) {
	var err error
	domain.Password, err = helpers.GenerateHash(domain.Password)
	if err != nil {
		return Domain{}, err
	}

	user, err := uc.repo.Store(ctx, domain)
	if err != nil {
		return Domain{}, err
	}

	return user, nil
}

func (uc *userUsecase) Login(ctx context.Context, domain *Domain) (Domain, error) {
	var err error

	userDomain, err := uc.repo.GetByEmail(ctx, domain)
	if err != nil {
		return Domain{}, err
	}

	if !userDomain.IsActivated {
		return Domain{}, errors.New("account is not activated")
	}

	if !helpers.ValidateHash(domain.Password, userDomain.Password) {
		return Domain{}, errors.New("invalid email or password")
	}

	if userDomain.Role == constants.Admin {
		userDomain.Token, err = uc.jwtService.GenerateToken(userDomain.ID, true, userDomain.Email, userDomain.Password)
	} else {
		userDomain.Token, err = uc.jwtService.GenerateToken(userDomain.ID, false, userDomain.Email, userDomain.Password)
	}

	if err != nil {
		return Domain{}, err
	}

	return userDomain, nil
}

func (uc *userUsecase) GetAll() ([]Domain, error) {
	usersFromRepo, err := uc.repo.GetAll()

	if err != nil {
		return []Domain{}, err
	}

	// encapsulate password
	for i := 0; i < len(usersFromRepo); i++ {
		usersFromRepo[i].Password = ""
	}

	return usersFromRepo, err
}

func (uc *userUsecase) GetById(ctx context.Context, id int, idClaims int) (Domain, error) {
	user, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return Domain{}, err
	}

	if id != idClaims {
		user.Password = ""
	}

	return user, nil
}

func (uc *userUsecase) Update(ctx context.Context, domain *Domain, id int) (Domain, error) {
	var err error
	domain.ID = id

	// if domain.Password != "" {
	// 	if domain.Password, err = helpers.GenerateHash(domain.Password); err != nil {
	// 		return Domain{}, err
	// 	}
	// }

	if err := uc.repo.Update(ctx, domain); err != nil {
		return Domain{}, err
	}

	newUserFromDB, err := uc.repo.GetById(ctx, id)

	return newUserFromDB, err
}

func (uc *userUsecase) Delete(ctx context.Context, id int) error {
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

func (uc *userUsecase) GetByEmail(ctx context.Context, email string) (Domain, error) {
	user, err := uc.repo.GetByEmail(ctx, &Domain{Email: email})
	if err != nil {
		return Domain{}, err
	}

	return user, nil
}

func (uc *userUsecase) ActivateUser(ctx context.Context, email string) (err error) {
	user, err := uc.repo.GetByEmail(ctx, &Domain{Email: email})
	if err != nil {
		return err
	}

	if err = uc.repo.Update(ctx, &Domain{ID: user.ID, IsActivated: true}); err != nil {
		return err
	}

	return
}

func (uc *userUsecase) ChangePassword(ctx context.Context, domain *Domain, new_pass string, id int) (err error) {
	domain.ID = id

	if domain.Password == new_pass {
		return errors.New("no changed detected")
	}

	userDom, err := uc.repo.GetById(ctx, domain.ID)
	if err != nil {
		return err
	}

	if !helpers.ValidateHash(domain.Password, userDom.Password) {
		return errors.New("incorrect password")
	}

	domain.Password, err = helpers.GenerateHash(new_pass)
	if err != nil {
		return err
	}

	return uc.repo.Update(ctx, domain)
}

func (uc *userUsecase) ChangeEmail(ctx context.Context, domain *Domain, id int) (err error) {
	domain.ID = id

	userDom, _ := uc.repo.GetByEmail(ctx, domain)
	if userDom.Password != "" {
		return errors.New("email is already in used")
	}

	return uc.repo.UpdateEmail(ctx, domain)
}
