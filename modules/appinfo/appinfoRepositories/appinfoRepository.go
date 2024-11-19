package appinfoRepositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/Supakornn/hexagonal-go/modules/appinfo"
	"github.com/jmoiron/sqlx"
)

type IAppinfoRepository interface {
	FindCategory(req *appinfo.CategoryFilter) ([]appinfo.Category, error)
	InsertCategory(req []*appinfo.Category) error
}

type appinfoRepository struct {
	db *sqlx.DB
}

func AppinfoRepository(db *sqlx.DB) IAppinfoRepository {
	return &appinfoRepository{db: db}
}

func (r *appinfoRepository) FindCategory(req *appinfo.CategoryFilter) ([]appinfo.Category, error) {
	query := `
	SELECT 
		"id",
		"title"
	FROM "categories"`

	filterValues := make([]any, 0)
	if req.Title != "" {
		query += `
	WHERE LOWER("title") LIKE $1`
		filterValues = append(filterValues, "%"+strings.ToLower(req.Title)+"%")
	}

	query += ";"

	category := make([]appinfo.Category, 0)
	if err := r.db.Select(&category, query, filterValues...); err != nil {
		return nil, fmt.Errorf("error get category: %w", err)
	}

	return category, nil
}

func (r *appinfoRepository) InsertCategory(req []*appinfo.Category) error {
	ctx := context.Background()
	query := `
	INSERT INTO "categories" (
	"title"
	)
	VALUES`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	valuesStack := make([]any, 0)

	for i, cat := range req {
		valuesStack = append(valuesStack, cat.Title)

		if i != len(req)-1 {
			query += fmt.Sprintf(`($%d),`, i+1)
		} else {
			query += fmt.Sprintf(`($%d)`, i+1)
		}
	}

	query += `
	RETURNING "id";`

	rows, err := tx.QueryxContext(ctx, query, valuesStack...)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error insert category: %w", err)
	}

	var i int

	for rows.Next() {
		if err := rows.Scan(&req[i].Id); err != nil {
			return fmt.Errorf("error scan category: %w", err)
		}

		i++
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
