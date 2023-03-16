package service

import (
	"errors"
	"github.com/smamykin/gofermart/internal/entity"
)

// region contracts

// region interfaces

type StorageInterface interface {
	UpsertUser(login string, pwd string) error
	GetUserByLogin(login string) (entity.User, error)
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

//endregion contracts

type UserService struct {
	Storage       StorageInterface
	HashGenerator HashGeneratorInterface
}

func (u *UserService) CreateNewUser(credentials Credentials) error {
	if "" == credentials.Pwd {
		return ErrPwdIsNotValid
	}

	if "" == credentials.Login {
		return ErrLoginIsNotValid
	}

	_, err := u.Storage.GetUserByLogin(credentials.Login)
	if err == nil {
		// the user exists already
		return ErrLoginIsNotValid
	}

	if err != ErrUserIsNotFound {
		return err
	}

	pwdHash, err := u.HashGenerator.Generate(credentials.Pwd)
	if err != nil {
		return err
	}

	return u.Storage.UpsertUser(credentials.Login, pwdHash)
}

func (u *UserService) GetUserIfPwdValid(credentials Credentials) (user entity.User, err error) {

	user, err = u.Storage.GetUserByLogin(credentials.Login)
	if err != nil {
		return user, err
	}

	isValid, err := u.HashGenerator.IsEqual(user.Pwd, credentials.Pwd)
	if err != nil {
		return user, err
	}

	if !isValid {
		return user, ErrPwdIsNotValid
	}

	return user, nil
}
