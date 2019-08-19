package pattern

type (
	Pattern struct {
		Name  string
		Parse ParserFn
	}
	ParserFn func(in <-chan byte, stopSig <-chan bool, out chan<- []byte)
)

func New(name string, fn ParserFn) *Pattern {
	return &Pattern{
		Name:  name,
		Parse: fn,
	}
}
