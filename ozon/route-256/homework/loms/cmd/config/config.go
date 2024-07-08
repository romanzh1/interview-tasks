package config

import (
	"github.com/caarlos0/env/v11"
)

type ServerConfig struct {
	GRPCPort string `env:"GRPC_PORT,required"`
	HTTPPort string `env:"HTTP_PORT,required"`
}

type LomsDatabaseConfig struct {
	Host     string `env:"DB_HOST,required"`
	Port     string `env:"DB_PORT,required"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	Database string `env:"DB_DATABASE,required"`
}

type Observability struct {
	TracerHost string `env:"TRACER_HOST,required"`
	TracerPort string `env:"TRACER_PORT,required"`
	MetricPort string `env:"METRIC_PORT,required"`
}

type Config struct {
	ServerConfig  ServerConfig
	LomsDatabase  LomsDatabaseConfig
	Observability Observability
}

func NewConfig() (Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (c Config) LomsDatabaseURL() string {
	return "postgres://" + c.LomsDatabase.User + ":" + c.LomsDatabase.Password + "@" + c.LomsDatabase.Host + ":" + c.LomsDatabase.Port + "/" + c.LomsDatabase.Database
}
