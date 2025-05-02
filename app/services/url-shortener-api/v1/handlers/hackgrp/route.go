package hackgrp

import (
	"net/http"

	"github.com/MinaMamdouh2/URL-Shortener/foundation/web"
)

// Routes adds specific routes for this group using Gin.
// Each group of handlers can define routes related to its domain and bind them to the router that is passed in.
// web.App is everything a gin engine is because of the embedding
func Routes(app *web.App) {
	app.Handle(http.MethodGet, "/hack", Hack)
}
