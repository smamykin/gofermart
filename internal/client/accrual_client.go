package client

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/smamykin/gofermart/internal/entity"
	"github.com/smamykin/gofermart/internal/service"
	"net/http"
)

func NewAccrualClient(apiEntrypoint string) *AccrualClient {
	return &AccrualClient{
		client:        resty.New(),
		apiEntrypoint: apiEntrypoint,
	}
}

type AccrualClient struct {
	client        *resty.Client
	apiEntrypoint string
}

func (a *AccrualClient) GetOrder(orderNumber string) (service.AccrualOrder, error) {
	order := AccrualOrderResponseModel{}
	resp, err := a.client.R().
		SetHeader("Accept", "application/json").
		SetPathParam("number", orderNumber).
		SetResult(order).
		Get(a.apiEntrypoint + "/api/orders/{number}")
	if err != nil {
		return service.AccrualOrder{}, err
	}

	if resp.StatusCode() == http.StatusNoContent {
		return service.AccrualOrder{
			Order:  orderNumber,
			Status: entity.AccrualStatusUnregistered,
		}, nil
	}

	if resp.StatusCode() != http.StatusOK {
		return service.AccrualOrder{}, fmt.Errorf(
			"unprocessable result from accrual service. Status: %s, Body: %s",
			resp.Status(),
			resp.String(),
		)
	}
	model, ok := resp.Result().(AccrualOrderResponseModel)
	if !ok {
		return service.AccrualOrder{}, fmt.Errorf(
			"unprocessable result from accrual service. Status: %s, Body: %s",
			resp.Status(),
			resp.String(),
		)
	}

	accrualStatus, err := getAccrualStatus(model.Status)
	if err != nil {
		return service.AccrualOrder{}, err
	}

	return service.AccrualOrder{
		Order:   orderNumber,
		Accrual: model.Accrual,
		Status:  accrualStatus,
	}, nil

}

func getAccrualStatus(status string) (entity.AccrualStatus, error) {
	switch status {
	case "REGISTERED":
		return entity.AccrualStatusRegistered, nil
	case "INVALID":
		return entity.AccrualStatusInvalid, nil
	case "PROCESSING":
		return entity.AccrualStatusProcessing, nil
	case "PROCESSED":
		return entity.AccrualStatusProcessed, nil
	default:
		return 0, fmt.Errorf("unknown status from accrual service: %s", status)
	}
}

type AccrualOrderResponseModel struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual,omitempty"`
}
