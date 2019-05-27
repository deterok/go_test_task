package endpoint

import (
	"context"

	"github.com/deterok/go_test_task/payments/pkg/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/shopspring/decimal"
)

// ─── ENDPOINTS ENVOKERS ─────────────────────────────────────────────────────────

// CreateAccountRequest collects the request parameters for the CreateAccount method.
type CreateAccountRequest struct {
	Name     string `json:"name"`
	Currency string `json:"currency"`
}

// CreateAccountResponse collects the response parameters for the CreateAccount method.
type CreateAccountResponse struct {
	Account *service.Account `json:"account"`
	Err     error            `json:"error,omitempty"`
}

// MakeCreateAccountEndpoint returns an endpoint that invokes CreateAccount on the service.
func MakeCreateAccountEndpoint(s service.PaymentsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateAccountRequest)
		a, err := s.CreateAccount(ctx, req.Name, req.Currency)
		return CreateAccountResponse{
			Account: a,
			Err:     err,
		}, nil
	}
}

// Failed implements Failer.
func (r CreateAccountResponse) Failed() error {
	return r.Err
}

// GetAccountRequest collects the request parameters for the GetAccount method.
type GetAccountRequest struct {
	ID int64 `json:"id"`
}

// GetAccountResponse collects the response parameters for the GetAccount method.
type GetAccountResponse struct {
	Account *service.Account `json:"account"`
	Err     error            `json:"error,omitempty"`
}

// MakeGetAccountEndpoint returns an endpoint that invokes GetAccount on the service.
func MakeGetAccountEndpoint(s service.PaymentsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAccountRequest)
		a, err := s.GetAccount(ctx, req.ID)
		return GetAccountResponse{
			Account: a,
			Err:     err,
		}, nil
	}
}

// Failed implements Failer.
func (r GetAccountResponse) Failed() error {
	return r.Err
}

// GetAccountsRequest collects the request parameters for the GetAccounts method.
type GetAccountsRequest struct{}

// GetAccountsResponse collects the response parameters for the GetAccounts method.
type GetAccountsResponse struct {
	Account []*service.Account `json:"account"`
	Err     error              `json:"error,omitempty"`
}

// MakeGetAccountsEndpoint returns an endpoint that invokes GetAccounts on the service.
func MakeGetAccountsEndpoint(s service.PaymentsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		a, err := s.GetAccounts(ctx)
		return GetAccountsResponse{
			Account: a,
			Err:     err,
		}, nil
	}
}

// Failed implements Failer.
func (r GetAccountsResponse) Failed() error {
	return r.Err
}

// GetAccountOperationsRequest collects the request parameters for the GetAccountOperations method.
type GetAccountOperationsRequest struct {
	AccountID int64 `json:"account_id"`
}

// GetAccountOperationsResponse collects the response parameters for the GetAccountOperations method.
type GetAccountOperationsResponse struct {
	Operations []*service.Operation `json:"operations"`
	Err        error                `json:"error,omitempty"`
}

// MakeGetAccountOperationsEndpoint returns an endpoint that invokes GetAccountOperations on the service.
func MakeGetAccountOperationsEndpoint(s service.PaymentsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAccountOperationsRequest)
		o, err := s.GetAccountOperations(ctx, req.AccountID)
		return GetAccountOperationsResponse{
			Operations: o,
			Err:        err,
		}, nil
	}
}

// Failed implements Failer.
func (r GetAccountOperationsResponse) Failed() error {
	return r.Err
}

// ─── ENDPOINTS IMPLIMENTATION ───────────────────────────────────────────────────

// GetAccount implements Service.
func (e Endpoints) GetAccount(ctx context.Context, id int64) (*service.Account, error) {
	request := GetAccountRequest{ID: id}
	response, err := e.GetAccountEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(GetAccountResponse).Account, response.(GetAccountResponse).Err
}

// GetAccounts implements Service.
func (e Endpoints) GetAccounts(ctx context.Context) ([]*service.Account, error) {
	request := GetAccountsRequest{}
	response, err := e.GetAccountsEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(GetAccountsResponse).Account, response.(GetAccountsResponse).Err
}

// GetAccountOperations implements Service.
func (e Endpoints) GetAccountOperations(ctx context.Context, accID int64) ([]*service.Operation, error) {
	request := GetAccountOperationsRequest{AccountID: accID}
	response, err := e.GetAccountOperationsEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(GetAccountOperationsResponse).Operations, response.(GetAccountOperationsResponse).Err
}

// CreateAccount implements Service.
func (e Endpoints) CreateAccount(ctx context.Context, name string, currency string, amount decimal.Decimal) (*service.Account, error) {
	request := CreateAccountRequest{
		Currency: currency,
		Name:     name,
	}
	response, err := e.CreateAccountEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.(CreateAccountResponse).Account, response.(CreateAccountResponse).Err
}
