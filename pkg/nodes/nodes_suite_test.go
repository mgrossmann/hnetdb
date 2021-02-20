package nodes_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNodes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nodes Suite")
}
