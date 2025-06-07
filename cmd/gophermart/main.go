package main

import (
	"context"
	"log"

	"github.com/arrowls/praktikum-diploma-1/internal/auth/handlers"
	"github.com/arrowls/praktikum-diploma-1/internal/config"
	"github.com/arrowls/praktikum-diploma-1/internal/database"
	"github.com/arrowls/praktikum-diploma-1/internal/di"
	"github.com/arrowls/praktikum-diploma-1/internal/middleware"
	orderHandlers "github.com/arrowls/praktikum-diploma-1/internal/orders/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	container := di.NewContainer()
	serverConfig := config.ProvideConfig(container)

	db := database.ProvideDatabase(container)
	defer db.Close(context.Background())

	authHandlers := handlers.ProvideAuthHandlers(container)
	orderHandler := orderHandlers.ProvideOrderHandlers(container)

	router := gin.Default()

	publicRouter := router.Group("/api/user")
	privateRouter := router.Group("/api/user")
	privateRouter.Use(middleware.AuthMiddleware())

	publicRouter.POST("/register", authHandlers.Register)
	publicRouter.POST("/login", authHandlers.Login)

	privateRouter.POST("/orders", orderHandler.AddOrder)
	privateRouter.GET("/orders", orderHandler.GetList)

	log.Fatal(router.Run(serverConfig.RunAddress))
}
