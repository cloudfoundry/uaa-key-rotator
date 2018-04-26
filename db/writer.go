package db

import (
	"database/sql"
	"fmt"
	"github.com/cloudfoundry/uaa-key-rotator/entity"
	"github.com/pkg/errors"
	"strings"
)

var updateGoogleMfaCredentialQuery = `update
user_google_mfa_credentials
set secret_key = ?,
scratch_codes = ?,
encryption_key_label = ?,
encrypted_validation_code = ?
where user_id = ?`

//go:generate counterfeiter . Updater
type Updater interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func Write(db Updater, credential entity.MfaCredential) error {
	rs, err := db.Query(rebindForSQLDialect(updateGoogleMfaCredentialQuery, "postgres"),
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

func rebindForSQLDialect(query, dialect string) string {
	if dialect == "mysql" {
		return query
	}
	if dialect != "postgres" {
		panic(fmt.Sprintf("Unrecognized DB dialect '%s'", dialect))
	}

	strParts := strings.Split(query, "?")
	for i := 1; i < len(strParts); i++ {
		strParts[i-1] = fmt.Sprintf("%s$%d", strParts[i-1], i)
	}
	return strings.Join(strParts, "")
}
