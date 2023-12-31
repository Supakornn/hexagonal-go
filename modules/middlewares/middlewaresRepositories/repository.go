package middlewaresRepositories

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/supakornn/hexagonal-go/modules/middlewares"
)

type IMiddlewaresRepository interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewaresrepository struct {
	db *sqlx.DB
}

func MiddlewaresRepo(db *sqlx.DB) IMiddlewaresRepository {
	return &middlewaresrepository{
		db: db,
	}
}

func (r *middlewaresrepository) FindAccessToken(userId, accessToken string) bool {
	query := `
	SELECT 
		(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
	FROM "oauth"
	WHERE "user_id" = $1
	AND "access_token" = $2;`

	var check bool
	if err := r.db.Get(&check, query, userId, accessToken); err != nil {
		return false
	}

	return true
}

func (r *middlewaresrepository) FindRole() ([]*middlewares.Role, error) {
	query := `
	SELECT 
		"id",
		"title"
	FROM "roles"
	ORDER BY "id" DESC;`

	roles := make([]*middlewares.Role, 0)

	if err := r.db.Select(&roles, query); err != nil {
		return nil, fmt.Errorf("role is empty")
	}

	return roles, nil
}
