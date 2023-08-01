package usersUsecases

import (
	"fmt"

	"github.com/supakornn/hexagonal-go/config"
	"github.com/supakornn/hexagonal-go/modules/users"
	"github.com/supakornn/hexagonal-go/modules/users/usersRepositories"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	InsertCustomer(req *users.Register) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
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

func (u *userUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	user, err := u.usersRepo.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid")
	}

	passport := &users.UserPassport{
		User: &users.User{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			RoleId:   user.RoleId,
		},
		Token: nil,
	}

	return passport, nil
}
