package ordersUsecases

import (
	"fmt"
	"math"

	"github.com/Supakornn/hexagonal-go/modules/entities"
	"github.com/Supakornn/hexagonal-go/modules/orders"
	"github.com/Supakornn/hexagonal-go/modules/orders/ordersRepositories"
	"github.com/Supakornn/hexagonal-go/modules/products/productsRepositories"
)

type IOrdersUsecase interface {
	FindOneOrder(orderId string) (*orders.Order, error)
	FindOrders(req *orders.OrderFilter) *entities.PaginateRes
	InsertOrder(req *orders.Order) (*orders.Order, error)
	UpdateOrder(req *orders.Order) (*orders.Order, error)
}

type ordersUsecase struct {
	orderRepository    ordersRepositories.IOrdersRepository
	productsRepository productsRepositories.IProductsRepository
}

func OrdersUsecase(orderRepository ordersRepositories.IOrdersRepository, productsRepository productsRepositories.IProductsRepository) IOrdersUsecase {
	return &ordersUsecase{
		orderRepository:    orderRepository,
		productsRepository: productsRepository,
	}
}

func (u *ordersUsecase) FindOneOrder(orderId string) (*orders.Order, error) {
	order, err := u.orderRepository.FindOneOrder(orderId)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (u *ordersUsecase) FindOrders(req *orders.OrderFilter) *entities.PaginateRes {
	orders, count := u.orderRepository.FindOrders(req)

	return &entities.PaginateRes{
		Data:      orders,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *ordersUsecase) InsertOrder(req *orders.Order) (*orders.Order, error) {
	for i := range req.Products {
		if req.Products[i].Product == nil {
			return nil, fmt.Errorf("product not found")
		}

		prod, err := u.productsRepository.FindOneProduct(req.Products[i].Product.Id)
		if err != nil {
			return nil, fmt.Errorf("product not found")
		}

		req.TotalPaid += req.Products[i].Product.Price * float64(req.Products[i].Qty)
		req.Products[i].Product = prod
	}

	orderId, err := u.orderRepository.InsertOrder(req)
	if err != nil {
		return nil, err
	}

	order, err := u.orderRepository.FindOneOrder(orderId)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (u *ordersUsecase) UpdateOrder(req *orders.Order) (*orders.Order, error) {
	if err := u.orderRepository.UpdateOrder(req); err != nil {
		return nil, err
	}

	order, err := u.orderRepository.FindOneOrder(req.Id)
	if err != nil {
		return nil, err
	}

	return order, nil

}
