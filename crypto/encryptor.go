package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

//go:generate counterfeiter . SaltGenerator
type SaltGenerator interface {
	GetSalt() ([]byte, error)
}

//go:generate counterfeiter . NonceGenerator
type NonceGenerator interface {
	GetNonce() ([]byte, error)
}

type UAAEncryptor struct {
	Passphrase     string
	SaltGenerator  SaltGenerator
	NonceGenerator NonceGenerator
}

type EncryptedValue struct {
	Salt        []byte
	Nonce       []byte
	CipherValue []byte
}

//go:generate counterfeiter . Encryptor
type Encryptor interface {
	Encrypt(plainText string) (EncryptedValue, error)
}

func (e UAAEncryptor) Encrypt(plainText string) (EncryptedValue, error) {
	salt, err := e.SaltGenerator.GetSalt()
	if err != nil {
		//TODO: handle error
		panic(err)
	}
	nonce, err := e.NonceGenerator.GetNonce()
	if err != nil {
		//TODO: handle error
		panic(err)
	}

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
