package v1

import (
	"os"

	"github.com/MinaMamdouh2/URL-Shortener/business/web/v1/mid"
	"github.com/MinaMamdouh2/URL-Shortener/foundation/web"
	"go.uber.org/zap"
)

// Config contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Build string
	// Used value semantics, since Channels in Go are reference types. Even when passed by value,
	// they behave like pointers
	Shutdown chan os.Signal
	// Used pointer semantics here since, loggers carry state (e.g. log levels, output targets, etc.)
	// Passing by pointer ensures we are sharing the same logger instance across the app.
	Log *zap.SugaredLogger
}

// RouteAdder defines behavior that sets the routes to bind for an instance
// of the service.
type RouteAdder interface {
	Add(app *web.App, cfg APIMuxConfig)
}

// WebAPI constructs a http.Handler with all application routes bound.
// In type systems and functions, we don't use interface type as a return type unless it is the error interface or maybe
// you need the empty interface because you are not using generics.
// Here we are violating this rule by returning the http.Handler but why we are doing that? for one reason he wants us
// to have this conversation with us.
// We shouldn't decouple anything for the callee ideally this should be what is returned here "*gin.Engine"
// We should return the concrete type and let the caller decide if it needs to abstract that or not.
// Now we are passing the config and pass it some concrete value that knows how to Add "routeAdder", this creates a level
// of abstraction here and "routeAdder" is going to call Add passing in the mux we wanna bind the route to and the config.
// We added the RouteAdder interface to add the ability to isolate the routes that are being added to the mux at the
// group level but also it gives you the ability to build different binaries for different sets of routes eventually.
// Once we created the web package in the foundation layer, the app embeds a mux so it's APIs are promoted up to the app.
// Now we can initialize the app here and return a pointer to it.
// Also by constructing the web app we sort of hidden the implementation of the mux we are using.
// This will later let us rip out the gin mux and going to use the standard library mux because we are gonna create that
// abstraction.
// Notice we are not creating an abstraction to an interface we are creating an abstraction to the concrete type.
func APIMux(cfg APIMuxConfig, routeAdder RouteAdder) *web.App {
	app := web.NewApp(cfg.Shutdown, mid.Logger(cfg.Log), mid.Errors(cfg.Log))

	routeAdder.Add(app, cfg)

	return app
}
