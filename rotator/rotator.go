package rotator

import (
	"github.com/cloudfoundry/uaa-key-rotator/crypto"
	"github.com/cloudfoundry/uaa-key-rotator/entity"
	"github.com/pkg/errors"
)

//go:generate counterfeiter . KeyService
type KeyService interface {
	Key(keyLabel string) (crypto.Decryptor, error)
	ActiveKey() (string, crypto.Encryptor, error)
}

//go:generate counterfeiter . MapEncryptedValueToDB
type MapEncryptedValueToDB interface {
	Map(value crypto.EncryptedValue) ([]byte, error)
	MapBase64ToCipherValue(value string) ([]byte, error)
}

type UAARotator struct {
	KeyService     KeyService
	SaltAccessor   crypto.CipherSaltAccessor
	NonceAccessor  crypto.CipherNonceAccessor
	CipherAccessor crypto.CipherAccessor
	DbMapper       MapEncryptedValueToDB
}

func (r UAARotator) Rotate(credential entity.MfaCredential) (entity.MfaCredential, error) {
	decryptor, err := r.KeyService.Key(credential.EncryptionKeyLabel)
	if err != nil {
		return entity.MfaCredential{}, errors.Wrap(err, "Unable to decrypt mfa record")
	}

	activeKeyLabel, encryptor, err := r.KeyService.ActiveKey()
	if err != nil {
		return entity.MfaCredential{}, errors.Wrap(err, "Unable to decrypt mfa record")
	}

	rotatedScratchCodes, err := r.rotateCipherValue(encryptor, decryptor, credential.ScratchCodes)
	if err != nil {
		return entity.MfaCredential{}, err
	}

	rotatedSecretKey, err := r.rotateCipherValue(encryptor, decryptor, credential.SecretKey)
	if err != nil {
		return entity.MfaCredential{}, err
	}

	rotatedEncryptedValidationCode, err := r.rotateCipherValue(encryptor, decryptor, credential.EncryptedValidationCode)
	if err != nil {
		return entity.MfaCredential{}, err
	}

	credential.ScratchCodes = string(rotatedScratchCodes)
	credential.SecretKey = string(rotatedSecretKey)
	credential.EncryptedValidationCode = string(rotatedEncryptedValidationCode)
	credential.EncryptionKeyLabel = activeKeyLabel

	return credential, nil
}

func (r UAARotator) rotateCipherValue(encryptor crypto.Encryptor, decryptor crypto.Decryptor, encodedCipherValue string) ([]byte, error) {
	uaaCipherValue, err := r.base64DecodeCipher(encodedCipherValue)
	if err != nil {
		return nil, err
	}

	salt, err := r.getSalt(uaaCipherValue)
	if err != nil {
		return nil, err
	}

	nonce, err := r.getNonce(uaaCipherValue)
	if err != nil {
		return nil, err
	}

	cipherValue, err := r.getCipherValue(uaaCipherValue)
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

func (r UAARotator) getCipherValue(uaaCipherValue []byte) ([]byte, error) {
	cipherValue, err := r.CipherAccessor.GetCipher(uaaCipherValue)
	if err != nil {
		return nil, errors.Wrap(err, "unable to access cipher value from 'uaa' cipher value provided")
	}
	return cipherValue, nil
}

func (r UAARotator) base64DecodeCipher(cipher string) ([]byte, error) {
	scratchCodes, err := r.DbMapper.MapBase64ToCipherValue(cipher)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to decode mfa credential value")
	}
	return scratchCodes, nil
}
