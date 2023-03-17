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
var ErrOrderAlreadyExists = errors.New("order already exists")

//endregion errors

//region DTO

type Credentials struct {
	Login string `json:"login"`
	Pwd   string `json:"password"`
}

//endregion DTO
