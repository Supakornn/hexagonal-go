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
	InsertProduct(req *products.Product) (*products.Product, error)
	UpdateProduct(req *products.Product) (*products.Product, error)
	DeleteProduct(productId string) error
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

func (u *productsUsecase) InsertProduct(req *products.Product) (*products.Product, error) {
	product, err := u.ProductsRepository.InsertProduct(req)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *productsUsecase) UpdateProduct(req *products.Product) (*products.Product, error) {
	product, err := u.ProductsRepository.UpdateProduct(req)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *productsUsecase) DeleteProduct(productId string) error {
	err := u.ProductsRepository.DeleteProduct(productId)
	if err != nil {
		return err
	}

	return nil
}
