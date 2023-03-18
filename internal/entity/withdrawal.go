package entity

import "time"

type Withdrawal struct {
	ID          int
	UserID      int
	OrderNumber string
	Amount      float64
	CreatedAt   time.Time
}
