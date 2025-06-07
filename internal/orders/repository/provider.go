package repository

import (
	"log"

	"github.com/arrowls/praktikum-diploma-1/internal/database"
	"github.com/arrowls/praktikum-diploma-1/internal/di"
	"github.com/arrowls/praktikum-diploma-1/internal/logger"
)

const diKey = "orders_repo"

func ProvideOrdersRepo(container di.ContainerInterface) *OrdersRepository {
	if repo, ok := container.Get(diKey).(*OrdersRepository); ok {
		return repo
	}

	repo := &OrdersRepository{
		db:     database.ProvideDatabase(container),
		logger: logger.ProvideLogger(container),
	}

	if err := container.Add(diKey, repo); err != nil {
		log.Fatal(err)
	}

	return repo
}
