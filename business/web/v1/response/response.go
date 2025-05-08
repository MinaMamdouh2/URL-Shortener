package response

import "errors"

// ErrorDocument is the form used for API responses from failures in the API.
// This represents the ErrorDocument that we send back on all error responses, so the caller will always get the error
// message which again could be 500 message that doesn't reveal anything and could also get a set of fields on what
// doesn't validate properly on the data model and that can be omitted if there is no validation error.
type ErrorDocument struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

// Error is used to pass an error during the request through the application with web specific content.
// This is the trusted Error type and the status that we wanna send back to the client.
type Error struct {
	Err    error
	Status int
}

// NewError wraps a provided error with an HTTP status code.
// This function should be used when handlers encounter expected errors.
func NewError(err error, status int) error {
	return &Error{err, status}
}

// Error implements the error interface. It uses the default message of the wrapped error.
// This is what will be shown in the services' logs.
func (re *Error) Error() string {
	return re.Err.Error()
}

// IsError checks if an error of type Error exists.
func IsError(err error) bool {
	var re *Error
	return errors.As(err, &re)
}

// GetError returns a copy of the Error pointer.
// Here we do 2 operations check if it a trusted and get a copy back if you need.
func GetError(err error) *Error {
	var re *Error
	if !errors.As(err, &re) {
		return nil
	}
	return re
}
