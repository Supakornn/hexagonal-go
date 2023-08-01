package middlewaresRepositories

import "github.com/jmoiron/sqlx"

type IMiddlewaresRepository interface {
}

type middlewaresrepository struct {
	db *sqlx.DB
}

func MiddlewaresRepo(db *sqlx.DB) IMiddlewaresRepository {
	return &middlewaresrepository{
		db: db,
	}
}
