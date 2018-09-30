package validator

import (
	"fmt"
	"strings"
)

type (
	// Tag is a validation tag.
	Tag struct {
		// Name is a tag name.
		Name string

		// Params is a parameter.
		Params []string

		// Optional is a flag. if true, empty value is always valid.
		Optional bool

		// Enable is a flag. if true, validate it.
		Enable bool

		// dig is a flag. if true, validate in map, slice, ptr and interface.
		dig bool

		// validateFn is a validator function.
		validateFn Func
	}
)

func (t Tag) String() string {
	if len(t.Params) > 0 {
		return fmt.Sprintf("%s%s%s%s", t.Name, "(", strings.Join(t.Params, "|"), ")")
	}
	return t.Name
}

func (t Tag) IsDig() bool {
	return t.dig
}
