package usecases

import "github.com/supakornn/hexagonal-go/modules/middlewares/repositories"

type IUsecase interface {
}

type usecase struct {
	middlewaresRepo repositories.IRepository
}

func MiddlewareUsecase(middlewaresRepo repositories.IRepository) IUsecase {
	return &usecase{
		middlewaresRepo: middlewaresRepo,
	}
}
