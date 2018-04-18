package crypto

import (
	"crypto/rand"
	"errors"
)

//go:generate counterfeiter . CipherSaltAccessor
type CipherSaltAccessor interface {
	GetSalt([]byte) ([]byte, error)
}

type UaaSaltGenerator struct{}

func (UaaSaltGenerator) GetSalt() ([]byte, error) {
	var salt = make([]byte, 32)
	_, err := rand.Read(salt)
	return salt, err
}

type UaaSaltAccessor struct{}

func (UaaSaltAccessor) GetSalt(cipher []byte) ([]byte, error) {
	if len(cipher) < 45 {
		return nil, errors.New("cipher should be more than 45 bytes")
	}
	return cipher[12:44], nil
}
