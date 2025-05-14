// The idea here that we are ready to receive traffic, if we are not responding properly to these calls on Readiness
// means we are not getting on traffic and liveness means we are getting restarted.
// Package checkgrp maintains the group of handlers for health checking.
package checkgrp

import (
	"context"
	"net/http"
	"os"

	"github.com/MinaMamdouh2/URL-Shortener/foundation/web"
	"go.uber.org/zap"
)

// Handlers manages the set of check endpoints.
// We are defining a new types and we are gonna make our handlers as methods of that type and that because we want to
// pass db connection.
// In middlewares we use closures to pass things in because we don't want to deal with construction of things in order
// to leverage a middleware but here a construction isn't a big deal at all and gives us more flexibility to update and
// manage states as we need
type Handlers struct {
	build string
	log   *zap.SugaredLogger
}

// New constructs a Handlers api for the check group.
func New(build string, log *zap.SugaredLogger) *Handlers {
	return &Handlers{
		build: build,
		log:   log,
	}
}

// Readiness checks if the database is ready and if not will return a 500 status.
// Do not respond by just returning an error because further up in the call stack it will interpret that as a non-trusted error.
func (h *Handlers) Readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := "ok"
	statusCode := http.StatusOK

	// Here the status is important
	data := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}
	h.log.Infow("readiness", "status", status)
	return web.Respond(ctx, w, data, statusCode)
}

// Liveness returns simple status info if the service is alive.
// If the app is deployed to a Kubernetes cluster, it will also return pod, node, and
// namespace details via the Downward API.
// The Kubernetes environment variables need to be set within your Pod/Deployment manifest.
func (h *Handlers) Liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	// We define this literal struct with all these fields, and we make sure we respond with 200 with these information.
	// Bill loves this idea because this endpoint gives you the ability to extract some information on the running service
	// without costing you anything
	data := struct {
		Status     string `json:"status,omitempty"`
		Build      string `json:"build,omitempty"`
		Host       string `json:"host,omitempty"`
		Name       string `json:"name,omitempty"`
		PodIP      string `json:"podIP,omitempty"`
		Node       string `json:"node,omitempty"`
		Namespace  string `json:"namespace,omitempty"`
		GOMAXPROCS string `json:"GOMAXPROCS,omitempty"`
	}{
		Status: "up",
		Build:  h.build,
		Host:   host,
		// These environment variables don't exist but we have a way of creating them in the yaml for K8s
		Name:       os.Getenv("KUBERNETES_NAME"),
		PodIP:      os.Getenv("KUBERNETES_POD_IP"),
		Node:       os.Getenv("KUBERNETES_NODE_NAME"),
		Namespace:  os.Getenv("KUBERNETES_NAMESPACE"),
		GOMAXPROCS: os.Getenv("GOMAXPROCS"),
	}

	h.log.Infow("liveness", "status", "OK")

	// This handler provides a free timer loop.
	// This is gonna get called on whatever interval you specify, instead of adding go routines that run on timers and
	// do things always keep at the back of your head as long this responds fast enough you get a free timer loop
	// Here the data is important
	return web.Respond(ctx, w, data, http.StatusOK)
}
