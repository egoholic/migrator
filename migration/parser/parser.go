package parser

import "errors"

type (
	Result struct {
		Up   []string
		Down []string
	}
	Vertex struct {
		Pattern *Pattern
		Edges   []*Vertex
	}
	Pattern struct {
		Name     string
		parserFn ParserFn
	}
	ParserFn func(in <-chan byte, stopSig <-chan bool, out chan<- []byte)
	Parser   struct {
		graph   *Vertex
		parent  *Vertex
		in      chan byte
		stopSig chan bool
		out     chan []byte
	}

	Iterator struct {
		bytes []byte
		len   int
		cur   int
	}
)

func (i Iterator) HasNext() bool {
	return i.len > i.cur+1
}

func (i Iterator) Next() (b byte, err error) {
	if i.HasNext() {
		b = i.bytes[i.cur]
		i.cur++
		return b, nil
	}
	return b, errors.New("")
}

func NewIterator(bytes []byte) *Iterator {
	return &Iterator{
		bytes: bytes,
		len:   len(bytes),
		cur:   0,
	}
}

func NewVertex(pattern *Pattern, edges ...*Vertex) *Vertex {
	return &Vertex{
		Pattern: pattern,
		Edges:   edges,
	}
}

func (v *Vertex) AddTransitionsTo(vertices ...*Vertex) {
	for _, v := range vertices {
		v.Edges = append(v.Edges, v)
	}
}

func (v *Vertex) IsToken() bool {
	return true
}

func New(graph *Vertex) *Parser {
	return &Parser{
		graph:   graph,
		parent:  graph,
		in:      make(chan byte),
		stopSig: make(chan bool),
		out:     make(chan []byte),
	}
}

func (p *Parser) Parse(raw []byte) (result [][]byte, err error) {
	iter := NewIterator(raw)
	for _, v := range p.parent.Edges {
		go v.Pattern.parserFn(p.in, p.stopSig, p.out)
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
	upPattern = &Pattern{
		Name:     "UP",
		parserFn: GenParserFn([]byte("-- <UP>\n")),
	}

	downPattern = &Pattern{
		Name:     "DOWN",
		parserFn: GenParserFn([]byte("-- <DOWN>\n")),
	}

	andPattern = &Pattern{
		Name:     "AND",
		parserFn: GenParserFn([]byte("-- <AND>\n")),
	}

	endPattern = &Pattern{
		Name:     "END",
		parserFn: GenParserFn([]byte("-- <END>\n")),
	}

	instructionPattern = &Pattern{
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

	up        = NewVertex(upPattern)
	upQuery   = NewVertex(instructionPattern)
	upAnd     = NewVertex(andPattern)
	down      = NewVertex(downPattern)
	downQuery = NewVertex(instructionPattern)
	downAnd   = NewVertex(andPattern)
	end       = NewVertex(endPattern)

	parsingGraph *Vertex
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

func GenParserFn(ptrn []byte) ParserFn {
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
