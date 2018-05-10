package db

import (
	"github.com/cloudfoundry/uaa-key-rotator/entity"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var updateGoogleMfaCredentialQuery = `update
user_google_mfa_credentials
set secret_key = ?,
scratch_codes = ?,
encryption_key_label = ?,
encrypted_validation_code = ?
where user_id = ?`

type DbAwareQuerier struct {
	DB       *sqlx.DB
	DBScheme string
}

func (q DbAwareQuerier) Close() error {
	return q.DB.Close()
}

func (q DbAwareQuerier) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	reboundQuery, err := RebindForSQLDialect(query, q.DBScheme)
	if err != nil {
		panic(err)
	}
	return q.DB.Queryx(reboundQuery, args...)
}

type GoogleMfaCredentialsDBUpdater struct {
	DB Queryer
}

func (gdb GoogleMfaCredentialsDBUpdater) Write(credential entity.MfaCredential) error {
	rs, err := gdb.DB.Queryx(updateGoogleMfaCredentialQuery,
		credential.SecretKey,
		credential.ScratchCodes,
		credential.EncryptionKeyLabel,
		credential.EncryptedValidationCode,
		credential.UserId,
	)
	if err != nil {
		return errors.Wrap(err, "Unable to update mfa db record")
	}
	defer rs.Close()
	return nil
}
