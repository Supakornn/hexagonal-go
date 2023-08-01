package usersUsecases

import (
	"github.com/supakornn/hexagonal-go/config"
	"github.com/supakornn/hexagonal-go/modules/users"
	"github.com/supakornn/hexagonal-go/modules/users/usersRepositories"
)

type IUserUsecase interface {
	InsertCustomer(req *users.Register) (*users.UserPassport, error)
}

type userUsecase struct {
	cfg       config.Iconfig
	usersRepo usersRepositories.IUserRepository
}

func UserUsecase(cfg config.Iconfig, usersRepo usersRepositories.IUserRepository) IUserUsecase {
	return &userUsecase{
		cfg:       cfg,
		usersRepo: usersRepo,
	}
}

func (u *userUsecase) InsertCustomer(req *users.Register) (*users.UserPassport, error) {
	if err := req.Bcrypt(); err != nil {
		return nil, err
	}
	result, err := u.usersRepo.InsertUser(req, false)
	if err != nil {
		return nil, err
	}
	return result, nil
}
