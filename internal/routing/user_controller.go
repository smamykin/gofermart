package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/smamykin/gofermart/internal/service"
	"github.com/smamykin/gofermart/pkg/contracts"
	"github.com/smamykin/gofermart/pkg/pwdhash"
	"net/http"
)

func NewUserController(store service.StorageInterface, logger contracts.LoggerInterface) *UserController {
	return &UserController{
		logger: logger,
		userService: service.UserService{
			Storage:       store,
			HashGenerator: &pwdhash.HashGenerator{},
		},
	}
}

type UserController struct {
	logger      contracts.LoggerInterface
	userService service.UserService
}

func (u *UserController) AddHandlers(engine *gin.Engine) {
	//uc := UserController{}

	engine.POST("/api/user/register", u.getRegisterHandler())
}

func (u *UserController) getRegisterHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		var credentials service.Credentials
		err := c.ShouldBindJSON(&credentials)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = u.userService.CreateNewUser(credentials)
		if err != nil {
			if err, ok := err.(service.BadCredentialsError); ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			u.logger.Err(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"endpoint": "register"})
	}
}
