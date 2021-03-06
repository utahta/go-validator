package validator_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/utahta/go-validator"
)

func TestValidateStruct_Simple(t *testing.T) {
	type (
		Cat struct {
			Name string `valid:"required"`
			Age  int    `valid:"required"`
		}

		SimpleTest struct {
			Cat Cat
			Str string
		}

		SimpleInterfaceTest struct {
			IF interface{} `valid:"required"`
		}

		SimplePointerTest struct {
			P *Cat `valid:"required"`
		}
	)

	testcases := []struct {
		name        string
		s           interface{}
		wantNoErr   bool
		wantMessage string
	}{
		// Cat
		{
			name:      "Valid Cat",
			s:         Cat{Name: "neko", Age: 5},
			wantNoErr: true,
		},
		{
			name:      "Valid *Cat",
			s:         &Cat{Name: "neko", Age: 5},
			wantNoErr: true,
		},

		// SimpleTest
		{
			name: "Valid SimpleTest",
			s: SimpleTest{
				Cat: Cat{Name: "neko", Age: 5},
			},
			wantNoErr: true,
		},
		{
			name: "Invalid",
			s: SimpleTest{
				Cat: Cat{},
			},
			wantMessage: "Cat.Name: '' does validate as 'required';Cat.Age: '0' does validate as 'required'",
		},

		// SimpleInterfaceTest
		{
			name: "Valid SimpleInterfaceTest",
			s: SimpleInterfaceTest{
				IF: Cat{Name: "neko", Age: 5},
			},
			wantNoErr: true,
		},
		{
			name: "Invalid SimpleInterfaceTest empty",
			s: SimpleInterfaceTest{
				IF: Cat{},
			},
			wantMessage: "IF: 'Cat' does validate as 'required';IF.Name: '' does validate as 'required';IF.Age: '0' does validate as 'required'",
		},
		{
			name: "Invalid SimpleInterfaceTest nil",
			s: SimpleInterfaceTest{
				IF: nil,
			},
			wantMessage: "IF: '<nil>' does validate as 'required'",
		},

		// SimplePointerTest
		{
			name: "Valid SimplePointerTest",
			s: SimplePointerTest{
				P: &Cat{Name: "neko", Age: 5},
			},
			wantNoErr: true,
		},
		{
			name: "Invalid SimplePointerTest empty",
			s: SimplePointerTest{
				P: &Cat{},
			},
			wantMessage: "P: 'Cat' does validate as 'required';P.Name: '' does validate as 'required';P.Age: '0' does validate as 'required'",
		},
		{
			name: "Invalid SimplePointerTest nil",
			s: SimplePointerTest{
				P: nil,
			},
			wantMessage: "P: '<nil>' does validate as 'required'",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateStruct(tc.s)

			if tc.wantNoErr {
				if err != nil {
					t.Error(err)
				}
				return
			}
			assertValidationError(t, tc.wantMessage, err)
		})
	}
}

func TestValidateStruct_Array(t *testing.T) {
	type (
		Cat struct {
			Name string `valid:"required,alpha"`
			Age  int    `valid:"required"`
		}

		ArrayStringTest struct {
			S []string `valid:"len(3); req,alpha"`
		}

		RequiredArrayStringTest struct {
			S []string `valid:"required ;"`
		}

		ArrayRequiredStringTest struct {
			S []string `valid:"; required"`
		}

		ArrayCatTest struct {
			S []Cat `valid:"required"`
		}

		OptionalArrayCatTest struct {
			S []Cat `valid:"optional"`
		}

		OptionalArrayCatTest2 struct {
			S []Cat `valid:"optional ;"`
		}

		ArrayOptionalCatTest struct {
			S []Cat `valid:"required; optional"`
		}

		OptionalArrayOptionalCatTest struct {
			S []Cat `valid:"optional ; optional"`
		}

		ArrayInterfaceTest struct {
			S []interface{} `valid:"required"`
		}

		ArrayPointerTest struct {
			S []*Cat `valid:"required"`
		}
	)

	testcases := []struct {
		name        string
		s           interface{}
		wantNoErr   bool
		wantMessage string
	}{
		// ArrayStringTest
		{
			name: "Valid ArrayStringTest",
			s: ArrayStringTest{
				S: []string{"a", "b", "c"},
			},
			wantNoErr: true,
		},
		{
			name: "Invalid ArrayStringTest length",
			s: ArrayStringTest{
				S: []string{"a", "b"},
			},
			wantMessage: "S: '<Array>' does validate as 'len(3)'",
		},
		{
			name: "Invalid ArrayStringTest.S[0] empty",
			s: ArrayStringTest{
				S: []string{"", "b", "c"},
			},
			wantMessage: "S[0]: '' does validate as 'req';S[0]: '' does validate as 'alpha'",
		},

		// RequiredArrayStringTest
		{
			name: "Valid RequiredArrayStringTest",
			s: RequiredArrayStringTest{
				S: []string{""},
			},
			wantNoErr: true,
		},
		{
			name: "invalid RequiredArrayStringTest empty",
			s: RequiredArrayStringTest{
				S: []string{},
			},
			wantMessage: "S: '<Array>' does validate as 'required'",
		},
		{
			name: "invalid RequiredArrayStringTest nil",
			s: RequiredArrayStringTest{
				S: nil,
			},
			wantMessage: "S: '<Array>' does validate as 'required'",
		},

		// ArrayRequiredStringTest
		{
			name: "Valid ArrayRequiredStringTest empty",
			s: ArrayRequiredStringTest{
				S: []string{},
			},
			wantNoErr: true,
		},
		{
			name: "Valid ArrayRequiredStringTest nil",
			s: ArrayRequiredStringTest{
				S: nil,
			},
			wantNoErr: true,
		},
		{
			name: "Invalid ArrayRequiredStringTest",
			s: ArrayRequiredStringTest{
				S: []string{""},
			},
			wantMessage: "S[0]: '' does validate as 'required'",
		},

		// ArrayCatTest
		{
			name: "Valid ArrayCatTest",
			s: ArrayCatTest{
				S: []Cat{{Name: "neko", Age: 5}},
			},
			wantNoErr: true,
		},
		{
			name: "Invalid ArrayCatTest.S empty",
			s: ArrayCatTest{
				S: []Cat{},
			},
			wantMessage: "S: '<Array>' does validate as 'required'",
		},

		// OptionalArrayCatTest
		{
			name: "Valid OptionalArrayCatTest empty",
			s: OptionalArrayCatTest{
				S: []Cat{},
			},
			wantNoErr: true,
		},
		{
			name: "Valid OptionalArrayCatTest nil",
			s: OptionalArrayCatTest{
				S: nil,
			},
			wantNoErr: true,
		},
		{
			name: "Valid OptionalArrayCatTest",
			s: OptionalArrayCatTest{
				S: []Cat{{Name: "neko", Age: 5}},
			},
			wantNoErr: true,
		},
		{
			name: "Invalid OptionalArrayCatTest.S[1]",
			s: OptionalArrayCatTest{
				S: []Cat{{Name: "neko", Age: 5}, {}},
			},
			wantMessage: "S[1].Name: '' does validate as 'required';S[1].Name: '' does validate as 'alpha';S[1].Age: '0' does validate as 'required'",
		},

		// OptionalArrayCatTest2
		{
			name: "Valid OptionalArrayCatTest2 empty",
			s: OptionalArrayCatTest2{
				S: []Cat{},
			},
			wantNoErr: true,
		},
		{
			name: "Valid OptionalArrayCatTest2 nil",
			s: OptionalArrayCatTest2{
				S: nil,
			},
			wantNoErr: true,
		},
		{
			name: "Valid OptionalArrayCatTest2",
			s: OptionalArrayCatTest2{
				S: []Cat{{Name: "neko", Age: 5}},
			},
			wantNoErr: true,
		},
		{
			name: "Invalid OptionalArrayCatTest2.S[1]",
			s: OptionalArrayCatTest2{
				S: []Cat{{Name: "neko", Age: 5}, {}},
			},
			wantMessage: "S[1].Name: '' does validate as 'required';S[1].Name: '' does validate as 'alpha';S[1].Age: '0' does validate as 'required'",
		},

		// ArrayOptionalCatTest
		{
			name: "Valid ArrayOptionalCatTest",
			s: ArrayOptionalCatTest{
				S: []Cat{{}},
			},
			wantNoErr: true,
		},
		{
			name: "Invalid ArrayOptionalCatTest",
			s: ArrayOptionalCatTest{
				S: []Cat{},
			},
			wantMessage: "S: '<Array>' does validate as 'required'",
		},
		{
			name: "Invalid ArrayOptionalCatTest.S[1]",
			s: ArrayOptionalCatTest{
				S: []Cat{{}, {Name: "123", Age: 5}},
			},
			wantMessage: "S[1].Name: '123' does validate as 'alpha'",
		},

		// OptionalArrayOptionalCatTest
		{
			name: "Valid OptionalArrayOptionalCatTest empty",
			s: OptionalArrayOptionalCatTest{
				S: []Cat{},
			},
			wantNoErr: true,
		},
		{
			name: "Valid OptionalArrayOptionalCatTest nil",
			s: OptionalArrayOptionalCatTest{
				S: nil,
			},
			wantNoErr: true,
		},
		{
			name: "Valid OptionalArrayOptionalCatTest",
			s: OptionalArrayOptionalCatTest{
				S: []Cat{{}},
			},
			wantNoErr: true,
		},
		{
			name: "Invalid OptionalArrayOptionalCatTest.S[0]",
			s: OptionalArrayOptionalCatTest{
				S: []Cat{{Name: "123", Age: 5}},
			},
			wantMessage: "S[0].Name: '123' does validate as 'alpha'",
		},

		// ArrayInterfaceTest
		{
			name: "Valid ArrayInterfaceTest",
			s: ArrayInterfaceTest{
				S: []interface{}{Cat{Name: "neko", Age: 5}},
			},
			wantNoErr: true,
		},
		{
			name: "Valid ArrayInterfaceTest.S[1] nil",
			s: ArrayInterfaceTest{
				S: []interface{}{Cat{Name: "neko", Age: 5}, nil},
			},
			wantNoErr: true,
		},
		{
			name: "Invalid ArrayInterfaceTest.S empty",
			s: ArrayInterfaceTest{
				S: []interface{}{},
			},
			wantMessage: "S: '<Array>' does validate as 'required'",
		},
		{
			name: "Invalid ArrayInterfaceTest.S[1] empty",
			s: ArrayInterfaceTest{
				S: []interface{}{Cat{Name: "neko", Age: 5}, Cat{}},
			},
			wantMessage: "S[1].Name: '' does validate as 'required';S[1].Name: '' does validate as 'alpha';S[1].Age: '0' does validate as 'required'",
		},

		// ArrayPointerTest
		{
			name: "Valid ArrayPointerTest",
			s: ArrayPointerTest{
				S: []*Cat{{Name: "neko", Age: 5}},
			},
			wantNoErr: true,
		},
		{
			name: "Valid ArrayPointerTest.S[1] nil",
			s: ArrayPointerTest{
				S: []*Cat{{Name: "neko", Age: 5}, nil},
			},
			wantNoErr: true,
		},
		{
			name: "Invalid ArrayPointerTest.S empty",
			s: ArrayPointerTest{
				S: []*Cat{},
			},
			wantMessage: "S: '<Array>' does validate as 'required'",
		},
		{
			name: "Invalid ArrayPointerTest.S[1] empty",
			s: ArrayPointerTest{
				S: []*Cat{{Name: "neko", Age: 5}, {}},
			},
			wantMessage: "S[1].Name: '' does validate as 'required';S[1].Name: '' does validate as 'alpha';S[1].Age: '0' does validate as 'required'",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateStruct(tc.s)

			if tc.wantNoErr {
				if err != nil {
					t.Error(err)
				}
				return
			}
			assertValidationError(t, tc.wantMessage, err)
		})
	}
}

func TestValidateStruct_Map(t *testing.T) {
	type (
		Cat struct {
			Name string `valid:"required,alpha"`
			Age  int    `valid:"required"`
		}

		MapStringTest struct {
			M map[string]string `valid:"len(3); req,alpha"`
		}

		RequiredMapStringTest struct {
			M map[string]string `valid:"required ;"`
		}

		MapRequiredStringTest struct {
			M map[string]string `valid:"; required"`
		}

		MapCatTest struct {
			M map[string]Cat `valid:"required"`
		}

		OptionalMapCatTest struct {
			M map[string]Cat `valid:"optional"`
		}

		OptionalMapCatTest2 struct {
			M map[string]Cat `valid:"optional ;"`
		}

		MapOptionalCatTest struct {
			M map[string]Cat `valid:"required; optional"`
		}

		OptionalMapOptionalCatTest struct {
			M map[string]Cat `valid:"optional ; optional"`
		}

		MapInterfaceTest struct {
			M map[string]interface{} `valid:"required"`
		}

		MapPointerTest struct {
			M map[string]*Cat `valid:"required"`
		}
	)

	testcases := []struct {
		name        string
		s           interface{}
		wantNoErr   bool
		wantMessage string
	}{
		// MapStringTest
		{
			name: "Valid MapStringTest",
			s: MapStringTest{
				M: map[string]string{"key1": "a", "key2": "b", "key3": "c"},
			},
			wantNoErr: true,
		},
		{
			name: "Invalid MapStringTest length",
			s: MapStringTest{
				M: map[string]string{"key1": "a", "key2": "b"},
			},
			wantMessage: "M: '<Map>' does validate as 'len(3)'",
		},
		{
			name: "Invalid MapStringTest.M[key1]",
			s: MapStringTest{
				M: map[string]string{"key1": "", "key2": "b", "key3": "c"},
			},
			wantMessage: "M[key1]: '' does validate as 'req';M[key1]: '' does validate as 'alpha'",
		},

		// RequiredMapStringTest
		{
			name: "Valid RequiredMapStringTest",
			s: RequiredMapStringTest{
				M: map[string]string{"key1": ""},
			},
			wantNoErr: true,
		},
		{
			name: "Invalid RequiredMapStringTest empty",
			s: RequiredMapStringTest{
				M: map[string]string{},
			},
			wantMessage: "M: '<Map>' does validate as 'required'",
		},
		{
			name: "Invalid RequiredMapStringTest nil",
			s: RequiredMapStringTest{
				M: nil,
			},
			wantMessage: "M: '<Map>' does validate as 'required'",
		},

		// MapRequiredStringTest
		{
			name: "Valid MapRequiredStringTest",
			s: MapRequiredStringTest{
				M: map[string]string{},
			},
			wantNoErr: true,
		},
		{
			name: "Valid MapRequiredStringTest",
			s: MapRequiredStringTest{
				M: nil,
			},
			wantNoErr: true,
		},
		{
			name: "Invalid MapRequiredStringTest.M[key1] empty",
			s: MapRequiredStringTest{
				M: map[string]string{"key1": ""},
			},
			wantMessage: "M[key1]: '' does validate as 'required'",
		},

		// MapCatTest
		{
			name: "Valid MapCatTest",
			s: MapCatTest{
				M: map[string]Cat{
					"key1": {Name: "neko", Age: 5},
				},
			},
			wantNoErr: true,
		},
		{
			name: "Invalid MapCatTest",
			s: MapCatTest{
				M: nil,
			},
			wantMessage: "M: '<Map>' does validate as 'required'",
		},
		{
			name: "Invalid MapCatTest.M[key1]",
			s: MapCatTest{
				M: map[string]Cat{"key1": {}},
			},
			wantMessage: "M[key1].Name: '' does validate as 'required';M[key1].Name: '' does validate as 'alpha';M[key1].Age: '0' does validate as 'required'",
		},

		// OptionalMapCatTest
		{
			name: "Valid OptionalMapCatTest empty",
			s: OptionalMapCatTest{
				M: map[string]Cat{},
			},
			wantNoErr: true,
		},
		{
			name: "Valid OptionalMapCatTest nil",
			s: OptionalMapCatTest{
				M: nil,
			},
			wantNoErr: true,
		},
		{
			name: "Invalid OptionalMapCatTest.M[key1]",
			s: OptionalMapCatTest{
				M: map[string]Cat{
					"key1": {Name: "123", Age: 5},
				},
			},
			wantMessage: "M[key1].Name: '123' does validate as 'alpha'",
		},

		// OptionalMapCatTest2
		{
			name: "Valid OptionalMapCatTest2 empty",
			s: OptionalMapCatTest2{
				M: map[string]Cat{},
			},
			wantNoErr: true,
		},
		{
			name: "Valid OptionalMapCatTest2 nil",
			s: OptionalMapCatTest2{
				M: nil,
			},
			wantNoErr: true,
		},
		{
			name: "Invalid OptionalMapCatTest2.M[key1]",
			s: OptionalMapCatTest2{
				M: map[string]Cat{
					"key1": {Name: "123", Age: 5},
				},
			},
			wantMessage: "M[key1].Name: '123' does validate as 'alpha'",
		},

		// MapOptionalCatTest
		{
			name: "Valid MapOptionalCatTest",
			s: MapOptionalCatTest{
				M: map[string]Cat{
					"key1": {},
				},
			},
			wantNoErr: true,
		},
		{
			name: "Invalid MapOptionalCatTest",
			s: MapOptionalCatTest{
				M: map[string]Cat{},
			},
			wantMessage: "M: '<Map>' does validate as 'required'",
		},
		{
			name: "Invalid MapOptionalCatTest.M[key1]",
			s: MapOptionalCatTest{
				M: map[string]Cat{
					"key1": {Name: "123", Age: 5},
				},
			},
			wantMessage: "M[key1].Name: '123' does validate as 'alpha'",
		},

		// OptionalMapOptionalCatTest
		{
			name: "Valid OptionalMapOptionalCatTest empty",
			s: OptionalMapOptionalCatTest{
				M: map[string]Cat{},
			},
			wantNoErr: true,
		},
		{
			name: "Valid OptionalMapOptionalCatTest nil",
			s: OptionalMapOptionalCatTest{
				M: nil,
			},
			wantNoErr: true,
		},
		{
			name: "Valid OptionalMapOptionalCatTest",
			s: OptionalMapOptionalCatTest{
				M: map[string]Cat{
					"key1": {},
				},
			},
			wantNoErr: true,
		},
		{
			name: "Invalid OptionalMapOptionalCatTest.M[key1]",
			s: OptionalMapOptionalCatTest{
				M: map[string]Cat{
					"key1": {Name: "123", Age: 5},
				},
			},
			wantMessage: "M[key1].Name: '123' does validate as 'alpha'",
		},

		// MapInterfaceTest
		{
			name: "Valid MapInterfaceTest",
			s: MapInterfaceTest{
				M: map[string]interface{}{
					"key1": Cat{Name: "neko", Age: 5},
				},
			},
			wantNoErr: true,
		},
		{
			name: "Valid MapInterfaceTest.M[key1] nil",
			s: MapInterfaceTest{
				M: map[string]interface{}{"key1": nil},
			},
			wantNoErr: true,
		},
		{
			name: "Invalid MapInterfaceTest.M empty",
			s: MapInterfaceTest{
				M: map[string]interface{}{},
			},
			wantMessage: "M: '<Map>' does validate as 'required'",
		},
		{
			name: "Invalid MapInterfaceTest.M nil",
			s: MapInterfaceTest{
				M: nil,
			},
			wantMessage: "M: '<Map>' does validate as 'required'",
		},
		{
			name: "Invalid MapInterfaceTest.M[key1] empty",
			s: MapInterfaceTest{
				M: map[string]interface{}{"key1": Cat{}},
			},
			wantMessage: "M[key1].Name: '' does validate as 'required';M[key1].Name: '' does validate as 'alpha';M[key1].Age: '0' does validate as 'required'",
		},

		// MapPointerTest
		{
			name: "Valid MapPointerTest",
			s: MapPointerTest{
				M: map[string]*Cat{
					"key1": {Name: "neko", Age: 5},
				},
			},
			wantNoErr: true,
		},
		{
			name: "Valid MapPointerTest.M[key1] nil",
			s: MapPointerTest{
				M: map[string]*Cat{"key1": nil},
			},
			wantNoErr: true,
		},
		{
			name: "Invalid MapPointerTest.M nil",
			s: MapPointerTest{
				M: nil,
			},
			wantMessage: "M: '<Map>' does validate as 'required'",
		},
		{
			name: "Invalid MapPointerTest.M[key1] empty",
			s: MapPointerTest{
				M: map[string]*Cat{"key1": {}},
			},
			wantMessage: "M[key1].Name: '' does validate as 'required';M[key1].Name: '' does validate as 'alpha';M[key1].Age: '0' does validate as 'required'",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateStruct(tc.s)

			if tc.wantNoErr {
				if err != nil {
					t.Error(err)
				}
				return
			}
			assertValidationError(t, tc.wantMessage, err)
		})
	}
}

func TestValidateStruct_OptionalStruct(t *testing.T) {
	type (
		Cat struct {
			Name string `valid:"alpha"`
		}
		OptionalCat struct {
			Cat Cat `valid:"optional"`
		}
	)

	testcases := []struct {
		name        string
		s           interface{}
		wantMessage string
		wantNoError bool
	}{
		{
			name:        "Valid",
			s:           OptionalCat{},
			wantNoError: true,
		},
		{
			name: "Invalid OptionalCat.Cat.Name",
			s: OptionalCat{
				Cat: Cat{Name: "123"},
			},
			wantMessage: "Cat.Name: '123' does validate as 'alpha'",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateStruct(tc.s)

			if tc.wantNoError {
				if err != nil {
					t.Error(err)
				}
				return
			}
			assertValidationError(t, tc.wantMessage, err)
		})
	}
}

func TestValidateStruct_Skip(t *testing.T) {
	type (
		Cat struct {
			Name string `valid:"alpha"`
		}

		SkipTest struct {
			Cat Cat `valid:"-"`
		}
	)

	testcases := []struct {
		name string
		s    interface{}
	}{
		{
			name: "Valid empty",
			s:    SkipTest{},
		},
		{
			name: "Valid skip",
			s:    SkipTest{Cat: Cat{Name: "123"}},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateStruct(tc.s)
			if err != nil {
				t.Fatal(err)
			}
		})
	}

	t.Run("Valid var", func(t *testing.T) {
		err := validator.ValidateVar(SkipTest{Cat: Cat{Name: "123"}}, "-")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestValidateStruct_InvalidTag(t *testing.T) {
	type (
		InvalidTag struct {
			Name string `valid:"unknown"`
		}

		InvalidTagTest struct {
			T InvalidTag
		}

		ArrayInvalidTagTest struct {
			S []InvalidTag `valid:"required; required"`
		}

		MapInvalidTagTest struct {
			M map[string]InvalidTag `valid:"required; required"`
		}
	)

	testcases := []struct {
		name        string
		s           interface{}
		wantNoError bool
		wantMessage string
	}{
		{
			name:        "Valid nil",
			s:           nil,
			wantNoError: true,
		},
		{
			name:        "Invalid struct type",
			s:           "test",
			wantMessage: "struct type required",
		},

		{
			name:        "Invalid InvalidTag",
			s:           InvalidTag{},
			wantMessage: "parse: tag unknown function not found",
		},
		{
			name: "Invalid InvalidTagTest",
			s: InvalidTagTest{
				T: InvalidTag{},
			},
			wantMessage: "parse: tag unknown function not found",
		},
		{
			name: "Invalid ArrayInvalidTagTest",
			s: ArrayInvalidTagTest{
				S: []InvalidTag{{}},
			},
			wantMessage: "parse: tag unknown function not found",
		},
		{
			name: "Invalid MapInvalidTagTest",
			s: MapInvalidTagTest{
				M: map[string]InvalidTag{
					"key1": {},
				},
			},
			wantMessage: "parse: tag unknown function not found",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateStruct(tc.s)
			if tc.wantNoError {
				if err != nil {
					t.Error(err)
				}
				return
			}

			if err == nil {
				t.Fatal("err is nil")
			}
			if err.Error() != tc.wantMessage {
				t.Fatalf("Message want %v, but got %v", tc.wantMessage, err)
			}
		})
	}
}

func TestValidateStruct_NotStruct(t *testing.T) {
	testcases := []struct {
		s          interface{}
		wantGotErr bool
		wantErrStr string
	}{
		{
			s:          nil,
			wantGotErr: false,
			wantErrStr: "",
		},
		{
			s:          "test",
			wantGotErr: true,
			wantErrStr: "struct type required",
		},
	}

	for _, tc := range testcases {
		err := validator.ValidateStruct(tc.s)
		if tc.wantGotErr {
			if err == nil {
				t.Fatal("want error, but got nil")
			}
			if tc.wantErrStr != err.Error() {
				t.Errorf("want %v, got %v", tc.wantErrStr, err)
			}
		} else if err != nil {
			t.Errorf("want nil, got %v", err)
		}
	}
}

func TestValidateStruct_PrivateField(t *testing.T) {
	type (
		privateFieldTest struct {
			name string `valid:"numeric"`
		}
	)

	testcases := []struct {
		name string
		s    interface{}
	}{
		{
			name: "Valid not numeric value",
			s: privateFieldTest{
				name: "abc",
			},
		},
		{
			name: "Valid empty value",
			s:    privateFieldTest{},
		},
	}

	for _, tc := range testcases {
		err := validator.ValidateStruct(tc.s)
		if err != nil {
			t.Errorf("want nil, but got %v", err)
		}
	}
}

func TestValidateStructContext(t *testing.T) {
	type (
		SimpleTest struct {
			Str string `valid:"required"`
		}
	)
	s := SimpleTest{Str: "str"}

	err := validator.ValidateStructContext(context.Background(), s)
	if err != nil {
		t.Errorf("want err nil, but got %v", err)
	}
}

func TestValidateVar(t *testing.T) {
	err := validator.ValidateVar("test", "req")
	if err != nil {
		t.Errorf("want err nil, but got %v", err)
	}
}

func TestValidateVarContext(t *testing.T) {
	err := validator.ValidateVarContext(context.Background(), "test", "req")
	if err != nil {
		t.Errorf("want err nil, but got %v", err)
	}
}

func TestWithFunc(t *testing.T) {
	v := validator.New()

	v.Apply(validator.WithFunc("test", func(_ context.Context, _ validator.Field, _ validator.FuncOption) (bool, error) {
		return false, fmt.Errorf("set func failure")
	}))

	err := v.ValidateVar("", "test")
	if err == nil {
		t.Fatal("want error, but got nil")
	}

	wantError := ": an internal error occurred in 'test': set func failure"
	if want, got := wantError, err.Error(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestWithFuncMap(t *testing.T) {
	v := validator.New()

	v.Apply(validator.WithFuncMap(map[string]validator.Func{
		"test": func(_ context.Context, _ validator.Field, _ validator.FuncOption) (bool, error) {
			return false, fmt.Errorf("set func failure")
		},
	}))

	err := v.ValidateVar("", "test")
	if err == nil {
		t.Fatal("want error, but got nil")
	}

	wantError := ": an internal error occurred in 'test': set func failure"
	if want, got := wantError, err.Error(); want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestWithAdapters(t *testing.T) {
	v := validator.New()

	var str string
	v.Apply(validator.WithAdapters(
		func(fn validator.Func) validator.Func {
			return func(ctx context.Context, f validator.Field, o validator.FuncOption) (bool, error) {
				str += "1"
				return fn(ctx, f, o)
			}
		},
		func(fn validator.Func) validator.Func {
			return func(ctx context.Context, f validator.Field, o validator.FuncOption) (bool, error) {
				str += "2"
				return fn(ctx, f, o)
			}
		},
		func(fn validator.Func) validator.Func {
			return func(ctx context.Context, f validator.Field, o validator.FuncOption) (bool, error) {
				str += "3"
				return fn(ctx, f, o)
			}
		},
	))

	err := v.ValidateVar("test", "req")
	if err != nil {
		t.Fatal(err)
	}

	if str != "123" {
		t.Errorf("want 123, got %v", str)
	}
}

func TestWithTagKey(t *testing.T) {
	type (
		TagKeyTest struct {
			Name string `tag_key_test:"required"`
		}
	)
	v := validator.New(validator.WithTagKey("tag_key_test"))

	err := v.ValidateStruct(TagKeyTest{})
	if err == nil {
		t.Fatal("want error, but got nil")
	}

	wantErrMessage := "Name: '' does validate as 'required'"
	if want, got := wantErrMessage, err.Error(); want != got {
		t.Errorf("want error message %v, but got %v", want, got)
	}
}

func TestWithSuppressErrorFieldValue(t *testing.T) {
	v := validator.New(validator.WithSuppressErrorFieldValue())

	err := v.ValidateVar("aaa", "numeric")
	if err == nil {
		t.Fatal("want error, but got nil")
	}

	if want, got := ": The value does validate as 'numeric'", err.Error(); want != got {
		t.Errorf("want error message %v, but got %v", want, got)
	}
}

func assertValidationError(t *testing.T, expectMessage string, err error) {
	if err == nil {
		t.Fatal("err want `error`, but got `nil`")
	}

	errs, ok := validator.ToErrors(err)
	if !ok {
		t.Fatal("ToErrors want `true`, but got `false`")
	}

	if expectMessage != errs.Error() {
		t.Fatalf("Message want `%v`\nbut got `%v`", expectMessage, errs.Error())
	}
}
