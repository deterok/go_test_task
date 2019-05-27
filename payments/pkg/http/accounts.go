package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/deterok/go_test_task/payments/pkg/endpoint"
)

// ─── GET ACCOUNT ─────────────────────────────────────────────────────────────────

func makeGetAccountHandler(m *mux.Router, endpoints endpoint.Endpoints, options []kithttp.ServerOption) {
	handler := handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))
	server := kithttp.NewServer(endpoints.GetAccountEndpoint, decodeGetAccountRequest, encodeGetAccountResponse, options...)
	m.Methods("GET").Path("/accounts/{id}").Handler(handler(server))
}

func decodeGetAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	strID := vars["id"]
	id, err := strconv.ParseInt(strID, 10, 64)

	req := endpoint.GetAccountRequest{
		ID: id,
	}

	if err != nil {
		return req, errors.Wrap(err, "order reciept failed")
	}

	return req, err
}

func encodeGetAccountResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// ─── GET ACCOUNTS ───────────────────────────────────────────────────────────────

func makeGetAccountsHandler(m *mux.Router, endpoints endpoint.Endpoints, options []kithttp.ServerOption) {
	handler := kithttp.NewServer(endpoints.GetAccountsEndpoint, decodeGetAccountsRequest, encodeGetAccountsResponse, options...)
	m.Methods("GET").Path("/accounts").Handler(handler)
}

func decodeGetAccountsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpoint.GetAccountsRequest{}, nil
}

func encodeGetAccountsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// ─── GET ACCOUNT OPERATIONS ───────────────────────────────────────────────────────

func makeGetAccountOperationsHandler(m *mux.Router, endpoints endpoint.Endpoints, options []kithttp.ServerOption) {
	handler := kithttp.NewServer(endpoints.GetAccountOperationsEndpoint, decodeGetAccountOperationsRequest, encodeGetAccountOperationsResponse, options...)
	m.Methods("GET").Path("/accounts/{id}/operations").Handler(handler)
}

func decodeGetAccountOperationsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	strID := vars["id"]
	id, err := strconv.ParseInt(strID, 10, 64)
	req := endpoint.GetAccountOperationsRequest{
		AccountID: id,
	}

	if err != nil {
		return req, errors.Wrap(err, "order getting failed")
	}

	return req, err
}

func encodeGetAccountOperationsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// ─── CREATE ACCOUNT ──────────────────────────────────────────────────────────────

func makeCreateAccountHandler(m *mux.Router, endpoints endpoint.Endpoints, options []kithttp.ServerOption) {
	handler := kithttp.NewServer(endpoints.CreateAccountEndpoint, decodeCreateAccountRequest, encodeCreateAccountResponse, options...)
	m.Methods("POST").Path("/accounts").Handler(handler)
}

func decodeCreateAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := endpoint.CreateAccountRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeCreateAccountResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}
