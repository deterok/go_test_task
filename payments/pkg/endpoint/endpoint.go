package endpoint

import (
	"github.com/go-kit/kit/endpoint"

	"github.com/deterok/go_test_task/payments/pkg/service"
)

// Endpoints collects all of the endpoints that compose a profile service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	CreateAccountEndpoint        endpoint.Endpoint
	GetAccountEndpoint           endpoint.Endpoint
	GetAccountsEndpoint          endpoint.Endpoint
	GetAccountOperationsEndpoint endpoint.Endpoint
	MakeDepositEndpoint         endpoint.Endpoint
	MakeTransferEndpoint         endpoint.Endpoint
}

// New returns a Endpoints struct that wraps the provided service, and wires in all of the
// expected endpoint middlewares
func New(s service.PaymentsService, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		CreateAccountEndpoint:        MakeCreateAccountEndpoint(s),
		GetAccountEndpoint:           MakeGetAccountEndpoint(s),
		GetAccountOperationsEndpoint: MakeGetAccountOperationsEndpoint(s),
		GetAccountsEndpoint:          MakeGetAccountsEndpoint(s),
		MakeDepositEndpoint:         MakeMakeDepositEndpoint(s),
		MakeTransferEndpoint:         MakeMakeTransferEndpoint(s),
	}
	for _, m := range mdw["CreateAccount"] {
		eps.CreateAccountEndpoint = m(eps.CreateAccountEndpoint)
	}
	for _, m := range mdw["GetAccount"] {
		eps.GetAccountEndpoint = m(eps.GetAccountEndpoint)
	}
	for _, m := range mdw["GetAccounts"] {
		eps.GetAccountsEndpoint = m(eps.GetAccountsEndpoint)
	}
	for _, m := range mdw["GetAccountOperations"] {
		eps.GetAccountOperationsEndpoint = m(eps.GetAccountOperationsEndpoint)
	}
	for _, m := range mdw["MakeDepositEndpoint"] {
		eps.MakeDepositEndpoint = m(eps.MakeDepositEndpoint)
	}
	for _, m := range mdw["MakeTransfer"] {
		eps.MakeTransferEndpoint = m(eps.MakeTransferEndpoint)
	}
	return eps
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}
