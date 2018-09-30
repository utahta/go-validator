package validator

import (
	"fmt"
	"reflect"
)

type (
	ParentField struct {
		origin reflect.Value
	}

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
	var fullName string
	if parent.fullName == "" {
		fullName = name
	} else if name == "" {
		fullName = parent.fullName
	} else {
		if name[0] == '[' {
			fullName = parent.fullName + name
		} else {
			fullName = parent.fullName + fieldNameDelim + name
		}
	}

	return Field{
		name:     name,
		fullName: fullName,
		origin:   origin,
		current:  current,
		parent:   ParentField{origin: parent.origin},
	}
}

func (f Field) Name() string {
	return f.name
}

func (f Field) FullName() string {
	return f.fullName
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

func (f ParentField) Interface() interface{} {
	if f.origin.IsValid() {
		return f.origin.Interface()
	}
	return nil
}
