package vertex

import (
	"sync"

	"github.com/egoholic/migrator/migration/parser/vertex/parseable"
)

const STRING_VERTEX_NAME = "<STRING>"

type (
	ParsingResult struct {
		Vertex *Vertex
		Parsed []rune
	}
	ParserFn            func(wg sync.WaitGroup, in <-chan rune, stopSig <-chan bool, out chan<- ParsingResult)
	ParserFnConstructor func(*Vertex, *parseable.Parseable) ParserFn
	Vertex              struct {
		Name        string
		Parseable   *parseable.Parseable
		Constructor ParserFnConstructor
		Edges       map[string]*Vertex
	}
)

// Vertex
// New() is too abstract. Prefer to use either NewToken() or NewString()
func New(name string, prl *parseable.Parseable, constructor ParserFnConstructor, vertices []*Vertex) *Vertex {
	v := &Vertex{
		Name:      name,
		Parseable: prl,
		Constructor: constructor,
	}
	v.MakeEdgesTo(vertices)
	return v
}

func NewToken(token string, vertices ...*Vertex) *Vertex {
	prl := parseable.NewParseableToken([]rune(token))
	return New(token, prl, NewParserFn, vertices)
}
func NewString(vertices ...*Vertex) *Vertex {
  prl := parseable.NewParseableString()
	return New(STRING_VERTEX_NAME, prl, NewParserFn, vertices)
}
func (v *Vertex) MakeEdgesTo(vertices ...*Vertex) {
	for _, anotherV := range vertices {
		v.Edges[anotherV.Name] = anotherV
	}
}
func(v *Vertex) ParserFn() ParserFn {
  return v.Constructor(v, v.Parseable)
}
// \Vertex

func NewParsingResult(v *Vertex) *ParsingResult {
	return &ParsingResult{
		Vertex: v,
		Parsed: nil,
	}
}

func NewParserFn(v *Vertex, prl parseable.Parseable) ParserFn {
	result := NewParsingResult(v)
	return
	func(wg sync.WaitGroup, in <-chan rune, stopSig <-chan bool, out chan<- ParsingResult) {
		defer wg.Done()
		for {
			select {
			case r := <- in:
				isMatches, isDone := prl.Match(r)
				if !isMatches {

				}

				if isDone {
					out <- NewParsingResult(v, prl.Parsed)
					return
				}
			case <- stopSig:
        return
		}
	}
}
