package service

import (
	"errors"
	"github.com/smamykin/gofermart/internal/entity"
)

// region interfaces

type UserRepositoryInterface interface {
	UpsertUser(login, pwd string) (user entity.User, err error)
	GetUserByLogin(login string) (entity.User, error)
}

type OrderRepositoryInterface interface {
	AddOrder(o entity.Order) (order entity.Order, err error)
	GetOrder(ID int) (order entity.Order, err error)
	GetOrderByOrderNumber(orderNumber string) (order entity.Order, err error)
	GetAllByUserID(userID int) ([]entity.Order, error)
	UpdateOrder(order entity.Order) (entity.Order, error)
	GetOrdersWithUnfinishedStatus() ([]entity.Order, error)
	GetAccrualSumByUserID(userID int) (sum float64, err error)
}

type WithdrawalRepositoryInterface interface {
	GetAmountSumByUserID(userID int) (sum float64, err error)
	AddWithdrawal(withdrawal entity.Withdrawal) (entity.Withdrawal, error)
	GetWithdrawal(ID int) (order entity.Withdrawal, err error)
	GetWithdrawalByOrderNumber(orderNumber string) (order entity.Withdrawal, err error)
}

type AccrualClientInterface interface {
	GetOrder(orderNumber string) (AccrualOrder, error)
}

type HashGeneratorInterface interface {
	Generate(stringToHash string) (string, error)
	IsEqual(hashedPassword string, plainTxtPwd string) (isValid bool, err error)
}

// endregion interfaces

// region errors

var ErrEntityIsNotFound = errors.New("entity is not found")
var ErrLoginIsNotValid = errors.New("login is incorrect")
var ErrPwdIsNotValid = errors.New("password is incorrect")
var ErrEntityAlreadyExists = errors.New("entity already exists")
var ErrInvalidOrderNumber = errors.New("order number is invalid")
var ErrNotEnoughAccrual = errors.New("not enough bonuses on the account")

//endregion errors

//region DTO

type Credentials struct {
	Login string `json:"login"`
	Pwd   string `json:"password"`
}
type AccrualOrder struct {
	Order   string
	Status  entity.AccrualStatus
	Accrual float64
}
type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

//endregion DTO
