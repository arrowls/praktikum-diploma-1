package repository

import (
	"log"

	"github.com/arrowls/praktikum-diploma-1/internal/database"
	"github.com/arrowls/praktikum-diploma-1/internal/di"
	"github.com/arrowls/praktikum-diploma-1/internal/logger"
)

const diKey = "accrual_repo"

func ProvideAccrualRepo(container di.ContainerInterface) *AccrualRepository {
	if repo, ok := container.Get(diKey).(*AccrualRepository); ok {
		return repo
	}

	repo := &AccrualRepository{
		db:     database.ProvideDatabase(container),
		logger: logger.ProvideLogger(container),
	}

	if err := container.Add(diKey, repo); err != nil {
		log.Fatal(err)
	}

	return repo
}
