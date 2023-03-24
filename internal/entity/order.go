package entity

import (
	"github.com/smamykin/gofermart/pkg/money"
	"time"
)

type Status int

func (s Status) String() string {
	switch s {
	case OrderStatusNew:
		return "NEW"
	case OrderStatusProcessing:
		return "PROCESSING"
	case OrderStatusInvalid:
		return "INVALID"
	case OrderStatusProcessed:
		return "PROCESSED"
	default:
		panic("unknown status")
	}
}

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
	AccrualStatusUnregistered
)

func (s AccrualStatus) String() string {
	switch s {
	case AccrualStatusUndefined:
		return "UNDEFINED"
	case AccrualStatusRegistered:
		return "REGISTERED"
	case AccrualStatusInvalid:
		return "INVALID"
	case AccrualStatusProcessing:
		return "PROCESSING"
	case AccrualStatusProcessed:
		return "PROCESSED"
	case AccrualStatusUnregistered:
		return "UNREGISTERED"
	default:
		panic("unknown status")
	}
}

type Order struct {
	ID            int
	UserID        int
	OrderNumber   string
	Status        Status
	AccrualStatus AccrualStatus
	Accrual       money.IntMoney
	CreatedAt     time.Time
}
