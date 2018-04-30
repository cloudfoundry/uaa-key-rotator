package db

import (
	"github.com/cloudfoundry/uaa-key-rotator/entity"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

//go:generate counterfeiter . Queryer
type Queryer interface {
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	Close() error
}

type GoogleMfaCredentialsDB struct {
	DB             Queryer
	ActiveKeyLabel string
}

func (gdb GoogleMfaCredentialsDB) RowsToRotate() (<-chan entity.MfaCredential, <-chan error) {
	var mfaCredentialChan = make(chan entity.MfaCredential)
	var errChan = make(chan error)

	go func() {
		rows, err := gdb.DB.Queryx(
			`select user_id, mfa_provider_id, zone_id, validation_code, scratch_codes, encryption_key_label, encrypted_validation_code, secret_key
		 from user_google_mfa_credentials
         where encryption_key_label <> ?`,
			gdb.ActiveKeyLabel)
		if err != nil {
			errChan <- errors.Wrap(err, "RowsToRotate failed to query table")
			return
		}

		defer rows.Close() // untested

		for rows.Next() {
			mfaCredential := entity.MfaCredential{}
			err = rows.StructScan(&mfaCredential)
			if err != nil {
				errChan <- errors.Wrap(err, "Unable to deserialize db response")
				return
			}
			mfaCredentialChan <- mfaCredential
		}
	}()

	return mfaCredentialChan, errChan
}
