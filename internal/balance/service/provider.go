package service

import (
	"log"

	"github.com/arrowls/praktikum-diploma-1/internal/balance/repository"
	"github.com/arrowls/praktikum-diploma-1/internal/di"
)

const diKey = "balance_service"

func ProvideBalanceService(container di.ContainerInterface) *BalanceService {
	if service, ok := container.Get(diKey).(*BalanceService); ok {
		return service
	}

	service := &BalanceService{repo: repository.ProvideBalanceRepo(container)}

	if err := container.Add(diKey, service); err != nil {
		log.Fatal(err)
	}

	return service
}
