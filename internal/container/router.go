package container

import (
	"github.com/smamykin/gofermart/pkg/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createRouter(controllers []controllerInterface, apiSecret []byte) *gin.Engine {
	r := gin.Default()
	public := r.Group("/")
	protected := r.Group("/")
	protected.Use(jwtAuthMiddleware(apiSecret))

	for _, c := range controllers {
		c.SetupRoutes(public, protected)
	}

	return r
}

type controllerInterface interface {
	SetupRoutes(public *gin.RouterGroup, protected *gin.RouterGroup)
}

func jwtAuthMiddleware(apiSecret []byte) func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := token.GetCurrentUserID(c, apiSecret)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Set("current_user_id", userID)
		c.Next()
	}
}
