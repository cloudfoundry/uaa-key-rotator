package rotator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cloudfoundry/uaa-key-rotator/rotator"
	"time"
)

var _ = Describe("KeyService", func() {
	var uaaKeyService rotator.UaaKeyService

	BeforeEach(func() {
		uaaKeyService = rotator.UaaKeyService{
			ActiveKeyLabel: "active-key-label",
			EncryptionKeys: []rotator.EncryptionKey{
				{Label: "active-key-label", Passphrase: "passphrase1"},
				{Label: "key-2", Passphrase: "passphrase2"},
			},
		}
	})

	It("should return the correct active key", func() {
		activeKeyLabel, _, err := uaaKeyService.ActiveKey()
		Expect(err).NotTo(HaveOccurred())
		Expect(activeKeyLabel).To(Equal(activeKeyLabel))
	})

	It("should be the identity to encrypt and then decrypt", func() {
		plainText := "some random plain text"
		activeKeyLabel, activeKeyEncryptor, err := uaaKeyService.ActiveKey()
		Expect(err).NotTo(HaveOccurred())

		activeKeyDecryptor, err := uaaKeyService.Key(activeKeyLabel)
		Expect(err).NotTo(HaveOccurred())

		encryptedValue, err := activeKeyEncryptor.Encrypt(plainText)
		Expect(err).NotTo(HaveOccurred())

		decryptedValue, err := activeKeyDecryptor.Decrypt(encryptedValue)
		Expect(err).NotTo(HaveOccurred())

		Expect(decryptedValue).To(Equal(plainText))
	})

	Context("when encrypting / decrypting with different keys", func() {
		It("should return a meaningful error", func() {
			plainText := "some random plain text"
			_, activeKeyEncryptor, err := uaaKeyService.ActiveKey()
			Expect(err).NotTo(HaveOccurred())

			activeKeyDecryptor, err := uaaKeyService.Key("key-2")
			Expect(err).NotTo(HaveOccurred())

			encryptedValue, err := activeKeyEncryptor.Encrypt(plainText)
			Expect(err).NotTo(HaveOccurred())

			decryptedValue, err := activeKeyDecryptor.Decrypt(encryptedValue)
			Expect(err).To(HaveOccurred())
			Expect(decryptedValue).NotTo(Equal(plainText))
		})
	})

	Context("when asking for a key that does not exist", func() {
		It("should return a meaningful error", func() {
			_, err := uaaKeyService.Key("key-does-not-exist")
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("unable to find key: key-does-not-exist"))
		})
	})

	Context("when asking for an active key that does not exist", func() {
		It("should return a meaningful error", func() {
			missingActiveKey := "active-key-does-not-exist" + time.Now().String()
			uaaKeyService.ActiveKeyLabel = missingActiveKey
			_, _, err := uaaKeyService.ActiveKey()
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("unable to find active key: " + missingActiveKey))
		})
	})
})
