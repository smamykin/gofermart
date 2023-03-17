package service

import (
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/smamykin/gofermart/internal/entity"
)

type OrderService struct {
	OrderRepository OrderRepositoryInterface
}

func (o *OrderService) AddOrder(userID int, orderNumber string) (order entity.Order, err error) {
	err = orderNumberValidation(orderNumber)
	if err != nil {
		return entity.Order{}, err
	}

	order, err = o.OrderRepository.GetOrderByOrderNumber(orderNumber)
	if err == nil {
		return order, ErrOrderAlreadyExists
	}

	order = entity.Order{
		UserID:        userID,
		OrderNumber:   orderNumber,
		Status:        entity.OrderStatusNew,
		AccrualStatus: entity.AccrualStatusUndefined,
		Accrual:       0,
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

func (o *OrderService) GetAllOrdersByUserID(userID int) (orders []entity.Order, err error) {
	return o.OrderRepository.GetAllByUserID(userID)
}

func orderNumberValidation(numberAsString string) error {
	err := goluhn.Validate(numberAsString)
	if err != nil {
		return ErrInvalidOrderNumber
	}
	return nil
}
