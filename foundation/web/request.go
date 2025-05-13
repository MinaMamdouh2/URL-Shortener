package web

import "github.com/gin-gonic/gin"

// This helper decouples your handlers from the router lib: they just call web.Param(r, "id")
// instead of digging into the route "e.g. Gin" internals internals.
// Param returns the web call parameters from the request.
func Param(c *gin.Context, key string) string {
	// Gin stores path parameters in c.Params, and c.Param does the lookup + defaulting.
	return c.Param(key)
}
