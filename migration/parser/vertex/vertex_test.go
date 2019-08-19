package vertex_test

import (
	. "github.com/egoholic/migrator/migration/parser/vertex"
	"github.com/egoholic/migrator/migration/parser/vertex/pattern"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func ExampleParserFn(in <-chan rune, stopSig <-chan bool, out chan<- []rune) {

}

var _ = Describe("vertex", func() {
	Describe("*Vertex", func() {
		Describe(".AddEdgesTo()", func() {
			It("adds edges", func() {
				v1 := New(pattern.New("-- P1", ExampleParserFn))
				v2 := New(pattern.New("-- P2", ExampleParserFn))
				Expect(v1.Edges).To(BeEmpty())
				v1.AddEdgesTo(v2)
				Expect(v1.Edges).To(ContainElement(v2))
			})
		})
		Describe(".IsToken()", func() {
			Context("when token", func() {
				It("returns true", func() {
				})
			})
			Context("when data", func() {
				It("returns false", func() {
				})
			})
		})
		Describe(".Pattern", func() {
			It("returns pattern", func() {

			})
		})
		Describe(".Edges", func() {
			It("returns edges", func() {

			})
		})
	})
})
