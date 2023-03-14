package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/smamykin/gofermart/internal/routing"
	"github.com/smamykin/gofermart/internal/storage"
	"os"
)

var logger = zerolog.New(os.Stdout)

func main() {
	var err error
	//todo get the dsn from env var
	db, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:54323/postgres")
	if err != nil {
		logger.Error().Msgf("cannot open db connection. error: %s\n", err.Error())
		return
	}
	defer db.Close()

	dbStorage, err := storage.NewDBStorage(db)
	if err != nil {
		logger.Error().Msgf("cannot create db storage. error: %s\n", err.Error())
		return
	}

	// Listen and Server in 0.0.0.0:8080
	err = routing.SetupRouter(dbStorage, &logger).Run(":8080")
	if err != nil {
		fmt.Println("error while running server")
		logger.Error().Msgf("error while running server. error: %s\n", err.Error())
	}
}
