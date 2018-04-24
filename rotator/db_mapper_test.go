package rotator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/uaa-key-rotator/crypto"
	. "github.com/cloudfoundry/uaa-key-rotator/rotator"
)

var _ = Describe("DbMapper", func() {
	It("should map to base64 encoded string", func() {
		encrypted := crypto.EncryptedValue{[]byte("salt"), []byte("nonce"), []byte("encryptedval")}

		mapped, err := DbMapper{}.Map(encrypted)

		Expect(err).NotTo(HaveOccurred())
		Expect(mapped).NotTo(BeNil())

		//base64 <<< 'noncesaltencryptedval'
		Expect(mapped).To(Equal([]byte("bm9uY2VzYWx0ZW5jcnlwdGVkdmFs")))
	})
})
