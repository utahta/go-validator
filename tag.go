package validator

import (
	"fmt"
	"strings"
)

type (
	// Tag is a validation tag.
	Tag struct {
		// Name is a tag name.
		name string

		// Params is a tag parameter.
		params []string

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

func (t Tag) Fullname() string {
	if len(t.params) > 0 {
		return fmt.Sprintf("%s%s%s%s", t.name, "(", strings.Join(t.params, "|"), ")")
	}
	return t.name
}

// String returns tag string.
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
