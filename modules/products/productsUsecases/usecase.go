package productsUsecases

import "github.com/supakornn/hexagonal-go/modules/products/productsRepositories"

type IProductsUsecase interface {
}

type productsUsecase struct {
	productRepository productsRepositories.IProductsRepository
}

func ProductsUsecase(productRepository productsRepositories.IProductsRepository) IProductsUsecase {
	return &productsUsecase{
		productRepository: productRepository,
	}
}
