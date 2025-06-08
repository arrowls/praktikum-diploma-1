package service

import (
	"log"
	"time"

	"github.com/arrowls/praktikum-diploma-1/internal/accrual/repository"
	"github.com/arrowls/praktikum-diploma-1/internal/config"
	"github.com/arrowls/praktikum-diploma-1/internal/di"
	"github.com/arrowls/praktikum-diploma-1/internal/logger"
)

const diKey = "accrual_service"

func ProvideAccrualService(container di.ContainerInterface) *AccrualService {
	if service, ok := container.Get(diKey).(*AccrualService); ok {
		return service
	}

	cfg := config.ProvideConfig(container)

	service := &AccrualService{
		repo:            repository.ProvideAccrualRepo(container),
		ticker:          time.NewTicker(time.Duration(cfg.PingInterval) * time.Second),
		baseInterval:    cfg.PingInterval,
		logger:          logger.ProvideLogger(container),
		accrualEndpoint: cfg.AccrualSystemAddress,
	}

	if err := container.Add(diKey, service); err != nil {
		log.Fatal(err)
	}

	return service
}
