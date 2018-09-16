package validator

import (
	"fmt"
	"strings"
	"sync"
)

var (
	cache sync.Map
)

func tagParse(rawTag string) ([]Tag, error) {
	if v, ok := cache.Load(rawTag); ok {
		return v.([]Tag), nil
	}

	var (
		tags      []Tag
		enable    = true
		optional  = false
		orParsing = false
	)

	s := newTagScanner(rawTag)
loop:
	for {
		token, lit := s.Scan()
		if lit == "optional" {
			optional = true
			continue
		}

		switch token {
		case eof:
			if lit == "" {
				break loop
			}

			if orParsing {
				tags[len(tags)-1].Params = append(tags[len(tags)-1].Params, lit)
			} else {
				tag, err := newTag(lit, enable, true)
				if err != nil {
					return nil, err
				}
				tags = append(tags, tag)
			}
			break loop

		case tagSeparator:
			if lit == "" {
				return nil, fmt.Errorf("parse: invalid literal in tag separator")
			}

			if orParsing {
				tags[len(tags)-1].Params = append(tags[len(tags)-1].Params, lit)
			} else {
				tag, err := newTag(lit, enable, true)
				if err != nil {
					return nil, err
				}
				tags = append(tags, tag)
			}
			orParsing = false

		case orSeparator:
			if lit == "" {
				return nil, fmt.Errorf("parse: invalid literal in or separator")
			}

			if orParsing {
				tags[len(tags)-1].Params = append(tags[len(tags)-1].Params, lit)
			} else {
				tags = append(tags, Tag{Name: "or", Params: []string{lit}, Enable: enable, dig: true})
			}
			orParsing = true

		case digSeparator:
			for i := range tags {
				tags[i].dig = false
			}
			if lit != "" {
				if orParsing {
					tags[len(tags)-1].Params = append(tags[len(tags)-1].Params, lit)
				} else {
					tag, err := newTag(lit, enable, false)
					if err != nil {
						return nil, err
					}
					tags = append(tags, tag)
				}
			}
			enable = false
			orParsing = false

		case illegal:
			return nil, fmt.Errorf("parse: illegal token")
		}
	}

	if optional {
		for i := range tags {
			tags[i].Optional = optional
		}
	}

	cache.Store(rawTag, tags)

	return tags, nil
}

func newTag(lit string, enable, dig bool) (Tag, error) {
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

	return Tag{
		Name:   name,
		Params: params,
		Enable: enable,
		dig:    dig,
	}, nil
}
