package container

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/smamykin/gofermart/internal/client"
	"github.com/smamykin/gofermart/internal/config"
	"github.com/smamykin/gofermart/internal/controller"
	"github.com/smamykin/gofermart/internal/repository"
	"github.com/smamykin/gofermart/internal/service"
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
	err = ensureSchemaExists(db)
	if err != nil {
		return c, fmt.Errorf("cannot create db schema. error: %w", err)
	}

	c = &Container{}
	c.config = cfg
	c.db = db
	c.userRepository = repository.NewUserRepository(c.DB())
	c.orderRepository = repository.NewOrderRepository(c.DB())
	c.withdrawalRepository = repository.NewWithdrawalRepository(c.DB())
	APISecret := []byte(c.Config().APISecret)
	c.logger = &logger.ZeroLogAdapter{Logger: zLogger}
	c.orderService = &service.OrderService{
		OrderRepository: c.OrderRepository(),
		AccrualClient:   client.NewAccrualClient(c.Config().AccrualEntrypoint),
		Logger:          c.logger,
	}
	c.controllers = []controllerInterface{
		controller.NewHealthcheckController(repository.CreateHealthcheckFunc(c.DB())),
		controller.NewUserController(
			c.logger,
			&service.UserService{
				UserRepository:       c.UserRepository(),
				OrderRepository:      c.OrderRepository(),
				WithdrawalRepository: c.WithdrawalRepository(),
				HashGenerator:        &pwdhash.HashGenerator{},
			},
			c.OrderService(),
			APISecret,
			c.Config().TokenLifespan,
		),
	}
	c.router = createRouter(c.controllers, APISecret)
	c.isOpen = true

	return c, nil
}

type Container struct {
	isOpen               bool
	config               config.Config
	controllers          []controllerInterface
	db                   *sql.DB
	router               *gin.Engine
	userRepository       *repository.UserRepository
	orderRepository      *repository.OrderRepository
	orderService         *service.OrderService
	logger               *logger.ZeroLogAdapter
	withdrawalRepository *repository.WithdrawalRepository
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

func (c *Container) UserRepository() service.UserRepositoryInterface {
	return c.userRepository
}

func (c *Container) OrderRepository() service.OrderRepositoryInterface {
	return c.orderRepository
}

func (c *Container) WithdrawalRepository() service.WithdrawalRepositoryInterface {
	return c.withdrawalRepository
}

func (c *Container) OrderService() *service.OrderService {
	return c.orderService
}

func (c *Container) IsOpen() bool {
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

func ensureSchemaExists(db *sql.DB) error {
	tableExistsSQL := "SELECT EXISTS ( SELECT FROM pg_tables WHERE tablename  = 'user');"
	var isTableExists bool
	err := db.QueryRow(tableExistsSQL).Scan(&isTableExists)
	if err != nil {
		return err
	}
	if isTableExists {
		return nil
	}

	_, err = db.Exec(`
		--- User
		CREATE TABLE "user" (
			"id" SERIAL PRIMARY KEY,
			"login" VARCHAR NOT NULL ,
			"pwd" VARCHAR NOT NULL
		);

		CREATE UNIQUE INDEX udx_login ON "user" (login);

		--- Order
		CREATE TABLE "order" (
		    "id" SERIAL PRIMARY KEY,
			"user_id" INTEGER NOT NULL,
		    "order_number" VARCHAR NOT NULL,
		    "status" INTEGER NOT NULL,
		    "accrual_status" INTEGER NOT NULL,
		    "accrual" DOUBLE PRECISION NOT NULL,
		    "created_at" TIMESTAMP NOT NULL,
			CONSTRAINT fk_order_user
			    FOREIGN KEY(user_id) 
			    REFERENCES "user"(id)
		);

		CREATE UNIQUE INDEX udx_order_number ON "order" (order_number);
		CREATE INDEX idx_order_user_id ON "order" (user_id);

		--- Withdrawal
		CREATE TABLE "withdrawal" (
		    "id" SERIAL PRIMARY KEY,
			"user_id" INTEGER NOT NULL,
		    "order_number" VARCHAR NOT NULL,
		    "amount" DOUBLE PRECISION NOT NULL,
		    "created_at" TIMESTAMP NOT NULL,
			CONSTRAINT fk_withdrawal_user
			    FOREIGN KEY(user_id) 
			    REFERENCES "user"(id)
		);

		CREATE UNIQUE INDEX udx_withdrawal_number ON "order" (order_number);
		CREATE INDEX idx_withdrawal_user_id ON "withdrawal" (user_id);
	`)

	if err != nil {
		return err
	}

	return nil
}
