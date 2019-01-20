package validator

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

const (
	tagName = "valid"
)

var (
	defaultValidator     *Validator
	defaultValidatorOnce sync.Once
)

type (
	// Validator is a validator
	Validator struct {
		// FuncMap is a map of validate functions.
		FuncMap FuncMap

		// Adapters is a validate function adapters.
		Adapters []Adapter

		// SuppressErrorFieldValue is a flag that suppress output of field value in error.
		SuppressErrorFieldValue bool

		tagCache    *tagCache
		structCache *structCache
	}
)

// New returns a Validator
func New() *Validator {
	funcMap := FuncMap{}
	for k, fn := range defaultFuncMap {
		funcMap[k] = apply(fn, defaultAdapters...)
	}

	return &Validator{
		FuncMap:     funcMap,
		Adapters:    defaultAdapters,
		tagCache:    newTagCache(),
		structCache: newStructCache(),
	}
}

// SetFunc sets a validate function.
func (v *Validator) SetFunc(rawTag string, fn Func) {
	v.FuncMap[rawTag] = apply(fn, v.Adapters...)
}

// SetAdapters sets a validate function adapters.
func (v *Validator) SetAdapters(adapter ...Adapter) {
	v.Adapters = append(v.Adapters, adapter...)
	for k, fn := range v.FuncMap {
		v.FuncMap[k] = apply(fn, adapter...)
	}
}

// ValidateStruct validates a struct that use tags for fields.
func (v *Validator) ValidateStruct(s interface{}) error {
	return v.ValidateStructContext(context.Background(), s)
}

// ValidateStructContext validates a struct that use tags for fields.
// Pass context to each validate function.
func (v *Validator) ValidateStructContext(ctx context.Context, s interface{}) error {
	if s == nil {
		return nil
	}
	value := reflect.ValueOf(s)
	return v.validateStruct(ctx, Field{origin: value, current: value})
}

func (v *Validator) validateStruct(ctx context.Context, field Field) error {
	val := field.current
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

		if err := v.validateVar(ctx, newFieldWithParent(fieldCaches[i].name, originField, valueField, field), fieldCaches[i].tagValue); err != nil {
			if es, ok := err.(Errors); ok {
				errs = append(errs, es...)
			} else {
				return err
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

// ValidateVar validates a value.
func (v *Validator) ValidateVar(s interface{}, rawTag string) error {
	return v.ValidateVarContext(context.Background(), s, rawTag)
}

// ValidateVarContext validates a value.
// Pass context to each validate function.
func (v *Validator) ValidateVarContext(ctx context.Context, s interface{}, rawTag string) error {
	value := reflect.ValueOf(s)
	return v.validateVar(ctx, Field{origin: value, current: v.extractVar(value)}, rawTag)
}

func (v *Validator) validateVar(ctx context.Context, field Field, rawTag string) error {
	if rawTag == "-" {
		return nil
	}

	chunk, err := v.tagParse(rawTag)
	if err != nil {
		return err
	}

	var errs Errors
	if err := v.validate(ctx, field, chunk); err != nil {
		if es, ok := err.(Errors); ok {
			errs = append(errs, es...)
		} else {
			return err
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (v *Validator) validate(ctx context.Context, field Field, chunk *tagChunk) error {
	if chunk.Optional && isEmpty(field) {
		return nil
	}

	var errs Errors
	for _, tag := range chunk.Tags {
		valid, err := tag.validateFn(ctx, field, FuncOption{Params: tag.Params, v: v})
		if !valid || err != nil {
			errs = append(errs, Error{Field: field, Tag: tag, Err: err, SuppressErrorFieldValue: v.SuppressErrorFieldValue})
		}
	}

	var val = field.current
	switch val.Kind() {
	case reflect.Map:
		for _, k := range val.MapKeys() {
			value := val.MapIndex(k)

			err := v.validate(ctx, newFieldWithParent(fmt.Sprintf("[%v]", k), value, v.extractVar(value), field), chunk.Next)
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

			err := v.validate(ctx, newFieldWithParent(fmt.Sprintf("[%d]", i), value, v.extractVar(value), field), chunk.Next)
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
		// unreachable

	case reflect.Struct:
		err := v.validateStruct(ctx, newFieldWithParent("", field.origin, val, field))
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

// DefaultValidator returns default validator.
func DefaultValidator() *Validator {
	defaultValidatorOnce.Do(func() {
		defaultValidator = New()
	})
	return defaultValidator
}

// ValidateStruct validates a struct that use tags for fields using default validator.
func ValidateStruct(s interface{}) error {
	return DefaultValidator().ValidateStruct(s)
}

// ValidateStructContext validates a struct that use tags for fields using default validator.
// Pass context to each validate function.
func ValidateStructContext(ctx context.Context, s interface{}) error {
	return DefaultValidator().ValidateStructContext(ctx, s)
}

// ValidateVar validates a value using default validator.
func ValidateVar(s interface{}, rawTag string) error {
	return DefaultValidator().ValidateVar(s, rawTag)
}

// ValidateVarContext validates a value using default validator.
// Pass context to each validate function.
func ValidateVarContext(ctx context.Context, s interface{}, rawTag string) error {
	return DefaultValidator().ValidateVarContext(ctx, s, rawTag)
}
