package usersRepositories

import "github.com/jmoiron/sqlx"

type IUserRepository interface {
}

type userRepository struct {
	db *sqlx.DB
}

func UsersRepository(db *sqlx.DB) IUserRepository {
	return &userRepository{
		db: db,
	}
}
