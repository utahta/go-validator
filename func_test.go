package validator_test

import (
	"testing"

	"github.com/utahta/go-validator"
)

func Test_required(t *testing.T) {
	t.Parallel()

	type (
		Cat struct {
			Name string `valid:"required"`
		}
	)
	const tag = "required"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid string", v.ValidateVar("str", tag), false},
		{"valid int", v.ValidateVar(1, tag), false},
		{"valid int8", v.ValidateVar(int8(1), tag), false},
		{"valid int16", v.ValidateVar(int16(1), tag), false},
		{"valid int32", v.ValidateVar(int32(1), tag), false},
		{"valid int64", v.ValidateVar(int64(1), tag), false},
		{"valid uint", v.ValidateVar(uint(1), tag), false},
		{"valid uint8", v.ValidateVar(uint8(1), tag), false},
		{"valid uint16", v.ValidateVar(uint16(1), tag), false},
		{"valid uint32", v.ValidateVar(uint32(1), tag), false},
		{"valid uint64", v.ValidateVar(uint64(1), tag), false},
		{"valid float32", v.ValidateVar(float32(1.0), tag), false},
		{"valid float64", v.ValidateVar(float64(1.0), tag), false},
		{"valid slice", v.ValidateVar([]int{1}, tag), false},
		{"valid array", v.ValidateVar([1]int{1}, tag), false},
		{"valid array2", v.ValidateVar([1]int{}, tag), false},
		{"valid map", v.ValidateVar(map[int]int{1: 1}, tag), false},
		{"valid ptr", v.ValidateVar(&Cat{Name: "str"}, tag), false},
		{"valid struct", v.ValidateVar(Cat{Name: "str"}, tag), false},
		{"valid bool", v.ValidateVar(true, tag), false},

		{"invalid string", v.ValidateVar("", tag), true},
		{"invalid int", v.ValidateVar(0, tag), true},
		{"invalid int8", v.ValidateVar(int8(0), tag), true},
		{"invalid int16", v.ValidateVar(int16(0), tag), true},
		{"invalid int32", v.ValidateVar(int32(0), tag), true},
		{"invalid int64", v.ValidateVar(int64(0), tag), true},
		{"invalid uint", v.ValidateVar(uint(0), tag), true},
		{"invalid uint8", v.ValidateVar(uint8(0), tag), true},
		{"invalid uint16", v.ValidateVar(uint16(0), tag), true},
		{"invalid uint32", v.ValidateVar(uint32(0), tag), true},
		{"invalid uint64", v.ValidateVar(uint64(0), tag), true},
		{"invalid float32", v.ValidateVar(float32(0.0), tag), true},
		{"invalid float64", v.ValidateVar(float64(0.0), tag), true},
		{"invalid slice", v.ValidateVar([]int{}, tag), true},
		{"invalid array", v.ValidateVar([0]int{}, tag), true},
		{"invalid map", v.ValidateVar(map[int]int{}, tag), true},
		{"invalid ptr", v.ValidateVar((*Cat)(nil), tag), true},
		{"invalid struct", v.ValidateVar(Cat{}, tag), true},
		{"invalid bool", v.ValidateVar(false, tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_empty(t *testing.T) {
	t.Parallel()

	type (
		Cat struct {
			Name string `valid:"empty"`
		}
	)
	const tag = "empty"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid string", v.ValidateVar("", tag), false},
		{"valid int", v.ValidateVar(0, tag), false},
		{"valid int8", v.ValidateVar(int8(0), tag), false},
		{"valid int16", v.ValidateVar(int16(0), tag), false},
		{"valid int32", v.ValidateVar(int32(0), tag), false},
		{"valid int64", v.ValidateVar(int64(0), tag), false},
		{"valid uint", v.ValidateVar(uint(0), tag), false},
		{"valid uint8", v.ValidateVar(uint8(0), tag), false},
		{"valid uint16", v.ValidateVar(uint16(0), tag), false},
		{"valid uint32", v.ValidateVar(uint32(0), tag), false},
		{"valid uint64", v.ValidateVar(uint64(0), tag), false},
		{"valid float32", v.ValidateVar(float32(0.0), tag), false},
		{"valid float64", v.ValidateVar(float64(0.0), tag), false},
		{"valid slice", v.ValidateVar([]int{}, tag), false},
		{"valid array", v.ValidateVar([0]int{}, tag), false},
		{"valid map", v.ValidateVar(map[int]int{}, tag), false},
		{"valid ptr", v.ValidateVar((*Cat)(nil), tag), false},
		{"valid struct", v.ValidateVar(Cat{}, tag), false},
		{"valid bool", v.ValidateVar(false, tag), false},

		{"invalid string", v.ValidateVar("str", tag), true},
		{"invalid int", v.ValidateVar(1, tag), true},
		{"invalid int8", v.ValidateVar(int8(1), tag), true},
		{"invalid int16", v.ValidateVar(int16(1), tag), true},
		{"invalid int32", v.ValidateVar(int32(1), tag), true},
		{"invalid int64", v.ValidateVar(int64(1), tag), true},
		{"invalid uint", v.ValidateVar(uint(1), tag), true},
		{"invalid uint8", v.ValidateVar(uint8(1), tag), true},
		{"invalid uint16", v.ValidateVar(uint16(1), tag), true},
		{"invalid uint32", v.ValidateVar(uint32(1), tag), true},
		{"invalid uint64", v.ValidateVar(uint64(1), tag), true},
		{"invalid float32", v.ValidateVar(float32(1.0), tag), true},
		{"invalid float64", v.ValidateVar(float64(1.0), tag), true},
		{"invalid slice", v.ValidateVar([]int{1}, tag), true},
		{"invalid array", v.ValidateVar([1]int{1}, tag), true},
		{"invalid array2", v.ValidateVar([1]int{}, tag), true},
		{"invalid map", v.ValidateVar(map[int]int{1: 1}, tag), true},
		{"invalid ptr", v.ValidateVar(&Cat{Name: "str"}, tag), true},
		{"invalid struct", v.ValidateVar(Cat{Name: "str"}, tag), true},
		{"invalid bool", v.ValidateVar(true, tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_alpha(t *testing.T) {
	t.Parallel()

	const tag = "alpha"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("abc", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid numeric", v.ValidateVar("123", tag), true},
		{"invalid alpha numeric", v.ValidateVar("abc123", tag), true},
		{"invalid multibyte", v.ValidateVar("てすと", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_alphanum(t *testing.T) {
	t.Parallel()

	const tag = "alphanum"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("abc123", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid symbol", v.ValidateVar("-", tag), true},
		{"invalid multibyte", v.ValidateVar("てすと", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_alphaunicode(t *testing.T) {
	t.Parallel()

	const tag = "alphaunicode"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("abcテスト", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid numeric", v.ValidateVar("123", tag), true},
		{"invalid symbol", v.ValidateVar("-", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_alphanumunicode(t *testing.T) {
	t.Parallel()

	const tag = "alphanumunicode"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("abc123テスト", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid symbol", v.ValidateVar("-", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_numeric(t *testing.T) {
	t.Parallel()

	const tag = "numeric"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("12345", tag), false},
		{"valid", v.ValidateVar("12345.678", tag), false},
		{"valid", v.ValidateVar(12345, tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid alpha", v.ValidateVar("abc", tag), true},
		{"invalid alphanum", v.ValidateVar("abc123", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_number(t *testing.T) {
	t.Parallel()

	const tag = "number"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("12345", tag), false},
		{"valid", v.ValidateVar(12345, tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid float", v.ValidateVar("12345.678", tag), true},
		{"invalid alpha", v.ValidateVar("abc", tag), true},
		{"invalid alphanum", v.ValidateVar("abc123", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_hexadecimal(t *testing.T) {
	t.Parallel()

	const tag = "hexadecimal"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("0099aaffAAFF", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid alpha", v.ValidateVar("GG", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_hexcolor(t *testing.T) {
	t.Parallel()

	const tag = "hexcolor"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("#FFFFFF", tag), false},
		{"valid", v.ValidateVar("#FFF", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid length=1", v.ValidateVar("#F", tag), true},
		{"invalid length=2", v.ValidateVar("#FF", tag), true},
		{"invalid length=4", v.ValidateVar("#FFFF", tag), true},
		{"invalid length=5", v.ValidateVar("#FFFFF", tag), true},
		{"invalid length=7", v.ValidateVar("#FFFFFFF", tag), true},
		{"invalid no prefix", v.ValidateVar("FFF", tag), true},
		{"invalid alpha", v.ValidateVar("#GG", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_rgb(t *testing.T) {
	t.Parallel()

	const tag = "rgb"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("rgb(0,128,255)", tag), false},
		{"valid", v.ValidateVar("rgb( 0, 128 , 255 )", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("rgb(-1,0,0)", tag), true},
		{"invalid value", v.ValidateVar("rgb(0,0,256)", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_rgba(t *testing.T) {
	t.Parallel()

	const tag = "rgba"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("rgba(0,128,255,0)", tag), false},
		{"valid", v.ValidateVar("rgba(0,128,255,0.5)", tag), false},
		{"valid", v.ValidateVar("rgba( 0, 128, 255, 1 )", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("rgba(0,0,0,-1)", tag), true},
		{"invalid value", v.ValidateVar("rgba(0,0,0,1.1)", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_hsl(t *testing.T) {
	t.Parallel()

	const tag = "hsl"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("hsl(0,0%,0%)", tag), false},
		{"valid", v.ValidateVar("hsl( 360, 100%, 100% )", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("hsl(-1,0%,0%)", tag), true},
		{"invalid value", v.ValidateVar("hsl(361,0%,0%)", tag), true},
		{"invalid value", v.ValidateVar("hsl(0,-1%,0%)", tag), true},
		{"invalid value", v.ValidateVar("hsl(0,101%,0%)", tag), true},
		{"invalid value", v.ValidateVar("hsl(0,0%,-1%)", tag), true},
		{"invalid value", v.ValidateVar("hsl(0,0%,101%)", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_hsla(t *testing.T) {
	t.Parallel()

	const tag = "hsla"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("hsla(0,0%,0%,0)", tag), false},
		{"valid", v.ValidateVar("hsla(0,0%,0%,0.5)", tag), false},
		{"valid", v.ValidateVar("hsla( 360, 100%, 100%,1 )", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("hsla(-1,0%,0%,0)", tag), true},
		{"invalid value", v.ValidateVar("hsla(361,0%,0%,0)", tag), true},
		{"invalid value", v.ValidateVar("hsla(0,-1%,0%,0)", tag), true},
		{"invalid value", v.ValidateVar("hsla(0,101%,0%,0)", tag), true},
		{"invalid value", v.ValidateVar("hsla(0,0%,-1%,0)", tag), true},
		{"invalid value", v.ValidateVar("hsla(0,0%,101%,0)", tag), true},
		{"invalid value", v.ValidateVar("hsla(0,0%,0%,-1)", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_email(t *testing.T) {
	t.Parallel()

	const tag = "email"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("test@localhost.local", tag), false},
		{"valid", v.ValidateVar("test.test@localhost.local.local", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("test@localhost", tag), true},
		{"invalid value", v.ValidateVar("@localhost.local", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_base64(t *testing.T) {
	t.Parallel()

	const tag = "base64"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("Aza9", tag), false},
		{"valid", v.ValidateVar("Az==", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("A", tag), true},
		{"invalid value", v.ValidateVar("Az", tag), true},
		{"invalid value", v.ValidateVar("Az-_", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_base64url(t *testing.T) {
	t.Parallel()

	const tag = "base64url"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("Az-_", tag), false},
		{"valid", v.ValidateVar("Az==", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("A", tag), true},
		{"invalid value", v.ValidateVar("Az", tag), true},
		{"invalid value", v.ValidateVar("Aza90", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_isbn10(t *testing.T) {
	t.Parallel()

	const tag = "isbn10"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("0111122229", tag), false},
		{"valid", v.ValidateVar("011112222X", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("011112222", tag), true},
		{"invalid value", v.ValidateVar("011112222A", tag), true},
		{"invalid value", v.ValidateVar("01111222", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_isbn13(t *testing.T) {
	t.Parallel()

	const tag = "isbn13"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("9780111122229", tag), false},
		{"valid", v.ValidateVar("9790111122229", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("9770111122229", tag), true},
		{"invalid value", v.ValidateVar("978011112222X", tag), true},
		{"invalid value", v.ValidateVar("978111122229", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_url(t *testing.T) {
	t.Parallel()

	const tag = "url"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("http://localhost", tag), false},
		{"valid", v.ValidateVar("file://localhost", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("//localhost", tag), true},
		{"invalid value", v.ValidateVar("localhost", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_uri(t *testing.T) {
	t.Parallel()

	const tag = "uri"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("http://localhost", tag), false},
		{"valid", v.ValidateVar("file://localhost", tag), false},
		{"valid", v.ValidateVar("//localhost", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("localhost", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_uuid(t *testing.T) {
	t.Parallel()

	const tag = "uuid"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("0099aaff-09af-09af-09af-003399aaccff", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("AA99aaff-09af-09af-09af-003399aaccff", tag), true},
		{"invalid value", v.ValidateVar("0099aaff-09af-09af-09af-003399aaccf", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_uuid3(t *testing.T) {
	t.Parallel()

	const tag = "uuid3"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("0099aaff-09af-39af-09af-003399aaccff", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("0099aaff-09af-09af-09af-003399aaccff", tag), true},
		{"invalid value", v.ValidateVar("AA99aaff-09af-39af-09af-003399aaccff", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_uuid4(t *testing.T) {
	t.Parallel()

	const tag = "uuid4"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("0099aaff-09af-49af-89af-003399aaccff", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("0099aaff-09af-39af-89af-003399aaccff", tag), true},
		{"invalid value", v.ValidateVar("AA99aaff-09af-49af-89af-003399aaccff", tag), true},
		{"invalid value", v.ValidateVar("0099aaff-09af-49af-09af-003399aaccff", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_uuid5(t *testing.T) {
	t.Parallel()

	const tag = "uuid5"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("0099aaff-09af-59af-89af-003399aaccff", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("0099aaff-09af-49af-89af-003399aaccff", tag), true},
		{"invalid value", v.ValidateVar("AA99aaff-09af-59af-89af-003399aaccff", tag), true},
		{"invalid value", v.ValidateVar("0099aaff-09af-59af-09af-003399aaccff", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_ascii(t *testing.T) {
	t.Parallel()

	const tag = "ascii"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("\x00\x7F", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("\x80", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_printableascii(t *testing.T) {
	t.Parallel()

	const tag = "printableascii"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("\x20\x7E", tag), false},
		{"valid", v.ValidateVar(" ~", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("\x00", tag), true},
		{"invalid value", v.ValidateVar("\x7F", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_multibyte(t *testing.T) {
	t.Parallel()

	const tag = "multibyte"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("\x80", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("\x00", tag), true},
		{"invalid value", v.ValidateVar("\x7F", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_datauri(t *testing.T) {
	t.Parallel()

	const tag = "datauri"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("data:image/png;base64,AA==", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("data:image/png;base64,", tag), true},
		{"invalid value", v.ValidateVar("data:image/png;base64,AA", tag), true},
		{"invalid value", v.ValidateVar("data:image/;base64,AA==", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_latitude(t *testing.T) {
	t.Parallel()

	const tag = "latitude"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("-90.0", tag), false},
		{"valid", v.ValidateVar("+90", tag), false},
		{"valid", v.ValidateVar("45.45", tag), false},
		{"valid", v.ValidateVar("-45.45", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("-91", tag), true},
		{"invalid value", v.ValidateVar("+90.1", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_longitude(t *testing.T) {
	t.Parallel()

	const tag = "longitude"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("-180.0", tag), false},
		{"valid", v.ValidateVar("+180", tag), false},
		{"valid", v.ValidateVar("90.45", tag), false},
		{"valid", v.ValidateVar("-90.45", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("-181", tag), true},
		{"invalid value", v.ValidateVar("+180.1", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_ssn(t *testing.T) {
	t.Parallel()

	const tag = "ssn"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("111223333", tag), false},
		{"valid", v.ValidateVar("111-22-3333", tag), false},
		{"valid", v.ValidateVar("111 22 3333", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("11122333", tag), true},
		{"invalid value", v.ValidateVar("1111223333", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_semver(t *testing.T) {
	t.Parallel()

	const tag = "semver"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("0.0.1", tag), false},
		{"valid", v.ValidateVar("v0.0.1", tag), false},
		{"valid", v.ValidateVar("0.0.1-rc", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("0.0.0.1", tag), true},
		{"invalid value", v.ValidateVar("0.1", tag), true},
		{"invalid value", v.ValidateVar("1", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_katakana(t *testing.T) {
	t.Parallel()

	const tag = "katakana"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("テスト", tag), false},
		{"valid", v.ValidateVar("ﾃｽﾄ", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("試験", tag), true},
		{"invalid value", v.ValidateVar("てすと", tag), true},
		{"invalid value", v.ValidateVar("123", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_hiragana(t *testing.T) {
	t.Parallel()

	const tag = "hiragana"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("てすと", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("試験", tag), true},
		{"invalid value", v.ValidateVar("テスト", tag), true},
		{"invalid value", v.ValidateVar("ﾃｽﾄ", tag), true},
		{"invalid value", v.ValidateVar("123", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_fullwidth(t *testing.T) {
	t.Parallel()

	const tag = "fullwidth"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("てすと", tag), false},
		{"valid", v.ValidateVar("テスト", tag), false},
		{"valid", v.ValidateVar("試験", tag), false},
		{"valid", v.ValidateVar("　ー！", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("ﾃｽﾄ", tag), true},
		{"invalid value", v.ValidateVar("abc", tag), true},
		{"invalid value", v.ValidateVar("123", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_halfwidth(t *testing.T) {
	t.Parallel()

	const tag = "halfwidth"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid", v.ValidateVar("ﾃｽﾄ", tag), false},
		{"valid", v.ValidateVar("abc", tag), false},
		{"valid", v.ValidateVar("123", tag), false},

		{"invalid empty", v.ValidateVar("", tag), true},
		{"invalid value", v.ValidateVar("テスト", tag), true},
		{"invalid value", v.ValidateVar("　ー！", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_length_minmax(t *testing.T) {
	t.Parallel()

	const tag = "len(2|3)"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid string min", v.ValidateVar("aa", tag), false},
		{"valid string max", v.ValidateVar("aaa", tag), false},
		{"valid multibyte string min", v.ValidateVar("ああ", tag), false},
		{"valid multibyte string max", v.ValidateVar("あああ", tag), false},
		{"valid int min", v.ValidateVar(2, tag), false},
		{"valid int max", v.ValidateVar(3, tag), false},
		{"valid uint min", v.ValidateVar(uint(2), tag), false},
		{"valid uint max", v.ValidateVar(uint(3), tag), false},
		{"valid float min", v.ValidateVar(float32(2.0), tag), false},
		{"valid float max", v.ValidateVar(float32(3.0), tag), false},
		{"valid slice min", v.ValidateVar([]int{2, 2}, tag), false},
		{"valid slice max", v.ValidateVar([]int{2, 2, 2}, tag), false},
		{"valid map min", v.ValidateVar(map[int]int{1: 2, 2: 2}, tag), false},
		{"valid map max", v.ValidateVar(map[int]int{1: 2, 2: 2, 3: 2}, tag), false},

		{"invalid string min", v.ValidateVar("a", tag), true},
		{"invalid string max", v.ValidateVar("aaaa", tag), true},
		{"invalid multibyte string min", v.ValidateVar("あ", tag), false},
		{"invalid multibyte string max", v.ValidateVar("ああああ", tag), false},
		{"invalid int min", v.ValidateVar(1, tag), true},
		{"invalid int max", v.ValidateVar(4, tag), true},
		{"invalid uint min", v.ValidateVar(uint(1), tag), true},
		{"invalid uint max", v.ValidateVar(uint(4), tag), true},
		{"invalid float min", v.ValidateVar(float32(1.0), tag), true},
		{"invalid float max", v.ValidateVar(float32(4.0), tag), true},
		{"invalid slice min", v.ValidateVar([]int{2}, tag), true},
		{"invalid slice max", v.ValidateVar([]int{2, 2, 2, 2}, tag), true},
		{"invalid map min", v.ValidateVar(map[int]int{1: 2}, tag), true},
		{"invalid map max", v.ValidateVar(map[int]int{1: 2, 2: 2, 3: 2, 4: 2}, tag), true},
		{"invalid bool", v.ValidateVar(true, tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, _ := tc.err.(validator.Errors)
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}

	t.Run("invalid param len", func(t *testing.T) {
		wantError := ": an internal error occurred in 'len(2|3|4)': invalid params len"
		err := v.ValidateVar(2, "len(2|3|4)")
		if err.Error() != wantError {
			t.Errorf("want `%v`, got `%v`", wantError, err)
		}
	})
}

func Test_length_equal(t *testing.T) {
	t.Parallel()

	const tag = "len(2)"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid string", v.ValidateVar("aa", tag), false},
		{"valid multibyte string", v.ValidateVar("ああ", tag), false},
		{"valid int", v.ValidateVar(2, tag), false},
		{"valid uint", v.ValidateVar(uint(2), tag), false},
		{"valid float", v.ValidateVar(float32(2.0), tag), false},
		{"valid slice", v.ValidateVar([]int{2, 2}, tag), false},
		{"valid map", v.ValidateVar(map[int]int{1: 2, 2: 2}, tag), false},

		{"invalid string", v.ValidateVar("a", tag), true},
		{"invalid multibyte string", v.ValidateVar("あ", tag), false},
		{"invalid int", v.ValidateVar(1, tag), true},
		{"invalid uint", v.ValidateVar(uint(1), tag), true},
		{"invalid float", v.ValidateVar(float32(1.0), tag), true},
		{"invalid slice", v.ValidateVar([]int{2}, tag), true},
		{"invalid map", v.ValidateVar(map[int]int{1: 2}, tag), true},
		{"invalid bool", v.ValidateVar(true, tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, _ := tc.err.(validator.Errors)
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}
}

func Test_length_invalidTag(t *testing.T) {
	t.Parallel()

	v := validator.New()

	testcases := []struct {
		name           string
		err            error
		wantErrMessage string
	}{
		{"equal", v.ValidateVar("a", "len(aaa)"), `: an internal error occurred in 'len(aaa)': strconv.ParseInt: parsing "aaa": invalid syntax`},
		{"min", v.ValidateVar("a", "len(aaa|3)"), `: an internal error occurred in 'len(aaa|3)': strconv.ParseInt: parsing "aaa": invalid syntax`},
		{"max", v.ValidateVar("a", "len(1|aaa)"), `: an internal error occurred in 'len(1|aaa)': strconv.ParseInt: parsing "aaa": invalid syntax`},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err == nil {
				t.Error("want error, but got nil")
				return
			}

			gotErrs, has := tc.err.(validator.Errors)
			if !has {
				t.Errorf("want validator Errors, but not %v", gotErrs)
				return
			}

			if want, got := tc.wantErrMessage, gotErrs[0].Error(); want != got {
				t.Errorf("want error message %v, but got %v", want, got)
			}
		})
	}
}

func Test_eqlength(t *testing.T) {
	t.Parallel()

	const tag = "eq(2)"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid string", v.ValidateVar("aa", tag), false},
		{"valid multibyte string", v.ValidateVar("ああ", tag), false},
		{"valid int", v.ValidateVar(2, tag), false},
		{"valid slice", v.ValidateVar([]int{2, 2}, tag), false},
		{"valid map", v.ValidateVar(map[int]int{1: 2, 2: 2}, tag), false},

		{"invalid string", v.ValidateVar("a", tag), true},
		{"invalid multibyte string", v.ValidateVar("あ", tag), true},
		{"invalid int", v.ValidateVar(1, tag), true},
		{"invalid slice", v.ValidateVar([]int{2}, tag), true},
		{"invalid map", v.ValidateVar(map[int]int{1: 2}, tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, _ := tc.err.(validator.Errors)
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}

	t.Run("invalid param len", func(t *testing.T) {
		wantError := ": an internal error occurred in 'eq(2|3)': invalid params len"
		err := v.ValidateVar("aa", "eq(2|3)")
		if err == nil {
			t.Error("want error, but got nil")
			return
		}
		if err.Error() != wantError {
			t.Errorf("want `%v`, got `%v`", wantError, err)
		}
	})
}

func Test_eqlength_invalidTag(t *testing.T) {
	t.Parallel()

	const tag = "eq(aaa)"
	v := validator.New()

	testcases := []struct {
		name           string
		err            error
		wantErrMessage string
	}{
		{"string", v.ValidateVar("a", tag), `: an internal error occurred in 'eq(aaa)': strconv.ParseInt: parsing "aaa": invalid syntax`},
		{"int", v.ValidateVar(int(1), tag), `: an internal error occurred in 'eq(aaa)': strconv.ParseInt: parsing "aaa": invalid syntax`},
		{"uint", v.ValidateVar(uint(1), tag), `: an internal error occurred in 'eq(aaa)': strconv.ParseUint: parsing "aaa": invalid syntax`},
		{"float", v.ValidateVar(float64(1), tag), `: an internal error occurred in 'eq(aaa)': strconv.ParseFloat: parsing "aaa": invalid syntax`},
		{"slice", v.ValidateVar([]int{2}, tag), `: an internal error occurred in 'eq(aaa)': strconv.ParseInt: parsing "aaa": invalid syntax`},
		{"map", v.ValidateVar(map[int]int{1: 2}, tag), `: an internal error occurred in 'eq(aaa)': strconv.ParseInt: parsing "aaa": invalid syntax`},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err == nil {
				t.Error("want error, but got nil")
				return
			}

			gotErrs, has := tc.err.(validator.Errors)
			if !has {
				t.Errorf("want validator Errors, but not %v", gotErrs)
				return
			}

			if want, got := tc.wantErrMessage, gotErrs[0].Error(); want != got {
				t.Errorf("want error message %v, but got %v", want, got)
			}
		})
	}
}

func Test_minlength(t *testing.T) {
	t.Parallel()

	const tag = "min(2)"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid string", v.ValidateVar("aa", tag), false},
		{"valid multibyte string", v.ValidateVar("ああ", tag), false},
		{"valid int", v.ValidateVar(2, tag), false},
		{"valid slice", v.ValidateVar([]int{2, 2}, tag), false},
		{"valid map", v.ValidateVar(map[int]int{1: 2, 2: 2}, tag), false},

		{"invalid string", v.ValidateVar("a", tag), true},
		{"invalid multibyte string", v.ValidateVar("あ", tag), true},
		{"invalid int", v.ValidateVar(1, tag), true},
		{"invalid slice", v.ValidateVar([]int{2}, tag), true},
		{"invalid map", v.ValidateVar(map[int]int{1: 2}, tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, _ := tc.err.(validator.Errors)
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}

	t.Run("invalid param len", func(t *testing.T) {
		wantError := ": an internal error occurred in 'min(2|3)': invalid params len"
		err := v.ValidateVar("aa", "min(2|3)")
		if err == nil {
			t.Error("want error, but got nil")
			return
		}
		if err.Error() != wantError {
			t.Errorf("want `%v`, got `%v`", wantError, err)
		}
	})
}

func Test_minlength_invalidTag(t *testing.T) {
	t.Parallel()

	const tag = "min(aaa)"
	v := validator.New()

	testcases := []struct {
		name           string
		err            error
		wantErrMessage string
	}{
		{"string", v.ValidateVar("a", tag), `: an internal error occurred in 'min(aaa)': strconv.ParseInt: parsing "aaa": invalid syntax`},
		{"int", v.ValidateVar(int(1), tag), `: an internal error occurred in 'min(aaa)': strconv.ParseInt: parsing "aaa": invalid syntax`},
		{"uint", v.ValidateVar(uint(1), tag), `: an internal error occurred in 'min(aaa)': strconv.ParseUint: parsing "aaa": invalid syntax`},
		{"float", v.ValidateVar(float64(1), tag), `: an internal error occurred in 'min(aaa)': strconv.ParseFloat: parsing "aaa": invalid syntax`},
		{"slice", v.ValidateVar([]int{2}, tag), `: an internal error occurred in 'min(aaa)': strconv.ParseInt: parsing "aaa": invalid syntax`},
		{"map", v.ValidateVar(map[int]int{1: 2}, tag), `: an internal error occurred in 'min(aaa)': strconv.ParseInt: parsing "aaa": invalid syntax`},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err == nil {
				t.Error("want error, but got nil")
				return
			}

			gotErrs, has := tc.err.(validator.Errors)
			if !has {
				t.Errorf("want validator Errors, but not %v", gotErrs)
				return
			}

			if want, got := tc.wantErrMessage, gotErrs[0].Error(); want != got {
				t.Errorf("want error message %v, but got %v", want, got)
			}
		})
	}
}

func Test_maxlength(t *testing.T) {
	t.Parallel()

	const tag = "max(2)"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid string", v.ValidateVar("aa", tag), false},
		{"valid multibyte string", v.ValidateVar("ああ", tag), false},
		{"valid int", v.ValidateVar(2, tag), false},
		{"valid slice", v.ValidateVar([]int{2, 2}, tag), false},
		{"valid map", v.ValidateVar(map[int]int{1: 2, 2: 2}, tag), false},

		{"invalid string", v.ValidateVar("aaa", tag), true},
		{"invalid multibyte string", v.ValidateVar("あああ", tag), true},
		{"invalid int", v.ValidateVar(3, tag), true},
		{"invalid slice", v.ValidateVar([]int{2, 2, 2}, tag), true},
		{"invalid map", v.ValidateVar(map[int]int{1: 2, 2: 2, 3: 2}, tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, _ := tc.err.(validator.Errors)
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}

	t.Run("invalid param len", func(t *testing.T) {
		wantError := ": an internal error occurred in 'max(2|3)': invalid params len"
		err := v.ValidateVar("aa", "max(2|3)")
		if err == nil {
			t.Error("want error, but got nil")
			return
		}
		if err.Error() != wantError {
			t.Errorf("want `%v`, got `%v`", wantError, err)
		}
	})
}

func Test_maxlength_invalidTag(t *testing.T) {
	t.Parallel()

	const tag = "max(aaa)"
	v := validator.New()

	testcases := []struct {
		name           string
		err            error
		wantErrMessage string
	}{
		{"string", v.ValidateVar("a", tag), `: an internal error occurred in 'max(aaa)': strconv.ParseInt: parsing "aaa": invalid syntax`},
		{"int", v.ValidateVar(int(1), tag), `: an internal error occurred in 'max(aaa)': strconv.ParseInt: parsing "aaa": invalid syntax`},
		{"uint", v.ValidateVar(uint(1), tag), `: an internal error occurred in 'max(aaa)': strconv.ParseUint: parsing "aaa": invalid syntax`},
		{"float", v.ValidateVar(float64(1), tag), `: an internal error occurred in 'max(aaa)': strconv.ParseFloat: parsing "aaa": invalid syntax`},
		{"slice", v.ValidateVar([]int{2}, tag), `: an internal error occurred in 'max(aaa)': strconv.ParseInt: parsing "aaa": invalid syntax`},
		{"map", v.ValidateVar(map[int]int{1: 2}, tag), `: an internal error occurred in 'max(aaa)': strconv.ParseInt: parsing "aaa": invalid syntax`},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err == nil {
				t.Error("want error, but got nil")
				return
			}

			gotErrs, has := tc.err.(validator.Errors)
			if !has {
				t.Errorf("want validator Errors, but not %v", gotErrs)
				return
			}

			if want, got := tc.wantErrMessage, gotErrs[0].Error(); want != got {
				t.Errorf("want error message %v, but got %v", want, got)
			}
		})
	}
}

func Test_or(t *testing.T) {
	t.Parallel()

	const tag = "or(alpha|numeric)"
	v := validator.New()

	testcases := []struct {
		name   string
		err    error
		hasErr bool
	}{
		{"valid alpha", v.ValidateVar("abc", tag), false},
		{"valid numeric", v.ValidateVar("123", tag), false},
		{"invalid", v.ValidateVar("===", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotErrs, hasErr := tc.err.(validator.Errors)
			if tc.hasErr != hasErr {
				t.Errorf("want hasErr %v, got %v", tc.hasErr, hasErr)
			}
			if tc.hasErr {
				if len(gotErrs) == 0 {
					t.Fatal("errors is empty")
				}
				if gotErrs[0].Tag().String() != tag {
					t.Errorf("want tag name %v, got %v", tag, gotErrs[0].Tag())
				}
			}
		})
	}

	t.Run("invalid tag", func(t *testing.T) {
		wantError := ": an internal error occurred in 'or(unknown|numeric)': parse: tag unknown function not found"
		err := v.ValidateVar("===", "or(unknown|numeric)")
		if err == nil {
			t.Error("want error, but got nil")
			return
		}
		if err.Error() != wantError {
			t.Errorf("want `%v`, got `%v`", wantError, err)
		}
	})
}
