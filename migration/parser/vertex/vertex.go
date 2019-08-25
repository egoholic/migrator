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

	StringVertex struct {
		Name        string
		Parseable   *parseable.Parseable
		Constructor ParserFnConstructor
		Edges       map[string]*Vertex
	}

	TokenVertex struct {
		Name        string
		Parseable   *parseable.Parseable
		Constructor ParserFnConstructor
		Edges       map[string]*Vertex
	}
)

func NewToken(token string, vertices ...*Vertex) *TokenVertex {
	prl := parseable.NewParseableToken([]rune(token))
	v := &TokenVertex{
		Name:      token,
		Parseable: prl,
		Constructor: NewParserFn,
	}
	v.MakeEdgesTo(vertices)
	return v
}
func NewString(vertices ...*Vertex) *StringVertex {
	prl := parseable.NewParseableString()
	v := &StringVertex{
		Name:      STRING_VERTEX_NAME,
		Parseable: prl,
		Constructor: NewParserFn,
	}
	v.MakeEdgesTo(vertices)
	return v

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
