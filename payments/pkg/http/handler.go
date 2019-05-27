package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/deterok/go_test_task/payments/pkg/endpoint"
)

// NewHTTPHandler returns a handler that makes a set of endpoints available on
// predefined paths.
func NewHTTPHandler(endpoints endpoint.Endpoints, options map[string][]kithttp.ServerOption) http.Handler {
	m := mux.NewRouter()
	makeCreateAccountHandler(m, endpoints, options["CreateAccount"])
	makeGetAccountHandler(m, endpoints, options["GetAccount"])
	makeGetAccountsHandler(m, endpoints, options["GetAccounts"])
	makeGetAccountOperationsHandler(m, endpoints, options["GetAccountOperations"])
	makeMakeDepositHandler(m, endpoints, options["MakeDeposit"])
	makeMakeTransferHandler(m, endpoints, options["MakeTransfer"])
	return m
}

func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}
func ErrorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

func err2code(err error) int {
	return http.StatusInternalServerError
}

type errorWrapper struct {
	Error string `json:"error"`
}
