package riter_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBiter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Riter Suite")
}
