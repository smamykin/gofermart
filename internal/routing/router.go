package routing

import (
	"database/sql"
	"github.com/rs/zerolog"
	"github.com/smamykin/gofermart/internal/controller"
	"github.com/smamykin/gofermart/internal/service"
	"github.com/smamykin/gofermart/internal/storage"
	"github.com/smamykin/gofermart/pkg/logger"
	"github.com/smamykin/gofermart/pkg/pwdhash"
	"github.com/smamykin/gofermart/pkg/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB, zLogger *zerolog.Logger) *gin.Engine {
	r := gin.Default()

	public := r.Group("/")
	protected := r.Group("/")
	protected.Use(JwtAuthMiddleware)

	dbStorage := storage.NewDBStorage(db)

	h := controller.NewHealthcheckController(dbStorage)
	u := controller.NewUserController(
		&logger.ZeroLogAdapter{Logger: zLogger},
		service.UserService{
			Storage:       dbStorage,
			HashGenerator: &pwdhash.HashGenerator{},
		},
	)

	public.GET("/ping", h.PingHandler)
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
