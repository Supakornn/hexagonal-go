package usersUsecases

import (
	"fmt"

	"github.com/supakornn/hexagonal-go/config"
	"github.com/supakornn/hexagonal-go/modules/users"
	"github.com/supakornn/hexagonal-go/modules/users/usersRepositories"
	"github.com/supakornn/hexagonal-go/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	InsertCustomer(req *users.Register) (*users.UserPassport, error)
	InsertAdmin(req *users.Register) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
	RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error)
	DeleteOauth(oauthId string) error
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

func (u *userUsecase) InsertAdmin(req *users.Register) (*users.UserPassport, error) {
	if err := req.Bcrypt(); err != nil {
		return nil, err
	}
	result, err := u.usersRepo.InsertUser(req, true)
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
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			RoleId:   user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}

	if err := u.usersRepo.InsertOauth(passport); err != nil {
		return nil, err
	}

	return passport, nil
}

func (u *userUsecase) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {
	claims, err := auth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	oauth, err := u.usersRepo.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	profile, err := u.usersRepo.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	newClaims := &users.UserClaims{
		Id:     profile.Id,
		RoleId: profile.RoleId,
	}

	accessToken, err := auth.NewAuth(auth.Access, u.cfg.Jwt(), newClaims)
	if err != nil {
		return nil, err
	}

	refreshToken := auth.RepeatToken(u.cfg.Jwt(), newClaims, claims.ExpiresAt.Unix())

	passport := &users.UserPassport{
		User: profile,
		Token: &users.UserToken{
			Id:           oauth.Id,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}
	if err := u.usersRepo.UpdateOauth(passport.Token); err != nil {
		return nil, err
	}
	return passport, nil
}

func (u *userUsecase) DeleteOauth(oauthId string) error {
	if err := u.usersRepo.DeleteOauth(oauthId); err != nil {
		return err
	}

	return nil
}
