package endpoint

import (
	"context"

	"github.com/deterok/go_test_task/payments/pkg/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/shopspring/decimal"
)

// MakeDepositRequest collects the request parameters for the MakeDeposit method.
type MakeDepositRequest struct {
	To       int64           `json:"to"`
	Currency string          `json:"currency"`
	Amount   decimal.Decimal `json:"amount"`
}

// MakeDepositResponse collects the response parameters for the MakeDeposit method.
type MakeDepositResponse struct {
	Operation *service.Operation `json:"operation"`
	Err       error              `json:"error,omitempty"`
}

// MakeMakeDepositEndpoint returns an endpoint that invokes MakeDeposit on the service.
func MakeMakeDepositEndpoint(s service.PaymentsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(MakeDepositRequest)
		o, err := s.MakeDeposit(ctx, req.To, req.Currency, req.Amount)
		return MakeDepositResponse{
			Operation: o,
			Err:       err,
		}, nil
	}
}

// Failed implements Failer.
func (r MakeDepositResponse) Failed() error {
	return r.Err
}

// MakeTransferRequest collects the request parameters for the MakeTransfer method.
type MakeTransferRequest struct {
	From     int64           `json:"from"`
	To       int64           `json:"to"`
	Currency string          `json:"currency"`
	Amount   decimal.Decimal `json:"amount"`
}

// MakeTransferResponse collects the response parameters for the MakeTransfer method.
type MakeTransferResponse struct {
	Operation *service.Operation `json:"operation"`
	Err       error              `json:"error,omitempty"`
}

// MakeMakeTransferEndpoint returns an endpoint that invokes MakeTransfer on the service.
func MakeMakeTransferEndpoint(s service.PaymentsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(MakeTransferRequest)
		o, err := s.MakeTransfer(ctx, req.From, req.To, req.Currency, req.Amount)
		return MakeTransferResponse{
			Operation: o,
			Err:       err,
		}, nil
	}
}

// Failed implements Failer.
func (r MakeTransferResponse) Failed() error {
	return r.Err
}

// MakeDeposit implements Service.
func (e Endpoints) MakeDeposit(ctx context.Context, to int64, currency string, amount decimal.Decimal) (*service.Operation, error) {
	request := MakeDepositRequest{
		Amount:   amount,
		Currency: currency,
		To:       to,
	}
	response, err := e.MakeDepositEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(MakeDepositResponse).Operation, response.(MakeDepositResponse).Err
}

// MakeTransfer implements Service.
func (e Endpoints) MakeTransfer(ctx context.Context, from int64, to int64, currency string, amount decimal.Decimal) (*service.Operation, error) {
	request := MakeTransferRequest{
		Amount:   amount,
		Currency: currency,
		From:     from,
		To:       to,
	}
	response, err := e.MakeTransferEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(MakeTransferResponse).Operation, response.(MakeTransferResponse).Err
}
