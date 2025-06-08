package service

import (
	"log"

	"github.com/arrowls/praktikum-diploma-1/internal/di"
	"github.com/arrowls/praktikum-diploma-1/internal/orders/repository"
)

const diKey = "order_service"

func ProvideOrderService(container di.ContainerInterface) *OrderService {
	if service, ok := container.Get(diKey).(*OrderService); ok {
		return service
	}

	service := &OrderService{repo: repository.ProvideOrdersRepo(container)}

	if err := container.Add(diKey, service); err != nil {
		log.Fatal(err)
	}

	return service
}
