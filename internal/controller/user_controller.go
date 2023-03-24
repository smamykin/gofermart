package controller

import (
	"errors"
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
	withdrawalService *service.WithdrawalService,
	apiSecret []byte,
	tokenLifespan time.Duration,
) *UserController {
	return &UserController{
		logger:            logger,
		userService:       userService,
		orderService:      orderService,
		withdrawalService: withdrawalService,
		apiSecret:         apiSecret,
		tokenLifespan:     tokenLifespan,
	}
}

type UserController struct {
	logger            contracts.LoggerInterface
	userService       *service.UserService
	orderService      *service.OrderService
	withdrawalService *service.WithdrawalService
	apiSecret         []byte
	tokenLifespan     time.Duration
}

func (u *UserController) SetupRoutes(public *gin.RouterGroup, protected *gin.RouterGroup) {
	public.POST("/api/user/register", u.registerHandler)
	public.POST("/api/user/login", u.loginHandler)
	protected.POST("/api/user/orders", u.addOrderHandler)
	protected.GET("/api/user/orders", u.orderListHandler)
	protected.GET("/api/user/balance", u.balanceHandler)
	protected.POST("/api/user/balance/withdraw", u.withdrawHandler)
	protected.GET("/api/user/withdrawals", u.withdrawalListHandler)
}

func (u *UserController) registerHandler(c *gin.Context) {
	var credentials service.Credentials
	err := c.ShouldBindJSON(&credentials)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := u.userService.CreateNewUser(credentials)
	if errors.Is(err, service.ErrPwdIsNotValid) || errors.Is(err, service.ErrLoginIsNotValid) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		u.logger.Err(err, "unknown error while creating new user")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tkn, err := token.Generate(user.ID, u.apiSecret, u.tokenLifespan)
	if err != nil {
		u.logger.Err(err, "error occurred while generating an auth token after user's registration")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", tkn))
	c.JSON(http.StatusOK, gin.H{"message": "success"})
	return
}

func (u *UserController) loginHandler(c *gin.Context) {
	var credentials service.Credentials
	err := c.ShouldBindJSON(&credentials)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.userService.GetUserIfPwdValid(credentials)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tkn, err := token.Generate(user.ID, u.apiSecret, u.tokenLifespan)
	if err != nil {
		u.logger.Err(err, "error occurred while generating an auth token after user's login")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", tkn))
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (u *UserController) addOrderHandler(c *gin.Context) {
	currentUserID, err := GetCurrentUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"message": "cannot get currentUserID"})
	}
	body, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "cannot read body"})
		return
	}

	orderNumber := string(body)
	if orderNumber == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "cannot fetch an order number from the  body"})
		return
	}
	order, err := u.orderService.AddOrder(currentUserID, orderNumber)
	if errors.Is(err, service.ErrEntityAlreadyExists) {
		if order.UserID == currentUserID {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "order already exists"})
			return
		}

		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "order already exists"})
		return
	}

	if errors.Is(err, service.ErrInvalidOrderNumber) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	if err != nil {
		u.logger.Err(err, "unknown error while adding the order")
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, &order)
}

func (u *UserController) orderListHandler(c *gin.Context) {
	userID, err := GetCurrentUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"message": "cannot get currentUserID"})
	}
	orders, err := u.orderService.GetAllOrdersByUserID(userID)
	if err != nil {
		u.logger.Err(err, "error while getting orders of the user")
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

func (u *UserController) balanceHandler(c *gin.Context) {
	userID, err := GetCurrentUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"message": "cannot get currentUserID"})
	}

	balance, err := u.userService.GetBalance(userID)
	if err != nil {
		u.logger.Err(err, "error while getting user's balance")
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, balance)
}

func (u *UserController) withdrawHandler(c *gin.Context) {
	userID, err := GetCurrentUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"message": "cannot get currentUserID"})
	}
	var withdrawalRequestModel WithdrawalRequestModel
	err = c.ShouldBindJSON(&withdrawalRequestModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	withdraw, err := u.withdrawalService.Withdraw(userID, withdrawalRequestModel.Amount, withdrawalRequestModel.OrderNumber)
	if errors.Is(err, service.ErrEntityAlreadyExists) {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "order already exists"})
		return
	}

	if errors.Is(err, service.ErrInvalidOrderNumber) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	if err != nil {
		u.logger.Err(err, "unknown error while withdrawing")
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, withdraw)
}

func (u *UserController) withdrawalListHandler(c *gin.Context) {
	userID, err := GetCurrentUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"message": "cannot get currentUserID"})
	}

	withdrawals, err := u.withdrawalService.GetAllWithdrawalByUserID(userID)
	if err != nil {
		u.logger.Err(err, "error while getting all withdrawals of the user")
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		return
	}
	var responseModels []WithdrawalResponseModel
	for _, withdrawal := range withdrawals {
		responseModel := WithdrawalResponseModel{
			OrderNumber: withdrawal.OrderNumber,
			Amount:      withdrawal.Amount,
			ProcessedAt: withdrawal.CreatedAt.Format(time.RFC3339),
		}
		responseModels = append(responseModels, responseModel)
	}

	if len(withdrawals) == 0 {
		c.JSON(http.StatusNoContent, responseModels)
		return
	}

	c.JSON(http.StatusOK, responseModels)
}

type OrderResponseModel struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at" `
}
type WithdrawalResponseModel struct {
	OrderNumber string  `json:"order"`
	Amount      float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type WithdrawalRequestModel struct {
	OrderNumber string  `json:"order"`
	Amount      float64 `json:"sum"`
}
