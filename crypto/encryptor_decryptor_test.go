package crypto_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/cloudfoundry/uaa-key-rotator/crypto"
	"bytes"
	"github.com/cloudfoundry/uaa-key-rotator/crypto/cryptofakes"
)

var _ = Describe("Encryptor", func() {
	var encryptor Encryptor
	var decryptor Decryptor

	var salt []byte
	var nonce []byte

	var saltGenerator *cryptofakes.FakeSaltGenerator
	var nonceGenerator *cryptofakes.FakeNonceGenerator

	var cipherSaltAccessor *cryptofakes.FakeCipherSaltAccessor
	var cipherNonceAccessor *cryptofakes.FakeCipherNonceAccessor

	BeforeEach(func() {
		salt = bytes.Repeat([]byte("s"), 32)
		nonce = bytes.Repeat([]byte("n"), 12)

		saltGenerator = &cryptofakes.FakeSaltGenerator{}
		nonceGenerator = &cryptofakes.FakeNonceGenerator{}

		cipherSaltAccessor = &cryptofakes.FakeCipherSaltAccessor{}
		cipherNonceAccessor = &cryptofakes.FakeCipherNonceAccessor{}

		saltGenerator.GetSaltReturns(salt)
		nonceGenerator.GetNonceReturns(nonce)

		cipherSaltAccessor.GetSaltReturns(salt, nil)
		cipherNonceAccessor.GetNonceReturns(nonce, nil)
	})

	JustBeforeEach(func() {
		passphrase := "passphrase"
		decryptor = Decryptor{
			Passphrase:          passphrase,
			CipherNonceAccessor: cipherNonceAccessor,
			CipherSaltAccessor:  cipherSaltAccessor,
		}
		encryptor = Encryptor{
			Passphrase:     passphrase,
			SaltGenerator:  saltGenerator,
			NonceGenerator: nonceGenerator,
		}
	})

	Describe("Encrypt", func() {
		It("should encrypt data", func() {
			plainTextValue := "data-to-encrypt"
			encryptedData, err := encryptor.Encrypt(plainTextValue)

			Expect(err).NotTo(HaveOccurred())
			Expect(encryptedData.Salt).To(Equal(salt))
			Expect(encryptedData.Nonce).To(Equal(nonce))
			Expect(encryptedData.CipherValue).ToNot(BeEmpty())
		})
	})

	Describe("Decrypt", func() {
		It("should be able to decrypt data that was previously encrypted", func() {
			plainText := "data-to-encrypt"
			encryptedData, err := encryptor.Encrypt(plainText)
			Expect(err).NotTo(HaveOccurred())
			Expect(encryptedData).ToNot(BeNil())

			decryptedData, err := decryptor.Decrypt(encryptedData.CipherValue)

			Expect(err).NotTo(HaveOccurred())
			Expect(decryptedData).To(Equal(plainText))
		})

		Context("when empty ciphervalue is provided", func() {
			It("should return a meaningful error", func() {
				_, err := decryptor.Decrypt(nil)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("Unable to decrypt due to empty CipherText."))
			})

		})
	})
})
