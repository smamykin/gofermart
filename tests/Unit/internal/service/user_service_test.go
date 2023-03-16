package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/smamykin/gofermart/internal/entity"
	"github.com/smamykin/gofermart/internal/service"
	mock "github.com/smamykin/gofermart/tests/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserService_CreateNewUser(t *testing.T) {
	type testCase struct {
		credentials          service.Credentials
		upsertUserWillReturn []error
		expected             error
	}
	tests := map[string]testCase{
		"general case": {
			service.Credentials{Login: "cheesecake", Pwd: "pancake"},
			[]error{nil},
			nil,
		},
		"if password is empty": {
			service.Credentials{Login: "cheesecake", Pwd: ""},
			[]error{},
			service.ErrPwdIsNotValid,
		},
		"if login is empty": {
			service.Credentials{Login: "", Pwd: "pancake"},
			[]error{},
			service.ErrLoginIsNotValid,
		},
		"if storage returns error": {
			service.Credentials{Login: "cheesecake", Pwd: "pancake"},
			[]error{errors.New("some error")},
			errors.New("some error"),
		},
		"if user exists already": {
			service.Credentials{Login: "already_exists", Pwd: "pancake"},
			[]error{},
			service.ErrLoginIsNotValid,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			us := service.UserService{
				Storage:       createStorageInterfaceMockForUpsertUser(ctrl, tt.credentials, tt.upsertUserWillReturn),
				HashGenerator: createHashGeneratorInterfaceMock(ctrl, true),
			}
			err := us.CreateNewUser(tt.credentials)
			require.Equal(t, tt.expected, err)
		})
	}
}

func TestUserService_GetUserIfPwdValid(t *testing.T) {

	type testCase struct {
		credentials       service.Credentials
		expectedUser      entity.User
		errorToReturn     error
		expectedErr       error
		IsEqualWillReturn bool
	}

	credentials := service.Credentials{Login: "cheesecake", Pwd: "pancake"}
	pwdHash, _ := hashFuncForTest(credentials.Pwd)
	expectedUser := entity.User{
		ID:    22,
		Login: credentials.Login,
		Pwd:   pwdHash,
	}
	tests := map[string]testCase{
		"general case":  {credentials, expectedUser, nil, nil, true},
		"no user":       {credentials, entity.User{}, service.ErrUserIsNotFound, service.ErrUserIsNotFound, true},
		"pwd not valid": {credentials, expectedUser, nil, service.ErrPwdIsNotValid, false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			us := service.UserService{
				Storage:       createStorageInterfaceMockForGetUserIfPwdValid(ctrl, expectedUser, tt.errorToReturn),
				HashGenerator: createHashGeneratorInterfaceMock(ctrl, tt.IsEqualWillReturn),
			}

			actualUser, err := us.GetUserIfPwdValid(credentials)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedUser, actualUser)
		})
	}
}

func createHashGeneratorInterfaceMock(ctrl *gomock.Controller, IsEqualWillReturn bool) service.HashGeneratorInterface {
	m := mock.NewMockHashGeneratorInterface(ctrl)
	m.EXPECT().Generate(gomock.Any()).DoAndReturn(hashFuncForTest).AnyTimes()
	m.EXPECT().IsEqual(gomock.Any(), gomock.Any()).DoAndReturn(func(string, string) (bool, error) {
		return IsEqualWillReturn, nil
	}).AnyTimes()
	return m
}

func createStorageInterfaceMockForGetUserIfPwdValid(ctrl *gomock.Controller, user entity.User, errToReturn error) service.StorageInterface {
	m := mock.NewMockStorageInterface(ctrl)
	m.EXPECT().GetUserByLogin(gomock.Any()).DoAndReturn(func(login string) (u entity.User, err error) {
		if errToReturn != nil {
			return u, errToReturn
		}

		return user, nil
	}).AnyTimes()
	return m
}

func createStorageInterfaceMockForUpsertUser(ctrl *gomock.Controller, credentials service.Credentials, upsertUserWillReturn []error) service.StorageInterface {
	m := mock.NewMockStorageInterface(ctrl)
	pwdHash, _ := hashFuncForTest(credentials.Pwd)
	call := m.EXPECT().
		UpsertUser(gomock.Eq(credentials.Login), gomock.Eq(pwdHash)).
		Times(len(upsertUserWillReturn))

	for _, err := range upsertUserWillReturn {
		call.Return(err)
	}
	m.EXPECT().GetUserByLogin(gomock.Any()).DoAndReturn(func(login string) (u entity.User, err error) {
		if login == "already_exists" {
			return entity.User{
				ID:    22,
				Login: login,
				Pwd:   "not matter",
			}, nil
		}

		return u, service.ErrUserIsNotFound

	}).AnyTimes()

	return m
}

var hashFuncForTest = func(stringToHash string) (string, error) {
	return stringToHash + "'", nil
}
