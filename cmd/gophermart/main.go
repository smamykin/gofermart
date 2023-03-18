package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/smamykin/gofermart/internal/container"
	"github.com/smamykin/gofermart/pkg/utils"
	"os"
)

var logger = zerolog.New(os.Stdout)

func main() {

	c, err := container.NewContainer(&logger)
	if err != nil {
		logger.Error().Err(err).Msgf("error while building container")
		return
	}
	defer c.Close()

	go utils.InvokeFunctionWithInterval(c.Config().UpdateStatusInterval, func() {
		err := c.OrderService().UpdateOrdersStatuses()
		if err != nil {
			logger.Error().Err(err).Msgf("cannot update status")
		}
	})

	err = c.Router().Run(c.Config().ServerAddr)
	if err != nil {
		logger.Error().Msgf("error while running server. error: %s\n", err.Error())
	}
}
