package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewHealthcheckController(HealthcheckStorage func(ctx context.Context) error) *HealthcheckController {
	return &HealthcheckController{
		HealthcheckStorage: HealthcheckStorage,
	}
}

type HealthcheckController struct {
	HealthcheckStorage func(ctx context.Context) error
}

func (h *HealthcheckController) SetupRoutes(public *gin.RouterGroup, protected *gin.RouterGroup) {
	public.GET("/ping", h.HealthcheckHandler)
}

func (h *HealthcheckController) HealthcheckHandler(c *gin.Context) {
	type metric struct {
		DBError string
	}
	err := h.HealthcheckStorage(c)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, metric{
			DBError: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, metric{})
}
