package hackgrp

import (
	"context"
	"encoding/json"
	"net/http"
)

// Hack handles the /hack route using Gin context.
func Hack(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	status := struct {
		Status string
	}{
		Status: "Ok",
	}

	return json.NewEncoder(w).Encode(status)
}
