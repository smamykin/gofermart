package repository

import (
	"github.com/smamykin/gofermart/internal/entity"
	"github.com/smamykin/gofermart/internal/service"
	"github.com/smamykin/gofermart/tests/Functional/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestWithdrawalRepository_AddWithdrawal(t *testing.T) {
	c := utils.GetContainer(t)
	db := c.DB()
	utils.TruncateTable(t, db)

	user := utils.InsertUser(t, db, entity.User{})
	sut := c.WithdrawalRepository()
	expectedWithdrawal := entity.Withdrawal{
		UserID:      user.ID,
		OrderNumber: "123",
		Amount:      33,
	}
	timeMin := time.Now()
	actualWithdrawal, err := sut.AddWithdrawal(expectedWithdrawal)
	timeMax := time.Now()
	require.NoError(t, err)

	expectedWithdrawal.ID = 1 // because the table was truncated beforehand, we can guess the id

	withdrawalFromDB, err := sut.GetWithdrawal(expectedWithdrawal.ID)
	require.NoError(t, err)

	assertWithdrawal(t, expectedWithdrawal, withdrawalFromDB, timeMin, timeMax)
	assertWithdrawal(t, expectedWithdrawal, actualWithdrawal, timeMin, timeMax)

	require.NoError(t, err)
}

func TestWithdrawalRepository_GetAmountSumByUserId(t *testing.T) {
	c := utils.GetContainer(t)
	db := c.DB()
	utils.TruncateTable(t, db)

	user := utils.InsertUser(t, db, entity.User{})
	anotherUser := utils.InsertUser(t, db, entity.User{})

	sut := c.WithdrawalRepository()

	//first check if there is no withdrawals at all
	actualSum, err := sut.GetAmountSumByUserID(user.ID)
	require.NoError(t, err)
	require.Equal(t, .0, actualSum)

	//second check if there are withdrawals
	withdrawal0, err := sut.AddWithdrawal(entity.Withdrawal{
		UserID:      user.ID,
		OrderNumber: "111",
		Amount:      11.1,
	})
	require.NoError(t, err)
	withdrawal1, err := sut.AddWithdrawal(entity.Withdrawal{
		UserID:      user.ID,
		OrderNumber: "222",
		Amount:      22.2,
	})
	require.NoError(t, err)
	_, err = sut.AddWithdrawal(entity.Withdrawal{
		UserID:      anotherUser.ID,
		OrderNumber: "333",
		Amount:      33.3,
	})
	require.NoError(t, err)

	expectedSum := withdrawal0.Amount + withdrawal1.Amount

	actualSum, err = sut.GetAmountSumByUserID(user.ID)
	require.NoError(t, err)
	require.Equal(t, expectedSum, actualSum)
}
func TestWithdrawalRepository_GetWithdrawalByOrderNumber(t *testing.T) {
	c := utils.GetContainer(t)
	db := c.DB()
	utils.TruncateTable(t, db)

	user := utils.InsertUser(t, db, entity.User{})
	sut := c.WithdrawalRepository()
	expectedWithdrawal, err := sut.AddWithdrawal(entity.Withdrawal{
		UserID:      user.ID,
		OrderNumber: "1",
		Amount:      0,
	})
	require.NoError(t, err)

	actualOrder, err := sut.GetWithdrawalByOrderNumber(expectedWithdrawal.OrderNumber)
	require.NoError(t, err)

	require.Equal(t, expectedWithdrawal, actualOrder)

	_, err = sut.GetWithdrawalByOrderNumber("unknown order number")
	require.Equal(t, err, service.ErrEntityIsNotFound)
}

func assertWithdrawal(t *testing.T, expected entity.Withdrawal, actual entity.Withdrawal, createAtMin time.Time, createAtMax time.Time) {
	require.WithinRange(t, actual.CreatedAt, createAtMin.Truncate(time.Second), createAtMax)
	now := time.Time{}
	expected.CreatedAt = now
	actual.CreatedAt = now
	require.Equal(t, expected, actual)
}
