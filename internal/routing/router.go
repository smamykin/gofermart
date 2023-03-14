package routing

import (
	"github.com/rs/zerolog"
	"github.com/smamykin/gofermart/internal/storage"
	"github.com/smamykin/gofermart/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(dbStorage *storage.DBStorage, zLogger *zerolog.Logger) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		type metric struct {
			dbError error
		}
		err := dbStorage.Healthcheck(c)
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

	NewUserController(dbStorage, &logger.ZeroLogAdapter{Logger: zLogger}).AddHandlers(r)

	return r
}
