package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env"
)

type config struct {
	ConnectAddr string `env:"RUN_ADDRESS"`
	DataBaseURI string `env:"DATABASE_URI"`
	AccrualURI  string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	LogLevel    string `env:"LOG_LEVEL"`
}

var cfg config

func parseFlags() error {
	flag.StringVar(&cfg.ConnectAddr, "a", "localhost:8080", "address to run HTTP server")
	flag.StringVar(&cfg.DataBaseURI, "d", "", "URI to database")
	flag.StringVar(&cfg.AccrualURI, "r", "", "URI to accrual system")
	flag.StringVar(&cfg.LogLevel, "l", "debug", "Log level")

	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("can't parse env; %w", err)
	}

	return nil
}
