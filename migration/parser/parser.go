package parser

import (
	"sync"

	"github.com/egoholic/migrator/migration/parser/riter"
	"github.com/egoholic/migrator/migration/parser/vertex"
)

type (
	Parser struct {
		entrance *vertex.Vertex
		current  *vertex.Vertex
		in       chan rune
		stopSig  chan bool
		out      chan []rune
		wg       sync.WaitGroup
		result   [][]rune
	}
)

func New(graph *vertex.Vertex) *Parser {
	return &Parser{
		entrance: graph,
		current:  graph,
		in:       make(chan rune),
		stopSig:  make(chan bool),
		out:      make(chan []rune),
		result:   [][]rune{},
	}
}
func (p *Parser) Parse(raw []rune) (result []pattern.ParsedOut, err error) {
	iter := riter.New(raw)
	if p.entrance.IsLenKnown {
		p.wg.Add(1)
		go p.entrance.ParserFn(p.wg, p.in, p.stopSig, p.out)
	}
}
