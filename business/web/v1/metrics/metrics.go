// We only need metrics for the web app, if we are doing a cli app we don't care about metrics, and since we need the
// metrics to really help maintain, manage, and debug the app, we could leverage the context for managing the metrics.
// Just like config, people can suddenly add metrics anywhere in the code to the point
// where the metrics are just a mess.
// so we wanna control that, we wanna have a package for at least we have an API for all the metrics we plan to gather
// that we if we do want to gather the metrics data in multiple places at least we know what is being gathered

// Package metrics constructs the metrics the application will track.
package metrics

import (
	"context"
	"expvar"
	"runtime"
)

// This holds the single instance of the metrics value needed for collecting metrics.
// The expvar package is already based on a singleton for the different metrics
// that are registered with the package so there isn't much choice here.
// The expvar package use singletons underneath, so anytime you have to work with a package API
// that is based on a singleton, it always creates complexity and problems and we kinda have this now.
// Remember the expvar package provides endpoints for the metrics so it is storing sort of globally
// the metrics counters that we are recording.
// So it is for someone to by pass what we are doing and we don't want that.
var m *metrics

// metrics represents the set of metrics we gather.
// These fields are safe to be accessed concurrently thanks to expvar. No extra abstraction is required.
// We are here declaring a variable package variable that is using pointer semantics that will hold this metrics a value
// of metrics type which have all the fields of metrics values we wanna capture using the expvar package.
type metrics struct {
	goroutines *expvar.Int
	requests   *expvar.Int
	errors     *expvar.Int
	panics     *expvar.Int
}

// init constructs the metrics value that will be used to capture metrics.
// The metrics value is stored in a package level variable since everything inside of expvar is registered as a singleton. The use of once will make
// sure this initialization only happens once.
// In general we don't want to have package level variables like "m", but when we can make it unexported and if we don't
// need anything from config and we don't care about the order of initialization we can get away with it.
// In this case we don't care when we construct metrics as long it happens before main and that will never change.
// Now if we did need some configuration to do this, we would need capital "Init" function and pass that in and from main
// we will need to make sure we are calling that.
func init() {
	m = &metrics{
		goroutines: expvar.NewInt("goroutines"),
		requests:   expvar.NewInt("requests"),
		errors:     expvar.NewInt("errors"),
		panics:     expvar.NewInt("panics"),
	}
}

// =============================================================================
type ctxKey int

const key ctxKey = 1

// "Set" sets the metrics data into the context.
// What Set does is that it takes our singleton and stores it in the context for each request.
// What it does it makes sure we are not using the singleton directly that the system is pulling the metrics from the
// context so it gives us some flexibility.
func Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, m)
}

// AddGoroutines refreshes the goroutine metric.
func AddGoroutines(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		g := int64(runtime.NumGoroutine())
		v.goroutines.Set(g)
		return g
	}

	return 0
}

// AddRequests increments the request metric by 1.
func AddRequests(ctx context.Context) int64 {
	v, ok := ctx.Value(key).(*metrics)
	if ok {
		v.requests.Add(1)
		return v.goroutines.Value()
	}

	return 0
}

// AddErrors increments the errors metric by 1.
func AddErrors(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.errors.Add(1)
		return v.errors.Value()
	}

	return 0
}

// AddPanics increments the panics metric by 1.
func AddPanics(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.panics.Add(1)
		return v.panics.Value()
	}

	return 0
}
