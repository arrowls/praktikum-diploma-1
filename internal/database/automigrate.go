package database

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func AutoMigrate(databaseURI string) error {
	migrationsPath, err := filepath.Abs("./migrations")
	if err != nil {
		return fmt.Errorf("migration dir not found")
	}
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		databaseURI,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
