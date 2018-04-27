package db_test

import (
	"errors"
	db2 "github.com/cloudfoundry/uaa-key-rotator/db"
	"github.com/cloudfoundry/uaa-key-rotator/db/dbfakes"
	"github.com/cloudfoundry/uaa-key-rotator/entity"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"
	"time"
	"github.com/cloudfoundry/uaa-key-rotator/db/testutils"
)

var _ = Describe("Writer", func() {

	var (
		defaultMfaCredential entity.MfaCredential
	)

	BeforeEach(func() {
		deleteResult, err := db.Exec(`delete from user_google_mfa_credentials`)
		Expect(err).NotTo(HaveOccurred())
		numOfRowsDeleted, err := deleteResult.RowsAffected()
		Expect(err).NotTo(HaveOccurred())
		Expect(numOfRowsDeleted).To(BeNumerically(">=", int64(0)))

		newUserID := getRandomTimestamp()
		defaultMfaCredential = insertGoogleMfaCredential(newUserID)
	})

	It("should update a single mfa record", func() {
		var err error

		updatedMfaCredential := entity.MfaCredential{
			UserId:                  defaultMfaCredential.UserId,
			MfaProviderId:           defaultMfaCredential.MfaProviderId,
			ZoneId:                  defaultMfaCredential.ZoneId,
			ValidationCode:          defaultMfaCredential.ValidationCode,
			ScratchCodes:            getRandomTimestamp(),
			SecretKey:               getRandomTimestamp(),
			EncryptionKeyLabel:      getRandomTimestamp(),
			EncryptedValidationCode: getRandomTimestamp(),
		}
		mfaCredentialId3 := insertGoogleMfaCredential("userid_3")

		err = db2.Write(db2.DbAwareQuerier{DB: db, DBScheme: testutils.Scheme}, updatedMfaCredential)
		Expect(err).NotTo(HaveOccurred())

		mfaCredentials, err := db2.ReadAll(db)
		Expect(err).NotTo(HaveOccurred())
		Expect(mfaCredentials).To(ConsistOf(updatedMfaCredential, mfaCredentialId3))
	})

	Describe("when db error occurs", func() {
		var mockDb *dbfakes.FakeQueryer
		BeforeEach(func() {
			mockDb = &dbfakes.FakeQueryer{}
			mockDb.QueryxReturns(nil, errors.New("some db error"))
		})
		It("should return meaningful error", func() {
			err := db2.Write(mockDb, entity.MfaCredential{})
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("Unable to update mfa db record: some db error"))
		})
	})
})

func getRandomTimestamp() string {
	return strconv.Itoa(int(time.Now().UnixNano()))
}
