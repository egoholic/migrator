package parseable

type (
	ParsedRunes struct {
		initSize int
		runes    []rune
	}
	Parseable interface {
		Match(rune) (bool, bool)
		Parsed() []rune
		Reset()
	}
	ParseableString struct {
		*ParsedRunes
	}
	ParseableToken struct {
		token []rune
		*ParsedRunes
	}
)

// ParsedRunes
func NewParsedRunes(initSize int) *ParsedRunes {
	return &ParsedRunes{
		initSize: initSize,
		runes:    make([]rune, initSize),
	}
}
func (p *ParsedRunes) Reset() {
	p.runes = make([]rune, p.initSize)
}
func (p *ParsedRunes) add(r rune) {
	p.runes = append(p.runes, r)
}
func (p *ParsedRunes) len() int {
	return len(p.runes)
}

// \ParsedRunes

// ParseableString
func NewParseableString() *ParseableString {
	return &ParseableString{NewParsedRunes(512)}
}
func (s *ParseableString) Match(r rune) bool {
	s.add(r)
	return true
}
func (s *ParseableString) Parsed() []rune {
	return s.runes
}

// \ ParseableString

// ParseableToken
func NewParseableToken(token []rune) *ParseableToken {
	return &ParseableToken{
		token:       token,
		ParsedRunes: NewParsedRunes(len(token)),
	}
}
func (t *ParseableToken) Match(r rune) bool {
	if t.token[t.ParsedRunes.len()] != r {
		return false
	}
	t.add(r)
	return true
}
func (t *ParseableToken) Parsed() []rune {
	return t.runes
}

// \ParseableToken
