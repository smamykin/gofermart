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

//func checksum(number int) int {
//	var luhn int
//
//	for i := 0; number > 0; i++ {
//		cur := number % 10
//
//		if i%2 == 0 { // even
//			cur = cur * 2
//			if cur > 9 {
//				cur = cur%10 + cur/10
//			}
//		}
//
//		luhn += cur
//		number = number / 10
//	}
//	return luhn % 10
//}
