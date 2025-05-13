package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// This helper decouples your handlers from the router lib: they just call web.Param(r, "id")
// instead of digging into the route "e.g. Gin" internals.
// Param returns the web call parameters from the request.
func Param(r *http.Request, key string) string {
	// r.Context().Value(paramKey), fetches the value stored under paramKey in the request’s context.
	// Earlier, we injected Gin’s c.Params slice via middleware.
	// Then we are doing ".([]gin.Param)"
	if ps, _ := r.Context().Value(paramKey).([]gin.Param); len(ps) > 0 {
		for _, p := range ps {
			if p.Key == key {
				return p.Value
			}
		}
	}
	return ""
}
