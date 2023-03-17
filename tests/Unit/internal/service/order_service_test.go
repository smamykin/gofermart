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

func TestOrderService_AddOrder(t *testing.T) {
	userID := 1
	orderNumber := "2755060072"
	expectedOrder := entity.Order{
		UserID:      userID,
		OrderNumber: orderNumber,
	}

	type testCase struct {
		userID        int
		orderNumber   string
		sut           service.OrderService
		expectedOrder entity.Order
		expectedErr   error
	}
	tests := map[string]testCase{
		"general case": {
			userID:      userID,
			orderNumber: orderNumber,
			sut: service.OrderService{
				OrderRepository: getOrderRepositoryMock(t, expectedOrder, service.ErrEntityIsNotFound, nil),
			},
			expectedOrder: expectedOrder,
			expectedErr:   nil,
		},
		"order already exists": {
			userID:      userID,
			orderNumber: orderNumber,
			sut: service.OrderService{
				OrderRepository: getOrderRepositoryMock(t, expectedOrder, nil, nil),
			},
			expectedOrder: expectedOrder,
			expectedErr:   service.ErrOrderAlreadyExists,
		},
		"unexpected error while getting the order": {
			userID:      userID,
			orderNumber: orderNumber,
			sut: service.OrderService{
				OrderRepository: getOrderRepositoryMock(t, expectedOrder, errors.New("some unexpected error"), nil),
			},
			expectedOrder: entity.Order{},
			expectedErr:   errors.New("some unexpected error"),
		},
		"unexpected error while add the order": {
			userID:      userID,
			orderNumber: orderNumber,
			sut: service.OrderService{
				OrderRepository: getOrderRepositoryMock(t, expectedOrder, service.ErrEntityIsNotFound, errors.New("some unexpected error")),
			},
			expectedOrder: entity.Order{},
			expectedErr:   errors.New("some unexpected error"),
		},
		"wrong order number (not Luhn's algorithm)": {
			userID:      userID,
			orderNumber: "123",
			sut: service.OrderService{
				OrderRepository: getOrderRepositoryMock(t, entity.Order{OrderNumber: "123"}, nil, nil),
			},
			expectedOrder: entity.Order{},
			expectedErr:   service.ErrInvalidOrderNumber,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actualOrder, actualErr := tt.sut.AddOrder(tt.userID, tt.orderNumber)

			require.Equal(t, tt.expectedErr, actualErr)
			require.Equal(t, tt.expectedOrder, actualOrder)
		})
	}
}

func getOrderRepositoryMock(t *testing.T, expectedOrder entity.Order, getOrderWillReturnErr error, addOrderWillReturnErr error) service.OrderRepositoryInterface {
	ctrl := gomock.NewController(t)
	m := mock.NewMockOrderRepositoryInterface(ctrl)

	call := m.EXPECT().AddOrder(gomock.Eq(expectedOrder))
	if getOrderWillReturnErr == service.ErrEntityIsNotFound {
		call.Times(1).Return(expectedOrder, addOrderWillReturnErr)

	} else {
		call.Times(0)
	}

	call = m.EXPECT().GetOrderByOrderNumber(gomock.Eq(expectedOrder.OrderNumber)).AnyTimes()
	call.Return(expectedOrder, getOrderWillReturnErr)

	return m
}
