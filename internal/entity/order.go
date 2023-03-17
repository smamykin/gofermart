package entity

import "time"

type Status int

const (
	OrderStatusNew Status = iota
	OrderStatusProcessing
	OrderStatusInvalid
	OrderStatusProcessed
)

type AccrualStatus int

const (
	AccrualStatusUndefined AccrualStatus = iota
	AccrualStatusRegistered
	AccrualStatusInvalid
	AccrualStatusProcessing
	AccrualStatusProcessed
)

type Order struct {
	ID            int
	UserID        int
	OrderNumber   string
	Status        Status
	AccrualStatus AccrualStatus
	Accrual       int
	CreatedAt     time.Time
}
