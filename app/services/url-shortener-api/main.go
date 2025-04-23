package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/MinaMamdouh2/URL-Shortener/business/web/v1/debug"
	"github.com/MinaMamdouh2/URL-Shortener/foundation/logger"
	"github.com/ardanlabs/conf/v3"
	"go.uber.org/zap"
)

// By default we are setting this variable as develop, but when we build the image, we can overwrite this variable.
var build = "develop"

func main() {

	log, err := logger.New("URL-Shortener-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	// -------------------------------------------------------------------------
	// What we do next is write a run function with the idea that the run function if it fails we are gonna have a
	// single place to write our error logs and return here
	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Errorw("startup", "ERROR", err)
		// Calling log.Sync here despite of using a deferred call for log.sync
		// since os.Exit() skips defers
		log.Sync()
		os.Exit(1)
	}

}

// Here we are passing the logger, this shows the precision we are taking about, we are constructing things in main and
// pass it through the run function.
func run(ctx context.Context, log *zap.SugaredLogger) error {
	// GOMAXPROCS
	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "BUILD", build)

	// -------------------------------------------------------------------------
	// Configuration

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout        time.Duration `conf:"default:5s,env:READ_TIMEOUT"`
			WriteTimeout       time.Duration `conf:"default:10s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s"`
			APIHost            string        `conf:"default:0.0.0.0:3000"`
			DebugHost          string        `conf:"default:0.0.0.0:4000,mask"`
			CORSAllowedOrigins []string      `conf:"default:*"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "Mina Mamdouh",
		},
	}

	const prefix = "URL_SHORTENER"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// -------------------------------------------------------------------------
	// App Starting

	log.Info("starting service ", "version ", build)
	defer log.Info("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info("startup", "config", out)

	// When we are logging the config, this line of code takes the build information we are logging and also puts it
	// into the metrics.
	// Also this will automatically execute the init function for the expvar which adds an endpoint
	// to the default server mux, we can find this build information when we hit "DebugHost/debug/vars"
	expvar.NewString("build").Set(build)

	// -------------------------------------------------------------------------
	// Start Debug Service
	// The first thing we are doing here is launching a go routine whose job is to block for that listenAndServe call.
	// Now we are gonna be telling to use the debug host import which is in our case the default 4000, and we are telling
	// it to call this package debug "debug.mux" which returns a Mux for this.
	// ***IMPORTANT NOTE***
	// Too many people do the following, they use the default server mux as their mux.
	// The "DefaultServeMux" is a singleton, the reason people are using this is because with that in place all you have
	// to do is include the '_ "net/http/pprof"' package and there is an init function in that package and it is adding
	// 5 default endpoints to the default server mux, so essentially this import with the blank identifier is asking
	// the compiler to execute that init code to add these endpoints to the default.
	// There is another endpoint that I want to add for debugging and it comes form another package "expvar" package.
	// The "expvar" package also has an init function that add another endpoint for the default server mux.
	// The problem is any package that you import can have an init function that adds endpoints to the default server mux.
	// We should only be exposing endpoints that we know about "This is a security bug".

	go func() {
		log.Info(ctx, "startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)
		// Here people usually use http.DefaultServeMux instead of our "debug.Mux()" but as we explained above
		// this is a security vulnerability
		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux()); err != nil {
			log.Error(ctx, "shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "msg", err)
		}
	}()

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
