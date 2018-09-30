package validator

import (
	"fmt"
	"strings"
)

func (v *Validator) tagParse(rawTag string) ([]Tag, error) {
	if tags, ok := v.tagCache.Load(rawTag); ok {
		return tags, nil
	}

	var (
		tags       []Tag
		optional   = false
		orParsing  = false
		digParsing = false
	)
	const optionalTagName = "optional"

	s := newTagScanner(rawTag)
loop:
	for {
		token, lit := s.Scan()
		if lit == optionalTagName {
			for i := range tags {
				if digParsing {
					if tags[i].isDig {
						tags[i].Optional = true
					}
				} else {
					tags[i].Optional = true
				}
			}
			optional = true
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
				tags[len(tags)-1].Params = append(tags[len(tags)-1].Params, lit)
			} else {
				tag, err := v.newTag(lit, !digParsing, optional)
				if err != nil {
					return nil, err
				}
				tags = append(tags, tag)
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
				tags[len(tags)-1].Params = append(tags[len(tags)-1].Params, lit)
			} else {
				tag, err := v.newTag(lit, !digParsing, optional)
				if err != nil {
					return nil, err
				}
				tags = append(tags, tag)
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
				tags[len(tags)-1].Params = append(tags[len(tags)-1].Params, lit)
			} else {
				tags = append(tags, Tag{Name: "or", Params: []string{lit}, Optional: optional, Enable: !digParsing, isDig: true, validateFn: v.FuncMap["or"]})
			}
			orParsing = true

		case digSeparator:
			if lit != "" && lit != optionalTagName {
				if orParsing {
					tags[len(tags)-1].Params = append(tags[len(tags)-1].Params, lit)
				} else {
					tag, err := v.newTag(lit, true, optional)
					if err != nil {
						return nil, err
					}
					tags = append(tags, tag)
				}
			}
			for i := range tags {
				tags[i].isDig = false
			}
			digParsing = true
			orParsing = false
			optional = false

		case illegal:
			return nil, fmt.Errorf("parse: illegal token")
		}
	}

	v.tagCache.Store(rawTag, tags)

	return tags, nil
}

// newTag returns Tag.
func (v *Validator) newTag(lit string, enable, optional bool) (Tag, error) {
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
		Optional:   optional,
		Enable:     enable,
		isDig:      true,
		validateFn: fn,
	}, nil
}
