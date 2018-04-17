package crypto

import (
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha256"
)


const (
	SaltSize            = 32
	NonceSize           = 12
	AES256KeyLength     = 32
	ShaNumberIterations = 65536
)

func GenerateKey(salt []byte, passphrase string) []byte {
	return pbkdf2.Key([]byte(passphrase), salt, ShaNumberIterations, AES256KeyLength, sha256.New)
}
