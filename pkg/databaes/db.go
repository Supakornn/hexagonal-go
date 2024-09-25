package databaes

import (
	"log"

	"github.com/Supakornn/go-api/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func DbConnect(cfg config.IDbConfig) *sqlx.DB {
	db, err := sqlx.Connect("pgx", cfg.Url())
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	db.DB.SetMaxOpenConns(cfg.MaxOpenConns())
	return db
}
