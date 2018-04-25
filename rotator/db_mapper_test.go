package rotator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/uaa-key-rotator/crypto"
	. "github.com/cloudfoundry/uaa-key-rotator/rotator"
)

var _ = Describe("DbMapper", func() {
	var dbMapper DbMapper
	BeforeEach(func() {
		dbMapper = DbMapper{}
	})

	Describe("Map", func() {
		It("should map to base64 encoded string", func() {
			encrypted := crypto.EncryptedValue{[]byte("salt"), []byte("nonce"), []byte("encryptedval")}

			mapped, err := dbMapper.Map(encrypted)

			Expect(err).NotTo(HaveOccurred())
			Expect(mapped).NotTo(BeNil())

			//base64 <<< 'noncesaltencryptedval'
			Expect(mapped).To(Equal([]byte("bm9uY2VzYWx0ZW5jcnlwdGVkdmFs")))
		})
	})

	Describe("MapBase64ToCipherValue", func() {
		It("should decode a base64 encoded string", func() {
			decodedValue, err := dbMapper.MapBase64ToCipherValue("bm9uY2VzYWx0ZW5jcnlwdGVkdmFs")
			Expect(err).NotTo(HaveOccurred())

			//base64 -D <<< 'bm9uY2VzYWx0ZW5jcnlwdGVkdmFs'
			Expect(decodedValue).To(Equal([]byte("noncesaltencryptedval")))

		})
	})
})
