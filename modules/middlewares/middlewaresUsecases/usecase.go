package middlewaresUsecases

import "github.com/supakornn/hexagonal-go/modules/middlewares/middlewaresRepositories"

type IMiddlewaresUsecase interface {
	FindAccessToken(userId, accessToken string) bool
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
