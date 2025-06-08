package repository

import (
	"log"

	"github.com/arrowls/praktikum-diploma-1/internal/database"
	"github.com/arrowls/praktikum-diploma-1/internal/di"
	"github.com/arrowls/praktikum-diploma-1/internal/logger"
)

const diKey = "balance_repo"

func ProvideBalanceRepo(container di.ContainerInterface) *BalanceRepository {
	if repo, ok := container.Get(diKey).(*BalanceRepository); ok {
		return repo
	}

	repo := &BalanceRepository{
		db:     database.ProvideDatabase(container),
		logger: logger.ProvideLogger(container),
	}

	if err := container.Add(diKey, repo); err != nil {
		log.Fatal(err)
	}

	return repo
}
