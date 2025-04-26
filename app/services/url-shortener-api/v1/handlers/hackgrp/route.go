package hackgrp

import (
	"github.com/MinaMamdouh2/URL-Shortener/foundation/web"
)

// Routes adds specific routes for this group using Gin.
// Each group of handlers can define routes related to its domain
// and bind them to the router that is passed in.
func Routes(app *web.App) {
	app.GET("/hack", Hack)
}
