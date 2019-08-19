package pattern

import "sync"

type (
	Pattern struct {
		Name  string
		Parse ParserFn
	}
	ParserFn func(wg sync.WaitGroup, in <-chan rune, stopSig <-chan bool, out chan<- []rune)
)

func New(name string, fn ParserFn) *Pattern {
	return &Pattern{
		Name:  name,
		Parse: fn,
	}
}
func NewTokenParserFn(token []rune) ParserFn {
	return func(wg sync.WaitGroup, in <-chan rune, stopSig <-chan bool, out chan<- []rune) {
		defer wg.Done()
		var (
			l    = len(token)
			last = token[l-1]
			cur  = 0
			psd  = make([]rune, l)
		)
		for {
			select {
			case b := <-in:
				if b != token[cur] {
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
func DataParserFn(wg sync.WaitGroup, in <-chan rune, stopSig <-chan bool, out chan<- []rune) {
	defer wg.Done()
	var psd = make([]rune, 512)
	for {
		select {
		case b := <-in:
			psd = append(psd, b)
		case <-stopSig:
			out <- psd
			return
		}
	}
}
