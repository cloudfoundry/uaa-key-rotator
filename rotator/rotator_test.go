package rotator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/cloudfoundry/uaa-key-rotator/rotator/rotatorfakes"
	"github.com/cloudfoundry/uaa-key-rotator/entity"
	"github.com/cloudfoundry/uaa-key-rotator/crypto/cryptofakes"
	"github.com/cloudfoundry/uaa-key-rotator/rotator"
	"github.com/pkg/errors"
	"time"
	"github.com/onsi/ginkgo/extensions/table"
	"github.com/cloudfoundry/uaa-key-rotator/crypto"
	"database/sql"
)

var _ = Describe("Rotator", func() {
	var updatedCredential entity.MfaCredential
	var rotatorError error

	var activeKeyLabel string
	var fakeKeyService *rotatorfakes.FakeKeyService

	var fakeDecryptor *cryptofakes.FakeDecryptor
	var fakeEncryptor *cryptofakes.FakeEncryptor

	var fakeSaltAccessor *cryptofakes.FakeCipherSaltAccessor
	var fakeNonceAccessor *cryptofakes.FakeCipherNonceAccessor

	var fakeDbMapper *rotatorfakes.FakeMapEncryptedValueToDB

	var fakeDecrpytedScratchCodes string
	var fakeDecryptedSecretKey string
	var fakeDecryptedValidationCode string
	var scratchCodes string
	var secretKey string
	var encryptedValidationCode string
	var secretKeySalt string
	var encryptedValidationCodeSalt string
	var scratchCodesSalt string
	var scratchCodesNonce string
	var encryptedValidationCodeNonce string
	var secretKeyNonce string

	var fakeEncryptedScratchCode crypto.EncryptedValue
	var fakeEncryptedSecretKey crypto.EncryptedValue
	var fakeEncryptedEncryptedValidationCode crypto.EncryptedValue
	var fakeRotatedScratchCode string
	var fakeRotatedSecretKey string
	var fakeRotatedEncryptedValidationCode string

	BeforeEach(func() {
		fakeKeyService = &rotatorfakes.FakeKeyService{}
		fakeDecryptor = &cryptofakes.FakeDecryptor{}
		fakeSaltAccessor = &cryptofakes.FakeCipherSaltAccessor{}
		fakeNonceAccessor = &cryptofakes.FakeCipherNonceAccessor{}

		scratchCodes = "encrypted-scratch-codes" + time.Now().String()
		secretKey = "secret-key" + time.Now().String()
		encryptedValidationCode = "encrypted-validation-code" + time.Now().String()

		fakeKeyService.KeyReturns(fakeDecryptor)
		fakeDecrpytedScratchCodes = "whatever-we-return-in-our-fake" + time.Now().String()
		fakeDecryptedSecretKey = "decrypted secret key" + time.Now().String()
		fakeDecryptedValidationCode = "validation code secret key" + time.Now().String()

		fakeDecryptor.DecryptReturnsOnCall(0, fakeDecrpytedScratchCodes, nil)
		fakeDecryptor.DecryptReturnsOnCall(1, fakeDecryptedSecretKey, nil)
		fakeDecryptor.DecryptReturnsOnCall(2, fakeDecryptedValidationCode, nil)

		scratchCodesSalt = "scratch-codes-salt" + time.Now().String()
		fakeSaltAccessor.GetSaltReturnsOnCall(0, []byte(scratchCodesSalt), nil)
		secretKeySalt = "secret-key-salt" + time.Now().String()
		fakeSaltAccessor.GetSaltReturnsOnCall(1, []byte(secretKeySalt), nil)
		encryptedValidationCodeSalt = "encrypted-validation-code-salt" + time.Now().String()
		fakeSaltAccessor.GetSaltReturnsOnCall(2, []byte(encryptedValidationCodeSalt), nil)

		scratchCodesNonce = "scratch-codes-nonce" + time.Now().String()
		fakeNonceAccessor.GetNonceReturnsOnCall(0, []byte(scratchCodesNonce), nil)
		secretKeyNonce = "secret-key-nonce" + time.Now().String()
		fakeNonceAccessor.GetNonceReturnsOnCall(1, []byte(secretKeyNonce), nil)
		encryptedValidationCodeNonce = "encrypted-validation-codes-nonce" + time.Now().String()
		fakeNonceAccessor.GetNonceReturnsOnCall(2, []byte(encryptedValidationCodeNonce), nil)

		fakeEncryptor = &cryptofakes.FakeEncryptor{}
		activeKeyLabel = "key-2"
		fakeKeyService.ActiveKeyReturns(activeKeyLabel, fakeEncryptor)

		fakeDbMapper = &rotatorfakes.FakeMapEncryptedValueToDB{}
		fakeEncryptedScratchCode = crypto.EncryptedValue{CipherValue: []byte("rotated_scratch_code")}
		fakeEncryptor.EncryptReturnsOnCall(0, fakeEncryptedScratchCode, nil)
		fakeEncryptedSecretKey = crypto.EncryptedValue{CipherValue: []byte("rotated_secret_key")}
		fakeEncryptor.EncryptReturnsOnCall(1, fakeEncryptedSecretKey, nil)
		fakeEncryptedEncryptedValidationCode = crypto.EncryptedValue{CipherValue: []byte("rotated_validation_code")}
		fakeEncryptor.EncryptReturnsOnCall(2, fakeEncryptedEncryptedValidationCode, nil)

		fakeRotatedScratchCode = "encrypted scratch code"
		fakeDbMapper.MapReturnsOnCall(0, []byte(fakeRotatedScratchCode), nil)

		fakeRotatedSecretKey = "rotated secret key"
		fakeDbMapper.MapReturnsOnCall(1, []byte(fakeRotatedSecretKey), nil)

		fakeRotatedEncryptedValidationCode = "rotated encrypted validation code"
		fakeDbMapper.MapReturnsOnCall(2, []byte(fakeRotatedEncryptedValidationCode), nil)

	})

	JustBeforeEach(func() {
		updatedCredential, rotatorError = rotator.Rotate(
			entity.MfaCredential{
				UserId:                  "some-user-id",
				MfaProviderId:           "some-provider-id",
				ZoneId:                  "some-zone-id",
				EncryptionKeyLabel:      "key-1",
				ValidationCode:          sql.NullInt64{Int64: 1},
				ScratchCodes:            scratchCodes,
				SecretKey:               secretKey,
				EncryptedValidationCode: encryptedValidationCode,
			},
			fakeKeyService,
			fakeSaltAccessor,
			fakeNonceAccessor,
			fakeDbMapper,
		)

		Expect(updatedCredential.ValidationCode).To(Equal(sql.NullInt64{Int64: 1}))
		Expect(updatedCredential.EncryptionKeyLabel).To(Equal(activeKeyLabel))
		Expect(updatedCredential.UserId).To(Equal("some-user-id"))
		Expect(string(updatedCredential.MfaProviderId)).To(Equal("some-provider-id"))
		Expect(string(updatedCredential.ZoneId)).To(Equal("some-zone-id"))

	})

	It("should rotate encrypted values from using one key to another", func() {
		Expect(rotatorError).NotTo(HaveOccurred())
		Expect(fakeKeyService.KeyCallCount()).To(Equal(1))

		decryptArgsForScratchCodes := fakeDecryptor.DecryptArgsForCall(0)
		Expect(string(decryptArgsForScratchCodes.CipherValue)).To(Equal(scratchCodes))
		Expect(string(decryptArgsForScratchCodes.Salt)).To(Equal(scratchCodesSalt))
		Expect(string(decryptArgsForScratchCodes.Nonce)).To(Equal(scratchCodesNonce))

		decryptArgsForSecretKey := fakeDecryptor.DecryptArgsForCall(1)
		Expect(string(decryptArgsForSecretKey.CipherValue)).To(Equal(secretKey))
		Expect(string(decryptArgsForSecretKey.Salt)).To(Equal(secretKeySalt))
		Expect(string(decryptArgsForSecretKey.Nonce)).To(Equal(secretKeyNonce))

		decryptArgsForEncryptedValidationCode := fakeDecryptor.DecryptArgsForCall(2)
		Expect(string(decryptArgsForEncryptedValidationCode.CipherValue)).To(Equal(encryptedValidationCode))
		Expect(string(decryptArgsForEncryptedValidationCode.Salt)).To(Equal(encryptedValidationCodeSalt))
		Expect(string(decryptArgsForEncryptedValidationCode.Nonce)).To(Equal(encryptedValidationCodeNonce))

		Expect(fakeKeyService.ActiveKeyCallCount()).To(Equal(1))
		Expect(fakeEncryptor.EncryptCallCount()).To(Equal(3))

		Expect(fakeEncryptor.EncryptArgsForCall(0)).To(Equal(fakeDecrpytedScratchCodes))
		Expect(fakeEncryptor.EncryptArgsForCall(1)).To(Equal(fakeDecryptedSecretKey))
		Expect(fakeEncryptor.EncryptArgsForCall(2)).To(Equal(fakeDecryptedValidationCode))

		Expect(fakeDbMapper.MapCallCount()).To(Equal(3))
		Expect(fakeDbMapper.MapArgsForCall(0)).To(Equal(fakeEncryptedScratchCode))
		Expect(fakeDbMapper.MapArgsForCall(1)).To(Equal(fakeEncryptedSecretKey))
		Expect(fakeDbMapper.MapArgsForCall(2)).To(Equal(fakeEncryptedEncryptedValidationCode))

		Expect(updatedCredential).To(MatchFields(IgnoreExtras, Fields{
			"ScratchCodes":            Equal(fakeRotatedScratchCode),
			"SecretKey":               Equal(fakeRotatedSecretKey),
			"EncryptedValidationCode": Equal(fakeRotatedEncryptedValidationCode),
		}))
	})

	table.DescribeTable("when accessing the salt returns an error", func(errorIndex int) {
		fakeSaltAccessor = &cryptofakes.FakeCipherSaltAccessor{}
		var errorStr = "some error" + time.Now().String()
		fakeSaltAccessor.GetSaltReturnsOnCall(errorIndex, nil, errors.New(errorStr))

		updatedCredential, rotatorError = rotator.Rotate(
			entity.MfaCredential{
				EncryptionKeyLabel:      "key-1",
				ScratchCodes:            scratchCodes,
				SecretKey:               secretKey,
				EncryptedValidationCode: encryptedValidationCode,
			},
			fakeKeyService,
			fakeSaltAccessor,
			fakeNonceAccessor,
			fakeDbMapper,
		)

		Expect(rotatorError).To(HaveOccurred())
		Expect(rotatorError).To(MatchError("unable to access salt from cipher value provided: " + errorStr))
	},
		table.Entry("when accessing salt for scratch codes fails", 0),
		table.Entry("when accessing salt for secret key fails", 1),
		table.Entry("when accessing salt for encrypted validation codes fails", 2),
	)

	table.DescribeTable("when accessing the nonce returns an error", func(errorIndex int) {
		fakeNonceAccessor = &cryptofakes.FakeCipherNonceAccessor{}
		var errorStr = "some error" + time.Now().String()
		fakeNonceAccessor.GetNonceReturnsOnCall(errorIndex, nil, errors.New(errorStr))

		updatedCredential, rotatorError = rotator.Rotate(
			entity.MfaCredential{
				EncryptionKeyLabel:      "key-1",
				ScratchCodes:            scratchCodes,
				SecretKey:               secretKey,
				EncryptedValidationCode: encryptedValidationCode,
			},
			fakeKeyService,
			fakeSaltAccessor,
			fakeNonceAccessor,
			fakeDbMapper,
		)

		Expect(rotatorError).To(HaveOccurred())
		Expect(rotatorError).To(MatchError("unable to access nonce from cipher value provided: " + errorStr))
	},
		table.Entry("when accessing nonce for scratch codes fails", 0),
		table.Entry("when accessing nonce for secret key fails", 1),
		table.Entry("when accessing nonce for encrypted validation codes fails", 2),
	)

	table.DescribeTable("when decrypting returns an error", func(errorIndex int) {
		fakeDecryptor = &cryptofakes.FakeDecryptor{}
		fakeKeyService.KeyReturns(fakeDecryptor)
		var errorStr = "some error" + time.Now().String()
		fakeDecryptor.DecryptReturnsOnCall(errorIndex, "", errors.New(errorStr))

		updatedCredential, rotatorError = rotator.Rotate(
			entity.MfaCredential{
				EncryptionKeyLabel:      "key-1",
				ScratchCodes:            scratchCodes,
				SecretKey:               secretKey,
				EncryptedValidationCode: encryptedValidationCode,
			},
			fakeKeyService,
			fakeSaltAccessor,
			fakeNonceAccessor,
			fakeDbMapper,
		)

		Expect(rotatorError).To(HaveOccurred())
		Expect(rotatorError).To(MatchError("unable to decrypt cipher value provided: " + errorStr))
	},
		table.Entry("when decrypting scratch codes fails", 0),
		table.Entry("when decrypting secret key fails", 1),
		table.Entry("when decrypting encrypted validation codes fails", 2),
	)

	table.DescribeTable("when encrypting returns an error", func(errorIndex int) {
		fakeEncryptor = &cryptofakes.FakeEncryptor{}
		fakeKeyService.ActiveKeyReturns(activeKeyLabel, fakeEncryptor)

		var errorStr = "some error" + time.Now().String()
		fakeEncryptor.EncryptReturnsOnCall(errorIndex, crypto.EncryptedValue{}, errors.New(errorStr))

		updatedCredential, rotatorError = rotator.Rotate(
			entity.MfaCredential{
				EncryptionKeyLabel:      "key-1",
				ScratchCodes:            scratchCodes,
				SecretKey:               secretKey,
				EncryptedValidationCode: encryptedValidationCode,
			},
			fakeKeyService,
			fakeSaltAccessor,
			fakeNonceAccessor,
			fakeDbMapper,
		)

		Expect(rotatorError).To(HaveOccurred())
		Expect(rotatorError).To(MatchError("unable to encrypt value provided: " + errorStr))
	},
		table.Entry("when encrypting scratch codes fails", 0),
		table.Entry("when encrypting secret key fails", 1),
		table.Entry("when encrypting encrypted validation codes fails", 2),
	)
})
