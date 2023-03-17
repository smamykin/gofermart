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
	AddOrder(order entity.Order) error
	GetOrderByOrderNumber(orderNumber int) error
}

type HashGeneratorInterface interface {
	Generate(stringToHash string) (string, error)
	IsEqual(hashedPassword string, plainTxtPwd string) (isValid bool, err error)
}

// endregion interfaces

// region errors

var ErrUserIsNotFound = errors.New("user is not found")
var ErrLoginIsNotValid = errors.New("login is incorrect")
var ErrPwdIsNotValid = errors.New("password is incorrect")

//endregion errors

//region DTO

type Credentials struct {
	Login string `json:"login"`
	Pwd   string `json:"password"`
}

//endregion DTO
