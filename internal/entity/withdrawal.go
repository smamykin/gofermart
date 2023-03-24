package entity

import (
	"github.com/smamykin/gofermart/pkg/money"
	"time"
)

type Withdrawal struct {
	ID          int
	UserID      int
	OrderNumber string
	Amount      money.IntMoney
	CreatedAt   time.Time
}
