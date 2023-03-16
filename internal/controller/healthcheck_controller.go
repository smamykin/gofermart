package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/smamykin/gofermart/internal/storage"
	"net/http"
)

func NewHealthcheckController(DBStorage *storage.DBStorage) *HealthcheckController {
	return &HealthcheckController{
		DBStorage: DBStorage,
	}
}

type HealthcheckController struct {
	DBStorage *storage.DBStorage
}

func (h *HealthcheckController) SetupRoutes(public *gin.RouterGroup, protected *gin.RouterGroup) {
	public.GET("/ping", h.HealthcheckHandler)
}

func (h *HealthcheckController) HealthcheckHandler(c *gin.Context) {
	type metric struct {
		DBError string
	}
	err := h.DBStorage.Healthcheck(c)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, metric{
			DBError: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, metric{})
}
