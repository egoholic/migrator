package vertex

import "github.com/egoholic/migrator/migration/parser/vertex/pattern"

type Vertex struct {
	Pattern *pattern.Pattern
	Edges   []*Vertex
}

func New(pattern *pattern.Pattern, edges ...*Vertex) *Vertex {
	return &Vertex{
		Pattern: pattern,
		Edges:   edges,
	}
}
func (v *Vertex) AddEdgesTo(vertices ...*Vertex) {
	for _, newV := range vertices {
		v.Edges = append(v.Edges, newV)
	}
}
func (v *Vertex) IsToken() bool {
	return true
}
