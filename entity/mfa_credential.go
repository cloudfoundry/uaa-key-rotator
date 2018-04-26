package entity

import (
	"bytes"
	"database/sql"
	"errors"
	"strings"
)

type MfaCredential struct {
	UserId                  string        `db:"user_id"`
	MfaProviderId           Char          `db:"mfa_provider_id"`
	ValidationCode          sql.NullInt64 `db:"validation_code"`
	ScratchCodes            string        `db:"scratch_codes"`
	SecretKey               string        `db:"secret_key"`
	EncryptionKeyLabel      string        `db:"encryption_key_label"`
	EncryptedValidationCode string        `db:"encrypted_validation_code"`
	ZoneId                  Char          `db:"zone_id"`
}

type Char string

func (g *Char) Scan(src interface{}) error {
	switch src.(type) {
	case string:
		s := src.(string)
		*g = Char(strings.Trim(s, " "))
	case []byte:
		b := src.([]byte)
		*g = Char(bytes.Trim(b, " "))
	default:
		return errors.New("Incompatible type for Char")
	}
	return nil
}
