package main

import (
	"github.com/GeneralKenobi/mailman/internal/api"
	"github.com/GeneralKenobi/mailman/internal/config"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
	"github.com/GeneralKenobi/mailman/pkg/shutdown"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	argsCfg := commandLineArgsConfig()

	cfg, err := config.Load(argsCfg.configFiles)
	if err != nil {
		mdctx.Fatalf(nil, "Error loading configuration: %v", err)
	}

	// TODO: Shutdown time configuration
	parentCtx := shutdown.NewParentContext(30 * time.Second)
	go api.ServeHttp(parentCtx.NewContext("http server"), cfg.HttpServer)

	stopSignalChannel := make(chan os.Signal)
	// SIGINT for ctrl+c, SIGTERM for k8s stopping the container.
	signal.Notify(stopSignalChannel, syscall.SIGINT, syscall.SIGTERM)
	caughtSignal := <-stopSignalChannel
	mdctx.Infof(nil, "Caught signal %v, shutting down", caughtSignal)
	parentCtx.Cancel()
	mdctx.Infof(nil, "Shutdown completed, exiting")
}
