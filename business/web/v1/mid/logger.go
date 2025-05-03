package mid

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MinaMamdouh2/URL-Shortener/foundation/web"
	"go.uber.org/zap"
)

// We will create a Logger that accepts a logger and have it to return a Middleware function and by doing this we are
// gonna leverage Go closures to get the logger into the "h"
func Logger(log *zap.SugaredLogger) web.Middleware {

	m := func(handler web.Handler) web.Handler {
		// The whole point here is this has to call the handler function.
		// so ideally since we are going to be wrapping it, we can say fine let's write the function.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// We log the start and end of the request. so we can make sure that
			// the goroutine we created is completed, so when we find out that
			// a request has started but not completed there maybe a data race, leak or a block.

			// How should we pass the logger here? The logger signature is locked in we can't change it and the h signature
			// is also locked in what can we do? we can't hide the logger into the context that is WRONG!!
			// 1- One way we can do that is turn "Logger" into a method, we could define a type create a "Logger" field
			// turn this into a method and pass the logger through the receiver but that's a lot of ceremony because now
			// you have to construct that value before you can pass the middleware to the framework.
			// 2- will see it above

			// Logging HERE - STARTED
			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
			}

			log.Infow("request started", "method", r.Method, "path", path, "remoteaddr", r.RemoteAddr)

			// Call the next handler
			err := handler(ctx, w, r)

			// Logging HERE - COMPLETED
			log.Infow("request completed", "method", r.Method, "path", path, "remoteaddr", r.RemoteAddr)

			return err
		}

		return h
	}
	return m
}
