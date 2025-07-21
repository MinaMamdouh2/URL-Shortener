package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type validator interface {
	Validate() error
}

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

// Decode reads the body of an HTTP request looking for a JSON document.
// The body is decoded into the provided value.
// If the provided value is a struct then it is checked for validation tags.
// If the value implements a validate function, it is executed.
// It is a simple function that is creating a abstraction around calling the JSON package for using JSON decoding.
func Decode(r *http.Request, val any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	// Here were using json.Unmarshal under the hood, if the JSON is malformed or types don't line up
	// you sent a string to an int field it returns an error
	if err := decoder.Decode(val); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}
	// We still wanna do some validation, so this is a little trick, we defined an interface named validator with one
	// active behavior called "Validate" and what we do is that we do the decoding and then we ask does this data model
	// that is passed into me, does it have a Validate method and if it does execute it.
	// So this gives the app layer developer to provide an extra layer of validation if they want.
	if v, ok := val.(validator); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("unable to validate payload: %w", err)
		}
	}

	return nil
}
