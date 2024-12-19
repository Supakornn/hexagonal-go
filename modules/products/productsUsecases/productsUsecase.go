package productsUsecases

import (
	"github.com/Supakornn/hexagonal-go/modules/products"
	"github.com/Supakornn/hexagonal-go/modules/products/productsRepositories"
)

type IProductsUsecase interface {
	FindOneProduct(productId string) (*products.Product, error)
}

type productsUsecase struct {
	ProductsRepository productsRepositories.IProductsRepository
}

func ProductUsecase(ProductsRepository productsRepositories.IProductsRepository) IProductsUsecase {
	return &productsUsecase{
		ProductsRepository: ProductsRepository,
	}
}

func (u *productsUsecase) FindOneProduct(productId string) (*products.Product, error) {
	product, err := u.ProductsRepository.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}

	return product, nil
}
