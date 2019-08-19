package parser

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
		graph *Vertex
	}
)

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

func New(graph *Vertex) *Parser {
	return &Parser{graph}
}

func (p *Parser) Parse(raw []byte) {
	in := make(chan byte)
	stopSig := make(chan bool)
	out := make(chan []byte)
	for _, b := range raw {
		in <- b
	}
	stopSig <- true
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
