package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)


//go:generate counterfeiter . SaltGenerator
type SaltGenerator interface {
	GetSalt() []byte
}

//go:generate counterfeiter . NonceGenerator
type NonceGenerator interface {
	GetNonce() []byte
}

type Encryptor struct {
	Passphrase     string
	SaltGenerator  SaltGenerator
	NonceGenerator NonceGenerator
}

type EncryptedValue struct {
	Salt        []byte
	Nonce       []byte
	CipherValue []byte
}

func (e Encryptor) Encrypt(plainText string) (EncryptedValue, error) {
	salt := e.SaltGenerator.GetSalt()
	nonce := e.NonceGenerator.GetNonce()
	aes, err := aes.NewCipher(GenerateKey(salt, e.Passphrase))
	if err != nil {
		return EncryptedValue{}, err
	}

	aesGcm, err := cipher.NewGCM(aes)
	if err != nil {
		return EncryptedValue{}, err
	}

	cipherValue := aesGcm.Seal(nil, nonce, []byte(plainText), nil)
	return EncryptedValue{salt, nonce, cipherValue}, nil
}

