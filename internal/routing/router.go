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

	controllers := setupDependencies(db, zLogger)

	public := r.Group("/")
	protected := r.Group("/")
	protected.Use(JwtAuthMiddleware)

	for _, c := range controllers {
		c.SetupRoutes(public, protected)
	}

	return r
}

type ControllerInterface interface {
	SetupRoutes(public *gin.RouterGroup, protected *gin.RouterGroup)
}

func setupDependencies(db *sql.DB, zLogger *zerolog.Logger) []ControllerInterface {
	dbStorage := storage.NewDBStorage(db)

	return []ControllerInterface{
		controller.NewHealthcheckController(dbStorage),
		controller.NewUserController(
			&logger.ZeroLogAdapter{Logger: zLogger},
			service.UserService{
				Storage:       dbStorage,
				HashGenerator: &pwdhash.HashGenerator{},
			},
		),
	}
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
