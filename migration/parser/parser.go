package parser

import (
	"github.com/egoholic/migrator/migration/parser/biter"
	"github.com/egoholic/migrator/migration/parser/vertex"
	"github.com/egoholic/migrator/migration/parser/vertex/pattern"
)

type (
	Result struct {
		Up   []string
		Down []string
	}

	Parser struct {
		graph   *vertex.Vertex
		parent  *vertex.Vertex
		in      chan byte
		stopSig chan bool
		out     chan []byte
	}
)

func New(graph *vertex.Vertex) *Parser {
	return &Parser{
		graph:   graph,
		parent:  graph,
		in:      make(chan byte),
		stopSig: make(chan bool),
		out:     make(chan []byte),
	}
}

func (p *Parser) Parse(raw []byte) (result [][]byte, err error) {
	iter := biter.New(raw)
	for _, v := range p.parent.Edges {
		go v.Pattern.Parse(p.in, p.stopSig, p.out)
	}

	for {
		select {
		case parsed := <-p.out:
			p.stopSig <- true
			result = append(result, parsed)
		default:
			for b, err := iter.Next(); err == nil; {
				p.in <- b
			}
		}
	}

	p.stopSig <- true
	return
}

var (
	upPattern          = pattern.New("UP", GenParserFn([]byte("-- <UP>\n")))
	downPattern        = pattern.New("DOWN", GenParserFn([]byte("-- <DOWN>\n")))
	andPattern         = pattern.New("AND", GenParserFn([]byte("-- <AND>\n")))
	endPattern         = pattern.New("END", GenParserFn([]byte("-- <END>\n")))
	instructionPattern = pattern.New("INSTRUCTION", instructionParserFn)

	up        = vertex.New(upPattern)
	upQuery   = vertex.New(instructionPattern)
	upAnd     = vertex.New(andPattern)
	down      = vertex.New(downPattern)
	downQuery = vertex.New(instructionPattern)
	downAnd   = vertex.New(andPattern)
	end       = vertex.New(endPattern)

	parsingGraph *vertex.Vertex
)

func init() {
	up.AddTransitionsTo(upQuery)
	upAnd.AddTransitionsTo(upQuery)
	upQuery.AddTransitionsTo(upAnd, down)
	down.AddTransitionsTo(downQuery)
	downAnd.AddTransitionsTo(downQuery)
	downQuery.AddTransitionsTo(downAnd, end)
	parsingGraph = up
}

func instructionParserFn(in <-chan byte, stopSig <-chan bool, out chan<- []byte) {
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
}
func GenParserFn(ptrn []byte) pattern.ParserFn {
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
