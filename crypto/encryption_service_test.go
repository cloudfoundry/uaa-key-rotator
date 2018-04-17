package crypto_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/cloudfoundry/uaa-key-rotator/crypto"
	"bytes"
	"github.com/onsi/ginkgo/extensions/table"
)

var _ = Describe("EncryptionService", func() {
	var encryptor Encryptor
	var salt []byte
	var nonce []byte

	BeforeEach(func() {
		salt = bytes.Repeat([]byte("s"), 32)
		nonce = bytes.Repeat([]byte("n"), 12)
	})

	JustBeforeEach(func() {
		encryptor = Encryptor{
			Passphrase: "passphrase",
			Salt:       salt,
			Nonce:      nonce,
		}
	})

	Describe("Encrypt", func() {
		It("should encrypt data", func() {
			plainTextValue := "data-to-encrypt"
			encryptedData, err := encryptor.Encrypt(plainTextValue)

			Expect(err).NotTo(HaveOccurred())
			Expect(encryptedData).ToNot(ContainSubstring(plainTextValue))

			Expect(encryptedData).To(HavePrefix(string(salt) + string(nonce)))
		})
		table.DescribeTable("Given invalid salt values", func(saltSize int) {
			encryptor.Salt = bytes.Repeat([]byte("s"), saltSize)
			_, err := encryptor.Encrypt("data-to-encrypt")

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("Salt should be exactly 32 bytes in length."))
		},
			table.Entry("salt value just below", 31),
			table.Entry("salt value just above", 33))

		table.DescribeTable("Given invalid nonce values", func(nonceSize int) {
			encryptor.Nonce = bytes.Repeat([]byte("s"), nonceSize)
			_, err := encryptor.Encrypt("data-to-encrypt")

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("Nonce should be exactly 12 bytes in length."))
		},
			table.Entry("nonce value just below", 11),
			table.Entry("nonce value just above", 13))
	})

	Describe("Decrypt", func() {
		table.DescribeTable("Given invalid salt values", func(saltSize int) {
			encryptor.Salt = bytes.Repeat([]byte("s"), saltSize)
			_, err := encryptor.Decrypt(bytes.Repeat([]byte("x"), 45))

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("Salt should be exactly 32 bytes in length."))
		},
			table.Entry("salt value just below", 31),
			table.Entry("salt value just above", 33))

		It("should be able to decrypt data that was previously encrypted", func() {
			plainText := "data-to-encrypt"
			encryptedData, err := encryptor.Encrypt(plainText)
			Expect(err).NotTo(HaveOccurred())
			Expect(encryptedData).ToNot(BeNil())

			decryptedData, err := encryptor.Decrypt(encryptedData)

			Expect(err).NotTo(HaveOccurred())
			Expect(decryptedData).To(Equal(plainText))
		})

		Context("when empty ciphervalue is provided", func() {
			It("should return a meaningful error", func() {
				_, err := encryptor.Decrypt(bytes.Repeat([]byte("x"), 44))
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("Unable to decrypt due to invalid CipherText."))
			})

		})
	})
})
