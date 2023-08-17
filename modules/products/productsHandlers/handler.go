package productsHandlers

import (
	"github.com/supakornn/hexagonal-go/config"
	"github.com/supakornn/hexagonal-go/modules/files/filesUsecases"
	"github.com/supakornn/hexagonal-go/modules/products/productsUsecases"
)

type IProductsHandler interface {
}

type productsHandler struct {
	cfg            config.Iconfig
	productUsecase productsUsecases.IProductsUsecase
	fileUsecase    filesUsecases.IFilesUsecase
}

func ProductsHandler(cfg config.Iconfig, productUsecase productsUsecases.IProductsUsecase, fileUsecase filesUsecases.IFilesUsecase) IProductsHandler {
	return &productsHandler{
		productUsecase: productUsecase,
		cfg:            cfg,
		fileUsecase:    fileUsecase,
	}
}
