package validator

import (
	"net/url"
	"reflect"
)

type (
	Func    func(Field, ...string) (bool, error)
	FuncMap map[string]Func
)

var defaultFuncMap = FuncMap{
	"required":        hasValue,
	"alpha":           isAlpha,
	"alphanum":        isAlphaNum,
	"alphaunicode":    isAlphaUnicode,
	"alphanumunicode": isAlphaNumUnicode,
	"email":           isEmail,
	"url":             isURL,
	"uri":             isURI,
	"numeric":         isNumeric,
	"number":          isNumber,
	"uuid":            isUUID,
	"uuid3":           isUUID3,
	"uuid4":           isUUID4,
	"uuid5":           isUUID5,
}

func hasValue(f Field, _ ...string) (bool, error) {
	v := f.val
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() != 0, nil
	case reflect.Map, reflect.Slice:
		return v.Len() != 0 && !v.IsNil(), nil
	case reflect.Bool:
		return v.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() != 0, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() != 0, nil
	case reflect.Float32, reflect.Float64:
		return v.Float() != 0, nil
	case reflect.Interface, reflect.Ptr:
		return !v.IsNil(), nil
	}
	return v.IsValid() && v.Interface() != reflect.Zero(v.Type()).Interface(), nil
}

func isAlpha(f Field, _ ...string) (bool, error) {
	return alphaRegex.MatchString(f.Value()), nil
}

func isAlphaNum(f Field, _ ...string) (bool, error) {
	return alphaNumericRegex.MatchString(f.Value()), nil
}

func isAlphaUnicode(f Field, _ ...string) (bool, error) {
	return alphaUnicodeRegex.MatchString(f.Value()), nil
}

func isAlphaNumUnicode(f Field, _ ...string) (bool, error) {
	return alphaNumericUnicodeRegex.MatchString(f.Value()), nil
}

func isEmail(f Field, _ ...string) (bool, error) {
	return emailRegex.MatchString(f.Value()), nil
}

func isURL(f Field, _ ...string) (bool, error) {
	u, err := url.ParseRequestURI(f.Value())
	return err == nil && len(u.Scheme) > 0, nil
}

func isURI(f Field, _ ...string) (bool, error) {
	_, err := url.ParseRequestURI(f.Value())
	return err == nil, nil
}

func isNumeric(f Field, _ ...string) (bool, error) {
	return numericRegex.MatchString(f.Value()), nil
}

func isNumber(f Field, _ ...string) (bool, error) {
	return numberRegex.MatchString(f.Value()), nil
}

func isUUID(f Field, _ ...string) (bool, error) {
	return uuidRegex.MatchString(f.Value()), nil
}

func isUUID3(f Field, _ ...string) (bool, error) {
	return uuid3Regex.MatchString(f.Value()), nil
}

func isUUID4(f Field, _ ...string) (bool, error) {
	return uuid4Regex.MatchString(f.Value()), nil
}
func isUUID5(f Field, _ ...string) (bool, error) {
	return uuid5Regex.MatchString(f.Value()), nil
}
