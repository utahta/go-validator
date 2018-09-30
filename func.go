package validator

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"unicode/utf8"
)

type (
	// Func is a type of validator function.
	Func func(Field, FuncOption) (bool, error)

	// FuncMap is a map of validator function type.
	FuncMap map[string]Func

	// FuncOption is a option.
	FuncOption struct {
		validator *Validator

		// Params has any parameters.
		Params []string

		// Optional is a optional flag. if true, empty value is always valid.
		Optional bool
	}
)

var defaultFuncMap = FuncMap{
	"required":        with(hasValue),
	"alpha":           with(isAlpha),
	"alphanum":        with(isAlphaNum),
	"alphaunicode":    with(isAlphaUnicode),
	"alphanumunicode": with(isAlphaNumUnicode),
	"numeric":         with(isNumeric),
	"number":          with(isNumber),
	"hexadecimal":     with(isHexadecimal),
	"hexcolor":        with(isHexcolor),
	"email":           with(isEmail),
	"base64":          with(isBase64),
	"url":             with(isURL),
	"uri":             with(isURI),
	"uuid":            with(isUUID),
	"uuid3":           with(isUUID3),
	"uuid4":           with(isUUID4),
	"uuid5":           with(isUUID5),

	// has params
	"len":        with(length),
	"length":     with(length),
	"range":      with(length),
	"min":        with(minLength),
	"max":        with(maxLength),
	"runelen":    with(strLength),
	"runelength": with(strLength),
	"strlen":     with(strLength),
	"strlength":  with(strLength),
	"strmin":     with(strMinLength),
	"strmax":     with(strMaxLength),
	"or":         with(or),
}

func with(fn Func) Func {
	return func(f Field, opt FuncOption) (bool, error) {
		if opt.Optional && isEmpty(f) {
			return true, nil
		}
		return fn(f, opt)
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