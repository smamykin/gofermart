package service

import "fmt"

type StorageInterface interface {
	UpsertUser(login string, pwd string) error
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

	pwdHash, err := u.HashGenerator.Generate(credentials.Pwd)
	if err != nil {
		return err
	}

	return u.Storage.UpsertUser(credentials.Login, pwdHash)
}
