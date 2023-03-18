package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/smamykin/gofermart/internal/entity"
	"github.com/smamykin/gofermart/internal/service"
	mock "github.com/smamykin/gofermart/tests/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func TestUserService_CreateNewUser(t *testing.T) {
	type testCase struct {
		credentials          service.Credentials
		expectedUser         entity.User
		upsertUserWillReturn []error
		expectedErr          error
	}
	tests := map[string]testCase{
		"general case": {
			service.Credentials{Login: "cheesecake", Pwd: "pancake"},
			entity.User{ID: rand.Int(), Login: "cheesecake", Pwd: "pancake"},
			[]error{nil},
			nil,
		},
		"if password is empty": {
			service.Credentials{Login: "cheesecake", Pwd: ""},
			entity.User{},
			[]error{},
			service.ErrPwdIsNotValid,
		},
		"if login is empty": {
			service.Credentials{Login: "", Pwd: "pancake"},
			entity.User{},
			[]error{},
			service.ErrLoginIsNotValid,
		},
		"if storage returns error": {
			service.Credentials{Login: "cheesecake", Pwd: "pancake"},
			entity.User{},
			[]error{errors.New("some error")},
			errors.New("some error"),
		},
		"if user exists already": {
			service.Credentials{Login: "already_exists", Pwd: "pancake"},
			entity.User{},
			[]error{},
			service.ErrLoginIsNotValid,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			us := service.UserService{
				UserRepository: createStorageInterfaceMockForUpsertUser(ctrl, tt.expectedUser, tt.credentials, tt.upsertUserWillReturn),
				HashGenerator:  createHashGeneratorInterfaceMock(ctrl, true),
			}
			actualUser, actualErr := us.CreateNewUser(tt.credentials)
			require.Equal(t, tt.expectedErr, actualErr)
			require.Equal(t, tt.expectedUser, actualUser)
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
		"no user":       {credentials, entity.User{}, service.ErrEntityIsNotFound, service.ErrEntityIsNotFound, true},
		"pwd not valid": {credentials, expectedUser, nil, service.ErrPwdIsNotValid, false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			us := service.UserService{
				UserRepository: createStorageInterfaceMockForGetUserIfPwdValid(ctrl, expectedUser, tt.errorToReturn),
				HashGenerator:  createHashGeneratorInterfaceMock(ctrl, tt.IsEqualWillReturn),
			}

			actualUser, err := us.GetUserIfPwdValid(credentials)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedUser, actualUser)
		})
	}
}

func TestUserService_GetBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	userID := rand.Int()
	accrualSum := 50.5
	withdrawalSum := 25.3
	expectedBalance := service.Balance{
		Current:   accrualSum - withdrawalSum,
		Withdrawn: withdrawalSum,
	}
	us := service.UserService{
		OrderRepository:      createOrderRepositoryMock(ctrl, userID, accrualSum, nil),
		WithdrawalRepository: createWithdrawalRepositoryMock(ctrl, userID, withdrawalSum, nil),
	}
	actualBalance, err := us.GetBalance(userID)
	require.NoError(t, err)
	require.Equal(t, expectedBalance, actualBalance)
}

func createWithdrawalRepositoryMock(ctrl *gomock.Controller, userID int, willReturnSum float64, willReturnErr error) service.WithdrawalRepositoryInterface {
	m := mock.NewMockWithdrawalRepositoryInterface(ctrl)
	m.EXPECT().GetAmountSumByUserId(gomock.Eq(userID)).Return(willReturnSum, willReturnErr).AnyTimes()
	return m
}

func createOrderRepositoryMock(ctrl *gomock.Controller, userID int, willReturnSum float64, willReturnErr error) service.OrderRepositoryInterface {
	m := mock.NewMockOrderRepositoryInterface(ctrl)
	m.EXPECT().GetAccrualSumByUserId(gomock.Eq(userID)).Return(willReturnSum, willReturnErr).AnyTimes()
	return m
}

func createHashGeneratorInterfaceMock(ctrl *gomock.Controller, IsEqualWillReturn bool) service.HashGeneratorInterface {
	m := mock.NewMockHashGeneratorInterface(ctrl)
	m.EXPECT().Generate(gomock.Any()).DoAndReturn(hashFuncForTest).AnyTimes()
	m.EXPECT().IsEqual(gomock.Any(), gomock.Any()).DoAndReturn(func(string, string) (bool, error) {
		return IsEqualWillReturn, nil
	}).AnyTimes()
	return m
}

func createStorageInterfaceMockForGetUserIfPwdValid(ctrl *gomock.Controller, user entity.User, errToReturn error) service.UserRepositoryInterface {
	m := mock.NewMockUserRepositoryInterface(ctrl)
	m.EXPECT().GetUserByLogin(gomock.Any()).DoAndReturn(func(login string) (u entity.User, err error) {
		if errToReturn != nil {
			return u, errToReturn
		}

		return user, nil
	}).AnyTimes()
	return m
}

func createStorageInterfaceMockForUpsertUser(ctrl *gomock.Controller, user entity.User, credentials service.Credentials, upsertUserWillReturn []error) service.UserRepositoryInterface {
	m := mock.NewMockUserRepositoryInterface(ctrl)
	pwdHash, _ := hashFuncForTest(credentials.Pwd)
	call := m.EXPECT().
		UpsertUser(gomock.Eq(credentials.Login), gomock.Eq(pwdHash)).
		Times(len(upsertUserWillReturn))

	for _, err := range upsertUserWillReturn {
		call.Return(user, err)
	}
	m.EXPECT().GetUserByLogin(gomock.Any()).DoAndReturn(func(login string) (u entity.User, err error) {
		if login == "already_exists" {
			return user, nil
		}

		return u, service.ErrEntityIsNotFound

	}).AnyTimes()

	return m
}

var hashFuncForTest = func(stringToHash string) (string, error) {
	return stringToHash + "'", nil
}
