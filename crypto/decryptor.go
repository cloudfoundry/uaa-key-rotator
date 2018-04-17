package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

type Decryptor struct {
	Passphrase          string
}

func (d Decryptor) Decrypt(encryptedValue EncryptedValue) (string, error) {
	if len(encryptedValue.CipherValue) == 0 {
		return "", errors.New("unable to decrypt due to empty CipherText")
	}

	aes, err := aes.NewCipher(GenerateKey(encryptedValue.Salt, d.Passphrase))
	if err != nil {
		return "", err
	}

	aesGcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}

	plainText, err := aesGcm.Open(nil, encryptedValue.Nonce, encryptedValue.CipherValue, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}