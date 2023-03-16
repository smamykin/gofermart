package container

import (
	"github.com/smamykin/gofermart/pkg/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createRouter(controllers []controllerInterface) *gin.Engine {
	r := gin.Default()
	public := r.Group("/")
	protected := r.Group("/")
	protected.Use(jwtAuthMiddleware)

	for _, c := range controllers {
		c.SetupRoutes(public, protected)
	}

	return r
}

type controllerInterface interface {
	SetupRoutes(public *gin.RouterGroup, protected *gin.RouterGroup)
}

func jwtAuthMiddleware(c *gin.Context) {
	err := token.Valid(c)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}
	c.Next()
}
