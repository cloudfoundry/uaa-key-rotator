package crypto

import "github.com/pkg/errors"
import (
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha256"
	"crypto/aes"
	"crypto/cipher"
)

type Encryptor struct {
	Passphrase string
	Salt       []byte
	Nonce      []byte
}

const (
	SaltSize            = 32
	NonceSize           = 12
	AES256KeyLength     = 32
	ShaNumberIterations = 65536
)

func (e Encryptor) Encrypt(plainText string) ([]byte, error) {
	if len(e.Salt) != SaltSize {
		return []byte(""), errors.New("Salt should be exactly 32 bytes in length.")
	}

	if len(e.Nonce) != NonceSize {
		return []byte(""), errors.New("Nonce should be exactly 12 bytes in length.")
	}

	key, err := e.generateKey()
	if err != nil {
		return nil, err
	}

	aes, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, err
	}

	cipherValue := aesGcm.Seal(nil, e.Nonce, []byte(plainText), nil)

	result := []byte{}
	result = append(result, e.Salt...)
	result = append(result, e.Nonce...)
	result = append(result, cipherValue...)

	return result, nil
}

func (e Encryptor) Decrypt(cipherValue []byte) (string, error) {
	if len(cipherValue) <= NonceSize + SaltSize {
		return "", errors.New("Unable to decrypt due to invalid CipherText.")
	}

	key, err := e.generateKey()
	if err != nil {
		return "", err
	}

	aes, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}

	plainText, err := aesGcm.Open(nil, e.Nonce, cipherValue[NonceSize+SaltSize:], nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

func (e Encryptor) generateKey() ([]byte, error) {
	if len(e.Salt) != SaltSize {
		return nil, errors.New("Salt should be exactly 32 bytes in length.")
	}

	return pbkdf2.Key([]byte(e.Passphrase), e.Salt, ShaNumberIterations, AES256KeyLength, sha256.New), nil
}
