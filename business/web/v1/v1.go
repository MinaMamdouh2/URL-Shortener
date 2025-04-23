package v1

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Config contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Build string
	// Used value semantics, since Channels in Go are reference types. Even when passed by value,
	// they behave like pointers
	Shutdown chan os.Signal
	// Used pointer semantics here since, loggers carry state (e.g. log levels, output targets, etc.)
	// Passing by pointer ensures we are sharing the same logger instance across the app.
	Log *zap.SugaredLogger
}

// WebAPI constructs a http.Handler with all application routes bound.
// In type systems and functions, we don't use interface type as a return type unless it is the error interface or maybe
// you need the empty interface because you are not using generics.
// Here we are violating this rule by returning the http.Handler but why we are doing that? for one reason he wants us
// to have this conversation with us.
// We shouldn't decouple anything for the callee ideally this should be what is returned here "*httptreemux.ContextMux"
// We should return the concrete type and let the caller decide if it needs to abstract that or not.
func APIMux(cfg APIMuxConfig) http.Handler {
	router := gin.New()

	router.GET("/hack", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Status": "Ok",
		})
	})

	return router
}
