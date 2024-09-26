package middlewareHandlers

import (
	"github.com/Supakornn/go-api/modules/middlewares/middlewareUsecases"
)

type IMiddlewaresHandler interface {
}

type middlewaresHandler struct {
	middlewareUsecase middlewareUsecases.IMiddlewaresUsecase
}

func MiddlewaresHandler(middlewareUsecase middlewareUsecases.IMiddlewaresUsecase) IMiddlewaresHandler {
	return &middlewaresHandler{
		middlewareUsecase: middlewareUsecase,
	}
}
