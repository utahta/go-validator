package validator

import (
	"fmt"
	"strings"
)

type (
	// Error represents a validation error
	Error struct {
		// Field is a validation field.
		Field Field

		// Tag is a validation tag.
		Tag Tag

		// Err is a internal error.
		Err error

		// CustomMessage is a custom error message. TODO:
		CustomMessage string

		// SuppressErrorFieldValue suppress output of field value.
		SuppressErrorFieldValue bool
	}

	// Errors represents validation errors
	Errors []Error
)

// ToErrors converts an error to the validation Errors.
func ToErrors(err error) (Errors, bool) {
	es, ok := err.(Errors)
	return es, ok
}

func (e Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: an error occurred in '%s': %v", e.Field.Name(), e.Tag, e.Err)
	}

	if e.SuppressErrorFieldValue {
		return fmt.Sprintf("%s: The value does validate as '%s'", e.Field.Name(), e.Tag)
	}
	return fmt.Sprintf("%s: '%s' does validate as '%s'", e.Field.Name(), e.Field.ShortString(), e.Tag)
}

func (es Errors) Error() string {
	var s []string
	for _, e := range es {
		s = append(s, e.Error())
	}
	return strings.Join(s, ";")
}
