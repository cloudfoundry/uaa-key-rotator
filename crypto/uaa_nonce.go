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

func (UaaNonceGenerator) GetNonce() []byte {
	var nonce = make([]byte, 12)
	_, err := rand.Read(nonce)
	if err != nil {
		//TODO: return an error
		panic(err)
	}
	return nonce
}

type UaaNonceAccessor struct{}

func (UaaNonceAccessor) GetNonce(cipher []byte) ([]byte, error) {
	if len(cipher) < 12 {
		return nil, errors.New("nonce should be exactly 12 bytes in length")
	}
	return cipher[:12], nil
}
