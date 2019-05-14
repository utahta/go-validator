package validator

import (
	"fmt"
	"strings"
)

type (
	// Tag represents the validation tag in struct field's tag.
	Tag struct {
		// Name is a tag name.
		// e.g. len(1|2) -> "len"
		name string

		// Params is a tag parameters.
		// e.g. len(1|2) -> []string{"1", "2"}
		params []string

		// validateFn is a validate function.
		validateFn Func
	}

	tagChunk struct {
		// Tags is a list of Tag.
		Tags []Tag

		// Optional is a flag. If true, the empty value is always valid.
		Optional bool

		Next *tagChunk
	}
)

// Fullname returns a tag value.
func (t Tag) Fullname() string {
	if len(t.params) > 0 {
		return fmt.Sprintf("%s%s%s%s", t.name, "(", strings.Join(t.params, "|"), ")")
	}
	return t.name
}

// String returns a tag value.
func (t Tag) String() string {
	return t.Fullname()
}

func (c *tagChunk) GetTags() []Tag {
	if c == nil {
		return nil
	}
	return c.Tags
}

func (c *tagChunk) IsOptional() bool {
	if c == nil {
		return false
	}
	return c.Optional
}
