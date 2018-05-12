package db

import (
	"github.com/cloudfoundry/uaa-key-rotator/entity"
	"github.com/pkg/errors"
)

var updateGoogleMfaCredentialQuery = `update
user_google_mfa_credentials
set secret_key = ?,
scratch_codes = ?,
encryption_key_label = ?,
encrypted_validation_code = ?
where user_id = ?`

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
