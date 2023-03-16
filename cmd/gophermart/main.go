package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/smamykin/gofermart/internal/container"
	"os"
)

var logger = zerolog.New(os.Stdout)

func main() {
	var err error
	dsn, ok := os.LookupEnv("DATABASE_DSN")
	if !ok {
		logger.Error().Msgf("env variable DATABASE_DSN is not defined")
		return
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Error().Msgf("cannot open db connection. error: %s\n", err.Error())
		return
	}
	defer db.Close()

	// Listen and Server in 0.0.0.0:8080
	err = container.NewContainer(db, &logger).Router().Run(":8080")
	if err != nil {
		fmt.Println("error while running server")
		logger.Error().Msgf("error while running server. error: %s\n", err.Error())
	}
}
