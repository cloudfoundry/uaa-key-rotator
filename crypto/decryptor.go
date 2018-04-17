package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

//go:generate counterfeiter . CipherSaltAccessor
type CipherSaltAccessor interface {
	GetSalt([]byte) ([]byte, error)
}

//go:generate counterfeiter . CipherNonceAccessor
type CipherNonceAccessor interface {
	GetNonce([]byte) ([]byte, error)
}

type Decryptor struct {
	Passphrase          string
	CipherSaltAccessor  CipherSaltAccessor
	CipherNonceAccessor CipherNonceAccessor
}

func (d Decryptor) Decrypt(cipherValue []byte) (string, error) {
	if len(cipherValue) == 0 {
		return "", errors.New("Unable to decrypt due to empty CipherText.")
	}
	salt, err := d.CipherSaltAccessor.GetSalt(cipherValue)

	aes, err := aes.NewCipher(GenerateKey(salt, d.Passphrase))
	if err != nil {
		return "", err
	}

	aesGcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}

	nonce, err := d.CipherNonceAccessor.GetNonce(cipherValue)
	if err != nil {
		return "", err
	}
	plainText, err := aesGcm.Open(nil, nonce, cipherValue, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}