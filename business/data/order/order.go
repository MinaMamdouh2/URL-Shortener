// Package order provides support for describing the ordering of data.
// Ordering is all about knowing that you wanna go ascending or descending.
package order

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/MinaMamdouh2/URL-Shortener/foundation/validate"
)

// Set of directions for data ordering.
const (
	ASC  = "ASC"
	DESC = "DESC"
)

var directions = map[string]string{
	ASC:  "ASC",
	DESC: "DESC",
}

// =============================================================================
// By represents a field used to order by and direction.
// We define a type called By, where you specify the field and the direction
type By struct {
	Field     string
	Direction string
}

// NewBy constructs a new By value with no checks.
func NewBy(field string, direction string) By {
	return By{
		Field:     field,
		Direction: direction,
	}
}

// =============================================================================
// Parse constructs a order.By value by parsing a string in the form
// of "field,direction".
// Parse is a way that we are going to have a valid by.
// Parse is taking a request, if we wanna move this to the foundation layer, we can't use request because we can't assume
// that's what our application layer going to be, but at the business layer we can cheat a bit because we are only using
// this function in the app layer against a request looking for the query string.
// e.g. on query "orderBy=name,desc"
func Parse(r *http.Request, defaultOrder By) (By, error) {
	//TODO: Test the case for "orderBy=name,asc&orderBy=createdAt,desc"
	v := r.URL.Query().Get("orderBy")

	if v == "" {
		return defaultOrder, nil
	}

	orderParts := strings.Split(v, ",")

	var by By
	switch len(orderParts) {
	// len == 1, only the field was provided e.g. "name", default direction will be "ASC	"
	case 1:
		by = NewBy(strings.Trim(orderParts[0], " "), ASC)
	// len == 2, the first field is the field e.g. "name", the second is the direction
	// TODO: use "TrimSpace"
	case 2:
		by = NewBy(strings.Trim(orderParts[0], " "), strings.Trim(orderParts[1], " "))
	// More than 3 parts, then it's malformed
	default:
		return By{}, validate.NewFieldsError(v, errors.New("unknown order field"))
	}
	// Checks wether the parsed direction is in the directions map
	if _, exists := directions[by.Direction]; !exists {
		return By{}, validate.NewFieldsError(v, fmt.Errorf("unknown direction: %s", by.Direction))
	}

	return by, nil
}
