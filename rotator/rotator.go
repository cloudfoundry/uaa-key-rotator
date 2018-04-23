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
	Map(value crypto.EncryptedValue) []byte
}

func Rotate(credential entity.MfaCredential,
	keyService KeyService,
	saltAccessor crypto.CipherSaltAccessor,
	nonceAccessor crypto.CipherNonceAccessor,
	dbMapper MapEncryptedValueToDB,
) (entity.MfaCredential, error) {
	decryptor := keyService.Key(credential.EncryptionKeyLabel)
	activeKeyLabel, activeKey := keyService.ActiveKey()

	rotatedScratchCodes, err := rotate(activeKey, decryptor, saltAccessor, nonceAccessor, dbMapper, []byte(credential.ScratchCodes))
	if err != nil {
		return entity.MfaCredential{}, err
	}

	rotatedSecretKey, err := rotate(activeKey, decryptor, saltAccessor, nonceAccessor, dbMapper, []byte(credential.SecretKey))
	if err != nil {
		return entity.MfaCredential{}, err
	}

	rotatedEncryptedValidationCode, err := rotate(activeKey, decryptor, saltAccessor, nonceAccessor, dbMapper, []byte(credential.EncryptedValidationCode))
	if err != nil {
		return entity.MfaCredential{}, err
	}

	credential.ScratchCodes = string(rotatedScratchCodes)
	credential.SecretKey = string(rotatedSecretKey)
	credential.EncryptedValidationCode = string(rotatedEncryptedValidationCode)
	credential.EncryptionKeyLabel = activeKeyLabel

	return credential, nil
}

func rotate(activeKey crypto.Encryptor,
	decryptor crypto.Decryptor,
	saltAccessor crypto.CipherSaltAccessor,
	nonceAccessor crypto.CipherNonceAccessor,
	dbMapper MapEncryptedValueToDB,
	cipherValue []byte) ([]byte, error) {

	salt, err := getSalt(saltAccessor, cipherValue)
	if err != nil {
		return nil, err
	}

	nonce, err := getNonce(nonceAccessor, cipherValue)
	if err != nil {
		return nil, err
	}

	decryptedValue, err := decrypt(decryptor, cipherValue, salt, nonce)
	if err != nil {
		return nil, err
	}

	reEncryptedValue, err := encrypt(activeKey, decryptedValue)
	if err != nil {
		return nil, err
	}

	return dbMapper.Map(reEncryptedValue), nil
}

func encrypt(activeKey crypto.Encryptor, decryptedValue string) (crypto.EncryptedValue, error) {
	reEncryptedValue, err := activeKey.Encrypt(decryptedValue)
	if err != nil {
		return crypto.EncryptedValue{}, errors.Wrap(err, "unable to encrypt value provided")
	}
	return reEncryptedValue, nil
}

func decrypt(decryptor crypto.Decryptor, cipherValue []byte, salt []byte, nonce []byte) (string, error) {
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

func getSalt(saltAccessor crypto.CipherSaltAccessor, cipherValue []byte) ([]byte, error) {
	scratchCodeSalt, err := saltAccessor.GetSalt(cipherValue)
	if err != nil {
		return nil, errors.Wrap(err, "unable to access salt from cipher value provided")
	}
	return scratchCodeSalt, err
}

func getNonce(nonceAccessor crypto.CipherNonceAccessor, cipherValue []byte) ([]byte, error) {
	nonce, err := nonceAccessor.GetNonce(cipherValue)
	if err != nil {
		return nil, errors.Wrap(err, "unable to access nonce from cipher value provided")
	}
	return nonce, nil

}
