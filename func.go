package validator

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"unicode/utf8"
)

type (
	// Func is a type of validate function.
	Func func(Field, FuncOption) (bool, error)

	// FuncMap is a map of validate function type.
	FuncMap map[string]Func

	// FuncOption is a option.
	FuncOption struct {
		validator *Validator

		// Params has any parameters.
		Params []string

		// Optional is a optional flag. if true, empty value is always valid.
		Optional bool
	}

	// Adapter is a validate function adapter.
	Adapter func(Func) Func
)

var (
	defaultFuncMap = FuncMap{
		"required":        hasValue,
		"req":             hasValue,
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

	defaultAdapters = []Adapter{
		withOptional(),
	}
)

func apply(fn Func, adapters ...Adapter) Func {
	for _, adapter := range adapters {
		fn = adapter(fn)
	}
	return fn
}

func withOptional() Adapter {
	return func(fn Func) Func {
		return func(f Field, opt FuncOption) (bool, error) {
			if opt.Optional && isEmpty(f) {
				return true, nil
			}
			return fn(f, opt)
		}
	}
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

func hasValue(f Field, _ FuncOption) (bool, error) {
	return !isEmpty(f), nil
}

func isAlpha(f Field, _ FuncOption) (bool, error) {
	return alphaRegex.MatchString(f.String()), nil
}

func isAlphaNum(f Field, _ FuncOption) (bool, error) {
	return alphaNumericRegex.MatchString(f.String()), nil
}

func isAlphaUnicode(f Field, _ FuncOption) (bool, error) {
	return alphaUnicodeRegex.MatchString(f.String()), nil
}

func isAlphaNumUnicode(f Field, _ FuncOption) (bool, error) {
	return alphaNumericUnicodeRegex.MatchString(f.String()), nil
}

func isEmail(f Field, _ FuncOption) (bool, error) {
	return emailRegex.MatchString(f.String()), nil
}

func isBase64(f Field, _ FuncOption) (bool, error) {
	return base64Regex.MatchString(f.String()), nil
}

func isBase64URL(f Field, _ FuncOption) (bool, error) {
	return base64URLRegex.MatchString(f.String()), nil
}

func isISBN10(f Field, _ FuncOption) (bool, error) {
	return isbn10Regex.MatchString(f.String()), nil
}

func isISBN13(f Field, _ FuncOption) (bool, error) {
	return isbn13Regex.MatchString(f.String()), nil
}

func isURL(f Field, _ FuncOption) (bool, error) {
	u, err := url.ParseRequestURI(f.String())
	return err == nil && len(u.Scheme) > 0, nil
}

func isURI(f Field, _ FuncOption) (bool, error) {
	_, err := url.ParseRequestURI(f.String())
	return err == nil, nil
}

func isNumeric(f Field, _ FuncOption) (bool, error) {
	switch f.current.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true, nil
	}
	return numericRegex.MatchString(f.String()), nil
}

func isNumber(f Field, _ FuncOption) (bool, error) {
	switch f.current.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true, nil
	}
	return numberRegex.MatchString(f.String()), nil
}

func isHexadecimal(f Field, _ FuncOption) (bool, error) {
	return hexadecimalRegex.MatchString(f.String()), nil
}

func isHexcolor(f Field, _ FuncOption) (bool, error) {
	return hexcolorRegex.MatchString(f.String()), nil
}

func isRGB(f Field, _ FuncOption) (bool, error) {
	return rgbRegex.MatchString(f.String()), nil
}

func isRGBA(f Field, _ FuncOption) (bool, error) {
	return rgbaRegex.MatchString(f.String()), nil
}

func isHSL(f Field, _ FuncOption) (bool, error) {
	return hslRegex.MatchString(f.String()), nil
}

func isHSLA(f Field, _ FuncOption) (bool, error) {
	return hslaRegex.MatchString(f.String()), nil
}

func isUUID(f Field, _ FuncOption) (bool, error) {
	return uuidRegex.MatchString(f.String()), nil
}

func isUUID3(f Field, _ FuncOption) (bool, error) {
	return uuid3Regex.MatchString(f.String()), nil
}

func isUUID4(f Field, _ FuncOption) (bool, error) {
	return uuid4Regex.MatchString(f.String()), nil
}

func isUUID5(f Field, _ FuncOption) (bool, error) {
	return uuid5Regex.MatchString(f.String()), nil
}

func isASCII(f Field, _ FuncOption) (bool, error) {
	return asciiRegex.MatchString(f.String()), nil
}

func isPrintableASCII(f Field, _ FuncOption) (bool, error) {
	return printableASCIIRegex.MatchString(f.String()), nil
}

func isMultibyte(f Field, _ FuncOption) (bool, error) {
	return multibyteRegex.MatchString(f.String()), nil
}

func isDataURI(f Field, _ FuncOption) (bool, error) {
	return dataURIRegex.MatchString(f.String()), nil
}

func isLatitude(f Field, _ FuncOption) (bool, error) {
	return latitudeRegex.MatchString(f.String()), nil
}

func isLongitude(f Field, _ FuncOption) (bool, error) {
	return longitudeRegex.MatchString(f.String()), nil
}

func isSSN(f Field, _ FuncOption) (bool, error) {
	return ssnRegex.MatchString(f.String()), nil
}

func isSemver(f Field, _ FuncOption) (bool, error) {
	return semverRegex.MatchString(f.String()), nil
}

func isKatakana(f Field, _ FuncOption) (bool, error) {
	return katakanaRegex.MatchString(f.String()), nil
}

func isHiragana(f Field, _ FuncOption) (bool, error) {
	return hiraganaRegex.MatchString(f.String()), nil
}

func isFullWidth(f Field, _ FuncOption) (bool, error) {
	return fullWidthRegex.MatchString(f.String()), nil
}

func isHalfWidth(f Field, _ FuncOption) (bool, error) {
	return halfWidthRegex.MatchString(f.String()), nil
}

func minLength(f Field, opt FuncOption) (bool, error) {
	var minStr string
	if len(opt.Params) == 1 {
		minStr = opt.Params[0]
	} else {
		return false, fmt.Errorf("invalid params len")
	}

	v := f.current
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Map, reflect.Slice:
		min, err := strconv.Atoi(minStr)
		if err != nil {
			return false, err
		}
		return min <= v.Len(), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		min, err := strconv.ParseInt(minStr, 10, 64)
		if err != nil {
			return false, err
		}
		return min <= v.Int(), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		min, err := strconv.ParseUint(minStr, 10, 64)
		if err != nil {
			return false, err
		}
		return min <= v.Uint(), nil

	case reflect.Float32, reflect.Float64:
		min, err := strconv.ParseFloat(minStr, 64)
		if err != nil {
			return false, err
		}
		return min <= v.Float(), nil
	}
	return false, nil
}

func maxLength(f Field, opt FuncOption) (bool, error) {
	var maxStr string
	if len(opt.Params) == 1 {
		maxStr = opt.Params[0]
	} else {
		return false, fmt.Errorf("invalid params len")
	}

	v := f.current
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Map, reflect.Slice:
		max, err := strconv.Atoi(maxStr)
		if err != nil {
			return false, err
		}
		return v.Len() <= max, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		max, err := strconv.ParseInt(maxStr, 10, 64)
		if err != nil {
			return false, err
		}
		return v.Int() <= max, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		max, err := strconv.ParseUint(maxStr, 10, 64)
		if err != nil {
			return false, err
		}
		return v.Uint() <= max, nil

	case reflect.Float32, reflect.Float64:
		max, err := strconv.ParseFloat(maxStr, 64)
		if err != nil {
			return false, err
		}
		return v.Float() <= max, nil
	}
	return false, nil
}

func eqLength(f Field, opt FuncOption) (bool, error) {
	if len(opt.Params) != 1 {
		return false, fmt.Errorf("invalid params len")
	}
	str := opt.Params[0]

	v := f.current
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Map, reflect.Slice:
		i, err := strconv.Atoi(str)
		if err != nil {
			return false, err
		}
		return v.Len() == i, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return false, err
		}
		return v.Int() == i, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		i, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return false, err
		}
		return v.Uint() == i, nil

	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return false, err
		}
		return v.Float() == i, nil
	}
	return false, nil
}

func length(f Field, opt FuncOption) (bool, error) {
	switch len(opt.Params) {
	case 1:
		eq, err := eqLength(f, FuncOption{Params: opt.Params})
		if err != nil {
			return false, err
		}
		return eq, nil

	case 2:
		min, err := minLength(f, FuncOption{Params: opt.Params[:1]})
		if err != nil {
			return false, err
		}
		max, err := maxLength(f, FuncOption{Params: opt.Params[1:]})
		if err != nil {
			return false, err
		}
		return min && max, nil

	}
	return false, fmt.Errorf("invalid params len")
}

func strMinLength(f Field, opt FuncOption) (bool, error) {
	if len(opt.Params) != 1 {
		return false, fmt.Errorf("invalid params len")
	}

	v := f.current
	switch v.Kind() {
	case reflect.String:
		min, err := strconv.Atoi(opt.Params[0])
		if err != nil {
			return false, err
		}
		return min <= utf8.RuneCountInString(v.String()), nil
	}
	return false, nil
}

func strMaxLength(f Field, opt FuncOption) (bool, error) {
	if len(opt.Params) != 1 {
		return false, fmt.Errorf("invalid params len")
	}

	v := f.current
	switch v.Kind() {
	case reflect.String:
		max, err := strconv.Atoi(opt.Params[0])
		if err != nil {
			return false, err
		}
		return utf8.RuneCountInString(v.String()) <= max, nil
	}
	return false, nil
}

func strEqLength(f Field, opt FuncOption) (bool, error) {
	if len(opt.Params) != 1 {
		return false, fmt.Errorf("invalid params len")
	}

	v := f.current
	switch v.Kind() {
	case reflect.String:
		i, err := strconv.Atoi(opt.Params[0])
		if err != nil {
			return false, err
		}
		return utf8.RuneCountInString(v.String()) == i, nil
	}
	return false, nil
}

func strLength(f Field, opt FuncOption) (bool, error) {
	switch len(opt.Params) {
	case 1:
		eq, err := strEqLength(f, FuncOption{Params: opt.Params})
		if err != nil {
			return false, err
		}
		return eq, nil

	case 2:
		min, err := strMinLength(f, FuncOption{Params: opt.Params[:1]})
		if err != nil {
			return false, err
		}
		max, err := strMaxLength(f, FuncOption{Params: opt.Params[1:]})
		if err != nil {
			return false, err
		}
		return min && max, nil

	}
	return false, fmt.Errorf("invalid params len")
}

func or(f Field, opt FuncOption) (bool, error) {
	if opt.validator == nil {
		return false, fmt.Errorf("validator is nil")
	}

	for _, rawTag := range opt.Params {
		err := opt.validator.ValidateVar(f.Interface(), rawTag)
		if err == nil {
			return true, nil
		}
		if _, ok := err.(Errors); !ok {
			return false, err
		}
	}
	return false, nil
}
