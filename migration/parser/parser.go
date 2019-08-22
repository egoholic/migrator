package parser

import (
	"sync"

	"github.com/egoholic/migrator/migration/parser/streamer"
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
		result:   [][]*vertex.ParsingResult,
	}
}
func (p *Parser) Parse(raw []rune) ([][]rune, error) {
	streamer = streamer.New(raw)
	pf = p.entrance.ParserFn()
	p.wg.Add(1)
	go pf(p.wg, p.in, p.stopSig, p.out)
	streamer.Stream(p.in)
	p.wg.Wait()
	return p.result, nil
}
