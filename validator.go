package validator

import (
	"fmt"
	"reflect"
)

const (
	tagName = "valid"
)

type (
	Validator struct {
	}
)

func New() *Validator {
	return &Validator{}
}

// ValidateStruct validates struct that use tags for fields.
func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validateStruct(s, nil)
}

func (v *Validator) validateStruct(s interface{}, parent *vField) error {
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

	var errs Errors
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		if typeField.PkgPath != "" {
			continue // private field
		}
		if typeField.Tag.Get(tagName) == "-" {
			continue
		}

		valueField := val.Field(i)
		if valueField.Kind() == reflect.Interface {
			valueField = valueField.Elem()
		}

		var err error
		if valueField.Kind() == reflect.Struct || (valueField.Kind() == reflect.Ptr && valueField.Elem().Kind() == reflect.Struct) {
			err = v.validateStruct(valueField.Interface(), &vField{name: typeField.Name, parent: parent})
		} else {
			err = v.validateVar(valueField.Interface(), typeField.Tag.Get(tagName), vField{name: typeField.Name, parent: parent})
		}
		if err != nil {
			if es, ok := err.(Errors); ok {
				errs = append(errs, es...)
			} else {
				return err
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (v *Validator) ValidateVar(s interface{}, tag string) error {
	return v.validateVar(s, tag, vField{})
}

func (v *Validator) validateVar(s interface{}, tag string, field vField) error {
	if tag == "-" {
		return nil
	}

	var errs Errors
	for _, t := range parseTag(tag) {
		err := v.validate(s, t, field)
		if err != nil {
			if e, ok := err.(Error); ok {
				errs = append(errs, e)
			} else {
				return err
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (v *Validator) validate(s interface{}, tag vTag, field vField) error {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	fmt.Printf("%v %v \n", field.Name(), val.Kind())

	return nil
}
