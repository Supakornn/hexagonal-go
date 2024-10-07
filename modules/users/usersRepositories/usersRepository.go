package usersRepositories

import (
	"context"
	"fmt"
	"time"

	"github.com/Supakornn/hexagonal-go/modules/users"
	"github.com/Supakornn/hexagonal-go/modules/users/usersPatterns"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error)
	FindOneUserByEmail(email string) (*users.UserCredentialCheck, error)
	InsertOauth(req *users.UserPassport) error
}

type userRepository struct {
	db *sqlx.DB
}

func UsersRepository(db *sqlx.DB) IUserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error) {
	result := usersPatterns.InsertUser(r.db, req, isAdmin)
	var err error

	if isAdmin {
		result, err = result.Admin()
		if err != nil {
			return nil, err
		}
	} else {
		result, err = result.Customer()
		if err != nil {
			return nil, err
		}
	}

	user, err := result.Result()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) FindOneUserByEmail(email string) (*users.UserCredentialCheck, error) {
	fmt.Printf("Querying for email: %s\n", email)

	query := `
	SELECT 
		"id",
		"email",
		"password",
		"username",
		"role_id"
	FROM "users"
	WHERE "email" = $1;`

	user := new(users.UserCredentialCheck)

	err := r.db.Get(user, query, email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (r *userRepository) InsertOauth(req *users.UserPassport) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
	INSERT INTO "oauth" (
		"userid",	
		"refresh_token",
		"access_token"
	)
	VALUES ($1, $2, $3)
	RETURNING "id";`

	if err := r.db.QueryRowContext(ctx,
		query,
		req.User.ID,
		req.Token.RefreshToken,
		req.Token.AccessToken,
	).Scan(&req.Token.Id); err != nil {
		fmt.Printf("insert oauth failed: %v\n", err)
	}

	return nil
}
