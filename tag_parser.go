package validator

import (
	"fmt"
	"strings"
)

func (v *Validator) parseTag(rawTag string) (*tagChunk, error) {
	if tags, ok := v.tagCache.Load(rawTag); ok {
		return tags, nil
	}

	var (
		rootChunk tagChunk
		chunk     *tagChunk
		orParsing = false
	)
	const optionalTagName = "optional"

	chunk = &rootChunk

	s := newTagScanner(rawTag)
loop:
	for {
		token, lit := s.Scan()
		if lit == optionalTagName {
			chunk.Optional = true
		}

		switch token {
		case eof:
			if lit == optionalTagName {
				continue
			}
			if lit == "" {
				break loop
			}

			if orParsing {
				idx := len(chunk.Tags) - 1
				chunk.Tags[idx].params = append(chunk.Tags[idx].params, lit)
			} else {
				tag, err := v.newTag(lit)
				if err != nil {
					return nil, err
				}
				chunk.Tags = append(chunk.Tags, tag)
			}
			break loop

		case tagSeparator:
			if lit == optionalTagName {
				continue
			}
			if lit == "" {
				return nil, fmt.Errorf("parse: invalid literal in tag separator")
			}

			if orParsing {
				idx := len(chunk.Tags) - 1
				chunk.Tags[idx].params = append(chunk.Tags[idx].params, lit)
			} else {
				tag, err := v.newTag(lit)
				if err != nil {
					return nil, err
				}
				chunk.Tags = append(chunk.Tags, tag)
			}
			orParsing = false

		case orSeparator:
			if lit == optionalTagName {
				continue
			}
			if lit == "" {
				return nil, fmt.Errorf("parse: invalid literal in or separator")
			}

			if orParsing {
				idx := len(chunk.Tags) - 1
				chunk.Tags[idx].params = append(chunk.Tags[idx].params, lit)
			} else {
				chunk.Tags = append(chunk.Tags, Tag{name: "or", params: []string{lit}, validateFn: v.funcMap["or"]})
			}
			orParsing = true

		case nextSeparator:
			if lit != "" && lit != optionalTagName {
				if orParsing {
					idx := len(chunk.Tags) - 1
					chunk.Tags[idx].params = append(chunk.Tags[idx].params, lit)
				} else {
					tag, err := v.newTag(lit)
					if err != nil {
						return nil, err
					}
					chunk.Tags = append(chunk.Tags, tag)
				}
			}
			orParsing = false
			chunk.Next = &tagChunk{}
			chunk = chunk.Next
		}
	}

	v.tagCache.Store(rawTag, &rootChunk)

	return &rootChunk, nil
}

// newTag returns Tag.
func (v *Validator) newTag(lit string) (Tag, error) {
	var (
		name   string
		params []string
	)

	idx := strings.Index(lit, "(")
	if idx < 0 {
		name = lit
	} else {
		name = lit[:idx]
		s := newTagParamsScanner(lit[idx+1 : len(lit)-1])
	loop:
		for {
			token, lit := s.Scan()
			switch token {
			case eof:
				params = append(params, lit)
				break loop

			case orSeparator:
				params = append(params, lit)
			}
		}
	}

	fn, ok := v.funcMap[name]
	if !ok {
		return Tag{}, fmt.Errorf("parse: tag %s function not found", name)
	}

	return Tag{
		name:       name,
		params:     params,
		validateFn: fn,
	}, nil
}
