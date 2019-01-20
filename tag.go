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

		// Params is a tag parameter.
		Params []string

		// validateFn is a validate function.
		validateFn Func
	}

	tagChunk struct {
		// Tags is a list of tag.
		Tags []Tag

		// Optional is a flag. if true, empty value is always valid.
		Optional bool

		Next *tagChunk
	}
)

// String returns tag string.
func (t Tag) String() string {
	if len(t.Params) > 0 {
		return fmt.Sprintf("%s%s%s%s", t.Name, "(", strings.Join(t.Params, "|"), ")")
	}
	return t.Name
}
