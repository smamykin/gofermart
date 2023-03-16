package config

import (
	"github.com/caarlos0/env/v7"
	"time"
)

func NewConfigFromEnv() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

type Config struct {
	Dsn           string        `env:"DATABASE_DSN,notEmpty" envDefault:"postgres://postgres:postgres@localhost:54323/postgres"`
	ApiSecret     string        `env:"API_SECRET,notEmpty" envDefault:"secret"`
	TokenLifespan time.Duration `env:"TOKEN_LIFESPAN,notEmpty" envDefault:"1h"`
	ServerAddr    string        `env:"SERVER_ADDR,notEmpty" envDefault:"localhost:8080"`
}
