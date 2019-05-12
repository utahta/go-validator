package validator

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"unicode/utf8"
)

type (
	// Func is a type of validate function.
	Func func(context.Context, Field, FuncOption) (bool, error)

	// FuncMap is a map of validate function type.
	FuncMap map[string]Func

	// FuncOption is a option.
	FuncOption struct {
		// Params has any parameters.
		Params []string

		v *Validator
	}

	// Adapter is a validate function adapter.
	Adapter func(Func) Func
)

var (
	defaultFuncMap = FuncMap{
		"required":        hasValue,
		"req":             hasValue,
		"empty":           isZeroValue,
		"zero":            isZeroValue,
		"alpha":           isAlpha,
		"alphanum":        isAlphaNum,
		"alphaunicode":    isAlphaUnicode,
		"alphanumunicode": isAlphaNumUnicode,
		"numeric":         isNumeric,
		"number":          isNumber,
		"hexadecimal":     isHexadecimal,
		"hexcolor":        isHexcolor,
		"rgb":             isRGB,
		"rgba":            isRGBA,
		"hsl":             isHSL,
		"hsla":            isHSLA,
		"email":           isEmail,
		"base64":          isBase64,
		"base64url":       isBase64URL,
		"isbn10":          isISBN10,
		"isbn13":          isISBN13,
		"isbn":            isISBN13,
		"url":             isURL,
		"uri":             isURI,
		"uuid":            isUUID,
		"uuid3":           isUUID3,
		"uuid4":           isUUID4,
		"uuid5":           isUUID5,
		"ascii":           isASCII,
		"printableascii":  isPrintableASCII,
		"multibyte":       isMultibyte,
		"datauri":         isDataURI,
		"latitude":        isLatitude,
		"longitude":       isLongitude,
		"ssn":             isSSN,
		"semver":          isSemver,
		"katakana":        isKatakana,
		"hiragana":        isHiragana,
		"fullwidth":       isFullWidth,
		"halfwidth":       isHalfWidth,

		// has parameters
		"len":        length,
		"length":     length,
		"range":      length,
		"min":        minLength,
		"max":        maxLength,
		"strlen":     strLength,
		"strlength":  strLength,
		"strmin":     strMinLength,
		"strmax":     strMaxLength,
		"runelen":    strLength,
		"runelength": strLength,
		"or":         or,
	}

	defaultAdapters []Adapter
)

// apply applies left to right.
func apply(fn Func, adapters ...Adapter) Func {
	for i := len(adapters) - 1; i >= 0; i-- {
		fn = adapters[i](fn)
	}
	return fn
}

// isEmpty return true if value is zero and else false.
func isEmpty(f Field) bool {
	v := f.current
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr, reflect.Chan, reflect.Func:
		return v.IsNil()
	}

	return v.IsValid() && reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func hasValue(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return !isEmpty(f), nil
}

func isZeroValue(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return isEmpty(f), nil
}

func isAlpha(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return alphaRegex.MatchString(f.String()), nil
}

func isAlphaNum(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return alphaNumericRegex.MatchString(f.String()), nil
}

func isAlphaUnicode(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return alphaUnicodeRegex.MatchString(f.String()), nil
}

func isAlphaNumUnicode(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return alphaNumericUnicodeRegex.MatchString(f.String()), nil
}

func isEmail(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return emailRegex.MatchString(f.String()), nil
}

func isBase64(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return base64Regex.MatchString(f.String()), nil
}

func isBase64URL(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return base64URLRegex.MatchString(f.String()), nil
}

func isISBN10(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return isbn10Regex.MatchString(f.String()), nil
}

func isISBN13(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return isbn13Regex.MatchString(f.String()), nil
}

func isURL(_ context.Context, f Field, _ FuncOption) (bool, error) {
	u, err := url.ParseRequestURI(f.String())
	return err == nil && len(u.Scheme) > 0, nil
}

func isURI(_ context.Context, f Field, _ FuncOption) (bool, error) {
	_, err := url.ParseRequestURI(f.String())
	return err == nil, nil
}

func isNumeric(_ context.Context, f Field, _ FuncOption) (bool, error) {
	switch f.current.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true, nil
	}
	return numericRegex.MatchString(f.String()), nil
}

func isNumber(_ context.Context, f Field, _ FuncOption) (bool, error) {
	switch f.current.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true, nil
	}
	return numberRegex.MatchString(f.String()), nil
}

func isHexadecimal(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return hexadecimalRegex.MatchString(f.String()), nil
}

func isHexcolor(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return hexcolorRegex.MatchString(f.String()), nil
}

func isRGB(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return rgbRegex.MatchString(f.String()), nil
}

func isRGBA(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return rgbaRegex.MatchString(f.String()), nil
}

func isHSL(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return hslRegex.MatchString(f.String()), nil
}

func isHSLA(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return hslaRegex.MatchString(f.String()), nil
}

func isUUID(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return uuidRegex.MatchString(f.String()), nil
}

func isUUID3(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return uuid3Regex.MatchString(f.String()), nil
}

func isUUID4(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return uuid4Regex.MatchString(f.String()), nil
}

func isUUID5(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return uuid5Regex.MatchString(f.String()), nil
}

func isASCII(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return asciiRegex.MatchString(f.String()), nil
}

func isPrintableASCII(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return printableASCIIRegex.MatchString(f.String()), nil
}

func isMultibyte(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return multibyteRegex.MatchString(f.String()), nil
}

func isDataURI(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return dataURIRegex.MatchString(f.String()), nil
}

func isLatitude(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return latitudeRegex.MatchString(f.String()), nil
}

func isLongitude(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return longitudeRegex.MatchString(f.String()), nil
}

func isSSN(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return ssnRegex.MatchString(f.String()), nil
}

func isSemver(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return semverRegex.MatchString(f.String()), nil
}

func isKatakana(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return katakanaRegex.MatchString(f.String()), nil
}

func isHiragana(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return hiraganaRegex.MatchString(f.String()), nil
}

func isFullWidth(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return fullWidthRegex.MatchString(f.String()), nil
}

func isHalfWidth(_ context.Context, f Field, _ FuncOption) (bool, error) {
	return halfWidthRegex.MatchString(f.String()), nil
}

func minLength(_ context.Context, f Field, opt FuncOption) (bool, error) {
	var minStr string
	if len(opt.Params) == 1 {
		minStr = opt.Params[0]
	} else {
		return false, fmt.Errorf("invalid params len")
	}

	v := f.current
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Map, reflect.Slice:
		min, err := parseInt64(minStr)
		if err != nil {
			return false, err
		}
		return min <= int64(v.Len()), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		min, err := parseInt64(minStr)
		if err != nil {
			return false, err
		}
		return min <= v.Int(), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		min, err := parseUint64(minStr)
		if err != nil {
			return false, err
		}
		return min <= v.Uint(), nil

	case reflect.Float32, reflect.Float64:
		min, err := parseFloat64(minStr)
		if err != nil {
			return false, err
		}
		return min <= v.Float(), nil
	}
	return false, nil
}

func maxLength(_ context.Context, f Field, opt FuncOption) (bool, error) {
	var maxStr string
	if len(opt.Params) == 1 {
		maxStr = opt.Params[0]
	} else {
		return false, fmt.Errorf("invalid params len")
	}

	v := f.current
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Map, reflect.Slice:
		max, err := parseInt64(maxStr)
		if err != nil {
			return false, err
		}
		return int64(v.Len()) <= max, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		max, err := parseInt64(maxStr)
		if err != nil {
			return false, err
		}
		return v.Int() <= max, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		max, err := parseUint64(maxStr)
		if err != nil {
			return false, err
		}
		return v.Uint() <= max, nil

	case reflect.Float32, reflect.Float64:
		max, err := parseFloat64(maxStr)
		if err != nil {
			return false, err
		}
		return v.Float() <= max, nil
	}
	return false, nil
}

func eqLength(_ context.Context, f Field, opt FuncOption) (bool, error) {
	if len(opt.Params) != 1 {
		return false, fmt.Errorf("invalid params len")
	}
	str := opt.Params[0]

	v := f.current
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Map, reflect.Slice:
		i, err := parseInt64(str)
		if err != nil {
			return false, err
		}
		return int64(v.Len()) == i, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := parseInt64(str)
		if err != nil {
			return false, err
		}
		return v.Int() == i, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		i, err := parseUint64(str)
		if err != nil {
			return false, err
		}
		return v.Uint() == i, nil

	case reflect.Float32, reflect.Float64:
		i, err := parseFloat64(str)
		if err != nil {
			return false, err
		}
		return v.Float() == i, nil
	}
	return false, nil
}

func length(ctx context.Context, f Field, opt FuncOption) (bool, error) {
	switch len(opt.Params) {
	case 1:
		eq, err := eqLength(ctx, f, FuncOption{Params: opt.Params})
		if err != nil {
			return false, err
		}
		return eq, nil

	case 2:
		min, err := minLength(ctx, f, FuncOption{Params: opt.Params[:1]})
		if err != nil {
			return false, err
		}
		max, err := maxLength(ctx, f, FuncOption{Params: opt.Params[1:]})
		if err != nil {
			return false, err
		}
		return min && max, nil

	}
	return false, fmt.Errorf("invalid params len")
}

func strMinLength(_ context.Context, f Field, opt FuncOption) (bool, error) {
	if len(opt.Params) != 1 {
		return false, fmt.Errorf("invalid params len")
	}

	v := f.current
	switch v.Kind() {
	case reflect.String:
		min, err := parseInt64(opt.Params[0])
		if err != nil {
			return false, err
		}
		return min <= int64(utf8.RuneCountInString(v.String())), nil
	}
	return false, nil
}

func strMaxLength(_ context.Context, f Field, opt FuncOption) (bool, error) {
	if len(opt.Params) != 1 {
		return false, fmt.Errorf("invalid params len")
	}

	v := f.current
	switch v.Kind() {
	case reflect.String:
		max, err := parseInt64(opt.Params[0])
		if err != nil {
			return false, err
		}
		return int64(utf8.RuneCountInString(v.String())) <= max, nil
	}
	return false, nil
}

func strEqLength(_ context.Context, f Field, opt FuncOption) (bool, error) {
	if len(opt.Params) != 1 {
		return false, fmt.Errorf("invalid params len")
	}

	v := f.current
	switch v.Kind() {
	case reflect.String:
		i, err := parseInt64(opt.Params[0])
		if err != nil {
			return false, err
		}
		return int64(utf8.RuneCountInString(v.String())) == i, nil
	}
	return false, nil
}

func strLength(ctx context.Context, f Field, opt FuncOption) (bool, error) {
	switch len(opt.Params) {
	case 1:
		eq, err := strEqLength(ctx, f, FuncOption{Params: opt.Params})
		if err != nil {
			return false, err
		}
		return eq, nil

	case 2:
		min, err := strMinLength(ctx, f, FuncOption{Params: opt.Params[:1]})
		if err != nil {
			return false, err
		}
		max, err := strMaxLength(ctx, f, FuncOption{Params: opt.Params[1:]})
		if err != nil {
			return false, err
		}
		return min && max, nil

	}
	return false, fmt.Errorf("invalid params len")
}

func or(_ context.Context, f Field, opt FuncOption) (bool, error) {
	if opt.v == nil {
		return false, fmt.Errorf("validator is nil")
	}

	for _, rawTag := range opt.Params {
		err := opt.v.ValidateVar(f.Interface(), rawTag)
		if err == nil {
			return true, nil
		}
		if _, ok := err.(Errors); !ok {
			return false, err
		}
	}
	return false, nil
}

func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func parseUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

func parseFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
