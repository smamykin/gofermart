package repository

import (
	"github.com/smamykin/gofermart/internal/entity"
	"github.com/smamykin/gofermart/internal/service"
	"github.com/smamykin/gofermart/tests/Functional/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestOrderRepository_AddOrder(t *testing.T) {
	c := utils.GetContainer(t)
	db := c.DB()
	utils.TruncateTable(t, db)

	user := utils.InsertUser(t, db, entity.User{})
	sut := c.OrderRepository()
	expectedOrder := entity.Order{
		UserID:        user.ID,
		OrderNumber:   "1",
		Status:        entity.OrderStatusNew,
		AccrualStatus: entity.AccrualStatusUndefined,
		Accrual:       0,
	}
	timeMin := time.Now()
	actualOrder, err := sut.AddOrder(expectedOrder)
	timeMax := time.Now()

	expectedOrder.ID = 1 // because the table was truncated beforehand, we can guess the id

	orderFromDB, err := sut.GetOrder(expectedOrder.ID)
	require.NoError(t, err)

	assertOrder(t, expectedOrder, orderFromDB, timeMin, timeMax)
	assertOrder(t, expectedOrder, actualOrder, timeMin, timeMax)

	require.NoError(t, err)
}

func TestOrderRepository_GetOrderByOrderNumber(t *testing.T) {
	c := utils.GetContainer(t)
	db := c.DB()
	utils.TruncateTable(t, db)

	user := utils.InsertUser(t, db, entity.User{})
	sut := c.OrderRepository()
	expectedOrder, err := sut.AddOrder(entity.Order{
		UserID:        user.ID,
		OrderNumber:   "1",
		Status:        entity.OrderStatusNew,
		AccrualStatus: entity.AccrualStatusUndefined,
		Accrual:       0,
	})
	require.NoError(t, err)

	actualOrder, err := sut.GetOrderByOrderNumber(expectedOrder.OrderNumber)
	require.NoError(t, err)

	require.Equal(t, expectedOrder, actualOrder)

	_, err = sut.GetOrderByOrderNumber("unknown order number")
	require.Equal(t, err, service.ErrEntityIsNotFound)
}

func TestOrderRepository_GetAllByUserID(t *testing.T) {
	c := utils.GetContainer(t)
	db := c.DB()
	utils.TruncateTable(t, db)

	userToGet := utils.InsertUser(t, db, entity.User{})
	userNotToGet := utils.InsertUser(t, db, entity.User{})
	sut := c.OrderRepository()
	orderToGet, err := c.OrderRepository().AddOrder(entity.Order{
		UserID:      userToGet.ID,
		OrderNumber: "123",
	})
	require.NoError(t, err)
	_, err = c.OrderRepository().AddOrder(entity.Order{
		UserID:      userNotToGet.ID,
		OrderNumber: "321",
	})
	require.NoError(t, err)

	actualOrders, err := sut.GetAllByUserID(userToGet.ID)
	require.NoError(t, err)

	require.Equal(t, []entity.Order{orderToGet}, actualOrders)

	actualOrders, err = sut.GetAllByUserID(999)
	require.Equal(t, []entity.Order{}, actualOrders)
	require.Equal(t, err, nil)
}

func assertOrder(t *testing.T, expected entity.Order, actual entity.Order, createAtMin time.Time, createAtMax time.Time) {
	require.WithinRange(t, actual.CreatedAt, createAtMin.Truncate(time.Second), createAtMax)
	now := time.Time{}
	expected.CreatedAt = now
	actual.CreatedAt = now
	require.Equal(t, expected, actual)
}
