package db_test

import (
	"database/sql"
	"errors"
	. "github.com/cloudfoundry/uaa-key-rotator/db"
	"github.com/cloudfoundry/uaa-key-rotator/db/dbfakes"
	"github.com/cloudfoundry/uaa-key-rotator/entity"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cloudfoundry/uaa-key-rotator/db/testutils"
	"time"
)

var _ = Describe("Postgresql", func() {
	var googleMfaCredentialsDB GoogleMfaCredentialsDBFetcher

	BeforeEach(func() {
		deleteResult, err := db.Exec(`delete from user_google_mfa_credentials`)
		Expect(err).NotTo(HaveOccurred())
		numOfRowsDeleted, err := deleteResult.RowsAffected()
		Expect(err).NotTo(HaveOccurred())
		Expect(numOfRowsDeleted).To(BeNumerically(">=", int64(0)))

		insertGoogleMfaCredential("1", "not-activeKeyLabel")
		insertGoogleMfaCredential("2", "not-activeKeyLabel")
		insertGoogleMfaCredential("3", "activeKeyLabel")
		insertGoogleMfaCredential("4", "activeKeyLabel")

		googleMfaCredentialsDB = GoogleMfaCredentialsDBFetcher{
			DB:             DbAwareQuerier{DB: db, DBScheme: testutils.Scheme},
			ActiveKeyLabel: "activeKeyLabel",
		}
	})

	It("should return every record (except active key) from the user_google_mfa_credentials table", func() {
		var mfaCredentials <-chan entity.MfaCredential
		var errChan <-chan error

		mfaCredentials, errChan = googleMfaCredentialsDB.RowsToRotate()
		Consistently(errChan).ShouldNot(Receive())

		var mfaCredential entity.MfaCredential
		Eventually(mfaCredentials, 5*time.Second).Should(Receive(&mfaCredential))
		Expect(mfaCredential).To(Equal(
			entity.MfaCredential{
				UserId:                  "1",
				MfaProviderId:           entity.Char("mfa_provider_id"),
				ZoneId:                  entity.Char("zone_id"),
				ValidationCode:          sql.NullInt64{Int64: 1234, Valid: true},
				ScratchCodes:            "scratch_codes",
				SecretKey:               "secret-key",
				EncryptionKeyLabel:      "not-activeKeyLabel",
				EncryptedValidationCode: "encrypted_validation_code",
			},
		))

		Eventually(mfaCredentials).Should(Receive(&mfaCredential))
		Expect(mfaCredential).To(Equal(
			entity.MfaCredential{
				UserId:                  "2",
				MfaProviderId:           entity.Char("mfa_provider_id"),
				ZoneId:                  entity.Char("zone_id"),
				ValidationCode:          sql.NullInt64{Int64: 1234, Valid: true},
				SecretKey:               "secret-key",
				ScratchCodes:            "scratch_codes",
				EncryptionKeyLabel:      "not-activeKeyLabel",
				EncryptedValidationCode: "encrypted_validation_code",
			},
		))

	})

	Describe("FakeDB", func() {
		var queryer *dbfakes.FakeQueryer

		Context("error during querying mfa table", func() {
			BeforeEach(func() {
				queryer = &dbfakes.FakeQueryer{}
				queryer.QueryxReturns(nil, errors.New("cannot query table"))
				googleMfaCredentialsDB = GoogleMfaCredentialsDBFetcher{
					DB:             queryer,
					ActiveKeyLabel: "activeKeyLabel",
				}
			})

			It("should return a meaningful error", func() {
				_, errChan := googleMfaCredentialsDB.RowsToRotate()
				var err error
				Eventually(errChan).Should(Receive(&err))
				Expect(err).To(MatchError("RowsToRotate failed to query table: cannot query table"))
			})
		})
	})

})
