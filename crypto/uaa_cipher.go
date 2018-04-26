package crypto

import "errors"

//go:generate counterfeiter . CipherAccessor
type CipherAccessor interface {
	GetCipher([]byte) ([]byte, error)
}

type UAACipherAccessor struct{}

func (a UAACipherAccessor) GetCipher(cipher []byte) ([]byte, error) {
	if len(cipher) < 45 {
		return nil, errors.New("uaa cipher should be at least 44 bytes in length")
	}
	return cipher[44:], nil
}
