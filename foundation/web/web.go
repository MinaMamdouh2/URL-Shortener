package web

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// A Handler is a type that handles a http request within our own little mini framework.
// We want our handlers to look like that, we don't want the http handler function signature that the mux uses
// "func (w http.ResponseWriter, r *http.Request)", we want a handler function that takes a context and returns an error.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures our context object for each of our http handlers.
// Feel free to add any configuration data/logic on this App struct.
// We are defining this type name App, and what we are gonna do is embed the mux inside of it.
// By embedding the mux in the App that means the App is everything the mux is, it's entire API promotes up to App.
// This is how we steal the mux then we can extend and add to it.
type App struct {
	*gin.Engine
	shutdown chan os.Signal
}

// NewApp creates an App value that handle a set of routes for the application.
// The new function constructs a new App and because this New function is in the web package we are calling it "NewApp"
// because the type is called App if it was called web we would have "New" instead of "NewApp" this is the idiom.
func NewApp(shutdown chan os.Signal) *App {
	engine := gin.New()
	return &App{
		Engine:   engine,
		shutdown: shutdown,
	}
}

// =============================================================================
// Handle sets a handler function for a given HTTP method and path pair to the application server mux.
// We now have overridden the mux's Handle method because now that the method exists that method overrides the promotion
// and it is saying Bill I don't want your mux handle signature I want the App Handle signature this was the case for
// 'httptreemux" but we are using gin so this is not the current case.
func (a *App) Handle(method string, path string, handler Handler) {
	// Now we are going to build the onion.
	// We know at the end of the day, we need a function that takes an Http response writer, and a pointer to http request
	// in order to use the contextMux handle function
	h := func(c *gin.Context) {
		err := handler(c.Request.Context(), c.Writer, c.Request)
		if err != nil {
			// Handle error: could enhance this to log, metrics, etc.
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}

	// We can create all the abstraction in the world but at the end of the day what is implementing the mux is the
	// the context mux
	a.Engine.Handle(method, path, h)
}
