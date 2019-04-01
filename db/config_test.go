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
			connectionURI, err := db2.ConnectionURI(rotatorConfig)
			Expect(err).NotTo(HaveOccurred())

			Expect(connectionURI).To(Equal("username:password@tcp(localhost:9876)/uaa?parseTime=true&timeout=240s&readTimeout=240s&writeTimeout=240s"))
		})

		Context("when tls is enabled", func() {
			BeforeEach(func() {
				rotatorConfig.DatabaseTlsEnabled = true
			})

			It("should generate connection uri", func() {
				connectionURI, err := db2.ConnectionURI(rotatorConfig)
				Expect(err).NotTo(HaveOccurred())

				Expect(connectionURI).To(Equal("username:password@tcp(localhost:9876)/uaa?parseTime=true&timeout=240s&readTimeout=240s&writeTimeout=240s&tls=true"))
			})

			Context("when skip ssl validation is enabled", func() {
				BeforeEach(func() {
					rotatorConfig.DatabaseSkipSSLValidation = true
				})

				It("should generate connection uri", func() {
					connectionURI, err := db2.ConnectionURI(rotatorConfig)
					Expect(err).NotTo(HaveOccurred())

					Expect(connectionURI).To(Equal("username:password@tcp(localhost:9876)/uaa?parseTime=true&timeout=240s&readTimeout=240s&writeTimeout=240s&tls=skip-verify"))
				})
			})
		})

		Context("when an invalid database port is used", func() {
			BeforeEach(func() {
				rotatorConfig.DatabasePort = "not-a-number"
			})
			It("should throw a meaningful error", func() {
				_, err := db2.ConnectionURI(rotatorConfig)
				Expect(err).To(HaveOccurred())

			})
		})
	})
	Describe("POSTGRES", func() {
		BeforeEach(func() {
			rotatorConfig.DatabaseScheme = "postgres"
		})

		It("should generate connection uri", func() {
			connectionURI, err := db2.ConnectionURI(rotatorConfig)
			Expect(err).NotTo(HaveOccurred())
			Expect(connectionURI).To(Equal("postgres://username:password@localhost:9876/uaa?connect_timeout=240&sslmode=disable"))
		})

		Context("when tls is enabled", func() {
			BeforeEach(func() {
				rotatorConfig.DatabaseTlsEnabled = true
			})

			It("should generate connection uri", func() {
				connectionURI, err := db2.ConnectionURI(rotatorConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(connectionURI).To(Equal("postgres://username:password@localhost:9876/uaa?connect_timeout=240&sslmode=verify-ca"))
			})

			Context("when skip ssl validation is enabled", func() {
				BeforeEach(func() {
					rotatorConfig.DatabaseSkipSSLValidation = true
				})

				It("should generate connection uri", func() {
					connectionURI, err := db2.ConnectionURI(rotatorConfig)
					Expect(err).NotTo(HaveOccurred())

					Expect(connectionURI).To(Equal("postgres://username:password@localhost:9876/uaa?connect_timeout=240&sslmode=require"))
				})
			})
		})
	})
})
