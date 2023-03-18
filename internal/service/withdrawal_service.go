package service

import "github.com/smamykin/gofermart/internal/entity"

type WithdrawalService struct {
	WithdrawalRepository WithdrawalRepositoryInterface
	OrderRepository      OrderRepositoryInterface
}

func (w *WithdrawalService) Withdraw(userID int, amount float64, orderNumber string) (withdrawal entity.Withdrawal, err error) {
	err = orderNumberValidation(orderNumber)
	if err != nil {
		return entity.Withdrawal{}, err
	}

	withdrawal, err = w.WithdrawalRepository.GetWithdrawalByOrderNumber(orderNumber)
	if err == nil {
		return withdrawal, ErrEntityAlreadyExists
	}

	if err != ErrEntityIsNotFound {
		return entity.Withdrawal{}, err
	}

	withdrawalSum, err := w.WithdrawalRepository.GetAmountSumByUserID(userID)
	if err != nil {
		return entity.Withdrawal{}, err
	}
	accrualSum, err := w.OrderRepository.GetAccrualSumByUserID(userID)
	if err != nil {
		return entity.Withdrawal{}, err
	}
	if (accrualSum - withdrawalSum - amount) < 0 {
		return entity.Withdrawal{}, ErrNotEnoughAccrual
	}

	withdrawal = entity.Withdrawal{
		UserID:      userID,
		OrderNumber: orderNumber,
		Amount:      amount,
	}

	withdrawal, err = w.WithdrawalRepository.AddWithdrawal(withdrawal)

	if err != nil {
		return entity.Withdrawal{}, err
	}

	return withdrawal, err
}
