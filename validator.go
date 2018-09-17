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
		FuncMap                 FuncMap
		SuppressErrorFieldValue bool //TODO

		tagCache    *tagCache
		structCache *structCache
	}
)

func New() *Validator {
	return &Validator{
		FuncMap:     defaultFuncMap,
		tagCache:    newTagCache(),
		structCache: newStructCache(),
	}
}

// SetFunc sets validator function.
func (v *Validator) SetFunc(rawTag string, fn Func) {
	v.FuncMap[rawTag] = with(fn)
}

// ValidateStruct validates struct that use tags for fields.
func (v *Validator) ValidateStruct(s interface{}) error {
	if s == nil {
		return nil
	}
	value := reflect.ValueOf(s)
	return v.validateStruct(Field{origin: value, current: value})
}

func (v *Validator) validateStruct(field Field) error {
	val := field.Value()
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("struct type required")
	}

	valueType := val.Type()
	fieldCaches, hasCache := v.structCache.Load(valueType)

	var errs Errors
	for i := 0; i < val.NumField(); i++ {
		if !hasCache {
			typeField := valueType.Field(i)
			fieldCaches = append(fieldCaches, fieldCache{
				isPrivate: typeField.PkgPath != "", // private field
				tagValue:  typeField.Tag.Get(tagName),
				name:      typeField.Name,
			})
		}
		if fieldCaches[i].isPrivate || fieldCaches[i].tagValue == "-" {
			continue
		}

		originField := val.Field(i)
		valueField := v.extractVar(originField)

		if err := v.validateVar(Field{name: fieldCaches[i].name, origin: originField, current: valueField, parent: field}, fieldCaches[i].tagValue); err != nil {
			if es, ok := err.(Errors); ok {
				errs = append(errs, es...)
			} else {
				return err
			}
		}

		if fieldCaches[i].tagValue == "" {
			if valueField.Kind() == reflect.Struct || (valueField.Kind() == reflect.Ptr && valueField.Elem().Kind() == reflect.Struct) {
				if err := v.validateStruct(Field{name: fieldCaches[i].name, origin: originField, current: valueField, parent: field}); err != nil {
					if es, ok := err.(Errors); ok {
						errs = append(errs, es...)
					} else {
						return err
					}
				}
			}
		}
	}

	if !hasCache {
		v.structCache.Store(valueType, fieldCaches)
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (v *Validator) ValidateVar(s interface{}, rawTag string) error {
	value := reflect.ValueOf(s)
	return v.validateVar(Field{origin: value, current: v.extractVar(value)}, rawTag)
}

func (v *Validator) validateVar(field Field, rawTag string) error {
	if rawTag == "-" || rawTag == "" {
		return nil
	}

	tags, err := v.tagParse(rawTag)
	if err != nil {
		return err
	}

	var errs Errors
	for _, t := range tags {
		if err := v.validate(field, t); err != nil {
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
	var errs Errors
	if tag.Enable {
		valid, err := tag.validateFn(field, FuncOption{validator: v, Params: tag.Params, Optional: tag.Optional})
		if err != nil {
			return fmt.Errorf("validateFn: %v in %s %s", err, field.FullName(), tag.String())
		}
		if !valid {
			errs = append(errs, Error{Field: field, Tag: tag, SuppressErrorFieldValue: v.SuppressErrorFieldValue})
		}
	}

	tag.Enable = true // for dig

	var val = field.Value()
	switch val.Kind() {
	case reflect.Map:
		for _, k := range val.MapKeys() {
			value := val.MapIndex(k)

			var err error
			if value.Kind() == reflect.Struct || (value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct) {
				err = v.validateStruct(Field{name: fmt.Sprintf("[%v]", k), origin: value, current: value, parent: field})
			} else if tag.IsDig() {
				err = v.validate(Field{name: fmt.Sprintf("[%v]", k), origin: value, current: v.extractVar(value), parent: field}, tag)
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
				err = v.validateStruct(Field{name: fmt.Sprintf("[%d]", i), origin: value, current: value, parent: field})
			} else if tag.IsDig() {
				err = v.validate(Field{name: fmt.Sprintf("[%d]", i), origin: value, current: v.extractVar(value), parent: field}, tag)
			}

			if err != nil {
				if es, ok := err.(Errors); ok {
					errs = append(errs, es...)
				} else {
					return err
				}
			}
		}

	case reflect.Interface, reflect.Ptr:
		if val.IsNil() {
			break
		}
		value := val.Elem()

		var err error
		if value.Kind() == reflect.Struct || (value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct) {
			err = v.validateStruct(Field{origin: field.origin, current: value, parent: field})
		} else if tag.IsDig() {
			err = v.validate(Field{origin: field.origin, current: value, parent: field}, tag)
		}

		if err != nil {
			if es, ok := err.(Errors); ok {
				errs = append(errs, es...)
			} else {
				return err
			}
		}

	case reflect.Struct:
		err := v.validateStruct(Field{origin: field.origin, current: val, parent: field})
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

func (v *Validator) extractVar(in reflect.Value) reflect.Value {
	val := in
	for {
		switch val.Kind() {
		case reflect.Ptr, reflect.Interface:
			if val.IsNil() {
				return val
			}
			val = val.Elem()

		default:
			return val
		}
	}
}
