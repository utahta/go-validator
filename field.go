package validator

import (
	"fmt"
	"reflect"
	"strings"
)

type (
	Field struct {
		name   string
		val    reflect.Value
		parent *Field
	}
)

const (
	fieldNameDelim = "."
)

func (f Field) Name() string {
	var s []string
	if f.parent != nil {
		if f.parent.Name() != "" {
			s = append(s, f.parent.Name())
		}
	}

	if f.name != "" {
		s = append(s, f.name)
	}
	return strings.Join(s, fieldNameDelim)
}

func (f Field) ShortValue() string {
	const maxSize = 10
	s := f.Value()
	if len(s) > maxSize {
		return s[:maxSize] + "..."
	}
	return s
}

func (f Field) Value() string {
	switch f.val.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return fmt.Sprint(f.val)

	case reflect.String:
		return f.val.String()

	case reflect.Map:
		return "Map"

	case reflect.Slice, reflect.Array:
		return "Array"

	case reflect.Interface:
		if f.val.IsNil() {
			return "<nil>"
		}
		return "Interface"

	case reflect.Ptr:
		if f.val.IsNil() {
			return "<nil>"
		}
		return "Ptr"
	}
	return ""
}
