package http

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/deterok/go_test_task/payments/pkg/endpoint"
)

// ─── MAKE Deposit ───────────────────────────────────────────────────────────────

// makeMakeDepositHandler creates the handler logic
func makeMakeDepositHandler(m *mux.Router, endpoints endpoint.Endpoints, options []kithttp.ServerOption) {
	handler := kithttp.NewServer(endpoints.MakeDepositEndpoint, decodeMakeDepositRequest, encodeMakeDepositResponse, options...)
	m.Methods("POST").Path("/operations/Deposit").Handler(handler)
}

// decodeMakeDepositRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeMakeDepositRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := endpoint.MakeDepositRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// encodeMakeDepositResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeMakeDepositResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}

// ─── MAKE TRANSFER ──────────────────────────────────────────────────────────────

func makeMakeTransferHandler(m *mux.Router, endpoints endpoint.Endpoints, options []kithttp.ServerOption) {
	handler := kithttp.NewServer(endpoints.MakeTransferEndpoint, decodeMakeTransferRequest, encodeMakeTransferResponse, options...)
	m.Methods("POST").Path("/operations/transfer").Handler(handler)
}

func decodeMakeTransferRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := endpoint.MakeTransferRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeMakeTransferResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return
}
