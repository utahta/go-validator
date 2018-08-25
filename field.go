package validator

import "strings"

type (
	vField struct {
		name   string
		parent *vField
	}
)

const (
	fieldNameDelim = "."
)

func (f vField) Name() string {
	var s []string
	if f.parent != nil {
		s = append(s, f.parent.Name())
	}

	s = append(s, f.name)
	return strings.Join(s, fieldNameDelim)
}
