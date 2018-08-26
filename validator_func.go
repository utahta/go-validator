package validator

import (
	"net/url"
	"reflect"
)

type (
	Func    func(Field, ...string) bool
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

func hasValue(f Field, _ ...string) bool {
	v := f.val
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() != 0
	case reflect.Map, reflect.Slice:
		return v.Len() != 0 && !v.IsNil()
	case reflect.Bool:
		return v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return v.Float() != 0
	case reflect.Interface, reflect.Ptr:
		return !v.IsNil()
	}
	return v.IsValid() && v.Interface() != reflect.Zero(v.Type()).Interface()
}

func isAlpha(f Field, _ ...string) bool {
	return alphaRegex.MatchString(f.Value())
}

func isAlphaNum(f Field, _ ...string) bool {
	return alphaNumericRegex.MatchString(f.Value())
}

func isAlphaUnicode(f Field, _ ...string) bool {
	return alphaUnicodeRegex.MatchString(f.Value())
}

func isAlphaNumUnicode(f Field, _ ...string) bool {
	return alphaNumericUnicodeRegex.MatchString(f.Value())
}

func isEmail(f Field, _ ...string) bool {
	return emailRegex.MatchString(f.Value())
}

func isURL(f Field, _ ...string) bool {
	u, err := url.ParseRequestURI(f.Value())
	if err != nil || len(u.Scheme) == 0 {
		return false
	}
	return true
}

func isURI(f Field, _ ...string) bool {
	_, err := url.ParseRequestURI(f.Value())
	return err == nil
}

func isNumeric(f Field, _ ...string) bool {
	return numericRegex.MatchString(f.Value())
}

func isNumber(f Field, _ ...string) bool {
	return numberRegex.MatchString(f.Value())
}

func isUUID(f Field, _ ...string) bool {
	return uuidRegex.MatchString(f.Value())
}

func isUUID3(f Field, _ ...string) bool {
	return uuid3Regex.MatchString(f.Value())
}

func isUUID4(f Field, _ ...string) bool {
	return uuid4Regex.MatchString(f.Value())
}
func isUUID5(f Field, _ ...string) bool {
	return uuid5Regex.MatchString(f.Value())
}
