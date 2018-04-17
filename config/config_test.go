package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cloudfoundry/uaa-key-rotator/config"
	"os"
	"io/ioutil"
	"errors"
)

var (
	activeKeyLabel      string
	activeKeyPassphrase string
)

var _ = Describe("Config", func() {
	var tempConfigFile *os.File
	configFileContent := `{ 
					"activeKeyLabel": "key1",
					"activeKeyPassphrase": "passphrase"
					}
					`

	JustBeforeEach(func() {
		var err error
		tempConfigFile, err = ioutil.TempFile(os.TempDir(), "config")
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(tempConfigFile.Name(), []byte(configFileContent), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should unmarshal a valid config file", func() {
		rotatorConfig, err := config.New(tempConfigFile)

		Expect(err).ToNot(HaveOccurred())
		Expect(rotatorConfig.ActiveKeyLabel).To(Equal("key1"))
		Expect(rotatorConfig.ActiveKeyPassphrase).To(Equal("passphrase"))
	})

	Context("Given invalid rotator config", func() {
		Context("when wellformed json with invalid values is provided", func() {
			BeforeEach(func() {
				configFileContent = `{ 
					"activeKeyLabel": "",
					"activeKeyPassphrase": ""
					}
					`
			})

			It("should throw an error", func() {
				_, err := config.New(tempConfigFile)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring("Invalid config.")))
			})
		})

		Context("when malformed json is provided", func() {
			BeforeEach(func() {
				configFileContent = `42 +`
			})

			It("should throw an error", func() {
				_, err := config.New(tempConfigFile)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("Malformed JSON provided.: invalid character '+' after top-level value"))
			})

		})
	})

	Context("when an invalid configFileReader is provided", func() {
		It("should return a meaningful error", func() {
			_, err := config.New(BadConfigFile{})
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("Unable to read config: cannot read"))
		})
	})
})

type BadConfigFile struct{}

func (BadConfigFile) Read(p []byte) (int, error) {
	return 0, errors.New("cannot read")
}
