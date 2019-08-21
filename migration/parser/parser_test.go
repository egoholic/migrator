package parser_test

import (
	. "github.com/egoholic/migrator/migration/parser"
	"github.com/egoholic/migrator/migration/parser/vertex"
	"github.com/egoholic/migrator/migration/parser/vertex/pattern"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	upPattern          = pattern.New("UP", true, pattern.NewTokenParserFn([]rune("-- <UP>\n")))
	downPattern        = pattern.New("DOWN", true, pattern.NewTokenParserFn([]rune("-- <DOWN>\n")))
	andPattern         = pattern.New("AND", true, pattern.NewTokenParserFn([]rune("-- <AND>\n")))
	endPattern         = pattern.New("END", true, pattern.NewTokenParserFn([]rune("-- <END>\n")))
	instructionPattern = pattern.New("<DATA>", true, pattern.DataParserFn)

	up        = vertex.New(upPattern)
	upQuery   = vertex.New(instructionPattern)
	upAnd     = vertex.New(andPattern)
	down      = vertex.New(downPattern)
	downQuery = vertex.New(instructionPattern)
	downAnd   = vertex.New(andPattern)
	end       = vertex.New(endPattern)

	ParsingGraph *vertex.Vertex
)

var _ = Describe("parser", func() {
	up.AddEdgesTo(upQuery)
	upAnd.AddEdgesTo(upQuery)
	upQuery.AddEdgesTo(upAnd, down)
	down.AddEdgesTo(downQuery)
	downAnd.AddEdgesTo(downQuery)
	downQuery.AddEdgesTo(downAnd, end)
	ParsingGraph = up

	Describe("*Parser", func() {
		var parser = New(ParsingGraph)
		Describe(".Parse()", func() {
			Context("when correct input", func() {
				It("returns AST", func() {
					result, err := parser.Parse([]rune(`-- <UP>
SELECT * FROM articles;
-- <AND>
SELECT * FROM publications;
-- <DOWN>
SELECT * FROM categories;
-- <AND>
SELECT * FROM comments;
-- <END>`))
					Expect(result).NotTo(BeNil())
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when incorrect input", func() {
				It("returns fails", func() {
					result, err := parser.Parse([]rune(`-- <UP>
SELECT * FROM articles;
-- <AND>
SELECT * FROM publications;
-- <DOWN>
SELECT * FROM categories;
-- <AND>
SELECT * FROM comments;
-- <END>`))
					Expect(result).To(BeNil())
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal(""))
				})
			})
		})
	})
})
