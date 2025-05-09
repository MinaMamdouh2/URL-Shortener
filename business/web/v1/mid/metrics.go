package mid

import (
	"context"
	"net/http"

	"github.com/MinaMamdouh2/URL-Shortener/business/web/v1/metrics"
	"github.com/MinaMamdouh2/URL-Shortener/foundation/web"
)

// Metrics updates program counters.
func Metrics() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// In the beginning of every request, we set the metrics in there
			ctx = metrics.Set(ctx)

			// Make a call
			err := handler(ctx, w, r)

			// Add the number of requests
			n := metrics.AddRequests(ctx)

			// for every 1k request update the number of go routines
			if n%1000 == 0 {
				metrics.AddGoroutines(ctx)
			}

			// If there is an error in flight we return the error.
			// Here we are counting all the errors the trusted and the not trusted ones.
			if err != nil {
				metrics.AddErrors(ctx)
			}

			return err
		}

		return h
	}

	return m
}
