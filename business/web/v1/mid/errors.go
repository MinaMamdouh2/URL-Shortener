package mid

import (
	"context"
	"net/http"

	"github.com/MinaMamdouh2/URL-Shortener/business/web/v1/response"
	"github.com/MinaMamdouh2/URL-Shortener/foundation/web"
	"go.uber.org/zap"
)

// Errors handles errors coming out of the call chain.
// It detects normal application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			err := handler(ctx, w, r)

			// If there is no error, go to the next middleware
			if err == nil {
				return nil
			}
			// First thing is to log the error
			log.Error(ctx, "message", "msg", err)
			// Second construct the error document that we are gonna send back and what status that we are gonna use.
			var er response.ErrorDocument
			var status int
			// This switch will determine what the ErrorDocument looks like and what status we are gonna be using
			switch {
			// This is the case for trusted error.
			case response.IsError(err):

				// The GetError will bring me back the concrete value so we can inspect it.
				reqErr := response.GetError(err)
				// Here we technically say we blindly trust the messaging coming back here.
				// Basically What the handler send we will respond with it back to the client.
				er = response.ErrorDocument{
					Error: reqErr.Error(),
				}
				status = reqErr.Status

			// We don't know what the error is, it is a non trusted error
			default:
				er = response.ErrorDocument{
					Error: http.StatusText(http.StatusInternalServerError),
				}
				status = http.StatusInternalServerError
			}

			// Here sending this error document out.
			// Here web respond could fail and the error will be sent back to the framework
			if err := web.Respond(ctx, w, er, status); err != nil {
				return err
			}

			// If we receive the shutdown err we need to return it, it will return back
			// to back to the base handler "app.Handle" to shut down the service.
			// Here we are checking if there is a shutdown error given to us and we are sending back to the framework.
			if web.IsShutdown(err) {
				return err
			}

			return nil
		}
		return h
	}
	return m
}
