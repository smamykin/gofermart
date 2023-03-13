package router

import (
	"github.com/rs/zerolog"
	"github.com/smamykin/gofermart/internal/controller"
	"github.com/smamykin/gofermart/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(storage *storage.DBStorage, logger *zerolog.Logger) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		type metric struct {
			dbError error
		}
		err := storage.Healthcheck(c)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, metric{
				dbError: err,
			})
			return
		}

		c.JSON(http.StatusOK, metric{
			dbError: nil,
		})
	})

	controller.NewUserController().AddHandlers(r)

	return r
}
