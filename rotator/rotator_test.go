package rotator_test

import (
	"database/sql"
	"github.com/cloudfoundry/uaa-key-rotator/crypto"
	"github.com/cloudfoundry/uaa-key-rotator/crypto/cryptofakes"
	"github.com/cloudfoundry/uaa-key-rotator/entity"
	"github.com/cloudfoundry/uaa-key-rotator/rotator"
	"github.com/cloudfoundry/uaa-key-rotator/rotator/rotatorfakes"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/pkg/errors"
	"time"
)

var _ = Describe("UAARotator", func() {
	var uaaRotator rotator.UAARotator
	var updatedCredential entity.MfaCredential
	var rotatorError error

	var scratchCodes string
	var secretKey string
	var encryptedValidationCode string

	var base64ScratchCodes string
	var base64SecretKey string
	var base64EncryptedValidationCode string

	var activeKeyLabel string
	var fakeKeyService *rotatorfakes.FakeKeyService

	var fakeDecryptor *cryptofakes.FakeDecryptor
	var fakeEncryptor *cryptofakes.FakeEncryptor

	var fakeSaltAccessor *cryptofakes.FakeCipherSaltAccessor
	var fakeNonceAccessor *cryptofakes.FakeCipherNonceAccessor
	var fakeCipherAccessor *cryptofakes.FakeCipherAccessor

	var fakeDbMapper *rotatorfakes.FakeMapEncryptedValueToDB

	var fakeDecrpytedScratchCodes string
	var fakeDecryptedSecretKey string
	var fakeDecryptedValidationCode string

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
		fakeCipherAccessor = &cryptofakes.FakeCipherAccessor{}

		base64ScratchCodes = "base64-encrypted-scratch-codes" + time.Now().String()
		scratchCodes = "encrypted-scratch-codes" + time.Now().String()
		base64SecretKey = "base64-secret-key" + time.Now().String()
		secretKey = "secret-key" + time.Now().String()
		base64EncryptedValidationCode = "base64-encrypted-validation-code" + time.Now().String()
		encryptedValidationCode = "encrypted-validation-code" + time.Now().String()

		fakeKeyService.KeyReturns(fakeDecryptor, nil)
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

		fakeCipherAccessor.GetCipherReturnsOnCall(0, []byte(scratchCodes), nil)
		fakeCipherAccessor.GetCipherReturnsOnCall(1, []byte(secretKey), nil)
		fakeCipherAccessor.GetCipherReturnsOnCall(2, []byte(encryptedValidationCode), nil)

		fakeEncryptor = &cryptofakes.FakeEncryptor{}
		activeKeyLabel = "key-2"
		fakeKeyService.ActiveKeyReturns(activeKeyLabel, fakeEncryptor, nil)

		fakeDbMapper = &rotatorfakes.FakeMapEncryptedValueToDB{}

		fakeDbMapper.MapBase64ToCipherValueReturnsOnCall(0, []byte(scratchCodes), nil)
		fakeDbMapper.MapBase64ToCipherValueReturnsOnCall(1, []byte(secretKey), nil)
		fakeDbMapper.MapBase64ToCipherValueReturnsOnCall(2, []byte(encryptedValidationCode), nil)

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

	Context("rotator is configured correctly", func() {
		JustBeforeEach(func() {
			uaaRotator = rotator.UAARotator{
				KeyService:     fakeKeyService,
				SaltAccessor:   fakeSaltAccessor,
				NonceAccessor:  fakeNonceAccessor,
				CipherAccessor: fakeCipherAccessor,
				DbMapper:       fakeDbMapper,
			}
			updatedCredential, rotatorError = uaaRotator.Rotate(
				entity.MfaCredential{
					UserId:                  "some-user-id",
					MfaProviderId:           "some-provider-id",
					ZoneId:                  "some-zone-id",
					EncryptionKeyLabel:      "key-1",
					ValidationCode:          sql.NullInt64{Int64: 1},
					ScratchCodes:            base64ScratchCodes,
					SecretKey:               base64SecretKey,
					EncryptedValidationCode: base64EncryptedValidationCode,
				},
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

	})

	Context("Attempting to rotate with an unknown key", func() {
		BeforeEach(func() {
			fakeKeyService.KeyReturns(nil, errors.New("Couldn't find key with label=key-1"))
		})

		JustBeforeEach(func() {
			uaaRotator = rotator.UAARotator{
				KeyService:     fakeKeyService,
				SaltAccessor:   fakeSaltAccessor,
				NonceAccessor:  fakeNonceAccessor,
				CipherAccessor: fakeCipherAccessor,
				DbMapper:       fakeDbMapper,
			}
			updatedCredential, rotatorError = uaaRotator.Rotate(
				entity.MfaCredential{
					UserId:                  "some-user-id",
					MfaProviderId:           "some-provider-id",
					ZoneId:                  "some-zone-id",
					EncryptionKeyLabel:      "key-1",
					ValidationCode:          sql.NullInt64{Int64: 1},
					ScratchCodes:            base64ScratchCodes,
					SecretKey:               base64SecretKey,
					EncryptedValidationCode: base64EncryptedValidationCode,
				},
			)
		})

		It("Should return a meaningful error", func() {
			Expect(rotatorError).To(HaveOccurred())
			Expect(rotatorError).To(MatchError("Unable to decrypt mfa record: Couldn't find key with label=key-1"))
		})
	})

	Context("Attempting to rotate with missing/invalid active key", func() {
		BeforeEach(func() {
			fakeKeyService.ActiveKeyReturns("", nil, errors.New("Configured active key is missing or invalid"))
		})

		JustBeforeEach(func() {
			uaaRotator = rotator.UAARotator{
				KeyService:     fakeKeyService,
				SaltAccessor:   fakeSaltAccessor,
				NonceAccessor:  fakeNonceAccessor,
				CipherAccessor: fakeCipherAccessor,
				DbMapper:       fakeDbMapper,
			}
			updatedCredential, rotatorError = uaaRotator.Rotate(
				entity.MfaCredential{
					UserId:                  "some-user-id",
					MfaProviderId:           "some-provider-id",
					ZoneId:                  "some-zone-id",
					EncryptionKeyLabel:      "key-1",
					ValidationCode:          sql.NullInt64{Int64: 1},
					ScratchCodes:            base64ScratchCodes,
					SecretKey:               base64SecretKey,
					EncryptedValidationCode: base64EncryptedValidationCode,
				},
			)
		})

		It("Should return a meaningful error", func() {
			Expect(rotatorError).To(HaveOccurred())
			Expect(rotatorError).To(MatchError("Unable to decrypt mfa record: Configured active key is missing or invalid"))
		})
	})

	table.DescribeTable("Attempting to base64 decode", func(errorIndex int) {
		fakeDbMapper.MapBase64ToCipherValueReturnsOnCall(errorIndex, []byte(scratchCodes), errors.New("some base64 decode error"))

		uaaRotator = rotator.UAARotator{
			KeyService:     fakeKeyService,
			SaltAccessor:   fakeSaltAccessor,
			NonceAccessor:  fakeNonceAccessor,
			CipherAccessor: fakeCipherAccessor,
			DbMapper:       fakeDbMapper,
		}
		updatedCredential, rotatorError = uaaRotator.Rotate(
			entity.MfaCredential{
				UserId:                  "some-user-id",
				MfaProviderId:           "some-provider-id",
				ZoneId:                  "some-zone-id",
				EncryptionKeyLabel:      "key-1",
				ValidationCode:          sql.NullInt64{Int64: 1},
				ScratchCodes:            base64ScratchCodes,
				SecretKey:               base64SecretKey,
				EncryptedValidationCode: base64EncryptedValidationCode,
			},
		)

		Expect(rotatorError).To(HaveOccurred())
		Expect(rotatorError).To(MatchError("Unable to decode mfa credential value: some base64 decode error"))
	},
		table.Entry("when base64 decoding ScratchCodes fails", 0),
		table.Entry("when base64 decoding SecretKey fails", 1),
		table.Entry("when base64 decoding EncryptedValidationCode fails", 2),
	)

	table.DescribeTable("when accessing the salt returns an error", func(errorIndex int) {
		fakeSaltAccessor = &cryptofakes.FakeCipherSaltAccessor{}
		var errorStr = "some error" + time.Now().String()
		fakeSaltAccessor.GetSaltReturnsOnCall(errorIndex, nil, errors.New(errorStr))

		uaaRotator = rotator.UAARotator{
			KeyService:     fakeKeyService,
			SaltAccessor:   fakeSaltAccessor,
			NonceAccessor:  fakeNonceAccessor,
			CipherAccessor: fakeCipherAccessor,
			DbMapper:       fakeDbMapper,
		}

		updatedCredential, rotatorError = uaaRotator.Rotate(
			entity.MfaCredential{
				EncryptionKeyLabel:      "key-1",
				ScratchCodes:            base64ScratchCodes,
				SecretKey:               base64SecretKey,
				EncryptedValidationCode: base64EncryptedValidationCode,
			},
		)

		Expect(rotatorError).To(HaveOccurred())
		Expect(rotatorError).To(MatchError("unable to access salt from cipher value provided: " + errorStr))
	},
		table.Entry("when accessing salt for scratch codes fails", 0),
		table.Entry("when accessing salt for secret key fails", 1),
		table.Entry("when accessing salt for encrypted validation codes fails", 2),
	)

	table.DescribeTable("when accessing the uaa cipher returns an error", func(errorIndex int) {
		fakeCipherAccessor = &cryptofakes.FakeCipherAccessor{}
		var errorStr = "some error" + time.Now().String()
		fakeCipherAccessor.GetCipherReturnsOnCall(errorIndex, nil, errors.New(errorStr))

		uaaRotator = rotator.UAARotator{
			KeyService:     fakeKeyService,
			SaltAccessor:   fakeSaltAccessor,
			NonceAccessor:  fakeNonceAccessor,
			CipherAccessor: fakeCipherAccessor,
			DbMapper:       fakeDbMapper,
		}

		updatedCredential, rotatorError = uaaRotator.Rotate(
			entity.MfaCredential{
				EncryptionKeyLabel:      "key-1",
				ScratchCodes:            base64ScratchCodes,
				SecretKey:               base64SecretKey,
				EncryptedValidationCode: base64EncryptedValidationCode,
			},
		)

		Expect(rotatorError).To(HaveOccurred())
		Expect(rotatorError).To(MatchError("unable to access cipher value from 'uaa' cipher value provided: " + errorStr))
	},
		table.Entry("when accessing uaa cipher for scratch codes fails", 0),
		table.Entry("when accessing uaa cipher for secret key fails", 1),
		table.Entry("when accessing uaa cipher for encrypted validation codes fails", 2),
	)

	table.DescribeTable("when accessing the nonce returns an error", func(errorIndex int) {
		fakeNonceAccessor = &cryptofakes.FakeCipherNonceAccessor{}
		var errorStr = "some error" + time.Now().String()
		fakeNonceAccessor.GetNonceReturnsOnCall(errorIndex, nil, errors.New(errorStr))

		uaaRotator = rotator.UAARotator{
			KeyService:     fakeKeyService,
			SaltAccessor:   fakeSaltAccessor,
			NonceAccessor:  fakeNonceAccessor,
			CipherAccessor: fakeCipherAccessor,
			DbMapper:       fakeDbMapper,
		}

		updatedCredential, rotatorError = uaaRotator.Rotate(
			entity.MfaCredential{
				EncryptionKeyLabel:      "key-1",
				ScratchCodes:            base64ScratchCodes,
				SecretKey:               base64SecretKey,
				EncryptedValidationCode: base64EncryptedValidationCode,
			},
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
		fakeKeyService.KeyReturns(fakeDecryptor, nil)
		var errorStr = "some error" + time.Now().String()
		fakeDecryptor.DecryptReturnsOnCall(errorIndex, "", errors.New(errorStr))

		uaaRotator = rotator.UAARotator{
			KeyService:     fakeKeyService,
			SaltAccessor:   fakeSaltAccessor,
			NonceAccessor:  fakeNonceAccessor,
			CipherAccessor: fakeCipherAccessor,
			DbMapper:       fakeDbMapper,
		}

		updatedCredential, rotatorError = uaaRotator.Rotate(
			entity.MfaCredential{
				EncryptionKeyLabel:      "key-1",
				ScratchCodes:            base64ScratchCodes,
				SecretKey:               base64SecretKey,
				EncryptedValidationCode: base64EncryptedValidationCode,
			},
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
		fakeKeyService.ActiveKeyReturns(activeKeyLabel, fakeEncryptor, nil)

		var errorStr = "some error" + time.Now().String()
		fakeEncryptor.EncryptReturnsOnCall(errorIndex, crypto.EncryptedValue{}, errors.New(errorStr))

		uaaRotator = rotator.UAARotator{
			KeyService:     fakeKeyService,
			SaltAccessor:   fakeSaltAccessor,
			NonceAccessor:  fakeNonceAccessor,
			CipherAccessor: fakeCipherAccessor,
			DbMapper:       fakeDbMapper,
		}

		updatedCredential, rotatorError = uaaRotator.Rotate(
			entity.MfaCredential{
				EncryptionKeyLabel:      "key-1",
				ScratchCodes:            scratchCodes,
				SecretKey:               secretKey,
				EncryptedValidationCode: encryptedValidationCode,
			},
		)

		Expect(rotatorError).To(HaveOccurred())
		Expect(rotatorError).To(MatchError("unable to encrypt value provided: " + errorStr))
	},
		table.Entry("when encrypting scratch codes fails", 0),
		table.Entry("when encrypting secret key fails", 1),
		table.Entry("when encrypting encrypted validation codes fails", 2),
	)
})
