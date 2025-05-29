package service

import (
	"log"

	"github.com/arrowls/praktikum-diploma-1/internal/auth/repository"
	"github.com/arrowls/praktikum-diploma-1/internal/di"
)

const diKey = "auth_service"

func ProvideAuthService(container di.ContainerInterface) *AuthService {
	if service, ok := container.Get(diKey).(*AuthService); ok {
		return service
	}

	service := &AuthService{repo: repository.ProvideAuthRepo(container)}

	if err := container.Add(diKey, service); err != nil {
		log.Fatal(err)
	}

	return service
}
