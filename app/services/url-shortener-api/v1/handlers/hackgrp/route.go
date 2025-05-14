package hackgrp

import (
	"net/http"

	"github.com/MinaMamdouh2/URL-Shortener/business/web/v1/auth"
	"github.com/MinaMamdouh2/URL-Shortener/business/web/v1/mid"
	"github.com/MinaMamdouh2/URL-Shortener/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Auth *auth.Auth
}

// Routes adds specific routes for this group using Gin.
// Each group of handlers can define routes related to its domain and bind them to the router that is passed in.
// web.App is everything a gin engine is because of the embedding
func Routes(app *web.App, cfg Config) {
	const version = "v1"
	authen := mid.Authenticate(cfg.Auth)
	ruleAdmin := mid.Authorize(cfg.Auth, auth.RuleAdminOnly)
	app.Handle(http.MethodGet, version, "/hack", Hack)
	app.Handle(http.MethodGet, version, "/hackauth", Hack, authen, ruleAdmin)
}
