package parser

import (
	"sync"

	"github.com/egoholic/migrator/migration/parser/riter"
	"github.com/egoholic/migrator/migration/parser/vertex"
)

type Parser struct {
	graph   *vertex.Vertex
	parent  *vertex.Vertex
	in      chan rune
	stopSig chan bool
	out     chan []rune
	wg      sync.WaitGroup
	result  [][]rune
}

func New(graph *vertex.Vertex) *Parser {
	return &Parser{
		graph:   graph,
		parent:  graph,
		in:      make(chan rune),
		stopSig: make(chan bool),
		out:     make(chan []rune),
		result:  [][]rune{},
	}
}
func (p *Parser) Parse(raw []rune) (result [][]rune, err error) {
	iter := riter.New(raw)
	for _, v := range p.parent.Edges {
		p.wg.Add(1)
		go v.Pattern.Parse(p.wg, p.in, p.stopSig, p.out)
	}
	p.wg.Add(1)
	go p.parse(iter)
	p.wg.Wait()
	p.stopSig <- true
	return
}
func (p *Parser) parse(iter *riter.Iterator) {
	defer p.wg.Done()
	for {
		select {
		case parsed := <-p.out:
			p.stopSig <- true
			p.result = append(p.result, parsed)
			// p.parent =
		default:
			if iter.HasNext() {
				b, _ := iter.Next()
				p.in <- b
			} else {
				return
			}
		}
	}
}
