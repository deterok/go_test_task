package service

import (
	payendpoint "github.com/deterok/go_test_task/payments/pkg/endpoint"
	payhttp "github.com/deterok/go_test_task/payments/pkg/http"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/oklog/oklog/pkg/group"
	opentracinggo "github.com/opentracing/opentracing-go"
)

func createService(endpoints payendpoint.Endpoints) (g *group.Group) {
	g = &group.Group{}
	initHttpHandler(endpoints, g)
	return g
}
func defaultHttpOptions(logger log.Logger, tracer opentracinggo.Tracer) map[string][]kithttp.ServerOption {
	options := map[string][]kithttp.ServerOption{
		"CreateAccount": {
			kithttp.ServerErrorEncoder(payhttp.ErrorEncoder),
			kithttp.ServerErrorLogger(logger),
		},
		"GetAccount": {
			kithttp.ServerErrorEncoder(payhttp.ErrorEncoder),
			kithttp.ServerErrorLogger(logger),
		},
		"GetAccountOperations": {
			kithttp.ServerErrorEncoder(payhttp.ErrorEncoder),
			kithttp.ServerErrorLogger(logger),
		},
		"GetAccounts": {
			kithttp.ServerErrorEncoder(payhttp.ErrorEncoder),
			kithttp.ServerErrorLogger(logger),
		},
		"MakeDeposit": {
			kithttp.ServerErrorEncoder(payhttp.ErrorEncoder),
			kithttp.ServerErrorLogger(logger),
		},
		"MakeTransfer": {
			kithttp.ServerErrorEncoder(payhttp.ErrorEncoder),
			kithttp.ServerErrorLogger(logger),
		},
	}
	return options
}

func addEndpointMiddlewareToAllMethods(mw map[string][]endpoint.Middleware, m endpoint.Middleware) {
	methods := []string{"CreateAccount", "GetAccount", "GetAccounts", "GetAccountOperations", "MakeDeposit", "MakeTransfer"}
	for _, v := range methods {
		mw[v] = append(mw[v], m)
	}
}
