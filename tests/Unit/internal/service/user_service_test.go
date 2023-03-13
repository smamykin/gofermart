package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/smamykin/gofermart/internal/entity"
	"github.com/smamykin/gofermart/internal/service"
	mock "github.com/smamykin/gofermart/tests/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserService_CreateNewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
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
			service.NewBadCredentialsError("password"),
		},
		"if login is empty": {
			service.Credentials{Login: "", Pwd: "pancake"},
			[]error{},
			service.NewBadCredentialsError("login"),
		},
		"if storage returns error": {
			service.Credentials{Login: "cheesecake", Pwd: "pancake"},
			[]error{errors.New("some error")},
			errors.New("some error"),
		},
		"if user exists already": {
			service.Credentials{Login: "already_exists", Pwd: "pancake"},
			[]error{},
			service.NewBadCredentialsError("login"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			us := service.UserService{
				Storage:       getStorageInterfaceMock(ctrl, tt.credentials, tt.upsertUserWillReturn),
				HashGenerator: getHashGeneratorInterfaceMock(ctrl),
			}
			err := us.CreateNewUser(tt.credentials)
			require.Equal(t, tt.expected, err)
		})
	}
}

var hashFuncForTest = func(stringToHash string) (string, error) {
	return stringToHash + "'", nil
}

func getHashGeneratorInterfaceMock(ctrl *gomock.Controller) service.HashGeneratorInterface {
	m := mock.NewMockHashGeneratorInterface(ctrl)
	m.EXPECT().Generate(gomock.Any()).DoAndReturn(hashFuncForTest).AnyTimes()
	return m
}

func getStorageInterfaceMock(ctrl *gomock.Controller, credentials service.Credentials, upsertUserWillReturn []error) service.StorageInterface {
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

		return u, service.ErrNoRows

	}).AnyTimes()

	return m
}
