package validator

import (
	"fmt"
	"strings"
)

func (v *Validator) tagParse(rawTag string) (*tagChunk, error) {
	if tags, ok := v.tagCache.Load(rawTag); ok {
		return tags, nil
	}

	var (
		rootTagChunk tagChunk
		chunk        *tagChunk
		orParsing    = false
	)
	const optionalTagName = "optional"

	chunk = &rootTagChunk

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
				chunk.Tags[len(chunk.Tags)-1].Params = append(chunk.Tags[len(chunk.Tags)-1].Params, lit)
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
				chunk.Tags[len(chunk.Tags)-1].Params = append(chunk.Tags[len(chunk.Tags)-1].Params, lit)
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
				chunk.Tags[len(chunk.Tags)-1].Params = append(chunk.Tags[len(chunk.Tags)-1].Params, lit)
			} else {
				chunk.Tags = append(chunk.Tags, Tag{Name: "or", Params: []string{lit}, validateFn: v.FuncMap["or"]})
			}
			orParsing = true

		case nextSeparator:
			if lit != "" && lit != optionalTagName {
				if orParsing {
					chunk.Tags[len(chunk.Tags)-1].Params = append(chunk.Tags[len(chunk.Tags)-1].Params, lit)
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

	if rootTagChunk.Next == nil {
		rootTagChunk.Next = &tagChunk{
			Tags:     rootTagChunk.Tags[:],
			Optional: rootTagChunk.Optional,
		}
	}
	v.tagCache.Store(rawTag, &rootTagChunk)

	return &rootTagChunk, nil
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
		s := newTagScanner(lit[idx+1 : len(lit)-1])
	loop:
		for {
			token, lit := s.Scan()
			switch token {
			case eof:
				params = append(params, lit)
				break loop

			case orSeparator:
				params = append(params, lit)

			default:
				return Tag{}, fmt.Errorf("parse: failed to new tag")
			}
		}
	}

	fn, ok := v.FuncMap[name]
	if !ok {
		return Tag{}, fmt.Errorf("parse: tag %s function not found", name)
	}

	return Tag{
		Name:       name,
		Params:     params,
		validateFn: fn,
	}, nil
}
