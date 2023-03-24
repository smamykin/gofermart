package service

import (
	"errors"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/smamykin/gofermart/internal/entity"
	"github.com/smamykin/gofermart/pkg/contracts"
	"github.com/smamykin/gofermart/pkg/money"
)

type OrderService struct {
	OrderRepository OrderRepositoryInterface
	AccrualClient   AccrualClientInterface
	Logger          contracts.LoggerInterface
}

func (o *OrderService) AddOrder(userID int, orderNumber string) (order entity.Order, err error) {
	err = orderNumberValidation(orderNumber)
	if err != nil {
		return entity.Order{}, err
	}

	order, err = o.OrderRepository.GetOrderByOrderNumber(orderNumber)
	if err == nil {
		return order, ErrEntityAlreadyExists
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
func (o *OrderService) UpdateOrdersStatuses() error {
	orders, err := o.OrderRepository.GetOrdersWithUnfinishedStatus()
	if err != nil {
		return err
	}
	for _, order := range orders {
		accrualOrder, err := o.AccrualClient.GetOrder(order.OrderNumber)
		if err != nil && !errors.Is(err, ErrEntityIsNotFound) {
			o.Logger.Err(err, "error occurred while getting order from accrual service")
			continue
		}

		if errors.Is(err, ErrEntityIsNotFound) {
			o.Logger.Warn(err, "the order is not found in the accrual service")
		}

		switch accrualOrder.Status {
		case entity.AccrualStatusRegistered, entity.AccrualStatusProcessing:
			order.Status = entity.OrderStatusProcessing
		case entity.AccrualStatusInvalid:
			order.Status = entity.OrderStatusInvalid
		case entity.AccrualStatusProcessed:
			order.Status = entity.OrderStatusProcessed
		case entity.AccrualStatusUnregistered:
			order.Status = entity.OrderStatusInvalid
		default:
			panic("unknown status")
		}

		order.AccrualStatus = accrualOrder.Status
		order.Accrual = money.FromFloat(accrualOrder.Accrual)

		_, err = o.OrderRepository.UpdateOrder(order)
		if err != nil {
			o.Logger.Err(err, "error occurred while updating order statuses")
		}
	}

	return nil
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
