package productsUsecases

import (
	"math"

	"github.com/Supakornn/hexagonal-go/modules/entities"
	"github.com/Supakornn/hexagonal-go/modules/products"
	"github.com/Supakornn/hexagonal-go/modules/products/productsRepositories"
)

type IProductsUsecase interface {
	FindOneProduct(productId string) (*products.Product, error)
	FindProducts(req *products.ProductFilter) *entities.PaginateRes
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

func (u *productsUsecase) FindProducts(req *products.ProductFilter) *entities.PaginateRes {
	products, count := u.ProductsRepository.FindProducts(req)
	return &entities.PaginateRes{
		Data:      products,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}
