package service

import (
	"flag"
	"fmt"
	"net"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/deterok/go_test_task/payments/pkg/endpoint"
	http "github.com/deterok/go_test_task/payments/pkg/http"
	service "github.com/deterok/go_test_task/payments/pkg/service"
	kitendpoint "github.com/go-kit/kit/endpoint"
	log "github.com/go-kit/kit/log"
	"github.com/oklog/oklog/pkg/group"
	opentracinggo "github.com/opentracing/opentracing-go"
)

var tracer opentracinggo.Tracer
var logger log.Logger

// Some flags
var (
	fs       = flag.NewFlagSet("payments", flag.ExitOnError)
	httpAddr = fs.String("http-addr", ":8081", "HTTP listen address")
	redisAddr = fs.String("redis-addr", "redis:6379", "Redis address")
	// Database
	dbDialect = fs.String("dialect", "postgres", "Database dialect")
	dbDSN = fs.String("db-dsn", "host=postgres sslmode=disable user=postgres", "Database DSN")
)

func Run() {
	fs.Parse(os.Args[1:])

	// Create a single logger, which we'll use and give to other components.
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	tracer = opentracinggo.GlobalTracer()

	db, err := gorm.Open(*dbDialect, *dbDSN)
	if err != nil {
		panic(err)
	}

	if err := service.InitModels(db); err != nil {
		panic(err)
	}

	lock := service.NewLockReposytory(*redisAddr)
	uowFacotry := service.NewUOWPaymentsFactory(db)
	svc := service.New(lock, uowFacotry, getServiceMiddleware(logger))
	eps := endpoint.New(svc, getEndpointMiddleware(logger))
	g := createService(eps)
	initCancelInterrupt(g)
	logger.Log("exit", g.Run())
}


func initHttpHandler(endpoints endpoint.Endpoints, g *group.Group) {
	options := defaultHttpOptions(logger, tracer)

	httpHandler := http.NewHTTPHandler(endpoints, options)
	httpListener, err := net.Listen("tcp", *httpAddr)
	if err != nil {
		logger.Log("transport", "HTTP", "during", "Listen", "err", err)
	}
	g.Add(func() error {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		return nethttp.Serve(httpListener, httpHandler)
	}, func(error) {
		httpListener.Close()
	})

}
func getServiceMiddleware(logger log.Logger) (mw []service.Middleware) {
	mw = []service.Middleware{}
	return
}
func getEndpointMiddleware(logger log.Logger) (mw map[string][]kitendpoint.Middleware) {
	mw = map[string][]kitendpoint.Middleware{}
	return
}


func initCancelInterrupt(g *group.Group) {
	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})
}
