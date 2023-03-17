package service

import (
	"github.com/smamykin/gofermart/internal/entity"
)

type UserService struct {
	Storage       UserRepositoryInterface
	HashGenerator HashGeneratorInterface
}

func (u *UserService) CreateNewUser(credentials Credentials) (user entity.User, err error) {
	if credentials.Pwd == "" {
		return user, ErrPwdIsNotValid
	}

	if credentials.Login == "" {
		return user, ErrLoginIsNotValid
	}

	_, err = u.Storage.GetUserByLogin(credentials.Login)
	if err == nil {
		// the user exists already
		return user, ErrLoginIsNotValid
	}

	if err != ErrUserIsNotFound {
		return user, err
	}

	pwdHash, err := u.HashGenerator.Generate(credentials.Pwd)
	if err != nil {
		return user, err
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
