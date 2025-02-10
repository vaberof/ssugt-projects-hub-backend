package main

import (
	"context"
	"github.com/shopspring/decimal"
	"log"
	"os"
	"os/signal"
	"ssugt-projects-hub/config"
	"ssugt-projects-hub/pkg/logging/logs"
	"syscall"
)

func main() {
	configureDecimal()
	config.Init()

	mainCtx, cancelMainCtx := context.WithCancel(context.Background())
	defer cancelMainCtx()

	logger := logs.New(os.Stdout, nil, "ssugt-projects-hub")
	log.SetFlags(log.Flags() | log.Llongfile)

	app := NewApp(mainCtx, logger)
	app.initDatabases()
	app.initServices()
	app.initServer()

	app.Start()
	stop := getStopSignal()
	<-stop

	app.Stop()
}

func configureDecimal() {
	decimal.DivisionPrecision = 8
	decimal.MarshalJSONWithoutQuotes = true
}

func getStopSignal() <-chan os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)

	return stop
}
