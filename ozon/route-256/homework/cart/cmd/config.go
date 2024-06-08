package main

import (
	"github.com/caarlos0/env/v11"
)

type ServerConfig struct {
	Port string `env:"PORT,required"`
}

type ProductClient struct {
	Host    string `env:"PRODUCT_HOST,required"`
	Token   string `env:"PRODUCT_TOKEN,required"`
	Timeout string `env:"PRODUCT_TIMEOUT,required"`
}

type Config struct {
	ServerConfig  ServerConfig
	ProductClient ProductClient
}

func NewConfig() (Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
