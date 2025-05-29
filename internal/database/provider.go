package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"

	"github.com/arrowls/praktikum-diploma-1/internal/config"
	"github.com/arrowls/praktikum-diploma-1/internal/di"
)

const diKey = "database"

func ProvideDatabase(container di.ContainerInterface) *pgx.Conn {
	if db, ok := container.Get(diKey).(*pgx.Conn); ok {
		return db
	}

	cfg := config.ProvideConfig(container)
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.DatabaseURI)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}
