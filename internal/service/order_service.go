package service

import "github.com/smamykin/gofermart/internal/entity"

type OrderService struct {
	orderRepository OrderRepositoryInterface
}

func (s *OrderService) AddOrder(userID int, orderNumber int) (order entity.Order, err error) {
	return
}
