// Ideally what we wanna do is generate a trace-id maybe we can use uuid as a start
// but since this is for basically logging debugging and tracing the app.
// We don't need to be passing this around, this is a value we could hide in a context
// because it is specific to log that information and if it is not there nothing breaks.
package web

import (
	"context"
	"time"
)

// Unexported type for injecting Gin param context into request context
// This is done since I want to be able to use abstractions like "Param" which takes "r *http.Request"
// making the helper route agnostic
type ctxParamKey int

const paramKey ctxParamKey = 1

// Unexported type so we enforce people to use our API.
type ctxKey int

// Anytime we are gonna store something in the context, or we wanna retrieve something from the context.
// we need a key, the key should be based on a type and a value, and we wanna our own unique key type
// that way no one can override our key values.
// The value 1 is irrelevant only the type that matters
const key ctxKey = 1

// Values represent state for each request.
// Concrete type called "Values" and what we are gonna do is that everything that is going to be foundational
// with the web package, we wanna sort of track from the context perspective will be in this one type.
type Values struct {
	// We are gonna store the TraceID for every request
	TraceID string
	// We are gonna get the current time of the request when it came in
	Now time.Time
	// We are gonna set the status code when we are sending the request back
	StatusCode int
}

// =============================================================================
// Bill sees the context api for storing and retrieving values is too complex. So what he does is he defines a getter
// and a setter api. Bill in the Go class always say "No getters and no setters in GO" and that's true
// but when it comes to setting and getting values out of the context you can do that.
// Also if the value is not in the context we can return some default value so we don't worry about error checking
// or anything blowing up.

// Also this way the code never panics due to missing context values
// No caller has to do "ok" checking - the logic is encapsulated here

// GetValues returns the values from the context.
func GetValues(ctx context.Context) *Values {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return &Values{
			TraceID: "00000000-0000-0000-0000-000000000000",
			Now:     time.Now(),
		}
	}

	return v
}

// GetTraceID returns the trace id from the context.
func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return "00000000-0000-0000-0000-000000000000"
	}

	return v.TraceID
}

// GetTime returns the time from the context.
func GetTime(ctx context.Context) time.Time {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return time.Now()
	}

	return v.Now
}

// =============================================================================
// Internal Setters (Unexported)
// These are used within the "web" framework only (not exposed to users)

func setStatusCode(ctx context.Context, statusCode int) {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return
	}

	v.StatusCode = statusCode
}

// Injects a "Values" object into a context
// Remember, values in context are immutable references meaning when you "add" a value to it
// you are actually creating a new derived context but the struct "*Values" itself is a pointer
// so its fields can be mutated safely
func setValues(ctx context.Context, v *Values) context.Context {
	return context.WithValue(ctx, key, v)
}
