package productsUsecases

import (
	"github.com/supakornn/hexagonal-go/modules/products"
	"github.com/supakornn/hexagonal-go/modules/products/productsRepositories"
)

type IProductsUsecase interface {
	FindOneProduct(productId string) (*products.Product, error)
}

type productsUsecase struct {
	productRepository productsRepositories.IProductsRepository
}

func ProductsUsecase(productRepository productsRepositories.IProductsRepository) IProductsUsecase {
	return &productsUsecase{
		productRepository: productRepository,
	}
}

func (u *productsUsecase) FindOneProduct(productId string) (*products.Product, error) {
	product, err := u.productRepository.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}

	return product, nil
}
