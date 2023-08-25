package productsUsecases

import (
	"math"

	"github.com/supakornn/hexagonal-go/modules/entities"
	"github.com/supakornn/hexagonal-go/modules/products"
	"github.com/supakornn/hexagonal-go/modules/products/productsRepositories"
)

type IProductsUsecase interface {
	FindOneProduct(productId string) (*products.Product, error)
	FindProduct(req *products.ProductFilter) *entities.PaginateRes
	AddProduct(req *products.Product) (*products.Product, error)
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
func (u *productsUsecase) FindProduct(req *products.ProductFilter) *entities.PaginateRes {
	products, count := u.productRepository.FindProduct(req)
	return &entities.PaginateRes{
		Data:      products,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *productsUsecase) AddProduct(req *products.Product) (*products.Product, error) {
	product, err := u.productRepository.InsertProduct(req)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// func (u *productsUsecase) DeleteProduct(productId string) error {
// 	if err := u.productsRepository.DeleteProduct(productId); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (u *productsUsecase) UpdateProduct(req *products.Product) (*products.Product, error) {
// 	product, err := u.productsRepository.UpdateProduct(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return product, nil
// }
