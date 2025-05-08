package web

import (
	"errors"
)

// Go has an error convention for error type, it needs to end in the word "Error", so this is our shutdownError type.
// Ideally it is also best to make it unexported if you can, so nobody can do type assertions against it and you can
// give it an API if you will.
// Pointer semantics when you are using a struct for your error type, value semantics when you are using a slice for
// error types.

// shutdownError is a type used to help with the graceful termination of the service.
type shutdownError struct {
	Message string
}

// NewShutdownError returns an error that causes the framework to signal a graceful shutdown.
// This is a factory function allows you to create errors on one line of code. sometimes on the return itself.
// This helps with readability in terms of constructing errors.
func NewShutdownError(message string) error {
	return &shutdownError{message}
}

// Error is the implementation of the error interface.
func (se *shutdownError) Error() string {
	return se.Message
}

// IsShutdown checks to see if the shutdown error is contained in the specified error value.
// The Go team added 2 functions along time ago to the errors package "As" & "Is".
// The "Is" let's you compare two error values to see if they are the same.
// - We tend to use is when we are dealing with error variable.
// The "As" provides 2 purposes
// - It allows you to check if an error of a single type is being stored
// - It gives you a copy of that.
func IsShutdown(err error) bool {
	// Here we are asking if is there an error stored inside the error interface that is a pointer of type shutdownError.
	// Also it checks if an error wraps a certain error type "shutdownError in our case"
	var se *shutdownError
	return errors.As(err, &se)
}
