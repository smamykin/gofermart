package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/smamykin/gofermart/internal/service"
	mock "github.com/smamykin/gofermart/tests/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserService_CreateNewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	type testCase struct {
		credentials       service.Credentials
		storageWillReturn []error
		expected          error
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
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			us := service.UserService{
				Storage:       getStorageInterfaceMock(ctrl, tt.credentials, tt.storageWillReturn),
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

func getStorageInterfaceMock(ctrl *gomock.Controller, credentials service.Credentials, willReturn []error) service.StorageInterface {
	m := mock.NewMockStorageInterface(ctrl)
	//todo should get the hash of the pwd
	pwdHash, _ := hashFuncForTest(credentials.Pwd)
	call := m.EXPECT().
		UpsertUser(gomock.Eq(credentials.Login), gomock.Eq(pwdHash)).
		Times(len(willReturn))

	for _, err := range willReturn {
		call.Return(err)
	}

	return m
}
