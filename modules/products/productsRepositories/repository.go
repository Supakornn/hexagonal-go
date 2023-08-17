package productsRepositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/supakornn/hexagonal-go/config"
	"github.com/supakornn/hexagonal-go/modules/files/filesUsecases"
)

type IProductsRepository interface{}

type prodictsRepository struct {
	db          *sqlx.DB
	cfg         config.Iconfig
	fileUsecase filesUsecases.IFilesUsecase
}

func ProductsRepository(db *sqlx.DB, cfg config.Iconfig, fileUsecase filesUsecases.IFilesUsecase) IProductsRepository {
	return &prodictsRepository{
		db:          db,
		cfg:         cfg,
		fileUsecase: fileUsecase,
	}
}
