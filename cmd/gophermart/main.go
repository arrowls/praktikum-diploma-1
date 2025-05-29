package main

import (
	"context"
	"log"

	"github.com/arrowls/praktikum-diploma-1/internal/auth/handlers"
	"github.com/arrowls/praktikum-diploma-1/internal/config"
	"github.com/arrowls/praktikum-diploma-1/internal/database"
	"github.com/arrowls/praktikum-diploma-1/internal/di"
	"github.com/gin-gonic/gin"
)

func main() {
	container := di.NewContainer()
	serverConfig := config.ProvideConfig(container)

	db := database.ProvideDatabase(container)
	defer db.Close(context.Background())
	authHandlers := handlers.ProvideAuthHandlers(container)

	router := gin.Default()

	router.POST("/api/user/register", authHandlers.Register)
	router.POST("/api/user/login", authHandlers.Login)

	log.Fatal(router.Run(serverConfig.RunAddress))
}
