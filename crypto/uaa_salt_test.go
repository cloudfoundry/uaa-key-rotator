package crypto_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/ginkgo/extensions/table"
	"bytes"
	. "github.com/cloudfoundry/uaa-key-rotator/crypto"
)

var _ = Describe("UaaSalt", func() {

	Describe("salt generator", func() {
		var saltGenerator UaaSaltGenerator
		BeforeEach(func() {
			saltGenerator = UaaSaltGenerator{}
		})

		It("should generate a salt 32 bytes long", func() {
			salt, err := saltGenerator.GetSalt()
			Expect(err).NotTo(HaveOccurred())
			Expect(salt).To(HaveLen(32))
		})

		It("should generate different salts every time", func() {
			salt1, err := saltGenerator.GetSalt()
			Expect(err).NotTo(HaveOccurred())

			salt2, err := saltGenerator.GetSalt()
			Expect(err).NotTo(HaveOccurred())

			Expect(salt1).ToNot(Equal(salt2))
		})
	})

	Describe("salt accessor", func() {
		var saltAccessor UaaSaltAccessor
		BeforeEach(func() {
			saltAccessor = UaaSaltAccessor{}

		})

		It("should access salt from cipher value", func() {

			cipherValue := []byte{}
			cipherValue = append(cipherValue, bytes.Repeat([]byte("n"), 12)...)
			cipherValue = append(cipherValue, bytes.Repeat([]byte("s"), 32)...)
			cipherValue = append(cipherValue, bytes.Repeat([]byte("x"), 52)...)

			salt, _ := saltAccessor.GetSalt(cipherValue)

			Expect(salt).To(Equal(bytes.Repeat([]byte("s"), 32)))

		})

		table.DescribeTable("Given invalid salt values", func(saltSize int) {
			cipherValue := []byte{}
			cipherValue = append(cipherValue, bytes.Repeat([]byte("n"), 12)...)
			cipherValue = append(cipherValue, bytes.Repeat([]byte("s"), saltSize)...)
			_, err := saltAccessor.GetSalt(cipherValue)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("cipher should be more than 45 bytes"))
		},
			table.Entry("nonce value really below", 1),
			table.Entry("nonce value just below", 31))
	})
})
