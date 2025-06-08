package handlers

import (
	"log"

	"github.com/arrowls/praktikum-diploma-1/internal/apperrors"
	"github.com/arrowls/praktikum-diploma-1/internal/di"
	ordererrors "github.com/arrowls/praktikum-diploma-1/internal/orders/errors"
	"github.com/arrowls/praktikum-diploma-1/internal/orders/service"
)

const diKey = "order_handlers"

func ProvideOrderHandlers(container di.ContainerInterface) *OrderHandlers {
	if handlers, ok := container.Get(diKey).(*OrderHandlers); ok {
		return handlers
	}

	errorHandler := ordererrors.OrderErrorHandler{
		NextHandler: apperrors.DefaultErrorHandler,
	}

	handlers := &OrderHandlers{
		service:      service.ProvideOrderService(container),
		errorHandler: errorHandler.Handle,
	}

	if err := container.Add(diKey, handlers); err != nil {
		log.Fatal(err)
	}

	return handlers
}
