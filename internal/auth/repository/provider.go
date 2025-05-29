package repository

import (
	"log"

	"github.com/arrowls/praktikum-diploma-1/internal/database"
	"github.com/arrowls/praktikum-diploma-1/internal/di"
)

const diKey = "auth_repo"

func ProvideAuthRepo(container di.ContainerInterface) *AuthRepository {
	if repo, ok := container.Get(diKey).(*AuthRepository); ok {
		return repo
	}

	repo := &AuthRepository{db: database.ProvideDatabase(container)}

	if err := container.Add(diKey, repo); err != nil {
		log.Fatal(err)
	}

	return repo
}
