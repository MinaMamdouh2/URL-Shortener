// Right now if we get a panic what is going to happen is the code is going to unwind all way back to the http package
// which will essentially send a 500 and terminate the go routine, that is really not what we want if something panics
// we wanna capture it immediately so we can make sure that the rest of the middleware works when we add this middleware
// we are solving the problem.
package mid

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/MinaMamdouh2/URL-Shortener/foundation/web"
)

// Panics recovers from panics and converts the panic to an error so it is reported in Metrics and handled in Errors.
// Even though this is not doing any logging for consistency we are still using a function that returns a middleware.
func Panics() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			// Defer a function to recover from a panic and set the err return variable after the fact.
			defer func() {
				if rec := recover(); rec != nil {
					// If we panic we want the stack trace, there will be good information there and we wanna log it
					// regardless
					trace := debug.Stack()
					// We are going to return the error from the defer by using the named return.
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))
				}
			}()
			// This is the call to our handler, before this call we are setting our defer which it's job to call the
			// builtin function recover "When you want to stop a panic, you call recover inside defer" that defer
			// has to be setup before we call the handler and we are doing that
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
