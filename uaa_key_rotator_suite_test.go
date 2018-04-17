package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"os"
)

func TestUaaKeyRotator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UaaKeyRotator Suite")
}

var uaaRotatorBuildPath string

var _ = BeforeSuite(func() {
	var err error
	uaaRotatorBuildPath, err = gexec.Build("github.com/cloudfoundry/uaa-key-rotator")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	Expect(os.Remove(uaaRotatorBuildPath)).To(Succeed())
})