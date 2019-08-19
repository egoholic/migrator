package riter_test

import (
	. "github.com/egoholic/migrator/migration/parser/riter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("riter", func() {
	Describe("*Iterator", func() {
		Describe(".Next()", func() {
			Context("when has next rune", func() {
				It("returns next rune", func() {
					runes := []rune("test")
					ri := New(runes)
					Expect(ri.Next()).To(Equal(runes[0]))
					Expect(ri.Next()).To(Equal(runes[1]))
					Expect(ri.Next()).To(Equal(runes[2]))
					Expect(ri.Next()).To(Equal(runes[3]))
					_, err := ri.Next()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("there is no more runes"))
				})
			})
			Context("when has no next rune", func() {
				It("returns next rune", func() {
					runes := []rune("test")
					ri := New(runes)
					_, err := ri.Next()
					_, err = ri.Next()
					_, err = ri.Next()
					_, err = ri.Next()
					b, err := ri.Next()
					Expect(b).To(Equal(int32(0)))
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("there is no more runes"))
				})
			})
		})
		Describe(".HasNext()", func() {
			Context("when has next", func() {
				It("returns true", func() {
					runes := []rune("test")
					ri := New(runes)
					Expect(ri.HasNext()).To(BeTrue())
					_, _ = ri.Next()
					Expect(ri.HasNext()).To(BeTrue())
					_, _ = ri.Next()
					Expect(ri.HasNext()).To(BeTrue())
					_, _ = ri.Next()
					Expect(ri.HasNext()).To(BeTrue())
				})
			})
			Context("when has no next", func() {
				It("returns true", func() {
					runes := []rune("test")
					ri := New(runes)
					_, _ = ri.Next()
					_, _ = ri.Next()
					_, _ = ri.Next()
					_, _ = ri.Next()
					Expect(ri.HasNext()).To(BeFalse())
				})
			})
		})
	})
})
