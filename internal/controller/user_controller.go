package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smamykin/gofermart/internal/service"
	"github.com/smamykin/gofermart/pkg/contracts"
	"github.com/smamykin/gofermart/pkg/token"
	"net/http"
	"time"
)

func NewUserController(
	logger contracts.LoggerInterface,
	userService *service.UserService,
	orderService *service.OrderService,
	apiSecret []byte,
	tokenLifespan time.Duration,
) *UserController {
	return &UserController{
		logger:        logger,
		userService:   userService,
		orderService:  orderService,
		apiSecret:     apiSecret,
		tokenLifespan: tokenLifespan,
	}
}

type UserController struct {
	logger        contracts.LoggerInterface
	userService   *service.UserService
	orderService  *service.OrderService
	apiSecret     []byte
	tokenLifespan time.Duration
}

func (u *UserController) SetupRoutes(public *gin.RouterGroup, protected *gin.RouterGroup) {
	public.POST("/api/user/register", u.registerHandler)
	public.POST("/api/user/login", u.loginHandler)
	protected.POST("/api/user/orders", u.addOrderHandler)
	protected.GET("/api/user/orders", u.orderListHandler)
}

func (u *UserController) registerHandler(c *gin.Context) {
	var credentials service.Credentials
	err := c.ShouldBindJSON(&credentials)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.userService.CreateNewUser(credentials)
	if err == nil {
		tkn, err := token.Generate(user.ID, u.apiSecret, u.tokenLifespan)
		if err != nil {
			u.logger.Err(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Authorization", fmt.Sprintf("Bearer %s", tkn))
		c.JSON(http.StatusOK, gin.H{"message": "success"})
		return
	}
	if err == service.ErrPwdIsNotValid || err == service.ErrLoginIsNotValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u.logger.Err(err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	tkn, err := token.Generate(user.ID, u.apiSecret, u.tokenLifespan)
	if err != nil {
		u.logger.Err(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", tkn))
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (u *UserController) addOrderHandler(c *gin.Context) {
	currentUserID := getCurrentUserIDFromContext(c)
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "cannot read body"})
		return
	}

	orderNumber := string(body)
	if orderNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "cannot fetch an order number from the  body"})
		return
	}
	order, err := u.orderService.AddOrder(currentUserID, orderNumber)
	if err != nil {
		if err == service.ErrOrderAlreadyExists {
			if order.UserID == currentUserID {
				c.JSON(http.StatusOK, gin.H{"message": "order already exists"})
				return
			}

			c.JSON(http.StatusConflict, gin.H{"message": "order already exists"})
			return
		}

		if err == service.ErrInvalidOrderNumber {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}

		u.logger.Err(err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, &order)
}

func (u *UserController) orderListHandler(c *gin.Context) {
	userID := getCurrentUserIDFromContext(c)
	orders, err := u.orderService.GetAllOrdersByUserID(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		return
	}
	var orderResponseModels []OrderResponseModel
	for _, order := range orders {
		orderResponseModel := OrderResponseModel{
			Number:     order.OrderNumber,
			Status:     order.Status.String(),
			Accrual:    order.Accrual,
			UploadedAt: order.CreatedAt.Format(time.RFC3339),
		}
		orderResponseModels = append(orderResponseModels, orderResponseModel)
	}

	if len(orders) == 0 {
		c.JSON(http.StatusNoContent, orderResponseModels)
		return
	}

	c.JSON(http.StatusOK, orderResponseModels)
}

func getCurrentUserIDFromContext(c *gin.Context) int {
	currentUserIDAsAny, _ := c.Get("current_user_id")
	return currentUserIDAsAny.(int)
}

type OrderResponseModel struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at" `
}
