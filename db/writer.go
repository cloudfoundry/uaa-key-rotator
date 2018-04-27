package db

import (
	"fmt"
	"github.com/cloudfoundry/uaa-key-rotator/entity"
	"github.com/pkg/errors"
	"strings"
	"github.com/jmoiron/sqlx"
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
	return q.DB.Queryx(RebindForSQLDialect(query, q.DBScheme), args...)
}

func Write(db Queryer, credential entity.MfaCredential) error {
	rs, err := db.Queryx(updateGoogleMfaCredentialQuery,
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

func RebindForSQLDialect(query, dialect string) string {
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
