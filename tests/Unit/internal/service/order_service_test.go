package service

import (
	"github.com/smamykin/gofermart/internal/entity"
	"github.com/smamykin/gofermart/internal/service"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOrderService_AddOrder(t *testing.T) {
	sut := service.OrderService{}
	userId := 1
	orderNumber := 1
	expectedOrder := entity.Order{}

	actualOrder, err := sut.AddOrder(userId, orderNumber)

	require.NoError(t, err)
	require.Equal(t, expectedOrder, actualOrder)
}
