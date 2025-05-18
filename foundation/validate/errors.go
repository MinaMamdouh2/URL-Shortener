// Validate package job is to validate data models
package validate

import (
	"encoding/json"
	"errors"
)

// FieldError is used to indicate an error with a specific request field.
// We could argue that this is a mistake, that we have json tags here at the foundation layer. This is the information
// that you could send back to the client in JSON, we could argue that this is a mistake since the app layer may want
// to use XML and we are enforcing a strong policy here.
type FieldError struct {
	Field string `json:"field"`
	Err   string `json:"error"`
}

// FieldErrors represents a collection of field errors.
type FieldErrors []FieldError

// NewFieldsError creates an fields error.
func NewFieldsError(field string, err error) error {
	return FieldErrors{
		{
			Field: field,
			Err:   err.Error(),
		},
	}
}

// Error implements the error interface.
// The implementation of the error interface is on the FieldErrors, and we are using value semantics in the implementation
// when the underlying type is a slice and pointer semantics other wise
func (fe FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}

// Fields returns the fields that failed validation
func (fe FieldErrors) Fields() map[string]string {
	m := make(map[string]string)
	for _, fld := range fe {
		m[fld.Field] = fld.Err
	}
	return m
}

// IsFieldErrors checks if an error of type FieldErrors exists.
// Basically, it returns true if err wraps a "FieldErrors"
func IsFieldErrors(err error) bool {
	var fe FieldErrors
	return errors.As(err, &fe)
}

// GetFieldErrors returns a copy of the FieldErrors pointer.
// Unwraps & returns the concrete "FieldErrors" slice or nil if not present
func GetFieldErrors(err error) FieldErrors {
	var fe FieldErrors
	if !errors.As(err, &fe) {
		return nil
	}
	return fe
}
