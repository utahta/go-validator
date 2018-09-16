package validator_test

import (
	"testing"

	"github.com/utahta/go-validator"
)

func Test_required(t *testing.T) {
	t.Parallel()

	const tag = "required"
	v := validator.New()

	testcases := []struct {
		name    string
		err     error
		wantErr bool
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
		{"valid map", v.ValidateVar(map[int]int{1: 1}, tag), false},
		{"valid ptr", v.ValidateVar(&Str{Value: "str"}, tag), false},
		{"valid struct", v.ValidateVar(Str{Value: "str"}, tag), false},

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
		{"invalid array", v.ValidateVar([1]int{}, tag), true},
		{"invalid map", v.ValidateVar(map[int]int{}, tag), true},
		{"invalid ptr", v.ValidateVar((*Str)(nil), tag), true},
		{"invalid struct", v.ValidateVar(Str{}, tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if _, got := tc.err.(validator.Errors); tc.wantErr != got {
				t.Errorf("want %v, got %v", tc.wantErr, got)
			}
		})
	}
}

func Test_length_minmax(t *testing.T) {
	t.Parallel()

	const tag = "len(2|3)"
	v := validator.New()

	testcases := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{"valid string min", v.ValidateVar("aa", tag), false},
		{"valid string max", v.ValidateVar("aaa", tag), false},
		{"valid int min", v.ValidateVar(2, tag), false},
		{"valid int max", v.ValidateVar(3, tag), false},
		{"valid slice min", v.ValidateVar([]int{2, 2}, tag), false},
		{"valid slice max", v.ValidateVar([]int{2, 2, 2}, tag), false},
		{"valid map min", v.ValidateVar(map[int]int{1: 2, 2: 2}, tag), false},
		{"valid map max", v.ValidateVar(map[int]int{1: 2, 2: 2, 3: 2}, tag), false},

		{"invalid string min", v.ValidateVar("a", tag), true},
		{"invalid string max", v.ValidateVar("aaaa", tag), true},
		{"invalid int min", v.ValidateVar(1, tag), true},
		{"invalid int max", v.ValidateVar(4, tag), true},
		{"invalid slice min", v.ValidateVar([]int{2}, tag), true},
		{"invalid slice max", v.ValidateVar([]int{2, 2, 2, 2}, tag), true},
		{"invalid map min", v.ValidateVar(map[int]int{1: 2}, tag), true},
		{"invalid map max", v.ValidateVar(map[int]int{1: 2, 2: 2, 3: 2, 4: 2}, tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if _, got := tc.err.(validator.Errors); tc.wantErr != got {
				t.Errorf("want %v, got %v", tc.wantErr, got)
			}
		})
	}
}

func Test_length_equal(t *testing.T) {
	t.Parallel()

	const tag = "len(2)"
	v := validator.New()

	testcases := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{"valid string", v.ValidateVar("aa", tag), false},
		{"valid int", v.ValidateVar(2, tag), false},
		{"valid slice", v.ValidateVar([]int{2, 2}, tag), false},
		{"valid map", v.ValidateVar(map[int]int{1: 2, 2: 2}, tag), false},

		{"invalid string", v.ValidateVar("a", tag), true},
		{"invalid int", v.ValidateVar(1, tag), true},
		{"invalid slice", v.ValidateVar([]int{2}, tag), true},
		{"invalid map", v.ValidateVar(map[int]int{1: 2}, tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if errs, got := tc.err.(validator.Errors); tc.wantErr != got {
				t.Errorf("want %v, got %v %v", tc.wantErr, got, errs)
			}
		})
	}
}

func Test_strlength_minmax(t *testing.T) {
	t.Parallel()

	const tag = "strlen(2|3)"
	v := validator.New()

	testcases := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{"valid string min", v.ValidateVar("ああ", tag), false},
		{"valid string max", v.ValidateVar("あああ", tag), false},

		{"invalid string min", v.ValidateVar("あ", tag), true},
		{"invalid string max", v.ValidateVar("ああああ", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if _, got := tc.err.(validator.Errors); tc.wantErr != got {
				t.Errorf("want %v, got %v", tc.wantErr, got)
			}
		})
	}
}

func Test_strlength_equal(t *testing.T) {
	t.Parallel()

	const tag = "strlen(2)"
	v := validator.New()

	testcases := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{"valid string", v.ValidateVar("ああ", tag), false},

		{"invalid string", v.ValidateVar("あ", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if _, got := tc.err.(validator.Errors); tc.wantErr != got {
				t.Errorf("want %v, got %v", tc.wantErr, got)
			}
		})
	}
}

func Test_minlength(t *testing.T) {
	t.Parallel()

	const tag = "min(2)"
	v := validator.New()

	testcases := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{"valid string", v.ValidateVar("aa", tag), false},
		{"valid int", v.ValidateVar(2, tag), false},
		{"valid slice", v.ValidateVar([]int{2, 2}, tag), false},
		{"valid map", v.ValidateVar(map[int]int{1: 2, 2: 2}, tag), false},

		{"invalid string", v.ValidateVar("a", tag), true},
		{"invalid int", v.ValidateVar(1, tag), true},
		{"invalid slice", v.ValidateVar([]int{2}, tag), true},
		{"invalid map", v.ValidateVar(map[int]int{1: 2}, tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if errs, got := tc.err.(validator.Errors); tc.wantErr != got {
				t.Errorf("want %v, got %v %v", tc.wantErr, got, errs)
			}
		})
	}
}

func Test_maxlength(t *testing.T) {
	t.Parallel()

	const tag = "max(2)"
	v := validator.New()

	testcases := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{"valid string", v.ValidateVar("aa", tag), false},
		{"valid int", v.ValidateVar(2, tag), false},
		{"valid slice", v.ValidateVar([]int{2, 2}, tag), false},
		{"valid map", v.ValidateVar(map[int]int{1: 2, 2: 2}, tag), false},

		{"invalid string", v.ValidateVar("aaa", tag), true},
		{"invalid int", v.ValidateVar(3, tag), true},
		{"invalid slice", v.ValidateVar([]int{2, 2, 2}, tag), true},
		{"invalid map", v.ValidateVar(map[int]int{1: 2, 2: 2, 3: 2}, tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if errs, got := tc.err.(validator.Errors); tc.wantErr != got {
				t.Errorf("want %v, got %v %v", tc.wantErr, got, errs)
			}
		})
	}
}

func Test_strminlength(t *testing.T) {
	t.Parallel()

	const tag = "strmin(2)"
	v := validator.New()

	testcases := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{"valid string", v.ValidateVar("ああ", tag), false},

		{"invalid string", v.ValidateVar("あ", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if errs, got := tc.err.(validator.Errors); tc.wantErr != got {
				t.Errorf("want %v, got %v %v", tc.wantErr, got, errs)
			}
		})
	}
}

func Test_strmaxlength(t *testing.T) {
	t.Parallel()

	const tag = "strmax(2)"
	v := validator.New()

	testcases := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{"valid string", v.ValidateVar("ああ", tag), false},

		{"invalid string", v.ValidateVar("あああ", tag), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if errs, got := tc.err.(validator.Errors); tc.wantErr != got {
				t.Errorf("want %v, got %v %v", tc.wantErr, got, errs)
			}
		})
	}
}

func Test_or(t *testing.T) {
	t.Parallel()

	v := validator.New()

	testcases := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{"valid alpha", v.ValidateVar("abc", "alpha|numeric"), false},
		{"valid numeric", v.ValidateVar("123", "alpha|numeric"), false},
		{"invalid", v.ValidateVar("===", "alpha|numeric"), true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if _, got := tc.err.(validator.Errors); tc.wantErr != got {
				t.Errorf("want %v, got %v %v", tc.wantErr, got, tc.err)
			}
		})
	}
}
