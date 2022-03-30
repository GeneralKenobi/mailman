package main

import (
	"github.com/GeneralKenobi/mailman/internal/api/httpgin"
	"github.com/GeneralKenobi/mailman/internal/config"
	"github.com/GeneralKenobi/mailman/internal/db/postgres"
	"github.com/GeneralKenobi/mailman/internal/email/mock"
	"github.com/GeneralKenobi/mailman/internal/job/mailingentry"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
	"github.com/GeneralKenobi/mailman/pkg/shutdown"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	configure()
	parentCtx := shutdown.NewParentContext(time.Duration(config.Get().Global.ShutdownTimeoutSeconds) * time.Second)
	bootstrap(parentCtx)
	shutdownAfterStopSignal(parentCtx)
}

func configure() {
	argsCfg := commandLineArgsConfig()

	err := config.Load(argsCfg.configFiles)
	if err != nil {
		mdctx.Fatalf(nil, "Error loading configuration: %v", err)
	}

	err = mdctx.SetLogLevelFromString(argsCfg.logLevel)
	if err != nil {
		mdctx.Fatalf(nil, "Error setting log level: %v", err)
	}
}

func bootstrap(parentCtx shutdown.ParentContext) {
	// DB
	dbCtx, err := postgres.NewContext(parentCtx.NewContext("postgres"))
	if err != nil {
		mdctx.Fatalf(nil, "Error connecting to DB: %v", err)
	}

	// Email service
	emailer := mock.NewEmailer()

	// Scheduled jobs
	mailingEntryCleanupJob := mailingentry.NewCleanupJob(dbCtx)
	go mailingEntryCleanupJob.RunScheduled(parentCtx.NewContext("scheduled stale mailing entry cleanup"))

	// HTTP server
	httpServer := httpgin.NewServer(dbCtx, emailer)
	go httpServer.Run(parentCtx.NewContext("http server"))
}

func shutdownAfterStopSignal(parentCtx shutdown.ParentContext) {
	stopSignalChannel := make(chan os.Signal)
	// SIGINT for ctrl+c, SIGTERM for k8s stopping the container.
	signal.Notify(stopSignalChannel, syscall.SIGINT, syscall.SIGTERM)

	caughtSignal := <-stopSignalChannel
	mdctx.Infof(nil, "Caught signal %v, shutting down", caughtSignal)

	parentCtx.Cancel()
	mdctx.Infof(nil, "Shutdown completed, exiting")
}
