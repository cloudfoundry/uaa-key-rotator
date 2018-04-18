package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cloudfoundry/uaa-key-rotator/config"
	"os"
	"io/ioutil"
	"errors"
	"github.com/onsi/ginkgo/extensions/table"
	"encoding/json"
	"github.com/onsi/gomega/gbytes"
)

var (
	activeKeyLabel      string
	activeKeyPassphrase string
)

var _ = Describe("Config", func() {
	var tempConfigFile *os.File
	configFileContent := `{ 
					"activeKeyLabel": "key1",
					"activeKeyPassphrase": "passphrase",
					"databaseHostname": "localhost",
					"databasePort": "5432",
					"databaseName": "uaadb",
					"databaseScheme": "postgresql",
					"databaseUsername": "admin",
					"databasePassword": "afdsafda"
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
		Expect(rotatorConfig.DatabaseHostname).To(Equal("localhost"))
		Expect(rotatorConfig.DatabasePort).To(Equal("5432"))
		Expect(rotatorConfig.DatabaseScheme).To(Equal("postgresql"))
		Expect(rotatorConfig.DatabaseName).To(Equal("uaadb"))
		Expect(rotatorConfig.DatabaseUsername).To(Equal("admin"))
		Expect(rotatorConfig.DatabasePassword).To(Equal("afdsafda"))
	})

	Context("when given an invalid json config file", func() {
		var requiredFields map[string]string

		BeforeEach(func() {
			requiredFields = map[string]string{
				"activeKeyLabel":      "active-key-value",
				"activeKeyPassphrase": "active-key-passphrase",
				"databaseHostname":    "db-hostname",
				"databasePort":        "db-port",
				"databaseScheme":      "db-scheme",
				"databaseName":        "db-name",
				"databaseUsername":    "db-username",
				"databasePassword":    "db-password",
			}
		})

		table.DescribeTable("invalid fields", func(invalidKey, invalidValue, errorDescription string) {
			cfg := cloneMap(requiredFields)

			cfg[invalidKey] = invalidValue

			jsonBytes, err := json.Marshal(cfg)
			Expect(err).NotTo(HaveOccurred())

			_, err = config.New(gbytes.BufferWithBytes(jsonBytes))
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(errorDescription))

		},
			table.Entry("invalid active key label", "activeKeyLabel", "", "Invalid config.: ActiveKeyLabel: zero value"),
			table.Entry("invalid active key passphrase", "activeKeyPassphrase", "", "Invalid config.: ActiveKeyPassphrase: zero value"),
			table.Entry("invalid db hostname", "databaseHostname", "", "Invalid config.: DatabaseHostname: zero value"),
			table.Entry("invalid db port", "databasePort", "", "Invalid config.: DatabasePort: zero value"),
			table.Entry("invalid db scheme", "databaseScheme", "", "Invalid config.: DatabaseScheme: zero value"),
			table.Entry("invalid db username", "databaseUsername", "", "Invalid config.: DatabaseUsername: zero value"),
			table.Entry("invalid db password", "databasePassword", "", "Invalid config.: DatabasePassword: zero value"),
			table.Entry("invalid db ", "databaseName", "", "Invalid config.: DatabaseName: zero value"),
		)

	})

	Context("Given invalid rotator config", func() {
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

func cloneMap(requiredFields map[string]string) map[string]string {
	configMap := map[string]string{}

	for key, value := range requiredFields {
		configMap[key] = value
	}

	return configMap
}
