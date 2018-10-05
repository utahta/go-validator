package validator_test

import (
	"fmt"
	"testing"

	"github.com/utahta/go-validator"
)

type (
	Str struct {
		Value string `valid:"required"`
	}

	Int struct {
		Value int `valid:"required"`
	}

	StrNoTag struct {
		Value string
	}

	SimpleTest struct {
		Value string `valid:"required"`
		Str   Str
		Int   Int
	}

	ArrayStringTest struct {
		Values []string `valid:"len(3); req,alpha"`
	}

	RequiredArrayStringTest struct {
		Values []string `valid:"required ;"`
	}

	ArrayRequiredStringTest struct {
		Values []string `valid:"; required"`
	}

	MapStringTest struct {
		Values map[string]string `valid:"len(3); req,alpha"`
	}

	RequiredMapStringTest struct {
		Values map[string]string `valid:"required ;"`
	}

	MapRequiredStringTest struct {
		Values map[string]string `valid:"; required"`
	}

	OptionalStringTest struct {
		Value string `valid:"optional,alpha"`
	}

	OptionalArrayStringTest struct {
		Values []string `valid:"optional; alpha"`
	}

	OptionalArrayOptionalStringTest struct {
		Values []string `valid:"optional; optional,alpha"`
	}

	OptionalMapStringTest struct {
		Values map[string]string `valid:"optional; alpha"`
	}

	OptionalMapOptionalStringTest struct {
		Values map[string]string `valid:"optional; optional,alpha"`
	}

	SimpleStructTest struct {
		Str      Str      `valid:"required ;"`
		StrNoTag StrNoTag `valid:"required"`
	}

	InterfaceTest struct {
		IF  fmt.Stringer            `valid:"required"`
		IFs []fmt.Stringer          `valid:"required"`
		IFm map[string]fmt.Stringer `valid:"required"`
	}

	PtrTest struct {
		Ptr  *Str            `valid:"required"`
		Ptrs []*Str          `valid:"required"`
		Ptrm map[string]*Str `valid:"required"`
	}
)

// for interface test
func (s Str) String() string {
	return s.Value
}

func TestValidator_ValidateStruct(t *testing.T) {
	testcases := []struct {
		name        string
		s           interface{}
		wantNoErr   bool
		wantMessage string
	}{
		// SimpleTest
		{
			name: "valid SimpleTest",
			s: SimpleTest{
				Value: "simple_test",
				Str:   Str{Value: "str_value"},
				Int:   Int{Value: 1},
			},
			wantNoErr: true,
		},
		{
			name: "invalid SimpleTest.Value_Str_Int",
			s: SimpleTest{
				Value: "",
				Str:   Str{Value: ""},
				Int:   Int{Value: 0},
			},
			wantMessage: "Value: '' does validate as 'required';Str.Value: '' does validate as 'required';Int.Value: '0' does validate as 'required'",
		},

		// ArrayStringTest
		{
			name: "valid ArrayStringTest",
			s: ArrayStringTest{
				Values: []string{"a", "b", "c"},
			},
			wantNoErr: true,
		},
		{
			name: "invalid ArrayStringTest length",
			s: ArrayStringTest{
				Values: []string{"a", "b"},
			},
			wantMessage: "Values: '<Array>' does validate as 'len(3)'",
		},
		{
			name: "invalid ArrayStringTest.Values[0]",
			s: ArrayStringTest{
				Values: []string{"", "b", "c"},
			},
			wantMessage: "Values[0]: '' does validate as 'req';Values[0]: '' does validate as 'alpha'",
		},

		// RequiredArrayStringTest
		{
			name: "valid RequiredArrayStringTest",
			s: RequiredArrayStringTest{
				Values: []string{""},
			},
			wantNoErr: true,
		},
		{
			name: "invalid RequiredArrayStringTest empty",
			s: RequiredArrayStringTest{
				Values: []string{},
			},
			wantMessage: "Values: '<Array>' does validate as 'required'",
		},
		{
			name: "invalid RequiredArrayStringTest nil",
			s: RequiredArrayStringTest{
				Values: nil,
			},
			wantMessage: "Values: '<Array>' does validate as 'required'",
		},

		// ArrayRequiredStringTest
		{
			name: "valid ArrayRequiredStringTest",
			s: ArrayRequiredStringTest{
				Values: []string{},
			},
			wantNoErr: true,
		},
		{
			name: "valid ArrayRequiredStringTest",
			s: ArrayRequiredStringTest{
				Values: nil,
			},
			wantNoErr: true,
		},
		{
			name: "invalid ArrayRequiredStringTest",
			s: ArrayRequiredStringTest{
				Values: []string{""},
			},
			wantMessage: "Values[0]: '' does validate as 'required'",
		},

		// MapStringTest
		{
			name: "valid MapStringTest",
			s: MapStringTest{
				Values: map[string]string{"key1": "a", "key2": "b", "key3": "c"},
			},
			wantNoErr: true,
		},
		{
			name: "invalid MapStringTest length",
			s: MapStringTest{
				Values: map[string]string{"key1": "a", "key2": "b"},
			},
			wantMessage: "Values: '<Map>' does validate as 'len(3)'",
		},
		{
			name: "invalid MapStringTest.Values[0]",
			s: MapStringTest{
				Values: map[string]string{"key1": "", "key2": "b", "key3": "c"},
			},
			wantMessage: "Values[key1]: '' does validate as 'req';Values[key1]: '' does validate as 'alpha'",
		},

		// RequiredMapStringTest
		{
			name: "valid RequiredMapStringTest",
			s: RequiredMapStringTest{
				Values: map[string]string{"key1": ""},
			},
			wantNoErr: true,
		},
		{
			name: "invalid RequiredMapStringTest empty",
			s: RequiredMapStringTest{
				Values: map[string]string{},
			},
			wantMessage: "Values: '<Map>' does validate as 'required'",
		},
		{
			name: "invalid RequiredMapStringTest nil",
			s: RequiredMapStringTest{
				Values: nil,
			},
			wantMessage: "Values: '<Map>' does validate as 'required'",
		},

		// MapRequiredStringTest
		{
			name: "valid MapRequiredStringTest",
			s: MapRequiredStringTest{
				Values: map[string]string{},
			},
			wantNoErr: true,
		},
		{
			name: "valid MapRequiredStringTest",
			s: MapRequiredStringTest{
				Values: nil,
			},
			wantNoErr: true,
		},
		{
			name: "invalid MapRequiredStringTest",
			s: MapRequiredStringTest{
				Values: map[string]string{"key1": ""},
			},
			wantMessage: "Values[key1]: '' does validate as 'required'",
		},

		// OptionalTest
		{
			name: "valid OptionalStringTest",
			s: OptionalStringTest{
				Value: "abc",
			},
			wantNoErr: true,
		},
		{
			name: "valid OptionalStringTest empty",
			s: OptionalStringTest{
				Value: "",
			},
			wantNoErr: true,
		},
		{
			name: "invalid OptionalStringTest",
			s: OptionalStringTest{
				Value: "123",
			},
			wantMessage: "Value: '123' does validate as 'alpha'",
		},

		// OptionalArrayStringTest
		{
			name: "valid OptionalArrayStringTest",
			s: OptionalArrayStringTest{
				Values: []string{},
			},
			wantNoErr: true,
		},
		{
			name: "valid OptionalArrayStringTest",
			s: OptionalArrayStringTest{
				Values: []string{"abc"},
			},
			wantNoErr: true,
		},
		{
			name: "invalid OptionalArrayStringTest.Values[0]",
			s: OptionalArrayStringTest{
				Values: []string{""},
			},
			wantMessage: "Values[0]: '' does validate as 'alpha'",
		},

		// OptionalArrayOptionalStringTest
		{
			name: "valid OptionalArrayOptionalStringTest",
			s: OptionalArrayOptionalStringTest{
				Values: []string{},
			},
			wantNoErr: true,
		},
		{
			name: "valid OptionalArrayOptionalStringTest",
			s: OptionalArrayOptionalStringTest{
				Values: []string{""},
			},
			wantNoErr: true,
		},
		{
			name: "invalid OptionalArrayOptionalStringTest.Values[0]",
			s: OptionalArrayOptionalStringTest{
				Values: []string{"123"},
			},
			wantMessage: "Values[0]: '123' does validate as 'alpha'",
		},

		// OptionalMapStringTest
		{
			name: "valid OptionalMapStringTest",
			s: OptionalMapStringTest{
				Values: map[string]string{},
			},
			wantNoErr: true,
		},
		{
			name: "valid OptionalMapStringTest",
			s: OptionalMapStringTest{
				Values: map[string]string{
					"key1": "abc",
				},
			},
			wantNoErr: true,
		},
		{
			name: "invalid OptionalMapStringTest.Values[key1]",
			s: OptionalMapStringTest{
				Values: map[string]string{
					"key1": "",
				},
			},
			wantMessage: "Values[key1]: '' does validate as 'alpha'",
		},

		// OptionalMapOptionalStringTest
		{
			name: "valid OptionalMapOptionalStringTest",
			s: OptionalMapOptionalStringTest{
				Values: map[string]string{},
			},
			wantNoErr: true,
		},
		{
			name: "valid OptionalMapOptionalStringTest",
			s: OptionalMapOptionalStringTest{
				Values: map[string]string{
					"key1": "",
				},
			},
			wantNoErr: true,
		},
		{
			name: "invalid OptionalMapOptionalStringTest",
			s: OptionalMapOptionalStringTest{
				Values: map[string]string{
					"key1": "123",
				},
			},
			wantMessage: "Values[key1]: '123' does validate as 'alpha'",
		},

		// SimpleStructTest
		{
			name: "valid SimpleStructTest",
			s: SimpleStructTest{
				Str:      Str{Value: "str"},
				StrNoTag: StrNoTag{Value: "str"},
			},
			wantNoErr: true,
		},
		{
			name: "invalid SimpleStructTest",
			s: SimpleStructTest{
				Str:      Str{},
				StrNoTag: StrNoTag{Value: ""},
			},
			wantMessage: "Str: 'Str' does validate as 'required';Str.Value: '' does validate as 'required';StrNoTag: 'StrNoTag' does validate as 'required'",
		},

		// InterfaceTest
		{
			name: "valid InterfaceTest",
			s: InterfaceTest{
				IF:  Str{"a"},
				IFs: []fmt.Stringer{Str{"a"}},
				IFm: map[string]fmt.Stringer{
					"key1": Str{"a"},
				},
			},
			wantNoErr: true,
		},
		{
			name: "invalid InterfaceTest_nil",
			s: InterfaceTest{
				IF:  nil,
				IFs: nil,
				IFm: nil,
			},
			wantMessage: "IF: '<nil>' does validate as 'required';IFs: '<Array>' does validate as 'required';IFm: '<Map>' does validate as 'required'",
		},
		{
			name: "invalid InterfaceTest_empty",
			s: InterfaceTest{
				IF:  Str{""},
				IFs: []fmt.Stringer{},
				IFm: map[string]fmt.Stringer{},
			},
			wantMessage: "IF: 'Str' does validate as 'required';IF.Value: '' does validate as 'required';IFs: '<Array>' does validate as 'required';IFm: '<Map>' does validate as 'required'",
		},
		{
			name: "invalid InterfaceTest_ptr_empty",
			s: InterfaceTest{
				IF:  &Str{""},
				IFs: []fmt.Stringer{&Str{""}},
				IFm: map[string]fmt.Stringer{
					"key1": &Str{""},
				},
			},
			wantMessage: "IF: 'Str' does validate as 'required';IF.Value: '' does validate as 'required';IFs[0]: 'Str' does validate as 'required';IFs[0].Value: '' does validate as 'required';IFm[key1]: 'Str' does validate as 'required';IFm[key1].Value: '' does validate as 'required'",
		},

		// PtrTest
		{
			name: "valid PtrTest",
			s: PtrTest{
				Ptr:  &Str{"a"},
				Ptrs: []*Str{{"a"}},
				Ptrm: map[string]*Str{
					"key1": {"a"},
				},
			},
			wantNoErr: true,
		},
		{
			name: "invalid PtrTest_nil",
			s: PtrTest{
				Ptr:  nil,
				Ptrs: []*Str{nil},
				Ptrm: map[string]*Str{
					"key1": nil,
				},
			},
			wantMessage: "Ptr: '<nil>' does validate as 'required';Ptrs[0]: '<nil>' does validate as 'required';Ptrm[key1]: '<nil>' does validate as 'required'",
		},
		{
			name: "invalid PtrTest_empty",
			s: PtrTest{
				Ptr:  &Str{""},
				Ptrs: []*Str{{""}},
				Ptrm: map[string]*Str{
					"key1": {""},
				},
			},
			wantMessage: "Ptr: 'Str' does validate as 'required';Ptr.Value: '' does validate as 'required';Ptrs[0].Value: '' does validate as 'required';Ptrm[key1].Value: '' does validate as 'required'",
		},
	}

	v := validator.New()
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := v.ValidateStruct(tc.s)

			if tc.wantNoErr {
				if err != nil {
					t.Error(err)
				}
				return
			}

			if err == nil {
				t.Fatal("expected `error`, but got `nil`")
			}

			errs, ok := validator.ToErrors(err)
			if !ok {
				t.Fatal("expected `true`, but got `false`")
			}

			if tc.wantMessage != errs.Error() {
				t.Fatalf("expected `%v`\nbut got `%v`", tc.wantMessage, errs.Error())
			}
		})
	}
}
