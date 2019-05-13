package validator

import (
	"reflect"
	"testing"
)

func TestFieldError_Error(t *testing.T) {
	testcase := []struct {
		name        string
		err         Error
		wantMessage string
	}{
		{
			name: "error",
			err: &fieldError{
				field: Field{name: "field", current: reflect.ValueOf("text")},
				tag:   Tag{Name: "tag"},
			},
			wantMessage: "field: 'text' does validate as 'tag'",
		},
		{
			name: "error suppress field value",
			err: &fieldError{
				field:                   Field{name: "field", current: reflect.ValueOf("text")},
				tag:                     Tag{Name: "tag"},
				suppressErrorFieldValue: true,
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

func TestFieldError_Field(t *testing.T) {
	err := &fieldError{field: Field{name: "field", current: reflect.ValueOf("text")}}
	if want, got := "field", err.Field().Name(); want != got {
		t.Errorf("want %v, but got %v", want, got)
	}
	if want, got := "text", err.Field().String(); want != got {
		t.Errorf("want %v, but got %v", want, got)
	}
}

func TestFieldError_Tag(t *testing.T) {
	err := &fieldError{tag: Tag{Name: "tmp", Params: []string{"1", "2", "3"}}}
	if want, got := "tmp", err.Tag().Name; want != got {
		t.Errorf("want %v, but got %v", want, got)
	}
	if want, got := "tmp(1|2|3)", err.Tag().String(); want != got {
		t.Errorf("want %v, but got %v", want, got)
	}
}
