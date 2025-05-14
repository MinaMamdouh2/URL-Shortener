package checkgrp

import (
	"net/http"

	"github.com/MinaMamdouh2/URL-Shortener/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	hdl := New(cfg.Build)
	// Bill didn't want to include middleware within the scope of these routes, we are binding them to port 3000
	// It is important they are on port 3000 because our application traffic goes through that port.
	// The problem that if I have the middleware getting involved with these 2 routes, we will get a lot of noise
	// in terms of logging, metrics ...
	// We added a method called "HandleNoMiddleware" as opposed just the handle.
	app.HandleNoMiddleware(http.MethodGet, version, "/readiness", hdl.Readiness)
	app.HandleNoMiddleware(http.MethodGet, version, "/liveness", hdl.Liveness)

}
