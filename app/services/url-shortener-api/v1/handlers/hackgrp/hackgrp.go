package hackgrp

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"

	"github.com/MinaMamdouh2/URL-Shortener/business/web/v1/response"
)

// if it is OK but we don't want the handler dealing with error handling, we would rather push the error handling back
// into the business layer if we could or somewhere outside the handler for consistency
// Hack handles the /hack route using Gin context.
func Hack(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	if n := rand.Intn(100) % 2; n == 0 {
		return response.NewError(errors.New("TRUSTED ERROR"), http.StatusBadRequest)
	}
	status := struct {
		Status string
	}{
		Status: "Ok",
	}

	return json.NewEncoder(w).Encode(status)
}
