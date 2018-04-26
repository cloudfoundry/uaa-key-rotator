package crypto

import (
	"crypto/sha256"
	"golang.org/x/crypto/pbkdf2"
)

const (
	AES256KeyLength     = 32
	ShaNumberIterations = 65536
)

func GenerateKey(salt []byte, passphrase string) []byte {
	return pbkdf2.Key([]byte(passphrase), salt, ShaNumberIterations, AES256KeyLength, sha256.New)
}
