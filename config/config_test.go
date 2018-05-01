package config_test

import (
	"encoding/json"
	"errors"
	"github.com/cloudfoundry/uaa-key-rotator/config"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"io/ioutil"
	"os"
)

var (
	activeKeyLabel string
)

var _ = Describe("Config", func() {
	var tempConfigFile *os.File
	configFileContent := `{ 
					"activeKeyLabel": "key1",
					"encryptionKeys": [
						{"label": "active-key", "passphrase": "secret"}
					],
					"databaseHostname": "localhost",
					"databasePort": "5432",
					"databaseName": "uaadb",
					"databaseScheme": "postgresql",
					"databaseUsername": "admin",
					"databasePassword": "afdsafda",
					"databaseTlsEnabled": true,
					"databaseSkipSSLValidation": true,
					"databaseTLSProtocols": "TLSv1.2,TLSv1.1"
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
		Expect(rotatorConfig.EncryptionKeys).To(ConsistOf(config.EncryptionKey{
			Label: "active-key", Passphrase: "secret",
		}))
		Expect(rotatorConfig.DatabaseHostname).To(Equal("localhost"))
		Expect(rotatorConfig.DatabasePort).To(Equal("5432"))
		Expect(rotatorConfig.DatabaseScheme).To(Equal("postgres"))
		Expect(rotatorConfig.DatabaseName).To(Equal("uaadb"))
		Expect(rotatorConfig.DatabaseUsername).To(Equal("admin"))
		Expect(rotatorConfig.DatabasePassword).To(Equal("afdsafda"))
		Expect(rotatorConfig.DatabaseTlsEnabled).To(BeTrue())
		Expect(rotatorConfig.DatabaseSkipSSLValidation).To(BeTrue())
		Expect(rotatorConfig.DatabaseTLSProtocols).To(Equal("TLSv1.2,TLSv1.1"))
	})

	Context("when given an invalid json config file", func() {
		var requiredFields map[string]interface{}

		BeforeEach(func() {
			requiredFields = map[string]interface{}{
				"activeKeyLabel": "active-key-value",
				"encryptionKeys": []map[string]string{
					{"label": "active-key", "passphrase": "secret2"},
				},
				"databaseHostname": "db-hostname",
				"databasePort":     "db-port",
				"databaseScheme":   "db-scheme",
				"databaseName":     "db-name",
				"databaseUsername": "db-username",
				"databasePassword": "db-password",
				"databaseTlsEnabled": true,
				"databaseSkipSSLValidation": true,
			}
		})

		table.DescribeTable("invalid fields", func(invalidKey string, invalidValue interface{}, errorDescription string) {
			cfg := cloneMap(requiredFields)

			cfg[invalidKey] = invalidValue

			jsonBytes, err := json.Marshal(cfg)
			Expect(err).NotTo(HaveOccurred())

			_, err = config.New(gbytes.BufferWithBytes(jsonBytes))
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(errorDescription))

		},
			table.Entry("invalid active key label", "activeKeyLabel", "", "Invalid config.: ActiveKeyLabel: zero value"),
			table.Entry("invalid encryption keys", "encryptionKeys", []map[string]string{}, "Invalid config.: EncryptionKeys: zero value"),
			table.Entry("invalid encryption keys", "encryptionKeys", []map[string]string{{"foobar": "value", "passphrase": "secret"}}, "Invalid config.: EncryptionKeys[0].Label: zero value"),
			table.Entry("invalid encryption keys", "encryptionKeys", []map[string]string{{"label": "value", "asdfasf": "secret"}}, "Invalid config.: EncryptionKeys[0].Passphrase: zero value"),
			table.Entry("invalid db hostname", "databaseHostname", "", "Invalid config.: DatabaseHostname: zero value"),
			table.Entry("invalid db port", "databasePort", "", "Invalid config.: DatabasePort: zero value"),
			table.Entry("invalid db scheme", "databaseScheme", "", "Invalid config.: DatabaseScheme: zero value"),
			table.Entry("invalid db username", "databaseUsername", "", "Invalid config.: DatabaseUsername: zero value"),
			table.Entry("invalid db ", "databaseName", "", "Invalid config.: DatabaseName: zero value"),
			table.Entry("invalid db tls enabled", "databaseTlsEnabled", "", "Malformed JSON provided.: json: cannot unmarshal string into Go struct field RotatorConfig.databaseTlsEnabled of type bool"),
			table.Entry("invalid db databaseSkipSSLValidation", "databaseSkipSSLValidation", "", "Malformed JSON provided.: json: cannot unmarshal string into Go struct field RotatorConfig.databaseSkipSSLValidation of type bool"),
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

func cloneMap(requiredFields map[string]interface{}) map[string]interface{} {
	configMap := map[string]interface{}{}

	for key, value := range requiredFields {
		configMap[key] = value
	}

	return configMap
}
