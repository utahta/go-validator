package validator

import (
	"fmt"
	"reflect"
)

type Validator struct {
}

func New() *Validator {
	return &Validator{}
}

// ValidateStruct validates struct that use tags for fields.
func (v *Validator) ValidateStruct(s interface{}) error {
	if s == nil {
		return nil
	}

	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		//FIXME define error
		return fmt.Errorf("invalid argument error")
	}

	return nil
}

func (v *Validator) ValidateVar(s interface{}, tags string) error {
	return nil
}
