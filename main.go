package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/DanielTrondoli/go-kit_microservice_example/account"
	"github.com/DanielTrondoli/go-kit_microservice_example/repository"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func main() {

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "account",
			"time", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	var db = repository.NewMemoryRepo([]string{"id", "email", "password"}, logger)

	var httpAddr = flag.String("http", ":9000", "http listen address")

	flag.Parse()
	ctx := context.Background()

	var srv account.Service
	{
		repository := db
		srv = account.NewService(repository, logger)
	}

	errsChan := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errsChan <- fmt.Errorf("%s", <-c)
	}()

	endpoints := account.MakeEndPoints(srv)

	go func() {
		fmt.Println("listening on port", *httpAddr)
		handler := account.NewHTTPServer(ctx, endpoints)
		errsChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	level.Error(logger).Log("exit", <-errsChan)
}
