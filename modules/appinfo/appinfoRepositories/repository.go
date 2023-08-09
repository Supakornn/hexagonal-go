package appinfoRepositories

import "github.com/jmoiron/sqlx"

type IAppinfoRepository interface{}

type appinfoRepsitory struct {
	db *sqlx.DB
}

func AppinfoRepository(db *sqlx.DB) IAppinfoRepository {
	return &appinfoRepsitory{db: db}
}
