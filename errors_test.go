package validator

import (
	"reflect"
	"testing"
)

func TestError_Error(t *testing.T) {
	testcase := []struct {
		name        string
		err         Error
		wantMessage string
	}{
		{
			name: "error",
			err: Error{
				Field: Field{name: "field", current: reflect.ValueOf("text")},
				Tag:   Tag{Name: "tag"},
			},
			wantMessage: "field: 'text' does validate as 'tag'",
		},
		{
			name: "error suppress field value",
			err: Error{
				Field: Field{name: "field", current: reflect.ValueOf("text")},
				Tag:   Tag{Name: "tag"},
				SuppressErrorFieldValue: true,
			},
			wantMessage: "field: The value does validate as 'tag'",
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err.Error() != tc.wantMessage {
				t.Fatalf("error want `%v`, got `%v`", tc.wantMessage, tc.err.Error())
			}
		})
	}
}
