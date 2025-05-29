package handlers

import (
	"log"

	"github.com/arrowls/praktikum-diploma-1/internal/apperrors"
	autherrors "github.com/arrowls/praktikum-diploma-1/internal/auth/errors"
	"github.com/arrowls/praktikum-diploma-1/internal/auth/service"
	"github.com/arrowls/praktikum-diploma-1/internal/di"
)

const diKey = "auth_handlers"

func ProvideAuthHandlers(container di.ContainerInterface) *AuthHandlers {
	if handlers, ok := container.Get(diKey).(*AuthHandlers); ok {
		return handlers
	}

	errorHandler := autherrors.AuthErrorHandler{
		NextHandler: apperrors.DefaultErrorHandler,
	}

	handlers := &AuthHandlers{
		service:      service.ProvideAuthService(container),
		errorHandler: errorHandler.Handle,
	}

	if err := container.Add(diKey, handlers); err != nil {
		log.Fatal(err)
	}

	return handlers
}
