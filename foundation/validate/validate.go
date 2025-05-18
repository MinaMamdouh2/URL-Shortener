// Package validate contains the support for validating models.
package validate

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	// This one comes from the Go team, and it has the ability to do model level validations against tags and also allows
	// you to support multiple languages for response and if you need multi language error support for messaging then
	// you can't do this using init, you will have to construct in main and initialize those languages that you want
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// validate holds the settings and caches for validating request struct values.
var validate *validator.Validate

// translator is a cache of locale and translation information.
var translator ut.Translator

func init() {

	// Instantiate a validator.
	validate = validator.New()

	// Create a translator for english so the error messages are more human-readable than technical.
	translator, _ = ut.New(en.New(), en.New()).GetTranslator("en")

	// Register the english error messages for use.
	en_translations.RegisterDefaultTranslations(validate, translator)

	// Use JSON tag names for errors instead of Go struct names.
	// Why are we using this?
	// By default when you do, "FirstName string `json:"first_name" validate:"required"`"
	// and you call validate.Struct(&user) with FirstName == "", the error object's "Field()" will return
	// the Go struct's name, "FirstName". That's not ideal for API responses, because clients only know
	// about "first_name".
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// "fld.Tag.Get("json")"
		// Reads the struct tag for the “json” key.
		// E.g. for `json:"first_name,omitempty"`, it returns "first_name,omitempty".
		// strings.SplitN(..., ",", 2)[0]
		// Splits on the first comma.
		// Takes the part before any comma, which is the actual JSON key: "first_name".
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// If the JSON tag is "-", that field is explicitly omitted from JSON.
		if name == "-" {
			// Returning "" causes the validator to skip naming that field in error reports.
			return ""
		}
		return name
	})
}

// Check validates the provided model against it's declared tags.
func Check(val any) error {

	// This api has a struct function which takes a data model with those validation tags and looks to see that the data
	// inside that val is holding true to these tags, look at the docs for more info.
	if err := validate.Struct(val); err != nil {
		// 1) Did the validator return field specific errors?
		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err // Some other error-pass it through
		}

		// 2) Convert each ValidationError into FieldError
		var fields FieldErrors
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),               // e.g. "email"
				Err:   verror.Translate(translator), // e.g. "email is a required field"
			}
			fields = append(fields, field)
		}
		// 3) Return the slice of FieldErrors as error
		return fields
	}
	// No validation error
	return nil
}
