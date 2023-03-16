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
	Dsn           string        `env:"DATABASE_URI,notEmpty" envDefault:"postgres://postgres:postgres@localhost:54323/postgres"`
	APISecret     string        `env:"API_SECRET,notEmpty" envDefault:"secret"`
	TokenLifespan time.Duration `env:"TOKEN_LIFESPAN,notEmpty" envDefault:"1h"`
	ServerAddr    string        `env:"RUN_ADDRESS,notEmpty" envDefault:":8080"`
}
