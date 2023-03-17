package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/smamykin/gofermart/internal/entity"
	"github.com/smamykin/gofermart/internal/service"
	"github.com/smamykin/gofermart/pkg/contracts"
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
				OrderRepository: getOrderRepositoryMockForAddOrder(t, expectedOrder, service.ErrEntityIsNotFound, nil),
			},
			expectedOrder: expectedOrder,
			expectedErr:   nil,
		},
		"order already exists": {
			userID:      userID,
			orderNumber: orderNumber,
			sut: service.OrderService{
				OrderRepository: getOrderRepositoryMockForAddOrder(t, expectedOrder, nil, nil),
			},
			expectedOrder: expectedOrder,
			expectedErr:   service.ErrOrderAlreadyExists,
		},
		"unexpected error while getting the order": {
			userID:      userID,
			orderNumber: orderNumber,
			sut: service.OrderService{
				OrderRepository: getOrderRepositoryMockForAddOrder(t, expectedOrder, errors.New("some unexpected error"), nil),
			},
			expectedOrder: entity.Order{},
			expectedErr:   errors.New("some unexpected error"),
		},
		"unexpected error while add the order": {
			userID:      userID,
			orderNumber: orderNumber,
			sut: service.OrderService{
				OrderRepository: getOrderRepositoryMockForAddOrder(t, expectedOrder, service.ErrEntityIsNotFound, errors.New("some unexpected error")),
			},
			expectedOrder: entity.Order{},
			expectedErr:   errors.New("some unexpected error"),
		},
		"wrong order number (not Luhn's algorithm)": {
			userID:      userID,
			orderNumber: "123",
			sut: service.OrderService{
				OrderRepository: getOrderRepositoryMockForAddOrder(t, entity.Order{OrderNumber: "123"}, nil, nil),
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

func TestOrderService_UpdateOrdersStatuses(t *testing.T) {
	ordersFromDB := []entity.Order{
		{
			ID:            11,
			Status:        entity.OrderStatusNew,
			AccrualStatus: entity.AccrualStatusUndefined,
			OrderNumber:   "111",
		},
		{
			ID:            22,
			Status:        entity.OrderStatusProcessing,
			AccrualStatus: entity.AccrualStatusRegistered,
			OrderNumber:   "222",
		},
		{
			ID:            33,
			Status:        entity.OrderStatusProcessing,
			AccrualStatus: entity.AccrualStatusProcessing,
			OrderNumber:   "333",
		},
		{
			ID:            44,
			Status:        entity.OrderStatusProcessing,
			AccrualStatus: entity.AccrualStatusProcessing,
			OrderNumber:   "444",
		},
		{
			ID:            55,
			Status:        entity.OrderStatusNew,
			AccrualStatus: entity.AccrualStatusUndefined,
			OrderNumber:   "555",
		},
		{
			ID:            66, // the call of client mock on this order will return error
			Status:        entity.OrderStatusNew,
			AccrualStatus: entity.AccrualStatusUndefined,
			OrderNumber:   "666",
		},
	}
	ordersFromClient := []service.AccrualOrder{
		{
			Order:  "111",
			Status: entity.AccrualStatusRegistered,
		},
		{
			Order:  "222",
			Status: entity.AccrualStatusProcessing,
		},
		{
			Order:  "333",
			Status: entity.AccrualStatusInvalid,
		},
		{
			Order:   "444",
			Status:  entity.AccrualStatusProcessed,
			Accrual: 500,
		},
		{
			Order:  "555",
			Status: entity.AccrualStatusUnregistered,
		},
		{
			Order:  "666", // the call of client mock on this order will return error
			Status: 666,
		},
	}

	ordersExpectedToUpdate := []entity.Order{
		{
			ID:            11,
			Status:        entity.OrderStatusProcessing,
			AccrualStatus: entity.AccrualStatusRegistered,
			OrderNumber:   "111",
		},
		{
			ID:            22,
			Status:        entity.OrderStatusProcessing,
			AccrualStatus: entity.AccrualStatusProcessing,
			OrderNumber:   "222",
		},
		{
			ID:            33,
			Status:        entity.OrderStatusInvalid,
			AccrualStatus: entity.AccrualStatusInvalid,
			OrderNumber:   "333",
		},
		{
			ID:            44,
			Status:        entity.OrderStatusProcessed,
			AccrualStatus: entity.AccrualStatusProcessed,
			OrderNumber:   "444",
		},
		{
			ID:            55,
			Status:        entity.OrderStatusInvalid,
			AccrualStatus: entity.AccrualStatusUnregistered,
			OrderNumber:   "555",
		},
	}

	sut := service.OrderService{
		OrderRepository: getOrderRepositoryMockForUpdateOrderStatuses(t, ordersFromDB, ordersExpectedToUpdate),
		AccrualClient:   getAccrualClientMock(t, ordersFromClient),
		Logger:          getLogMock(t),
	}
	err := sut.UpdateOrdersStatuses()
	require.NoError(t, err)
}

func getLogMock(t *testing.T) contracts.LoggerInterface {
	ctrl := gomock.NewController(t)
	m := mock.NewMockLoggerInterface(ctrl)
	m.EXPECT().Err(gomock.Any()).AnyTimes()
	return m
}

func getAccrualClientMock(t *testing.T, ordersFromClient []service.AccrualOrder) service.AccrualClientInterface {
	ctrl := gomock.NewController(t)
	m := mock.NewMockAccrualClientInterface(ctrl)

	var calls []*gomock.Call
	for _, accrualOrder := range ordersFromClient {
		if accrualOrder.Status == 666 {
			call := m.EXPECT().GetOrder(gomock.Eq(accrualOrder.Order)).Return(accrualOrder, errors.New("some error"))
			calls = append(calls, call)
			continue
		}
		call := m.EXPECT().GetOrder(gomock.Eq(accrualOrder.Order)).Return(accrualOrder, nil)
		calls = append(calls, call)
	}
	gomock.InOrder(calls...)

	return m
}

func getOrderRepositoryMockForUpdateOrderStatuses(t *testing.T, ordersFromDB []entity.Order, ordersToUpdate []entity.Order) service.OrderRepositoryInterface {
	ctrl := gomock.NewController(t)
	m := mock.NewMockOrderRepositoryInterface(ctrl)

	var calls []*gomock.Call
	for _, orderToUpdate := range ordersToUpdate {
		call := m.EXPECT().UpdateOrder(gomock.Eq(orderToUpdate)).Return(orderToUpdate, nil)
		calls = append(calls, call)
	}
	gomock.InOrder(calls...)

	m.EXPECT().GetOrdersWithUnfinishedStatus().Return(ordersFromDB, nil)

	return m
}

func getOrderRepositoryMockForAddOrder(t *testing.T, expectedOrder entity.Order, getOrderWillReturnErr error, addOrderWillReturnErr error) service.OrderRepositoryInterface {
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
