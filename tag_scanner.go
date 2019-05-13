package validator

type (
	tagScanner struct {
		buf string
		pos int
	}

	tagParamsScanner struct {
		buf string
		pos int
	}

	tagToken int
)

const (
	illegal tagToken = iota
	eof
	tagSeparator
	orSeparator
	nextSeparator
)

func newTagScanner(buf string) *tagScanner {
	return &tagScanner{buf: buf}
}

func (s *tagScanner) Scan() (tagToken, string) {
	var (
		lit        []byte
		depthParen int
	)
	for {
		if s.eof() {
			break
		}

		ch := s.read()
		switch ch {
		case ' ', '\t', '\r', '\n':
			if depthParen == 0 {
				continue
			}

		case ',':
			if depthParen == 0 {
				return tagSeparator, string(lit)
			}

		case '|':
			if depthParen == 0 {
				return orSeparator, string(lit)
			}

		case ';':
			if depthParen == 0 {
				return nextSeparator, string(lit)
			}

		case '(':
			depthParen++

		case ')':
			depthParen--
		}

		lit = append(lit, ch)
	}

	return eof, string(lit)
}

func (s *tagScanner) read() (ch byte) {
	ch = s.buf[s.pos]
	s.pos++
	return
}

func (s *tagScanner) eof() bool {
	return len(s.buf) == s.pos
}

func newTagParamsScanner(buf string) *tagParamsScanner {
	return &tagParamsScanner{buf: buf}
}

func (s *tagParamsScanner) Scan() (tagToken, string) {
	var lit []byte
	for {
		if s.eof() {
			break
		}

		ch := s.read()
		switch ch {
		case '|':
			return orSeparator, string(lit)

		case '\\':
			switch s.read() {
			case '|':
				ch = '|'
			default:
				s.unread()
			}
		}

		lit = append(lit, ch)
	}

	return eof, string(lit)
}

func (s *tagParamsScanner) read() (ch byte) {
	ch = s.buf[s.pos]
	s.pos++
	return
}

func (s *tagParamsScanner) unread() {
	if s.pos > 0 {
		s.pos--
	}
}

func (s *tagParamsScanner) eof() bool {
	return len(s.buf) == s.pos
}
