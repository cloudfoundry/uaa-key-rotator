package crypto_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/ginkgo/extensions/table"
	"bytes"
	uaa_crypto "github.com/cloudfoundry/uaa-key-rotator/crypto"
)

var _ = Describe("UaaNonce", func() {

	Describe("nonce generator", func() {
		var nonceGenerator uaa_crypto.UaaNonceGenerator
		BeforeEach(func() {
			nonceGenerator = uaa_crypto.UaaNonceGenerator{}
		})

		It("should generate a nonce 12 bytes long", func() {
			nonce, err := nonceGenerator.GetNonce()
			Expect(err).ToNot(HaveOccurred())
			Expect(nonce).To(HaveLen(12))

		})
	})

	Describe("nonce accessor", func() {
		var nonceAccessor uaa_crypto.UaaNonceAccessor
		BeforeEach(func() {
			nonceAccessor = uaa_crypto.UaaNonceAccessor{}

		})

		It("should access nonce from cipher value", func() {
			cipherValue := []byte{}
			cipherValue = append(cipherValue, bytes.Repeat([]byte("n"), 12)...)
			cipherValue = append(cipherValue, bytes.Repeat([]byte("s"), 32)...)
			cipherValue = append(cipherValue, bytes.Repeat([]byte("x"), 52)...)
			salt, _ := nonceAccessor.GetNonce(cipherValue)

			Expect(salt).To(Equal(bytes.Repeat([]byte("n"), 12)))
		})

		table.DescribeTable("Given invalid nonce values", func(nonceSize int) {
			cipherValue := []byte{}
			cipherValue = append(cipherValue, bytes.Repeat([]byte("n"), nonceSize)...)
			_, err := nonceAccessor.GetNonce(cipherValue)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("nonce should be exactly 12 bytes in length"))
		},
			table.Entry("nonce value really below", 1),
			table.Entry("nonce value just below", 11))
	})
})
