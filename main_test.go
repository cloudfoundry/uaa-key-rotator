package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Main", func() {
	var session *gexec.Session

	JustBeforeEach(func() {
		uaaRotatorCmd := exec.Command(uaaRotatorBuildPath)

		var err error
		session, err = gexec.Start(uaaRotatorCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should be able to run", func() {
		Eventually(session).Should(gbytes.Say("hello"))
	})
})
