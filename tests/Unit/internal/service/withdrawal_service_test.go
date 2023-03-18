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

func TestWithdrawalService_Withdraw(t *testing.T) {
	userID := 1
	orderNumber := "2755060072"
	amount := 300.0
	expectedWithdrawal := entity.Withdrawal{
		UserID:      userID,
		OrderNumber: orderNumber,
		Amount:      amount,
	}

	type testCase struct {
		userID             int
		orderNumber        string
		Amount             float64
		sut                service.WithdrawalService
		expectedWithdrawal entity.Withdrawal
		expectedErr        error
	}
	tests := map[string]testCase{
		"general case": {
			userID:      userID,
			orderNumber: orderNumber,
			Amount:      amount,
			sut: service.WithdrawalService{
				WithdrawalRepository: getWithdrawalRepositoryMock(t, expectedWithdrawal, service.ErrEntityIsNotFound, nil, 0),
				OrderRepository:      getOrderRepositoryMockForWithdraw(t, userID, amount+1),
			},
			expectedWithdrawal: expectedWithdrawal,
			expectedErr:        nil,
		},
		"order already exists": {
			userID:      userID,
			orderNumber: orderNumber,
			Amount:      amount,
			sut: service.WithdrawalService{
				WithdrawalRepository: getWithdrawalRepositoryMock(t, expectedWithdrawal, nil, nil, 0),
				OrderRepository:      getOrderRepositoryMockForWithdraw(t, userID, amount+1),
			},
			expectedWithdrawal: expectedWithdrawal,
			expectedErr:        service.ErrEntityAlreadyExists,
		},
		"unexpected error while getting the order": {
			userID:      userID,
			orderNumber: orderNumber,
			Amount:      amount,
			sut: service.WithdrawalService{
				WithdrawalRepository: getWithdrawalRepositoryMock(t, expectedWithdrawal, errors.New("some unexpected error"), nil, 0),
				OrderRepository:      getOrderRepositoryMockForWithdraw(t, userID, amount+1),
			},
			expectedWithdrawal: entity.Withdrawal{},
			expectedErr:        errors.New("some unexpected error"),
		},
		"unexpected error while add the order": {
			userID:      userID,
			orderNumber: orderNumber,
			Amount:      amount,
			sut: service.WithdrawalService{
				WithdrawalRepository: getWithdrawalRepositoryMock(t, expectedWithdrawal, service.ErrEntityIsNotFound, errors.New("some unexpected error"), 0),
				OrderRepository:      getOrderRepositoryMockForWithdraw(t, userID, amount+1),
			},
			expectedWithdrawal: entity.Withdrawal{},
			expectedErr:        errors.New("some unexpected error"),
		},
		"wrong order number (not Luhn's algorithm)": {
			userID:      userID,
			orderNumber: "123",
			Amount:      amount,
			sut: service.WithdrawalService{
				WithdrawalRepository: getWithdrawalRepositoryMock(t, entity.Withdrawal{OrderNumber: "123", Amount: amount}, nil, nil, 0),
				OrderRepository:      getOrderRepositoryMockForWithdraw(t, userID, amount+1),
			},
			expectedWithdrawal: entity.Withdrawal{},
			expectedErr:        service.ErrInvalidOrderNumber,
		},
		"there is no enough accrual": {
			userID:      userID,
			orderNumber: orderNumber,
			Amount:      amount,
			sut: service.WithdrawalService{
				WithdrawalRepository: getWithdrawalRepositoryMock(t, expectedWithdrawal, service.ErrEntityIsNotFound, nil, 0),
				OrderRepository:      getOrderRepositoryMockForWithdraw(t, userID, amount-1),
			},
			expectedWithdrawal: entity.Withdrawal{},
			expectedErr:        service.ErrNotEnoughAccrual,
		},
		"there is too many withdrawals": {
			userID:      userID,
			orderNumber: orderNumber,
			Amount:      amount,
			sut: service.WithdrawalService{
				WithdrawalRepository: getWithdrawalRepositoryMock(t, expectedWithdrawal, service.ErrEntityIsNotFound, nil, amount),
				OrderRepository:      getOrderRepositoryMockForWithdraw(t, userID, amount+1),
			},
			expectedWithdrawal: entity.Withdrawal{},
			expectedErr:        service.ErrNotEnoughAccrual,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actualWithdrawal, actualErr := tt.sut.Withdraw(tt.userID, tt.Amount, tt.orderNumber)

			require.Equal(t, tt.expectedErr, actualErr)
			require.Equal(t, tt.expectedWithdrawal, actualWithdrawal)
		})
	}
}

func getWithdrawalRepositoryMock(
	t *testing.T,
	withdrawal entity.Withdrawal,
	getWithdrawalWillReturnErr error,
	addWithdrawalWillReturnError error,
	getAmountSumByUserIdWillReturn float64,
) service.WithdrawalRepositoryInterface {
	ctrl := gomock.NewController(t)
	m := mock.NewMockWithdrawalRepositoryInterface(ctrl)

	call := m.EXPECT().AddWithdrawal(gomock.Eq(withdrawal)).AnyTimes().Return(withdrawal, addWithdrawalWillReturnError)

	call = m.EXPECT().GetWithdrawalByOrderNumber(gomock.Eq(withdrawal.OrderNumber)).AnyTimes()
	call.Return(withdrawal, getWithdrawalWillReturnErr)

	m.EXPECT().GetAmountSumByUserID(gomock.Eq(withdrawal.UserID)).AnyTimes().Return(getAmountSumByUserIdWillReturn, nil)

	return m
}

func getOrderRepositoryMockForWithdraw(
	t *testing.T,
	userID int,
	getAmountSumByUserIdWillReturn float64,
) service.OrderRepositoryInterface {
	ctrl := gomock.NewController(t)
	m := mock.NewMockOrderRepositoryInterface(ctrl)

	m.EXPECT().GetAccrualSumByUserID(gomock.Eq(userID)).AnyTimes().Return(getAmountSumByUserIdWillReturn, nil)

	return m
}
