package rotator

import (
	"github.com/cloudfoundry/uaa-key-rotator/entity"
	"github.com/cloudfoundry/uaa-key-rotator/crypto"
	"github.com/pkg/errors"
)

//go:generate counterfeiter . KeyService
type KeyService interface {
	Key(keyLabel string) crypto.Decryptor
	ActiveKey() (string, crypto.Encryptor)
}

//go:generate counterfeiter . MapEncryptedValueToDB
type MapEncryptedValueToDB interface {
	Map(value crypto.EncryptedValue) ([]byte, error)
}

type UAARotator struct {
	KeyService    KeyService
	SaltAccessor  crypto.CipherSaltAccessor
	NonceAccessor crypto.CipherNonceAccessor
	DbMapper      MapEncryptedValueToDB
}

func (r UAARotator) Rotate(credential entity.MfaCredential) (entity.MfaCredential, error) {
	decryptor := r.KeyService.Key(credential.EncryptionKeyLabel)
	activeKeyLabel, encryptor := r.KeyService.ActiveKey()

	rotatedScratchCodes, err := r.rotate(encryptor, decryptor, []byte(credential.ScratchCodes))
	if err != nil {
		return entity.MfaCredential{}, err
	}

	rotatedSecretKey, err := r.rotate(encryptor, decryptor, []byte(credential.SecretKey))
	if err != nil {
		return entity.MfaCredential{}, err
	}

	rotatedEncryptedValidationCode, err := r.rotate(encryptor, decryptor, []byte(credential.EncryptedValidationCode))
	if err != nil {
		return entity.MfaCredential{}, err
	}

	credential.ScratchCodes = string(rotatedScratchCodes)
	credential.SecretKey = string(rotatedSecretKey)
	credential.EncryptedValidationCode = string(rotatedEncryptedValidationCode)
	credential.EncryptionKeyLabel = activeKeyLabel

	return credential, nil
}

func (r UAARotator) rotate(encryptor crypto.Encryptor, decryptor crypto.Decryptor, cipherValue []byte) ([]byte, error) {
	salt, err := r.getSalt(cipherValue)
	if err != nil {
		return nil, err
	}

	nonce, err := r.getNonce(cipherValue)
	if err != nil {
		return nil, err
	}

	decryptedValue, err := r.decrypt(decryptor, cipherValue, salt, nonce)
	if err != nil {
		return nil, err
	}

	reEncryptedValue, err := r.encrypt(encryptor, decryptedValue)
	if err != nil {
		return nil, err
	}

	return r.DbMapper.Map(reEncryptedValue)
}

func (r UAARotator) encrypt(activeKey crypto.Encryptor, decryptedValue string) (crypto.EncryptedValue, error) {
	reEncryptedValue, err := activeKey.Encrypt(decryptedValue)
	if err != nil {
		return crypto.EncryptedValue{}, errors.Wrap(err, "unable to encrypt value provided")
	}
	return reEncryptedValue, nil
}

func (r UAARotator) decrypt(decryptor crypto.Decryptor, cipherValue []byte, salt []byte, nonce []byte) (string, error) {
	decrpytedValue, err := decryptor.Decrypt(
		crypto.EncryptedValue{
			CipherValue: []byte(cipherValue),
			Salt:        salt,
			Nonce:       nonce,
		})
	if err != nil {
		return "", errors.Wrap(err, "unable to decrypt cipher value provided")

	}
	return decrpytedValue, nil
}

func (r UAARotator) getSalt(cipherValue []byte) ([]byte, error) {
	scratchCodeSalt, err := r.SaltAccessor.GetSalt(cipherValue)
	if err != nil {
		return nil, errors.Wrap(err, "unable to access salt from cipher value provided")
	}
	return scratchCodeSalt, err
}

func (r UAARotator) getNonce(cipherValue []byte) ([]byte, error) {
	nonce, err := r.NonceAccessor.GetNonce(cipherValue)
	if err != nil {
		return nil, errors.Wrap(err, "unable to access nonce from cipher value provided")
	}
	return nonce, nil

}
