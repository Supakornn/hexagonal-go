package usersUsecases

import (
	"fmt"

	"github.com/Supakornn/hexagonal-go/config"
	"github.com/Supakornn/hexagonal-go/modules/users"
	"github.com/Supakornn/hexagonal-go/modules/users/usersRepositories"
	"github.com/Supakornn/hexagonal-go/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
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

func (u *userUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	user, err := u.userRepository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is incorrect")
	}

	accessToken, err := auth.NewAuth(auth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})

	refreshToken, err := auth.NewAuth(auth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})

	passport := &users.UserPassport{
		User: &users.User{
			ID:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			RoleId:   user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}

	if err := u.userRepository.InsertOauth(passport); err != nil {
		return nil, err
	}

	return passport, nil
}
