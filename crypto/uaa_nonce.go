package crypto

import (
	"crypto/rand"
	"errors"
)

//go:generate counterfeiter . CipherNonceAccessor
type CipherNonceAccessor interface {
	GetNonce([]byte) ([]byte, error)
}

type UaaNonceGenerator struct{}

func (UaaNonceGenerator) GetNonce() ([]byte, error) {
	var nonce = make([]byte, 12)
	_, err := rand.Read(nonce)
	return nonce, err
}

type UaaNonceAccessor struct{}

func (UaaNonceAccessor) GetNonce(cipher []byte) ([]byte, error) {
	if len(cipher) < 12 {
		return nil, errors.New("nonce should be exactly 12 bytes in length")
	}
	return cipher[:12], nil
}
