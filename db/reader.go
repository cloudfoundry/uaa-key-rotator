package db

import (
	"github.com/cloudfoundry/uaa-key-rotator/entity"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

//go:generate counterfeiter .
// Queryer
type Queryer interface {
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
}

func ReadAll(db Queryer) ([]entity.MfaCredential, error) {
	rows, err := db.Queryx(`select user_id, mfa_provider_id, zone_id, validation_code, scratch_codes, encryption_key_label, encrypted_validation_code, secret_key from user_google_mfa_credentials`)
	if err != nil {
		return nil, errors.Wrap(err, "ReadAll failed to query table")
	}
	defer rows.Close() // untested

	var mfaCredentials []entity.MfaCredential

	for rows.Next() {
		mfaCredential := entity.MfaCredential{}
		err = rows.StructScan(&mfaCredential)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to deserialize db response") //untested
		}
		mfaCredentials = append(mfaCredentials, mfaCredential)
	}

	return mfaCredentials, nil
}
