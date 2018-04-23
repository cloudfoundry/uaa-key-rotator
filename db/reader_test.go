package db_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/cloudfoundry/uaa-key-rotator/db"
	"github.com/cloudfoundry/uaa-key-rotator/db/dbfakes"
	"errors"
	"database/sql"
	"github.com/cloudfoundry/uaa-key-rotator/entity"
)

var _ = Describe("Postgresql", func() {

	BeforeEach(func() {
		deleteResult, err := db.Exec(`delete from user_google_mfa_credentials`)
		Expect(err).NotTo(HaveOccurred())
		numOfRowsDeleted, err := deleteResult.RowsAffected()
		Expect(err).NotTo(HaveOccurred())
		Expect(numOfRowsDeleted).To(BeNumerically(">=", int64(0)))

		insertGoogleMfaCredential("1")
		insertGoogleMfaCredential("2")
	})

	It("should return every record from the user_google_mfa_credentials table", func() {
		var mfaCredentials []entity.MfaCredential
		var err error

		mfaCredentials, err = ReadAll(db)
		Expect(err).NotTo(HaveOccurred())

		Expect(mfaCredentials).To(HaveLen(2))
		Expect(mfaCredentials).To(ConsistOf(
			entity.MfaCredential{
				UserId:                  "1",
				MfaProviderId:           entity.Char("mfa_provider_id"),
				ZoneId:                  entity.Char("zone_id"),
				ValidationCode:          sql.NullInt64{Int64: 1234, Valid: true},
				ScratchCodes:            "scratch_codes",
				SecretKey:            	 "secret-key",
				EncryptionKeyLabel:      "activeKeyLabel",
				EncryptedValidationCode: "encrypted_validation_code",
			},
			entity.MfaCredential{
				UserId:                  "2",
				MfaProviderId:           entity.Char("mfa_provider_id"),
				ZoneId:                  entity.Char("zone_id"),
				ValidationCode:          sql.NullInt64{Int64: 1234, Valid: true},
				SecretKey:            	 "secret-key",
				ScratchCodes:            "scratch_codes",
				EncryptionKeyLabel:      "activeKeyLabel",
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
			})

			It("should return a meaningful error", func() {
				_, err := ReadAll(queryer)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("ReadAll failed to query table: cannot query table"))
			})
		})
	})

})
