package crypto_test

import (
	"bytes"
	uaa_crypto "github.com/cloudfoundry/uaa-key-rotator/crypto"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("UaaCipher", func() {
	Describe("cipher accessor", func() {
		var cipherAccessor uaa_crypto.UAACipherAccessor
		BeforeEach(func() {
			cipherAccessor = uaa_crypto.UAACipherAccessor{}
		})

		It("should access cipher from 'UAA' cipher value", func() {
			cipherValue := []byte{}
			cipherValue = append(cipherValue, bytes.Repeat([]byte("n"), 12)...)
			cipherValue = append(cipherValue, bytes.Repeat([]byte("s"), 32)...)
			cipherValue = append(cipherValue, bytes.Repeat([]byte("x"), 52)...)
			cipher, _ := cipherAccessor.GetCipher(cipherValue)

			Expect(cipher).To(Equal(bytes.Repeat([]byte("x"), 52)))
		})

		table.DescribeTable("Given invalid cipher values", func(uaaCipherSize int) {
			cipherValue := []byte{}
			cipherValue = append(cipherValue, bytes.Repeat([]byte("n"), uaaCipherSize)...)
			_, err := cipherAccessor.GetCipher(cipherValue)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("uaa cipher should be at least 44 bytes in length"))
		},
			table.Entry("uaa cipher value really below", 1),
			table.Entry("uaa cipher value just below", 43))
	})
})
