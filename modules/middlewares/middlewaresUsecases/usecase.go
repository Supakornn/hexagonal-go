package middlewaresUsecases

import "github.com/supakornn/hexagonal-go/modules/middlewares/middlewaresRepositories"

type IMiddlewaresUsecase interface {
}

type middlewaresusecase struct {
	middlewaresRepo middlewaresRepositories.IMiddlewaresRepository
}

func MiddlewaresUsecase(middlewaresRepo middlewaresRepositories.IMiddlewaresRepository) IMiddlewaresUsecase {
	return &middlewaresusecase{
		middlewaresRepo: middlewaresRepo,
	}
}
