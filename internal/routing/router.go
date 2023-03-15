package routing

import (
	"github.com/rs/zerolog"
	"github.com/smamykin/gofermart/internal/storage"
	"github.com/smamykin/gofermart/pkg/logger"
	"github.com/smamykin/gofermart/pkg/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(dbStorage *storage.DBStorage, zLogger *zerolog.Logger) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	public := r.Group("/")
	protected := r.Group("/")
	protected.Use(JwtAuthMiddleware)

	// Ping test
	public.GET("/ping", func(c *gin.Context) {
		type metric struct {
			DBError string
		}
		err := dbStorage.Healthcheck(c)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, metric{
				DBError: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, metric{})
	})

	NewUserController(dbStorage, &logger.ZeroLogAdapter{Logger: zLogger}).AddHandlers(public, protected)

	return r
}

func JwtAuthMiddleware(c *gin.Context) {
	err := token.TokenValid(c)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}
	c.Next()
}
