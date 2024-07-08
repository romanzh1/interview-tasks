package config

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
	Rps     int    `env:"PRODUCT_RPS,required"`
}

type LomsClient struct {
	Host string `env:"LOMS_HOST,required"`
}

type Observability struct {
	TracerHost string `env:"TRACER_HOST,required"`
	TracerPort string `env:"TRACER_PORT,required"`
	MetricPort string `env:"METRIC_PORT,required"`
	PprofPort  string `env:"PPROF_PORT,required"`
}

type Config struct {
	ServerConfig  ServerConfig
	ProductClient ProductClient
	LomsClient    LomsClient
	Observability Observability
}

func NewConfig() (Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
