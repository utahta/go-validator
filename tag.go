package validator

import "strings"

type (
	vTag struct {
		Name string
	}
)

const (
	tagSeparator      = ","
	tagParamSeparator = "|"
)

func parseTag(tag string) []vTag {
	var tags []vTag
	for _, t := range strings.Split(tag, tagSeparator) {
		tags = append(tags, vTag{
			Name: t,
		})
	}
	return tags
}
