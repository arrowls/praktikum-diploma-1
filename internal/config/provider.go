package config

import (
	"log"

	"github.com/arrowls/praktikum-diploma-1/internal/di"
)

const diKey = "config"

func ProvideConfig(container di.ContainerInterface) ServerConfig {
	if cfg, ok := container.Get(diKey).(ServerConfig); ok {
		return cfg
	}

	cfg := NewServerConfig()
	if err := container.Add(diKey, cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}
