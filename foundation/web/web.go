package web

import (
	"os"

	"github.com/gin-gonic/gin"
)

// App is the entrypoint into our application and what configures our context object for each of our http handlers.
//
//	Feel free to add any configuration data/logic on this App struct.
//
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
