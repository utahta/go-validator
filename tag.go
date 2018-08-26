package validator

import (
	"fmt"
	"strings"
)

type (
	// Tag is a tag that parsed tag string.
	Tag struct {
		Name   string
		Params []string
	}
)

const (
	tagSeparator      = ","
	tagParamLParen    = "("
	tagParamRParen    = ")"
	tagParamSeparator = "|"
)

func (t Tag) String() string {
	if len(t.Params) > 0 {
		return fmt.Sprintf("%s%s%s%s", t.Name, tagParamLParen, strings.Join(t.Params, tagParamSeparator), tagParamRParen)
	}
	return t.Name
}

func parseTag(rawTag string) ([]Tag, error) {
	var tags []Tag
	for _, t := range strings.Split(rawTag, tagSeparator) {
		nt, err := newTag(t)
		if err != nil {
			return nil, err
		}
		tags = append(tags, nt)
	}
	return tags, nil
}

func newTag(rawTag string) (Tag, error) {
	rawTag = strings.TrimSpace(rawTag)
	if rawTag == "" {
		//TODO: define error
		return Tag{}, fmt.Errorf("newTag: tag name is empty")
	}

	var (
		name   = rawTag
		params []string
	)
	if il := strings.Index(rawTag, tagParamLParen); il >= 0 {
		ir := strings.LastIndex(rawTag, tagParamRParen)
		if len(rawTag)-1 != ir {
			//TODO: define error
			return Tag{}, fmt.Errorf("newTag: right paren not found")
		}

		name = name[:il]
		params = strings.Split(rawTag[il:ir], tagParamSeparator)
	}

	return Tag{
		Name:   name,
		Params: params,
	}, nil
}
