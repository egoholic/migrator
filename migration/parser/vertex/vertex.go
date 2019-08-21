package vertex

import (
	"sync"
)

const STRING_VERTEX_NAME = "<STRING>"

type (
	ParserFn func(wg sync.WaitGroup, in <-chan rune, stopSig <-chan bool, out chan<- []rune)

	Vertex struct {
		Name       string
		Parser     ParserFn
		IsLenKnown bool
		Edges      map[string]*Vertex
	}
)

// New() is too abstract. Prefer to use either NewToken() or NewString()
func New(name string, isLenKnown bool, parser ParserFn, transfersTo []*Vertex) *Vertex {
	v := &Vertex{
		Name:       name,
		Parser:     parser,
		IsLenKnown: isLenKnown,
	}
	for _, t := range transfersTo {
		v.Edges[t.Name] = t
	}
	return v
}

func NewToken(token string, edges ...*Vertex) *Vertex {
	return New(token, true, newTokenParserFn([]rune(token)), edges)
}
func NewString(edges ...*Vertex) *Vertex {
	return New(STRING_VERTEX_NAME, false, stringParserFn, edges)
}
func (v *Vertex) AddEdgesTo(transfersTo ...*Vertex) {
	for _, t := range transfersTo {
		v.Edges[t.Name] = t
	}
}

func newTokenParserFn(token []rune) func(sync.WaitGroup, <-chan rune, <-chan bool, chan<- []rune) {
	return func(wg sync.WaitGroup, in <-chan rune, stopSig <-chan bool, out chan<- []rune) {
		defer wg.Done()
		var (
			l    = len(token)
			last = token[l-1]
			cur  = 0
			psd  = []rune{}
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
func stringParserFn(wg sync.WaitGroup, in <-chan rune, stopSig <-chan bool, out chan<- []rune) {
	defer wg.Done()
	psd := make([]rune, 512)

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
