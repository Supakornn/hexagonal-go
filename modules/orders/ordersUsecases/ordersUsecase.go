package ordersUsecases

import (
	"github.com/Supakornn/hexagonal-go/modules/orders"
	"github.com/Supakornn/hexagonal-go/modules/orders/ordersRepositories"
	"github.com/Supakornn/hexagonal-go/modules/products/productsRepositories"
)

type IOrdersUsecase interface {
	FindOneOrder(orderId string) (*orders.Order, error)
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

func (o *ordersUsecase) FindOneOrder(orderId string) (*orders.Order, error) {
	order, err := o.orderRepository.FindOneOrder(orderId)
	if err != nil {
		return nil, err
	}

	return order, nil
}
