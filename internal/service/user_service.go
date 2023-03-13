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

var ErrNoRows = errors.New("error no rows")

type UserService struct {
	Storage       StorageInterface
	HashGenerator HashGeneratorInterface
}

type Credentials struct {
	Login string `json:"login"`
	Pwd   string `json:"password"`
}

func (u UserService) CreateNewUser(credentials Credentials) error {
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

	if err != ErrNoRows {
		return err
	}

	pwdHash, err := u.HashGenerator.Generate(credentials.Pwd)
	if err != nil {
		return err
	}

	return u.Storage.UpsertUser(credentials.Login, pwdHash)
}
