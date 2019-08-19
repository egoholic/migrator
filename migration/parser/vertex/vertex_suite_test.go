package vertex_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVertex(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Vertex Suite")
}
