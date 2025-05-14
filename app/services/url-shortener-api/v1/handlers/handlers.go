package handlers

import (
	"github.com/MinaMamdouh2/URL-Shortener/app/services/url-shortener-api/v1/handlers/checkgrp"
	"github.com/MinaMamdouh2/URL-Shortener/app/services/url-shortener-api/v1/handlers/hackgrp"
	v1 "github.com/MinaMamdouh2/URL-Shortener/business/web/v1"
	"github.com/MinaMamdouh2/URL-Shortener/foundation/web"
)

type Routes struct{}

// Add implements the RouterAdder interface
func (Routes) Add(app *web.App, apiCfg v1.APIMuxConfig) {
	hackgrp.Routes(app, hackgrp.Config{
		Auth: apiCfg.Auth,
	})

	checkgrp.Routes(app, checkgrp.Config{
		Build: apiCfg.Build,
	})
}
