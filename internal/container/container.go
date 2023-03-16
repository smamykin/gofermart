package container

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/smamykin/gofermart/internal/controller"
	"github.com/smamykin/gofermart/internal/service"
	"github.com/smamykin/gofermart/internal/storage"
	"github.com/smamykin/gofermart/pkg/logger"
	"github.com/smamykin/gofermart/pkg/pwdhash"
)

func NewContainer(db *sql.DB, zLogger *zerolog.Logger) *Container {
	dbStorage := storage.NewDBStorage(db)

	controllers := []controllerInterface{
		controller.NewHealthcheckController(dbStorage),
		controller.NewUserController(
			&logger.ZeroLogAdapter{Logger: zLogger},
			service.UserService{
				Storage:       dbStorage,
				HashGenerator: &pwdhash.HashGenerator{},
			},
		),
	}

	router := createRouter(controllers)

	return &Container{controllers, router}
}

type Container struct {
	controllers []controllerInterface
	router      *gin.Engine
}

func (c *Container) Controllers() []controllerInterface {
	return c.controllers
}

func (c *Container) Router() *gin.Engine {
	return c.router
}
