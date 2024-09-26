package usersUsecases

import (
	"github.com/Supakornn/go-api/config"
	"github.com/Supakornn/go-api/modules/users/usersRepositories"
)

type IUserUsecase interface {
}

type userUsecase struct {
	cfg            config.IConfig
	userRepository usersRepositories.IUserRepository
}

func UserUsecase(cfg config.IConfig, userRepository usersRepositories.IUserRepository) IUserUsecase {
	return &userUsecase{
		cfg:            cfg,
		userRepository: userRepository,
	}
}
