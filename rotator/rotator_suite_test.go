package rotator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRotator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UAARotator Suite")
}
