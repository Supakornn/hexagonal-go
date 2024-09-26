package usersHandlers

import (
	"github.com/Supakornn/go-api/config"
	"github.com/Supakornn/go-api/modules/users/usersUsecases"
)

type IUserHandler interface {
}

type userHandler struct {
	cfg         config.IConfig
	userUsecase usersUsecases.IUserUsecase
}

func UserHandler() IUserHandler {
	return &userHandler{}
}
