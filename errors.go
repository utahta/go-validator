package validator

import (
	"fmt"
	"strings"
)

type (
	// Error is an interface that represents a validation error
	Error interface {
		// Field returns a validating field.
		Field() Field

		// Tag returns a validation tag.
		Tag() Tag

		// Error returns an error message string.
		Error() string
	}

	fieldError struct {
		field Field
		tag   Tag

		// err is an internal error.
		err error

		// customMessage is a custom error message. TODO:
		customMessage string

		// suppressErrorFieldValue suppress output of field value.
		suppressErrorFieldValue bool
	}

	// Errors represents validation errors
	Errors []Error
)

// ToErrors converts an error to the validation Errors.
func ToErrors(err error) (Errors, bool) {
	es, ok := err.(Errors)
	return es, ok
}

func (e *fieldError) Field() Field {
	return e.field
}

func (e *fieldError) Tag() Tag {
	return e.tag
}

func (e *fieldError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: an internal error occurred in '%s': %v", e.field.Name(), e.tag, e.err)
	}

	if e.suppressErrorFieldValue {
		return fmt.Sprintf("%s: The value does validate as '%s'", e.field.Name(), e.tag)
	}
	return fmt.Sprintf("%s: '%s' does validate as '%s'", e.field.Name(), e.field.ShortString(), e.tag)
}

func (es Errors) Error() string {
	var s []string
	for _, e := range es {
		s = append(s, e.Error())
	}
	return strings.Join(s, ";")
}
