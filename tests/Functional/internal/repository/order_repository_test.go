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
	require.NoError(t, err)

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

func TestOrderRepository_GetOrdersWithUnfinishedStatus(t *testing.T) {
	c := utils.GetContainer(t)
	db := c.DB()
	utils.TruncateTable(t, db)

	sut := c.OrderRepository()
	//first check if there is no orders
	actualOrders, err := sut.GetOrdersWithUnfinishedStatus()
	require.NoError(t, err)

	require.Equal(t, []entity.Order{}, actualOrders)

	//second check if there are orders
	user := utils.InsertUser(t, db, entity.User{})
	order0, err := c.OrderRepository().AddOrder(entity.Order{
		UserID:      user.ID,
		OrderNumber: "123",
	})
	require.NoError(t, err)
	order1, err := c.OrderRepository().AddOrder(entity.Order{
		UserID:      user.ID,
		OrderNumber: "321",
	})
	require.NoError(t, err)

	actualOrders, err = sut.GetOrdersWithUnfinishedStatus()
	require.NoError(t, err)

	require.Equal(t, []entity.Order{order0, order1}, actualOrders)
}

func TestOrderRepository_UpdateOrder(t *testing.T) {
	c := utils.GetContainer(t)
	db := c.DB()
	utils.TruncateTable(t, db)

	sut := c.OrderRepository()
	user := utils.InsertUser(t, db, entity.User{})
	order, err := c.OrderRepository().AddOrder(entity.Order{
		UserID:      user.ID,
		OrderNumber: "123",
	})
	require.NoError(t, err)

	//first check if there are orders
	actualOrders, err := sut.UpdateOrder(order)
	require.NoError(t, err)

	require.Equal(t, order, actualOrders)

	//second check if there is no orders
	_, err = sut.UpdateOrder(entity.Order{
		UserID:      user.ID,
		OrderNumber: "99999",
	})

	require.Equal(t, err, service.ErrEntityIsNotFound)
}

func TestOrderRepository_GetAccrualSumByUserId(t *testing.T) {
	c := utils.GetContainer(t)
	db := c.DB()
	utils.TruncateTable(t, db)

	user := utils.InsertUser(t, db, entity.User{})
	anotherUser := utils.InsertUser(t, db, entity.User{})

	sut := c.OrderRepository()

	//first check if there is no withdrawals at all
	actualSum, err := sut.GetAccrualSumByUserId(user.ID)
	require.NoError(t, err)
	require.Equal(t, .0, actualSum)

	//second check if there are withdrawals
	order0, err := sut.AddOrder(entity.Order{
		UserID:      user.ID,
		OrderNumber: "111",
		Accrual:     11.1,
	})
	require.NoError(t, err)
	order1, err := sut.AddOrder(entity.Order{
		UserID:      user.ID,
		OrderNumber: "222",
		Accrual:     22.2,
	})
	require.NoError(t, err)

	_, err = sut.AddOrder(entity.Order{
		UserID:      anotherUser.ID,
		OrderNumber: "333",
		Accrual:     33.3,
	})
	require.NoError(t, err)

	expectedSum := order0.Accrual + order1.Accrual

	actualSum, err = sut.GetAccrualSumByUserId(user.ID)
	require.NoError(t, err)
	require.Equal(t, expectedSum, actualSum)
}

func assertOrder(t *testing.T, expected entity.Order, actual entity.Order, createAtMin time.Time, createAtMax time.Time) {
	require.WithinRange(t, actual.CreatedAt, createAtMin.Truncate(time.Second), createAtMax)
	now := time.Time{}
	expected.CreatedAt = now
	actual.CreatedAt = now
	require.Equal(t, expected, actual)
}
