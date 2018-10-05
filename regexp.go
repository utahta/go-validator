package validator

import "regexp"

const (
	alphaRegexString               = "^[a-zA-Z]+$"
	alphaNumericRegexString        = "^[a-zA-Z0-9]+$"
	alphaUnicodeRegexString        = `^[\p{L}]+$`
	alphaNumericUnicodeRegexString = `^[\p{L}\p{N}]+$`
	numericRegexString             = `^[-+]?[0-9]+(?:\.[0-9]+)?$`
	numberRegexString              = "^[0-9]+$"
	hexadecimalRegexString         = "^[0-9a-fA-F]+$"
	hexcolorRegexString            = "^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{6})$"
	rgbRegexString                 = `^rgb\(\s*(?:(?:0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])\s*,\s*(?:0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])\s*,\s*(?:0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])|(?:0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])%\s*,\s*(?:0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])%\s*,\s*(?:0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])%)\s*\)$`
	rgbaRegexString                = `^rgba\(\s*(?:(?:0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])\s*,\s*(?:0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])\s*,\s*(?:0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])|(?:0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])%\s*,\s*(?:0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])%\s*,\s*(?:0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])%)\s*,\s*(?:(?:0.[1-9]*)|[01])\s*\)$`
	hslRegexString                 = `^hsl\(\s*(?:0|[1-9]\d?|[12]\d\d|3[0-5]\d|360)\s*,\s*(?:(?:0|[1-9]\d?|100)%)\s*,\s*(?:(?:0|[1-9]\d?|100)%)\s*\)$`
	hslaRegexString                = `^hsla\(\s*(?:0|[1-9]\d?|[12]\d\d|3[0-5]\d|360)\s*,\s*(?:(?:0|[1-9]\d?|100)%)\s*,\s*(?:(?:0|[1-9]\d?|100)%)\s*,\s*(?:(?:0.[1-9]*)|[01])\s*\)$`
	emailRegexString               = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:\\(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22)))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	base64RegexString              = `^(?:[A-Za-z0-9+\/]{4})*(?:[A-Za-z0-9+\/]{2}==|[A-Za-z0-9+\/]{3}=|[A-Za-z0-9+\/]{4})$`
	base64URLRegexString           = `^(?:[A-Za-z0-9-_]{4})*(?:[A-Za-z0-9-_]{2}==|[A-Za-z0-9-_]{3}=|[A-Za-z0-9-_]{4})$`
	isbn10RegexString              = "^(?:[0-9]{9}X|[0-9]{10})$"
	isbn13RegexString              = "^(?:(?:97(?:8|9))[0-9]{10})$"
	uuidRegexString                = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
	uuid3RegexString               = "^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"
	uuid4RegexString               = "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	uuid5RegexString               = "^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	asciiRegexString               = "^[\x00-\x7F]+$"
	printableASCIIRegexString      = "^[\x20-\x7E]+$"
	multibyteRegexString           = "^[^\x00-\x7F]+$"
	dataURIRegexString             = `^data:.+\/(.+);base64,(?:[A-Za-z0-9+\/]{4})*(?:[A-Za-z0-9+\/]{2}==|[A-Za-z0-9+\/]{3}=|[A-Za-z0-9+\/]{4})$`
	latitudeRegexString            = `^[-+]?([1-8]?\d(\.\d+)?|90(\.0+)?)$`
	longitudeRegexString           = `^[-+]?(180(\.0+)?|((1[0-7]\d)|([1-9]?\d))(\.\d+)?)$`
	ssnRegexString                 = `^\d{3}[- ]?\d{2}[- ]?\d{4}$`
	semverRegexString              = `^v?(?:0|[1-9]\d*)\.(?:0|[1-9]\d*)\.(?:0|[1-9]\d*)(-(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\+[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*)?$`
	katakanaRegexString            = `^[\p{Katakana}]+$`
	hiraganaRegexString            = `^[\p{Hiragana}]+$`
	fullWidthRegexString           = "[^\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]"
	halfWidthRegexString           = "[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]"
)

var (
	alphaRegex               = regexp.MustCompile(alphaRegexString)
	alphaNumericRegex        = regexp.MustCompile(alphaNumericRegexString)
	alphaUnicodeRegex        = regexp.MustCompile(alphaUnicodeRegexString)
	alphaNumericUnicodeRegex = regexp.MustCompile(alphaNumericUnicodeRegexString)
	numericRegex             = regexp.MustCompile(numericRegexString)
	numberRegex              = regexp.MustCompile(numberRegexString)
	hexadecimalRegex         = regexp.MustCompile(hexadecimalRegexString)
	hexcolorRegex            = regexp.MustCompile(hexcolorRegexString)
	rgbRegex                 = regexp.MustCompile(rgbRegexString)
	rgbaRegex                = regexp.MustCompile(rgbaRegexString)
	hslRegex                 = regexp.MustCompile(hslRegexString)
	hslaRegex                = regexp.MustCompile(hslaRegexString)
	emailRegex               = regexp.MustCompile(emailRegexString)
	base64Regex              = regexp.MustCompile(base64RegexString)
	base64URLRegex           = regexp.MustCompile(base64URLRegexString)
	isbn10Regex              = regexp.MustCompile(isbn10RegexString)
	isbn13Regex              = regexp.MustCompile(isbn13RegexString)
	uuidRegex                = regexp.MustCompile(uuidRegexString)
	uuid3Regex               = regexp.MustCompile(uuid3RegexString)
	uuid4Regex               = regexp.MustCompile(uuid4RegexString)
	uuid5Regex               = regexp.MustCompile(uuid5RegexString)
	asciiRegex               = regexp.MustCompile(asciiRegexString)
	printableASCIIRegex      = regexp.MustCompile(printableASCIIRegexString)
	multibyteRegex           = regexp.MustCompile(multibyteRegexString)
	dataURIRegex             = regexp.MustCompile(dataURIRegexString)
	latitudeRegex            = regexp.MustCompile(latitudeRegexString)
	longitudeRegex           = regexp.MustCompile(longitudeRegexString)
	ssnRegex                 = regexp.MustCompile(ssnRegexString)
	semverRegex              = regexp.MustCompile(semverRegexString)
	katakanaRegex            = regexp.MustCompile(katakanaRegexString)
	hiraganaRegex            = regexp.MustCompile(hiraganaRegexString)
	fullWidthRegex           = regexp.MustCompile(fullWidthRegexString)
	halfWidthRegex           = regexp.MustCompile(halfWidthRegexString)
)
