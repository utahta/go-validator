package validator

type (
	tagScanner struct {
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
		lit        string
		depthParen int
	)
	for {
		if s.eof() {
			break
		}

		ch := s.read()
		switch ch {
		case ' ', '\t', '\r', '\n':
			continue

		case ',':
			if depthParen == 0 {
				return tagSeparator, lit
			}

		case '|':
			if depthParen == 0 {
				return orSeparator, lit
			}

		case ';':
			if depthParen == 0 {
				return nextSeparator, lit
			}

		case '(':
			depthParen++

		case ')':
			depthParen--
		}

		lit += string(ch)
	}

	return eof, lit
}

func (s *tagScanner) read() (ch byte) {
	ch = s.buf[s.pos]
	s.pos++
	return
}

func (s *tagScanner) eof() bool {
	return len(s.buf) == s.pos
}
