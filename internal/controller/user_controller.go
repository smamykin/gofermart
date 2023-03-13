package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/smamykin/gofermart/internal/service"
	"net/http"
)

func NewUserController() *UserController {
	return &UserController{
		userService: service.UserService{},
	}
}

type UserController struct {
	userService service.UserService
}

func (u *UserController) AddHandlers(engine *gin.Engine) {
	//uc := UserController{}

	engine.POST("/api/user/register", u.getRegisterHandler())
}

func (u *UserController) getRegisterHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		var credentials service.Credentials
		err := c.BindJSON(&credentials)
		if err != nil {
			return
		}
		u.userService.CreateNewUser(credentials)
		c.JSON(http.StatusOK, gin.H{"endpoint": "register"})
	}
}
