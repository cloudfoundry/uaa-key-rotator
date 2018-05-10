package db_test

import (
	"testing"

	"database/sql"
	db2 "github.com/cloudfoundry/uaa-key-rotator/db"
	"github.com/cloudfoundry/uaa-key-rotator/db/testutils"
	"github.com/cloudfoundry/uaa-key-rotator/entity"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Db Suite")
}

var db *sqlx.DB

var _ = BeforeSuite(func() {
	By("validating and testing the db connection", func() {
		db = testutils.TestDBConnection()
	})

	By("migrating UAA database", testutils.MigrateUaaDatabase)
})

func insertGoogleMfaCredential(userId string, activeKeyLabel string) entity.MfaCredential {
	mfaCredential := entity.MfaCredential{
		UserId:                  userId,
		SecretKey:               "secret-key",
		ScratchCodes:            "scratch_codes",
		MfaProviderId:           "mfa_provider_id",
		ZoneId:                  "zone_id",
		EncryptionKeyLabel:      activeKeyLabel,
		EncryptedValidationCode: "encrypted_validation_code",
		ValidationCode:          sql.NullInt64{Int64: 1234, Valid: true},
	}

	insertSQL, err := db2.RebindForSQLDialect(`insert into user_google_mfa_credentials(
		user_id, 
		secret_key, 
		validation_code, 
		scratch_codes, 
		mfa_provider_id, 
		zone_id, 
		encryption_key_label, 
		encrypted_validation_code) values(
		?, ?, ?, ?, ?, ?, ?, ?
		)`,
		testutils.Scheme)
	Expect(err).NotTo(HaveOccurred())

	insertResult, err := db.Exec(insertSQL, mfaCredential.UserId,
		mfaCredential.SecretKey,
		mfaCredential.ValidationCode,
		mfaCredential.ScratchCodes,
		string(mfaCredential.MfaProviderId),
		string(mfaCredential.ZoneId),
		mfaCredential.EncryptionKeyLabel,
		mfaCredential.EncryptedValidationCode)

	Expect(err).NotTo(HaveOccurred())
	numOfRowsInserted, err := insertResult.RowsAffected()
	Expect(err).NotTo(HaveOccurred())
	Expect(numOfRowsInserted).To(Equal(int64(1)))
	return mfaCredential
}
