package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/smamykin/gofermart/internal/repository"
	"net/http"
)

func NewHealthcheckController(userRepository *repository.UserRepository) *HealthcheckController {
	return &HealthcheckController{
		UserRepository: userRepository,
	}
}

type HealthcheckController struct {
	UserRepository *repository.UserRepository
}

func (h *HealthcheckController) SetupRoutes(public *gin.RouterGroup, protected *gin.RouterGroup) {
	public.GET("/ping", h.HealthcheckHandler)
}

func (h *HealthcheckController) HealthcheckHandler(c *gin.Context) {
	type metric struct {
		DBError string
	}
	err := h.UserRepository.Healthcheck(c)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, metric{
			DBError: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, metric{})
}
