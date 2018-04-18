package db_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/cloudfoundry/uaa-key-rotator/db"
	"github.com/cloudfoundry/uaa-key-rotator/db/dbfakes"
	"errors"
)

var _ = Describe("Postgresql", func() {

	BeforeEach(func() {
		deleteResult, err := db.Exec(`delete from user_google_mfa_credentials`)
		Expect(err).NotTo(HaveOccurred())
		numOfRowsDeleted, err := deleteResult.RowsAffected()
		Expect(err).NotTo(HaveOccurred())
		Expect(numOfRowsDeleted).To(BeNumerically(">=", int64(0)))

		insertResult, err := db.Exec(`insert into user_google_mfa_credentials(
		user_id, 
		secret_key, 
		validation_code, 
		scratch_codes, 
		mfa_provider_id, 
		zone_id, 
		encryption_key_label, 
		encrypted_validation_code) values(
		'1', 'secret-key', 1234, 'scratch_codes', 'mfa_provider_id', 'zone_id', 'activeKeyLabel', 'encrypted_validation_code'
		)`)
		Expect(err).NotTo(HaveOccurred())
		numOfRowsInserted, err := insertResult.RowsAffected()
		Expect(err).NotTo(HaveOccurred())
		Expect(numOfRowsInserted).To(Equal(int64(1)))
	})

	It("should return every record from the user_google_mfa_credentials table", func() {
		var mfaCredentials []MfaCredentials
		var err error

		mfaCredentials, err = ReadAll(db)
		Expect(err).NotTo(HaveOccurred())

		Expect(mfaCredentials).To(ContainElement(MfaCredentials{
			UserId:                  "1",
			MfaProviderId:           Char("mfa_provider_id"),
			ZoneId:                  Char("zone_id"),
			ValidationCode:          1234,
			ScratchCodes:            "scratch_codes",
			EncryptionKeyLabel:      "activeKeyLabel",
			EncryptedValidationCode: "encrypted_validation_code",
		}))
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
