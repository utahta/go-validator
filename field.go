package validator

import (
	"fmt"
	"reflect"
)

type (
	// ParentField represents a parent of Field.
	ParentField struct {
		origin reflect.Value
	}

	// Field represents a value.
	Field struct {
		name     string
		fullName string
		origin   reflect.Value
		current  reflect.Value
		parent   ParentField
	}
)

const (
	fieldNameDelim = "."
)

func newFieldWithParent(name string, origin, current reflect.Value, parent Field) Field {
	if name == "" {
		name = parent.name
	} else if parent.name != "" {
		if name[0] == '[' {
			name = parent.name + name
		} else {
			name = parent.name + fieldNameDelim + name
		}
	}

	return Field{
		name:    name,
		origin:  origin,
		current: current,
		parent:  ParentField{origin: parent.origin},
	}
}

// Name is a field name. e.g Foo.Bar.Value
func (f Field) Name() string {
	return f.name
}

// Interface returns an interface{}
func (f Field) Interface() interface{} {
	return f.origin.Interface()
}

// Value returns a current field value.
func (f Field) Value() reflect.Value {
	return f.current
}

// Parent returns a parent field.
func (f Field) Parent() ParentField {
	return f.parent
}

// ShortString returns a string with 32 characters or more omitted.
func (f Field) ShortString() string {
	const maxSize = 32
	s := f.String()
	if len(s) > maxSize {
		return s[:maxSize] + "..."
	}
	return s
}

// String returns a string.
func (f Field) String() string {
	val := f.current
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

// Interface returns an interface{}
func (f ParentField) Interface() interface{} {
	if f.origin.IsValid() {
		return f.origin.Interface()
	}
	return nil
}
