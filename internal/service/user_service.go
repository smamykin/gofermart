package service

import (
	"errors"
	"fmt"
	"github.com/smamykin/gofermart/internal/entity"
)

type StorageInterface interface {
	UpsertUser(login string, pwd string) error
	GetUserByLogin(login string) (entity.User, error)
}

type HashGeneratorInterface interface {
	Generate(stringToHash string) (string, error)
	IsEqual(hashedPassword string, plainTxtPwd string) (isValid bool, err error)
}

func NewBadCredentialsError(fieldName string) error {
	return BadCredentialsError{fmt.Sprintf("%s is incorrect", fieldName)}
}

type BadCredentialsError struct {
	msg string
}

func (b BadCredentialsError) Error() string {
	return b.msg
}

var ErrUserNotFound = errors.New("user is not found")
var ErrPwdNotValid = errors.New("password is incorrect")

type UserService struct {
	Storage       StorageInterface
	HashGenerator HashGeneratorInterface
}

type Credentials struct {
	Login string `json:"login"`
	Pwd   string `json:"password"`
}

func (u *UserService) CreateNewUser(credentials Credentials) error {
	if "" == credentials.Pwd {
		return NewBadCredentialsError("password")
	}

	if "" == credentials.Login {
		return NewBadCredentialsError("login")
	}

	_, err := u.Storage.GetUserByLogin(credentials.Login)
	if err == nil {
		// the user exists already
		return NewBadCredentialsError("login")
	}

	if err != ErrUserNotFound {
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
		return user, ErrPwdNotValid
	}

	return user, nil
}
