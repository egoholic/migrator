package tokenizer

type (
	Pattern struct {
		Name     string
		parserFn ParserFn
	}
	Combinator struct {
		Combinators []*Combinator
		Pattern     *Pattern
	}
	Token struct {
		Title string
		Body  []byte
	}
	Tokens   []Token
	ParserFn func(in <-chan byte, stopSig <-chan bool, out chan<- []byte)
)

var (
	upPattern = &Pattern{
		Name:     "UP",
		parserFn: GenParser([]byte("-- <UP>\n")),
	}

	downPattern = &Pattern{
		Name:     "DOWN",
		parserFn: GenParser([]byte("-- <DOWN>\n")),
	}

	andPattern = &Pattern{
		Name:     "AND",
		parserFn: GenParser([]byte("-- <AND>\n")),
	}

	endPattern = &Pattern{
		Name:     "END",
		parserFn: GenParser([]byte("-- <END>\n")),
	}

	instructionToken = &Pattern{
		Name: "INSTRUCTION",
		parserFn: func(in <-chan byte, stopSig <-chan bool, out chan<- []byte) {
			var psd = make([]byte, 512)
			for {
				select {
				case b := <-in:
					psd = append(psd, b)
				case <-stopSig:
					out <- psd
					return
				}
			}
		},
	}
)

func NewCombinator(combinators) *Combinator {
	return &Combinator{}
}

func (c *Combinator) IsRoot() bool {
	return c.Pattern == nil
}

func GenParser(ptrn []byte) ParserFn {
	return func(in <-chan byte, stopSig <-chan bool, out chan<- []byte) {
		var (
			l    = len(ptrn)
			last = ptrn[l-1]
			cur  = 0
			psd  = make([]byte, l)
		)
		for {
			select {
			case b := <-in:
				if b != ptrn[cur] {
					return
				}
				psd = append(psd, b)
				cur++
				if b == last {
					out <- psd
					return
				}
			case <-stopSig:
				return
			}
		}
	}
}
