package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v7"
)

const (
	runAddressDefault = ":8080"
)

type ServerConfig struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

var serverConfig ServerConfig

func NewServerConfig() ServerConfig {
	flag.StringVar(&serverConfig.RunAddress, "a", runAddressDefault, "server endpoint url")
	flag.StringVar(&serverConfig.DatabaseURI, "d", "", "databese URI")
	flag.StringVar(&serverConfig.AccrualSystemAddress, "r", "", "accrual system endpoint")

	flag.Parse()

	if err := env.Parse(&serverConfig); err != nil {
		fmt.Printf("Failed to parse env: %v\n", err)
		os.Exit(1)
	}

	return serverConfig
}
