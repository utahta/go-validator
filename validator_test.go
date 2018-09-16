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

	OptionalTest struct {
		Value string `valid:"optional,alpha"`
	}

	SimpleStructTest struct {
		Str      Str      `valid:"required ;"`
		StrNoTag StrNoTag `valid:"required"`
	}

	DigArrayTest struct {
		Strs []string `valid:"len(3) ; required,alpha"`
		Ints []int    `valid:"required ;"`
		Nums []string `valid:"; numeric"`
	}

	DigMapTest struct {
		Strm map[string]string `valid:"len(3) ; required,alpha"`
		Intm map[int]int       `valid:"required ;"`
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

	ArrayTest struct {
	}

	MapTest struct {
	}
)

// for interface test
func (s Str) String() string {
	return s.Value
}

func TestValidator_ValidateStruct(t *testing.T) {
	testcases := []struct {
		name            string
		s               interface{}
		expectedNoErr   bool
		expectedMessage string
	}{
		// SimpleTest
		{
			name: "valid SimpleTest",
			s: SimpleTest{
				Value: "simple_test",
				Str:   Str{Value: "str_value"},
				Int:   Int{Value: 1},
			},
			expectedNoErr: true,
		},
		{
			name: "invalid SimpleTest.Value_Str_Int",
			s: SimpleTest{
				Value: "",
				Str:   Str{Value: ""},
				Int:   Int{Value: 0},
			},
			expectedMessage: "Value: '' does validate as 'required';Str.Value: '' does validate as 'required';Int.Value: '0' does validate as 'required'",
		},

		// OptionalTest
		{
			name: "valid OptionalTest",
			s: OptionalTest{
				Value: "abc",
			},
			expectedNoErr: true,
		},
		{
			name: "valid OptionalTest empty",
			s: OptionalTest{
				Value: "",
			},
			expectedNoErr: true,
		},
		{
			name: "invalid OptionalTest",
			s: OptionalTest{
				Value: "123",
			},
			expectedMessage: "Value: '123' does validate as 'alpha'",
		},

		// SimpleStructTest
		{
			name: "valid SimpleStructTest",
			s: SimpleStructTest{
				Str:      Str{Value: "str"},
				StrNoTag: StrNoTag{Value: "str"},
			},
			expectedNoErr: true,
		},
		{
			name: "invalid SimpleStructTest",
			s: SimpleStructTest{
				Str:      Str{},
				StrNoTag: StrNoTag{Value: ""},
			},
			expectedMessage: "Str: 'Str' does validate as 'required';Str.Value: '' does validate as 'required';StrNoTag: 'StrNoTag' does validate as 'required'",
		},

		// DigArrayTest
		{
			name: "valid DigArrayTest",
			s: DigArrayTest{
				Strs: []string{"a", "b", "c"},
				Ints: []int{0},
				Nums: []string{"0"},
			},
			expectedNoErr: true,
		},
		{
			name: "invalid DigArrayTest.Strs_Ints_Nums",
			s: DigArrayTest{
				Strs: []string{},
				Ints: []int{},
				Nums: []string{},
			},
			expectedMessage: "Strs: '<Array>' does validate as 'len(3)';Ints: '<Array>' does validate as 'required'",
		},
		{
			name: "invalid DigArrayTest.Strs.[1]_Nums[1]",
			s: DigArrayTest{
				Strs: []string{"a", "", "c"},
				Ints: []int{1, 0, 2},
				Nums: []string{"0", "a", "2"},
			},
			expectedMessage: "Strs.[1]: '' does validate as 'required';Strs.[1]: '' does validate as 'alpha';Nums.[1]: 'a' does validate as 'numeric'",
		},

		// DigMapTest
		{
			name: "valid DigMapTest",
			s: DigMapTest{
				Strm: map[string]string{"key1": "a", "key2": "b", "key3": "c"},
				Intm: map[int]int{0: 0},
			},
			expectedNoErr: true,
		},
		{
			name: "invalid DigMapTest.Strm_Intm",
			s: DigMapTest{
				Strm: map[string]string{},
				Intm: map[int]int{},
			},
			expectedMessage: "Strm: '<Map>' does validate as 'len(3)';Intm: '<Map>' does validate as 'required'",
		},
		{
			name: "invalid DigMapTest.Strm.[1]",
			s: DigMapTest{
				Strm: map[string]string{"key1": "a", "key2": "", "key3": "c"},
				Intm: map[int]int{0: 1, 1: 0, 2: 2},
			},
			expectedMessage: "Strm.[key2]: '' does validate as 'required';Strm.[key2]: '' does validate as 'alpha'",
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
			expectedNoErr: true,
		},
		{
			name: "invalid InterfaceTest_nil",
			s: InterfaceTest{
				IF:  nil,
				IFs: nil,
				IFm: nil,
			},
			expectedMessage: "IF: '<nil>' does validate as 'required';IFs: '<Array>' does validate as 'required';IFm: '<Map>' does validate as 'required'",
		},
		{
			name: "invalid InterfaceTest_empty",
			s: InterfaceTest{
				IF:  Str{""},
				IFs: []fmt.Stringer{},
				IFm: map[string]fmt.Stringer{},
			},
			expectedMessage: "IF: 'Str' does validate as 'required';IF.Value: '' does validate as 'required';IFs: '<Array>' does validate as 'required';IFm: '<Map>' does validate as 'required'",
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
			expectedMessage: "IF: 'Str' does validate as 'required';IF.Value: '' does validate as 'required';IFs.[0]: 'Str' does validate as 'required';IFs.[0].Value: '' does validate as 'required';IFm.[key1]: 'Str' does validate as 'required';IFm.[key1].Value: '' does validate as 'required'",
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
			expectedNoErr: true,
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
			expectedMessage: "Ptr: '<nil>' does validate as 'required';Ptrs.[0]: '<nil>' does validate as 'required';Ptrm.[key1]: '<nil>' does validate as 'required'",
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
			expectedMessage: "Ptr: 'Str' does validate as 'required';Ptr.Value: '' does validate as 'required';Ptrs.[0].Value: '' does validate as 'required';Ptrm.[key1].Value: '' does validate as 'required'",
		},
	}

	v := validator.New()
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := v.ValidateStruct(tc.s)

			if tc.expectedNoErr {
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

			if tc.expectedMessage != errs.Error() {
				t.Fatalf("expected `%v`\nbut got `%v`", tc.expectedMessage, errs.Error())
			}
		})
	}
}
