package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/MinaMamdouh2/URL-Shortener/foundation/logger"
	"go.uber.org/zap"
)

func main() {

	log, err := logger.New("URL-Shortener-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		// Calling log.Sync here despite of using a deferred call for log.sync
		// since os.Exit() skips defers
		log.Sync()
		os.Exit(1)
	}

}

func run(log *zap.SugaredLogger) error {
	// GOMAXPROCS
	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Shuting down protocol
	shutdown := make(chan os.Signal, 1)
	// We are waiting for "SIGINT" which is "Ctrl+c" or a "SIGTERM" which what
	// will get back from Kubernetes
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	sig := <-shutdown

	log.Infow("shutdown", "status", "shutdown started", "signal", sig)
	defer log.Infow("shutdown", "status", "shutdown completed", "signal", sig)

	return nil
}
