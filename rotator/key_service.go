package rotator

import (
	"github.com/cloudfoundry/uaa-key-rotator/crypto"
	"github.com/pkg/errors"
	"fmt"
)

type EncryptionKey struct {
	Label      string
	Passphrase string
}

type UaaKeyService struct {
	ActiveKeyLabel string
	EncryptionKeys []EncryptionKey
}
var _ KeyService = UaaKeyService{}

func (s UaaKeyService) Key(keyLabel string) (crypto.Decryptor, error) {
	var key EncryptionKey
	var found bool

	if found, key = s.getEncryptionKey(keyLabel); !found {
		return crypto.UAADecryptor{}, errors.New(fmt.Sprintf("unable to find key: %s", keyLabel))

	}

	return crypto.UAADecryptor{
		Passphrase: key.Passphrase,
	}, nil
}

func (s UaaKeyService) ActiveKey() (string, crypto.Encryptor, error) {
	var key EncryptionKey
	var found bool

	if found, key = s.getEncryptionKey(s.ActiveKeyLabel); !found {
		return "", nil, errors.New(fmt.Sprintf("unable to find active key: %s", s.ActiveKeyLabel))
	}

	return s.ActiveKeyLabel, crypto.UAAEncryptor{
		Passphrase:     key.Passphrase,
		SaltGenerator:  crypto.UaaSaltGenerator{},
		NonceGenerator: crypto.UaaNonceGenerator{},
	}, nil
}

func (s UaaKeyService) getEncryptionKey(label string) (bool, EncryptionKey) {
	for _, key := range s.EncryptionKeys {
		if key.Label == label {
			return true, key
		}
	}

	return false, EncryptionKey{}
}
