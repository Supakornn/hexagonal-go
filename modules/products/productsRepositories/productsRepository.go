package productsRepositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Supakornn/hexagonal-go/config"
	"github.com/Supakornn/hexagonal-go/modules/entities"
	"github.com/Supakornn/hexagonal-go/modules/files/filesUsecases"
	"github.com/Supakornn/hexagonal-go/modules/products"
	"github.com/Supakornn/hexagonal-go/modules/products/productsPatterns"
	"github.com/jmoiron/sqlx"
)

type IProductsRepository interface {
	FindOneProduct(productId string) (*products.Product, error)
	FindProducts(req *products.ProductFilter) ([]*products.Product, int)
	InsertProduct(req *products.Product) (*products.Product, error)
	UpdateProduct(req *products.Product) (*products.Product, error)
	DeleteProduct(productId string) error
}

type productsRepositorty struct {
	db           *sqlx.DB
	cfg          config.IConfig
	filesUsecase filesUsecases.IFilesUsecase
}

func ProductsRepository(db *sqlx.DB, cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase) IProductsRepository {
	return &productsRepositorty{
		db:           db,
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}

func (r *productsRepositorty) FindOneProduct(productId string) (*products.Product, error) {
	query := `
    SELECT
        to_jsonb("t")
    FROM (
        SELECT
            "p"."id",
            "p"."title",
            "p"."description",
            "p"."price",
            (
                SELECT
                    to_jsonb("ct")
                FROM (
                        SELECT
                            "c"."id",
                            "c".title
                        FROM "categories" "c"
                                LEFT JOIN "products_categories" "pc" ON "pc"."category_id" = "c"."id"
                        WHERE "pc"."product_id" = "p"."id"
                    ) AS "ct"
            )AS "category",
            "p"."created_at",
            "p"."updated_at",
            (
                SELECT
                    COALESCE(array_to_json(array_agg("it")), '[]'::json)
                FROM(
                        SELECT
                            "i"."id",
                            "i"."filename",
                            "i"."url"
                        FROM "images" "i"
                        WHERE "i"."product_id" = "p"."id"
                    ) AS "it"
            ) AS "images"
        FROM "products" "p"
        WHERE "p"."id" = $1
        LIMIT 1
    ) AS "t";`

	productBytes := make([]byte, 0)
	product := &products.Product{
		Images: make([]*entities.Image, 0),
	}

	if err := r.db.Get(&productBytes, query, productId); err != nil {
		return nil, fmt.Errorf("find one product failed: %v", err)
	}

	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, fmt.Errorf("unmarshal product failed: %v", err)
	}

	return product, nil
}

func (r *productsRepositorty) FindProducts(req *products.ProductFilter) ([]*products.Product, int) {
	builder := productsPatterns.FindProductBuilder(r.db, req)
	engineer := productsPatterns.FindProductEngineer(builder)

	result := engineer.FindProduct().Result()
	count := engineer.CountProduct().Count()

	return result, count
}

func (r *productsRepositorty) InsertProduct(req *products.Product) (*products.Product, error) {
	builder := productsPatterns.InsertProductBuilder(r.db, req)
	productId, err := productsPatterns.InsertProductEngineer(builder).InsertProduct()
	if err != nil {
		return nil, err
	}

	product, err := r.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productsRepositorty) UpdateProduct(req *products.Product) (*products.Product, error) {
	builder := productsPatterns.UpdateProductBuilder(r.db, req, r.filesUsecase)
	engineer := productsPatterns.UpdateProductEngineer(builder)

	if err := engineer.UpdateProduct(); err != nil {
		return nil, err
	}

	product, err := r.FindOneProduct(req.Id)
	if err != nil {
		return nil, err
	}

	return product, nil

}

func (r *productsRepositorty) DeleteProduct(productId string) error {
	deleteImagesQuery := `DELETE FROM "images" WHERE "product_id" = $1;`
	if _, err := r.db.ExecContext(context.Background(), deleteImagesQuery, productId); err != nil {
		return fmt.Errorf("delete images failed: %v", err)
	}

	query := `DELETE FROM "products" WHERE "id" = $1;`
	if _, err := r.db.ExecContext(context.Background(), query, productId); err != nil {
		return fmt.Errorf("delete product failed: %v", err)
	}
	return nil
}
