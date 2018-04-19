package db

import (
	"strings"
	"bytes"
	"github.com/pkg/errors"
	"github.com/jmoiron/sqlx"
)

type MfaCredential struct {
	UserId                  string `db:"user_id"`
	MfaProviderId           Char   `db:"mfa_provider_id"`
	ValidationCode          int    `db:"validation_code"`
	ScratchCodes            string `db:"scratch_codes"`
	EncryptionKeyLabel      string `db:"encryption_key_label"`
	EncryptedValidationCode string `db:"encrypted_validation_code"`
	ZoneId                  Char   `db:"zone_id"`
}

//go:generate counterfeiter . Queryer
type Queryer interface {
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
}

func ReadAll(db Queryer) ([]MfaCredential, error) {
	rows, err := db.Queryx(`select user_id, mfa_provider_id, zone_id, validation_code, scratch_codes, encryption_key_label, encrypted_validation_code from user_google_mfa_credentials`)
	if err != nil {
		return nil, errors.Wrap(err, "ReadAll failed to query table")
	}
	defer rows.Close() // untested

	var mfaCredentials []MfaCredential

	for rows.Next() {
		mfaCredential := MfaCredential{}
		err = rows.StructScan(&mfaCredential)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to deserialize db response") //untested
		}
		mfaCredentials = append(mfaCredentials, mfaCredential)
	}

	return mfaCredentials, nil
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
