package service

import (
	"github.com/smamykin/gofermart/internal/entity"
)

type OrderService struct {
	OrderRepository OrderRepositoryInterface
}

func (o *OrderService) AddOrder(userID int, orderNumber string) (order entity.Order, err error) {
	order = entity.Order{
		UserID:        userID,
		OrderNumber:   orderNumber,
		Status:        entity.OrderStatusNew,
		AccrualStatus: entity.AccrualStatusUndefined,
		Accrual:       0,
	}
	_, err = o.OrderRepository.GetOrderByOrderNumber(orderNumber)
	if err == nil {
		return entity.Order{}, ErrOrderAlreadyExists
	}

	if err != ErrEntityIsNotFound {
		return entity.Order{}, err
	}

	order, err = o.OrderRepository.AddOrder(order)

	if err != nil {
		return entity.Order{}, err
	}

	return order, nil
}
