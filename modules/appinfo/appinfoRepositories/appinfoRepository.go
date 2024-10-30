package appinfoRepositories

import "github.com/jmoiron/sqlx"

type IAppinfoRepository interface {
}

type appinfoRepository struct {
	db *sqlx.DB
}

func AppinfoRepository(db *sqlx.DB) IAppinfoRepository {
	return &appinfoRepository{db: db}
}
