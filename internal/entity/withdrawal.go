package entity

import "time"

type Withdrawal struct {
	ID        int
	UserID    int
	Amount    float64
	CreatedAt time.Time
}
