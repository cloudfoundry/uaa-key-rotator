package main_test

import (
	"encoding/json"
	"github.com/cloudfoundry/uaa-key-rotator/config"
	dbRotator "github.com/cloudfoundry/uaa-key-rotator/db"
	"github.com/cloudfoundry/uaa-key-rotator/db/testutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"time"
)

var _ = Describe("Main", func() {
	var session *gexec.Session
	var rotatorConfig *config.RotatorConfig
	var rotatorConfigFile *os.File
	var activeKey config.EncryptionKey

	BeforeEach(func() {
		activeKey = config.EncryptionKey{
			Label:      "active-key",
			Passphrase: "new-secret",
		}

		rotatorConfig = &config.RotatorConfig{
			ActiveKeyLabel: activeKey.Label,
			EncryptionKeys: []config.EncryptionKey{
				activeKey,
				oldKey,
			},
			DatabaseHostname: testutils.Hostname,
			DatabaseName:     testutils.DBName,
			DatabasePort:     testutils.Port,
			DatabaseScheme:   "postgres",
			DatabaseUsername: "postgres",
			DatabasePassword: "",
		}

		jsonConfig, err := json.Marshal(rotatorConfig)
		rotatorConfigFile, err = ioutil.TempFile(os.TempDir(), "rotator_config")
		Expect(err).NotTo(HaveOccurred())

		Expect(ioutil.WriteFile(rotatorConfigFile.Name(), jsonConfig, os.ModePerm)).To(Succeed())
	})

	JustBeforeEach(func() {
		uaaRotatorCmd := exec.Command(uaaRotatorBuildPath, "-config", rotatorConfigFile.Name())

		var err error
		session, err = gexec.Start(uaaRotatorCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should exit gracefully when an interrupt signal is received", func() {
		Eventually(session).Should(gbytes.Say("rotator has started"))

		session.Signal(os.Interrupt)

		Eventually(session).Should(gbytes.Say("shutting down gracefully..."))
		Eventually(session).ShouldNot(gexec.Exit(0))
	})

	It("should rotate encrypted data from an old key to the new 'active' key", func() {
		Eventually(session, 2*time.Minute).Should(gbytes.Say("rotator has finished"))
		session.Signal(syscall.SIGTERM)
		Eventually(session).ShouldNot(gbytes.Say("shutting down gracefully..."))

		mfaCredentials, err := dbRotator.ReadAll(db)
		Expect(err).NotTo(HaveOccurred())
		Expect(mfaCredentials).To(HaveLen(1))

		decryptedRotatedSecretKey := decryptCipherValue(mfaCredentials[0].SecretKey, activeKey.Passphrase)
		Expect(decryptedRotatedSecretKey).To(Equal("secret-key"))
	})
})
