package middlewaresUsecases

import (
	"github.com/supakornn/hexagonal-go/modules/middlewares"
	"github.com/supakornn/hexagonal-go/modules/middlewares/middlewaresRepositories"
)

type IMiddlewaresUsecase interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewaresusecase struct {
	middlewaresRepo middlewaresRepositories.IMiddlewaresRepository
}

func MiddlewaresUsecase(middlewaresRepo middlewaresRepositories.IMiddlewaresRepository) IMiddlewaresUsecase {
	return &middlewaresusecase{
		middlewaresRepo: middlewaresRepo,
	}
}

func (u *middlewaresusecase) FindAccessToken(userId, accessToken string) bool {
	return u.middlewaresRepo.FindAccessToken(userId, accessToken)
}

func (u *middlewaresusecase) FindRole() ([]*middlewares.Role, error) {
	roles, err := u.middlewaresRepo.FindRole()
	if err != nil {
		return nil, err
	}
	return roles, nil
}
