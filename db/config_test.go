package db_test

import (
	"github.com/cloudfoundry/uaa-key-rotator/config"
	db2 "github.com/cloudfoundry/uaa-key-rotator/db"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	var rotatorConfig *config.RotatorConfig

	BeforeEach(func() {
		rotatorConfig = &config.RotatorConfig{
			DatabaseUsername: "username",
			DatabasePassword: "password",
			DatabaseHostname: "localhost",
			DatabasePort:     "9876",
			DatabaseName:     "uaa",
		}
	})

	Describe("MYSQL", func() {
		BeforeEach(func() {
			rotatorConfig.DatabaseScheme = "mysql"
		})

		It("should generate connection uri", func() {
			connectionURI := db2.ConnectionURI(rotatorConfig)
			Expect(connectionURI).To(Equal("username:password@tcp(localhost:9876)/uaa?parseTime=true&timeout=120s&readTimeout=120s&writeTimeout=120s"))
		})

		Context("when tls is enabled", func() {
			BeforeEach(func() {
				rotatorConfig.DatabaseTlsEnabled = true
			})

			It("should generate connection uri", func() {
				connectionURI := db2.ConnectionURI(rotatorConfig)
				Expect(connectionURI).To(Equal("username:password@tcp(localhost:9876)/uaa?parseTime=true&timeout=120s&readTimeout=120s&writeTimeout=120s&useSSL=true&trustServerCertificate=false"))
			})

			Context("when skip ssl validation is enabled", func() {
				BeforeEach(func() {
					rotatorConfig.DatabaseSkipSSLValidation = true
				})

				It("should generate connection uri", func() {
					connectionURI := db2.ConnectionURI(rotatorConfig)
					Expect(connectionURI).To(Equal("username:password@tcp(localhost:9876)/uaa?parseTime=true&timeout=120s&readTimeout=120s&writeTimeout=120s&useSSL=true&trustServerCertificate=true"))
				})
			})

			Context("when tls_protocols are provided", func() {
				BeforeEach(func() {
					rotatorConfig.DatabaseTLSProtocols = "TLSv1.2,TLSv1.1"
				})

				It("should generate connection uri", func() {
					connectionURI := db2.ConnectionURI(rotatorConfig)
					Expect(connectionURI).To(Equal("username:password@tcp(localhost:9876)/uaa?parseTime=true&timeout=120s&readTimeout=120s&writeTimeout=120s&useSSL=true&trustServerCertificate=false&enabledSslProtocolSuites=TLSv1.2,TLSv1.1"))
				})
			})
		})
	})
	Describe("POSTGRES", func() {
		BeforeEach(func() {
			rotatorConfig.DatabaseScheme = "postgres"
		})

		It("should generate connection uri", func() {
			connectionURI := db2.ConnectionURI(rotatorConfig)
			Expect(connectionURI).To(Equal("postgres://username:password@localhost:9876/uaa?connect_timeout=120&sslmode=disable"))
		})

		Context("when tls is enabled", func() {
			BeforeEach(func() {
				rotatorConfig.DatabaseTlsEnabled = true
			})

			It("should generate connection uri", func() {
				connectionURI := db2.ConnectionURI(rotatorConfig)
				Expect(connectionURI).To(Equal("postgres://username:password@localhost:9876/uaa?connect_timeout=120&sslmode=verify-ca"))
			})

			Context("when skip ssl validation is enabled", func() {
				BeforeEach(func() {
					rotatorConfig.DatabaseSkipSSLValidation = true
				})

				It("should generate connection uri", func() {
					connectionURI := db2.ConnectionURI(rotatorConfig)
					Expect(connectionURI).To(Equal("postgres://username:password@localhost:9876/uaa?connect_timeout=120&sslmode=require"))
				})
			})
		})
	})
})
