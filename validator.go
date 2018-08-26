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
		FuncMap FuncMap
	}
)

func New() *Validator {
	return &Validator{
		FuncMap: defaultFuncMap,
	}
}

// ValidateStruct validates struct that use tags for fields.
func (v *Validator) ValidateStruct(s interface{}) error {
	if s == nil {
		return nil
	}
	return v.validateStruct(Field{val: reflect.ValueOf(s)})
}

func (v *Validator) validateStruct(field Field) error {
	val := field.val
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
			err = v.validateStruct(Field{name: typeField.Name, val: valueField, parent: &field})
		} else {
			err = v.validateVar(Field{name: typeField.Name, val: valueField, parent: &field}, typeField.Tag.Get(tagName))
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

func (v *Validator) ValidateVar(s interface{}, rawTag string) error {
	return v.validateVar(Field{val: reflect.ValueOf(s)}, rawTag)
}

func (v *Validator) validateVar(field Field, rawTag string) error {
	if rawTag == "-" || rawTag == "" {
		return nil
	}

	tags, err := parseTag(rawTag)
	if err != nil {
		return err
	}

	var errs Errors
	for _, tag := range tags {
		if err := v.validate(field, tag); err != nil {
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

func (v *Validator) validate(field Field, tag Tag) error {
	validate, ok := v.FuncMap[tag.Name]
	if !ok {
		//FIXME define error
		return fmt.Errorf("unknown tag")
	}

	var errs Errors
	if !validate(field, tag.Params...) {
		errs = append(errs, Error{Field: field, Tag: tag})
	}

	var val = field.val
	switch val.Kind() {
	case reflect.Map:
		for _, k := range val.MapKeys() {
			value := val.MapIndex(k)

			var err error
			if value.Kind() == reflect.Struct || (value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct) {
				err = v.validateStruct(Field{name: fmt.Sprintf("[%v]", k), val: value, parent: &field})
			} else {
				err = v.validate(Field{name: fmt.Sprintf("[%v]", k), val: value, parent: &field}, tag)
			}

			if err != nil {
				if es, ok := err.(Errors); ok {
					errs = append(errs, es...)
				} else {
					return err
				}
			}
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			value := val.Index(i)

			var err error
			if value.Kind() == reflect.Struct || (value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct) {
				err = v.validateStruct(Field{name: fmt.Sprintf("[%d]", i), val: value, parent: &field})
			} else {
				err = v.validate(Field{name: fmt.Sprintf("[%d]", i), val: value, parent: &field}, tag)
			}

			if err != nil {
				if es, ok := err.(Errors); ok {
					errs = append(errs, es...)
				} else {
					return err
				}
			}
		}

	case reflect.Interface:
		if val.IsNil() {
			break
		}
		value := val.Elem()

		var err error
		if value.Kind() == reflect.Struct || (value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct) {
			err = v.validateStruct(Field{val: value, parent: &field})
		} else {
			err = v.validate(Field{val: value, parent: &field}, tag)
		}

		if err != nil {
			if es, ok := err.(Errors); ok {
				errs = append(errs, es...)
			} else {
				return err
			}
		}

	case reflect.Ptr:
		if val.IsNil() {
			break
		}
		value := val.Elem()

		var err error
		if value.Kind() == reflect.Struct || (value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct) {
			err = v.validateStruct(Field{val: value, parent: &field})
		} else {
			err = v.validate(Field{val: value, parent: &field}, tag)
		}

		if err != nil {
			if es, ok := err.(Errors); ok {
				errs = append(errs, es...)
			} else {
				return err
			}
		}

	case reflect.Struct:
		err := v.validateStruct(Field{val: val, parent: &field})
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
