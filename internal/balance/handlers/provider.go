package handlers

import (
	"log"

	"github.com/arrowls/praktikum-diploma-1/internal/apperrors"
	balanceerrors "github.com/arrowls/praktikum-diploma-1/internal/balance/errors"
	"github.com/arrowls/praktikum-diploma-1/internal/balance/service"
	"github.com/arrowls/praktikum-diploma-1/internal/di"
)

const diKey = "balance_handlers"

func ProvideBalanceHandlers(container di.ContainerInterface) *BalanceHandlers {
	if handlers, ok := container.Get(diKey).(*BalanceHandlers); ok {
		return handlers
	}

	errorHandler := &balanceerrors.BalanceErrorHandler{
		NextHandler: apperrors.DefaultErrorHandler,
	}

	handlers := &BalanceHandlers{
		service:      service.ProvideBalanceService(container),
		errorHandler: errorHandler.Handle,
	}

	if err := container.Add(diKey, handlers); err != nil {
		log.Fatal(err)
	}

	return handlers
}
