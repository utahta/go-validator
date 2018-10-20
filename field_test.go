package validator

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestField_Value(t *testing.T) {
	value := reflect.ValueOf("value")
	field := Field{current: value}
	if value != field.Value() {
		t.Fatalf("invalid field value")
	}
}

func TestField_Parent(t *testing.T) {
	value := reflect.ValueOf("value")
	parent := ParentField{value}
	field := Field{parent: parent}

	if parent != field.Parent() {
		t.Fatal("invalid parent of field")
	}
	if parent.Interface() != value.Interface() {
		t.Fatal("invalid parent field")
	}

	field = Field{parent: ParentField{reflect.ValueOf(nil)}}
	if field.Parent().Interface() != nil {
		t.Fatalf("want parent field nil, but got %v", field.Parent().Interface())
	}
}

func TestField_ShortString(t *testing.T) {
	testcases := []struct {
		field      Field
		wantString string
	}{
		{
			field:      Field{current: reflect.ValueOf(strings.Repeat("a", 31))},
			wantString: fmt.Sprintf("%s", strings.Repeat("a", 31)),
		},
		{
			field:      Field{current: reflect.ValueOf(strings.Repeat("a", 32))},
			wantString: fmt.Sprintf("%s", strings.Repeat("a", 32)),
		},
		{
			field:      Field{current: reflect.ValueOf(strings.Repeat("a", 33))},
			wantString: fmt.Sprintf("%s...", strings.Repeat("a", 32)),
		},
		{
			field:      Field{current: reflect.ValueOf(strings.Repeat("a", 64))},
			wantString: fmt.Sprintf("%s...", strings.Repeat("a", 32)),
		},
	}

	for _, tc := range testcases {
		if tc.field.ShortString() != tc.wantString {
			t.Errorf("want `%s`, but got `%s`", tc.wantString, tc.field.ShortString())
		}
	}
}

func TestField_String(t *testing.T) {
	var unknown interface{}

	testcases := []struct {
		field      Field
		wantString string
	}{
		{
			field:      Field{current: reflect.ValueOf([]interface{}{"a"}).Index(0)},
			wantString: "<Interface>",
		},
		{
			field:      Field{current: reflect.ValueOf(&bytes.Buffer{})},
			wantString: "<Ptr>",
		},
		{
			field:      Field{current: reflect.ValueOf(unknown)},
			wantString: "<Unknown>",
		},
	}

	for _, tc := range testcases {
		if tc.field.String() != tc.wantString {
			t.Errorf("want %s, but got %s", tc.wantString, tc.field.String())
		}
	}
}
