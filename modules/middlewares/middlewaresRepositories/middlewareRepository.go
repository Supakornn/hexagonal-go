package middlewaresRepositories

import (
	"fmt"

	"github.com/Supakornn/hexagonal-go/modules/middlewares"
	"github.com/jmoiron/sqlx"
)

type IMiddlewaresRepository interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewaresRepository struct {
	db *sqlx.DB
}

func MiddlewaresRepository(db *sqlx.DB) IMiddlewaresRepository {
	return &middlewaresRepository{
		db: db,
	}
}

func (r *middlewaresRepository) FindAccessToken(userId, accessToken string) bool {
	query := `
	SELECT 
		(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
	FROM "oauth"
	WHERE "userid" = $1
	AND "access_token" = $2;`

	var check bool
	if err := r.db.Get(&check, query, userId, accessToken); err != nil {
		return false
	}

	return check
}

func (r *middlewaresRepository) FindRole() ([]*middlewares.Role, error) {
	query := `
	SELECT
		"id",
		"title" 
	FROM "roles"
	ORDER BY "id" DESC;`

	role := make([]*middlewares.Role, 0)
	if err := r.db.Select(&role, query); err != nil {
		return nil, fmt.Errorf("roles are not found")
	}

	return role, nil
}
