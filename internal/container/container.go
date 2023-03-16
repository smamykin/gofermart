package container

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/smamykin/gofermart/internal/config"
	"github.com/smamykin/gofermart/internal/controller"
	"github.com/smamykin/gofermart/internal/service"
	"github.com/smamykin/gofermart/internal/storage"
	"github.com/smamykin/gofermart/pkg/logger"
	"github.com/smamykin/gofermart/pkg/pwdhash"
)

var ErrClosedContainer = errors.New("the container is closed")

func NewContainer(zLogger *zerolog.Logger) (c *Container, err error) {

	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		return c, fmt.Errorf("error while getting configuration. error: %w", err)
	}

	db, err := sql.Open("pgx", cfg.Dsn)
	if err != nil {
		return c, fmt.Errorf("cannot open db connection. error: %w", err)
	}

	c = &Container{}
	c.config = cfg
	c.db = db
	c.storage = storage.NewDBStorage(c.db)
	c.controllers = []controllerInterface{
		controller.NewHealthcheckController(c.storage),
		controller.NewUserController(
			&logger.ZeroLogAdapter{Logger: zLogger},
			service.UserService{
				Storage:       c.storage,
				HashGenerator: &pwdhash.HashGenerator{},
			},
			[]byte(cfg.ApiSecret),
			cfg.TokenLifespan,
		),
	}
	c.router = createRouter(c.controllers, []byte(c.Config().ApiSecret))
	c.isOpen = true

	return c, nil
}

type Container struct {
	isOpen      bool
	config      config.Config
	controllers []controllerInterface
	db          *sql.DB
	router      *gin.Engine
	storage     *storage.DBStorage
}

func (c *Container) Controllers() []controllerInterface {
	return c.controllers
}

func (c *Container) Router() *gin.Engine {
	return c.router
}

func (c *Container) Config() config.Config {
	return c.config
}

func (c *Container) DB() *sql.DB {
	return c.db
}

func (c *Container) Storage() service.StorageInterface {
	return c.storage
}

func (c Container) IsOpen() bool {
	return c.isOpen
}

func (c *Container) Close() error {
	if !c.IsOpen() {
		return ErrClosedContainer
	}
	err := c.db.Close()
	if err != nil {
		return err
	}
	c.isOpen = false

	return nil
}
