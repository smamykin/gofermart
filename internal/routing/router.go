package routing

import (
	"github.com/rs/zerolog"
	"github.com/smamykin/gofermart/internal/service"
	"github.com/smamykin/gofermart/internal/storage"
	"github.com/smamykin/gofermart/pkg/logger"
	"github.com/smamykin/gofermart/pkg/pwdhash"
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

	u := NewUserController(
		&logger.ZeroLogAdapter{Logger: zLogger},
		service.UserService{
			Storage:       dbStorage,
			HashGenerator: &pwdhash.HashGenerator{},
		},
	)
	public.POST("/api/user/register", u.RegisterHandler)
	public.POST("/api/user/login", u.LoginHandler)
	protected.POST("/api/user/orders", u.OrderHandler)

	return r
}

func JwtAuthMiddleware(c *gin.Context) {
	err := token.Valid(c)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}
	c.Next()
}
