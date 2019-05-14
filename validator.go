package validator

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

var (
	defaultValidator     *Validator
	defaultValidatorOnce sync.Once
)

type (
	// Validator is a validator that validates each fields using struct field's tag.
	Validator struct {
		// funcMap represents a map of validating functions.
		funcMap FuncMap

		// adapters represents a slice of validating function adapter.
		adapters []Adapter

		// tagKey is the key in the struct field's tag. the default value is `valid`.
		tagKey string

		// suppressErrorFieldValue is a flag that suppresses field value by error.
		suppressErrorFieldValue bool

		tagCache    *tagCache
		structCache *structCache
	}

	Option func(v *Validator)
)

// New returns a Validator
func New(opts ...Option) *Validator {
	funcMap := FuncMap{}
	for k, fn := range defaultFuncMap {
		funcMap[k] = apply(fn, defaultAdapters...)
	}

	v := &Validator{
		funcMap:     funcMap,
		adapters:    defaultAdapters,
		tagKey:      "valid",
		tagCache:    newTagCache(),
		structCache: newStructCache(),
	}
	v.Apply(opts...)
	return v
}

// WithFunc is a validator option that sets a validating function.
func WithFunc(k string, fn Func) Option {
	return func(v *Validator) {
		v.funcMap[k] = apply(fn, v.adapters...)
	}
}

// WithFuncMap is a validator option that sets validating functions.
func WithFuncMap(funcMap FuncMap) Option {
	return func(v *Validator) {
		for k, fn := range funcMap {
			v.funcMap[k] = apply(fn, v.adapters...)
		}
	}
}

// WithAdapters is a validator option that sets validator function adapters.
func WithAdapters(adapters ...Adapter) Option {
	return func(v *Validator) {
		v.adapters = append(v.adapters, adapters...)
		for k, fn := range v.funcMap {
			v.funcMap[k] = apply(fn, adapters...)
		}
	}
}

// WithTagKey is a validator option that sets the key in the struct field's tag.
func WithTagKey(k string) Option {
	return func(v *Validator) {
		v.tagKey = k
	}
}

// WithSuppressErrorFieldValue is a validator option that enables suppress validating field value by error.
// If enabled this option, the field value always replaces `The value`.
func WithSuppressErrorFieldValue() Option {
	return func(v *Validator) {
		v.suppressErrorFieldValue = true
	}
}

// Apply applies validator options.
func (v *Validator) Apply(opts ...Option) {
	for _, o := range opts {
		o(v)
	}
}

// ValidateStruct validates a struct that uses the struct field's tag.
func (v *Validator) ValidateStruct(s interface{}) error {
	return v.ValidateStructContext(context.Background(), s)
}

// ValidateStructContext validates a struct that uses the struct field's tag.
// Pass context to each validating functions.
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
	if !hasCache {
		for i := 0; i < val.NumField(); i++ {
			typeField := valueType.Field(i)
			cache := fieldCache{
				index:     i,
				isPrivate: typeField.PkgPath != "", // private field
				tagValue:  typeField.Tag.Get(v.tagKey),
				name:      typeField.Name,
			}
			if cache.isPrivate {
				continue
			}
			if !v.canValidate(cache.tagValue, v.extractVar(val.Field(i)).Kind()) {
				continue
			}

			chunk, err := v.parseTag(cache.tagValue)
			if err != nil {
				return err
			}
			cache.tagChunk = chunk

			fieldCaches = append(fieldCaches, cache)
		}
		v.structCache.Store(valueType, fieldCaches)
	}

	var errs Errors
	for i := 0; i < len(fieldCaches); i++ {
		originField := val.Field(fieldCaches[i].index)
		valueField := v.extractVar(originField)

		if err := v.validate(ctx, newFieldWithParent(fieldCaches[i].name, originField, valueField, field), fieldCaches[i].tagChunk); err != nil {
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

// ValidateVar validates a value.
func (v *Validator) ValidateVar(s interface{}, rawTag string) error {
	return v.ValidateVarContext(context.Background(), s, rawTag)
}

// ValidateVarContext validates a value.
// Pass context to each validating functions.
func (v *Validator) ValidateVarContext(ctx context.Context, s interface{}, rawTag string) error {
	value := reflect.ValueOf(s)
	return v.validateVar(ctx, Field{origin: value, current: v.extractVar(value)}, rawTag)
}

func (v *Validator) validateVar(ctx context.Context, field Field, rawTag string) error {
	if !v.canValidate(rawTag, field.current.Kind()) {
		return nil
	}

	chunk, err := v.parseTag(rawTag)
	if err != nil {
		return err
	}

	return v.validate(ctx, field, chunk)
}

func (v *Validator) validate(ctx context.Context, field Field, chunk *tagChunk) error {
	if chunk.IsOptional() && isEmpty(field) {
		return nil
	}

	var errs Errors
	for _, tag := range chunk.GetTags() {
		valid, err := tag.validateFn(ctx, field, FuncOption{Params: tag.params, v: v})
		if !valid || err != nil {
			errs = append(errs, &fieldError{field: field, tag: tag, err: err, suppressErrorFieldValue: v.suppressErrorFieldValue})
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
		// do nothing

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

func (v *Validator) canValidate(rawTag string, kind reflect.Kind) bool {
	if rawTag == "-" {
		return false
	}

	if rawTag == "" {
		switch kind {
		case reflect.String, reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64:
			// these kinds do not perform recursive process so let's skip validation
			return false
		}
	}
	return true
}

// DefaultValidator returns a default validator.
func DefaultValidator() *Validator {
	defaultValidatorOnce.Do(func() {
		defaultValidator = New()
	})
	return defaultValidator
}

// ValidateStruct validates a struct that uses the struct field's tag using default validator.
func ValidateStruct(s interface{}) error {
	return DefaultValidator().ValidateStruct(s)
}

// ValidateStructContext validates a struct that uses the struct field's tag using default validator.
// Pass context to each validating functions.
func ValidateStructContext(ctx context.Context, s interface{}) error {
	return DefaultValidator().ValidateStructContext(ctx, s)
}

// ValidateVar validates a value using default validator.
func ValidateVar(s interface{}, rawTag string) error {
	return DefaultValidator().ValidateVar(s, rawTag)
}

// ValidateVarContext validates a value using default validator.
// Pass context to each validating functions.
func ValidateVarContext(ctx context.Context, s interface{}, rawTag string) error {
	return DefaultValidator().ValidateVarContext(ctx, s, rawTag)
}
