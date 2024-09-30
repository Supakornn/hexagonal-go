package usersUsecases

import (
	"github.com/Supakornn/hexagonal-go/config"
	"github.com/Supakornn/hexagonal-go/modules/users"
	"github.com/Supakornn/hexagonal-go/modules/users/usersRepositories"
)

type IUserUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
}

type userUsecase struct {
	cfg            config.IConfig
	userRepository usersRepositories.IUserRepository
}

func UsersUsecase(cfg config.IConfig, userRepository usersRepositories.IUserRepository) IUserUsecase {
	return &userUsecase{
		cfg:            cfg,
		userRepository: userRepository,
	}
}

func (u *userUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	result, err := u.userRepository.InsertUser(req, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}
