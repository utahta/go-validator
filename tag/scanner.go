package tag

type (
	scanner struct {
		buf string
		pos int
	}

	token int
)

const (
	illegal token = iota
	eof
	tagSeparator
	orSeparator
	digSeparator
	paramLParen
	paramRParen
)

func newScanner(buf string) *scanner {
	return &scanner{buf: buf}
}

func (s *scanner) Scan() (token, string) {
	var (
		lit        string
		depthParen int
	)
	for {
		if s.eof() {
			return eof, lit
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
			return digSeparator, lit

		case '(':
			depthParen++

		case ')':
			depthParen--
		}

		lit += string(ch)
	}

	return illegal, ""
}

func (s *scanner) read() (ch byte) {
	ch = s.buf[s.pos]
	s.pos++
	return
}

func (s *scanner) eof() bool {
	return len(s.buf) == s.pos
}
