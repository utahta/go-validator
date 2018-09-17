package validator

import (
	"fmt"
	"reflect"
	"strings"
)

type (
	ParentField interface {
		Name() string

		FullName() string

		Interface() interface{}

		Parent() ParentField
	}

	Field struct {
		name    string
		origin  reflect.Value
		current reflect.Value
		parent  ParentField
	}
)

const (
	fieldNameDelim = "."
)

func (f Field) Name() string {
	return f.name
}

func (f Field) FullName() string {
	var s []string
	if f.parent != nil && f.parent.FullName() != "" {
		s = append(s, f.parent.FullName())
	}

	if f.name != "" {
		s = append(s, f.name)
	}
	return strings.Join(s, fieldNameDelim)
}

func (f Field) Interface() interface{} {
	return f.origin.Interface()
}

func (f Field) Value() reflect.Value {
	return f.current
}

func (f Field) Parent() ParentField {
	return f.parent
}

func (f Field) ShortString() string {
	const maxSize = 32
	s := f.String()
	if len(s) > maxSize {
		return s[:maxSize] + "..."
	}
	return s
}

func (f Field) String() string {
	val := f.Value()
	switch val.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return fmt.Sprint(val)

	case reflect.String:
		return val.String()

	case reflect.Struct:
		return val.Type().Name()

	case reflect.Map:
		return "<Map>"

	case reflect.Slice, reflect.Array:
		return "<Array>"

	case reflect.Interface:
		if val.IsNil() {
			return "<nil>"
		}
		return "<Interface>"

	case reflect.Ptr:
		if val.IsNil() {
			return "<nil>"
		}
		return "<Ptr>"
	}
	return "<Unknown>"
}
