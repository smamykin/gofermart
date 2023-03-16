package routing

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smamykin/gofermart/internal/service"
	"github.com/smamykin/gofermart/pkg/contracts"
	"github.com/smamykin/gofermart/pkg/pwdhash"
	"github.com/smamykin/gofermart/pkg/token"
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

func (u *UserController) AddHandlers(public *gin.RouterGroup, protected *gin.RouterGroup) {
	//uc := UserController{}

	public.POST("/api/user/register", u.registerHandler)
	public.POST("/api/user/login", u.loginHandler)
	protected.POST("/api/user/orders", u.orderHandler)
}

func (u *UserController) registerHandler(c *gin.Context) {
	var credentials service.Credentials
	err := c.ShouldBindJSON(&credentials)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = u.userService.CreateNewUser(credentials)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
		return
	}

	if err, ok := err.(service.BadCredentialsError); ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u.logger.Err(err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	return
}

func (u *UserController) loginHandler(c *gin.Context) {
	var credentials service.Credentials
	err := c.ShouldBindJSON(&credentials)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.userService.GetUserIfPwdValid(credentials)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tkn, err := token.Generate(user.ID)
	if err != nil {
		u.logger.Err(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", tkn))
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (u *UserController) orderHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "success"})

}
