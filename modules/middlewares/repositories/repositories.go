package repositories

import "github.com/jmoiron/sqlx"

type IRepository interface {
}

type repository struct {
	db *sqlx.DB
}

func MiddlewarsRepo(db *sqlx.DB) IRepository {
	return &repository{
		db: db,
	}
}
